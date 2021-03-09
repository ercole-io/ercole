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
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
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
			Name:         oracleDb["name"].(string),
			Type:         model.TechnologyOracleDatabase,
			Version:      oracleDb["version"].(string),
			Hostname:     oracleDb["hostname"].(string),
			Environment:  oracleDb["environment"].(string),
			Charset:      oracleDb["charset"].(string),
			Memory:       oracleDb["memory"].(float64),
			DatafileSize: oracleDb["datafileSize"].(float64),
			SegmentsSize: oracleDb["segmentsSize"].(float64),
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
			Name:         instance.Name,
			Type:         model.TechnologyOracleMySQL,
			Version:      instance.Version,
			Hostname:     instance.Hostname,
			Environment:  instance.Environment,
			Charset:      instance.CharsetServer,
			Memory:       instance.BufferPoolSize / 1024,
			DatafileSize: 0,
			SegmentsSize: segmentsSize / 1024,
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}
