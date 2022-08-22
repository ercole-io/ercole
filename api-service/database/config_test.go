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

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
)

func (m *MongodbSuite) TestFindConfig() {
	defer m.db.Client.Database(m.dbname).Collection("config").DeleteMany(context.TODO(), bson.M{})
	m.db.Client.Database(m.dbname).Collection("config").InsertOne(context.TODO(), config.Configuration{})

	m.T().Run("find_element", func(t *testing.T) {
		out, err := m.db.FindConfig()
		m.Require().NoError(err)
		expectedRes := config.Configuration{}

		assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(out))
	})

	m.T().Run("change_element", func(t *testing.T) {
		expectedOut := config.Configuration{
			APIService: config.APIService{Port: 9999},
		}
		err := m.db.ChangeConfig(expectedOut)
		m.Require().NoError(err)

		out, err := m.db.FindConfig()
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
