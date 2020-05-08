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

	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAssets_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"Compliance": false,
			"Cost":       0,
			"Count":      0,
			"Name":       "Oracle/Database",
			"Used":       8,
		},
		{
			"Compliance": true,
			"Cost":       0,
			"Count":      2,
			"Name":       "Oracle/Exadata",
			"Used":       2,
		},
	}

	getAssetsUsageRes := map[string]float32{
		"Oracle/Database": 8,
		"Oracle/Exadata":  2,
	}

	db.EXPECT().
		GetAssetsUsage("Count", true, "Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(getAssetsUsageRes, nil)

	res, err := as.ListAssets(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestListAssets_SuccessEmpty(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{}

	getAssetsUsageRes := map[string]float32{
		"Oracle/Database": 0,
	}

	db.EXPECT().
		GetAssetsUsage("Count", true, "Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(getAssetsUsageRes, nil)

	res, err := as.ListAssets(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestListAssets_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetAssetsUsage("Count", true, "Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(nil, aerrMock)

	_, err := as.ListAssets(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.Equal(t, aerrMock, err)
}
