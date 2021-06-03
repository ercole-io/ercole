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

	dto "github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
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
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchHosts(
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
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
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchHosts(
		"summary",
		dto.SearchHostsFilters{
			Search:      []string{"foo", "bar", "foobarx"},
			SortBy:      "Memory",
			SortDesc:    true,
			Location:    "Italy",
			Environment: "PROD",
			OlderThan:   utils.P("2019-12-05T14:02:03Z"),
			PageNumber:  1,
			PageSize:    1,
		},
	)
	assert.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchHostsAsLMS(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	hosts := []map[string]interface{}{
		{
			"coresPerProcessor":        1,
			"dbInstanceName":           "ERCOLE",
			"environment":              "TST",
			"licenseMetricAllocated":   "processor",
			"operatingSystem":          "Red Hat Enterprise Linux",
			"options":                  "Diagnostics Pack",
			"physicalCores":            2,
			"physicalServerName":       "erclin7dbx",
			"pluggableDatabaseName":    "",
			"processorModel":           "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
			"processorSpeed":           "2.53GHz",
			"processors":               2,
			"productLicenseAllocated":  "EE",
			"productVersion":           "12",
			"threadsPerCore":           2,
			"usedManagementPacks":      "Diagnostics Pack",
			"usingLicenseCount":        0.5,
			"virtualServerName":        "itl-csllab-112.sorint.localpippo",
			"virtualizationTechnology": "VMware",
			"_id":                      utils.Str2oid("5efc38ab79f92e4cbf283b03"),
		},
		{
			"coresPerProcessor":        4,
			"dbInstanceName":           "rudeboy-fb3160a04ffea22b55555bbb58137f77",
			"environment":              "SVIL",
			"licenseMetricAllocated":   "processor",
			"operatingSystem":          "Red Hat Enterprise Linux",
			"options":                  "",
			"physicalCores":            8,
			"physicalServerName":       "",
			"pluggableDatabaseName":    "",
			"processorModel":           "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz",
			"processorSpeed":           "2.93GHz",
			"processors":               2,
			"productLicenseAllocated":  "EE",
			"productVersion":           "11",
			"threadsPerCore":           2,
			"usedManagementPacks":      "",
			"usingLicenseCount":        4,
			"virtualServerName":        "publicitate-36d06ca83eafa454423d2097f4965517",
			"virtualizationTechnology": "",
			"_id":                      utils.Str2oid("5efc38ab79f92e4cbf283b04"),
		},
	}

	filters := dto.SearchHostsFilters{
		Search:         []string{"foobar"},
		SortBy:         "Processors",
		SortDesc:       true,
		Location:       "Italy",
		Environment:    "TST",
		OlderThan:      utils.P("2020-06-10T11:54:59Z"),
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}

	t.Run("with no agreements", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().
				SearchHosts("lms", gomock.Any()).
				DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
					assert.EqualValues(t, filters, actual)

					return hosts, nil
				}),
			db.EXPECT().
				ListOracleDatabaseAgreements().
				Return([]dto.OracleDatabaseAgreementFE{}, nil),
		)

		sp, err := as.SearchHostsAsLMS(filters)
		assert.NoError(t, err)

		assert.Equal(t, "erclin7dbx", sp.GetCellValue("Database_&_EBS_DB_Tier", "B4"))
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", sp.GetCellValue("Database_&_EBS_DB_Tier", "C4"))
		assert.Equal(t, "VMware", sp.GetCellValue("Database_&_EBS_DB_Tier", "D4"))
		assert.Equal(t, "ERCOLE", sp.GetCellValue("Database_&_EBS_DB_Tier", "E4"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "F4"))
		assert.Equal(t, "TST", sp.GetCellValue("Database_&_EBS_DB_Tier", "G4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue("Database_&_EBS_DB_Tier", "H4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue("Database_&_EBS_DB_Tier", "I4"))
		assert.Equal(t, "12", sp.GetCellValue("Database_&_EBS_DB_Tier", "N4"))
		assert.Equal(t, "EE", sp.GetCellValue("Database_&_EBS_DB_Tier", "O4"))
		assert.Equal(t, "processor", sp.GetCellValue("Database_&_EBS_DB_Tier", "P4"))
		assert.Equal(t, "0.5", sp.GetCellValue("Database_&_EBS_DB_Tier", "Q4"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AC4"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AD4"))
		assert.Equal(t, "1", sp.GetCellValue("Database_&_EBS_DB_Tier", "AE4"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AF4"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AG4"))
		assert.Equal(t, "2.53GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AH4"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue("Database_&_EBS_DB_Tier", "AJ4"))

		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "B5"))
		assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sp.GetCellValue("Database_&_EBS_DB_Tier", "C5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "D5"))
		assert.Equal(t, "rudeboy-fb3160a04ffea22b55555bbb58137f77", sp.GetCellValue("Database_&_EBS_DB_Tier", "E5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "F5"))
		assert.Equal(t, "SVIL", sp.GetCellValue("Database_&_EBS_DB_Tier", "G5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "H5"))
		assert.Equal(t, "", sp.GetCellValue("Database_&_EBS_DB_Tier", "I5"))
		assert.Equal(t, "11", sp.GetCellValue("Database_&_EBS_DB_Tier", "N5"))
		assert.Equal(t, "EE", sp.GetCellValue("Database_&_EBS_DB_Tier", "O5"))
		assert.Equal(t, "processor", sp.GetCellValue("Database_&_EBS_DB_Tier", "P5"))
		assert.Equal(t, "4", sp.GetCellValue("Database_&_EBS_DB_Tier", "Q5"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AC5"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AD5"))
		assert.Equal(t, "4", sp.GetCellValue("Database_&_EBS_DB_Tier", "AE5"))
		assert.Equal(t, "8", sp.GetCellValue("Database_&_EBS_DB_Tier", "AF5"))
		assert.Equal(t, "2", sp.GetCellValue("Database_&_EBS_DB_Tier", "AG5"))
		assert.Equal(t, "2.93GHz", sp.GetCellValue("Database_&_EBS_DB_Tier", "AH5"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue("Database_&_EBS_DB_Tier", "AJ5"))
	})

	t.Run("with agreements", func(t *testing.T) {
		agreements := []dto.OracleDatabaseAgreementFE{
			{
				ID:              utils.Str2oid("aaaaaaaaaaaa"),
				AgreementID:     "",
				CSI:             "csi001",
				LicenseTypeID:   "",
				ItemDescription: "",
				Metric:          "",
				ReferenceNumber: "",
				Unlimited:       false,
				CatchAll:        false,
				Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
					{Hostname: "publicitate-36d06ca83eafa454423d2097f4965517"},
				},
				AvailableCount: 0,
				LicensesCount:  0,
				UsersCount:     0,
			},
			{
				ID:              utils.Str2oid("aaaaaaaaaaaa"),
				AgreementID:     "",
				CSI:             "csi002",
				LicenseTypeID:   "",
				ItemDescription: "",
				Metric:          "",
				ReferenceNumber: "",
				Unlimited:       false,
				CatchAll:        false,
				Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
					{Hostname: "publicitate-36d06ca83eafa454423d2097f4965517"},
				},
				AvailableCount: 0,
				LicensesCount:  0,
				UsersCount:     0,
			},
		}

		gomock.InOrder(
			db.EXPECT().
				SearchHosts("lms", gomock.Any()).
				DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
					assert.EqualValues(t, filters, actual)

					return hosts, nil
				}),
			db.EXPECT().
				ListOracleDatabaseAgreements().
				Return(agreements, nil),
		)

		sp, err := as.SearchHostsAsLMS(filters)
		assert.NoError(t, err)

		sheet := "Database_&_EBS_DB_Tier"
		assert.Equal(t, "erclin7dbx", sp.GetCellValue(sheet, "B4"))
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", sp.GetCellValue(sheet, "C4"))
		assert.Equal(t, "VMware", sp.GetCellValue(sheet, "D4"))
		assert.Equal(t, "ERCOLE", sp.GetCellValue(sheet, "E4"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "F4"))
		assert.Equal(t, "TST", sp.GetCellValue(sheet, "G4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue(sheet, "H4"))
		assert.Equal(t, "Diagnostics Pack", sp.GetCellValue(sheet, "I4"))
		assert.Equal(t, "12", sp.GetCellValue(sheet, "N4"))
		assert.Equal(t, "EE", sp.GetCellValue(sheet, "O4"))
		assert.Equal(t, "processor", sp.GetCellValue(sheet, "P4"))
		assert.Equal(t, "0.5", sp.GetCellValue(sheet, "Q4"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "R4"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz", sp.GetCellValue(sheet, "AC4"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AD4"))
		assert.Equal(t, "1", sp.GetCellValue(sheet, "AE4"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AF4"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AG4"))
		assert.Equal(t, "2.53GHz", sp.GetCellValue(sheet, "AH4"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue(sheet, "AJ4"))

		assert.Equal(t, "", sp.GetCellValue(sheet, "B5"))
		assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sp.GetCellValue(sheet, "C5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "D5"))
		assert.Equal(t, "rudeboy-fb3160a04ffea22b55555bbb58137f77", sp.GetCellValue(sheet, "E5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "F5"))
		assert.Equal(t, "SVIL", sp.GetCellValue(sheet, "G5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "H5"))
		assert.Equal(t, "", sp.GetCellValue(sheet, "I5"))
		assert.Equal(t, "11", sp.GetCellValue(sheet, "N5"))
		assert.Equal(t, "EE", sp.GetCellValue(sheet, "O5"))
		assert.Equal(t, "processor", sp.GetCellValue(sheet, "P5"))
		assert.Equal(t, "4", sp.GetCellValue(sheet, "Q5"))
		assert.Equal(t, "csi001, csi002", sp.GetCellValue(sheet, "R5"))
		assert.Equal(t, "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz", sp.GetCellValue(sheet, "AC5"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AD5"))
		assert.Equal(t, "4", sp.GetCellValue(sheet, "AE5"))
		assert.Equal(t, "8", sp.GetCellValue(sheet, "AF5"))
		assert.Equal(t, "2", sp.GetCellValue(sheet, "AG5"))
		assert.Equal(t, "2.93GHz", sp.GetCellValue(sheet, "AH5"))
		assert.Equal(t, "Red Hat Enterprise Linux", sp.GetCellValue(sheet, "AJ5"))
	})
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
				"AlertCategory":           model.AlertCategoryEngine,
				"AlertCode":               "NEW_SERVER",
				"AlertSeverity":           "INFO",
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
