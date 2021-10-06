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

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetOracleDatabaseArchivelogStatusStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"ArchiveLog": false,
			"Count":      8,
		},
		map[string]interface{}{
			"ArchiveLog": true,
			"Count":      20,
		},
	}

	db.EXPECT().GetOracleDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleDatabaseArchivelogStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleDatabaseEnvironmentStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":       3,
			"Environment": "TST",
		},
		map[string]interface{}{
			"Count":       5,
			"Environment": "SVIL",
		},
	}

	db.EXPECT().GetOracleDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleDatabaseEnvironmentStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleDatabaseVersionStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":   5,
			"Version": "11.2.0.3.0 Enterprise Edition",
		},
		map[string]interface{}{
			"Count":   3,
			"Version": "12.2.0.1.0 Enterprise Edition",
		},
	}

	db.EXPECT().GetOracleDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleDatabaseVersionStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTopReclaimableOracleDatabaseStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Dbname":                     "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Hostname":                   "publicitate-36d06ca83eafa454423d2097f4965517",
			"ReclaimableSegmentAdvisors": 20.5,
			"_id":                        utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		map[string]interface{}{
			"Dbname":                     "rudeboy-fb3160a04ffea22b55555bbb58137f77",
			"Hostname":                   "publicitate-36d06ca83eafa454423d2097f4965517",
			"ReclaimableSegmentAdvisors": 4.5,
			"_id":                        utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
	}

	db.EXPECT().GetTopReclaimableOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetTopReclaimableOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetTopReclaimableOracleDatabaseStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTopReclaimableOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetTopReclaimableOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleDatabasePatchStatusStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":  8,
			"Status": "KO",
		},
		map[string]interface{}{
			"Count":  8,
			"Status": "OK",
		},
	}

	db.EXPECT().GetOracleDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleDatabasePatchStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTopWorkloadOracleDatabaseStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Dbname":   "ERCOLE",
			"Hostname": "itl-csllab-112.sorint.localpippo",
			"Workload": 2,
			"_id":      utils.Str2oid("5e8c234b24f648a08585bd2e"),
		},
		map[string]interface{}{
			"Dbname":   "urcole",
			"Hostname": "itl-csllab-112.sorint.localpippo",
			"Workload": 1,
			"_id":      utils.Str2oid("5e8c234b24f648a08585bd2e"),
		},
	}

	db.EXPECT().GetTopWorkloadOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetTopWorkloadOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetTopWorkloadOracleDatabaseStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTopWorkloadOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetTopWorkloadOracleDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleDatabaseRACStatusStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count": 8,
			"RAC":   false,
		},
		map[string]interface{}{
			"Count": 16,
			"RAC":   false,
		},
	}

	db.EXPECT().GetOracleDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleDatabaseRACStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleDatabaseDataguardStatusStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Count":     8,
			"Dataguard": false,
		},
		map[string]interface{}{
			"Count":     8,
			"Dataguard": true,
		},
	}

	db.EXPECT().GetOracleDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetOracleDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetOracleDatabaseDataguardStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetOracleDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetOracleDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetOracleDatabasesStatistics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	gomock.InOrder(
		db.EXPECT().GetTotalOracleDatabaseMemorySizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(1.1), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseSegmentSizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(2.2), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseDatafileSizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(3.3), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseWorkStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(4.4), nil).Times(1),
	)

	res, err := as.GetOracleDatabasesStatistics(
		dto.GlobalFilter{
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
		},
	)

	expected := dto.OracleDatabasesStatistics{
		TotalMemorySize:   1.1,
		TotalSegmentsSize: 2.2,
		TotalDatafileSize: 3.3,
		TotalWork:         4.4,
	}

	require.NoError(t, err)
	assert.Equal(t, expected, *res)
}

func TestGetOracleDatabasesStatistics_Fail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	gomock.InOrder(
		db.EXPECT().GetTotalOracleDatabaseMemorySizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(1.1), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseSegmentSizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(2.2), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseDatafileSizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(3.3), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseWorkStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(0), aerrMock).Times(1),
	)

	res, err := as.GetOracleDatabasesStatistics(
		dto.GlobalFilter{
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
		},
	)

	assert.Equal(t, aerrMock, err)
	assert.Nil(t, res)
}

func TestGetOracleDatabasesStatistics_Fail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	gomock.InOrder(
		db.EXPECT().GetTotalOracleDatabaseMemorySizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(1.1), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseSegmentSizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(2.2), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseDatafileSizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(0), aerrMock).Times(1),
	)

	res, err := as.GetOracleDatabasesStatistics(
		dto.GlobalFilter{
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
		},
	)

	assert.Equal(t, aerrMock, err)
	assert.Nil(t, res)
}

func TestGetOracleDatabasesStatistics_Fail3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	gomock.InOrder(
		db.EXPECT().GetTotalOracleDatabaseMemorySizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(1.1), nil).Times(1),
		db.EXPECT().GetTotalOracleDatabaseSegmentSizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(0), aerrMock).Times(1),
	)

	res, err := as.GetOracleDatabasesStatistics(
		dto.GlobalFilter{
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
		},
	)

	assert.Equal(t, aerrMock, err)
	assert.Nil(t, res)
}

func TestGetOracleDatabasesStatistics_Fail4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	gomock.InOrder(
		db.EXPECT().GetTotalOracleDatabaseMemorySizeStats(
			"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		).Return(float64(0), aerrMock).Times(1),
	)

	res, err := as.GetOracleDatabasesStatistics(
		dto.GlobalFilter{
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
		},
	)

	assert.Equal(t, aerrMock, err)
	assert.Nil(t, res)
}
