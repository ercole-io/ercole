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

	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvironmentStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":       3,
			"Environment": "PROD",
		},
		map[string]interface{}{
			"Count":       3,
			"Environment": "PROD",
		},
	}

	db.EXPECT().GetEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetEnvironmentStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTypeStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":                         2,
			"HardwareAbstractionTechnology": "PH",
		},
		map[string]interface{}{
			"Count":                         4,
			"HardwareAbstractionTechnology": "VMWARE",
		},
	}

	db.EXPECT().GetTypeStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetTypeStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetTypeStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTypeStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetTypeStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOperatingSystemStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":           1,
			"OperatingSystem": "RHEL5",
		},
		map[string]interface{}{
			"Count":           6,
			"OperatingSystem": "RHEL7",
		},
	}

	db.EXPECT().GetOperatingSystemStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOperatingSystemStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOperatingSystemStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOperatingSystemStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOperatingSystemStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTopUnusedOracleDatabaseInstanceResourceStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Hostname": "publicitate-36d06ca83eafa454423d2097f4965517",
			"Unused":   30,
			"_id":      utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"Hostname": "itl-csllab-112.sorint.localpippo",
			"Unused":   2,
			"_id":      "5e8c234b24f648a08585bd2e",
		},
	}

	db.EXPECT().GetTopUnusedOracleDatabaseInstanceResourceStats(
		"Italy", "PROD", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetTopUnusedOracleDatabaseInstanceResourceStats(
		"Italy", "PROD", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetTopUnusedOracleDatabaseInstanceResourceStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTopUnusedOracleDatabaseInstanceResourceStats(
		"Italy", "PROD", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetTopUnusedOracleDatabaseInstanceResourceStats(
		"Italy", "PROD", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}
