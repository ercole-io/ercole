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

package service

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"sort"
	"strings"

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoadOracleDatabaseAgreementPartsList loads the list of Oracle/Database agreement parts and store it to as.OracleDatabaseAgreementParts.
func (as *APIService) LoadOracleDatabaseAgreementPartsList() {
	// read the list content
	listContentRaw, err := ioutil.ReadFile(as.Config.ResourceFilePath + "/oracle_database_agreement_parts_list.json")
	if err != nil {
		as.Log.Warnf("Unable to read %s: %v\n", as.Config.ResourceFilePath+"/oracle_database_agreement_parts_list.json", err)
		return
	}

	// unmarshal to OracleDatabaseAgreementParts
	err = json.Unmarshal(listContentRaw, &as.OracleDatabaseAgreementParts)
	if err != nil {
		as.Log.Warnf("Unable to unmarshal %s: %v\n", as.Config.ResourceFilePath+"/oracle_database_agreement_parts_list.json", err)
		return
	}
}

// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
func (as *APIService) GetOracleDatabaseAgreementPartsList() ([]model.OracleDatabaseAgreementPart, utils.AdvancedErrorInterface) {
	return as.OracleDatabaseAgreementParts, nil
}

// AddOracleDatabaseAgreements return the list of Oracle/Database agreement parts
func (as *APIService) AddOracleDatabaseAgreements(req apimodel.OracleDatabaseAgreementsAddRequest) (interface{}, utils.AdvancedErrorInterface) {

	var parts []*model.OracleDatabaseAgreementPart
	var err utils.AdvancedErrorInterface
	if parts, err = resolvePartIds(req.PartsID, as.OracleDatabaseAgreementParts); err != nil {
		return nil, err
	}

	if err := checkHosts(as, req.Hosts); err != nil {
		return nil, err
	}

	agreements := make([]model.OracleDatabaseAgreement, len(req.PartsID))
	for i, part := range parts {
		agreements[i].AgreementID = req.AgreementID
		agreements[i].CSI = req.CSI
		agreements[i].CatchAll = req.CatchAll
		agreements[i].Count = req.Count
		agreements[i].Hosts = req.Hosts
		agreements[i].ID = primitive.NewObjectIDFromTimestamp(as.TimeNow())
		agreements[i].ReferenceNumber = req.ReferenceNumber
		agreements[i].Unlimited = req.Unlimited

		agreements[i].PartID = part.PartID
		agreements[i].ItemDescription = part.ItemDescription
		agreements[i].Metrics = part.Metrics
	}

	res := make([]interface{}, len(agreements))
	for i, agr := range agreements {
		var aerr utils.AdvancedErrorInterface
		if res[i], aerr = as.Database.InsertOracleDatabaseAgreement(agr); aerr != nil {
			return nil, aerr
		}
	}

	return res, nil
}

func resolvePartIds(partsID []string, agreementParts []model.OracleDatabaseAgreementPart) ([]*model.OracleDatabaseAgreementPart, utils.AdvancedErrorInterface) {
	var parts []*model.OracleDatabaseAgreementPart = make([]*model.OracleDatabaseAgreementPart, len(partsID))

	for i, pID := range partsID {
		var err utils.AdvancedErrorInterface

		if parts[i], err = isValidPartID(pID, agreementParts); err != nil {
			return nil, utils.AerrOracleDatabaseAgreementInvalidPartID
		}
	}

	return parts, nil
}

func isValidPartID(partID string, agreementParts []model.OracleDatabaseAgreementPart) (*model.OracleDatabaseAgreementPart, utils.AdvancedErrorInterface) {
	for i, part := range agreementParts {
		if partID == part.PartID {
			return &agreementParts[i], nil
		}
	}

	return nil, utils.AerrOracleDatabaseAgreementInvalidPartID
}

