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

import "github.com/ercole-io/ercole/v2/api-service/dto"

// SearchMySQLDatabases search databases
func (as *APIService) SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error) {
	instances, err := as.Database.SearchMySQLInstances(filter)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

//TODO
//func (as *APIService) SearchMySQLDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
//	databases, aerr := as.Database.SearchMySQLDatabases(false, strings.Split(filter.Search, " "),
//		filter.SortBy, filter.SortDesc,
//		-1, -1,
//		filter.Location, filter.Environment, filter.OlderThan)
//	if aerr != nil {
//		return nil, aerr
//	}
//
//	file, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_databases.xlsx")
//	if err != nil {
//		return nil, err
//	}
//
//	for i, val := range databases {
//		i += 2 // offset for headers
//		file.SetCellValue("Databases", fmt.Sprintf("A%d", i), val["name"])
//		file.SetCellValue("Databases", fmt.Sprintf("B%d", i), val["uniqueName"])
//		file.SetCellValue("Databases", fmt.Sprintf("C%d", i), val["version"])
//		file.SetCellValue("Databases", fmt.Sprintf("D%d", i), val["hostname"])
//		file.SetCellValue("Databases", fmt.Sprintf("E%d", i), val["status"])
//		file.SetCellValue("Databases", fmt.Sprintf("F%d", i), val["environment"])
//		file.SetCellValue("Databases", fmt.Sprintf("G%d", i), val["location"])
//		file.SetCellValue("Databases", fmt.Sprintf("H%d", i), val["charset"])
//		file.SetCellValue("Databases", fmt.Sprintf("I%d", i), val["blockSize"])
//		file.SetCellValue("Databases", fmt.Sprintf("J%d", i), val["cpuCount"])
//		file.SetCellValue("Databases", fmt.Sprintf("K%d", i), val["work"])
//		file.SetCellValue("Databases", fmt.Sprintf("L%d", i), val["memory"])
//		file.SetCellValue("Databases", fmt.Sprintf("M%d", i), val["datafileSize"])
//		file.SetCellValue("Databases", fmt.Sprintf("N%d", i), val["segmentsSize"])
//		file.SetCellValue("Databases", fmt.Sprintf("O%d", i), val["archivelog"])
//		file.SetCellValue("Databases", fmt.Sprintf("P%d", i), val["dataguard"])
//		file.SetCellValue("Databases", fmt.Sprintf("Q%d", i), val["rac"])
//		file.SetCellValue("Databases", fmt.Sprintf("R%d", i), val["ha"])
//	}
//
//	return file, nil
//}
