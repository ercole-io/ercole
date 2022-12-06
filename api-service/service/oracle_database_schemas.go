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

func (as *APIService) ListOracleDatabaseSchemas(filter dto.GlobalFilter) ([]dto.OracleDatabaseSchema, error) {
	result := make([]dto.OracleDatabaseSchema, 0)

	dbSchemas, err := as.Database.FindAllOracleDatabaseSchemas(filter)
	if err != nil {
		return nil, err
	}

	result = append(result, dbSchemas...)

	pdbSchemas, err := as.Database.FindAllOraclePDBSchemas(filter)
	if err != nil {
		return nil, err
	}

	result = append(result, pdbSchemas...)

	return result, nil
}

func (as *APIService) CreateOracleDatabaseSchemasXlsx(filter dto.GlobalFilter) (*excelize.File, error) {
	grants, err := as.Database.FindAllOracleDatabaseSchemas(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Schemas"
	headers := []string{
		"Hostname",
		"DatabaseName",
		"Indexes",
		"LOB",
		"Tables",
		"Total",
		"User",
		"Account Status",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range grants {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.DatabaseName)
		sheets.SetCellValue(sheet, nextAxis(), val.Indexes)
		sheets.SetCellValue(sheet, nextAxis(), val.LOB)
		sheets.SetCellValue(sheet, nextAxis(), val.Tables)
		sheets.SetCellValue(sheet, nextAxis(), val.Total)
		sheets.SetCellValue(sheet, nextAxis(), val.User)
		sheets.SetCellValue(sheet, nextAxis(), val.AccountStatus)
	}

	return sheets, err
}
