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

func (as *APIService) GetOracleDiskGroups(hostname string, dbname string) ([]dto.OracleDatabaseDiskGroupDto, error) {
	result, err := as.Database.GetOracleDiskGroups(hostname, dbname)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (as *APIService) ListOracleDiskGroups(filter dto.GlobalFilter) ([]dto.OracleDatabaseDiskGroupDto, error) {
	result, err := as.Database.ListOracleDiskGroups(filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (as *APIService) CreateOracleDiskGroupsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	result, err := as.Database.ListOracleDiskGroups(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Disk Groups"
	headers := []string{
		"Hostname",
		"Databases",
		"Disk group name",
		"Total space",
		"Free space",
		"Percentage of free space",
		"Used space",
		"Percentage of used space",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range result {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.Databases)
		sheets.SetCellValue(sheet, nextAxis(), val.OracleDatabaseDiskGroup.DiskGroupName)
		sheets.SetCellValue(sheet, nextAxis(), val.OracleDatabaseDiskGroup.TotalSpace)
		sheets.SetCellValue(sheet, nextAxis(), val.OracleDatabaseDiskGroup.FreeSpace)
		sheets.SetCellValue(sheet, nextAxis(), val.PercentageFreeSpace)
		sheets.SetCellValue(sheet, nextAxis(), val.UsedSpace)
		sheets.SetCellValue(sheet, nextAxis(), val.PercentageUsedSpace)
	}

	return sheets, err
}
