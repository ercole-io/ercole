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
	"errors"
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

func (as *APIService) AddOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) (*dto.OracleDatabaseAgreementFE, error) {
	if err := checkHosts(as, agreement.Hosts); err != nil {
		return nil, err
	}

	if err := checkLicenseTypeIDExists(as, &agreement); err != nil {
		return nil, err
	}

	agreement.ID = as.NewObjectID()
	err := as.Database.InsertOracleDatabaseAgreement(agreement)
	if err != nil {
		return nil, err
	}

	agrs, err := as.GetOracleDatabaseAgreements(dto.NewGetOracleDatabaseAgreementsFilter())
	if err != nil {
		return nil, err
	}
	for _, agr := range agrs {
		if agr.ID == agreement.ID {
			return &agr, nil
		}
	}

	return nil, utils.NewError(errors.New("Can't find agreement which has just been saved"))
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
		return utils.NewError(err, "")
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

func checkLicenseTypeIDExists(as *APIService, agreement *model.OracleDatabaseAgreement) error {
	_, err := as.GetOracleDatabaseLicenseType(agreement.LicenseTypeID)
	if err != nil {
		return err
	}

	return nil
}

func (as *APIService) UpdateOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) (*dto.OracleDatabaseAgreementFE, error) {
	if err := checkHosts(as, agreement.Hosts); err != nil {
		return nil, err
	}

	if err := checkLicenseTypeIDExists(as, &agreement); err != nil {
		return nil, err
	}

	if err := as.Database.UpdateOracleDatabaseAgreement(agreement); err != nil {
		return nil, err
	}

	agrs, err := as.GetOracleDatabaseAgreements(dto.NewGetOracleDatabaseAgreementsFilter())
	if err != nil {
		return nil, err
	}
	for _, agr := range agrs {
		if agr.ID == agreement.ID {
			return &agr, nil
		}
	}

	return nil, utils.NewError(errors.New("Can't find agreement which has just been saved"))
}

func (as *APIService) GetOracleDatabaseAgreements(filter dto.GetOracleDatabaseAgreementsFilter) ([]dto.OracleDatabaseAgreementFE, error) {
	if as.mockGetOracleDatabaseAgreements != nil {
		return as.mockGetOracleDatabaseAgreements(filter)
	}

	agreements, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	hosts, err := as.Database.ListHostUsingOracleDatabaseLicenses()
	if err != nil {
		return nil, err
	}

	if err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	filteredAgrs := make([]dto.OracleDatabaseAgreementFE, 0)
	for _, agr := range agreements {
		if checkOracleDatabaseAgreementMatchFilter(agr, filter) {
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
func checkOracleDatabaseAgreementMatchFilter(agr dto.OracleDatabaseAgreementFE, filters dto.GetOracleDatabaseAgreementsFilter) bool {
	return strings.Contains(strings.ToLower(agr.AgreementID), strings.ToLower(filters.AgreementID)) &&
		strings.Contains(strings.ToLower(agr.LicenseTypeID), strings.ToLower(filters.LicenseTypeID)) &&
		strings.Contains(strings.ToLower(agr.ItemDescription), strings.ToLower(filters.ItemDescription)) &&
		strings.Contains(strings.ToLower(agr.CSI), strings.ToLower(filters.CSI)) &&
		(filters.Metric == "" || strings.EqualFold(agr.Metric, filters.Metric)) &&
		strings.Contains(strings.ToLower(agr.ReferenceNumber), strings.ToLower(filters.ReferenceNumber)) &&
		(filters.Unlimited == "" || agr.Unlimited == (filters.Unlimited == "true")) &&
		(filters.CatchAll == "" || agr.CatchAll == (filters.CatchAll == "true")) &&
		(filters.LicensesCountLTE == -1 || agr.LicensesCount <= float64(filters.LicensesCountLTE)) &&
		(filters.LicensesCountGTE == -1 || agr.LicensesCount >= float64(filters.LicensesCountGTE)) &&
		(filters.UsersCountLTE == -1 || agr.UsersCount <= float64(filters.UsersCountLTE)) &&
		(filters.UsersCountGTE == -1 || agr.UsersCount >= float64(filters.UsersCountGTE)) &&
		(filters.AvailableCountLTE == -1 || agr.AvailableCount <= float64(filters.AvailableCountLTE)) &&
		(filters.AvailableCountGTE == -1 || agr.AvailableCount >= float64(filters.AvailableCountGTE))
}

func (as *APIService) DeleteOracleDatabaseAgreement(id primitive.ObjectID) error {
	return as.Database.RemoveOracleDatabaseAgreement(id)
}

func (as *APIService) AddHostToOracleDatabaseAgreement(id primitive.ObjectID, hostname string) error {
	if err := checkHosts(as, []string{hostname}); err != nil {
		return err
	}

	agreement, err := as.Database.GetOracleDatabaseAgreement(id)
	if err != nil {
		return err
	}

	for _, host := range agreement.Hosts {
		if host == hostname {
			return nil
		}
	}

	agreement.Hosts = append(agreement.Hosts, hostname)

	return as.Database.UpdateOracleDatabaseAgreement(*agreement)
}

func (as *APIService) DeleteHostFromOracleDatabaseAgreement(id primitive.ObjectID, hostname string) error {
	if err := checkHosts(as, []string{hostname}); err != nil {
		return err
	}

	agreement, err := as.Database.GetOracleDatabaseAgreement(id)
	if err != nil {
		return err
	}

	for i := range agreement.Hosts {
		host := agreement.Hosts[i]
		if host == hostname {
			agreement.Hosts = append(
				agreement.Hosts[0:i],
				agreement.Hosts[i+1:len(agreement.Hosts)]...)
			break
		}
	}

	return as.Database.UpdateOracleDatabaseAgreement(*agreement)
}
