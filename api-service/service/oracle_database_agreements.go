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

//TODO Instead of use 25 everywhere for NamedUserPlus licenses, use const

// TODO When insert or update unlimited agr, set count == 0

// AddAssociatedPartToOracleDbAgreement add associated part to OracleDatabaseAgreement or create a new one
func (as *APIService) AddAssociatedPartToOracleDbAgreement(request dto.AssociatedPartInOracleDbAgreementRequest,
) (string, utils.AdvancedErrorInterface) {
	if err := checkHosts(as, request.Hosts); err != nil {
		return "", err
	}

	agreement, err := as.Database.GetOracleDatabaseAgreement(request.AgreementID)
	if err == utils.AerrOracleDatabaseAgreementNotFound {
		agreement = &model.OracleDatabaseAgreement{
			AgreementID: request.AgreementID,
			CSI:         request.CSI,
			Parts:       make([]model.AssociatedPart, 0),
		}

	} else if err != nil {
		return "", err
	}

	if err := addAssociatedPart(as, agreement, request); err != nil {
		return "", err
	}

	if agreement.ID == primitive.NilObjectID {
		agreement.ID = as.NewObjectID()

		res, err := as.Database.InsertOracleDatabaseAgreement(*agreement)
		if err != nil {
			return "", err
		}

		agreement.ID = res.InsertedID.(primitive.ObjectID)
	} else {
		err := as.Database.UpdateOracleDatabaseAgreement(*agreement)
		if err != nil {
			return "", err
		}
	}

	return agreement.ID.Hex(), nil
}

