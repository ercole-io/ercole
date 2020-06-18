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

	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatabaseArchivelogStatusStats_Success(t *testing.T) {
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

	db.EXPECT().GetDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetDatabaseArchivelogStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetDatabaseArchivelogStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetDatabaseEnvironmentStats_Success(t *testing.T) {
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

	db.EXPECT().GetDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetDatabaseEnvironmentStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetDatabaseEnvironmentStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetDatabaseVersionStats_Success(t *testing.T) {
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

	db.EXPECT().GetDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetDatabaseVersionStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetDatabaseVersionStats(
		"Italy", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTopReclaimableDatabaseStats_Success(t *testing.T) {
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

	db.EXPECT().GetTopReclaimableDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetTopReclaimableDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetTopReclaimableDatabaseStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTopReclaimableDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetTopReclaimableDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetDatabasePatchStatusStats_Success(t *testing.T) {
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

	db.EXPECT().GetDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetDatabasePatchStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetDatabasePatchStatusStats(
		"Italy", utils.P("2019-06-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTopWorkloadDatabaseStats_Success(t *testing.T) {
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

	db.EXPECT().GetTopWorkloadDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetTopWorkloadDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetTopWorkloadDatabaseStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTopWorkloadDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetTopWorkloadDatabaseStats(
		"Italy", 10, utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetDatabaseRACStatusStats_Success(t *testing.T) {
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

	db.EXPECT().GetDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetDatabaseRACStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetDatabaseRACStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetDatabaseDataguardStatusStats_Success(t *testing.T) {
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

	db.EXPECT().GetDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetDatabaseDataguardStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetDatabaseDataguardStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTotalDatabaseWorkStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseWorkStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(5), nil).Times(1)

	res, err := as.GetTotalDatabaseWorkStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, float64(5), res)
}

func TestGetTotalDatabaseWorkStatsStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseWorkStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(0), aerrMock).Times(1)

	res, err := as.GetTotalDatabaseWorkStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Zero(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTotalDatabaseMemorySizeStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(17.151), nil).Times(1)

	res, err := as.GetTotalDatabaseMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, float64(17.151), res)
}

func TestGetTotalDatabaseMemorySizeStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(0), aerrMock).Times(1)

	res, err := as.GetTotalDatabaseMemorySizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Zero(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTotalDatabaseDatafileSizeStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseDatafileSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(158), nil).Times(1)

	res, err := as.GetTotalDatabaseDatafileSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, float64(158), res)
}

func TestGetTotalDatabaseDatafileSizeStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseDatafileSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(0), aerrMock).Times(1)

	res, err := as.GetTotalDatabaseDatafileSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Zero(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetTotalDatabaseSegmentSizeStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseSegmentSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(117), nil).Times(1)

	res, err := as.GetTotalDatabaseSegmentSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, float64(117), res)
}

func TestGetTotalDatabaseSegmentSizeStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetTotalDatabaseSegmentSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(float64(0), aerrMock).Times(1)

	res, err := as.GetTotalDatabaseSegmentSizeStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Zero(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetDatabaseLicenseComplianceStatusStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]interface{}{
		"Compliant": false,
		"Count":     0,
		"Used":      14,
	}

	db.EXPECT().GetDatabaseLicenseComplianceStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetDatabaseLicenseComplianceStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetDatabaseLicenseComplianceStatusStats_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetDatabaseLicenseComplianceStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetDatabaseLicenseComplianceStatusStats(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}
