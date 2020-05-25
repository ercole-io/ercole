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
	"math"
	"testing"

	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetTotalExadataMemorySizeStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetTotalExadataMemorySizeStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTotalExadataMemorySizeStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float32(0.0), out)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTotalExadataMemorySizeStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTotalExadataMemorySizeStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-1128) < 0.00001)
	})
}

func (m *MongodbSuite) TestGetTotalExadataCPUStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTotalExadataCPUStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"Enabled": 0,
			"Total":   0,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetTotalExadataCPUStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"Enabled": 0,
			"Total":   0,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTotalExadataCPUStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"Enabled": 0,
			"Total":   0,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTotalExadataCPUStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"Enabled": 136,
			"Total":   176,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetAverageExadataStorageUsageStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetAverageExadataStorageUsageStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetAverageExadataStorageUsageStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float32(0.0), out)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetAverageExadataStorageUsageStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetAverageExadataStorageUsageStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-52.4) < 0.00001)
	})
}

func (m *MongodbSuite) TestGetExadataStorageErrorCountStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetExadataStorageErrorCountStatusStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetExadataStorageErrorCountStatusStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetExadataStorageErrorCountStatusStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetExadataStorageErrorCountStatusStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Failing": false,
				"Count":   3,
			},
			{
				"Failing": true,
				"Count":   2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetExadataPatchStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetExadataPatchStatusStats("France", "", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetExadataPatchStatusStats("", "FOOBAR", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetExadataPatchStatusStats("", "", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetExadataPatchStatusStats("", "", utils.P("2019-06-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Count":  2,
				"Status": false,
			},
			{
				"Count":  4,
				"Status": true,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