func checkHosts(as *APIService, hosts []string) utils.AdvancedErrorInterface {
	notInClusterHosts, aerr := as.SearchHosts("hostnames", "", database.SearchHostsFilters{
		GTECPUCores:    -1,
		LTECPUCores:    -1,
		LTECPUThreads:  -1,
		LTEMemoryTotal: -1,
		GTECPUThreads:  -1,
		GTESwapTotal:   -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
	}, "", false, -1, -1, "", "", utils.MAX_TIME)
	if aerr != nil {
		return aerr
	}

	notInClusterHostnames := make([]string, len(notInClusterHosts))
	for i, h := range notInClusterHosts {
		notInClusterHostnames[i] = h["hostname"].(string)
	}

hosts_loop:
	for _, host := range hosts {
		for _, notInClusterHostname := range notInClusterHostnames {
			if host == notInClusterHostname {
				continue hosts_loop
			}
		}

		return utils.AerrHostNotFound
	}

	return nil
}

// UpdateOracleDatabaseAgreement update an Oracle Database Agreement
func (as *APIService) UpdateOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) utils.AdvancedErrorInterface {
	if _, err := as.Database.FindOracleDatabaseAgreement(agreement.ID); err != nil {
		return err
	}

	var part *model.OracleDatabaseAgreementPart
	var err utils.AdvancedErrorInterface
	if part, err = isValidPartID(agreement.PartID, as.OracleDatabaseAgreementParts); err != nil {
		return err
	}

	agreement.PartID = part.PartID
	agreement.ItemDescription = part.ItemDescription
	agreement.Metrics = part.Metrics

	return as.Database.UpdateOracleDatabaseAgreement(agreement)
}

// SearchOracleDatabaseAgreements search Oracle/Database agreements
func (as *APIService) SearchOracleDatabaseAgreements(search string, filters apimodel.SearchOracleDatabaseAgreementsFilter) ([]apimodel.OracleDatabaseAgreementFE, utils.AdvancedErrorInterface) {
	agrs, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	hosts, err := as.Database.ListHostUsingOracleDatabaseLicenses()
	if err != nil {
		return nil, err
	}

	as.AssignOracleDatabaseAgreementsToHosts(agrs, hosts)

	filteredAgrs := make([]apimodel.OracleDatabaseAgreementFE, 0)
	for _, agr := range agrs {

		if CheckOracleDatabaseAgreementMatchFilter(agr, filters) {
			filteredAgrs = append(filteredAgrs, agr)
		}

	}

	return filteredAgrs, nil
}

