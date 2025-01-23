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

func (as *APIService) SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error) {
	instances, err := as.Database.SearchMySQLInstances(filter)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (as *APIService) SearchMySQLInstancesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	instances, err := as.Database.SearchMySQLInstances(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Instances"
	headers := []string{
		"Name",
		"Location",
		"Version",
		"Edition",
		"Platform",
		"Architecture",
		"Engine",
		"RedoLogEnabled",
		"Charset Server",
		"Charset System",
		"PageSize",
		"Threads Concurrency",
		"BufferPool Size",
		"LogBuffer Size",
		"SortBuffer Size",
		"ReadOnly",
		"Databases",
		"Table Schemas",
	}

	file, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)
	for _, val := range instances {
		nextAxis := axisHelp.NewRow()

		file.SetCellValue(sheet, nextAxis(), val.Name)
		file.SetCellValue(sheet, nextAxis(), val.Location)
		file.SetCellValue(sheet, nextAxis(), val.Version)
		file.SetCellValue(sheet, nextAxis(), val.Edition)
		file.SetCellValue(sheet, nextAxis(), val.Platform)
		file.SetCellValue(sheet, nextAxis(), val.Architecture)
		file.SetCellValue(sheet, nextAxis(), val.Engine)
		file.SetCellValue(sheet, nextAxis(), val.RedoLogEnabled)
		file.SetCellValue(sheet, nextAxis(), val.CharsetServer)
		file.SetCellValue(sheet, nextAxis(), val.CharsetSystem)
		file.SetCellValue(sheet, nextAxis(), val.PageSize)
		file.SetCellValue(sheet, nextAxis(), val.ThreadsConcurrency)
		file.SetCellValue(sheet, nextAxis(), val.BufferPoolSize)
		file.SetCellValue(sheet, nextAxis(), val.LogBufferSize)
		file.SetCellValue(sheet, nextAxis(), val.SortBufferSize)
		file.SetCellValue(sheet, nextAxis(), val.ReadOnly)

		databases := make([]string, len(val.Databases))
		for i := range val.Databases {
			databases[i] = val.Databases[i].Name
		}

		file.SetCellValue(sheet, nextAxis(), strings.Join(databases, ", "))

		tableSchemas := make([]string, len(val.TableSchemas))
		for i := range val.TableSchemas {
			tableSchemas[i] = val.TableSchemas[i].Name
		}

		file.SetCellValue(sheet, nextAxis(), strings.Join(tableSchemas, ", "))
	}

	return file, nil
}

func (as *APIService) GetMySQLUsedLicenses(hostname string, filter dto.GlobalFilter) ([]dto.MySQLUsedLicense, error) {
	usedLicenses, err := as.Database.GetMySQLUsedLicenses(hostname, filter)
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
		clusters[cluster.Hostname] = cluster
	}

	contracts, err := as.Database.GetMySQLContracts([]string{filter.Location})
	if err != nil {
		return nil, err
	}

	hostWithCluster := make(map[string]string, len(contracts))

	for _, contract := range contracts {
		if contract.Type == model.MySQLContractTypeCluster {
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

	for i := range usedLicenses {
		usedLicense := &usedLicenses[i]
		_, ok := hostWithCluster[usedLicense.Hostname]

		if ok {
			usedLicense.ContractType = model.MySQLContractTypeCluster
			usedLicense.Clustername = hostWithCluster[usedLicense.Hostname]

			continue
		}

		usedLicense.ContractType = model.MySQLContractTypeHost
	}

	return usedLicenses, nil
}

func (as *APIService) GetMySQLDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {
	licenses := make(map[string]*dto.LicenseCompliance)
	purchasedContracts := make(map[string]int)

	any := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses, err := as.GetMySQLUsedLicenses("", any)
	if err != nil {
		return nil, err
	}

	if len(usedLicenses) == 0 {
		return []dto.LicenseCompliance{}, nil
	}

	contracts, err := as.Database.GetMySQLContracts([]string{})
	if err != nil {
		return nil, err
	}

	for _, usedLicense := range usedLicenses {
		license, ok := licenses[usedLicense.LicenseTypeID]

		var errC error

		if !ok {
			license, errC = getNewMySqlLicenseCompliance(usedLicense)
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
					purchasedContracts[contract.ContractID] = int(contract.NumberOfLicenses)
					license.Purchased = float64(contract.NumberOfLicenses)

					if contract.Type == model.MySQLContractTypeCluster {
						cluster, err := as.GetCluster(usedLicense.Clustername, utils.MAX_TIME)
						if err != nil {
							continue
						}

						if usedLicense.Clustername != "" {
							for _, clusterContract := range contract.Clusters {
								if clusterContract == cluster.Name {
									isInCluster = true

									break
								}
							}
						}

						license.Consumed = float64(cluster.CPU)
					}
				}

				if contract.Type == model.MySQLContractTypeHost {
					if !isInCluster {
						license.Consumed += usedLicense.UsedLicenses
					}
				}
			}
		}

		if usedLicense.ContractType == model.MySQLContractTypeHost && !isContract {
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

func getNewMySqlLicenseCompliance(usedLicense dto.MySQLUsedLicense) (*dto.LicenseCompliance, error) {
	var licenseCompliance dto.LicenseCompliance

	licenseCompliance.LicenseTypeID = model.MySqlPartNumber
	licenseCompliance.ItemDescription = model.MySqlItemDescription

	if usedLicense.ContractType == model.MySQLContractTypeHost {
		licenseCompliance.Metric = model.MySQLContractTypeHost
	} else if usedLicense.ContractType == model.MySQLContractTypeCluster {
		licenseCompliance.Metric = model.MySQLContractTypeCluster
	} else {
		return nil, errors.New("Unknown MySqlContractType: " + usedLicense.ContractType)
	}

	return &licenseCompliance, nil
}
