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

// Package service is a package that provides methods for querying data

package service

import (
	"errors"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) SearchSqlServerInstances(f dto.SearchSqlServerInstancesFilter) (*dto.SqlServerInstanceResponse, error) {
	return as.Database.SearchSqlServerInstances(strings.Split(f.Search, " "), f.SortBy, f.SortDesc,
		f.PageNumber, f.PageSize, f.Location, f.Environment, f.OlderThan)
}

func (as *APIService) SearchSqlServerInstancesAsXLSX(filter dto.SearchSqlServerInstancesFilter) (*excelize.File, error) {
	instances, err := as.Database.SearchSqlServerInstances(strings.Split(filter.Search, " "),
		filter.SortBy, filter.SortDesc,
		-1, -1,
		filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	sheet := "Instances"
	headers := []string{
		"Hostname",
		"Location",
		"Name",
		"Status",
		"Edition",
		"CollationName",
		"Version",
	}

	file, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)
	for _, val := range instances.Content {
		nextAxis := axisHelp.NewRow()

		file.SetCellValue(sheet, nextAxis(), val.Hostname)
		file.SetCellValue(sheet, nextAxis(), val.Location)
		file.SetCellValue(sheet, nextAxis(), val.Name)
		file.SetCellValue(sheet, nextAxis(), val.Status)
		file.SetCellValue(sheet, nextAxis(), val.Edition)
		file.SetCellValue(sheet, nextAxis(), val.CollationName)
		file.SetCellValue(sheet, nextAxis(), val.Version)
	}

	return file, nil
}

func (as *APIService) GetSqlServerUsedLicenses(hostname string, filter dto.GlobalFilter) (*dto.SqlServerDatabaseUsedLicenseSearchResponse, error) {
	usedLicenses, err := as.Database.SearchSqlServerDatabaseUsedLicenses(hostname, "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	any := dto.GlobalFilter{Location: "", Environment: "", OlderThan: utils.MAX_TIME}

	clustersList, err := as.Database.GetClusters(any)
	if err != nil {
		return nil, err
	}

	clusters := make(map[string]dto.Cluster, len(clustersList))
	for _, cluster := range clustersList {
		clusters[cluster.Name] = cluster
	}

	contracts, err := as.Database.ListSqlServerDatabaseContracts()
	if err != nil {
		return nil, err
	}

	hostWithCluster := make(map[string]string, len(contracts))

	for _, contract := range contracts {
		if contract.Type == model.SqlServerContractTypeCluster {
			for _, clusterName := range contract.Clusters {
				cluster, err := as.Database.GetCluster(clusterName, utils.MAX_TIME)
				if err != nil {
					continue
				}

				hostWithCluster[cluster.Hostname] = clusterName
				for _, vm := range clusters[clusterName].VMs {
					hostWithCluster[vm.Hostname] = clusterName
				}
			}
		}
	}

	for i := range usedLicenses.Content {
		usedLicense := &usedLicenses.Content[i]
		_, ok := hostWithCluster[usedLicense.Hostname]

		if ok {
			usedLicense.ContractType = model.SqlServerContractTypeCluster
			usedLicense.Clustername = hostWithCluster[usedLicense.Hostname]

			continue
		}

		usedLicense.ContractType = model.SqlServerContractTypeHost
	}

	return usedLicenses, nil
}

func (as *APIService) GetSqlServerDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {
	licenses := make(map[string]*dto.LicenseCompliance)
	purchasedContracts := make(map[string]int)

	any := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses, err := as.GetSqlServerUsedLicenses("", any)
	if err != nil {
		return nil, err
	}

	if len(usedLicenses.Content) == 0 {
		return []dto.LicenseCompliance{}, nil
	}

	contracts, err := as.Database.ListSqlServerDatabaseContracts()
	if err != nil {
		return nil, err
	}

	lts, err := as.GetSqlServerDatabaseLicenseTypesAsMap()
	if err != nil {
		return nil, err
	}

	for _, usedLicense := range usedLicenses.Content {
		license, ok := licenses[usedLicense.LicenseTypeID]

		var errC error

		if !ok {
			license, errC = getNewSqlServerLicenseCompliance(lts, usedLicense)
			if errC != nil {
				as.Log.Errorf(errC.Error())
				continue
			}

			licenses[usedLicense.LicenseTypeID] = license
		}

		var isContract bool

		var isInCluster bool

		for _, contract := range contracts {
			if contract.LicenseTypeID == usedLicense.LicenseTypeID && contract.Type == usedLicense.ContractType {
				isContract = true
				_, ok := purchasedContracts[contract.ContractID]

				if !ok {
					purchasedContracts[contract.ContractID] = contract.LicensesNumber
					license.Purchased = float64(contract.LicensesNumber)

					if contract.Type == model.SqlServerContractTypeCluster {
						cluster, err := as.GetCluster(usedLicense.Clustername, utils.MAX_TIME)
						if err != nil {
							continue
						}

						if usedLicense.Clustername != "" {
							for _, clusterContract := range contract.Clusters {
								if clusterContract == cluster.Name {
									ltHost := lts[contract.LicenseTypeID]
									ltCluster := lts[usedLicense.LicenseTypeID]

									if ltCluster.Version == "STD" && ltHost.Version == "ENT" {
										isInCluster = false
									} else {
										isInCluster = true
									}

									break
								}
							}
						}

						license.Consumed = float64(cluster.CPU)
					}
				}

				if contract.Type == model.SqlServerContractTypeHost {
					if !isInCluster {
						license.Consumed += usedLicense.UsedLicenses
					}
				}
			}
		}

		if usedLicense.ContractType == model.SqlServerContractTypeHost && !isContract {
			license.Consumed += usedLicense.UsedLicenses
		}
	}

	result := make([]dto.LicenseCompliance, 0, len(licenses))

	for _, license := range licenses {
		if license.Purchased != 0 {
			if license.Purchased >= license.Consumed {
				license.Covered = license.Consumed
			} else {
				license.Covered = license.Purchased
			}

			license.Available = license.Purchased - license.Covered
		}

		if license.Consumed == 0 {
			license.Compliance = 1
		} else {
			license.Compliance = license.Covered / license.Consumed
		}

		result = append(result, *license)
	}

	return result, nil
}

func getNewSqlServerLicenseCompliance(licenseTypes map[string]model.SqlServerDatabaseLicenseType, usedLicense dto.SqlServerDatabaseUsedLicense) (*dto.LicenseCompliance, error) {
	lt := licenseTypes[usedLicense.LicenseTypeID]

	var licenseCompliance dto.LicenseCompliance

	licenseCompliance.LicenseTypeID = lt.ID
	licenseCompliance.ItemDescription = lt.ItemDescription

	if usedLicense.ContractType == model.SqlServerContractTypeHost {
		licenseCompliance.Metric = model.SqlServerContractTypeHost
	} else if usedLicense.ContractType == model.SqlServerContractTypeCluster {
		licenseCompliance.Metric = model.SqlServerContractTypeCluster
	} else {
		return nil, errors.New("Unknown SqlServerContractType: " + usedLicense.ContractType)
	}

	return &licenseCompliance, nil
}
