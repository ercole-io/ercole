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

	agreements, err := as.Database.GetMySQLAgreements()
	if err != nil {
		return nil, err
	}

	hostCoveredForCluster := make(map[string]bool, len(agreements))
	hostCoveredAsHost := make(map[string]bool, len(agreements))

	for _, agreement := range agreements {
		if agreement.Type == model.MySQLAgreementTypeCluster {
			for _, cluster := range agreement.Clusters {
				for _, vm := range clusters[cluster].VMs {
					hostCoveredForCluster[vm.Hostname] = true
				}
			}

			continue
		}

		if agreement.Type == model.MySQLAgreementTypeHost {
			for _, host := range agreement.Hosts {
				hostCoveredAsHost[host] = true
			}

			continue
		}

		as.Log.Errorf("Unknown MySQLAgreementType: %s", agreement.Type)
	}

	for i := range usedLicenses {
		usedLicense := &usedLicenses[i]
		if hostCoveredForCluster[usedLicense.Hostname] {
			usedLicense.AgreementType = model.MySQLAgreementTypeCluster
			usedLicense.Covered = true

			continue
		}

		usedLicense.AgreementType = model.MySQLAgreementTypeHost
	}

	return usedLicenses, nil
}

func (as *APIService) GetMySQLDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {
	any := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	licenses, err := as.GetMySQLUsedLicenses("", any)
	if err != nil {
		return nil, err
	}

	if len(licenses) == 0 {
		return []dto.LicenseCompliance{}, nil
	}

	perCluster := dto.LicenseCompliance{
		LicenseTypeID:   "",
		ItemDescription: "MySQL Enterprise per cluster",
		Metric:          "",
		Cost:            0,
		Consumed:        0,
		Covered:         0,
		Compliance:      0,
		Unlimited:       false,
	}

	perHost := dto.LicenseCompliance{
		LicenseTypeID:   "",
		ItemDescription: "MySQL Enterprise per host",
		Metric:          "",
		Cost:            0,
		Consumed:        0,
		Covered:         0,
		Compliance:      0,
		Unlimited:       false,
	}

	for _, license := range licenses {
		var lc *dto.LicenseCompliance
		if license.AgreementType == model.MySQLAgreementTypeHost {
			lc = &perHost
		} else if license.AgreementType == model.MySQLAgreementTypeCluster {
			lc = &perCluster
		} else {
			as.Log.Errorf("Unknown MySQLAgreementType: %s", license.AgreementType)
			continue
		}

		if license.Covered {
			lc.Covered += 1
		}

		lc.Consumed += 1
	}

	result := make([]dto.LicenseCompliance, 0, 2)

	for _, lc := range []*dto.LicenseCompliance{&perCluster, &perHost} {
		if lc.Consumed == 0 {
			lc.Compliance = 1
		} else {
			lc.Compliance = lc.Covered / lc.Consumed
		}

		result = append(result, *lc)
	}

	return result, nil
}
