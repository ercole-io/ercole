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
	instances, aerr := as.Database.SearchMySQLInstances(filter)
	if aerr != nil {
		return nil, aerr
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
