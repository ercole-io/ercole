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

package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestInsertGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		expected := model.GroupType{
			ID:    utils.Str2oid("000000000000000000000001"),
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}
		db.EXPECT().InsertGroup(expected).
			Return(nil).Times(1)

		group := model.GroupType{
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}
		actual, err := as.InsertGroup(group)
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		group := model.GroupType{
			ID:    utils.Str2oid("000000000000000000000002"),
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}
		db.EXPECT().InsertGroup(group).
			Return(errMock).Times(1)

		actual, err := as.InsertGroup(group)
		assert.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestUpdateGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		group := model.GroupType{
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}
		db.EXPECT().UpdateGroup(group).
			Return(nil).Times(1)

		actual, err := as.UpdateGroup(group)
		require.NoError(t, err)
		assert.Equal(t, group, *actual)
	})

	t.Run("Error", func(t *testing.T) {
		group := model.GroupType{
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}
		db.EXPECT().UpdateGroup(group).
			Return(errMock).Times(1)

		actual, err := as.UpdateGroup(group)
		require.EqualError(t, err, "MockError")
		assert.Nil(t, actual)
	})
}

func TestGetGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		expected := []model.GroupType{
			{
				ID:    [12]byte{},
				Name:  "Test",
				Roles: []string{"role1", "role2"},
			},
		}
		db.EXPECT().GetGroups().
			Return(expected, nil).Times(1)

		actual, err := as.GetGroups()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetGroups().
			Return(nil, errMock).Times(1)

		actual, err := as.GetGroups()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestGetGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		expected := model.GroupType{
			ID:    utils.Str2oid("000000000000000000000002"),
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		}
		db.EXPECT().GetGroup(utils.Str2oid("000000000000000000000002")).
			Return(&expected, nil).Times(1)

		actual, err := as.GetGroup(utils.Str2oid("000000000000000000000002"))
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetGroup(utils.Str2oid("000000000000000000000002")).
			Return(nil, errMock).Times(1)

		actual, err := as.GetGroup(utils.Str2oid("000000000000000000000002"))
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestDeleteGroup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteGroup(id).
			Return(nil).Times(1)

		err := as.DeleteGroup(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteGroup(id).
			Return(errMock).Times(1)

		err := as.DeleteGroup(id)
		require.EqualError(t, err, "MockError")
	})
}