// AssignOracleDatabaseAgreementsToHosts assign in-place the agreements to every licensingObjects by modifying them
func (as *APIService) AssignOracleDatabaseAgreementsToHosts(
	agrs []apimodel.OracleDatabaseAgreementFE,
	hosts []apimodel.HostUsingOracleDatabaseLicenses) {

	sortOracleDatabaseAgreements(agrs)
	sortHostsUsingLicenses(hosts)

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Agreements = %s\nHosts= %s\n", utils.ToJSON(agrs), utils.ToJSON(hosts))
	}

	// Build data structure for fast access to the informations
	licensingObjectsMap := buildOracleDatabaseLicensingObjectsMap(hosts)
	partsMap := buildOracleDatabaseAgreementPartMap(as.OracleDatabaseAgreementParts)

	// Assign every agreements to the associated host
	for i := range agrs {
		agr := &agrs[i]
		//sort associated hosts by count, considering that parts may have multiple aliases
		sortAssociatedHostsInOracleDatabaseAgreement(*agr, licensingObjectsMap, partsMap)

		// Debug print
		if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
			as.Log.Debugf("Distributing licenses of agreement #%d to host. Agreement = %s\n", i, utils.ToJSON(agr))
		}

		//distribute licenses for each host
		for j := range agr.Hosts {
			host := &agr.Hosts[j]
			//Assign the
			for _, alias := range partsMap[agr.PartID].Aliases {
				// If we have finished the licenses, break
				if agr.Count <= 0 && !agr.Unlimited {
					break
				}
				// If no host require a license with licenseName == alias, skip
				if _, ok := licensingObjectsMap[alias]; !ok {
					continue
				}
				// If the host don't use the license, skip
				if _, ok := licensingObjectsMap[alias][host.Hostname]; !ok {
					continue
				}
				// If the host don't require the license, skip
				if licensingObjectsMap[alias][host.Hostname].LicenseCount <= 0 {
					continue
				}

				// fill all required license, if the host need
				if agr.Unlimited {
					host.CoveredLicensesCount = licensingObjectsMap[alias][host.Hostname].LicenseCount
					licensingObjectsMap[alias][host.Hostname].LicenseCount = 0

					// Debug print
					if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
						as.Log.Debugf("Distributing (ULA) %f licenses to host %s. aggCount=%f associatedHostCovered=%f hostCount=%f licenseName=%s\n",
							licensingObjectsMap[alias][host.Hostname].LicenseCount,
							host.Hostname,
							agr.Count,
							host.CoveredLicensesCount,
							licensingObjectsMap[alias][host.Hostname].LicenseCount,
							alias)
					}
				} else {
					if agr.Metrics == "Processor Perpetual" || agr.Metrics == "Computer Perpetual" {
						coverableLicenses := math.Min(agr.Count, licensingObjectsMap[alias][host.Hostname].LicenseCount)
						licensingObjectsMap[alias][host.Hostname].LicenseCount -= coverableLicenses
						host.CoveredLicensesCount += coverableLicenses
						agr.Count -= coverableLicenses
						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Processor Perpetual/Computer Perpetual) %f licenses to host %s. aggCount=%f associatedHostCovered=%f hostCount=%f licenseName=%s\n", coverableLicenses, host.Hostname, agr.Count, host.CoveredLicensesCount, licensingObjectsMap[alias][host.Hostname].LicenseCount, alias)
						}
					} else if agr.Metrics == "Named User Plus Perpetual" {
						coverableLicenses := math.Floor(math.Min(agr.Count*25, licensingObjectsMap[alias][host.Hostname].LicenseCount) / 25)
						licensingObjectsMap[alias][host.Hostname].LicenseCount -= coverableLicenses * 25
						host.CoveredLicensesCount += coverableLicenses * 25
						agr.Count -= coverableLicenses

						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Named User Plus Perpetual) %f(user) licenses to host %s. aggCount=%f(user) associatedHostCovered=%f hostCount=%f licenseName=%s\n", coverableLicenses, host.Hostname, agr.Count, host.CoveredLicensesCount, licensingObjectsMap[alias][host.Hostname].LicenseCount, alias)
						}
					}
				}
			}
			// If we have finished the licenses, break
			if agr.Count <= 0 && !agr.Unlimited {
				break
			}
		}
	}

	//Resort licensingObjects
	sortHostsUsingLicenses(hosts)
	licensingObjectsMap = buildOracleDatabaseLicensingObjectsMap(hosts) //the map is rebuilded because the references are updated during the sort

	// Debug print
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Resorted LicensingObjects: %#v\n", hosts)
	}

	//Distribute remaining licenses in catch-all agreement to the licensingObjects
	for i := range hosts {
		obj := &hosts[i]

		//The object is already full covered
		if obj.LicenseCount <= 0 {
			continue
		}

		// //Debug print
		// if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		// 	as.Log.Debugf("Finding valid agreement for licensingObject #%d. obj = %s\n", i, utils.ToJSON(obj))
		// }

		//Find an agreement that can cover the object
		for j := range agrs {
			agr := &agrs[j]

			//non catch-all agreement cannot cover the object
			if !agr.CatchAll {
				continue
			}
			//non catch-all agreement cannot cover the object
			if agr.Count <= 0 && !agr.Unlimited {
				continue
			}

			//Try to fill this obj
			for _, alias := range partsMap[agr.PartID].Aliases {
				// If we have finished the licenses, break
				if agr.Count <= 0 && !agr.Unlimited {
					break
				}
				//Ignore this license because it isn't the right
				if obj.LicenseName != alias {
					continue
				}

				// fill all required license, if the host need
				if agr.Unlimited {
					// Debug print
					if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
						as.Log.Debugf("Distributing (ULA) %f licenses to obj %s. aggCount=%f objCount=0 licenseName=%s\n", obj.LicenseCount, obj.Name, agr.Count, alias)
					}

					obj.LicenseCount = 0
				} else {
					if agr.Metrics == "Processor Perpetual" || agr.Metrics == "Computer Perpetual" {
						coverableLicenses := math.Min(agr.Count, obj.LicenseCount)
						obj.LicenseCount -= coverableLicenses
						agr.Count -= coverableLicenses
						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Processor Perpetual/Computer Perpetual) %f licenses to obj %s. aggCount=%f objCount=%f licenseName=%s\n", coverableLicenses, obj.Name, agr.Count, obj.LicenseCount, alias)
						}
					} else if agr.Metrics == "Named User Plus Perpetual" {
						coverableLicenses := math.Floor(math.Min(agr.Count*25, obj.LicenseCount) / 25)
						obj.LicenseCount -= coverableLicenses * 25
						agr.Count -= coverableLicenses

						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Named User Plus Perpetual) %f(user) licenses to obj %s. aggCount=%f(user) objCount=%f licenseName=%s\n", coverableLicenses, obj.Name, agr.Count, obj.LicenseCount, alias)
						}
					}
				}
			}
		}
	}

	// Debug print
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Associations finished. LicensingObjects: %#v\n", hosts)
	}

	type coverStatus struct {
		Covered                float64 //==purchased
		TotalCoverableLicenses float64 //==consumed
	}

	//Calculate total number of covered/uncovered for each
	allLicensesCoverStatus := make(map[string]coverStatus)
	for _, obj := range hosts {
		allLicensesCoverStatus[obj.LicenseName] = coverStatus{
			TotalCoverableLicenses: allLicensesCoverStatus[obj.LicenseName].TotalCoverableLicenses + obj.OriginalCount,
			Covered:                allLicensesCoverStatus[obj.LicenseName].Covered + (obj.OriginalCount - obj.LicenseCount),
		}
	}

	// Debug print
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Cover status: %#v\n", allLicensesCoverStatus)
	}

	//Calculate TotalCoveredLicenses and available
	for i := range agrs {
		agr := &agrs[i]
		uncoveredLicenseAssociatedHostSum := 0.0
		uncoveredLicenseUnassociatedObjSum := 0.0
		//calculate available
		for _, alias := range partsMap[agr.PartID].Aliases {
			uncoveredLicenseUnassociatedObjSum += allLicensesCoverStatus[alias].TotalCoverableLicenses - allLicensesCoverStatus[alias].Covered
			for j := range agr.Hosts {
				host := &agr.Hosts[j]
				// If no host require a license with licenseName == alias, skip
				if _, ok := licensingObjectsMap[alias]; !ok {
					continue
				}
				// If the host don't use the license, skip
				if _, ok := licensingObjectsMap[alias][host.Hostname]; !ok {
					continue
				}
				host.TotalCoveredLicensesCount = licensingObjectsMap[alias][host.Hostname].OriginalCount - licensingObjectsMap[alias][host.Hostname].LicenseCount
				host.ConsumedLicensesCount = licensingObjectsMap[alias][host.Hostname].OriginalCount
				uncoveredLicenseAssociatedHostSum += licensingObjectsMap[alias][host.Hostname].LicenseCount //non-covered part
			}
		}

		if !agr.CatchAll {
			agr.AvailableCount = -uncoveredLicenseAssociatedHostSum
		} else {
			agr.AvailableCount = -uncoveredLicenseUnassociatedObjSum
		}
	}

}

