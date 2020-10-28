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
	"math"
	"sort"
	"strings"

	"github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/api-service/dto"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddAssociatedPartToOracleDbAgreement add associated part to OracleDatabaseAgreement or create a new one
func (as *APIService) AddAssociatedPartToOracleDbAgreement(request dto.AssociatedPartInOracleDbAgreementRequest,
) (primitive.ObjectID, utils.AdvancedErrorInterface) {
	if err := checkHosts(as, request.Hosts); err != nil {
		return primitive.NilObjectID, err
	}

	agreement, err := as.Database.GetOracleDatabaseAgreement(request.AgreementID)
	if err == utils.AerrOracleDatabaseAgreementNotFound {
		agreement = &model.OracleDatabaseAgreement{
			AgreementID: request.AgreementID,
			CSI:         request.CSI,
			Parts:       make([]model.AssociatedPart, 0),
		}

	} else if err != nil {
		return primitive.NilObjectID, err
	}

	if err := addAssociatedPart(as, agreement, request); err != nil {
		return primitive.NilObjectID, err
	}

	if agreement.ID == primitive.NilObjectID {
		res, err := as.Database.InsertOracleDatabaseAgreement(*agreement)
		if err != nil {
			return primitive.NilObjectID, err
		}

		agreement.ID = res.InsertedID.(primitive.ObjectID)
	} else {
		err := as.Database.UpdateOracleDatabaseAgreement(*agreement)
		if err != nil {
			return primitive.NilObjectID, err
		}
	}

	return agreement.ID, nil
}

func addAssociatedPart(as *APIService, agreement *model.OracleDatabaseAgreement,
	req dto.AssociatedPartInOracleDbAgreementRequest) utils.AdvancedErrorInterface {
	part, err := as.GetOraclePart(req.PartID)
	if err != nil {
		return err
	}

	associatedPart := model.AssociatedPart{
		OracleDatabasePart: *part,
		ReferenceNumber:    req.ReferenceNumber,
		Unlimited:          req.Unlimited,
		Count:              req.Count,
		CatchAll:           req.CatchAll,
		Hosts:              req.Hosts,
	}
	agreement.Parts = append(agreement.Parts, associatedPart)

	return nil
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

// UpdateAssociatedPartOfOracleDbAgreement update associated part in OracleDatabaseAgreement
func (as *APIService) UpdateAssociatedPartOfOracleDbAgreement(request dto.AssociatedPartInOracleDbAgreementRequest,
) utils.AdvancedErrorInterface {
	if err := checkHosts(as, request.Hosts); err != nil {
		return err
	}

	associatedPartID := utils.Str2oid(request.ID)
	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedPart(associatedPartID)
	if err != nil {
		return err
	}

	err = updateAssociatedPart(as, agreement, request)
	if err != nil {
		return err
	}

	return as.Database.UpdateOracleDatabaseAgreement(*agreement)
}

func updateAssociatedPart(as *APIService, agreement *model.OracleDatabaseAgreement,
	req dto.AssociatedPartInOracleDbAgreementRequest) utils.AdvancedErrorInterface {

	var associatedPart *model.AssociatedPart
	reqID := utils.Str2oid(req.ID)

	for i := range agreement.Parts {
		if agreement.Parts[i].ID == reqID {
			associatedPart = &agreement.Parts[i]
			break
		}
	}

	if associatedPart == nil {
		return utils.AerrOracleDatabaseAssociatedPartNotFound
	}

	part, err := as.GetOraclePart(req.PartID)
	if err != nil {
		return err
	}
	associatedPart.OracleDatabasePart = *part
	associatedPart.ReferenceNumber = req.ReferenceNumber
	associatedPart.Unlimited = req.Unlimited
	associatedPart.Count = req.Count
	associatedPart.CatchAll = req.CatchAll
	associatedPart.Hosts = req.Hosts

	return nil
}

// SearchAssociatedPartsInOracleDatabaseAgreements search Oracle/Database agreements
func (as *APIService) SearchAssociatedPartsInOracleDatabaseAgreements(
	search string,
	filters dto.SearchOracleDatabaseAgreementsFilter,
) ([]dto.OracleDatabaseAgreementFE, utils.AdvancedErrorInterface) {
	agreements, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}
	// TODO Insert Agr Part info
	// TODO Check thet parts exists, if not warn?

	hosts, err := as.Database.ListHostUsingOracleDatabaseLicenses()
	if err != nil {
		return nil, err
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	filteredAgrs := make([]dto.OracleDatabaseAgreementFE, 0)
	for _, agr := range agreements {

		if checkOracleDatabaseAgreementMatchFilter(agr, filters) {
			filteredAgrs = append(filteredAgrs, agr)
		}

	}

	return filteredAgrs, nil
}

// checkOracleDatabaseAgreementMatchFilter check that agr match the filters
func checkOracleDatabaseAgreementMatchFilter(agr dto.OracleDatabaseAgreementFE, filters dto.SearchOracleDatabaseAgreementsFilter) bool {
	return strings.Contains(strings.ToLower(agr.AgreementID), strings.ToLower(filters.AgreementID)) &&
		strings.Contains(strings.ToLower(agr.PartID), strings.ToLower(filters.PartID)) &&
		strings.Contains(strings.ToLower(agr.ItemDescription), strings.ToLower(filters.ItemDescription)) &&
		strings.Contains(strings.ToLower(agr.CSI), strings.ToLower(filters.CSI)) &&
		(filters.Metric == "" || strings.ToLower(agr.Metric) == strings.ToLower(filters.Metric)) &&
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

// assignOracleDatabaseAgreementsToHosts assign available licenses in each agreements to hosts using licenses
func (as *APIService) assignOracleDatabaseAgreementsToHosts(
	agrs []dto.OracleDatabaseAgreementFE,
	hosts []dto.HostUsingOracleDatabaseLicenses) {

	sortOracleDatabaseAgreements(agrs)
	sortHostsUsingLicenses(hosts)

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Agreements = %s\nHosts= %s\n", utils.ToJSON(agrs), utils.ToJSON(hosts))
	}

	hostsMap := buildHostUsingLicensesMap(hosts)
	partsMap := buildAgreementPartMap(as.OracleDatabaseAgreementParts)

	assignLicensesInAgreementsToAssociatedHost(as, agrs, hostsMap, partsMap)

	// sort again and rebuild map because the references are updated during the sort
	sortHostsUsingLicenses(hosts)
	hostsMap = buildHostUsingLicensesMap(hosts)

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Resorted LicensingObjects: %#v\n", hosts)
	}

	distributeLicensesInCatchAllAgrs(as, agrs, hosts, partsMap)

	allLicensesCoverStatus := calculateCoverStatusByLicenseName(hosts)
	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Cover status: %#v\n", allLicensesCoverStatus)
	}

	calculateTotalCoveredLicensesAndAvailable(agrs, hostsMap, partsMap, allLicensesCoverStatus)
}

