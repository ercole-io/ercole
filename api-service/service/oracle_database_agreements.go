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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TODO Instead of use 25 everywhere for NamedUserPlus licenses, use const

// TODO When insert or update unlimited agr, set count == 0

// AddAssociatedLicenseTypeToOracleDbAgreement add associated part to OracleDatabaseAgreement or create a new one
func (as *APIService) AddAssociatedLicenseTypeToOracleDbAgreement(request dto.AssociatedLicenseTypeInOracleDbAgreementRequest,
) (string, error) {
	if err := checkHosts(as, request.Hosts); err != nil {
		return "", err
	}

	agreement, err := as.Database.GetOracleDatabaseAgreement(request.AgreementID)
	if err == utils.ErrOracleDatabaseAgreementNotFound {
		agreement = &model.OracleDatabaseAgreement{
			AgreementID:  request.AgreementID,
			CSI:          request.CSI,
			LicenseTypes: make([]model.AssociatedLicenseType, 0),
		}

	} else if err != nil {
		return "", err
	}

	if err := addAssociatedLicenseType(as, agreement, request); err != nil {
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

func addAssociatedLicenseType(as *APIService, agreement *model.OracleDatabaseAgreement,
	req dto.AssociatedLicenseTypeInOracleDbAgreementRequest) error {
	part, err := as.GetOracleDatabaseLicenseType(req.LicenseTypeID)
	if err != nil {
		return err
	}

	associatedLicenseType := model.AssociatedLicenseType{
		ID:              as.NewObjectID(),
		LicenseTypeID:   part.ID,
		ReferenceNumber: req.ReferenceNumber,
		Unlimited:       req.Unlimited,
		Count:           req.Count,
		CatchAll:        req.CatchAll,
		Restricted:      req.Restricted,
		Hosts:           req.Hosts,
	}
	agreement.LicenseTypes = append(agreement.LicenseTypes, associatedLicenseType)

	return nil
}

func checkHosts(as *APIService, hosts []string) error {
	notInClusterHosts, err := as.SearchHosts("hostnames",
		dto.SearchHostsFilters{
			Search:         []string{""},
			OlderThan:      utils.MAX_TIME,
			PageNumber:     -1,
			PageSize:       -1,
			LTEMemoryTotal: -1,
			GTEMemoryTotal: -1,
			LTESwapTotal:   -1,
			GTESwapTotal:   -1,
			LTECPUCores:    -1,
			GTECPUCores:    -1,
			LTECPUThreads:  -1,
			GTECPUThreads:  -1,
		})
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "")
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

		return utils.ErrHostNotFound
	}

	return nil
}

// UpdateAssociatedLicenseTypeOfOracleDbAgreement update associated part in OracleDatabaseAgreement
func (as *APIService) UpdateAssociatedLicenseTypeOfOracleDbAgreement(request dto.AssociatedLicenseTypeInOracleDbAgreementRequest,
) error {
	if err := checkHosts(as, request.Hosts); err != nil {
		return err
	}

	associateLicenseTypeID := utils.Str2oid(request.ID)
	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedLicenseType(associateLicenseTypeID)
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
	req dto.AssociatedLicenseTypeInOracleDbAgreementRequest) error {

	reqID := utils.Str2oid(req.ID)
	associatedLicenseType := agreement.AssociatedLicenseTypeByID(reqID)
	if associatedLicenseType == nil {
		return utils.ErrOracleDatabaseAssociatedPartNotFound
	}

	licenseType, err := as.GetOracleDatabaseLicenseType(req.LicenseTypeID)
	if err != nil {
		return err
	}
	associatedLicenseType.LicenseTypeID = licenseType.ID
	associatedLicenseType.ReferenceNumber = req.ReferenceNumber
	associatedLicenseType.Unlimited = req.Unlimited
	associatedLicenseType.Count = req.Count
	associatedLicenseType.CatchAll = req.CatchAll
	associatedLicenseType.Restricted = req.Restricted
	associatedLicenseType.Hosts = req.Hosts

	return nil
}

