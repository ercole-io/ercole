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

func (as *APIService) ListOracleDatabasePartitionings(filter dto.GlobalFilter) ([]dto.OracleDatabasePartitioning, error) {
	result := make([]dto.OracleDatabasePartitioning, 0)

	dbPartitionings, err := as.Database.FindAllOracleDatabasePartitionings(filter)
	if err != nil {
		return nil, err
	}

	result = append(result, dbPartitionings...)

	pdbPartitionings, err := as.Database.FindAllOraclePDBPartitionings(filter)
	if err != nil {
		return nil, err
	}

	result = append(result, pdbPartitionings...)

	return result, nil
}

func (as *APIService) CreateOracleDatabasePartitioningsXlsx(filter dto.GlobalFilter) (*excelize.File, error) {
	dbPartitionings, err := as.Database.FindAllOracleDatabasePartitionings(filter)
	if err != nil {
		return nil, err
	}

	pdbPartitionings, err := as.Database.FindAllOraclePDBPartitionings(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Partitioning"
	headers := []string{
		"Hostname",
		"Database Name",
		"PDB Name",
		"Owner",
		"Segment Name",
		"Partition Name",
		"Segment Type",
		"Mb",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range dbPartitionings {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.DatabaseName)
		sheets.SetCellValue(sheet, nextAxis(), "")
		sheets.SetCellValue(sheet, nextAxis(), val.Owner)
		sheets.SetCellValue(sheet, nextAxis(), val.SegmentName)
		sheets.SetCellValue(sheet, nextAxis(), val.PartitionName)
		sheets.SetCellValue(sheet, nextAxis(), val.SegmentType)
		sheets.SetCellValue(sheet, nextAxis(), val.Mb)
	}

	for _, valPdb := range pdbPartitionings {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), valPdb.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), valPdb.DatabaseName)
		sheets.SetCellValue(sheet, nextAxis(), valPdb.Pdb)
		sheets.SetCellValue(sheet, nextAxis(), valPdb.Owner)
		sheets.SetCellValue(sheet, nextAxis(), valPdb.SegmentName)
		sheets.SetCellValue(sheet, nextAxis(), valPdb.PartitionName)
		sheets.SetCellValue(sheet, nextAxis(), valPdb.SegmentType)
		sheets.SetCellValue(sheet, nextAxis(), valPdb.Mb)
	}

	return sheets, err
}
