// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) GetClusterVeritasLicenses(filter dto.GlobalFilter) ([]dto.ClusterVeritasLicense, error) {
	clusterVeritasLicenses, err := as.Database.FindClusterVeritasLicenses(filter)
	if err != nil {
		return nil, err
	}

	clusterVeritasLicensesMap := make(map[string][]dto.ClusterVeritasLicense)

	for _, clusterLicenses := range clusterVeritasLicenses {
		if clusterLicenses.ID != "" {
			clusterVeritasLicensesMap[clusterLicenses.ID] = append(clusterVeritasLicensesMap[clusterLicenses.ID], clusterLicenses)
		}
	}

	for k, v := range clusterVeritasLicensesMap {
		if strings.Contains(k, "_DR") {
			existingHostsDR := v[0].Hostnames

			realclusterID := strings.Replace(k, "_DR", "", -1)
			realclusterLicenses := clusterVeritasLicensesMap[realclusterID]

			for _, realLicense := range realclusterLicenses {
				if !containsClusterVeritasLicense(v, realLicense) {
					clusterVeritasLicenses = append(clusterVeritasLicenses, dto.ClusterVeritasLicense{
						ID:            k,
						LicenseTypeID: realLicense.LicenseTypeID,
						Description:   realLicense.Description,
						Metric:        realLicense.Metric,
						Count:         clusterVeritasLicensesMap[k][0].Count,
						Hostnames:     existingHostsDR,
					})
				}
			}
		}
	}

	clusterVeritasLicenses = removeDuplicates(clusterVeritasLicenses)

	return clusterVeritasLicenses, nil
}

func removeDuplicates(licenses []dto.ClusterVeritasLicense) []dto.ClusterVeritasLicense {
	check := make(map[string]bool)
	unique := []dto.ClusterVeritasLicense{}

	for _, license := range licenses {
		if license.ID == "" {
			continue
		}

		key := license.LicenseTypeID + "|" + license.Description + "|" + license.Metric + "|" + license.ID
		if !check[key] {
			check[key] = true

			unique = append(unique, license)
		}
	}

	return unique
}

func containsClusterVeritasLicense(licenses []dto.ClusterVeritasLicense, license dto.ClusterVeritasLicense) bool {
	for _, v := range licenses {
		if v.LicenseTypeID == license.LicenseTypeID &&
			v.Description == license.Description &&
			v.Metric == license.Metric {
			return true
		}
	}

	return false
}

func (as *APIService) GetClusterVeritasLicensesXlsx(filter dto.GlobalFilter) (*excelize.File, error) {
	clusterVeritasLicenses, err := as.Database.FindClusterVeritasLicenses(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Cluster Veritas Licenses"
	headers := []string{
		"ID",
		"Hostnames",
		"LicenseTypeID",
		"Description",
		"Metric",
		"Count",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range clusterVeritasLicenses {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.ID)
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.Hostnames, ", "))
		sheets.SetCellValue(sheet, nextAxis(), val.LicenseTypeID)
		sheets.SetCellValue(sheet, nextAxis(), val.Description)
		sheets.SetCellValue(sheet, nextAxis(), val.Metric)
		sheets.SetCellValue(sheet, nextAxis(), val.Count)
	}

	return sheets, err
}
