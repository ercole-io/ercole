// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestListGcpRecommendations(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	ts := ThunderService{
		Config:      config.Configuration{},
		Database:    db,
		TimeNow:     utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Log:         logger.NewLogger("TEST"),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	activeProfile := []model.GcpProfile{
		{
			ID:       utils.Str2oid("000000000000000000000001"),
			Name:     "profile1",
			Selected: true,
		},
		{
			ID:       utils.Str2oid("000000000000000000000002"),
			Name:     "profile2",
			Selected: true,
		}}

	db.EXPECT().GetActiveGcpProfiles().Return(activeProfile, nil)

	recommendation := []model.GcpRecommendation{
		{
			SeqValue:    1,
			CreatedAt:   ts.TimeNow(),
			ProfileID:   utils.Str2oid("000000000000000000000001"),
			InstanceID:  uint64(1),
			Category:    "",
			Suggestion:  "",
			ProjectID:   "",
			ProjectName: "",
			ObjectType:  "",
			Details:     map[string]string{},
		},
		{
			SeqValue:    1,
			CreatedAt:   ts.TimeNow(),
			ProfileID:   utils.Str2oid("000000000000000000000002"),
			InstanceID:  uint64(2),
			Category:    "",
			Suggestion:  "",
			ProjectID:   "",
			ProjectName: "",
			ObjectType:  "",
			Details:     map[string]string{},
		}}

	profileIDs := []primitive.ObjectID{activeProfile[0].ID, activeProfile[1].ID}

	db.EXPECT().ListGcpRecommendationsByProfiles(profileIDs).Return(recommendation, nil)

	actual, err := ts.ListGcpRecommendations()

	require.NoError(t, err)

	assert.Equal(t, recommendation, actual)
}