// SearchAssociatedLicenseTypesInOracleDatabaseAgreements search OracleDatabase associated parts agreements
func (as *APIService) SearchAssociatedLicenseTypesInOracleDatabaseAgreements(filters dto.SearchOracleDatabaseAgreementsFilter,
) ([]dto.OracleDatabaseAgreementFE, error) {
	agreements, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	hosts, err := as.Database.ListHostUsingOracleDatabaseLicenses()
	if err != nil {
		return nil, err
	}

	err2 := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	if err2 != nil {
		return nil, utils.NewAdvancedErrorPtr(err2, "DB ERROR")
	}

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
	hosts []dto.HostUsingOracleDatabaseLicenses) error {

	licenseTypes, err := as.Database.GetOracleDatabaseLicenseTypes()
	if err != nil {
		return err
	}

	sortOracleDatabaseAgreements(agrs)
	sortHostsByLicenses(hosts)

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Agreements = %s\nHosts= %s\n", utils.ToJSON(agrs), utils.ToJSON(hosts))
	}

	hostsMap := buildHostUsingLicensesMap(hosts)
	licenseTypesMap := buildLicenseTypesMap(licenseTypes)

	fillAgreementsInfo(as, agrs, licenseTypesMap)

	assignAgreementsLicensesToItsAssociatedHosts(as, agrs, hostsMap, licenseTypesMap)

	// sort again and rebuild map because the references are updated during the sort
	sortHostsByLicenses(hosts)
	hostsMap = buildHostUsingLicensesMap(hosts)

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Resorted LicensingObjects: %#v\n", hosts)
	}

	assignLicensesFromCatchAllAgreements(as, agrs, hosts)

	calculateTotalCoveredLicensesAndAvailableCount(as, agrs, hosts, hostsMap, licenseTypesMap)

	return nil
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

func sortHostsByLicenses(obj []dto.HostUsingOracleDatabaseLicenses) {
	sort.Slice(obj, func(i, j int) bool {
		if obj[i].LicenseCount != obj[j].LicenseCount {
			return obj[i].LicenseCount > obj[j].LicenseCount

		} else if obj[i].Name != obj[j].Name {
			return obj[i].Name > obj[j].Name

		} else {
			return obj[i].LicenseTypeID > obj[j].LicenseTypeID
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
		if _, ok := res[host.LicenseTypeID]; !ok {
			res[host.LicenseTypeID] = make(map[string]*dto.HostUsingOracleDatabaseLicenses)
		}
		res[host.LicenseTypeID][host.Name] = &hosts[i]
	}

	return res
}

//TODO  use GetOracleDatabaseAgreementPartsMap ?
func buildLicenseTypesMap(licenseTypes []model.OracleDatabaseLicenseType) map[string]*model.OracleDatabaseLicenseType {
	ltMap := make(map[string]*model.OracleDatabaseLicenseType)

	for i, licenseType := range licenseTypes {
		ltMap[licenseType.ID] = &licenseTypes[i]
	}

	return ltMap
}

func fillAgreementsInfo(as *APIService, agrs []dto.OracleDatabaseAgreementFE, licenseTypes map[string]*model.OracleDatabaseLicenseType) {

	for i := range agrs {
		agr := &agrs[i]

		if licenseType, ok := licenseTypes[agr.LicenseTypeID]; ok {
			agr.ItemDescription = licenseType.ItemDescription
			agr.Metric = licenseType.Metric

			switch agr.Metric {
			case model.LicenseTypeMetricProcessorPerpetual:
				agr.LicensesCount = agr.Count
			case model.LicenseTypeMetricNamedUserPlusPerpetual:
				agr.UsersCount = agr.Count
			}
		} else {
			as.Log.Errorf("Unknown PartID: [%s] in agreement: [%#v]", agr.LicenseTypeID, agr)
		}
	}
}

// Assign available licenses in each agreement to each host associated in each agreement
// if this host is using that kind of license.
func assignAgreementsLicensesToItsAssociatedHosts(
	as *APIService,
	agreements []dto.OracleDatabaseAgreementFE,
	hostsMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses,
	licenseTypes map[string]*model.OracleDatabaseLicenseType) {

	for i := range agreements {
		agreement := &agreements[i]
		sortHostsInAgreementByLicenseCount(agreement, hostsMap, licenseTypes)

		if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
			as.Log.Debugf("Distributing licenses of agreement #%d to host. Agreement = %s\n", i, utils.ToJSON(agreement))
		}

		for j := range agreement.Hosts {
			associatedHost := &agreement.Hosts[j]

			if agreement.AvailableCount <= 0 && !agreement.Unlimited {
				break
			}

			ltID := agreement.LicenseTypeID
			if _, ok := hostsMap[ltID]; !ok {
				// no host use this license
				continue
			}

			hostUsingLicenses, ok := hostsMap[ltID][associatedHost.Hostname]
			if !ok {
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
					ltID)
			}

			if agreement.AvailableCount <= 0 && !agreement.Unlimited {
				break
			}
		}
	}
}

