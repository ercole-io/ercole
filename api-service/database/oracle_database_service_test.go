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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetOracleServiceList() {
	name := "ECXSERVER"
	enabled := false

	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.T().Run("success with empty table", func(t *testing.T) {
		actual, err := m.db.GetOracleServiceList(dto.GlobalFilter{OlderThan: utils.MAX_TIME})
		m.Require().NoError(err)

		expected := make([]dto.OracleDatabaseServiceDto, 0)
		assert.ElementsMatch(m.T(), expected, actual)
	})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_33.json"))

	m.T().Run("success with some values", func(t *testing.T) {
		expected := []interface{}{
			dto.OracleDatabaseServiceDto{
				Hostname:     "newdb",
				Databasename: "pippodb",
				OracleDatabaseService: model.OracleDatabaseService{
					Name:            &name,
					Enabled:         &enabled,
					FailoverMethod:  nil,
					FailoverType:    nil,
					FailoverRetries: nil,
					FailoverDelay:   "",
				},
			},
		}

		ctx := context.TODO()
		_, err := m.db.Client.Database(m.dbname).
			Collection("hosts").
			InsertMany(ctx, expected)
		m.Require().NoError(err)

		actual, err := m.db.GetOracleServiceList(dto.GlobalFilter{OlderThan: utils.MAX_TIME})
		m.Require().NoError(err)

		assert.ElementsMatch(m.T(), expected, actual)
	})

}
