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

package database

import (
	"context"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchSqlServerDatabaseUsedLicenses() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_28.json"))

	defer m.db.Client.Database(m.dbname).Collection("ms_sqlserver_database_license_types").DeleteMany(context.TODO(), bson.M{})
	m.db.InsertSqlServerDatabaseLicenseType(model.SqlServerDatabaseLicenseType{
		ID:              "123-45678",
		ItemDescription: "Enterproise Edition",
		Edition:         "ENT",
		Version:         "2019",
	})

	emptyResponse := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{},
		Metadata: dto.PagingMetadata{
			Empty:         true,
			First:         true,
			Last:          true,
			Number:        0,
			Size:          0,
			TotalElements: 0,
			TotalPages:    0,
		},
	}

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "Italy", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(emptyResponse), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "TEST", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(emptyResponse), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(emptyResponse), utils.ToJSON(out))
	})

	m.T().Run("should_do_pagination", func(t *testing.T) {
		out, err := m.db.SearchSqlServerDatabaseUsedLicenses("", "", false, 0, 2, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expected := dto.SqlServerDatabaseUsedLicenseSearchResponse{
			Content: []dto.SqlServerDatabaseUsedLicense{
				{Hostname: "test-db2", DbName: "MSSQLSERVER", LicenseTypeID: "123-45678", UsedLicenses: 2, Ignored: false},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          1,
				TotalElements: 1,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expected), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorted", func(t *testing.T) {
		out, err := m.db.SearchSqlServerDatabaseUsedLicenses("", "licenseName", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expected := dto.SqlServerDatabaseUsedLicenseSearchResponse{
			Content: []dto.SqlServerDatabaseUsedLicense{
				{Hostname: "test-db2", DbName: "MSSQLSERVER", LicenseTypeID: "123-45678", UsedLicenses: 2, Ignored: false},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          1,
				TotalElements: 1,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expected), utils.ToJSON(out))
	})

	m.T().Run("should_not_filter", func(t *testing.T) {
		out, err := m.db.SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expected := dto.SqlServerDatabaseUsedLicenseSearchResponse{
			Content: []dto.SqlServerDatabaseUsedLicense{
				{Hostname: "test-db2", DbName: "MSSQLSERVER", LicenseTypeID: "123-45678", UsedLicenses: 2, Ignored: false},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          1,
				TotalElements: 1,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expected), utils.ToJSON(out))
	})

	m.T().Run("should_filter_by_hostname", func(t *testing.T) {
		out, err := m.db.SearchSqlServerDatabaseUsedLicenses("test-db3", "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expected := dto.SqlServerDatabaseUsedLicenseSearchResponse{
			Content: []dto.SqlServerDatabaseUsedLicense{},
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expected), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestSqlServerLicenseHostIgnoredField_Success() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_28.json"))

	m.T().Run("update_ignored_success", func(t *testing.T) {

		hostname, instancename := "test-db2", "MSSQLSERVER"
		ignored := true

		err := m.db.UpdateSqlServerLicenseIgnoredField(hostname, instancename, ignored)
		require.NoError(t, err)

		hostData, _ := m.db.FindHostData("test-db2")

		var resultIgnored bool
		for i := range hostData.Features.Microsoft.SQLServer.Instances {
			db := &hostData.Features.Microsoft.SQLServer.Instances[i]
			if db.Name == instancename {
				lic := &db.License
				resultIgnored = lic.Ignored
			}
		}

		require.Equal(t, ignored, resultIgnored)
	})
}

func (m *MongodbSuite) TestSqlServerLicenseHostIgnoredField_Fail() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.T().Run("update_ignored_fail", func(t *testing.T) {

		hostname, instancename := "buu", "ERCOLAO"
		ignored := false

		err := m.db.UpdateSqlServerLicenseIgnoredField(hostname, instancename, ignored)
		require.Error(t, err)
	})
}
