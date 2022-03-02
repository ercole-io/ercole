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

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestGetHostCores() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_chartservice_mongohostdata_01.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_chartservice_mongohostdata_02.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		location := ""
		environment := "TST"
		olderThan := utils.MAX_TIME
		newerThan := utils.MIN_TIME

		out, err := m.db.GetHostCores(location, environment, olderThan, newerThan)
		m.Require().NoError(err)
		expectedOut := []dto.HostCores{
			{
				Date:  utils.P("2020-04-15T00:00:00Z"),
				Cores: 1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		location := "Germany"
		environment := ""
		olderThan := utils.MAX_TIME
		newerThan := utils.MIN_TIME

		out, err := m.db.GetHostCores(location, environment, olderThan, newerThan)
		m.Require().NoError(err)
		expectedOut := []dto.HostCores{
			{
				Date:  utils.P("2020-04-15T00:00:00Z"),
				Cores: 1,
			},
			{
				Date:  utils.P("2020-05-13T00:00:00Z"),
				Cores: 1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		location := ""
		environment := ""
		olderThan := utils.P("2020-04-16T00:00:00Z")
		newerThan := utils.MIN_TIME

		out, err := m.db.GetHostCores(location, environment, olderThan, newerThan)
		m.Require().NoError(err)
		expectedOut := []dto.HostCores{
			{
				Date:  utils.P("2020-04-15T00:00:00Z"),
				Cores: 1,
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_newer_than", func(t *testing.T) {
		location := ""
		environment := ""
		olderThan := utils.MAX_TIME
		newerThan := utils.P("2020-05-12T00:00:00Z")

		out, err := m.db.GetHostCores(location, environment, olderThan, newerThan)
		m.Require().NoError(err)
		expectedOut := []dto.HostCores{
			{
				Date:  utils.P("2020-05-13T00:00:00Z"),
				Cores: 1,
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
