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
	time "time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetOciRecommendationErrors_DBError(t *testing.T) {
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
		var expectedRes []model.OciRecommendationError
		var strProfiles = []string{"TestProfile1", "TestProfile4"}

		db.EXPECT().GetOciRecommendationErrors(strProfiles).
			Return(nil, utils.NewError(utils.ErrNotFound, "DB ERROR")).Times(1)

		actual, err := as.GetOciRecommendationErrors(strProfiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrNotFound)

		assert.Equal(t, expectedRes, actual)
	})
}

func TestGetOciRecommendationErrors(t *testing.T) {
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
		expected := []model.OciRecommendationError{
			{
				SeqValue:  0,
				ProfileID: "",
				Category:  "",
				CreatedAt: time.Now().UTC(),
				Error:     "",
			},
		}
		var strProfiles = []string{"TestProfile1", "TestProfile4"}
		db.EXPECT().GetOciRecommendationErrors(strProfiles).
			Return(expected, nil).Times(1)

		actual, err := as.GetOciRecommendationErrors(strProfiles)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		var strProfiles = []string{"TestProfile1", "TestProfile4"}
		db.EXPECT().GetOciRecommendationErrors(strProfiles).
			Return(nil, errMock).Times(1)

		actual, err := as.GetOciRecommendationErrors(strProfiles)
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}
