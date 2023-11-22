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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestGetLastAwsObjects(t *testing.T) {
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
		expected := []model.AwsObject{
			{
				SeqValue:    uint64(1),
				ProfileID:   primitive.NilObjectID,
				CreatedAt:   time.Date(2022, 6, 1, 0, 0, 1, 0, time.UTC),
				ProfileName: "test1",
				ObjectsCount: []model.AwsObjectCount{
					{Name: "obj1", Count: 1},
				},
			},
		}

		db.EXPECT().GetLastAwsSeqValue().Return(uint64(1), nil)

		db.EXPECT().GetAwsObjectsBySeqValue(uint64(1)).Return(expected, nil)

		actual, err := as.GetLastAwsObjects()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}
