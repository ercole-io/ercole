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

func (m *MongodbSuite) TestInsertGroup() {
	defer m.db.Client.Database(m.dbname).Collection(groupCollection).DeleteMany(context.TODO(), bson.M{})

	group := model.GroupType{
		ID:    utils.Str2oid("000000000000000000000001"),
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	m.T().Run("should_insert", func(t *testing.T) {
		err := m.db.InsertGroup(group)
		m.Require().NoError(err)
	})
}

func (m *MongodbSuite) TestUpdateGroup() {
	defer m.db.Client.Database(m.dbname).Collection(groupCollection).DeleteMany(context.TODO(), bson.M{})

	group := model.GroupType{
		ID:    utils.Str2oid("000000000000000000000001"),
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	m.T().Run("error not found", func(t *testing.T) {
		err := m.db.UpdateGroup(group)
		var aerr *utils.AdvancedError
		assert.ErrorAs(t, err, &aerr)
		assert.ErrorIs(t, aerr.Err, utils.ErrGroupNotFound)
	})

	m.T().Run("should_update", func(t *testing.T) {
		_, err := m.db.Client.Database(m.dbname).Collection(groupCollection).
			InsertOne(
				context.TODO(),
				group,
			)
		require.NoError(t, err)

		err = m.db.UpdateGroup(group)
		assert.NoError(t, err)
	})
}

func (m *MongodbSuite) TestGetGroups() {
	m.T().Run("should_load_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(groupCollection).DeleteMany(context.TODO(), bson.M{})

		groups := []model.GroupType{
			{
				ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				Name:  "Test",
				Roles: []string{"role1", "role2"},
			},
			{
				ID:    utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				Name:  "Testdue",
				Roles: []string{"role1", "role2"},
			},
		}
		groupsInt := []interface{}{
			groups[0],
			groups[1],
		}
		_, err := m.db.Client.Database(m.dbname).Collection(groupCollection).InsertMany(context.TODO(), groupsInt)
		require.Nil(m.T(), err)

		actual, err := m.db.GetGroups()
		m.Require().NoError(err)

		assert.Equal(t, groups, actual)
	})

	m.T().Run("should_load_empty", func(t *testing.T) {
		actual, err := m.db.GetGroups()
		m.Require().NoError(err)

		groups := make([]model.GroupType, 0)
		assert.Equal(t, groups, actual)
	})
}

func (m *MongodbSuite) TestGetGroup() {
	m.T().Run("should_load_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(groupCollection).DeleteMany(context.TODO(), bson.M{})

		group := model.GroupType{
			ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}

		groups := []model.GroupType{
			{
				ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				Name:  "Test",
				Roles: []string{"role1", "role2"},
			},
			{
				ID:    utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				Name:  "Testdue",
				Roles: []string{"role1", "role2"},
			},
		}
		groupsInt := []interface{}{
			groups[0],
			groups[1],
		}
		_, err := m.db.Client.Database(m.dbname).Collection(groupCollection).InsertMany(context.TODO(), groupsInt)
		require.Nil(m.T(), err)

		actual, err := m.db.GetGroup(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"))
		m.Require().NoError(err)

		assert.Equal(t, &group, actual)
	})

	m.T().Run("should_load_empty", func(t *testing.T) {
		actual, err := m.db.GetGroups()
		m.Require().NoError(err)

		groups := make([]model.GroupType, 0)
		assert.Equal(t, groups, actual)
	})
}

func (m *MongodbSuite) TestDeleteGroup() {
	defer m.db.Client.Database(m.dbname).Collection(groupCollection).DeleteMany(context.TODO(), bson.M{})

	id := utils.Str2oid("000000000000000000000001")

	m.T().Run("error not found", func(t *testing.T) {
		err := m.db.DeleteGroup(id)
		var aerr *utils.AdvancedError
		assert.ErrorAs(t, err, &aerr)
		assert.ErrorIs(t, aerr.Err, utils.ErrGroupNotFound)
	})

	m.T().Run("should_delete", func(t *testing.T) {
		group := model.GroupType{
			ID:    utils.Str2oid("000000000000000000000001"),
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}
		_, err := m.db.Client.Database(m.dbname).Collection(groupCollection).
			InsertOne(
				context.TODO(),
				group,
			)
		require.NoError(t, err)

		err = m.db.DeleteGroup(group.ID)
		assert.NoError(t, err)
	})
}
