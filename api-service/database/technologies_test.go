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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetHostsCountUsingTechnologies() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_22.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.GetHostsCountUsingTechnologies("Foobarland", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.GetHostsCountUsingTechnologies("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.GetHostsCountUsingTechnologies("", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_return_correct_res", func(t *testing.T) {
		out, err := m.db.GetHostsCountUsingTechnologies("", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = map[string]interface{}{
			model.TechnologyOracleExadata:  1,
			model.TechnologyOracleDatabase: 2,
			model.TechnologyOracleMySQL:    1,
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
