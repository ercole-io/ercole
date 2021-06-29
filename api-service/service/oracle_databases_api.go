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
	"github.com/ercole-io/ercole/v2/utils/exutils"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
)

// SearchOracleDatabaseAddms search addms
func (as *APIService) SearchOracleDatabaseAddms(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error) {
	return as.Database.SearchOracleDatabaseAddms(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// SearchOracleDatabaseSegmentAdvisors search segment advisors
func (as *APIService) SearchOracleDatabaseSegmentAdvisors(search string, sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]dto.OracleDatabaseSegmentAdvisor, error) {
	return as.Database.SearchOracleDatabaseSegmentAdvisors(strings.Split(search, " "), sortBy, sortDesc, location, environment, olderThan)
}

func (as *APIService) SearchOracleDatabaseSegmentAdvisorsAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	segmentAdvisors, err := as.Database.SearchOracleDatabaseSegmentAdvisors([]string{}, "", false, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	sheet := "Segment_Advisor"
	headers := []string{
		"ReclaimableGB",
		"GB Total",
		"Retrieve",
		"Hostname",
		"DB Names",
		"Segment Owner",
		"Segment Name",
		"Segment Type",
		"Partition Name",
		"Recommendation",
	}
	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}
	axisHelp := exutils.NewAxisHelper(1)

	for _ , val := range segmentAdvisors {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet,nextAxis(), val.Reclaimable)
		sheets.SetCellValue(sheet,nextAxis(), val.SegmentsSize)
		if val.SegmentsSize == 0 {
			nextAxis()
		}else {
			sheets.SetCellValue(sheet,nextAxis(), val.Reclaimable/val.SegmentsSize)
		}
		sheets.SetCellValue(sheet,nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet,nextAxis(), val.Dbname)
		sheets.SetCellValue(sheet,nextAxis(), val.SegmentOwner)
		sheets.SetCellValue(sheet,nextAxis(), val.SegmentName)
		sheets.SetCellValue(sheet,nextAxis(), val.SegmentType)
		sheets.SetCellValue(sheet,nextAxis(), val.PartitionName)
		sheets.SetCellValue(sheet,nextAxis(), val.Recommendation)
	}

	return sheets, err
}

// SearchOracleDatabasePatchAdvisors search patch advisors
func (as *APIService) SearchOracleDatabasePatchAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) ([]map[string]interface{}, error) {
	return as.Database.SearchOracleDatabasePatchAdvisors(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, windowTime, location, environment, olderThan, status)
}

// SearchOracleDatabases search databases
func (as *APIService) SearchOracleDatabases(f dto.SearchOracleDatabasesFilter) ([]map[string]interface{}, error) {
	return as.Database.SearchOracleDatabases(f.Full, strings.Split(f.Search, " "), f.SortBy, f.SortDesc,
		f.PageNumber, f.PageSize, f.Location, f.Environment, f.OlderThan)
}

func (as *APIService) SearchOracleDatabasesAsXLSX(filter dto.SearchOracleDatabasesFilter) (*excelize.File, error) {
	databases, err := as.Database.SearchOracleDatabases(false, strings.Split(filter.Search, " "),
		filter.SortBy, filter.SortDesc,
		-1, -1,
		filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	file, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_databases.xlsx")
	if err != nil {
		return nil, err
	}

	for i, val := range databases {
		i += 2 // offset for headers
		file.SetCellValue("Databases", fmt.Sprintf("A%d", i), val["name"])
		file.SetCellValue("Databases", fmt.Sprintf("B%d", i), val["uniqueName"])
		file.SetCellValue("Databases", fmt.Sprintf("C%d", i), val["version"])
		file.SetCellValue("Databases", fmt.Sprintf("D%d", i), val["hostname"])
		file.SetCellValue("Databases", fmt.Sprintf("E%d", i), val["status"])
		file.SetCellValue("Databases", fmt.Sprintf("F%d", i), val["environment"])
		file.SetCellValue("Databases", fmt.Sprintf("G%d", i), val["location"])
		file.SetCellValue("Databases", fmt.Sprintf("H%d", i), val["charset"])
		file.SetCellValue("Databases", fmt.Sprintf("I%d", i), val["blockSize"])
		file.SetCellValue("Databases", fmt.Sprintf("J%d", i), val["cpuCount"])
		file.SetCellValue("Databases", fmt.Sprintf("K%d", i), val["work"])
		file.SetCellValue("Databases", fmt.Sprintf("L%d", i), val["memory"])
		file.SetCellValue("Databases", fmt.Sprintf("M%d", i), val["datafileSize"])
		file.SetCellValue("Databases", fmt.Sprintf("N%d", i), val["segmentsSize"])
		file.SetCellValue("Databases", fmt.Sprintf("O%d", i), val["archivelog"])
		file.SetCellValue("Databases", fmt.Sprintf("P%d", i), val["dataguard"])
		file.SetCellValue("Databases", fmt.Sprintf("Q%d", i), val["rac"])
		file.SetCellValue("Databases", fmt.Sprintf("R%d", i), val["ha"])
	}

	return file, nil
}

// SearchOracleDatabaseUsedLicenses return the list of used licenses
func (as *APIService) SearchOracleDatabaseUsedLicenses(sortBy string, sortDesc bool, page int, pageSize int,
	location string, environment string, olderThan time.Time,
) (*dto.OracleDatabaseUsedLicenseSearchResponse, error) {
	return as.Database.SearchOracleDatabaseUsedLicenses(sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}
