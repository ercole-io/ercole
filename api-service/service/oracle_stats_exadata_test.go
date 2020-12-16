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

func TestGetTotalOracleExadataMemorySizeStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalOracleExadataMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(1316), nil).Times(1)

	res, err := as.GetTotalOracleExadataMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, float64(1316), res)
}

func TestGetTotalOracleExadataMemorySizeStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalOracleExadataMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(0), aerrMock).Times(1)

	res, err := as.GetTotalOracleExadataMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Zero(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTotalOracleExadataCPUStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]interface{}{
		"Enabled": 156,
		"Total":   216,
	}

	db.EXPECT().GetTotalOracleExadataCPUStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetTotalOracleExadataCPUStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetTotalOracleExadataCPUStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalOracleExadataCPUStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetTotalOracleExadataCPUStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetAverageOracleExadataStorageUsageStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetAverageOracleExadataStorageUsageStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(54.666668), nil).Times(1)

	res, err := as.GetAverageOracleExadataStorageUsageStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, float64(54.666668), res)
}

func TestGetAverageOracleExadataStorageUsageStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetAverageOracleExadataStorageUsageStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(0), aerrMock).Times(1)

	res, err := as.GetAverageOracleExadataStorageUsageStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Zero(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleExadataStorageErrorCountStatusStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":   10,
			"Failing": false,
		},
		map[string]interface{}{
			"Count":   8,
			"Failing": true,
		},
	}

	db.EXPECT().GetOracleExadataStorageErrorCountStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleExadataStorageErrorCountStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleExadataStorageErrorCountStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleExadataStorageErrorCountStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleExadataStorageErrorCountStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleExadataPatchStatusStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":  5,
			"Status": true,
		},
		map[string]interface{}{
			"Count":  2,
			"Status": false,
		},
	}

	db.EXPECT().GetOracleExadataPatchStatusStats(
		"Italy", "PROD", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleExadataPatchStatusStats(
		"Italy", "PROD", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleExadataPatchStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleExadataPatchStatusStats(
		"Italy", "PROD", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleExadataPatchStatusStats(
		"Italy", "PROD", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}
