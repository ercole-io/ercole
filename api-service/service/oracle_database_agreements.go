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

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ercole-io/ercole/v2/utils/exutils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
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
	notInClusterHosts, err := as.SearchHosts("hostnames", //TODO Why not in cluster hosts?
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

func (as *APIService) GetOracleDatabaseAgreementsAsXLSX(filter dto.GetOracleDatabaseAgreementsFilter) (*excelize.File, error) {
	agreements, err := as.GetOracleDatabaseAgreements(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Agreements"
	headers := []string{
		"Agreement Number",
		"Part Number",
		"Description",
		"Metric",
		"CSI",
		"Reference Number",
		"ULA",
		"Licenses Per Core",
		"Licenses Per User",
		"Available Licenses Core",
		"Available Licenses User",
		"Basket",
		"Restricted",

		"Hostname",
		"Used Licenses",
		"Covered by this agreement",
		"Covered by all agreements",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}
	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range agreements {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.AgreementID)
		sheets.SetCellValue(sheet, nextAxis(), val.LicenseTypeID)
		sheets.SetCellValue(sheet, nextAxis(), val.ItemDescription)
		sheets.SetCellValue(sheet, nextAxis(), val.Metric)
		sheets.SetCellValue(sheet, nextAxis(), val.CSI)
		sheets.SetCellValue(sheet, nextAxis(), val.ReferenceNumber)
		sheets.SetCellValue(sheet, nextAxis(), val.Unlimited)
		sheets.SetCellValue(sheet, nextAxis(), val.LicensesPerCore)
		sheets.SetCellValue(sheet, nextAxis(), val.LicensesPerUser)
		sheets.SetCellValue(sheet, nextAxis(), val.AvailableLicensesPerCore)
		sheets.SetCellValue(sheet, nextAxis(), val.AvailableLicensesPerUser)
		sheets.SetCellValue(sheet, nextAxis(), val.Basket)
		sheets.SetCellValue(sheet, nextAxis(), val.Restricted)

		for _, val2 := range val.Hosts {
			sheets.DuplicateRow(sheet, axisHelp.GetIndexRow())
			duplicateRowNextAxis := axisHelp.NewRowSincePreviousColumn()

			sheets.SetCellValue(sheet, duplicateRowNextAxis(), val2.Hostname)
			sheets.SetCellValue(sheet, duplicateRowNextAxis(), val2.ConsumedLicensesCount)
			sheets.SetCellValue(sheet, duplicateRowNextAxis(), val2.CoveredLicensesCount)
			sheets.SetCellValue(sheet, duplicateRowNextAxis(), val2.TotalCoveredLicensesCount)
		}
	}
	return sheets, err
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

	assignLicensesFromBasketAgreements(as, agrs, hosts)

	calculateTotalCoveredAndConsumedLicenses(as, agrs, hostsMap)

	return nil
}

// sortOracleDatabaseAgreements sort the list of dto.OracleDatabaseAgreementsFE
// by Basket (falses first), Unlimited (falses first), decreasing UsersCount, decreasing LicensesCount
func sortOracleDatabaseAgreements(obj []dto.OracleDatabaseAgreementFE) {
	sort.Slice(obj, func(i, j int) bool {

		if obj[i].Basket != obj[j].Basket {
			return obj[j].Basket

		} else if obj[i].Unlimited != obj[j].Unlimited {
			return obj[j].Unlimited

		} else if obj[i].LicensesPerUser != obj[j].LicensesPerUser {
			return obj[i].LicensesPerUser > obj[j].LicensesPerUser

		} else {
			return obj[i].LicensesPerCore > obj[j].LicensesPerCore
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

			if !hasAvailableLicenses(agreement) {
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
				as.Log.Debugf(`Distributing %f licenses to host %s. agr.Metrics=%s \
					agr.AvailableLicensesPerCore=%f agr.AvailableLicensesPerUser=%f \
					hostInAgr.CoveredLicensesCount=%f hostUsingLicenses.LicenseCount=%f licenseName=%s\n`,
					hostUsingLicenses.LicenseCount,
					associatedHost.Hostname,
					agreement.Metric,
					agreement.AvailableLicensesPerCore,
					agreement.AvailableLicensesPerUser,
					associatedHost.CoveredLicensesCount,
					hostUsingLicenses.LicenseCount,
					ltID)
			}
		}
	}
}

func hasAvailableLicenses(agreement *dto.OracleDatabaseAgreementFE) bool {
	if agreement.Unlimited {
		return true
	}

	if agreement.AvailableLicensesPerCore > 0 &&
		agreement.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual {
		return true

	}

	if agreement.AvailableLicensesPerUser > 0 &&
		agreement.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
		return true
	}

	return false
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

	if agreement.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual {
		var coverableLicenses float64
		if agreement.Unlimited {
			coverableLicenses = host.LicenseCount
			agreement.AvailableLicensesPerCore = 0
		} else {
			coverableLicenses = math.Min(agreement.AvailableLicensesPerCore, host.LicenseCount)
			agreement.AvailableLicensesPerCore -= coverableLicenses
		}

		associatedHost.CoveredLicensesCount += coverableLicenses
		agreement.CoveredLicenses += coverableLicenses
		host.LicenseCount -= coverableLicenses

	} else {
		var coverableLicenses float64

		if agreement.Unlimited {
			coverableLicenses = host.LicenseCount
			agreement.AvailableLicensesPerUser = 0
		} else {
			coverableLicenses = math.Floor(math.Min(agreement.AvailableLicensesPerUser, host.LicenseCount*25) / 25)
			agreement.AvailableLicensesPerUser -= coverableLicenses * 25
		}

		associatedHost.CoveredLicensesCount += coverableLicenses * 25
		agreement.CoveredLicenses += coverableLicenses * 25
		host.LicenseCount -= coverableLicenses
	}
}

// If an agreement is basket distributes its licenses to every hosts that use that kind of license
func assignLicensesFromBasketAgreements(
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

			if !agr.Basket {
				continue
			}

			if host.LicenseTypeID != agr.LicenseTypeID {
				continue
			}

			if !hasAvailableLicenses(agr) {
				continue
			}

			doAssignLicenseFromBasketAgreement(as, agr, host)

			if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
				as.Log.Debugf("Distributing with metric [%s] [ULA? %t] %f licenses to obj %s. objCount=0 licenseTypeID=%s\n",
					agr.Metric,
					agr.Unlimited,
					host.LicenseCount,
					host.Name,
					agr.LicenseTypeID)
			}
		}
	}

	if as.Config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		as.Log.Debugf("Associations finished. LicensingObjects: %#v\n", hosts)
	}
}

