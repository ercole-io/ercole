// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package service is a package that provides methods for querying data
package job

import (
	"context"
	"fmt"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"github.com/oracle/oci-go-sdk/database"
	"github.com/oracle/oci-go-sdk/v45/common"
)

type OciDatabase struct {
	CompartmentID   string
	CompartmentName string
	DBSystemID      string
	HomeID          string
	Name            string
	UniqueName      string
	ProfileID       string
	Hostname        []string
}

type OciCompartmnentAndDB struct {
	CompartmentID   string
	CompartmentName string
	ProfileID       string
	Databases       []string
}

type HostnameAndStatus struct {
	hostname string
	status   string
	nodeId   string
}

type DBWorks struct {
	hostname   string
	uniqueName string
	cpuThreads int
	archived   bool
	profileID  string
	work       []int
}

func (job *OciDataRetrieveJob) GetOciSISRightsizing(profiles []string, seqValue uint64) {
	var listRec []model.OciRecommendation
	listRec = make([]model.OciRecommendation, 0)
	errors := make([]model.OciRecommendationError, 0)

	var ore model.OciRecommendationError

	var listCompartments []model.OciCompartment

	dbList := make(map[string]OciDatabase)

	var recommendation model.OciRecommendation

	ercoleDatabases, err := job.Database.GetErcoleDatabases()
	if err != nil {
		recError := ore.SetOciRecommendationError(seqValue, "", model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
		errors = append(errors, recError)
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}

		return
	}

	ercoleActiveDatabases, err := job.Database.GetErcoleActiveDatabases()
	if err != nil {
		recError := ore.SetOciRecommendationError(seqValue, "", model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
		errors = append(errors, recError)
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}

		return
	}

	var reorderedDBList map[string]OciCompartmnentAndDB

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		dbClient, err := database.NewDatabaseClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		// retrieve metrics data for each compartment
		for _, compartment := range listCompartments {
			// First retrieve the list of homeID
			req := database.ListDbHomesRequest{
				CompartmentId: &compartment.CompartmentID,
			}
			resp, err := dbClient.ListDbHomes(context.Background(), req)

			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

				continue
			}

			var dbRefId string

			var dbIdType string

			if len(resp.Items) != 0 {
				for _, dbHome := range resp.Items {
					if dbHome.VmClusterId != nil {
						dbRefId = *dbHome.VmClusterId
						dbIdType = "CLUSTER"
					}

					if dbHome.DbSystemId != nil {
						dbRefId = *dbHome.DbSystemId
						dbIdType = "SYSTEM_ID"
					}

					var dbTmp OciDatabase
					dbTmp.HomeID = *dbHome.Id
					dbTmp.CompartmentID = compartment.CompartmentID
					dbTmp.CompartmentName = compartment.Name
					dbTmp.ProfileID = profileId
					dbVal, err1 := job.getDatabaseName(dbClient, compartment.CompartmentID, *dbHome.Id)

					dbTmp.DBSystemID = dbVal.DBSystemID
					dbTmp.HomeID = dbVal.HomeID
					dbTmp.Name = dbVal.Name
					dbTmp.UniqueName = dbVal.UniqueName

					dbList[dbRefId] = dbTmp

					if err1 != nil {
						recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
						errors = append(errors, recError)

						continue
					}

					hostnamesAndStatus, err2 := job.getHostamesAndStatus(dbClient, compartment.CompartmentID, dbVal.DBSystemID, dbIdType)

					if err2 != nil {
						recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
						errors = append(errors, recError)

						continue
					}

					for _, val := range hostnamesAndStatus {
						if val.status == "STOPPED" {
							recommendation.Details = make([]model.RecDetail, 0)
							recommendation.SeqValue = seqValue
							recommendation.ProfileID = profileId
							recommendation.Category = model.OciUnusedServiceDecommisioning //TYPE 3
							recommendation.Suggestion = model.OciDeleteDatabaseInstanceNotActive
							recommendation.CompartmentID = compartment.CompartmentID
							recommendation.CompartmentName = compartment.Name
							recommendation.ResourceID = val.nodeId
							recommendation.Name = val.hostname
							recommendation.ObjectType = model.OciObjectTypeDatabase
							detail1 := model.RecDetail{Name: "Hostname", Value: val.hostname}
							detail2 := model.RecDetail{Name: "CPU Core Count", Value: ""}

							recommendation.Details = append(recommendation.Details, detail1, detail2)
							recommendation.CreatedAt = time.Now().UTC()

							listRec = append(listRec, recommendation)
						} else {
							var dbTmp OciDatabase
							var hosts []string
							dbTmp = dbList[dbRefId]
							hosts = dbTmp.Hostname

							hosts = append(hosts, val.hostname)
							dbTmp.Hostname = hosts
							dbList[dbRefId] = dbTmp
						}
					}
				}
			}
		}
	}

	reorderedDBList = job.getOciReorderedDBList(dbList)
	listRec = job.verifyErcoleAndOciDatabasesConfiguration(ercoleActiveDatabases, reorderedDBList, listRec, seqValue)
	listRec = job.manageErcoleDatabases(ercoleDatabases, reorderedDBList, listRec, seqValue)

	if len(listRec) > 0 {
		errDb := job.Database.AddOciRecommendations(listRec)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}

	if len(errors) > 0 {
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}
}

