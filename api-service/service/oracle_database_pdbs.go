// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) ListOracleDatabasePdbs(filter dto.GlobalFilter) ([]dto.OracleDatabasePluggableDatabase, error) {
	result, err := as.Database.FindAllOracleDatabasePdbs(filter)
	if err != nil {
		return nil, err
	}

	for i, n := range result {
		if n.OracleDatabasePluggableDatabase.SegmentAdvisors == nil {
			result[i].OracleDatabasePluggableDatabase.SegmentAdvisors = []model.OracleDatabaseSegmentAdvisor{}
		}
	}

	for i, n := range result {
		if n.OracleDatabasePluggableDatabase.Partitionings == nil {
			result[i].OracleDatabasePluggableDatabase.Partitionings = []model.OracleDatabasePartitioning{}
		}
	}

	return result, nil
}

func (as *APIService) CreateOracleDatabasePdbsXlsx(filter dto.GlobalFilter) (*excelize.File, error) {
	grants, err := as.Database.FindAllOracleDatabasePdbs(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Pluggable dbs"
	headers := []string{
		"Hostname",
		"DB name",
		"PDB name",
		"Status",
		"Allocable",
		"DatafileSize",
		"SegmentsSize",
		"Migrable to Postgres",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range grants {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.Dbname)
		sheets.SetCellValue(sheet, nextAxis(), val.Name)
		sheets.SetCellValue(sheet, nextAxis(), val.Status)
		sheets.SetCellValue(sheet, nextAxis(), val.Allocable)
		sheets.SetCellValue(sheet, nextAxis(), val.DatafileSize)
		sheets.SetCellValue(sheet, nextAxis(), val.SegmentsSize)
		sheets.SetCellValue(sheet, nextAxis(), val.Color)
	}

	return sheets, err
}

func (as *APIService) GetOraclePDBChanges(filter dto.GlobalFilter, hostname string, start time.Time, end time.Time) ([]dto.OraclePdbChange, error) {
	return as.Database.FindOraclePDBChangesByHostname(filter, hostname, start, end)
}
