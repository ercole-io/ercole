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

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetEnvironmentStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetEnvironmentStats("France", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetEnvironmentStats("", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetEnvironmentStats("", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"environment": "DEV",
				"count":       1,
			},
			{
				"environment": "PROD",
				"count":       1,
			},
			{
				"environment": "TST",
				"count":       2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetTypeStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTypeStats("France", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTypeStats("", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTypeStats("", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"hardwareAbstractionTechnology": "OVM",
				"count":                         1,
			},
			{
				"hardwareAbstractionTechnology": "PH",
				"count":                         1,
			},
			{
				"hardwareAbstractionTechnology": "VMWARE",
				"count":                         2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetOperatingSystemStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_11.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		m.db.OperatingSystemAggregationRules = []config.AggregationRule{}
		out, err := m.db.GetOperatingSystemStats("France", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		m.db.OperatingSystemAggregationRules = []config.AggregationRule{}
		out, err := m.db.GetOperatingSystemStats("", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		m.db.OperatingSystemAggregationRules = []config.AggregationRule{}
		out, err := m.db.GetOperatingSystemStats("", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"operatingSystem": "Red Hat Enterprise Linux 7.6",
				"count":           3,
			},
			{
				"operatingSystem": "Ubuntu 18.04.4",
				"count":           1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_aggregate_correctly", func(t *testing.T) {
		m.db.OperatingSystemAggregationRules = []config.AggregationRule{
			{
				Regex: "^Red Hat Enterprise Linux 7.*$",
				Group: "RHEL7",
			},
			{
				Regex: "^Ubuntu 18\\.04.*$",
				Group: "Ubuntu Server",
			},
		}
		out, err := m.db.GetOperatingSystemStats("", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"operatingSystem": "RHEL7",
				"count":           3,
			},
			{
				"operatingSystem": "Ubuntu Server",
				"count":           1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
