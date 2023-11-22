// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"github.com/ercole-io/ercole/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestAddAzureProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ThunderService{
		Config:      config.Configuration{},
		Database:    db,
		TimeNow:     utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Log:         logger.NewLogger("TEST"),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	t.Run("Success", func(t *testing.T) {
		expected := model.AzureProfile{
			ID:             utils.Str2oid("000000000000000000000001"),
			TenantId:       "TestProfileAdd",
			ClientId:       "TestProfileAdd",
			SubscriptionId: "TestProfileAdd",
			Region:         "eu-frankfurt-testAdd",
			ClientSecret:   &strClientSecretTestAdd,
			Selected:       false,
		}

		db.EXPECT().AddAzureProfile(expected).
			Return(nil).Times(1)

		profile := model.AzureProfile{
			TenantId:       "TestProfileAdd",
			ClientId:       "TestProfileAdd",
			SubscriptionId: "TestProfileAdd",
			Region:         "eu-frankfurt-testAdd",
			ClientSecret:   &strClientSecretTestAdd,
		}
		actual, err := as.AddAzureProfile(profile)
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		profile := model.AzureProfile{
			ID: utils.Str2oid("000000000000000000000002"),
		}
		db.EXPECT().AddAzureProfile(profile).
			Return(errMock).Times(1)

		actual, err := as.AddAzureProfile(profile)
		assert.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestUpdateAzureProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ThunderService{
		Config:      config.Configuration{},
		Database:    db,
		TimeNow:     utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Log:         logger.NewLogger("TEST"),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		profile := model.AzureProfile{}
		db.EXPECT().UpdateAzureProfile(profile).
			Return(nil).Times(1)

		actual, err := as.UpdateAzureProfile(profile)
		require.NoError(t, err)
		assert.Equal(t, profile, *actual)
	})

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	t.Run("Success1", func(t *testing.T) {
		profile := model.AzureProfile{
			ID:             utils.Str2oid("000000000000000000000001"),
			TenantId:       "TestProfileAdd",
			ClientId:       "TestProfileAdd",
			SubscriptionId: "TestProfileAdd",
			Region:         "eu-frankfurt-testAdd",
			ClientSecret:   &strClientSecretTestAdd,
			Selected:       false,
		}
		db.EXPECT().UpdateAzureProfile(profile).
			Return(nil).Times(1)

		actual, err := as.UpdateAzureProfile(profile)
		require.NoError(t, err)
		assert.Equal(t, profile, *actual)
	})

	t.Run("Error", func(t *testing.T) {
		profile := model.AzureProfile{}
		db.EXPECT().UpdateAzureProfile(profile).
			Return(errMock).Times(1)

		actual, err := as.UpdateAzureProfile(profile)
		require.EqualError(t, err, "MockError")
		assert.Nil(t, actual)
	})
}

func TestGetAzureProfiles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ThunderService{
		Config:      config.Configuration{},
		Database:    db,
		TimeNow:     utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Log:         logger.NewLogger("TEST"),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		expected := []model.AzureProfile{
			{
				ID:             [12]byte{},
				TenantId:       "",
				ClientId:       "",
				ClientSecret:   nil,
				SubscriptionId: "",
				Region:         "",
				Selected:       false,
			},
		}
		db.EXPECT().GetAzureProfiles(true).
			Return(expected, nil).Times(1)

		actual, err := as.GetAzureProfiles()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetAzureProfiles(true).
			Return(nil, errMock).Times(1)

		actual, err := as.GetAzureProfiles()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestGetMapAzureProfiles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ThunderService{
		Config:      config.Configuration{},
		Database:    db,
		TimeNow:     utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Log:         logger.NewLogger("TEST"),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		var expectedMap = make(map[primitive.ObjectID]model.AzureProfile)

		expected := model.AzureProfile{
			ID:             [12]byte{},
			TenantId:       "",
			ClientId:       "",
			ClientSecret:   nil,
			SubscriptionId: "",
			Region:         "",
			Selected:       false,
		}
		expectedMap[utils.Str2oid("000000000000000000000000")] = expected
		db.EXPECT().GetMapAzureProfiles().
			Return(expectedMap, nil).Times(1)

		actual, err := as.GetMapAzureProfiles()
		require.NoError(t, err)

		assert.Equal(t, expectedMap, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetMapAzureProfiles().
			Return(nil, errMock).Times(1)

		actual, err := as.GetMapAzureProfiles()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})

}

func TestDeleteAzureProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ThunderService{
		Config:      config.Configuration{},
		Database:    db,
		TimeNow:     utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Log:         logger.NewLogger("TEST"),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteAzureProfile(id).
			Return(nil).Times(1)

		err := as.DeleteAzureProfile(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteAzureProfile(id).
			Return(errMock).Times(1)

		err := as.DeleteAzureProfile(id)
		require.EqualError(t, err, "MockError")
	})
}

func TestSelectAzureProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ThunderService{
		Config:      config.Configuration{},
		Database:    db,
		TimeNow:     utils.Btc(utils.P("2022-06-14T16:02:10Z")),
		Log:         logger.NewLogger("TEST"),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		db.EXPECT().SelectAzureProfile("000000000000000000000001", true).
			Return(nil).Times(1)

		err := as.SelectAzureProfile("000000000000000000000001", true)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().SelectAzureProfile("000000000000000000000001", true).
			Return(errMock).Times(1)

		err := as.SelectAzureProfile("000000000000000000000001", true)
		require.EqualError(t, err, "MockError")
	})
}
