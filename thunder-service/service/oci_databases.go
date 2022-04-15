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
package service

import (
	"context"
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	multierror "github.com/hashicorp/go-multierror"
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
	Hostname        []string
}

type OciCompartmnentAndDB struct {
	CompartmentID   string
	CompartmentName string
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
	work       []int
}

func (as *ThunderService) GetOciSISRightsizing(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var listRec []model.OciErcoleRecommendation
	listRec = make([]model.OciErcoleRecommendation, 0)

	var merr error

	var listCompartments []model.OciCompartment

	dbList := make(map[string]OciDatabase)

	var recommendation model.OciErcoleRecommendation

	ercoleDatabases, err := as.Database.GetErcoleDatabases()
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	ercoleActiveDatabases, err := as.Database.GetErcoleActiveDatabases()
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	var reorderedDBList map[string]OciCompartmnentAndDB

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := as.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		listCompartments, err = as.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		dbClient, err := database.NewDatabaseClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
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
				merr = multierror.Append(merr, err)
				continue
			}

			if len(resp.Items) != 0 {
				for _, dbHome := range resp.Items {
					var dbTmp OciDatabase
					dbTmp.HomeID = *dbHome.Id
					dbTmp.CompartmentID = compartment.CompartmentID
					dbTmp.CompartmentName = compartment.Name
					dbList[*dbHome.DbSystemId] = dbTmp
					dbVal, err1 := getDatabaseName(dbClient, compartment.CompartmentID, *dbHome.Id)

					if err1 != nil {
						merr = multierror.Append(merr, err)
						continue
					}

					hostnamesAndStatus, err2 := getHostamesAndStatus(dbClient, compartment.CompartmentID, dbVal.DBSystemID)

					if err2 != nil {
						merr = multierror.Append(merr, err)
						continue
					}

					for _, val := range hostnamesAndStatus {
						if val.status == "STOPPED" {
							recommendation.Details = make([]model.RecDetail, 0)
							recommendation.Category = model.RecommendationTypeUnusedServiceDecommisioning //TYPE 3
							recommendation.Suggestion = model.DeleteDatabaseInstanceNotActive
							recommendation.CompartmentID = compartment.CompartmentID
							recommendation.CompartmentName = compartment.Name
							recommendation.ResourceID = val.nodeId
							recommendation.Name = val.hostname
							recommendation.ObjectType = model.ObjectTypeDatabase
							detail1 := model.RecDetail{Name: "Hostname", Value: val.hostname}
							detail2 := model.RecDetail{Name: "CPU Core Count", Value: ""}

							recommendation.Details = append(recommendation.Details, detail1, detail2)

							listRec = append(listRec, recommendation)
						} else {
							dbVal.Hostname = append(dbVal.Hostname, val.hostname)
							dbList[*dbHome.DbSystemId] = dbVal
						}
					}
				}
			}
		}
	}

	reorderedDBList = getOciReorderedDBList(dbList)
	listRec = verifyErcoleAndOciDatabasesConfiguration(ercoleActiveDatabases, reorderedDBList, listRec)
	listRec = manageErcoleDatabases(ercoleDatabases, reorderedDBList, listRec)

	return listRec, nil
}

func getDatabaseName(dbClient database.DatabaseClient, compartmentId string, homeId string) (OciDatabase, error) {
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

	for _, dbVal := range resp.Items {
		retDB.DBSystemID = *dbVal.DbSystemId
		retDB.HomeID = homeId
		retDB.Name = *dbVal.DbName
		retDB.UniqueName = *dbVal.DbUniqueName
	}

	return retDB, nil
}

func getHostamesAndStatus(dbClient database.DatabaseClient, compartmentId string, dbSystemId string) ([]HostnameAndStatus, error) {
	var hostnamesAndStatus []HostnameAndStatus

	var tmpHostAndSt HostnameAndStatus

	req := database.ListDbNodesRequest{
		DbSystemId:    common.String(dbSystemId),
		CompartmentId: common.String(compartmentId),
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

func manageErcoleDatabases(ercoleDatabases []dto.ErcoleDatabase, reorderedDBList map[string]OciCompartmnentAndDB, listRec []model.OciErcoleRecommendation) []model.OciErcoleRecommendation { //([]string, error)
	var recommendation model.OciErcoleRecommendation

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
			recommendation.Category = model.RecommendationTypeSISRightsizing
			recommendation.Suggestion = model.ResizeOversizedDatabaseInstance
			recommendation.Name = dbWork.hostname + "-" + dbWork.uniqueName
			recommendation.ResourceID = ""
			recommendation.ObjectType = model.ObjectTypeDatabase
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

			listRec = append(listRec, recommendation)
		}
	}

	return listRec
}

func verifyErcoleAndOciDatabasesConfiguration(ercoleDatabases []dto.ErcoleDatabase, reorderedDBList map[string]OciCompartmnentAndDB, listRec []model.OciErcoleRecommendation) []model.OciErcoleRecommendation {
	var listDBTmp []string

	var recommendation model.OciErcoleRecommendation

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
						recommendation.Category = model.RecommendationTypeSISRightsizing
						recommendation.Suggestion = model.ResizeOversizedDatabaseInstance
						recommendation.CompartmentID = v.CompartmentID
						recommendation.CompartmentName = v.CompartmentName
						recommendation.Name = eDBlist.Hostname + "-" + fList.UniqueName
						recommendation.ResourceID = ""
						recommendation.ObjectType = model.ObjectTypeDatabase
						detail1 := model.RecDetail{Name: "Instance Name", Value: recommendation.Name}
						detail2 := model.RecDetail{Name: "Ercole Installed", Value: "NO"}
						detail3 := model.RecDetail{Name: "AWR Enabled", Value: "NO"}

						recommendation.Details = append(recommendation.Details, detail1, detail2, detail3)

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
			recommendation.Category = model.RecommendationTypeSISRightsizing
			recommendation.Suggestion = model.ResizeOversizedDatabaseInstance
			recommendation.CompartmentID = v.CompartmentID
			recommendation.CompartmentName = v.CompartmentName
			recommendation.Name = k
			recommendation.ResourceID = ""
			recommendation.ObjectType = model.ObjectTypeDatabase
			detail1 := model.RecDetail{Name: "Instance Name", Value: recommendation.Name}
			detail2 := model.RecDetail{Name: "Ercole Installed", Value: "NO"}
			detail3 := model.RecDetail{Name: "AWR Enabled", Value: "NO"}

			recommendation.Details = append(recommendation.Details, detail1, detail2, detail3)

			listRec = append(listRec, recommendation)
		}
	}

	return listRec
}

func getOciReorderedDBList(ociDBList map[string]OciDatabase) map[string]OciCompartmnentAndDB {
	var listDBTmp []string

	var ociCompDBTmp OciCompartmnentAndDB

	reorderedDBList := make(map[string]OciCompartmnentAndDB)

	for _, db := range ociDBList {
		for _, host := range db.Hostname {
			ociCompDBTmp = reorderedDBList[host]
			ociCompDBTmp.CompartmentID = db.CompartmentID
			ociCompDBTmp.CompartmentName = db.CompartmentName
			listDBTmp = ociCompDBTmp.Databases
			listDBTmp = append(listDBTmp, db.UniqueName)
			ociCompDBTmp.Databases = listDBTmp
			reorderedDBList[host] = ociCompDBTmp
		}
	}

	return reorderedDBList
}
