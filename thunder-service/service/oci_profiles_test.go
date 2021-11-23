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

func TestAddOciProfile(t *testing.T) {
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

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	t.Run("Success", func(t *testing.T) {
		expected := model.OciProfile{
			ID:             utils.Str2oid("000000000000000000000001"),
			Profile:        "TestProfileAdd",
			TenancyOCID:    "ocid1.tenancy.testAdd",
			UserOCID:       "ocid1.user.testAdd",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
			Region:         "eu-frankfurt-testAdd",
			PrivateKey:     &strPrivateKeyTestAdd,
		}

		db.EXPECT().AddOciProfile(expected).
			Return(nil).Times(1)

		profile := model.OciProfile{
			Profile:        "TestProfileAdd",
			TenancyOCID:    "ocid1.tenancy.testAdd",
			UserOCID:       "ocid1.user.testAdd",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
			Region:         "eu-frankfurt-testAdd",
			PrivateKey:     &strPrivateKeyTestAdd,
		}
		actual, err := as.AddOciProfile(profile)
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		profile := model.OciProfile{
			ID: utils.Str2oid("000000000000000000000002"),
		}
		db.EXPECT().AddOciProfile(profile).
			Return(errMock).Times(1)

		actual, err := as.AddOciProfile(profile)
		assert.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestUpdateOciProfile(t *testing.T) {
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
		profile := model.OciProfile{}
		db.EXPECT().UpdateOciProfile(profile).
			Return(nil).Times(1)

		actual, err := as.UpdateOciProfile(profile)
		require.NoError(t, err)
		assert.Equal(t, profile, *actual)
	})

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	t.Run("Success1", func(t *testing.T) {
		profile := model.OciProfile{
			ID:             utils.Str2oid("000000000000000000000001"),
			Profile:        "TestProfileAdd",
			TenancyOCID:    "ocid1.tenancy.testAdd",
			UserOCID:       "ocid1.user.testAdd",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
			Region:         "eu-frankfurt-testAdd",
			PrivateKey:     &strPrivateKeyTestAdd,
		}
		db.EXPECT().UpdateOciProfile(profile).
			Return(nil).Times(1)

		actual, err := as.UpdateOciProfile(profile)
		require.NoError(t, err)
		assert.Equal(t, profile, *actual)
	})

	t.Run("Error", func(t *testing.T) {
		profile := model.OciProfile{}
		db.EXPECT().UpdateOciProfile(profile).
			Return(errMock).Times(1)

		actual, err := as.UpdateOciProfile(profile)
		require.EqualError(t, err, "MockError")
		assert.Nil(t, actual)
	})
}

func TestGetOciProfiles(t *testing.T) {
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
		expected := []model.OciProfile{
			{
				ID:             [12]byte{},
				Profile:        "",
				TenancyOCID:    "",
				UserOCID:       "",
				KeyFingerprint: "",
				Region:         "",
				PrivateKey:     nil,
			},
		}
		db.EXPECT().GetOciProfiles(true).
			Return(expected, nil).Times(1)

		actual, err := as.GetOciProfiles()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetOciProfiles(true).
			Return(nil, errMock).Times(1)

		actual, err := as.GetOciProfiles()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestGetMapOciProfiles(t *testing.T) {
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
		var expectedMap = make(map[primitive.ObjectID]model.OciProfile)

		expected := model.OciProfile{
			ID:             [12]byte{},
			Profile:        "",
			TenancyOCID:    "",
			UserOCID:       "",
			KeyFingerprint: "",
			Region:         "",
			PrivateKey:     nil,
		}
		expectedMap[utils.Str2oid("000000000000000000000000")] = expected
		db.EXPECT().GetMapOciProfiles().
			Return(expectedMap, nil).Times(1)

		actual, err := as.GetMapOciProfiles()
		require.NoError(t, err)

		assert.Equal(t, expectedMap, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetMapOciProfiles().
			Return(nil, errMock).Times(1)

		actual, err := as.GetMapOciProfiles()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})

}

func TestDeleteOciProfile(t *testing.T) {
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
		db.EXPECT().DeleteOciProfile(id).
			Return(nil).Times(1)

		err := as.DeleteOciProfile(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteOciProfile(id).
			Return(errMock).Times(1)

		err := as.DeleteOciProfile(id)
		require.EqualError(t, err, "MockError")
	})
}