// sortOracleDatabaseAgreements sort the list of dto.OracleDatabaseAgreementsFE
// by CatchAll (falses first), Unlimited (falses first), decreasing UsersCount, decreasing LicensesCount
func sortOracleDatabaseAgreements(obj []dto.OracleDatabaseAgreementFE) {
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

// sortHostsUsingLicenses sort the list of hosts by decreasing license count,
// alphabetical name, alphabetical license name
func sortHostsUsingLicenses(obj []dto.HostUsingOracleDatabaseLicenses) {
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

// buildHostUsingLicensesMap return a map of license name to map of object name to pointer to  dto.HostUsingOracleDatabaseLicenses for fast object lookup
// Assume that doesn't exist a cluster and a host with the same name
func buildHostUsingLicensesMap(hosts []dto.HostUsingOracleDatabaseLicenses) map[string]map[string]*dto.HostUsingOracleDatabaseLicenses {
	res := make(map[string]map[string]*dto.HostUsingOracleDatabaseLicenses)

	for i, host := range hosts {
		if _, ok := res[host.LicenseName]; !ok {
			res[host.LicenseName] = make(map[string]*dto.HostUsingOracleDatabaseLicenses)
		}
		res[host.LicenseName][host.Name] = &hosts[i]
	}

	return res
}

// buildAgreementPartMap return a map of partID to OracleDatabaseAgreementPart
func buildAgreementPartMap(parts []model.OracleDatabasePart) map[string]*model.OracleDatabasePart {
	partsMap := make(map[string]*model.OracleDatabasePart)

	for i, part := range parts {
		partsMap[part.PartID] = &parts[i]
	}

	return partsMap
}

// Assign available licenses in each agreement to each host associated in each agreement
// if this host is using that kind of license.
func assignLicensesInAgreementsToAssociatedHost(
	as *APIService,
	agreements []dto.OracleDatabaseAgreementFE,
	hostsMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses,
	partsMap map[string]*model.OracleDatabasePart) {

	for i := range agreements {
		agreement := &agreements[i]
		sortHostsInAgreementByLicenseCount(agreement, hostsMap, partsMap)

		if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
			as.Log.Debugf("Distributing licenses of agreement #%d to host. Agreement = %s\n", i, utils.ToJSON(agreement))
		}

		//distribute licenses for each host
		for j := range agreement.Hosts {
			hostInAgr := &agreement.Hosts[j]

			for _, alias := range partsMap[agreement.PartID].Aliases {

				if agreement.Count <= 0 && !agreement.Unlimited {
					break
				}

				if _, ok := hostsMap[alias]; !ok {
					// no host use this license
					continue
				}

				var hostUsingLicenses *dto.HostUsingOracleDatabaseLicenses
				var ok bool
				if hostUsingLicenses, ok = hostsMap[alias][hostInAgr.Hostname]; !ok {
					// host doesn't use this license
					continue
				}

				if hostUsingLicenses.LicenseCount <= 0 {
					continue
				}

				switch {
				case agreement.Unlimited:
					hostInAgr.CoveredLicensesCount = hostUsingLicenses.LicenseCount
					hostUsingLicenses.LicenseCount = 0

				case agreement.Metric == model.AgreementPartMetricProcessorPerpetual ||
					agreement.Metric == model.AgreementPartMetricComputerPerpetual:

					coverableLicenses := math.Min(agreement.Count, hostUsingLicenses.LicenseCount)
					hostUsingLicenses.LicenseCount -= coverableLicenses
					hostInAgr.CoveredLicensesCount += coverableLicenses
					agreement.Count -= coverableLicenses

				case agreement.Metric == model.AgreementPartMetricNamedUserPlusPerpetual:
					coverableLicenses := math.Floor(math.Min(agreement.Count*25, hostUsingLicenses.LicenseCount) / 25)
					hostUsingLicenses.LicenseCount -= coverableLicenses * 25
					hostInAgr.CoveredLicensesCount += coverableLicenses * 25
					agreement.Count -= coverableLicenses

				default:
					as.Log.Errorf("Distributing licenses. Unknown metric type: [%s]", agreement.Metric)
				}

				if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
					as.Log.Debugf(`Distributing %f licenses to host %s. agr.Metrics=%s agr.Count=%f \
					hostInAgr.CoveredLicensesCount=%f hostUsingLicenses.LicenseCount=%f licenseName=%s\n`,
						hostUsingLicenses.LicenseCount,
						hostInAgr.Hostname,
						agreement.Metric,
						agreement.Count,
						hostInAgr.CoveredLicensesCount,
						hostUsingLicenses.LicenseCount,
						alias)
				}
			}

			if agreement.Count <= 0 && !agreement.Unlimited {
				break
			}
		}

		//TODO should I move distributeLicensesInCatchAllAgrs here?
	}
}