func (job *OciDataRetrieveJob) getDatabaseName(dbClient database.DatabaseClient, compartmentId string, homeId string) (OciDatabase, error) {
	var retDB OciDatabase

	req := database.ListDatabasesRequest{
		CompartmentId: common.String(compartmentId),
		DbHomeId:      common.String(homeId),
	}

	// Send the request using the service client
	resp, err := dbClient.ListDatabases(context.Background(), req)
	if err != nil {
		return retDB, err
	}

	var DbRefId string

	for _, dbVal := range resp.Items {
		if dbVal.VmClusterId != nil {
			DbRefId = *dbVal.VmClusterId
		}

		if dbVal.DbSystemId != nil {
			DbRefId = *dbVal.DbSystemId
		}

		retDB.DBSystemID = DbRefId
		retDB.HomeID = homeId
		retDB.Name = *dbVal.DbName
		retDB.UniqueName = *dbVal.DbUniqueName
	}

	return retDB, nil
}

func (job *OciDataRetrieveJob) getHostamesAndStatus(dbClient database.DatabaseClient, compartmentId string, dbRefId string, dbIdType string) ([]HostnameAndStatus, error) {
	var hostnamesAndStatus []HostnameAndStatus

	var tmpHostAndSt HostnameAndStatus

	var req database.ListDbNodesRequest

	if dbIdType == "SYSTEM_ID" {
		req = database.ListDbNodesRequest{
			DbSystemId:    common.String(dbRefId),
			CompartmentId: common.String(compartmentId),
		}
	} else {
		req = database.ListDbNodesRequest{
			VmClusterId:   common.String(dbRefId),
			CompartmentId: common.String(compartmentId),
		}
	}

	// Send the request using the service client
	resp, err := dbClient.ListDbNodes(context.Background(), req)
	if err != nil {
		return nil, err
	}

	for _, node := range resp.Items {
		tmpHostAndSt.hostname = *node.Hostname
		tmpHostAndSt.status = fmt.Sprintf("%v", node.LifecycleState)
		tmpHostAndSt.nodeId = *node.Id
		hostnamesAndStatus = append(hostnamesAndStatus, tmpHostAndSt)
	}

	return hostnamesAndStatus, nil
}

func (job *OciDataRetrieveJob) manageErcoleDatabases(ercoleDatabases []dto.ErcoleDatabase, reorderedDBList map[string]OciCompartmnentAndDB, listRec []model.OciRecommendation, seqValue uint64) []model.OciRecommendation { //([]string, error)
	var recommendation model.OciRecommendation

	eDBWorkList := make(map[string]DBWorks)

	for _, eDB := range ercoleDatabases {
		var wkList []int

		var wkDB DBWorks

		for _, sDB := range eDB.Features.Oracle.Database.Databases {
			if _, ok := reorderedDBList[eDB.Hostname]; ok {
				key := eDB.Hostname + "-" + sDB.UniqueName
				wkDB = eDBWorkList[key]
				wkDB.hostname = eDB.Hostname
				wkDB.uniqueName = sDB.UniqueName
				wkDB.cpuThreads = eDB.Info.CpuThreads
				wkDB.profileID = reorderedDBList[eDB.Hostname].ProfileID

				if !eDB.Archived {
					wkDB.archived = false
				}

				wkList = wkDB.work
				wkList = append(wkList, sDB.Work)
				wkDB.work = wkList
				eDBWorkList[key] = wkDB
			}
		}
	}

	for _, dbWork := range eDBWorkList {
		cnt := 0
		opt := true

		for _, wk := range dbWork.work {
			if wk == 0 {
				cnt += 1
			}

			if wk > dbWork.cpuThreads/2 {
				opt = false
			}
		}

		if cnt > 5 || opt {
			recommendation.Details = make([]model.RecDetail, 0)
			recommendation.SeqValue = seqValue
			recommendation.ProfileID = dbWork.profileID
			recommendation.Category = model.OciSISRightsizing
			recommendation.Suggestion = model.OciResizeOversizedDatabaseInstance
			recommendation.Name = dbWork.hostname + "-" + dbWork.uniqueName
			recommendation.ResourceID = ""
			recommendation.ObjectType = model.OciObjectTypeDatabase
			detail1 := model.RecDetail{Name: "Instance Name", Value: recommendation.Name}
			detail2 := model.RecDetail{Name: "Ercole Installed", Value: "YES"}
			detail4 := model.RecDetail{Name: "Ercole Host Cpu Thread", Value: fmt.Sprintf("%d", dbWork.cpuThreads)}

			var detail3 model.RecDetail
			if cnt > 5 {
				detail3 = model.RecDetail{Name: "AWR Enabled", Value: "NO"}
			} else {
				detail3 = model.RecDetail{Name: "AWR Enabled", Value: "YES"}
			}

			recommendation.Details = append(recommendation.Details, detail1, detail2, detail3, detail4)
			recommendation.CreatedAt = time.Now().UTC()

			listRec = append(listRec, recommendation)
		}
	}

	return listRec
}

