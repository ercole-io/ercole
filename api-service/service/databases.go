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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (as *APIService) SearchDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	type getter func(filter dto.GlobalFilter) ([]dto.Database, error)
	getters := []getter{as.getOracleDatabases, as.getMySQLDatabases}

	dbs := make([]dto.Database, 0)
	for _, get := range getters {
		thisDbs, err := get(filter)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, thisDbs...)
	}

	return dbs, nil
}

func (as *APIService) getOracleDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	sodf := dto.SearchOracleDatabasesFilter{
		GlobalFilter: filter,
		Full:         false,
		PageNumber:   -1,
		PageSize:     -1,
	}
	oracleDbs, err := as.SearchOracleDatabases(sodf)
	if err != nil {
		return nil, err
	}

	dbs := make([]dto.Database, 0)
	for _, oracleDb := range oracleDbs {
		db := dto.Database{
			Name:             oracleDb["name"].(string),
			Type:             model.TechnologyOracleDatabase,
			Version:          oracleDb["version"].(string),
			Hostname:         oracleDb["hostname"].(string),
			Environment:      oracleDb["environment"].(string),
			Charset:          oracleDb["charset"].(string),
			Memory:           oracleDb["memory"].(float64),
			DatafileSize:     oracleDb["datafileSize"].(float64),
			SegmentsSize:     oracleDb["segmentsSize"].(float64),
			Archivelog:       oracleDb["archivelog"].(bool),
			HighAvailability: oracleDb["ha"].(bool),
			DisasterRecovery: oracleDb["dataguard"].(bool),
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (as *APIService) getMySQLDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	mysqlInstances, err := as.Database.SearchMySQLInstances(filter)
	if err != nil {
		return nil, err
	}

	dbs := make([]dto.Database, 0)
	for _, instance := range mysqlInstances {

		segmentsSize := 0.0
		for _, ts := range instance.TableSchemas {
			segmentsSize += ts.Allocation
		}

		db := dto.Database{
			Name:             instance.Name,
			Type:             model.TechnologyOracleMySQL,
			Version:          instance.Version,
			Hostname:         instance.Hostname,
			Environment:      instance.Environment,
			Charset:          instance.CharsetServer,
			Memory:           instance.BufferPoolSize / 1024,
			DatafileSize:     0,
			SegmentsSize:     segmentsSize / 1024,
			Archivelog:       instance.LogBin,
			HighAvailability: instance.HighAvailability,
			DisasterRecovery: instance.IsMaster || instance.IsSlave,
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (as *APIService) SearchDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	databases, err := as.SearchDatabases(filter)
	if err != nil {
		return nil, err
	}

	file, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_generic.xlsx")
	if err != nil {
		return nil, err
	}

	sheet := "Databases"
	file.SetSheetName("Sheet1", sheet)
	headers := []string{
		"Name",
		"Type",
		"Version",
		"Hostname",
		"Environment",
		"Charset",
		"Memory",
		"Datafile Size",
		"Segments Size",
	}

	for i, val := range headers {
		column := rune('A' + i)
		file.SetCellValue(sheet, fmt.Sprintf("%c1", column), val)
	}

	axisHelp := utils.NewAxisHelper(1)
	for _, val := range databases {
		nextAxis := axisHelp.NewRow()

		file.SetCellValue(sheet, nextAxis(), val.Name)
		file.SetCellValue(sheet, nextAxis(), val.Type)
		file.SetCellValue(sheet, nextAxis(), val.Version)
		file.SetCellValue(sheet, nextAxis(), val.Hostname)
		file.SetCellValue(sheet, nextAxis(), val.Environment)
		file.SetCellValue(sheet, nextAxis(), val.Charset)
		file.SetCellValue(sheet, nextAxis(), val.Memory)
		file.SetCellValue(sheet, nextAxis(), val.DatafileSize)
		file.SetCellValue(sheet, nextAxis(), val.SegmentsSize)
	}

	return file, nil
}

func (as *APIService) GetDatabasesStatistics(filter dto.GlobalFilter) (*dto.DatabasesStatistics, error) {
	dbs, err := as.SearchDatabases(filter)
	if err != nil {
		return nil, err
	}

	stats := new(dto.DatabasesStatistics)
	for _, db := range dbs {
		stats.TotalMemorySize += db.Memory * 1024 * 1024 * 1024         // From GBytes to bytes
		stats.TotalSegmentsSize += db.SegmentsSize * 1024 * 1024 * 1024 // From GBytes to bytes
	}

	return stats, nil
}

func (as *APIService) GetDatabasesUsedLicenses(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	type getter func(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error)
	getters := []getter{as.getOracleDatabasesUsedLicenses, as.getMySQLUsedLicenses}

	usedLicenses := make([]dto.DatabaseUsedLicense, 0)
	for _, get := range getters {
		thisDbs, err := get(filter)
		if err != nil {
			return nil, err
		}
		usedLicenses = append(usedLicenses, thisDbs...)
	}

	return usedLicenses, nil
}

func (as *APIService) getOracleDatabasesUsedLicenses(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	oracleLics, err := as.Database.SearchOracleDatabaseUsedLicenses("", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	licenseTypes, err := as.GetOracleDatabaseLicenseTypesAsMap()
	if err != nil {
		return nil, err
	}

	genericLics := make([]dto.DatabaseUsedLicense, 0, len(oracleLics.Content))
	for _, o := range oracleLics.Content {
		lt := licenseTypes[o.LicenseTypeID]

		g := dto.DatabaseUsedLicense{
			Hostname:      o.Hostname,
			DbName:        o.DbName,
			LicenseTypeID: o.LicenseTypeID,
			Description:   lt.ItemDescription,
			Metric:        lt.Metric,
			UsedLicenses:  o.UsedLicenses,
		}

		genericLics = append(genericLics, g)
	}

	return genericLics, nil
}

func (as *APIService) getMySQLUsedLicenses(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	mysqlLics, err := as.GetMySQLUsedLicenses(filter)
	if err != nil {
		return nil, err
	}

	genericLics := make([]dto.DatabaseUsedLicense, 0, len(mysqlLics))

	for _, lic := range mysqlLics {
		g := dto.DatabaseUsedLicense{
			Hostname:      lic.Hostname,
			DbName:        lic.InstanceName,
			LicenseTypeID: "",
			Description:   "MySQL " + lic.InstanceEdition,
			Metric:        lic.AgreementType,
			UsedLicenses:  1,
		}

		genericLics = append(genericLics, g)
	}

	return genericLics, nil
}

func (as *APIService) GetDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {
	licenses := make([]dto.LicenseCompliance, 0)

	oracle, err := as.GetOracleDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}
	licenses = append(licenses, oracle...)

	mysql, err := as.GetMySQLDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}
	licenses = append(licenses, mysql...)
	return licenses, nil
}
