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

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongodbSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test for mongodb database(api-service)")
	}

	mongodbHandlerSuiteTest := &MongodbSuite{}

	suite.Run(t, mongodbHandlerSuiteTest)
}

func (m *MongodbSuite) TestFilterByLocationAndEnvironmentSteps() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_01.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_05.json"))

	m.RunTestQuery(
		"no_location_no_environment",
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps("", ""),
			mu.APProject(bson.M{
				"Hostname": 1,
				"_id":      0,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut string = `[
				{ "Hostname": "test-small" },
				{ "Hostname": "test-small2" },
				{ "Hostname": "test-small3" }
			]`

			assert.JSONEq(m.T(), expectedOut, utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"location_no_environment",
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps("Italy", ""),
			mu.APProject(bson.M{
				"Hostname": 1,
				"_id":      0,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut string = `[
				{ "Hostname": "test-small2" },
				{ "Hostname": "test-small3" }
			]`

			assert.JSONEq(m.T(), expectedOut, utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"no_location_environment",
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps("", "TST"),
			mu.APProject(bson.M{
				"Hostname": 1,
				"_id":      0,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut string = `[
				{ "Hostname": "test-small" },
				{ "Hostname": "test-small3" }
			]`

			assert.JSONEq(m.T(), expectedOut, utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"location_environment",
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps("Italy", "DEV"),
			mu.APProject(bson.M{
				"Hostname": 1,
				"_id":      0,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut string = `[
				{ "Hostname": "test-small2" }
			]`

			assert.JSONEq(m.T(), expectedOut, utils.ToJSON(out))
		},
	)
}

func (m *MongodbSuite) TestFilterByOldnessSteps() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_01.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_05.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_06.json"))

	m.RunTestQuery(
		"latest_mongohostdata",
		mu.MAPipeline(
			FilterByOldnessSteps(utils.MAX_TIME),
			mu.APProject(bson.M{
				"_id": 1,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut interface{} = []interface{}{
				map[string]interface{}{"_id": utils.Str2oid("5ea2d26d20d55cbdc35022b4")},
				map[string]interface{}{"_id": utils.Str2oid("5ea2dec620d55cbdc35022b9")},
			}

			assert.JSONEq(m.T(), utils.ToJSON(expectedOut), utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"no_mongohostdata_min_time",
		mu.MAPipeline(
			FilterByOldnessSteps(utils.MIN_TIME),
			mu.APProject(bson.M{
				"_id": 1,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut interface{} = []interface{}{}

			assert.JSONEq(m.T(), utils.ToJSON(expectedOut), utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"all_mongohostdata_latest",
		mu.MAPipeline(
			FilterByOldnessSteps(utils.P("2020-04-24T12:46:36+00:00")),
			mu.APProject(bson.M{
				"_id": 1,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut interface{} = []interface{}{
				map[string]interface{}{"_id": utils.Str2oid("5ea2d26d20d55cbdc35022b4")},
				map[string]interface{}{"_id": utils.Str2oid("5ea2dec620d55cbdc35022b9")},
			}

			assert.JSONEq(m.T(), utils.ToJSON(expectedOut), utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"no_mongohostdata_too_early",
		mu.MAPipeline(
			FilterByOldnessSteps(utils.P("2020-04-24T12:46:36+02:00")),
			mu.APProject(bson.M{
				"_id": 1,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut interface{} = []interface{}{}

			assert.JSONEq(m.T(), utils.ToJSON(expectedOut), utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"only_latest_test_small_hostdata",
		mu.MAPipeline(
			FilterByOldnessSteps(utils.P("2020-04-24T13:50:36+02:00")),
			mu.APProject(bson.M{
				"_id": 1,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut interface{} = []interface{}{
				map[string]interface{}{"_id": utils.Str2oid("5ea2d26d20d55cbdc35022b4")},
			}

			assert.JSONEq(m.T(), utils.ToJSON(expectedOut), utils.ToJSON(out))
		},
	)

	m.RunTestQuery(
		"not_latest_test_small2",
		mu.MAPipeline(
			FilterByOldnessSteps(utils.P("2020-04-24T13:58:36+02:00")),
			mu.APProject(bson.M{
				"_id": 1,
			}),
		),
		func(out []map[string]interface{}) {
			var expectedOut interface{} = []interface{}{
				map[string]interface{}{"_id": utils.Str2oid("5ea2d26d20d55cbdc35022b4")},
				map[string]interface{}{"_id": utils.Str2oid("5ea2d3c920d55cbdc35022b7")},
			}

			assert.JSONEq(m.T(), utils.ToJSON(expectedOut), utils.ToJSON(out))
		},
	)
}
