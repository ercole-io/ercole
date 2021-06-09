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

package database

import (
	"context"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchMySQLInstances() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_20.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_21.json"))
	first := dto.MySQLInstance{
		Hostname:    "erc-mysql",
		Location:    "Germany",
		Environment: "TST",
		MySQLInstance: model.MySQLInstance{
			Name:               "mysql:3306",
			Version:            "8.0.23",
			Edition:            "COMMUNITY",
			Platform:           "Linux",
			Architecture:       "x86_64",
			Engine:             "InnoDB",
			RedoLogEnabled:     "ON",
			CharsetServer:      "utf8mb4",
			CharsetSystem:      "utf8",
			PageSize:           16,
			ThreadsConcurrency: 0,
			BufferPoolSize:     128,
			LogBufferSize:      16,
			SortBufferSize:     1,
			ReadOnly:           false,
			Databases: []model.MySQLDatabase{{
				Name:      "mysql",
				Charset:   "utf8mb4",
				Collation: "utf8mb4_0900_ai_ci",
				Encrypted: false},
				{
					Name:      "information_schema",
					Charset:   "utf8",
					Collation: "utf8_general_ci",
					Encrypted: false},
				{
					Name:      "performance_schema",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Encrypted: false},
				{
					Name:      "sys",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Encrypted: false}},
			TableSchemas: []model.MySQLTableSchema{
				{
					Name:       "innodb_temporary",
					Engine:     "InnoDB",
					Allocation: 12},
				{
					Name:       "innodb_undo_001",
					Engine:     "InnoDB",
					Allocation: 16},
				{
					Name:       "innodb_undo_002",
					Engine:     "InnoDB",
					Allocation: 16},
				{
					Name:       "mysql",
					Engine:     "InnoDB",
					Allocation: 24},
				{
					Name:       "sys",
					Engine:     "InnoDB",
					Allocation: 0.078}},
			SegmentAdvisors: []model.MySQLSegmentAdvisor{
				{
					TableSchema: "innodb_temporary",
					TableName:   "innodb_temporary",
					Engine:      "InnoDB",
					Allocation:  12,
					Data:        12,
					Index:       0,
					Free:        0},
				{
					TableSchema: "innodb_undo_001",
					TableName:   "innodb_undo_001",
					Engine:      "InnoDB",
					Allocation:  16,
					Data:        16,
					Index:       0,
					Free:        0},
				{
					TableSchema: "innodb_undo_002",
					TableName:   "innodb_undo_002",
					Engine:      "InnoDB",
					Allocation:  16,
					Data:        16,
					Index:       0,
					Free:        0},
				{
					TableSchema: "mysql",
					TableName:   "mysql",
					Engine:      "InnoDB",
					Allocation:  24,
					Data:        24,
					Index:       0,
					Free:        0},
				{
					TableSchema: "sys",
					TableName:   "sys_config",
					Engine:      "InnoDB",
					Allocation:  0.078,
					Data:        0.062,
					Index:       0.016,
					Free:        0,
				},
			},
		},
	}
	second := dto.MySQLInstance{
		Hostname:    "erc-mysql-prod",
		Location:    "Cuba",
		Environment: "PROD",
		MySQLInstance: model.MySQLInstance{
			Name:               "mysql:3306",
			Version:            "8.0.23",
			Edition:            "COMMUNITY",
			Platform:           "Linux",
			Architecture:       "x86_64",
			Engine:             "InnoDB",
			RedoLogEnabled:     "ON",
			CharsetServer:      "utf8mb4",
			CharsetSystem:      "utf8",
			PageSize:           16,
			ThreadsConcurrency: 0,
			BufferPoolSize:     128,
			LogBufferSize:      16,
			SortBufferSize:     1,
			ReadOnly:           false,
			Databases: []model.MySQLDatabase{
				{
					Name:      "mysql",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Encrypted: false},
				{
					Name:      "information_schema",
					Charset:   "utf8",
					Collation: "utf8_general_ci",
					Encrypted: false},
				{
					Name:      "performance_schema",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Encrypted: false},
				{
					Name:      "sys",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
					Encrypted: false},
			},
			TableSchemas: []model.MySQLTableSchema{
				{
					Name:       "innodb_temporary",
					Engine:     "InnoDB",
					Allocation: 12},
				{
					Name:       "innodb_undo_001",
					Engine:     "InnoDB",
					Allocation: 16},
				{
					Name:       "innodb_undo_002",
					Engine:     "InnoDB",
					Allocation: 16},
				{
					Name:       "mysql",
					Engine:     "InnoDB",
					Allocation: 24},
				{
					Name:       "sys",
					Engine:     "InnoDB",
					Allocation: 0.078},
			},
			SegmentAdvisors: []model.MySQLSegmentAdvisor{
				{
					TableSchema: "innodb_temporary",
					TableName:   "innodb_temporary",
					Engine:      "InnoDB",
					Allocation:  12,
					Data:        12,
					Index:       0,
					Free:        0},
				{
					TableSchema: "innodb_undo_001",
					TableName:   "innodb_undo_001",
					Engine:      "InnoDB",
					Allocation:  16,
					Data:        16,
					Index:       0,
					Free:        0},
				{
					TableSchema: "innodb_undo_002",
					TableName:   "innodb_undo_002",
					Engine:      "InnoDB",
					Allocation:  16,
					Data:        16,
					Index:       0,
					Free:        0},
				{
					TableSchema: "mysql",
					TableName:   "mysql",
					Engine:      "InnoDB",
					Allocation:  24,
					Data:        24,
					Index:       0,
					Free:        0},
				{
					TableSchema: "sys",
					TableName:   "sys_config",
					Engine:      "InnoDB",
					Allocation:  0.078,
					Data:        0.062,
					Index:       0.016,
					Free:        0},
			},
		},
	}

	m.T().Run("should_load_all", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "",
			Environment: "",
			OlderThan:   utils.MAX_TIME,
		}
		actual, err := m.db.SearchMySQLInstances(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLInstance{first, second}
		assert.Equal(t, expected, actual)
	})

	m.T().Run("should_filter_by_location", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "Cuba",
			Environment: "",
			OlderThan:   utils.MAX_TIME,
		}
		actual, err := m.db.SearchMySQLInstances(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLInstance{second}
		assert.Equal(t, expected, actual)
	})

	m.T().Run("should_filter_by_environment", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "",
			Environment: "TST",
			OlderThan:   utils.MAX_TIME,
		}
		actual, err := m.db.SearchMySQLInstances(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLInstance{first}
		assert.Equal(t, expected, actual)
	})

	m.T().Run("should_filter_by_older_than", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "",
			Environment: "",
			OlderThan:   utils.P("2021-03-03T09:00:32.981Z"),
		}
		actual, err := m.db.SearchMySQLInstances(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLInstance{second}
		assert.Equal(t, expected, actual)
	})
}

