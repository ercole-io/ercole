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

	"github.com/amreo/ercole-services/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetDatabaseEnvironmentStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetDatabaseEnvironmentStats("France", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetDatabaseEnvironmentStats("", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetDatabaseEnvironmentStats("", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Environment": "PRD",
				"Count":       2,
			},
			{
				"Environment": "TST",
				"Count":       2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetDatabaseVersionStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetDatabaseVersionStats("France", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetDatabaseVersionStats("", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetDatabaseVersionStats("", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Version": "12.2.0.1.0 Enterprise Edition",
				"Count":   2,
			},
			{
				"Version": "16.2.0.1.0 Enterprise Edition",
				"Count":   1,
			},
			{
				"Version": "18.2.0.1.0 Enterprise Edition",
				"Count":   1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetTopReclaimableDatabaseStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTopReclaimableDatabaseStats("France", 15, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTopReclaimableDatabaseStats("", 15, utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_limit_the_result", func(t *testing.T) {
		out, err := m.db.GetTopReclaimableDatabaseStats("", 1, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Hostname":                   "test-db3",
				"Dbname":                     "foobar4",
				"ReclaimableSegmentAdvisors": 534.34,
				"_id":                        "5ec2518bbc4991e955e2cb3f",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_all_results", func(t *testing.T) {
		out, err := m.db.GetTopReclaimableDatabaseStats("", 15, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Hostname":                   "test-db3",
				"Dbname":                     "foobar4",
				"ReclaimableSegmentAdvisors": 534.34,
				"_id":                        "5ec2518bbc4991e955e2cb3f",
			},
			{
				"Hostname":                   "test-db3",
				"Dbname":                     "foobar3",
				"ReclaimableSegmentAdvisors": 4.3,
				"_id":                        "5ec2518bbc4991e955e2cb3f",
			},
			{
				"Hostname":                   "test-db2",
				"Dbname":                     "foobar1",
				"ReclaimableSegmentAdvisors": 0.5,
				"_id":                        "5ebbaaf747c3fcf9dc0a1f51",
			},
			{
				"Hostname":                   "test-db2",
				"Dbname":                     "foobar2",
				"ReclaimableSegmentAdvisors": 0.5,
				"_id":                        "5ebbaaf747c3fcf9dc0a1f51",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetTopWorkloadDatabaseStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTopWorkloadDatabaseStats("France", 15, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTopWorkloadDatabaseStats("", 15, utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_limit_the_result", func(t *testing.T) {
		out, err := m.db.GetTopWorkloadDatabaseStats("", 1, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Hostname": "test-db3",
				"Dbname":   "foobar3",
				"Workload": 99,
				"_id":      "5ec2518bbc4991e955e2cb3f",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_all_results", func(t *testing.T) {
		out, err := m.db.GetTopWorkloadDatabaseStats("", 15, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Hostname": "test-db3",
				"Dbname":   "foobar3",
				"Workload": 99,
				"_id":      "5ec2518bbc4991e955e2cb3f",
			},
			{
				"Hostname": "test-db3",
				"Dbname":   "foobar4",
				"Workload": 10,
				"_id":      "5ec2518bbc4991e955e2cb3f",
			},
			{
				"Hostname": "test-db2",
				"Dbname":   "foobar2",
				"Workload": 6.4,
				"_id":      "5ebbaaf747c3fcf9dc0a1f51",
			},
			{
				"Hostname": "test-db2",
				"Dbname":   "foobar1",
				"Workload": 1,
				"_id":      "5ebbaaf747c3fcf9dc0a1f51",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetDatabasePatchStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetDatabasePatchStatusStats("France", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetDatabasePatchStatusStats("", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetDatabasePatchStatusStats("", utils.P("2019-10-10T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Status": "KO",
				"Count":  3,
			},
			{
				"Status": "OK",
				"Count":  1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetDatabaseDataguardStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetDatabaseDataguardStatusStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetDatabaseDataguardStatusStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetDatabaseDataguardStatusStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetDatabaseDataguardStatusStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Dataguard": false,
				"Count":     3,
			},
			{
				"Dataguard": true,
				"Count":     1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetDatabaseRACStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetDatabaseRACStatusStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetDatabaseRACStatusStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetDatabaseRACStatusStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetDatabaseRACStatusStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"RAC":   false,
				"Count": 3,
			},
			{
				"RAC":   true,
				"Count": 1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetDatabaseArchivelogStatusStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetDatabaseArchivelogStatusStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetDatabaseArchivelogStatusStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetDatabaseArchivelogStatusStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_result", func(t *testing.T) {
		out, err := m.db.GetDatabaseArchivelogStatusStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"Archivelog": false,
				"Count":      3,
			},
			{
				"Archivelog": true,
				"Count":      1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetTotalDatabaseWorkStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseWorkStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float32(0.0), out)
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseWorkStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseWorkStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseWorkStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-116.4) < 0.00001)
	})
}

func (m *MongodbSuite) TestGetTotalDatabaseMemorySizeStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseMemorySizeStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float32(0.0), out)
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseMemorySizeStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseMemorySizeStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseMemorySizeStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-34.642) < 0.00001)
	})
}

func (m *MongodbSuite) TestGetTotalDatabaseSegmentSizeStats() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseSegmentSizeStats("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.Equal(t, float32(0.0), out)
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseSegmentSizeStats("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseSegmentSizeStats("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-0.0) < 0.00001)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.GetTotalDatabaseSegmentSizeStats("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.True(t, math.Abs(float64(out)-48) < 0.00001)
	})
}
