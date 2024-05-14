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
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddGcpProfile(t *testing.T) {
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

	newProfile := model.GcpProfile{
		ID:   utils.Str2oid("000000000000000000000001"),
		Name: "profile-test",
	}

	db.EXPECT().AddGcpProfile(newProfile).Return(nil).AnyTimes()

	newProfileRequest := dto.GcpProfileRequest{
		Name: "profile-test",
	}

	err := ts.AddGcpProfile(newProfileRequest)

	require.NoError(t, err)
}

func TestGetGcpProfiles(t *testing.T) {
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

	expected := []model.GcpProfile{
		{
			ID:   utils.Str2oid("000000000000000000000001"),
			Name: "profile-test-01",
		},
		{
			ID:   utils.Str2oid("000000000000000000000002"),
			Name: "profile-test-02",
		},
	}

	db.EXPECT().ListGcpProfiles().Return(expected, nil)

	actual, err := ts.GetGcpProfiles()

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestUpdateGcpProfile(t *testing.T) {
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

	id := utils.Str2oid("000000000000000000000001")

	expect := model.GcpProfile{
		ID:   id,
		Name: "profile-test-01",
	}

	db.EXPECT().UpdateGcpProfile(id, expect).Return(nil).AnyTimes()

	err := ts.UpdateGcpProfile("000000000000000000000001", dto.GcpProfileRequest{
		Name: "profile-test-01",
	})

	require.NoError(t, err)

}

func TestRemoveGcpProfile(t *testing.T) {
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

	id := utils.Str2oid("000000000000000000000001")

	db.EXPECT().RemoveGcpProfile(id).Return(nil)

	err := ts.RemoveGcpProfile("000000000000000000000001")

	require.NoError(t, err)
}
