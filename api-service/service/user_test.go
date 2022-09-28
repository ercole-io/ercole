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

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUsers_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []model.User{
		{
			Username:  "username",
			LastLogin: &thisMoment,
			Groups:    []string{"ercole"},
		},
	}
	db.EXPECT().ListUsers().Return(expected, nil)

	res, err := as.ListUsers()
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetUser_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := &model.User{
		Username:  "username",
		LastLogin: &thisMoment,
		Groups:    []string{"ercole"},
	}

	db.EXPECT().GetUser("username").Return(expected, nil)

	res, err := as.GetUser("username")
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestUpdateUserGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		user := model.User{
			Username: "username",
			Groups:   []string{"group1", "group2"},
		}
		db.EXPECT().UpdateUserGroups(user).Return(nil)

		err := as.UpdateUserGroups(user)

		assert.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		user := model.User{
			Groups: []string{"group1", "group2"},
		}
		db.EXPECT().UpdateUserGroups(user).Return(errMock)

		err := as.UpdateUserGroups(user)

		require.EqualError(t, err, "MockError")
	})
}

func TestRemoveUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		db.EXPECT().RemoveUser("username").Return(nil)

		err := as.RemoveUser("username")
		assert.Nil(t, err)
	})
}
