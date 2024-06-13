// Copyright (c) 2022 Sorint.lab S.p.A.
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

func (as *APIService) AddOracleDatabaseContract(contract model.OracleDatabaseContract) (*dto.OracleDatabaseContractFE, error) {
	if err := checkHosts(as, contract.Hosts); err != nil {
		return nil, err
	}

	if err := checkLicenseTypeIDExists(as, &contract); err != nil {
		return nil, err
	}

	contract.ID = as.NewObjectID()

	err := as.Database.InsertOracleDatabaseContract(contract)
	if err != nil {
		return nil, err
	}

	agrs, err := as.GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	if err != nil {
		return nil, err
	}

	for _, agr := range agrs {
		if agr.ID == contract.ID {
			return &agr, nil
		}
	}

	return nil, utils.NewError(errors.New("Can't find contract which has just been saved"))
}

func checkLicenseTypeIDExists(as *APIService, contract *model.OracleDatabaseContract) error {
	_, err := as.GetOracleDatabaseLicenseType(contract.LicenseTypeID)
	if err != nil {
		return err
	}

	return nil
}

func (as *APIService) UpdateOracleDatabaseContract(contract model.OracleDatabaseContract) (*dto.OracleDatabaseContractFE, error) {
	if err := checkHosts(as, contract.Hosts); err != nil {
		return nil, err
	}

	if err := checkLicenseTypeIDExists(as, &contract); err != nil {
		return nil, err
	}

	if err := as.Database.UpdateOracleDatabaseContract(contract); err != nil {
		return nil, err
	}

	agrs, err := as.GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	if err != nil {
		return nil, err
	}

	for _, agr := range agrs {
		if agr.ID == contract.ID {
			return &agr, nil
		}
	}

	return nil, utils.NewError(errors.New("Can't find contract which has just been saved"))
}

func (as *APIService) GetOracleDatabaseContracts(filter dto.GetOracleDatabaseContractsFilter) ([]dto.OracleDatabaseContractFE, error) {
	if as.mockGetOracleDatabaseContracts != nil {
		return as.mockGetOracleDatabaseContracts(filter)
	}

	contracts, err := as.Database.ListOracleDatabaseContracts()
	if err != nil {
		return nil, err
	}

	usages, err := as.getLicensesUsage()
	if err != nil {
		return nil, err
	}

	if err := as.assignOracleDatabaseContractsToHosts(contracts, usages); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	filteredAgrs := make([]dto.OracleDatabaseContractFE, 0)

	for _, agr := range contracts {
		if checkOracleDatabaseContractMatchFilter(agr, filter) {
			filteredAgrs = append(filteredAgrs, agr)
		}
	}

	return filteredAgrs, nil
}