func (job *OciDataRetrieveJob) verifyErcoleAndOciDatabasesConfiguration(ercoleDatabases []dto.ErcoleDatabase, reorderedDBList map[string]OciCompartmnentAndDB, listRec []model.OciRecommendation, seqValue uint64) []model.OciRecommendation {
	var listDBTmp []string

	var recommendation model.OciRecommendation

	dbNotFound := make(map[string][]string)

	for k, v := range reorderedDBList {
		findHostname := false

		for _, eDBlist := range ercoleDatabases {
			if eDBlist.Hostname == k {
				findHostname = true

				for _, fList := range eDBlist.Features.Oracle.Database.Databases {
					findDB := false

					for _, vv := range v.Databases {
						if fList.UniqueName == vv {
							findDB = true
						}
					}

					if !findDB {
						listDBTmp = dbNotFound[eDBlist.Hostname]
						listDBTmp = append(listDBTmp, fList.UniqueName)
						dbNotFound[eDBlist.Hostname] = listDBTmp
						recommendation.Details = make([]model.RecDetail, 0)
						recommendation.SeqValue = seqValue
						recommendation.ProfileID = v.ProfileID
						recommendation.Category = model.OciSISRightsizing
						recommendation.Suggestion = model.OciResizeOversizedDatabaseInstance
						recommendation.CompartmentID = v.CompartmentID
						recommendation.CompartmentName = v.CompartmentName
						recommendation.Name = eDBlist.Hostname + "-" + fList.UniqueName
						recommendation.ResourceID = ""
						recommendation.ObjectType = model.OciObjectTypeDatabase
						detail1 := model.RecDetail{Name: "Instance Name", Value: recommendation.Name}
						detail2 := model.RecDetail{Name: "Ercole Installed", Value: "NO"}
						detail3 := model.RecDetail{Name: "AWR Enabled", Value: "NO"}

						recommendation.Details = append(recommendation.Details, detail1, detail2, detail3)
						recommendation.CreatedAt = time.Now().UTC()

						listRec = append(listRec, recommendation)
					}
				}
			}
		}

		if !findHostname {
			listDBTmp = dbNotFound[k]
			listDBTmp = append(listDBTmp, "placeholder")
			dbNotFound[k] = listDBTmp
			recommendation.Details = make([]model.RecDetail, 0)
			recommendation.SeqValue = seqValue
			recommendation.ProfileID = v.ProfileID
			recommendation.Category = model.OciSISRightsizing
			recommendation.Suggestion = model.OciResizeOversizedDatabaseInstance
			recommendation.CompartmentID = v.CompartmentID
			recommendation.CompartmentName = v.CompartmentName
			recommendation.Name = k
			recommendation.ResourceID = ""
			recommendation.ObjectType = model.OciObjectTypeDatabase
			detail1 := model.RecDetail{Name: "Instance Name", Value: recommendation.Name}
			detail2 := model.RecDetail{Name: "Ercole Installed", Value: "NO"}
			detail3 := model.RecDetail{Name: "AWR Enabled", Value: "NO"}

			recommendation.Details = append(recommendation.Details, detail1, detail2, detail3)
			recommendation.CreatedAt = time.Now().UTC()

			listRec = append(listRec, recommendation)
		}
	}

	return listRec
}

func (job *OciDataRetrieveJob) getOciReorderedDBList(ociDBList map[string]OciDatabase) map[string]OciCompartmnentAndDB {
	var listDBTmp []string

	var ociCompDBTmp OciCompartmnentAndDB

	reorderedDBList := make(map[string]OciCompartmnentAndDB)

	for _, db := range ociDBList {
		for _, host := range db.Hostname {
			ociCompDBTmp = reorderedDBList[host]
			ociCompDBTmp.CompartmentID = db.CompartmentID
			ociCompDBTmp.CompartmentName = db.CompartmentName
			ociCompDBTmp.ProfileID = db.ProfileID
			listDBTmp = ociCompDBTmp.Databases
			listDBTmp = append(listDBTmp, db.UniqueName)
			ociCompDBTmp.Databases = listDBTmp
			reorderedDBList[host] = ociCompDBTmp
		}
	}

	return reorderedDBList
}
