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
	//Check and resolve every part id
	var parts []*model.OracleDatabaseAgreementPart = make([]*model.OracleDatabaseAgreementPart, len(req.PartsID))
	for i, pid := range req.PartsID {
		found := false
		for j, vpid := range as.OracleDatabaseAgreementParts {
			if pid == vpid.PartID {
				found = true
				parts[i] = &as.OracleDatabaseAgreementParts[j]
				break
			}
		}
		if !found {
			return nil, utils.AerrOracleDatabaseAgreementInvalidPartID
		}
	}

	//Get the list of hosts not in cluster
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
		return nil, aerr
	}
	notInClusterHostnames := make([]string, len(notInClusterHosts))
	for i, h := range notInClusterHosts {
		notInClusterHostnames[i] = h["hostname"].(string)
	}

	//Check every host in req.Hosts
	for _, host := range req.Hosts {
		found := false
		for _, vhost := range notInClusterHostnames {
			if host == vhost {
				found = true
				break
			}
		}
		if !found {
			return nil, utils.AerrHostNotFound
		}
	}

	//expode req in multple agreement
	aggs := make([]model.OracleDatabaseAgreement, len(req.PartsID))
	for i, part := range parts {
		aggs[i].AgreementID = req.AgreementID
		aggs[i].CSI = req.CSI
		aggs[i].CatchAll = req.CatchAll
		aggs[i].Count = req.Count
		aggs[i].Hosts = req.Hosts
		aggs[i].ID = primitive.NewObjectIDFromTimestamp(as.TimeNow())
		aggs[i].ItemDescription = part.ItemDescription
		aggs[i].Metrics = part.Metrics
		aggs[i].PartID = part.PartID
		aggs[i].ReferenceNumber = req.ReferenceNumber
		aggs[i].Unlimited = req.Unlimited
	}

	//insert it to the database
	res := make([]interface{}, len(aggs))
	for i, agg := range aggs {
		if res[i], aerr = as.Database.InsertOracleDatabaseAgreement(agg); aerr != nil {
			return nil, aerr
		}
	}

	return res, nil
}

// SearchOracleDatabaseAgreements search Oracle/Database agreements
func (as *APIService) SearchOracleDatabaseAgreements(search string, filters apimodel.SearchOracleDatabaseAgreementsFilters) ([]apimodel.OracleDatabaseAgreementsFE, utils.AdvancedErrorInterface) {
	//Get the list of aggreements
	aggs, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	//Get the list of licensingObjecst
	objs, err := as.Database.ListOracleDatabaseLicensingObjects()
	if err != nil {
		return nil, err
	}

	//Compute the algorithm
	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(aggs, objs)

	//Filter them!
	filteredAggs := make([]apimodel.OracleDatabaseAgreementsFE, 0)
	for _, agg := range aggs {
		if !CheckOracleDatabaseAgreementMatchFilter(agg, filters) {
			continue
		}

		filteredAggs = append(filteredAggs, agg)
	}

	return filteredAggs, nil
}

