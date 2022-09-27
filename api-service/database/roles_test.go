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
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestInsertRole() {
	defer m.db.Client.Database(m.dbname).Collection(roleCollection).DeleteMany(context.TODO(), bson.M{})

	role := model.RoleType{
		ID:   utils.Str2oid("000000000000000000000001"),
		Name: "Test",
	}

	m.T().Run("should_insert", func(t *testing.T) {
		err := m.db.InsertRole(role)
		m.Require().NoError(err)
	})
}

func (m *MongodbSuite) TestUpdateRole() {
	defer m.db.Client.Database(m.dbname).Collection(roleCollection).DeleteMany(context.TODO(), bson.M{})

	role := model.RoleType{
		ID:   utils.Str2oid("000000000000000000000001"),
		Name: "Test",
	}

	m.T().Run("error not found", func(t *testing.T) {
		err := m.db.UpdateRole(role)
		var aerr *utils.AdvancedError
		assert.ErrorAs(t, err, &aerr)
		assert.ErrorIs(t, aerr.Err, utils.ErrNotFound)
	})

	m.T().Run("should_update", func(t *testing.T) {
		_, err := m.db.Client.Database(m.dbname).Collection(roleCollection).
			InsertOne(
				context.TODO(),
				role,
			)
		require.NoError(t, err)

		err = m.db.UpdateRole(role)
		assert.NoError(t, err)
	})
}

func (m *MongodbSuite) TestGetRoles() {
	m.T().Run("should_load_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(roleCollection).DeleteMany(context.TODO(), bson.M{})

		roles := []model.RoleType{
			{
				ID:   utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				Name: "Test",
			},
			{
				ID:   utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
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

		roles := make([]model.RoleType, 0)
		assert.Equal(t, roles, actual)
	})
}

func (m *MongodbSuite) TestGetRole() {
	m.T().Run("should_load_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(roleCollection).DeleteMany(context.TODO(), bson.M{})

		role := model.RoleType{
			ID:   utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			Name: "Test",
		}

		roles := []model.RoleType{
			{
				ID:   utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				Name: "Test",
			},
			{
				ID:   utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				Name: "Testdue",
			},
		}
		rolesInt := []interface{}{
			roles[0],
			roles[1],
		}
		_, err := m.db.Client.Database(m.dbname).Collection(roleCollection).InsertMany(context.TODO(), rolesInt)
		require.Nil(m.T(), err)

		actual, err := m.db.GetRole(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"))
		m.Require().NoError(err)

		assert.Equal(t, &role, actual)
	})

	m.T().Run("should_load_empty", func(t *testing.T) {
		actual, err := m.db.GetRoles()
		m.Require().NoError(err)

		roles := make([]model.RoleType, 0)
		assert.Equal(t, roles, actual)
	})
}

func (m *MongodbSuite) TestDeleteRole() {
	defer m.db.Client.Database(m.dbname).Collection(roleCollection).DeleteMany(context.TODO(), bson.M{})

	id := utils.Str2oid("000000000000000000000001")

	m.T().Run("error not found", func(t *testing.T) {
		err := m.db.DeleteRole(id)
		var aerr *utils.AdvancedError
		assert.ErrorAs(t, err, &aerr)
		assert.ErrorIs(t, aerr.Err, utils.ErrRoleNotFound)
	})

	m.T().Run("should_delete", func(t *testing.T) {
		role := model.RoleType{
			ID:   utils.Str2oid("000000000000000000000001"),
			Name: "Test",
		}
		_, err := m.db.Client.Database(m.dbname).Collection(roleCollection).
			InsertOne(
				context.TODO(),
				role,
			)
		require.NoError(t, err)

		err = m.db.DeleteRole(role.ID)
		assert.NoError(t, err)
	})
}
