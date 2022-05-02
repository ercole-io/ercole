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
		file.SetCellValue(sheet, nextAxis(), val.Name)
		file.SetCellValue(sheet, nextAxis(), val.Status)
		file.SetCellValue(sheet, nextAxis(), val.Edition)
		file.SetCellValue(sheet, nextAxis(), val.CollationName)
		file.SetCellValue(sheet, nextAxis(), val.Version)
	}

	return file, nil
}