// GreedilyAssignOracleDatabaseAgreementsToLicensingObjects assign in-place the agreements greedly to every licensingObjects by modifying them
func (as *APIService) GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(aggs []apimodel.OracleDatabaseAgreementsFE, licensingObjects []apimodel.OracleDatabaseLicensingObjects) {
	//TODO: optimize this algorithm!

	// Sort the arrays for optimizing the greedy take of the right object
	SortOracleDatabaseAgreements(aggs)
	SortOracleDatabaseAgreementLicensingObjects(licensingObjects)

	// Debug print
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Agreements = %s\nLicensingObjects = %s\n", utils.ToJSON(aggs), utils.ToJSON(licensingObjects))
	}

	// Build data structure for fast access to the informations
	licensingObjectsMap := BuildOracleDatabaseLicensingObjectsMap(licensingObjects)
	partsMap := BuildOracleDatabaseAgreementPartMap(as.OracleDatabaseAgreementParts)

	// Assign every agreements to the associated host
	for i := range aggs {
		agg := &aggs[i]
		//sort associated hosts by count, considering that parts may have multiple aliases
		SortAssociatedHostsInOracleDatabaseAgreement(*agg, licensingObjectsMap, partsMap)

		// Debug print
		if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
			as.Log.Debugf("Distributing licenses of agreement #%d to host. Agreement = %s\n", i, utils.ToJSON(agg))
		}

		//distribute licenses for each host
		for j := range agg.Hosts {
			host := &agg.Hosts[j]
			//Assign the
			for _, alias := range partsMap[agg.PartID].Aliases {
				// If we have finished the licenses, break
				if agg.Count <= 0 && !agg.Unlimited {
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
				if licensingObjectsMap[alias][host.Hostname].Count <= 0 {
					continue
				}

				// fill all required license, if the host need
				if agg.Unlimited {
					host.CoveredLicensesCount = licensingObjectsMap[alias][host.Hostname].Count
					licensingObjectsMap[alias][host.Hostname].Count = 0

					// Debug print
					if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
						as.Log.Debugf("Distributing (ULA) %f licenses to host %s. aggCount=%f associatedHostCovered=%f hostCount=%f licenseName=%s\n", licensingObjectsMap[alias][host.Hostname].Count, host.Hostname, agg.Count, host.CoveredLicensesCount, licensingObjectsMap[alias][host.Hostname].Count, alias)
					}
				} else {
					if agg.Metrics == "Processor Perpetual" || agg.Metrics == "Computer Perpetual" {
						coverableLicenses := math.Min(agg.Count, licensingObjectsMap[alias][host.Hostname].Count)
						licensingObjectsMap[alias][host.Hostname].Count -= coverableLicenses
						host.CoveredLicensesCount += coverableLicenses
						agg.Count -= coverableLicenses
						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Processor Perpetual/Computer Perpetual) %f licenses to host %s. aggCount=%f associatedHostCovered=%f hostCount=%f licenseName=%s\n", coverableLicenses, host.Hostname, agg.Count, host.CoveredLicensesCount, licensingObjectsMap[alias][host.Hostname].Count, alias)
						}
					} else if agg.Metrics == "Named User Plus Perpetual" {
						coverableLicenses := math.Floor(math.Min(agg.Count*25, licensingObjectsMap[alias][host.Hostname].Count) / 25)
						licensingObjectsMap[alias][host.Hostname].Count -= coverableLicenses * 25
						host.CoveredLicensesCount += coverableLicenses * 25
						agg.Count -= coverableLicenses

						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Named User Plus Perpetual) %f(user) licenses to host %s. aggCount=%f(user) associatedHostCovered=%f hostCount=%f licenseName=%s\n", coverableLicenses, host.Hostname, agg.Count, host.CoveredLicensesCount, licensingObjectsMap[alias][host.Hostname].Count, alias)
						}
					}
				}
			}
			// If we have finished the licenses, break
			if agg.Count <= 0 && !agg.Unlimited {
				break
			}
		}
	}

	//Resort licensingObjects
	SortOracleDatabaseAgreementLicensingObjects(licensingObjects)
	licensingObjectsMap = BuildOracleDatabaseLicensingObjectsMap(licensingObjects) //the map is rebuilded because the references are updated during the sort

	// Debug print
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Resorted LicensingObjects: %#v\n", licensingObjects)
	}

	//Distribute remaining licenses in catch-all agreement to the licensingObjects
	for i := range licensingObjects {
		obj := &licensingObjects[i]

		//The object is already full covered
		if obj.Count <= 0 {
			continue
		}

		// //Debug print
		// if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		// 	as.Log.Debugf("Finding valid agreement for licensingObject #%d. obj = %s\n", i, utils.ToJSON(obj))
		// }

		//Find a agreement that can cover the object
		for j := range aggs {
			agg := &aggs[j]

			//non catch-all agreement cannot cover the object
			if !agg.CatchAll {
				continue
			}
			//non catch-all agreement cannot cover the object
			if agg.Count <= 0 && !agg.Unlimited {
				continue
			}

			//Try to fill this obj
			for _, alias := range partsMap[agg.PartID].Aliases {
				// If we have finished the licenses, break
				if agg.Count <= 0 && !agg.Unlimited {
					break
				}
				//Ignore this license because it isn't the right
				if obj.LicenseName != alias {
					continue
				}

				// fill all required license, if the host need
				if agg.Unlimited {
					// Debug print
					if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
						as.Log.Debugf("Distributing (ULA) %f licenses to obj %s. aggCount=%f objCount=0 licenseName=%s\n", obj.Count, obj.Name, agg.Count, alias)
					}

					obj.Count = 0
				} else {
					if agg.Metrics == "Processor Perpetual" || agg.Metrics == "Computer Perpetual" {
						coverableLicenses := math.Min(agg.Count, obj.Count)
						obj.Count -= coverableLicenses
						agg.Count -= coverableLicenses
						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Processor Perpetual/Computer Perpetual) %f licenses to obj %s. aggCount=%f objCount=%f licenseName=%s\n", coverableLicenses, obj.Name, agg.Count, obj.Count, alias)
						}
					} else if agg.Metrics == "Named User Plus Perpetual" {
						coverableLicenses := math.Floor(math.Min(agg.Count*25, obj.Count) / 25)
						obj.Count -= coverableLicenses * 25
						agg.Count -= coverableLicenses

						// Debug print
						if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
							as.Log.Debugf("Distributing (Named User Plus Perpetual) %f(user) licenses to obj %s. aggCount=%f(user) objCount=%f licenseName=%s\n", coverableLicenses, obj.Name, agg.Count, obj.Count, alias)
						}
					}
				}
			}
		}
	}

	// Debug print
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Associations finished. LicensingObjects: %#v\n", licensingObjects)
	}

	type coverStatus struct {
		Covered                float64 //==purchased
		TotalCoverableLicenses float64 //==consumed
	}

	//Calculate total number of covered/uncovered for each
	allLicensesCoverStatus := make(map[string]coverStatus)
	for _, obj := range licensingObjects {
		allLicensesCoverStatus[obj.LicenseName] = coverStatus{
			TotalCoverableLicenses: allLicensesCoverStatus[obj.LicenseName].TotalCoverableLicenses + obj.OriginalCount,
			Covered:                allLicensesCoverStatus[obj.LicenseName].Covered + (obj.OriginalCount - obj.Count),
		}
	}

	// Debug print
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Cover status: %#v\n", allLicensesCoverStatus)
	}

	//Calculate TotalCoveredLicenses and available
	for i := range aggs {
		agg := &aggs[i]
		uncoveredLicenseAssociatedHostSum := 0.0
		uncoveredLicenseUnassociatedObjSum := 0.0
		//calculate available
		for _, alias := range partsMap[agg.PartID].Aliases {
			uncoveredLicenseUnassociatedObjSum += allLicensesCoverStatus[alias].TotalCoverableLicenses - allLicensesCoverStatus[alias].Covered
			for j, host := range agg.Hosts {
				// If no host require a license with licenseName == alias, skip
				if _, ok := licensingObjectsMap[alias]; !ok {
					continue
				}
				// If the host don't use the license, skip
				if _, ok := licensingObjectsMap[alias][host.Hostname]; !ok {
					continue
				}
				agg.Hosts[j].TotalCoveredLicensesCount = licensingObjectsMap[alias][host.Hostname].OriginalCount - licensingObjectsMap[alias][host.Hostname].Count
				uncoveredLicenseAssociatedHostSum += licensingObjectsMap[alias][host.Hostname].Count //non-covered part
			}
		}

		if !agg.CatchAll {
			agg.AvailableCount = -uncoveredLicenseAssociatedHostSum
		} else {
			agg.AvailableCount = -uncoveredLicenseUnassociatedObjSum
		}
	}
}

