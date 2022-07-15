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
	time "time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAwsRecommendation_DBError(t *testing.T) {
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
		var expectedRes []model.AwsRecommendation
		var strProfiles = []string{"TestProfile1", "TestProfile4"}

		db.EXPECT().GetSelectedAwsProfiles().Return(strProfiles, nil)

		db.EXPECT().GetAwsRecommendations(strProfiles).
			Return(nil, utils.NewError(utils.ErrNotFound, "DB ERROR")).Times(1)

		actual, err := as.GetAwsRecommendations()
		require.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrNotFound)

		assert.Equal(t, expectedRes, actual)
	})
}

func TestGetAwsRecommendations(t *testing.T) {
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
		expected := []model.AwsRecommendation{
			{
				SeqValue:   0,
				ProfileID:  "",
				Category:   "",
				Suggestion: "",
				Name:       "",
				ResourceID: "",
				ObjectType: "",
				Details: []map[string]interface{}{
					{
						"Name":  "Name1",
						"Value": "Value1",
					},
				},
				CreatedAt: time.Date(2022, 6, 1, 0, 0, 1, 0, time.UTC),
			},
		}
		var strProfiles = []string{"TestProfile1", "TestProfile4"}
		db.EXPECT().GetSelectedAwsProfiles().Return(strProfiles, nil)

		db.EXPECT().GetAwsRecommendations(strProfiles).
			Return(expected, nil).Times(1)

		actual, err := as.GetAwsRecommendations()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		var strProfiles = []string{"TestProfile1", "TestProfile4"}
		db.EXPECT().GetSelectedAwsProfiles().Return(strProfiles, nil)

		db.EXPECT().GetAwsRecommendations(strProfiles).
			Return(nil, errMock).Times(1)

		actual, err := as.GetAwsRecommendations()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}
