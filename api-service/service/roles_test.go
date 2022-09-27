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

func TestInsertRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		expected := model.RoleType{
			ID:   utils.Str2oid("000000000000000000000001"),
			Name: "Test",
		}
		db.EXPECT().InsertRole(expected).
			Return(nil).Times(1)

		role := model.RoleType{
			Name: "Test",
		}
		actual, err := as.InsertRole(role)
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		role := model.RoleType{
			ID:   utils.Str2oid("000000000000000000000002"),
			Name: "Test",
		}
		db.EXPECT().InsertRole(role).
			Return(errMock).Times(1)

		actual, err := as.InsertRole(role)
		assert.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestUpdateRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		role := model.RoleType{
			Name: "Test",
		}
		db.EXPECT().UpdateRole(role).
			Return(nil).Times(1)

		actual, err := as.UpdateRole(role)
		require.NoError(t, err)
		assert.Equal(t, role, *actual)
	})

	t.Run("Error", func(t *testing.T) {
		role := model.RoleType{
			Name: "Test",
		}
		db.EXPECT().UpdateRole(role).
			Return(errMock).Times(1)

		actual, err := as.UpdateRole(role)
		require.EqualError(t, err, "MockError")
		assert.Nil(t, actual)
	})
}

func TestGetRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		expected := []model.RoleType{
			{
				ID:   [12]byte{},
				Name: "Test",
			},
		}
		db.EXPECT().GetRoles().
			Return(expected, nil).Times(1)

		actual, err := as.GetRoles()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetRoles().
			Return(nil, errMock).Times(1)

		actual, err := as.GetRoles()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestGetRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		expected := model.RoleType{
			ID:   utils.Str2oid("000000000000000000000002"),
			Name: "Test",
		}
		db.EXPECT().GetRole(utils.Str2oid("000000000000000000000002")).
			Return(&expected, nil).Times(1)

		actual, err := as.GetRole(utils.Str2oid("000000000000000000000002"))
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetRole(utils.Str2oid("000000000000000000000002")).
			Return(nil, errMock).Times(1)

		actual, err := as.GetRole(utils.Str2oid("000000000000000000000002"))
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestDeleteRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteRole(id).
			Return(nil).Times(1)

		err := as.DeleteRole(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteRole(id).
			Return(errMock).Times(1)

		err := as.DeleteRole(id)
		require.EqualError(t, err, "MockError")
	})
}