// sortHostsInAgreementByLicenseCount sort the associated hosts by license count
// considering that parts may have multiple aliases
func sortHostsInAgreementByLicenseCount(agr *dto.OracleDatabaseAgreementFE,
	hostsMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses,
	partsMap map[string]*model.OracleDatabasePart) {

	sort.Slice(agr.Hosts, func(i, j int) bool {
		var iMaxLicenseCount float64 = 0
		var jMaxLicenseCount float64 = 0

		for _, alias := range partsMap[agr.PartID].Aliases {
			if mapHostnamesLicenses, ok := hostsMap[alias]; ok {

				if hostUsingLicenses, ok := mapHostnamesLicenses[agr.Hosts[i].Hostname]; ok {
					iMaxLicenseCount = math.Max(iMaxLicenseCount, hostUsingLicenses.LicenseCount)
				}

				if hostUsingLicenses, ok := mapHostnamesLicenses[agr.Hosts[j].Hostname]; ok {
					jMaxLicenseCount = math.Max(jMaxLicenseCount, hostUsingLicenses.LicenseCount)
				}
			}
		}
		return iMaxLicenseCount > jMaxLicenseCount
	})
}

// If an agreement is catchAll (or basket..) distributes its licenses to every hosts that use that kind of license
func distributeLicensesInCatchAllAgrs(
	as *APIService,
	agrs []dto.OracleDatabaseAgreementFE,
	hosts []dto.HostUsingOracleDatabaseLicenses,
	partsMap map[string]*model.OracleDatabasePart) {

	for i := range hosts {
		host := &hosts[i]

		if host.LicenseCount <= 0 {
			continue
		}

		// //Debug print
		// if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		// 	as.Log.Debugf("Finding valid agreement for licensingObject #%d. obj = %s\n", i, utils.ToJSON(obj))
		// }

		for j := range agrs {
			agr := &agrs[j]

			if !agr.CatchAll {
				continue
			}

			if agr.Count <= 0 && !agr.Unlimited {
				continue
			}

			for _, alias := range partsMap[agr.PartID].Aliases {
				if agr.Count <= 0 && !agr.Unlimited {
					break
				}

				if host.LicenseName != alias {
					continue
				}

				switch {

				case agr.Unlimited:
					host.LicenseCount = 0

				case agr.Metric == model.AgreementPartMetricProcessorPerpetual || agr.Metric == model.AgreementPartMetricComputerPerpetual:
					coverableLicenses := math.Min(agr.Count, host.LicenseCount)
					host.LicenseCount -= coverableLicenses
					agr.Count -= coverableLicenses

				case agr.Metric == model.AgreementPartMetricNamedUserPlusPerpetual:
					coverableLicenses := math.Floor(math.Min(agr.Count*25, host.LicenseCount) / 25)
					host.LicenseCount -= coverableLicenses * 25
					agr.Count -= coverableLicenses
				}

				if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
					as.Log.Debugf("Distributing with metric [%s] [ULA? %t] %f licenses to obj %s. aggCount=%f objCount=0 licenseName=%s\n",
						agr.Metric,
						agr.Unlimited,
						host.LicenseCount,
						host.Name,
						agr.Count,
						alias)
				}
			}
		}
	}

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Associations finished. LicensingObjects: %#v\n", hosts)
	}
}

