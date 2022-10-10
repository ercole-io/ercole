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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetRoles() {
	m.T().Run("should_load_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(roleCollection).DeleteMany(context.TODO(), bson.M{})

		roles := []model.Role{
			{
				Name: "Test",
			},
			{
				Name: "Testdue",
			},
		}
		rolesInt := []interface{}{
			roles[0],
			roles[1],
		}
		_, err := m.db.Client.Database(m.dbname).Collection(roleCollection).InsertMany(context.TODO(), rolesInt)
		require.Nil(m.T(), err)

		actual, err := m.db.GetRoles()
		m.Require().NoError(err)

		assert.Equal(t, roles, actual)
	})

	m.T().Run("should_load_empty", func(t *testing.T) {
		actual, err := m.db.GetRoles()
		m.Require().NoError(err)

		roles := make([]model.Role, 0)
		assert.Equal(t, roles, actual)
	})
}

func (m *MongodbSuite) TestGetRole() {
	m.T().Run("should_load_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(roleCollection).DeleteMany(context.TODO(), bson.M{})

		role := model.Role{
			Name: "Test",
		}

		roles := []model.Role{
			{
				Name: "Test",
			},
			{
				Name: "Testdue",
			},
		}
		rolesInt := []interface{}{
			roles[0],
			roles[1],
		}
		_, err := m.db.Client.Database(m.dbname).Collection(roleCollection).InsertMany(context.TODO(), rolesInt)
		require.Nil(m.T(), err)

		actual, err := m.db.GetRole("Test")
		m.Require().NoError(err)

		assert.Equal(t, &role, actual)
	})

	m.T().Run("should_load_empty", func(t *testing.T) {
		actual, err := m.db.GetRoles()
		m.Require().NoError(err)

		roles := make([]model.Role, 0)
		assert.Equal(t, roles, actual)
	})
}
