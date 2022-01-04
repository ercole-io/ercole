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
	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetRecommendation_InvalidProfileId(t *testing.T) {
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

	t.Run("BadRequest", func(t *testing.T) {
		var expectedMap = make(map[primitive.ObjectID]model.OciProfile)
		var expectedRes []model.OciRecommendation
		expectedRes = make([]model.OciRecommendation, 0)

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

		profiles := []string{"pippo"}
		actual, err := as.GetOciRecommendations(profiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrInvalidProfileId)

		assert.Equal(t, expectedRes, actual)
	})
}

func TestGetRecommendation_ProfileNotFound(t *testing.T) {
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

	t.Run("BadRequest", func(t *testing.T) {
		//var err1 error
		var expectedRes = make([]model.OciRecommendation, 0)
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

		profiles := []string{"0000000000000000000000ab"}
		actual, err := as.GetOciRecommendations(profiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrNotFound)

		assert.Equal(t, expectedRes, actual)
	})
}

func TestGetRecommendation_DBError(t *testing.T) {
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

	t.Run("DB Error", func(t *testing.T) {
		var expectedRes []model.OciRecommendation

		db.EXPECT().GetMapOciProfiles().
			Return(nil, utils.NewError(utils.ErrNotFound, "DB ERROR")).Times(1)

		profiles := []string{"pippo"}
		actual, err := as.GetOciRecommendations(profiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrNotFound)

		assert.Equal(t, expectedRes, actual)
	})
}
