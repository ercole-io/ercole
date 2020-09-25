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

	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO
// func (m *MongodbSuite) TestSearchLicenses() {
// 	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
// 	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

// 	defer m.db.Client.Database(m.dbname).Collection("licenses").DeleteMany(context.TODO(), bson.M{})
// 	m.InsertLicense(model.LicenseCount{Name: "Oracle ENT"})
// 	m.InsertLicense(model.LicenseCount{Name: "Diagnostics Pack"})
// 	m.InsertLicense(model.LicenseCount{Name: "Real Application Clusters"})

// 	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("Italy", "", utils.MAX_TIME)

// 		m.Require().NoError(err)
// 		var expectedOut []interface{} = []interface{}{
// 			map[string]interface{}{"_id": "Oracle ENT", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
// 			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
// 			map[string]interface{}{"_id": "Real Application Clusters", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0}}

// 		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
// 	})

// 	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("TEST", utils.MAX_TIME)

// 		m.Require().NoError(err)
// 		var expectedOut []interface{} = []interface{}{
// 			map[string]interface{}{"_id": "Oracle ENT", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
// 			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
// 			map[string]interface{}{"_id": "Real Application Clusters", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0}}

// 		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
// 	})

// 	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("", "", utils.MIN_TIME)

// 		m.Require().NoError(err)
// 		var expectedOut []interface{} = []interface{}{
// 			map[string]interface{}{"_id": "Oracle ENT", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
// 			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
// 			map[string]interface{}{"_id": "Real Application Clusters", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0}}

// 		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
// 	})

// 	m.T().Run("should_do_pagination", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("", "", utils.MAX_TIME)
// 		m.Require().NoError(err)

// 		var expectedOut interface{} = []interface{}{
// 			map[string]interface{}{
// 				"content": []interface{}{
// 					map[string]interface{}{"_id": "Oracle ENT", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
// 				},
// 				"metadata": map[string]interface{}{
// 					"empty":         false,
// 					"first":         true,
// 					"last":          false,
// 					"number":        0,
// 					"size":          1,
// 					"totalElements": 3,
// 					"totalPages":    3},
// 			}}

// 		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
// 	})

// 	m.T().Run("should_be_sorted", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("summary", "used", true, -1, -1, "", "", utils.MAX_TIME)
// 		m.Require().NoError(err)

// 		var expectedOut interface{} = []interface{}{
// 			map[string]interface{}{"_id": "Real Application Clusters", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 2},
// 			map[string]interface{}{"_id": "Oracle ENT", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
// 			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
// 		}

// 		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
// 	})

// 	m.T().Run("should_not_filter", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("summary", "", false, -1, -1, "", "", utils.MAX_TIME)
// 		m.Require().NoError(err)

// 		var expectedOut interface{} = []interface{}{
// 			map[string]interface{}{"_id": "Oracle ENT", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
// 			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
// 			map[string]interface{}{"_id": "Real Application Clusters", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 2},
// 		}

// 		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
// 	})

// 	m.T().Run("should_not_filter_full", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("full", "", false, -1, -1, "", "", utils.MAX_TIME)
// 		m.Require().NoError(err)

// 		var expectedOut interface{} = []interface{}{
// 			map[string]interface{}{"_id": "Oracle ENT", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1,
// 				"hosts": []interface{}{map[string]interface{}{"databases": []interface{}{"foobar3", "foobar4"}, "hostname": "test-db3"}},
// 			},
// 			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1,
// 				"hosts": []interface{}{map[string]interface{}{"databases": []interface{}{"foobar3", "foobar4"}, "hostname": "test-db3"}},
// 			},
// 			map[string]interface{}{"_id": "Real Application Clusters", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 2,
// 				"hosts": []interface{}{map[string]interface{}{"databases": []interface{}{"foobar4"}, "hostname": "test-db3"}},
// 			},
// 		}

// 		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
// 	})

// 	m.T().Run("should_raise_error_wrong_mode", func(t *testing.T) {
// 		out, err := m.db.SearchLicenses("wrong_mode", "", false, -1, -1, "", "", utils.MAX_TIME)

// 		m.Require().Nil(out)
// 		m.Require().Error(err)
// 	})
// }

func (m *MongodbSuite) TestListLicenses() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.ListLicenses("", false, -1, -1, "Italy", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.ListLicenses("", false, -1, -1, "", "TEST", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.ListLicenses("", false, -1, -1, "", "", utils.MIN_TIME)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_do_pagination", func(t *testing.T) {
		out, err := m.db.ListLicenses("", false, 0, 2, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
					map[string]interface{}{"dbName": "foobar3", "hostname": "test-db3", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
				},
				"metadata": map[string]interface{}{
					"empty":         false,
					"first":         true,
					"last":          false,
					"number":        0,
					"size":          2,
					"totalElements": 5,
					"totalPages":    2},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorted", func(t *testing.T) {
		out, err := m.db.ListLicenses("licenseName", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Real Application Clusters", "usedLicenses": 1.5},
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_not_filter", func(t *testing.T) {
		out, err := m.db.ListLicenses("", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedOut []interface{} = []interface{}{
			//foobar3
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar3", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
			//foobar4
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Oracle ENT", "usedLicenses": 0.5},
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Diagnostics Pack", "usedLicenses": 0.5},
			map[string]interface{}{"hostname": "test-db3", "dbName": "foobar4", "licenseName": "Real Application Clusters", "usedLicenses": 1.5},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
