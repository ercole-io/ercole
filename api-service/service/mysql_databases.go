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

// Package service is a package that provides methods for querying data
package service

import (
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
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

	file, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_generic.xlsx")
	if err != nil {
		return nil, err
	}

	sheet := "Instances"
	file.SetSheetName("Sheet1", sheet)
	headers := []string{
		"Name",
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

	for i, val := range headers {
		column := rune('A' + i)
		file.SetCellValue(sheet, fmt.Sprintf("%c1", column), val)
	}

	axisHelp := utils.NewAxisHelper(1)
	for _, val := range instances {
		nextAxis := axisHelp.NewRow()

		file.SetCellValue(sheet, nextAxis(), val.Name)
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

func (as *APIService) GetMySQLUsedLicenses(filter dto.GlobalFilter) ([]dto.MySQLUsedLicense, error) {
	usedLicenses, err := as.Database.GetMySQLUsedLicenses(filter)
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

	agreements, err := as.Database.GetMySQLAgreements()
	if err != nil {
		return nil, err
	}

	hosts := make(map[string]bool)
	for _, agreement := range agreements {
		if agreement.Type == model.MySQLAgreementTypeCluster {
			for _, cluster := range agreement.Clusters {
				for _, vm := range clusters[cluster].VMs {
					hosts[vm.Hostname] = true
				}
			}
		}
	}

	for i := range usedLicenses {
		usedLicense := &usedLicenses[i]
		if hosts[usedLicense.Hostname] {
			usedLicense.AgreementType = model.MySQLAgreementTypeCluster
		} else {
			usedLicense.AgreementType = model.MySQLAgreementTypeHost
		}
	}

	return usedLicenses, nil
}

//TODO
func (as *APIService) GetMySQLDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {
	//MySQL Enterprise per cluster

	result := make([]dto.LicenseCompliance, 0) //, len(licenses))

	//perServer := dto.LicenseCompliance{
	//	LicenseTypeID:   "",
	//	ItemDescription: "MySQL Enterprise per server",
	//	Metric:          "",
	//	Consumed:        0,
	//	Covered:         0,
	//	Compliance:      0,
	//	Unlimited:       false,
	//}
	//result= append(result, perServer)
	//for _, license := range licenses {
	//	if license.Consumed == 0 {
	//		license.Compliance = 1
	//	} else {
	//		license.Compliance = license.Covered / license.Consumed
	//	}

	//	license.ItemDescription = parts[license.LicenseTypeID].ItemDescription
	//	license.Metric = parts[license.LicenseTypeID].Metric

	//	result = append(result, *license)
	//}

	return result, nil
}
