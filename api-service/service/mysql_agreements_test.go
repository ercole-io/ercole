// Copyright (c) 2020 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddMySQLAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		expected := model.MySQLAgreement{
			ID:               utils.Str2oid("000000000000000000000001"),
			Type:             "server",
			NumberOfLicenses: 42,
			Clusters:         []string{"pippo"},
			Hosts:            []string{"pluto"},
		}
		db.EXPECT().AddMySQLAgreement(expected).
			Return(utils.Str2oid("000000000000000000000001"), nil).Times(1)

		agreement := model.MySQLAgreement{
			Type:             "server",
			NumberOfLicenses: 42,
			Clusters:         []string{"pippo"},
			Hosts:            []string{"pluto"},
		}
		actual, err := as.AddMySQLAgreement(agreement)
		require.NoError(t, err)

		assert.Equal(t, utils.Str2oid("000000000000000000000001"), actual)
	})

	t.Run("Error", func(t *testing.T) {
		agreement := model.MySQLAgreement{
			ID: utils.Str2oid("000000000000000000000002"),
		}
		db.EXPECT().AddMySQLAgreement(agreement).
			Return(utils.Str2oid(""), errMock).Times(1)

		actual, err := as.AddMySQLAgreement(agreement)
		require.EqualError(t, err, "MockError")

		expected := primitive.NilObjectID
		assert.Equal(t, expected, actual)
	})
}

func TestUpdateMySQLAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		agreement := model.MySQLAgreement{}
		db.EXPECT().UpdateMySQLAgreement(agreement).
			Return(nil).Times(1)

		err := as.UpdateMySQLAgreement(agreement)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		agreement := model.MySQLAgreement{}
		db.EXPECT().UpdateMySQLAgreement(agreement).
			Return(errMock).Times(1)

		err := as.UpdateMySQLAgreement(agreement)
		require.EqualError(t, err, "MockError")
	})
}

func TestGetMySQLAgreements(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		expected := []model.MySQLAgreement{
			{
				ID:               [12]byte{},
				Type:             "",
				NumberOfLicenses: 0,
				Clusters:         []string{},
				Hosts:            []string{},
			},
		}
		db.EXPECT().GetMySQLAgreements().
			Return(expected, nil).Times(1)

		actual, err := as.GetMySQLAgreements()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetMySQLAgreements().
			Return(nil, errMock).Times(1)

		actual, err := as.GetMySQLAgreements()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestDeleteMySQLAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteMySQLAgreement(id).
			Return(nil).Times(1)

		err := as.DeleteMySQLAgreement(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteMySQLAgreement(id).
			Return(errMock).Times(1)

		err := as.DeleteMySQLAgreement(id)
		require.EqualError(t, err, "MockError")
	})
}