// sortHostsInAgreementByLicenseCount sort the associated hosts by license count
func sortHostsInAgreementByLicenseCount(agr *dto.OracleDatabaseAgreementFE,
	hostsMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses,
	licenseTypes map[string]*model.OracleDatabaseLicenseType) {

	sort.Slice(agr.Hosts, func(i, j int) bool {

		ltID := agr.LicenseTypeID
		mapHostnamesLicenses := hostsMap[ltID]

		iLicenseCount := 0.0
		if iHostUsingLicenses, ok := mapHostnamesLicenses[agr.Hosts[i].Hostname]; ok {
			iLicenseCount = iHostUsingLicenses.LicenseCount
		}

		jLicenseCount := 0.0
		if jHostUsingLicenses, ok := mapHostnamesLicenses[agr.Hosts[j].Hostname]; ok {
			jLicenseCount = jHostUsingLicenses.LicenseCount
		}

		return iLicenseCount > jLicenseCount
	})
}

// Use all the licenses available in agreement to cover host and associatedHost if provided
func doAssignAgreementLicensesToAssociatedHost(
	as *APIService,
	agreement *dto.OracleDatabaseAgreementFE,
	host *dto.HostUsingOracleDatabaseLicenses,
	associatedHost *dto.OracleDatabaseAgreementAssociatedHostFE) {

	switch {
	case agreement.Metric == model.LicenseTypeMetricProcessorPerpetual ||
		agreement.Metric == model.LicenseTypeMetricComputerPerpetual:

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

	case agreement.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual:

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
	hosts []dto.HostUsingOracleDatabaseLicenses) {

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

			if host.LicenseTypeID != agr.LicenseTypeID {
				continue
			}

			doAssignLicenseFromCatchAllAgreement(as, agr, host)

			if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
				as.Log.Debugf("Distributing with metric [%s] [ULA? %t] %f licenses to obj %s. aggCount=%f objCount=0 licenseTypeID=%s\n",
					agr.Metric,
					agr.Unlimited,
					host.LicenseCount,
					host.Name,
					agr.Count,
					agr.LicenseTypeID)
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
	case agreement.Metric == model.LicenseTypeMetricProcessorPerpetual ||
		agreement.Metric == model.LicenseTypeMetricComputerPerpetual:

		var coverableLicenses float64
		if agreement.Unlimited {
			coverableLicenses = hostUsingLicenses.LicenseCount
			agreement.AvailableCount = 0
		} else {
			coverableLicenses = math.Min(agreement.AvailableCount, hostUsingLicenses.LicenseCount)
			agreement.AvailableCount -= coverableLicenses
		}

		hostUsingLicenses.LicenseCount -= coverableLicenses

	case agreement.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual:

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
	licenseTypes map[string]*model.OracleDatabaseLicenseType) {

	licensesCoverStatusByLicenseTypeID := calculateCoverStatusByLicenseType(hosts)

	for i := range agrs {
		agreement := &agrs[i]

		uncoveredLicensesByAssociatedHosts := 0.0
		uncoveredLicensesByAllHosts := 0.0

		ltID := agreement.LicenseTypeID

		for j := range agreement.Hosts {
			associatedHost := &agreement.Hosts[j]
			if _, ok := hostsMap[ltID]; !ok {
				continue
			}

			host, ok := hostsMap[ltID][associatedHost.Hostname]
			if !ok {
				continue
			}

			switch {
			case agreement.Metric == model.LicenseTypeMetricProcessorPerpetual ||
				agreement.Metric == model.LicenseTypeMetricComputerPerpetual:
				associatedHost.TotalCoveredLicensesCount = host.OriginalCount - host.LicenseCount
				associatedHost.ConsumedLicensesCount = host.OriginalCount
				uncoveredLicensesByAssociatedHosts += host.LicenseCount

			case agreement.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual:
				associatedHost.TotalCoveredLicensesCount = (host.OriginalCount - host.LicenseCount) * 25
				associatedHost.ConsumedLicensesCount = host.OriginalCount * 25
				uncoveredLicensesByAssociatedHosts += host.LicenseCount * 25

			default:
				as.Log.Errorf("Unknown metric type: [%s]", agreement.Metric)
			}
		}

		uncoveredLicensesByAllHosts += licensesCoverStatusByLicenseTypeID[ltID].Consumed - licensesCoverStatusByLicenseTypeID[ltID].Covered

		var uncoveredLicenses float64
		if agreement.CatchAll {
			uncoveredLicenses = uncoveredLicensesByAllHosts
		} else {
			uncoveredLicenses = uncoveredLicensesByAssociatedHosts
		}

		if uncoveredLicenses > 0 {
			if (agreement.AvailableCount > 0 && agreement.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual) ||
				(agreement.AvailableCount > 25 && agreement.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual) {

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

func calculateCoverStatusByLicenseType(hosts []dto.HostUsingOracleDatabaseLicenses) map[string]coverStatus {
	licensesStatus := make(map[string]coverStatus)

	for _, host := range hosts {
		licensesStatus[host.LicenseTypeID] = coverStatus{
			Consumed: licensesStatus[host.LicenseTypeID].Consumed + host.OriginalCount,
			Covered:  licensesStatus[host.LicenseTypeID].Covered + (host.OriginalCount - host.LicenseCount),
		}
	}

	return licensesStatus
}

// checkOracleDatabaseAgreementMatchFilter check that agr match the filters
func checkOracleDatabaseAgreementMatchFilter(agr dto.OracleDatabaseAgreementFE, filters dto.SearchOracleDatabaseAgreementsFilter) bool {
	return strings.Contains(strings.ToLower(agr.AgreementID), strings.ToLower(filters.AgreementID)) &&
		strings.Contains(strings.ToLower(agr.LicenseTypeID), strings.ToLower(filters.LicenseTypeID)) &&
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

// DeleteAssociatedLicenseTypeFromOracleDatabaseAgreement delete associated part from OracleDatabaseAgreement
func (as *APIService) DeleteAssociatedLicenseTypeFromOracleDatabaseAgreement(associateLicenseTypeID primitive.ObjectID,
) error {
	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedLicenseType(associateLicenseTypeID)
	if err != nil {
		return err
	}

	if len(agreement.LicenseTypes) <= 1 {
		return as.Database.RemoveOracleDatabaseAgreement(agreement.ID)
	}

	for i := range agreement.LicenseTypes {
		if agreement.LicenseTypes[i].ID == associateLicenseTypeID {
			agreement.LicenseTypes = append(agreement.LicenseTypes[:i], agreement.LicenseTypes[i+1])
			break
		}
	}

	return as.Database.UpdateOracleDatabaseAgreement(*agreement)
}

// AddHostToAssociatedLicenseType add an host to AssociatedLicenseType
func (as *APIService) AddHostToAssociatedLicenseType(associateLicenseTypeID primitive.ObjectID, hostname string,
) error {

	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedLicenseType(associateLicenseTypeID)
	if err != nil {
		return err
	}

	associatedLicenseType := agreement.AssociatedLicenseTypeByID(associateLicenseTypeID)

	for _, host := range associatedLicenseType.Hosts {
		if host == hostname {
			return nil
		}
	}

	if err := checkHosts(as, []string{hostname}); err != nil {
		return err
	}

	associatedLicenseType.Hosts = append(associatedLicenseType.Hosts, hostname)

	return as.Database.UpdateOracleDatabaseAgreement(*agreement)
}

// RemoveHostFromAssociatedLicenseType remove host from AssociatedLicenseType
func (as *APIService) RemoveHostFromAssociatedLicenseType(associateLicenseTypeID primitive.ObjectID, hostname string,
) error {

	agreement, err := as.Database.GetOracleDatabaseAgreementByAssociatedLicenseType(associateLicenseTypeID)
	if err != nil {
		return err
	}

	associatedLicenseType := agreement.AssociatedLicenseTypeByID(associateLicenseTypeID)

	for i, host := range associatedLicenseType.Hosts {
		if host == hostname {
			associatedLicenseType.Hosts = append(associatedLicenseType.Hosts[:i], associatedLicenseType.Hosts[i+1:]...)

			return as.Database.UpdateOracleDatabaseAgreement(*agreement)
		}
	}

	return nil
}