// CheckOracleDatabaseAgreementMatchFilter check that agr match the filters
func CheckOracleDatabaseAgreementMatchFilter(agr apimodel.OracleDatabaseAgreementFE, filters apimodel.SearchOracleDatabaseAgreementsFilter) bool {
	return strings.Contains(strings.ToLower(agr.AgreementID), strings.ToLower(filters.AgreementID)) &&
		strings.Contains(strings.ToLower(agr.PartID), strings.ToLower(filters.PartID)) &&
		strings.Contains(strings.ToLower(agr.ItemDescription), strings.ToLower(filters.ItemDescription)) &&
		strings.Contains(strings.ToLower(agr.CSI), strings.ToLower(filters.CSI)) &&
		(filters.Metrics == "" || strings.ToLower(agr.Metrics) == strings.ToLower(filters.Metrics)) &&
		strings.Contains(strings.ToLower(agr.ReferenceNumber), strings.ToLower(filters.ReferenceNumber)) &&
		(filters.Unlimited == "NULL" || agr.Unlimited == (filters.Unlimited == "true")) &&
		(filters.CatchAll == "NULL" || agr.CatchAll == (filters.CatchAll == "true")) &&
		(filters.LicensesCountLTE == -1 || agr.LicensesCount <= float64(filters.LicensesCountLTE)) &&
		(filters.LicensesCountGTE == -1 || agr.LicensesCount >= float64(filters.LicensesCountGTE)) &&
		(filters.UsersCountLTE == -1 || agr.UsersCount <= float64(filters.UsersCountLTE)) &&
		(filters.UsersCountGTE == -1 || agr.UsersCount >= float64(filters.UsersCountGTE)) &&
		(filters.AvailableCountLTE == -1 || agr.AvailableCount <= float64(filters.AvailableCountLTE)) &&
		(filters.AvailableCountGTE == -1 || agr.AvailableCount >= float64(filters.AvailableCountGTE))
}