// Use all the licenses available in agreement to cover host and associatedHost if provided
func doAssignLicenseFromBasketAgreement(
	as *APIService,
	agreement *dto.OracleDatabaseAgreementFE,
	hostUsingLicenses *dto.HostUsingOracleDatabaseLicenses) {

	var coverableLicenses float64

	if agreement.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual {
		if agreement.Unlimited {
			coverableLicenses = hostUsingLicenses.LicenseCount
			agreement.AvailableLicensesPerCore = 0
		} else {
			coverableLicenses = math.Min(agreement.AvailableLicensesPerCore, hostUsingLicenses.LicenseCount)
			agreement.AvailableLicensesPerCore -= coverableLicenses
		}
		agreement.CoveredLicenses += coverableLicenses

	} else {
		if agreement.Unlimited {
			coverableLicenses = hostUsingLicenses.LicenseCount
			agreement.AvailableLicensesPerUser = 0
		} else {
			coverableLicenses = math.Floor(math.Min(agreement.AvailableLicensesPerUser, hostUsingLicenses.LicenseCount*25) / 25)
			agreement.AvailableLicensesPerUser -= coverableLicenses * 25
		}
		agreement.CoveredLicenses += coverableLicenses * 25

	}

	hostUsingLicenses.LicenseCount -= coverableLicenses
}

func calculateTotalCoveredAndConsumedLicenses(
	as *APIService,
	agrs []dto.OracleDatabaseAgreementFE,
	hostsMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses) {

	for i := range agrs {
		agreement := &agrs[i]

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

			if agreement.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual {
				associatedHost.TotalCoveredLicensesCount = host.OriginalCount - host.LicenseCount
				associatedHost.ConsumedLicensesCount = host.OriginalCount
			} else {
				associatedHost.TotalCoveredLicensesCount = (host.OriginalCount - host.LicenseCount) * 25
				associatedHost.ConsumedLicensesCount = host.OriginalCount * 25
			}
		}
	}
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
		(filters.Basket == "" || agr.Basket == (filters.Basket == "true")) &&
		(filters.LicensesPerCoreLTE == -1 || agr.LicensesPerCore <= float64(filters.LicensesPerCoreLTE)) &&
		(filters.LicensesPerCoreGTE == -1 || agr.LicensesPerCore >= float64(filters.LicensesPerCoreGTE)) &&
		(filters.LicensesPerUserLTE == -1 || agr.LicensesPerUser <= float64(filters.LicensesPerUserLTE)) &&
		(filters.LicensesPerUserGTE == -1 || agr.LicensesPerUser >= float64(filters.LicensesPerUserGTE)) &&
		(filters.AvailableLicensesPerCoreLTE == -1 || agr.AvailableLicensesPerCore <= float64(filters.AvailableLicensesPerCoreLTE)) &&
		(filters.AvailableLicensesPerCoreGTE == -1 || agr.AvailableLicensesPerCore >= float64(filters.AvailableLicensesPerCoreGTE)) &&
		(filters.AvailableLicensesPerUserLTE == -1 || agr.AvailableLicensesPerUser <= float64(filters.AvailableLicensesPerCoreLTE)) &&
		(filters.AvailableLicensesPerUserGTE == -1 || agr.AvailableLicensesPerUser >= float64(filters.AvailableLicensesPerCoreGTE))
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