func addAssociatedPart(as *APIService, agreement *model.OracleDatabaseAgreement,
	req dto.AssociatedPartInOracleDbAgreementRequest) utils.AdvancedErrorInterface {
	part, err := as.GetOraclePart(req.PartID)
	if err != nil {
		return err
	}

	associatedPart := model.AssociatedPart{
		ID:                 as.NewObjectID(),
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

// SearchAssociatedPartsInOracleDatabaseAgreements search OracleDatabase associated parts agreements
func (as *APIService) SearchAssociatedPartsInOracleDatabaseAgreements(filters dto.SearchOracleDatabaseAgreementsFilter,
) ([]dto.OracleDatabaseAgreementFE, utils.AdvancedErrorInterface) {
	agreements, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	parts := buildAgreementPartMap(as.OracleDatabaseAgreementParts)
	for i := range agreements {
		agr := &agreements[i]

		if part, ok := parts[agr.PartID]; ok {
			agr.ItemDescription = part.ItemDescription
			agr.Metric = part.Metric

			switch agr.Metric {
			case model.AgreementPartMetricProcessorPerpetual:
				agr.LicensesCount = agr.Count
			case model.AgreementPartMetricNamedUserPlusPerpetual:
				agr.UsersCount = agr.Count
			}
		} else {
			as.Log.Errorf("Unknown PartID: [%s] in agreement: [%#v]", agr.PartID, agr)
		}
	}

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

	assignAgreementsLicensesToItsAssociatedHosts(as, agrs, hostsMap, partsMap)

	// sort again and rebuild map because the references are updated during the sort
	sortHostsUsingLicenses(hosts)
	hostsMap = buildHostUsingLicensesMap(hosts)

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Resorted LicensingObjects: %#v\n", hosts)
	}

	assignLicensesFromCatchAllAgreements(as, agrs, hosts, partsMap)

	calculateTotalCoveredLicensesAndAvailableCount(as, agrs, hosts, hostsMap, partsMap)
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

// buildHostUsingLicensesMap return a map of license name to map of object name to pointer to
// dto.HostUsingOracleDatabaseLicenses for fast object lookup
// Assume that doesn't exist a cluster and a host with the same name
func buildHostUsingLicensesMap(hosts []dto.HostUsingOracleDatabaseLicenses,
) map[string]map[string]*dto.HostUsingOracleDatabaseLicenses {

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
func assignAgreementsLicensesToItsAssociatedHosts(
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

		for j := range agreement.Hosts {
			associatedHost := &agreement.Hosts[j]

			for _, alias := range partsMap[agreement.PartID].Aliases {

				if agreement.AvailableCount <= 0 && !agreement.Unlimited {
					break
				}

				if _, ok := hostsMap[alias]; !ok {
					// no host use this license
					continue
				}

				var hostUsingLicenses *dto.HostUsingOracleDatabaseLicenses
				var ok bool
				if hostUsingLicenses, ok = hostsMap[alias][associatedHost.Hostname]; !ok {
					// host doesn't use this license
					continue
				}

				if hostUsingLicenses.LicenseCount <= 0 {
					continue
				}

				doAssignAgreementLicensesToAssociatedHost(as, agreement, hostUsingLicenses, associatedHost)

				if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
					as.Log.Debugf(`Distributing %f licenses to host %s. agr.Metrics=%s agr.AvailableCount=%f \
					hostInAgr.CoveredLicensesCount=%f hostUsingLicenses.LicenseCount=%f licenseName=%s\n`,
						hostUsingLicenses.LicenseCount,
						associatedHost.Hostname,
						agreement.Metric,
						agreement.AvailableCount,
						associatedHost.CoveredLicensesCount,
						hostUsingLicenses.LicenseCount,
						alias)
				}
			}

			if agreement.AvailableCount <= 0 && !agreement.Unlimited {
				break
			}
		}
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

// Use all the licenses available in agreement to cover host and associatedHost if provided
func doAssignAgreementLicensesToAssociatedHost(
	as *APIService,
	agreement *dto.OracleDatabaseAgreementFE,
	host *dto.HostUsingOracleDatabaseLicenses,
	associatedHost *dto.OracleDatabaseAgreementAssociatedHostFE) {

	switch {
	case agreement.Metric == model.AgreementPartMetricProcessorPerpetual ||
		agreement.Metric == model.AgreementPartMetricComputerPerpetual:

		var coverableLicenses float64
		if agreement.Unlimited {
			coverableLicenses = host.LicenseCount
			agreement.AvailableCount = 0
		} else {
			coverableLicenses = math.Min(agreement.AvailableCount, host.LicenseCount)
			agreement.AvailableCount -= coverableLicenses
		}

		associatedHost.CoveredLicensesCount += coverableLicenses

		host.LicenseCount -= coverableLicenses

	case agreement.Metric == model.AgreementPartMetricNamedUserPlusPerpetual:

		var coverableLicenses float64
		if agreement.Unlimited {
			coverableLicenses = host.LicenseCount
			agreement.AvailableCount = 0
		} else {
			coverableLicenses = math.Floor(math.Min(agreement.AvailableCount, host.LicenseCount*25) / 25)
			agreement.AvailableCount -= coverableLicenses * 25
		}

		associatedHost.CoveredLicensesCount += coverableLicenses * 25

		host.LicenseCount -= coverableLicenses

	default:
		as.Log.Errorf("Distributing licenses. Unknown metric type: [%s]", agreement.Metric)
	}
}

// If an agreement is catchAll (or basket..) distributes its licenses to every hosts that use that kind of license
func assignLicensesFromCatchAllAgreements(
	as *APIService,
	agrs []dto.OracleDatabaseAgreementFE,
	hosts []dto.HostUsingOracleDatabaseLicenses,
	partsMap map[string]*model.OracleDatabasePart) {

	for i := range hosts {
		host := &hosts[i]

		if host.LicenseCount <= 0 {
			continue
		}

		for j := range agrs {
			agr := &agrs[j]

			if !agr.CatchAll {
				continue
			}

			if agr.AvailableCount <= 0 && !agr.Unlimited {
				continue
			}

			for _, alias := range partsMap[agr.PartID].Aliases {
				if agr.AvailableCount <= 0 && !agr.Unlimited {
					break
				}

				if host.LicenseName != alias {
					continue
				}

				doAssignLicenseFromCatchAllAgreement(as, agr, host)

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

// Use all the licenses available in agreement to cover host and associatedHost if provided
func doAssignLicenseFromCatchAllAgreement(
	as *APIService,
	agreement *dto.OracleDatabaseAgreementFE,
	hostUsingLicenses *dto.HostUsingOracleDatabaseLicenses) {

	switch {
	case agreement.Metric == model.AgreementPartMetricProcessorPerpetual ||
		agreement.Metric == model.AgreementPartMetricComputerPerpetual:

		var coverableLicenses float64
		if agreement.Unlimited {
			coverableLicenses = hostUsingLicenses.LicenseCount
			agreement.AvailableCount = 0
		} else {
			coverableLicenses = math.Min(agreement.AvailableCount, hostUsingLicenses.LicenseCount)
			agreement.AvailableCount -= coverableLicenses
		}

		hostUsingLicenses.LicenseCount -= coverableLicenses

	case agreement.Metric == model.AgreementPartMetricNamedUserPlusPerpetual:

		var coverableLicenses float64
		if agreement.Unlimited {
			coverableLicenses = hostUsingLicenses.LicenseCount
			agreement.AvailableCount = 0
		} else {
			coverableLicenses = math.Floor(math.Min(agreement.AvailableCount, hostUsingLicenses.LicenseCount*25) / 25)
			agreement.AvailableCount -= coverableLicenses * 25
		}

		hostUsingLicenses.LicenseCount -= coverableLicenses

	default:
		as.Log.Errorf("Distributing licenses. Unknown metric type: [%s]", agreement.Metric)
	}
}

func calculateTotalCoveredLicensesAndAvailableCount(
	as *APIService,
	agrs []dto.OracleDatabaseAgreementFE,
	hosts []dto.HostUsingOracleDatabaseLicenses,
	hostsMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses,
	partsMap map[string]*model.OracleDatabasePart) {

	licensesCoverStatusByName := calculateCoverStatusByLicenseName(hosts)

	for i := range agrs {
		agreement := &agrs[i]

		uncoveredLicensesByAssociatedHosts := 0.0
		uncoveredLicensesByAllHosts := 0.0

		for _, alias := range partsMap[agreement.PartID].Aliases {

			for j := range agreement.Hosts {
				associatedHost := &agreement.Hosts[j]
				if _, ok := hostsMap[alias]; !ok {
					continue
				}

				host, ok := hostsMap[alias][associatedHost.Hostname]
				if !ok {
					continue
				}

				switch {
				case agreement.Metric == model.AgreementPartMetricProcessorPerpetual ||
					agreement.Metric == model.AgreementPartMetricComputerPerpetual:
					associatedHost.TotalCoveredLicensesCount = host.OriginalCount - host.LicenseCount
					associatedHost.ConsumedLicensesCount = host.OriginalCount
					uncoveredLicensesByAssociatedHosts += host.LicenseCount

				case agreement.Metric == model.AgreementPartMetricNamedUserPlusPerpetual:
					associatedHost.TotalCoveredLicensesCount = (host.OriginalCount - host.LicenseCount) * 25
					associatedHost.ConsumedLicensesCount = host.OriginalCount * 25
					uncoveredLicensesByAssociatedHosts += host.LicenseCount * 25

				default:
					as.Log.Errorf("Unknown metric type: [%s]", agreement.Metric)
				}
			}

			uncoveredLicensesByAllHosts += licensesCoverStatusByName[alias].Consumed - licensesCoverStatusByName[alias].Covered
		}

		var uncoveredLicenses float64
		if agreement.CatchAll {
			uncoveredLicenses = uncoveredLicensesByAllHosts
		} else {
			uncoveredLicenses = uncoveredLicensesByAssociatedHosts
		}

		if uncoveredLicenses > 0 {
			if (agreement.AvailableCount > 0 && agreement.Metric != model.AgreementPartMetricNamedUserPlusPerpetual) ||
				(agreement.AvailableCount > 25 && agreement.Metric == model.AgreementPartMetricNamedUserPlusPerpetual) {

				as.Log.Errorf("Agreement has still some available licenses but hosts are uncovered. Agreement: [%v]",
					agreement)
			}

			agreement.AvailableCount -= uncoveredLicenses
		}
	}
}

type coverStatus struct {
	Covered  float64 //==purchased
	Consumed float64 //==consumed
}

// Calculate total number of covered/uncovered for each host
func calculateCoverStatusByLicenseName(hosts []dto.HostUsingOracleDatabaseLicenses) map[string]coverStatus {
	licensesStatus := make(map[string]coverStatus)

	for _, host := range hosts {
		licensesStatus[host.LicenseName] = coverStatus{
			Consumed: licensesStatus[host.LicenseName].Consumed + host.OriginalCount,
			Covered:  licensesStatus[host.LicenseName].Covered + (host.OriginalCount - host.LicenseCount),
		}
	}

	return licensesStatus
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

// DeleteAssociatedPartFromOracleDatabaseAgreement delete associated part from OracleDatabaseAgreement
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

// AddHostToAssociatedPart add an host to AssociatedPart
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

// RemoveHostFromAssociatedPart remove host from AssociatedPart
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
