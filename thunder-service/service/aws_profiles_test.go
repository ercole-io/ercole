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

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestAddAwsProfile(t *testing.T) {
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

	var strSecretAccessKeyTestAdd = "PrivateKeyTestAdd"
	t.Run("Success", func(t *testing.T) {
		expected := model.AwsProfile{
			ID:              utils.Str2oid("000000000000000000000001"),
			AccessKeyId:     "TestProfileAdd",
			Region:          "eu-frankfurt-testAdd",
			SecretAccessKey: &strSecretAccessKeyTestAdd,
			Selected:        false,
		}

		db.EXPECT().AddAwsObject(expected, "aws_profiles").
			Return(nil)

		profile := model.AwsProfile{
			AccessKeyId:     "TestProfileAdd",
			Region:          "eu-frankfurt-testAdd",
			SecretAccessKey: &strSecretAccessKeyTestAdd,
		}
		err := as.AddAwsProfile(profile)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		profile := model.AwsProfile{
			ID: utils.Str2oid("000000000000000000000002"),
		}
		db.EXPECT().AddAwsObject(profile, "aws_profiles").
			Return(errMock)

		err := as.AddAwsProfile(profile)
		assert.EqualError(t, err, "MockError")
	})
}

func TestUpdateAwsProfile(t *testing.T) {
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
		profile := model.AwsProfile{}
		db.EXPECT().UpdateAwsProfile(profile).
			Return(nil).Times(1)

		actual, err := as.UpdateAwsProfile(profile)
		require.NoError(t, err)
		assert.Equal(t, profile, *actual)
	})

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	t.Run("Success1", func(t *testing.T) {
		profile := model.AwsProfile{
			ID:              utils.Str2oid("000000000000000000000001"),
			AccessKeyId:     "TestProfileAdd",
			Region:          "eu-frankfurt-testAdd",
			SecretAccessKey: &strPrivateKeyTestAdd,
			Selected:        false,
		}
		db.EXPECT().UpdateAwsProfile(profile).
			Return(nil).Times(1)

		actual, err := as.UpdateAwsProfile(profile)
		require.NoError(t, err)
		assert.Equal(t, profile, *actual)
	})

	t.Run("Error", func(t *testing.T) {
		profile := model.AwsProfile{}
		db.EXPECT().UpdateAwsProfile(profile).
			Return(errMock).Times(1)

		actual, err := as.UpdateAwsProfile(profile)
		require.EqualError(t, err, "MockError")
		assert.Nil(t, actual)
	})
}

func TestGetAwsProfiles(t *testing.T) {
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
		expected := []model.AwsProfile{
			{
				ID:              [12]byte{},
				AccessKeyId:     "",
				Region:          "",
				SecretAccessKey: nil,
				Selected:        false,
			},
		}
		db.EXPECT().GetAwsProfiles(true).
			Return(expected, nil).Times(1)

		actual, err := as.GetAwsProfiles()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetAwsProfiles(true).
			Return(nil, errMock).Times(1)

		actual, err := as.GetAwsProfiles()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestGetMapAwsProfiles(t *testing.T) {
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
		var expectedMap = make(map[primitive.ObjectID]model.AwsProfile)

		expected := model.AwsProfile{
			ID:              [12]byte{},
			AccessKeyId:     "",
			Region:          "",
			SecretAccessKey: nil,
			Selected:        false,
		}
		expectedMap[utils.Str2oid("000000000000000000000000")] = expected
		db.EXPECT().GetMapAwsProfiles().
			Return(expectedMap, nil).Times(1)

		actual, err := as.GetMapAwsProfiles()
		require.NoError(t, err)

		assert.Equal(t, expectedMap, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetMapAwsProfiles().
			Return(nil, errMock).Times(1)

		actual, err := as.GetMapAwsProfiles()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})

}

func TestDeleteAwsProfile(t *testing.T) {
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
		db.EXPECT().DeleteAwsProfile(id).
			Return(nil).Times(1)

		err := as.DeleteAwsProfile(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteAwsProfile(id).
			Return(errMock).Times(1)

		err := as.DeleteAwsProfile(id)
		require.EqualError(t, err, "MockError")
	})
}

func TestSelectAwsProfile(t *testing.T) {
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
		db.EXPECT().SelectAwsProfile("000000000000000000000001", true).
			Return(nil).Times(1)

		err := as.SelectAwsProfile("000000000000000000000001", true)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().SelectAwsProfile("000000000000000000000001", true).
			Return(errMock).Times(1)

		err := as.SelectAwsProfile("000000000000000000000001", true)
		require.EqualError(t, err, "MockError")
	})
}
