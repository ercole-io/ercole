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
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) GetOracleServiceList(filter dto.GlobalFilter) ([]dto.OracleDatabaseServiceDto, error) {
	result, err := as.Database.GetOracleServiceList(filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (as *APIService) CreateGetOracleServiceListXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	result, err := as.Database.GetOracleServiceList(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Services"
	headers := []string{
		"Hostname",
		"DB Name",
		"Service Name",
		"Container Name",
		"Enabled",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range result {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.Databasename)

		if val.OracleDatabaseService.Name != nil {
			sheets.SetCellValue(sheet, nextAxis(), *val.OracleDatabaseService.Name)
		} else {
			sheets.SetCellValue(sheet, nextAxis(), "")
		}

		sheets.SetCellValue(sheet, nextAxis(), val.OracleDatabaseService.ContainerName)

		if val.OracleDatabaseService.Enabled != nil {
			sheets.SetCellValue(sheet, nextAxis(), *val.OracleDatabaseService.Enabled)
		} else {
			sheets.SetCellValue(sheet, nextAxis(), "")
		}
	}

	return sheets, err
}