// CheckOracleDatabaseAgreementMatchFilter check that agg match the filters
func CheckOracleDatabaseAgreementMatchFilter(agg apimodel.OracleDatabaseAgreementsFE, filters apimodel.SearchOracleDatabaseAgreementsFilters) bool {
	return strings.Contains(strings.ToLower(agg.AgreementID), strings.ToLower(filters.AgreementID)) &&
		strings.Contains(strings.ToLower(agg.PartID), strings.ToLower(filters.PartID)) &&
		strings.Contains(strings.ToLower(agg.ItemDescription), strings.ToLower(filters.ItemDescription)) &&
		strings.Contains(strings.ToLower(agg.CSI), strings.ToLower(filters.CSI)) &&
		(filters.Metrics == "" || strings.ToLower(agg.Metrics) == strings.ToLower(filters.Metrics)) &&
		strings.Contains(strings.ToLower(agg.ReferenceNumber), strings.ToLower(filters.ReferenceNumber)) &&
		(filters.Unlimited == "NULL" || agg.Unlimited == (filters.Unlimited == "true")) &&
		(filters.CatchAll == "NULL" || agg.CatchAll == (filters.CatchAll == "true")) &&
		(filters.LicensesCountLTE == -1 || agg.LicensesCount <= float64(filters.LicensesCountLTE)) &&
		(filters.LicensesCountGTE == -1 || agg.LicensesCount >= float64(filters.LicensesCountGTE)) &&
		(filters.UsersCountLTE == -1 || agg.UsersCount <= float64(filters.UsersCountLTE)) &&
		(filters.UsersCountGTE == -1 || agg.UsersCount >= float64(filters.UsersCountGTE)) &&
		(filters.AvailableCountLTE == -1 || agg.AvailableCount <= float64(filters.AvailableCountLTE)) &&
		(filters.AvailableCountGTE == -1 || agg.AvailableCount >= float64(filters.AvailableCountGTE))
}

