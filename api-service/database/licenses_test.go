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
	"fmt"
	"testing"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestListLicenses() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))
	m.InsertLicense(model.LicenseCount{Name: "Oracle ENT"})
	m.InsertLicense(model.LicenseCount{Name: "Diagnostics Pack"})
	m.InsertLicense(model.LicenseCount{Name: "Real Application Clusters"})

	m.T().Run("should_not_filter", func(t *testing.T) {
		mode := "summary"
		sortBy := ""
		sortDesc := false
		page := -1
		pageSize := -1
		location := ""
		environment := ""
		olderThan := utils.MAX_TIME
		out, err := m.db.ListLicenses(mode, sortBy, sortDesc, page, pageSize, location, environment, olderThan)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{
			map[string]interface{}{"_id": "Oracle ENT", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
			map[string]interface{}{"_id": "Real Application Clusters", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 2},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_not_filter_full", func(t *testing.T) {
		mode := "full"
		sortBy := ""
		sortDesc := false
		page := -1
		pageSize := -1
		location := ""
		environment := ""
		olderThan := utils.MAX_TIME
		out, err := m.db.ListLicenses(mode, sortBy, sortDesc, page, pageSize, location, environment, olderThan)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{
			map[string]interface{}{"_id": "Oracle ENT", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1,
				"hosts": []interface{}{map[string]interface{}{"databases": []interface{}{"foobar3", "foobar4"}, "hostname": "test-db3"}},
			},
			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1,
				"hosts": []interface{}{map[string]interface{}{"databases": []interface{}{"foobar3", "foobar4"}, "hostname": "test-db3"}},
			},
			map[string]interface{}{"_id": "Real Application Clusters", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 2,
				"hosts": []interface{}{map[string]interface{}{"databases": []interface{}{"foobar4"}, "hostname": "test-db3"}},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_raise_error_wrong_mode", func(t *testing.T) {
		mode := "this_is_a_wrong_mode"
		sortBy := ""
		sortDesc := false
		page := -1
		pageSize := -1
		location := ""
		environment := ""
		olderThan := utils.MAX_TIME
		_, err := m.db.ListLicenses(mode, sortBy, sortDesc, page, pageSize, location, environment, olderThan)

		m.Require().Error(err)
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		mode := "summary"
		sortBy := ""
		sortDesc := false
		page := 0
		pageSize := 1
		location := ""
		environment := ""
		olderThan := utils.MAX_TIME
		out, err := m.db.ListLicenses(mode, sortBy, sortDesc, page, pageSize, location, environment, olderThan)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{"_id": "Oracle ENT", "compliance": false, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 1},
				},
				"metadata": map[string]interface{}{
					"empty":         false,
					"first":         true,
					"last":          false,
					"number":        0,
					"size":          1,
					"totalElements": 3,
					"totalPages":    3},
			}}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		mode := "summary"
		sortBy := ""
		sortDesc := false
		page := -1
		pageSize := -1
		location := "Italy"
		environment := ""
		olderThan := utils.MAX_TIME
		out, err := m.db.ListLicenses(mode, sortBy, sortDesc, page, pageSize, location, environment, olderThan)

		m.Require().NoError(err)
		var expectedOut []interface{} = []interface{}{
			map[string]interface{}{"_id": "Oracle ENT", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
			map[string]interface{}{"_id": "Real Application Clusters", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0}}

		fmt.Printf("%v\n", out)
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		mode := "summary"
		sortBy := ""
		sortDesc := false
		page := -1
		pageSize := -1
		location := ""
		environment := "TEST"
		olderThan := utils.MAX_TIME
		out, err := m.db.ListLicenses(mode, sortBy, sortDesc, page, pageSize, location, environment, olderThan)

		m.Require().NoError(err)
		var expectedOut []interface{} = []interface{}{
			map[string]interface{}{"_id": "Oracle ENT", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
			map[string]interface{}{"_id": "Real Application Clusters", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0}}

		fmt.Printf("%v\n", out)
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		mode := "summary"
		sortBy := ""
		sortDesc := false
		page := -1
		pageSize := -1
		location := ""
		environment := ""
		olderThan := utils.MIN_TIME
		out, err := m.db.ListLicenses(mode, sortBy, sortDesc, page, pageSize, location, environment, olderThan)

		m.Require().NoError(err)
		var expectedOut []interface{} = []interface{}{
			map[string]interface{}{"_id": "Oracle ENT", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
			map[string]interface{}{"_id": "Diagnostics Pack", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0},
			map[string]interface{}{"_id": "Real Application Clusters", "compliance": true, "costPerProcessor": 0, "count": 0, "paidCost": 0, "totalCost": 0, "unlimited": false, "used": 0}}

		fmt.Printf("%v\n", out)
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

}
