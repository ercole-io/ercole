// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDefaultDatabaseTags_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			APIService: config.APIService{
				DefaultDatabaseTags: []string{
					"foo",
					"bar",
				},
			},
		},
	}
	res, err := as.GetDefaultDatabaseTags()
	require.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, res)
}

func TestGetErcoleFeatures_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]bool{
		"Oracle/Database": true,
		"Oracle/Exadata":  false,
	}

	getAssetsUsageRes := map[string]float32{
		"Oracle/Database": 8,
		"Oracle/Exadata":  0,
	}

	db.EXPECT().
		GetAssetsUsage("", false, "", "", utils.MAX_TIME).
		Return(getAssetsUsageRes, nil)

	res, err := as.GetErcoleFeatures()

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetErcoleFeatures_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetAssetsUsage("", false, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	_, err := as.GetErcoleFeatures()

	require.Equal(t, aerrMock, err)
}
