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
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

// SearchOracleExadata search exadata
func (as *APIService) SearchOracleExadata(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]dto.OracleExadataResponse, error) {
	return as.Database.SearchOracleExadata(full, strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

func (as *APIService) SearchOracleExadataAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	exadatas, err := as.Database.SearchOracleExadata(true, []string{}, "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	var sheets *excelize.File
	templateSheet := "template"
	sheets, err = exutils.NewXLSX(as.Config, templateSheet, "")
	if err != nil {
		return nil, err
	}

	for _, exa := range exadatas[0].Content {

		indexNewSheet := sheets.NewSheet(exa.Hostname)
		errs := sheets.CopySheet(1, indexNewSheet)
		if errs != nil {
			return nil, errs
		}

		axisHelp := exutils.NewAxisHelper(1)
		axisHelp.FillRow(sheets, exa.Hostname, exa.Hostname)
		headers := []string{
			"Hostname",
			"Model",
			"CPU",
			"Memory",
			"Version",
			"Power/Temp",
		}
		axisHelp.NewRowAndFill(sheets, exa.Hostname, headers...)
		axisHelp.NewRowAndFill(sheets, exa.Hostname, "DB Servers")
		for _, server := range exa.DbServers {
			nextAxis := axisHelp.NewRow()
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.Hostname)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.Model)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.TotalCPUCount)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.Memory)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.SwVersion)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.TotalPowerSupply)
		}
		axisHelp.NewRowAndFill(sheets, exa.Hostname, "IBSwitch")
		for _, ibSwitch := range exa.IbSwitches {
			nextAxis := axisHelp.NewRow()
			sheets.SetCellValue(exa.Hostname, nextAxis(), ibSwitch.Hostname)
			sheets.SetCellValue(exa.Hostname, nextAxis(), ibSwitch.Model)
			sheets.SetCellValue(exa.Hostname, nextAxis(), nil)
			sheets.SetCellValue(exa.Hostname, nextAxis(), nil)
			sheets.SetCellValue(exa.Hostname, nextAxis(), ibSwitch.SwVersion)
			sheets.SetCellValue(exa.Hostname, nextAxis(), nil)
		}
		axisHelp.NewRowAndFill(sheets, exa.Hostname, "Storage")
		for _, server := range exa.StorageServers {
			nextAxis := axisHelp.NewRow()
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.Hostname)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.Model)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.TotalCPUCount)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.Memory)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.SwVersion)
			sheets.SetCellValue(exa.Hostname, nextAxis(), server.TotalPowerSupply)
		}
	}

	sheets.DeleteSheet(templateSheet)

	return sheets, nil
}
