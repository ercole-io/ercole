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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetRoles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		expected := []model.Role{
			{
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
		expected := model.Role{
			Name: "Test",
		}
		db.EXPECT().GetRole("Test").
			Return(&expected, nil).Times(1)

		actual, err := as.GetRole("Test")
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetRole("Bart").
			Return(nil, errMock).Times(1)

		actual, err := as.GetRole("Bart")
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestAddRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expected := model.Role{
		Name:        "Test",
		Description: "test role",
		Location:    "IT",
		Permission:  "admin",
	}

	t.Run("Success", func(t *testing.T) {
		db.EXPECT().ListAllLocations("", "", utils.MAX_TIME).Return([]string{"IT"}, nil)
		db.EXPECT().AddRole(expected).Return(nil)

		err := as.AddRole(expected)
		require.NoError(t, err)
	})

	t.Run("Error location", func(t *testing.T) {
		db.EXPECT().ListAllLocations("", "", utils.MAX_TIME).Return([]string{}, nil)
		db.EXPECT().AddRole(expected).Return(utils.ErrInvalidLocation).AnyTimes()

		err := as.AddRole(expected)
		require.EqualError(t, err, "Invalid location")
	})
}

func TestUpdateRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expected := model.Role{
		Name:        "Test",
		Description: "test role",
		Location:    "IT",
		Permission:  "admin",
	}

	documents := bson.D{
		primitive.E{Key: "description", Value: "new description"},
		primitive.E{Key: "location", Value: expected.Location},
		primitive.E{Key: "permission", Value: expected.Permission},
	}

	t.Run("Success", func(t *testing.T) {
		db.EXPECT().ListAllLocations("", "", utils.MAX_TIME).Return([]string{"IT"}, nil)
		db.EXPECT().UpdateRole("Test", documents).Return(nil)

		expected.Description = "new description"
		err := as.UpdateRole(expected)
		require.NoError(t, err)
	})

	t.Run("Error location", func(t *testing.T) {
		db.EXPECT().ListAllLocations("", "", utils.MAX_TIME).Return([]string{}, nil)

		db.EXPECT().UpdateRole("Test", documents).Return(utils.ErrInvalidLocation).AnyTimes()

		expected.Description = "new description"
		err := as.UpdateRole(expected)
		require.EqualError(t, err, "Invalid location")
	})
}