// SortOracleDatabaseAgreementLicensingObjects sort the list of apimodel.OracleDatabaseLicensingObjects by count
func SortOracleDatabaseAgreementLicensingObjects(obj []apimodel.OracleDatabaseLicensingObjects) {
	sort.Slice(obj, func(i, j int) bool {
		if obj[i].Count != obj[j].Count {
			return obj[i].Count > obj[j].Count
		} else if obj[i].Name != obj[j].Name {
			return obj[i].Name > obj[j].Name
		} else {
			return obj[i].LicenseName > obj[j].LicenseName
		}
	})
}

// SortOracleDatabaseAgreements sort the list of apimodel.OracleDatabaseAgreementsFE for GreedilyAssignOracleDatabaseAgreementsToLicensingObjects algorithm
func SortOracleDatabaseAgreements(obj []apimodel.OracleDatabaseAgreementsFE) {
	sort.Slice(obj, func(i, j int) bool {
		if !obj[i].CatchAll && obj[j].CatchAll {
			return true
		} else if obj[i].CatchAll && !obj[j].CatchAll {
			return false
		} else if !obj[i].Unlimited && obj[j].Unlimited {
			return true
		} else if obj[i].Unlimited && !obj[j].Unlimited {
			return false
		} else if obj[i].UsersCount != obj[j].UsersCount {
			return obj[i].UsersCount > obj[j].UsersCount
		} else {
			return obj[i].LicensesCount > obj[j].LicensesCount
		}
	})
}

// SortAssociatedHostsInOracleDatabaseAgreement sort the associated hosts by license count. It  that parts may have multiple aliases
func SortAssociatedHostsInOracleDatabaseAgreement(agg apimodel.OracleDatabaseAgreementsFE, licensingObjectsMap map[string]map[string]*apimodel.OracleDatabaseLicensingObjects, partsMap map[string]*model.OracleDatabaseAgreementPart) {
	sort.Slice(agg.Hosts, func(i, j int) bool {
		var maxLicensingObjectICount float64 = 0
		var maxLicensingObjectJCount float64 = 0
		for _, alias := range partsMap[agg.PartID].Aliases {
			if _, ok := licensingObjectsMap[alias]; ok {
				if _, ok := licensingObjectsMap[alias][agg.Hosts[i].Hostname]; ok {
					maxLicensingObjectICount = math.Max(maxLicensingObjectICount, licensingObjectsMap[alias][agg.Hosts[i].Hostname].Count)
				}
				if _, ok := licensingObjectsMap[alias][agg.Hosts[j].Hostname]; ok {
					maxLicensingObjectJCount = math.Max(maxLicensingObjectJCount, licensingObjectsMap[alias][agg.Hosts[j].Hostname].Count)
				}
			}
		}
		return maxLicensingObjectICount > maxLicensingObjectJCount
	})
}

// BuildOracleDatabaseLicensingObjectsMap return a map of license name to map of object name to pointer to  apimodel.OracleDatabaseLicensingObjects for fast object lookup
// BuildOracleDatabaseLicensingObjectsMap assume that doesn't exist a cluster and a host with the same name
func BuildOracleDatabaseLicensingObjectsMap(objs []apimodel.OracleDatabaseLicensingObjects) map[string]map[string]*apimodel.OracleDatabaseLicensingObjects {
	res := make(map[string]map[string]*apimodel.OracleDatabaseLicensingObjects)

	for i, obj := range objs {
		if _, ok := res[obj.LicenseName]; !ok {
			res[obj.LicenseName] = make(map[string]*apimodel.OracleDatabaseLicensingObjects)
		}
		res[obj.LicenseName][obj.Name] = &objs[i]
	}

	return res
}

// BuildOracleDatabaseAgreementPartMap return a map of partID to OracleDatabaseAgreementPart
func BuildOracleDatabaseAgreementPartMap(parts []model.OracleDatabaseAgreementPart) map[string]*model.OracleDatabaseAgreementPart {
	res := make(map[string]*model.OracleDatabaseAgreementPart)

	for i, part := range parts {
		res[part.PartID] = &parts[i]
	}

	return res
}
