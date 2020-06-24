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

	database "github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchHosts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Cluster":                       "Angola-1dac9f7418db9b52c259ce4ba087cdb6",
			"CreatedAt":                     utils.P("2020-04-07T08:52:59.844+02:00"),
			"Databases":                     "8888888-d41d8cd98f00b204e9800998ecf8427e",
			"Environment":                   "PROD",
			"Hostname":                      "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemTotal":                      3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"VirtualizationNode":            "suspended-290dce22a939f3868f8f23a6e1f57dd8",
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"Version":                       "1.6.1",
			"HardwareAbstraction":           "VIRT",
			"_id":                           utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Cluster":                       "Puzzait",
			"CreatedAt":                     utils.P("2020-04-07T08:52:59.869+02:00"),
			"Databases":                     "",
			"Environment":                   "PROD",
			"Hostname":                      "test-virt",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemTotal":                      3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"Version":                       "1.6.1",
			"HardwareAbstraction":           "VIRT",
			"_id":                           utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	db.EXPECT().SearchHosts(
		"summary", []string{"foo", "bar", "foobarx"}, database.SearchHostsFilters{}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchHosts(
		"summary", "foo bar foobarx", database.SearchHostsFilters{}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchHosts_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchHosts(
		"summary", []string{"foo", "bar", "foobarx"}, database.SearchHostsFilters{}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchHosts(
		"summary", "foo bar foobarx", database.SearchHostsFilters{}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetHost_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]interface{}{
		"Alerts": []interface{}{
			map[string]interface{}{
				"AlertAffectedTechnology": nil,
				"AlertCategory":           model.AlertCategorySystem,
				"AlertCode":               "NEW_SERVER",
				"AlertSeverity":           "NOTICE",
				"AlertStatus":             "NEW",
				"Date":                    utils.P("2020-04-07T08:52:59.871+02:00"),
				"Description":             "The server 'test-virt' was added to ercole",
				"OtherInfo": map[string]interface{}{
					"Hostname": "test-virt",
				},
				"_id": utils.Str2oid("5e8c234b24f648a08585bd42"),
			},
		},
		"Archived":    false,
		"Cluster":     "Puzzait",
		"CreatedAt":   utils.P("2020-04-07T08:52:59.869+02:00"),
		"Databases":   "",
		"Environment": "PROD",
		"Extra": map[string]interface{}{
			"Clusters":  []interface{}{},
			"Databases": []interface{}{},
			"Filesystems": []interface{}{
				map[string]interface{}{
					"Available":  "4.6G",
					"Filesystem": "/dev/mapper/vg_os-lv_root",
					"FsType":     "xfs",
					"MountedOn":  "/",
					"Size":       "8.0G",
					"Used":       "3.5G",
					"UsedPerc":   "43%",
				},
			},
		},
		"History": []interface{}{
			map[string]interface{}{
				"CreatedAt": utils.P("2020-04-07T08:52:59.869+02:00"),
				"_id":       utils.Str2oid("5e8c234b24f648a08585bd41"),
			},
		},
		"HostDataSchemaVersion": 3,
		"Hostname":              "test-virt",
		"Info": map[string]interface{}{
			"AixCluster":                    false,
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Environment":                   "PROD",
			"Hostname":                      "test-virt",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemoryTotal":                   3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"HardwareAbstraction":           "VIRT",
		},
		"Location":           "Italy",
		"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
		"SchemaVersion":      1,
		"Schemas":            "",
		"ServerVersion":      "latest",
		"Version":            "1.6.1",
		"_id":                utils.Str2oid("5e8c234b24f648a08585bd41"),
	}

	db.EXPECT().GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetHost_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetHost(
		"foobar", utils.P("2019-12-05T14:02:03Z"), false,
	)
	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestListLocations_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []string{
		"Italy",
		"Germany",
	}

	db.EXPECT().ListLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.ListLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestListLocations_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().ListLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.ListLocations(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestListEnvironments_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []string{
		"TST",
		"SVIL",
		"PROD",
	}

	db.EXPECT().ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestListEnvironments_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.ListEnvironments(
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)
	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestArchiveHost_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().ArchiveHost("foobar").Return(nil).Times(1)

	err := as.ArchiveHost("foobar")
	require.NoError(t, err)
}

func TestArchiveHost_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().ArchiveHost("foobar").Return(nil).Times(1)

	err := as.ArchiveHost("foobar")
	require.NoError(t, err)
}
