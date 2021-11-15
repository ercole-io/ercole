// Copyright (c) 2021 Sorint.lab S.p.A.
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestSearchOracleDatabaseUsedLicenses() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

	emptyResponse := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{},
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
		out, err := m.db.SearchOracleDatabaseUsedLicenses("", false, -1, -1, "Italy", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(emptyResponse), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseUsedLicenses("", false, -1, -1, "", "TEST", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(emptyResponse), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseUsedLicenses("", false, -1, -1, "", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(emptyResponse), utils.ToJSON(out))
	})

	m.T().Run("should_do_pagination", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseUsedLicenses("", false, 0, 2, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expected := dto.OracleDatabaseUsedLicenseSearchResponse{
			Content: []dto.OracleDatabaseUsedLicense{
				{Hostname: "test-db3", DbName: "foobar3", LicenseTypeID: "A90611", UsedLicenses: 0.5, Ignored: false},
				{Hostname: "test-db3", DbName: "foobar3", LicenseTypeID: "A90649", UsedLicenses: 0.5, Ignored: false},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          false,
				Number:        0,
				Size:          2,
				TotalElements: 5,
				TotalPages:    2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expected), utils.ToJSON(out))
	})

	//TODO
	//m.T().Run("should_be_sorted", func(t *testing.T) {
	//	out, err := m.db.SearchOracleDatabaseUsedLicenses("licenseName", true, -1, -1, "", "", utils.MAX_TIME)
	//	m.Require().NoError(err)

	//	var expectedOut interface{} = []interface{}{
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Real Application Clusters", "usedLicenses": 1.5},
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
	//	}

	//	assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	//})

	//m.T().Run("should_not_filter", func(t *testing.T) {
	//	out, err := m.db.SearchOracleDatabaseUsedLicenses("", false, -1, -1, "", "", utils.MAX_TIME)
	//	m.Require().NoError(err)

	//	var expectedOut []interface{} = []interface{}{
	//		//foobar3
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
	//		//foobar4
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
	//		map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Real Application Clusters", "usedLicenses": 1.5},
	//	}

	//	assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	//})
}

func (m *MongodbSuite) TestLicenseHostIgnoredField_Success() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))

	m.T().Run("update_ignored_success", func(t *testing.T) {

		hostdata := []model.HostDataBE{
			{
				Hostname: "test-db",
				Features: model.Features{
					Oracle: &model.OracleFeature{
						Database: &model.OracleDatabaseFeature{
							Databases: []model.OracleDatabase{
								{
									InstanceName: "ERCOLE1",
									Licenses: []model.OracleDatabaseLicense{
										{
											LicenseTypeID: "A90611",
											Ignored:       true,
										},
									},
								},
							},
						},
					},
				},
			},
		}

		hostname, dbname, licenseTypeID := "test-db", "ERCOLE1", "A90611"
		ignored := true

		err := m.db.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored)
		require.NoError(t, err)

		var resultIgnored bool
		for i := range hostdata[0].Features.Oracle.Database.Databases {
			db := &hostdata[0].Features.Oracle.Database.Databases[i]
			if db.InstanceName == dbname {
				for j := range db.Licenses {
					lic := &db.Licenses[j]
					if lic.LicenseTypeID == licenseTypeID {
						resultIgnored = lic.Ignored
					}
				}
			}
		}

		require.Equal(t, ignored, resultIgnored)

		err = m.db.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, !ignored)
		require.NoError(t, err)
	})
}

func (m *MongodbSuite) TestLicenseHostIgnoredField_Fail() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.T().Run("update_ignored_fail", func(t *testing.T) {

		hostname, dbname, licenseTypeID := "buu", "ERCOLAO", "BBBIIII"
		ignored := false

		err := m.db.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored)
		require.Error(t, err)
	})
}