type coverStatus struct {
	Covered                float64 //==purchased
	TotalCoverableLicenses float64 //==consumed
}

// Calculate total number of covered/uncovered for each host
func calculateCoverStatusByLicenseName(hosts []dto.HostUsingOracleDatabaseLicenses) map[string]coverStatus {
	allLicensesCoverStatus := make(map[string]coverStatus)

	for _, host := range hosts {
		allLicensesCoverStatus[host.LicenseName] = coverStatus{
			TotalCoverableLicenses: allLicensesCoverStatus[host.LicenseName].TotalCoverableLicenses + host.OriginalCount,
			Covered:                allLicensesCoverStatus[host.LicenseName].Covered + (host.OriginalCount - host.LicenseCount),
		}
	}

	return allLicensesCoverStatus
}

// Calculate TotalCoveredLicenses and available
func calculateTotalCoveredLicensesAndAvailable(
	agrs []dto.OracleDatabaseAgreementFE,
	hostsMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses,
	partsMap map[string]*model.OracleDatabasePart,
	allLicensesCoverStatus map[string]coverStatus) {

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
				if _, ok := hostsMap[alias]; !ok {
					continue
				}

				// If the host don't use the license, skip
				if _, ok := hostsMap[alias][host.Hostname]; !ok {
					continue
				}
				host.TotalCoveredLicensesCount = hostsMap[alias][host.Hostname].OriginalCount - hostsMap[alias][host.Hostname].LicenseCount
				host.ConsumedLicensesCount = hostsMap[alias][host.Hostname].OriginalCount
				uncoveredLicenseAssociatedHostSum += hostsMap[alias][host.Hostname].LicenseCount //non-covered part
			}
		}

		if !agr.CatchAll {
			agr.AvailableCount = -uncoveredLicenseAssociatedHostSum
		} else {
			agr.AvailableCount = -uncoveredLicenseUnassociatedObjSum
		}
	}
}

// DeleteAssociatedPartFromOracleDatabaseAgreement remove an Oracle/Database agreement
func (as *APIService) DeleteAssociatedPartFromOracleDatabaseAgreement(associatedPartID primitive.ObjectID,
) utils.AdvancedErrorInterface {
	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedPart(associatedPartID)
	if err != nil {
		return err
	}

	if len(agreement.Parts) <= 1 {
		return as.Database.RemoveOracleDatabaseAgreement(agreement.ID)
	}

	for i := range agreement.Parts {
		if agreement.Parts[i].ID == associatedPartID {
			agreement.Parts = append(agreement.Parts[:i], agreement.Parts[i+1])
			break
		}
	}

	return as.Database.UpdateOracleDatabaseAgreement(*agreement)
}

// AddHostToAssociatedPart a new host to the list of associated hosts of the agreement
func (as *APIService) AddHostToAssociatedPart(associatedPartID primitive.ObjectID, hostname string,
) utils.AdvancedErrorInterface {

	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedPart(associatedPartID)
	if err != nil {
		return err
	}

	associatedPart := agreement.AssociatedPartByID(associatedPartID)

	for _, host := range associatedPart.Hosts {
		if host == hostname {
			return nil
		}
	}

	if err := checkHosts(as, []string{hostname}); err != nil {
		return err
	}

	associatedPart.Hosts = append(associatedPart.Hosts, hostname)

	return as.Database.UpdateOracleDatabaseAgreement(*agreement)
}

// RemoveHostFromAssociatedPart remove the host from the list of associated hosts of the agreement
func (as *APIService) RemoveHostFromAssociatedPart(associatedPartID primitive.ObjectID, hostname string,
) utils.AdvancedErrorInterface {

	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedPart(associatedPartID)
	if err != nil {
		return err
	}

	associatedPart := agreement.AssociatedPartByID(associatedPartID)

	for i, host := range associatedPart.Hosts {
		if host == hostname {
			associatedPart.Hosts = append(associatedPart.Hosts[:i], associatedPart.Hosts[i+1:]...)

			return as.Database.UpdateOracleDatabaseAgreement(*agreement)
		}
	}

	return nil
}
