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
	time "time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestGetAwsRDS(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	ts := ThunderService{
		Database: db,
		TimeNow:  time.Now,
	}

	createdAt := ts.TimeNow()
	objID := primitive.NewObjectID()
	expect := []model.AwsRDS{
		{
			SeqValue:    1,
			ProfileID:   objID,
			ProfileName: "test-profile",
			Instances:   nil,
			CreatedAt:   createdAt,
		}}

	db.EXPECT().GetAwsRDS().Return(expect, nil)

	actual, err := ts.GetAwsRDS()

	require.NoError(t, err)

	assert.Equal(t, expect, actual)
}
