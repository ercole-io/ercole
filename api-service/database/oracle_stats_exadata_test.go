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

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetTotalOracleExadataMemorySizeStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataMemorySizeStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataMemorySizeStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float64(0.0), out)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataMemorySizeStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataMemorySizeStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-1128) < 0.00001)
	})
}

func (m *MongodbSuite) TestGetTotalOracleExadataCPUStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataCPUStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"running": 0,
			"total":   0,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataCPUStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"running": 0,
			"total":   0,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataCPUStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"running": 0,
			"total":   0,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTotalOracleExadataCPUStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			"running": 136,
			"total":   176,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetAverageOracleExadataStorageUsageStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.T().Run("should_return_zero_when_collection_is_empty", func(t *testing.T) {
		out, err := m.db.GetAverageOracleExadataStorageUsageStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float64(0.0), out)
	})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))

	m.T().Run("should_return_zero_when_no_exadata_is_present", func(t *testing.T) {
		out, err := m.db.GetAverageOracleExadataStorageUsageStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float64(0.0), out)
	})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetAverageOracleExadataStorageUsageStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetAverageOracleExadataStorageUsageStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float64(0.0), out)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetAverageOracleExadataStorageUsageStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetAverageOracleExadataStorageUsageStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-52.4) < 0.00001)
	})
}

func (m *MongodbSuite) TestGetOracleExadataStorageErrorCountStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetOracleExadataStorageErrorCountStatusStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetOracleExadataStorageErrorCountStatusStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetOracleExadataStorageErrorCountStatusStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetOracleExadataStorageErrorCountStatusStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"failing": false,
				"count":   3,
			},
			{
				"failing": true,
				"count":   2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetOracleExadataPatchStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetOracleExadataPatchStatusStats("France", "", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetOracleExadataPatchStatusStats("", "FOOBAR", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetOracleExadataPatchStatusStats("", "", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetOracleExadataPatchStatusStats("", "", utils.P("2019-06-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"count":  2,
				"status": false,
			},
			{
				"count":  4,
				"status": true,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