// sortHostsUsingLicenses sort the list of hosts by count
func sortHostsUsingLicenses(obj []apimodel.HostUsingOracleDatabaseLicenses) {
	sort.Slice(obj, func(i, j int) bool {
		if obj[i].LicenseCount != obj[j].LicenseCount {
			return obj[i].LicenseCount > obj[j].LicenseCount
		} else if obj[i].Name != obj[j].Name {
			return obj[i].Name > obj[j].Name
		} else {
			return obj[i].LicenseName > obj[j].LicenseName
		}
	})
}

// sortOracleDatabaseAgreements sort the list of apimodel.OracleDatabaseAgreementsFE
// firstly, which aren't CatchAll
// secondly, which arent't Unlimited
// thirdly, which has more UsersCount
// fourth, which has more LicensesCount
func sortOracleDatabaseAgreements(obj []apimodel.OracleDatabaseAgreementFE) {
	sort.Slice(obj, func(i, j int) bool {

		if obj[i].CatchAll != obj[j].CatchAll {
			return obj[j].CatchAll

		} else if obj[i].Unlimited != obj[j].Unlimited {
			return obj[j].Unlimited

		} else if obj[i].UsersCount != obj[j].UsersCount {
			return obj[i].UsersCount > obj[j].UsersCount

		} else {
			return obj[i].LicensesCount > obj[j].LicensesCount
		}
	})
}

// sortAssociatedHostsInOracleDatabaseAgreement sort the associated hosts by license count. It  that parts may have multiple aliases
func sortAssociatedHostsInOracleDatabaseAgreement(agr apimodel.OracleDatabaseAgreementFE, licensingObjectsMap map[string]map[string]*apimodel.HostUsingOracleDatabaseLicenses, partsMap map[string]*model.OracleDatabaseAgreementPart) {
	sort.Slice(agr.Hosts, func(i, j int) bool {
		var maxLicensingObjectICount float64 = 0
		var maxLicensingObjectJCount float64 = 0
		for _, alias := range partsMap[agr.PartID].Aliases {
			if _, ok := licensingObjectsMap[alias]; ok {
				if _, ok := licensingObjectsMap[alias][agr.Hosts[i].Hostname]; ok {
					maxLicensingObjectICount = math.Max(maxLicensingObjectICount, licensingObjectsMap[alias][agr.Hosts[i].Hostname].LicenseCount)
				}
				if _, ok := licensingObjectsMap[alias][agr.Hosts[j].Hostname]; ok {
					maxLicensingObjectJCount = math.Max(maxLicensingObjectJCount, licensingObjectsMap[alias][agr.Hosts[j].Hostname].LicenseCount)
				}
			}
		}
		return maxLicensingObjectICount > maxLicensingObjectJCount
	})
}

