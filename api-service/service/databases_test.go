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

package service

import (
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchDatabases_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	oracleDbs := []map[string]interface{}{
		{
			"name":        "pippo",
			"version":     "",
			"hostname":    "",
			"environment": "",
			"charset":     "",

			"memory":       42.42,
			"datafileSize": 75.42,
			"segmentsSize": 99.42,
			"archivelog":   true,
			"ha":           false,
			"dataguard":    true,
		},
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	db.EXPECT().SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(oracleDbs, nil)

	mysqlInstances := []dto.MySQLInstance{
		{
			Hostname:    "pluto",
			Location:    "Cuba",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "mysql",
				Version:            "",
				Edition:            "",
				Platform:           "",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           1,
				ThreadsConcurrency: 2,
				BufferPoolSize:     43008,
				LogBufferSize:      4,
				SortBufferSize:     5,
				ReadOnly:           false,
				LogBin:             true,
				HighAvailability:   false,
				UUID:               "000000000000000000000000",
				IsMaster:           true,
				SlaveUUIDs:         []string{"111111111111111111111111"},
				IsSlave:            false,
				MasterUUID:         new(string),
				Databases:          []model.MySQLDatabase{{Name: "", Charset: "", Collation: "", Encrypted: false}},
				TableSchemas:       []model.MySQLTableSchema{{Name: "", Engine: "", Allocation: 24576}},
				SegmentAdvisors:    []model.MySQLSegmentAdvisor{{TableSchema: "", TableName: "", Engine: "", Allocation: 76, Data: 0, Index: 0, Free: 0}},
			},
		},
	}

	globalFilter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	actual, err := as.SearchDatabases(globalFilter)
	require.NoError(t, err)

	expected := []dto.Database{
		{
			Name:             "pippo",
			Type:             "Oracle/Database",
			Version:          "",
			Hostname:         "",
			Environment:      "",
			Charset:          "",
			Memory:           42.42,
			DatafileSize:     75.42,
			SegmentsSize:     99.42,
			Archivelog:       true,
			HighAvailability: false,
			DisasterRecovery: true,
		},
		{
			Name:             "mysql",
			Type:             "Oracle/MySQL",
			Version:          "",
			Hostname:         "pluto",
			Environment:      "TST",
			Charset:          "",
			Memory:           42.0,
			DatafileSize:     0,
			SegmentsSize:     24.0,
			Archivelog:       true,
			HighAvailability: false,
			DisasterRecovery: true,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestSearchDatabasesAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}

	oracleDbs := []map[string]interface{}{
		{
			"name":        "pippo",
			"version":     "",
			"hostname":    "",
			"environment": "",
			"charset":     "",

			"memory":       42.42,
			"datafileSize": 75.42,
			"segmentsSize": 99.42,
			"archivelog":   true,
			"ha":           false,
			"dataguard":    true,
		},
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	db.EXPECT().SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(oracleDbs, nil)

	mysqlInstances := []dto.MySQLInstance{
		{
			Hostname:    "pluto",
			Location:    "Cuba",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "mysql",
				Version:            "",
				Edition:            "",
				Platform:           "",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           1,
				ThreadsConcurrency: 2,
				BufferPoolSize:     43008,
				LogBufferSize:      4,
				SortBufferSize:     5,
				ReadOnly:           false,
				LogBin:             false,
				HighAvailability:   false,
				UUID:               "",
				IsMaster:           false,
				SlaveUUIDs:         []string{},
				IsSlave:            false,
				MasterUUID:         new(string),
				Databases:          []model.MySQLDatabase{{Name: "", Charset: "", Collation: "", Encrypted: false}},
				TableSchemas:       []model.MySQLTableSchema{{Name: "", Engine: "", Allocation: 24576}},
				SegmentAdvisors:    []model.MySQLSegmentAdvisor{{TableSchema: "", TableName: "", Engine: "", Allocation: 76, Data: 0, Index: 0, Free: 0}},
			},
		},
	}

	globalFilter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	actual, err := as.SearchDatabasesAsXLSX(globalFilter)
	require.NoError(t, err)

	expected := []dto.Database{
		{
			Name:         "pippo",
			Type:         "Oracle/Database",
			Version:      "",
			Hostname:     "",
			Environment:  "",
			Charset:      "",
			Memory:       42.42,
			DatafileSize: 75.42,
			SegmentsSize: 99.42,
		},
		{
			Name:         "mysql",
			Type:         "Oracle/MySQL",
			Version:      "",
			Hostname:     "pluto",
			Environment:  "TST",
			Charset:      "",
			Memory:       42.0,
			DatafileSize: 0,
			SegmentsSize: 24.0,
		},
	}

	assert.Equal(t, "Name", actual.GetCellValue("Databases", "A1"))
	assert.Equal(t, expected[0].Name, actual.GetCellValue("Databases", "A2"))
	assert.Equal(t, expected[1].Name, actual.GetCellValue("Databases", "A3"))

	assert.Equal(t, "Type", actual.GetCellValue("Databases", "B1"))
	assert.Equal(t, expected[0].Type, actual.GetCellValue("Databases", "B2"))
	assert.Equal(t, expected[1].Type, actual.GetCellValue("Databases", "B3"))

	assert.Equal(t, "Memory", actual.GetCellValue("Databases", "G1"))
	assert.Equal(t, "42.42", actual.GetCellValue("Databases", "G2"))
	assert.Equal(t, "42", actual.GetCellValue("Databases", "G3"))
}

func TestGetDatabasesStatistics_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	oracleDbs := []map[string]interface{}{
		{
			"name":        "pippo",
			"version":     "",
			"hostname":    "",
			"environment": "",
			"charset":     "",

			"memory":       42.42,
			"datafileSize": 75.42,
			"segmentsSize": 99.42,
			"archivelog":   true,
			"ha":           false,
			"dataguard":    true,
		},
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	db.EXPECT().SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(oracleDbs, nil)

	mysqlInstances := []dto.MySQLInstance{
		{
			Hostname:    "pluto",
			Location:    "Cuba",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "mysql",
				Version:            "",
				Edition:            "",
				Platform:           "",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           1,
				ThreadsConcurrency: 2,
				BufferPoolSize:     43008,
				LogBufferSize:      4,
				SortBufferSize:     5,
				ReadOnly:           false,
				Databases: []model.MySQLDatabase{
					{
						Name:      "",
						Charset:   "",
						Collation: "",
						Encrypted: false,
					},
				},
				TableSchemas: []model.MySQLTableSchema{
					{
						Name:       "",
						Engine:     "",
						Allocation: 24576,
					},
				},
				SegmentAdvisors: []model.MySQLSegmentAdvisor{
					{
						TableSchema: "",
						TableName:   "",
						Engine:      "",
						Allocation:  76,
						Data:        0,
						Index:       0,
						Free:        0,
					},
				},
			},
		},
	}

	globalFilter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	actual, err := as.GetDatabasesStatistics(globalFilter)
	require.NoError(t, err)

	expected := dto.DatabasesStatistics{
		TotalMemorySize:   84.42,
		TotalSegmentsSize: 123.42,
	}

	assert.Equal(t, expected, *actual)
}