func (as *APIService) GetOracleDatabaseContractsAsXLSX(filter dto.GetOracleDatabaseContractsFilter) (*excelize.File, error) {
	contracts, err := as.GetOracleDatabaseContracts(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Contracts"
	headers := []string{
		"Contract Number",
		"Part Number",
		"Description",
		"Metric",
		"CSI",
		"Reference Number",
		"Support Expiration",
		"Status",
		"Product Order Date",
		"ULA",
		"Licenses Per Core",
		"Licenses Per User",
		"Available Licenses Core",
		"Available Licenses User",
		"Basket",
		"Restricted",
		"Hostname",
		"Used Licenses",
		"Covered by this contract",
		"Covered by all contracts",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range contracts {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.ContractID)
		sheets.SetCellValue(sheet, nextAxis(), val.LicenseTypeID)
		sheets.SetCellValue(sheet, nextAxis(), val.ItemDescription)
		sheets.SetCellValue(sheet, nextAxis(), val.Metric)
		sheets.SetCellValue(sheet, nextAxis(), val.CSI)
		sheets.SetCellValue(sheet, nextAxis(), val.ReferenceNumber)

		if val.SupportExpiration != nil {
			sheets.SetCellValue(sheet, nextAxis(), val.SupportExpiration)
		} else {
			sheets.SetCellValue(sheet, nextAxis(), "")
		}

		sheets.SetCellValue(sheet, nextAxis(), val.Status)

		if val.ProductOrderDate != nil {
			sheets.SetCellValue(sheet, nextAxis(), val.ProductOrderDate)
		} else {
			sheets.SetCellValue(sheet, nextAxis(), "")
		}

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

// assignOracleDatabaseContractsToHosts assign available licenses in each contracts to hosts using licenses
func (as *APIService) assignOracleDatabaseContractsToHosts(
	agrs []dto.OracleDatabaseContractFE,
	usages []dto.HostUsingOracleDatabaseLicenses) error {
	licenseTypes, err := as.Database.GetOracleDatabaseLicenseTypes()
	if err != nil {
		return err
	}

	sortOracleDatabaseContracts(agrs)
	sortHostsByLicenses(usages)

	if as.Config.APIService.DebugOracleDatabaseContractsAssignmentAlgorithm {
		as.Log.Debugf("Contracts = %s\nHosts= %s\n", utils.ToJSON(agrs), utils.ToJSON(usages))
	}

	usagesMap := buildHostUsingLicensesMap(usages)
	licenseTypesMap := buildLicenseTypesMap(licenseTypes)

	fillContractsInfo(as, agrs, licenseTypesMap)

	err = assignContractsLicensesToItsAssociatedHosts(as, agrs, usagesMap)
	if err != nil {
		return err
	}

	// sort again and rebuild map because the references are updated during the sort
	// sortHostsByLicenses(usages)
	// usagesMap = buildHostUsingLicensesMap(usages)

	if as.Config.APIService.DebugOracleDatabaseContractsAssignmentAlgorithm {
		as.Log.Debugf("Resorted LicensingObjects: %#v\n", usages)
	}

	usages = []dto.HostUsingOracleDatabaseLicenses{}

	for k := range usagesMap {
		for j := range usagesMap[k] {
			usages = append(usages, dto.HostUsingOracleDatabaseLicenses{
				LicenseTypeID: usagesMap[k][j].LicenseTypeID,
				Name:          usagesMap[k][j].Name,
				Type:          usagesMap[k][j].Type,
				LicenseCount:  usagesMap[k][j].LicenseCount,
				OriginalCount: usagesMap[k][j].OriginalCount,
			})
		}
	}

	assignLicensesFromBasketContracts(as, agrs, usages)

	calculateTotalCoveredAndConsumedLicenses(agrs, usagesMap)

	return nil
}

// sortOracleDatabaseContracts sort the list of dto.OracleDatabaseContractsFE
// by Basket (falses first), Unlimited (falses first), decreasing UsersCount, decreasing LicensesCount
func sortOracleDatabaseContracts(obj []dto.OracleDatabaseContractFE) {
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
func buildHostUsingLicensesMap(usages []dto.HostUsingOracleDatabaseLicenses,
) map[string]map[string]*dto.HostUsingOracleDatabaseLicenses {
	res := make(map[string]map[string]*dto.HostUsingOracleDatabaseLicenses)

	for i, usage := range usages {
		if _, ok := res[usage.LicenseTypeID]; !ok {
			res[usage.LicenseTypeID] = make(map[string]*dto.HostUsingOracleDatabaseLicenses)
		}

		res[usage.LicenseTypeID][usage.Name] = &usages[i]
	}

	return res
}

func buildLicenseTypesMap(licenseTypes []model.OracleDatabaseLicenseType) map[string]*model.OracleDatabaseLicenseType {
	ltMap := make(map[string]*model.OracleDatabaseLicenseType)

	for i, licenseType := range licenseTypes {
		ltMap[licenseType.ID] = &licenseTypes[i]
	}

	return ltMap
}

func fillContractsInfo(as *APIService, agrs []dto.OracleDatabaseContractFE, licenseTypes map[string]*model.OracleDatabaseLicenseType) {
	for i := range agrs {
		agr := &agrs[i]

		if licenseType, ok := licenseTypes[agr.LicenseTypeID]; ok {
			agr.ItemDescription = licenseType.ItemDescription
			agr.Metric = licenseType.Metric
		} else {
			as.Log.Errorf("Unknown PartID: [%s] in contract: [%#v]", agr.LicenseTypeID, agr)
		}
	}
}

// Assign available licenses in each contract to each host associated in each contract
// if this host is using that kind of license.
func assignContractsLicensesToItsAssociatedHosts(
	as *APIService,
	contracts []dto.OracleDatabaseContractFE,
	usagesMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses) error {
	hostnamesPerLicense := make(map[string]map[string]bool)

	for i := range contracts {
		contract := &contracts[i]
		sortHostsInContractByLicenseCount(contract, usagesMap)

		if as.Config.APIService.DebugOracleDatabaseContractsAssignmentAlgorithm {
			as.Log.Debugf("Distributing licenses of contract #%d to host. Contract = %s\n", i, utils.ToJSON(contract))
		}

		for j := range contract.Hosts {
			associatedHost := &contract.Hosts[j]

			if !hasAvailableLicenses(contract) {
				associatedHost.TotalCoveredLicensesCount = 0
				associatedHost.CoveredLicensesCount = 0

				continue
			}

			var hostUsingLicenses *dto.HostUsingOracleDatabaseLicenses

			var usages map[string]*dto.HostUsingOracleDatabaseLicenses

			var ok bool

			ltID := contract.LicenseTypeID

			if usages, ok = usagesMap[ltID]; !ok {
				// no host use this license
				continue
			}

			err := as.assignContractsLicensesToHostBelongToCluster(usages, contract, associatedHost, hostnamesPerLicense)
			if err != nil {
				return err
			}

			hostUsingLicenses, ok = usagesMap[ltID][associatedHost.Hostname]
			if !ok {
				// host doesn't use this license
				continue
			}

			if hostUsingLicenses == nil || hostUsingLicenses.LicenseCount <= 0 {
				continue
			}

			doAssignContractLicensesToAssociatedHost(contract, hostUsingLicenses, associatedHost)

			if as.Config.APIService.DebugOracleDatabaseContractsAssignmentAlgorithm {
				as.Log.Debugf(`Distributing %f licenses to host %s. agr.Metrics=%s \
					agr.AvailableLicensesPerCore=%f agr.AvailableLicensesPerUser=%f \
					hostInAgr.CoveredLicensesCount=%f hostUsingLicenses.LicenseCount=%f licenseName=%s\n`,
					hostUsingLicenses.LicenseCount,
					associatedHost.Hostname,
					contract.Metric,
					contract.AvailableLicensesPerCore,
					contract.AvailableLicensesPerUser,
					associatedHost.CoveredLicensesCount,
					hostUsingLicenses.LicenseCount,
					ltID)
			}
		}
	}

	return nil
}

func (as *APIService) assignContractsLicensesToHostBelongToCluster(
	usages map[string]*dto.HostUsingOracleDatabaseLicenses,
	contract *dto.OracleDatabaseContractFE,
	associatedHost *dto.OracleDatabaseContractAssociatedHostFE,
	hostnamesPerLicense map[string]map[string]bool) error {
	for _, usage := range usages {
		if usage == nil || usage.Type != "cluster" || contract.Restricted {
			continue
		}

		cluster, err := as.GetCluster(usage.Name, utils.MAX_TIME)
		if err != nil {
			return err
		}

		for _, hostNameVM := range cluster.VMs {
			if hostNameVM.Hostname == associatedHost.Hostname {
				_, found := hostnamesPerLicense[usage.LicenseTypeID]
				if !found {
					hostnamesPerLicense[usage.LicenseTypeID] = make(map[string]bool)
				}

				alreadyUsed := hostnamesPerLicense[usage.LicenseTypeID][usage.Name]
				if alreadyUsed {
					continue
				}

				hostnamesPerLicense[usage.LicenseTypeID][usage.Name] = true
				contract.CoveredLicenses += usage.LicenseCount

				if contract.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
					contract.AvailableLicensesPerUser -= usage.LicenseCount
				} else {
					contract.AvailableLicensesPerCore -= usage.LicenseCount
				}

				break
			}
		}
	}

	return nil
}

func hasAvailableLicenses(contract *dto.OracleDatabaseContractFE) bool {
	if contract.Unlimited {
		return true
	}

	if contract.AvailableLicensesPerCore > 0 &&
		contract.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual {
		return true
	}

	if contract.AvailableLicensesPerUser > 0 &&
		contract.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
		return true
	}

	return false
}

// sortHostsInContractByLicenseCount sort the associated hosts by license count
func sortHostsInContractByLicenseCount(agr *dto.OracleDatabaseContractFE,
	usagesMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses) {
	sort.Slice(agr.Hosts, func(i, j int) bool {
		ltID := agr.LicenseTypeID
		mapHostnamesLicenses := usagesMap[ltID]

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

// Use all the licenses available in contract to cover host and associatedHost if provided
func doAssignContractLicensesToAssociatedHost(
	contract *dto.OracleDatabaseContractFE,
	host *dto.HostUsingOracleDatabaseLicenses,
	associatedHost *dto.OracleDatabaseContractAssociatedHostFE) {
	if contract.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual {
		var coverableLicenses float64
		if contract.Unlimited {
			coverableLicenses = host.LicenseCount
			contract.AvailableLicensesPerCore = 0
		} else {
			coverableLicenses = math.Min(contract.AvailableLicensesPerCore, host.LicenseCount)
			contract.AvailableLicensesPerCore -= coverableLicenses
		}

		associatedHost.CoveredLicensesCount += coverableLicenses
		contract.CoveredLicenses += coverableLicenses
		host.LicenseCount -= coverableLicenses
	} else {
		var coverableLicenses float64
		if contract.Unlimited {
			coverableLicenses = host.LicenseCount
			contract.AvailableLicensesPerUser = 0
		} else {
			availableInContract := math.Floor(contract.AvailableLicensesPerUser/model.FactorNamedUser) * model.FactorNamedUser
			coverableLicenses = math.Min(availableInContract, host.LicenseCount)
			contract.AvailableLicensesPerUser -= coverableLicenses
		}

		associatedHost.CoveredLicensesCount += coverableLicenses
		contract.CoveredLicenses += coverableLicenses
		host.LicenseCount -= coverableLicenses
	}
}

// If an contract is basket distributes its licenses to every hosts that use that kind of license
func assignLicensesFromBasketContracts(
	as *APIService,
	agrs []dto.OracleDatabaseContractFE,
	usages []dto.HostUsingOracleDatabaseLicenses) {
	for i := range usages {
		usage := &usages[i]

		if usage.LicenseCount <= 0 {
			continue
		}

		for j := range agrs {
			agr := &agrs[j]

			if !agr.Basket {
				continue
			}

			if usage.LicenseTypeID != agr.LicenseTypeID {
				continue
			}

			if !hasAvailableLicenses(agr) {
				continue
			}

			doAssignLicenseFromBasketContract(agr, usage)

			if as.Config.APIService.DebugOracleDatabaseContractsAssignmentAlgorithm {
				as.Log.Debugf("Distributing with metric [%s] [ULA? %t] %f licenses to obj %s. objCount=0 licenseTypeID=%s\n",
					agr.Metric,
					agr.Unlimited,
					usage.LicenseCount,
					usage.Name,
					agr.LicenseTypeID)
			}
		}
	}

	if as.Config.APIService.DebugOracleDatabaseContractsAssignmentAlgorithm {
		as.Log.Debugf("Associations finished. LicensingObjects: %#v\n", usages)
	}
}

// Use all the licenses available in contract to cover host and associatedHost if provided
func doAssignLicenseFromBasketContract(
	contract *dto.OracleDatabaseContractFE,
	usage *dto.HostUsingOracleDatabaseLicenses) {
	var coverableLicenses float64

	if contract.Metric != model.LicenseTypeMetricNamedUserPlusPerpetual {
		if contract.Unlimited {
			coverableLicenses = usage.LicenseCount
			contract.AvailableLicensesPerCore = 0
		} else {
			coverableLicenses = math.Min(contract.AvailableLicensesPerCore, usage.LicenseCount)
			contract.AvailableLicensesPerCore -= coverableLicenses
		}
	} else {
		if contract.Unlimited {
			coverableLicenses = usage.LicenseCount
			contract.AvailableLicensesPerUser = 0
		} else {
			coverableLicenses = math.Floor(math.Min(contract.AvailableLicensesPerUser, usage.LicenseCount)/model.FactorNamedUser) * model.FactorNamedUser
			contract.AvailableLicensesPerUser -= coverableLicenses
		}
	}

	contract.CoveredLicenses += coverableLicenses
	usage.LicenseCount -= coverableLicenses
}

func calculateTotalCoveredAndConsumedLicenses(
	agrs []dto.OracleDatabaseContractFE,
	usagesMap map[string]map[string]*dto.HostUsingOracleDatabaseLicenses) {
	for i := range agrs {
		contract := &agrs[i]

		ltID := contract.LicenseTypeID

		for j := range contract.Hosts {
			associatedHost := &contract.Hosts[j]

			if _, ok := usagesMap[ltID]; !ok {
				continue
			}

			host, ok := usagesMap[ltID][associatedHost.Hostname]
			if !ok {
				continue
			}

			associatedHost.TotalCoveredLicensesCount = host.OriginalCount - host.LicenseCount
			associatedHost.ConsumedLicensesCount = host.OriginalCount
		}
	}
}

// checkOracleDatabaseContractMatchFilter check that agr match the filters
func checkOracleDatabaseContractMatchFilter(agr dto.OracleDatabaseContractFE, filters dto.GetOracleDatabaseContractsFilter) bool {
	return strings.Contains(strings.ToLower(agr.ContractID), strings.ToLower(filters.ContractID)) &&
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

func (as *APIService) DeleteOracleDatabaseContract(id primitive.ObjectID) error {
	return as.Database.RemoveOracleDatabaseContract(id)
}

func (as *APIService) AddHostToOracleDatabaseContract(id primitive.ObjectID, hostname string) error {
	if err := checkHosts(as, []string{hostname}); err != nil {
		return err
	}

	contract, err := as.Database.GetOracleDatabaseContract(id)
	if err != nil {
		return err
	}

	for _, host := range contract.Hosts {
		if host == hostname {
			return nil
		}
	}

	contract.Hosts = append(contract.Hosts, hostname)

	return as.Database.UpdateOracleDatabaseContract(*contract)
}

func (as *APIService) DeleteHostFromOracleDatabaseContract(id primitive.ObjectID, hostname string) error {
	if err := checkHosts(as, []string{hostname}); err != nil {
		return err
	}

	contract, err := as.Database.GetOracleDatabaseContract(id)
	if err != nil {
		return err
	}

	for i := range contract.Hosts {
		host := contract.Hosts[i]
		if host == hostname {
			contract.Hosts = append(
				contract.Hosts[0:i],
				contract.Hosts[i+1:len(contract.Hosts)]...)

			break
		}
	}

	return as.Database.UpdateOracleDatabaseContract(*contract)
}

func (as *APIService) DeleteHostFromOracleDatabaseContracts(hostname string) error {
	if err := checkHosts(as, []string{hostname}); err != nil {
		return err
	}

	listContracts, err := as.Database.ListOracleDatabaseContracts()
	if err != nil {
		return err
	}

	for _, la := range listContracts {
		for i := range la.Hosts {
			host := la.Hosts[i]
			if host.Hostname == hostname {
				contract, err := as.Database.GetOracleDatabaseContract(la.ID)
				if err != nil {
					return err
				}

				contract.Hosts = append(
					contract.Hosts[0:i],
					contract.Hosts[i+1:len(contract.Hosts)]...)

				errContract := as.Database.UpdateOracleDatabaseContract(*contract)
				if errContract != nil {
					return err
				}
			}
		}
	}

	return nil
}