// buildOracleDatabaseLicensingObjectsMap return a map of license name to map of object name to pointer to  apimodel.HostUsingOracleDatabaseLicenses for fast object lookup
// Assume that doesn't exist a cluster and a host with the same name
func buildOracleDatabaseLicensingObjectsMap(objs []apimodel.HostUsingOracleDatabaseLicenses) map[string]map[string]*apimodel.HostUsingOracleDatabaseLicenses {
	res := make(map[string]map[string]*apimodel.HostUsingOracleDatabaseLicenses)

	for i, obj := range objs {
		if _, ok := res[obj.LicenseName]; !ok {
			res[obj.LicenseName] = make(map[string]*apimodel.HostUsingOracleDatabaseLicenses)
		}
		res[obj.LicenseName][obj.Name] = &objs[i]
	}

	return res
}

// buildOracleDatabaseAgreementPartMap return a map of partID to OracleDatabaseAgreementPart
func buildOracleDatabaseAgreementPartMap(parts []model.OracleDatabaseAgreementPart) map[string]*model.OracleDatabaseAgreementPart {
	partsMap := make(map[string]*model.OracleDatabaseAgreementPart)

	for i, part := range parts {
		partsMap[part.PartID] = &parts[i]
	}

	return partsMap
}

// AddAssociatedHostToOracleDatabaseAgreement a new host to the list of associated hosts of the agreement
func (as *APIService) AddAssociatedHostToOracleDatabaseAgreement(id primitive.ObjectID, hostname string) utils.AdvancedErrorInterface {
	var err utils.AdvancedErrorInterface

	//check the existence of the host
	if exist, err := as.Database.ExistNotInClusterHost(hostname); err != nil {
		return err
	} else if !exist {
		return utils.AerrNotInClusterHostNotFound
	}

	//check the existence and get the agreement
	var agr model.OracleDatabaseAgreement
	if agr, err = as.Database.FindOracleDatabaseAgreement(id); err != nil {
		return err
	}

	//check the host isn't already part of the list, and do nothing
	for _, host := range agr.Hosts {
		if host == hostname {
			return nil
		}
	}

	//add the host to the list
	agr.Hosts = append(agr.Hosts, hostname)

	//save the host in the database
	return as.Database.UpdateOracleDatabaseAgreement(agr)
}

// RemoveAssociatedHostToOracleDatabaseAgreement remove the host from the list of associated hosts of the agreement
func (as *APIService) RemoveAssociatedHostToOracleDatabaseAgreement(id primitive.ObjectID, hostname string) utils.AdvancedErrorInterface {
	var err utils.AdvancedErrorInterface

	var agr model.OracleDatabaseAgreement
	if agr, err = as.Database.FindOracleDatabaseAgreement(id); err != nil {
		return err
	}

	for i, host := range agr.Hosts {
		if host == hostname {
			agr.Hosts = append(agr.Hosts[:i], agr.Hosts[i+1:]...)

			return as.Database.UpdateOracleDatabaseAgreement(agr)
		}
	}

	return nil
}

// DeleteOracleDatabaseAgreement remove an Oracle/Database agreement
func (as *APIService) DeleteOracleDatabaseAgreement(id primitive.ObjectID) utils.AdvancedErrorInterface {
	return as.Database.RemoveOracleDatabaseAgreement(id)
}