func (m *MongodbSuite) TestGetMySQLUsedLicenses() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_20.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_23.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_24.json"))
	first := dto.MySQLUsedLicense{
		Hostname:        "erc-mysql-2",
		InstanceName:    "mysql:3306",
		InstanceEdition: "ENTERPRISE",
	}
	second := dto.MySQLUsedLicense{
		Hostname:        "erc-mysql-prod-2",
		InstanceName:    "mysql:3306",
		InstanceEdition: "ENTERPRISE",
	}

	m.T().Run("should_load_all", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "",
			Environment: "",
			OlderThan:   utils.MAX_TIME,
		}
		actual, err := m.db.GetMySQLUsedLicenses(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLUsedLicense{first, second}
		assert.Equal(t, expected, actual)
	})

	m.T().Run("should_filter_by_location", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "Cuba",
			Environment: "",
			OlderThan:   utils.MAX_TIME,
		}
		actual, err := m.db.GetMySQLUsedLicenses(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLUsedLicense{second}
		assert.Equal(t, expected, actual)
	})

	m.T().Run("should_filter_by_environment", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "",
			Environment: "TST",
			OlderThan:   utils.MAX_TIME,
		}
		actual, err := m.db.GetMySQLUsedLicenses(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLUsedLicense{first}
		assert.Equal(t, expected, actual)
	})

	m.T().Run("should_filter_by_older_than", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "",
			Environment: "",
			OlderThan:   utils.P("2021-03-03T09:00:32.981Z"),
		}
		actual, err := m.db.GetMySQLUsedLicenses(filter)
		m.Require().NoError(err)

		expected := []dto.MySQLUsedLicense{second}
		assert.Equal(t, expected, actual)
	})
}
