// Copyright (c) 2022 Sorint.lab S.p.A.
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

package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestSearchHosts() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	commonFilters := dto.NewSearchHostsFilters()

	//TODO: add search hosts filter tests!

	m.T().Run("lms_mode", func(t *testing.T) {
		out, err := m.db.SearchHosts("lms", commonFilters)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"createdAt":                utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"coresPerProcessor":        1,
				"dbInstanceName":           "ERCOLE",
				"environment":              "TST",
				"options":                  "Diagnostics Pack",
				"operatingSystem":          "Red Hat Enterprise Linux",
				"physicalCores":            2,
				"physicalServerName":       "Puzzait",
				"pluggableDatabaseName":    "",
				"processorModel":           "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"processorSpeed":           "2.53GHz",
				"processors":               2,
				"productLicenseAllocated":  "EE",
				"productVersion":           "12.2.0.1.0",
				"threadsPerCore":           2,
				"virtualServerName":        "test-db",
				"virtualizationTechnology": "VMware",
				"_id":                      utils.Str2oid("5e96ade270c184faca93fe36"),
				"usingLicenseCount":        0.5,
				"usedManagementPacks":      "Diagnostics Pack",
				"licenseMetricAllocated":   "processor",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("hostnames_mode", func(t *testing.T) {
		out, err := m.db.SearchHosts("hostnames", commonFilters)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"hostname": "test-db",
			},
			{
				"hostname": "test-small",
			},
			{
				"hostname": "test-virt",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetHostDataSummaries() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	commonFilters := dto.NewSearchHostsFilters()
	// asdf := dto.SearchHostsFilters{Search: []string{""}, SortBy: "", SortDesc: false, Location: "", Environment: "", OlderThan: time.Time{wall: 0x1f52f60d, ext: 95316340478, loc: (*time.Location)(0x1a63ec0)}, PageNumber: -1, PageSize: -1, Hostname: "", Database: "", Technology: "", HardwareAbstractionTechnology: "", Cluster: (*string)(0xc001940280), VirtualizationNode: "", OperatingSystem: "", Kernel: "", LTEMemoryTotal: -1, GTEMemoryTotal: -1, LTESwapTotal: -1, GTESwapTotal: -1, IsMemberOfCluster: (*bool)(nil), CPUModel: "", LTECPUCores: -1, GTECPUCores: -1, LTECPUThreads: -1, GTECPUThreads: -1}

	//TODO: add search hosts filter tests!

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		thisFilter := commonFilters
		thisFilter.Environment = "FOOBAR"
		out, err := m.db.GetHostDataSummaries(thisFilter)
		m.Require().NoError(err)
		expectedOut := []dto.HostDataSummary{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		thisFilter := commonFilters
		thisFilter.Location = "France"

		out, err := m.db.GetHostDataSummaries(thisFilter)
		m.Require().NoError(err)
		expectedOut := []dto.HostDataSummary{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		thisFilter := commonFilters
		thisFilter.OlderThan = utils.MIN_TIME

		out, err := m.db.GetHostDataSummaries(thisFilter)
		m.Require().NoError(err)
		expectedOut := []dto.HostDataSummary{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		thisFilter := commonFilters
		thisFilter.SortBy = "createdAt"
		thisFilter.SortDesc = true

		out, err := m.db.GetHostDataSummaries(thisFilter)
		m.Require().NoError(err)

		expectedOut := []dto.HostDataSummary{
			{
				ID:           "5eb0222a45d85f4193704944",
				CreatedAt:    utils.P("2020-05-04T14:09:46.608Z").UTC(),
				Hostname:     "test-virt",
				Location:     "Italy",
				Environment:  "PROD",
				AgentVersion: "1.6.1",
				Info: model.Host{
					Hostname:                      "test-virt",
					CPUModel:                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
					CPUFrequency:                  "2.50GHz",
					CPUSockets:                    2,
					CPUCores:                      1,
					CPUThreads:                    2,
					ThreadsPerCore:                2,
					CoresPerSocket:                1,
					HardwareAbstraction:           "VIRT",
					HardwareAbstractionTechnology: "VMWARE",
					Kernel:                        "Linux",
					KernelVersion:                 "3.10.0-862.9.1.el7.x86_64",
					OS:                            "Red Hat Enterprise Linux",
					OSVersion:                     "7.5",
					MemoryTotal:                   3,
					SwapTotal:                     4,
					OtherInfo:                     map[string]interface{}{},
				},
				ClusterMembershipStatus: model.ClusterMembershipStatus{
					OracleClusterware:       false,
					SunCluster:              false,
					HACMP:                   false,
					VeritasClusterServer:    false,
					VeritasClusterHostnames: []string(nil),
					OtherInfo:               map[string]interface{}{},
				},
				VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402",
				Cluster:            "Puzzait",
				Databases:          map[string][]string{},
			},
			{
				ID:           "5ea2d26d20d55cbdc35022b4",
				CreatedAt:    utils.P("2020-04-24T13:50:05.46+02:00").UTC(),
				Hostname:     "test-small",
				Location:     "Germany",
				Environment:  "TST",
				AgentVersion: "latest",
				Info: model.Host{
					Hostname:                      "test-small",
					CPUModel:                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
					CPUFrequency:                  "2.53GHz",
					CPUSockets:                    2,
					CPUCores:                      1,
					CPUThreads:                    2,
					ThreadsPerCore:                2,
					CoresPerSocket:                1,
					HardwareAbstraction:           "VIRT",
					HardwareAbstractionTechnology: "VMWARE",
					Kernel:                        "Linux",
					KernelVersion:                 "3.10.0-514.el7.x86_64",
					OS:                            "Red Hat Enterprise Linux",
					OSVersion:                     "7.6",
					MemoryTotal:                   3,
					SwapTotal:                     1,
					OtherInfo: map[string]interface {
					}{}},
				ClusterMembershipStatus: model.ClusterMembershipStatus{
					OracleClusterware:       false,
					SunCluster:              false,
					HACMP:                   false,
					VeritasClusterServer:    false,
					VeritasClusterHostnames: []string(nil),
					OtherInfo: map[string]interface {
					}{}},
				Databases: map[string][]string{},
			},
			{
				ID:           "5e96ade270c184faca93fe36",
				CreatedAt:    utils.P("2020-04-15T06:46:58.471Z").UTC(),
				Hostname:     "test-db",
				Location:     "Germany",
				Environment:  "TST",
				AgentVersion: "latest",
				Info: model.Host{
					Hostname:                      "test-db",
					CPUModel:                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
					CPUFrequency:                  "2.53GHz",
					CPUSockets:                    2,
					CPUCores:                      1,
					CPUThreads:                    2,
					ThreadsPerCore:                2,
					CoresPerSocket:                1,
					HardwareAbstraction:           "VIRT",
					HardwareAbstractionTechnology: "VMWARE",
					Kernel:                        "Linux",
					KernelVersion:                 "3.10.0-514.el7.x86_64",
					OS:                            "Red Hat Enterprise Linux",
					OSVersion:                     "7.6",
					MemoryTotal:                   3,
					SwapTotal:                     1,
					OtherInfo: map[string]interface {
					}{},
				},
				ClusterMembershipStatus: model.ClusterMembershipStatus{
					OracleClusterware:       false,
					SunCluster:              false,
					HACMP:                   false,
					VeritasClusterServer:    false,
					VeritasClusterHostnames: []string(nil),
					OtherInfo: map[string]interface {
					}{},
				},
				VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402",
				Cluster:            "Puzzait",
				Databases: map[string][]string{
					"Oracle/Database": {"ERCOLE"},
				},
			},
		}

		assert.Equal(t, expectedOut, out)
	})

	m.T().Run("should_search1", func(t *testing.T) {

		thisFilter := commonFilters
		thisFilter.Search = []string{"foobar"}
		out, err := m.db.GetHostDataSummaries(thisFilter)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search2", func(t *testing.T) {
		thisFilter := commonFilters
		thisFilter.Search = []string{"test-db", "ERCOLE"}

		out, err := m.db.GetHostDataSummaries(thisFilter)
		m.Require().NoError(err)

		expectedOut := []dto.HostDataSummary{
			{
				ID:           "5e96ade270c184faca93fe36",
				CreatedAt:    utils.P("2020-04-15T06:46:58.471Z").UTC(),
				Hostname:     "test-db",
				Location:     "Germany",
				Environment:  "TST",
				AgentVersion: "latest",
				Info: model.Host{
					Hostname:                      "test-db",
					CPUModel:                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
					CPUFrequency:                  "2.53GHz",
					CPUSockets:                    2,
					CPUCores:                      1,
					CPUThreads:                    2,
					ThreadsPerCore:                2,
					CoresPerSocket:                1,
					HardwareAbstraction:           "VIRT",
					HardwareAbstractionTechnology: "VMWARE",
					Kernel:                        "Linux",
					KernelVersion:                 "3.10.0-514.el7.x86_64",
					OS:                            "Red Hat Enterprise Linux",
					OSVersion:                     "7.6",
					MemoryTotal:                   3,
					SwapTotal:                     1,
					OtherInfo: map[string]interface {
					}{}},
				ClusterMembershipStatus: model.ClusterMembershipStatus{
					OracleClusterware:       false,
					SunCluster:              false,
					HACMP:                   false,
					VeritasClusterServer:    false,
					VeritasClusterHostnames: []string(nil),
					OtherInfo: map[string]interface {
					}{}},
				VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402",
				Cluster:            "Puzzait",
				Databases: map[string][]string{
					"Oracle/Database": {
						"ERCOLE"}}},
		}
		assert.Equal(t, expectedOut, out)
	})

	m.T().Run("should_search3", func(t *testing.T) {
		thisFilter := commonFilters
		thisFilter.Search = []string{"Puzzait"}

		out, err := m.db.GetHostDataSummaries(thisFilter)
		m.Require().NoError(err)

		expectedOut := []dto.HostDataSummary{
			{
				ID:           "5eb0222a45d85f4193704944",
				CreatedAt:    utils.P("2020-05-04T14:09:46.608Z").UTC(),
				Hostname:     "test-virt",
				Location:     "Italy",
				Environment:  "PROD",
				AgentVersion: "1.6.1",
				Info: model.Host{
					Hostname:                      "test-virt",
					CPUModel:                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
					CPUFrequency:                  "2.50GHz",
					CPUSockets:                    2,
					CPUCores:                      1,
					CPUThreads:                    2,
					ThreadsPerCore:                2,
					CoresPerSocket:                1,
					HardwareAbstraction:           "VIRT",
					HardwareAbstractionTechnology: "VMWARE",
					Kernel:                        "Linux",
					KernelVersion:                 "3.10.0-862.9.1.el7.x86_64",
					OS:                            "Red Hat Enterprise Linux",
					OSVersion:                     "7.5",
					MemoryTotal:                   3,
					SwapTotal:                     4,
					OtherInfo:                     map[string]interface{}{}},
				ClusterMembershipStatus: model.ClusterMembershipStatus{
					OracleClusterware:       false,
					SunCluster:              false,
					HACMP:                   false,
					VeritasClusterServer:    false,
					VeritasClusterHostnames: []string(nil),
					OtherInfo:               map[string]interface{}{},
				},
				VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402",
				Cluster:            "Puzzait",
				Databases:          map[string][]string{}},
		}
		assert.Equal(t, expectedOut, out)
	})
}

func (m *MongodbSuite) TestGetHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_14.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_15.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_16.json"))
	m.InsertAlert(model.Alert{
		ID:                      utils.Str2oid("5e96ade270c184faca93fe1b"),
		AlertCategory:           model.AlertCategoryEngine,
		AlertAffectedTechnology: nil,
		AlertCode:               model.AlertCodeNewServer,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusAck,
		Date:                    utils.P("2020-04-10T08:46:58.38Z"),
		Description:             "The server 'test-virt' was added to ercole",
		OtherInfo: map[string]interface{}{
			"hostname": "test-virt",
		},
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		_, err := m.db.GetHost("test-virt", utils.MIN_TIME, false)
		m.Assert().Equal(utils.ErrHostNotFound, err)
	})
	m.T().Run("should_not_find", func(t *testing.T) {
		_, err := m.db.GetHost("foobar", utils.MAX_TIME, false)
		m.Assert().Equal(utils.ErrHostNotFound, err)
	})

	m.T().Run("should_detect_cluster_physical_host_and_alerts", func(t *testing.T) {
		out, err := m.db.GetHost("test-virt", utils.MAX_TIME, false)
		m.Require().NoError(err)

		expectedResult := dto.HostData{
			AgentVersion: "1.6.1",
			Alerts: []model.Alert{
				{
					AlertAffectedTechnology: nil,
					AlertCategory:           "ENGINE",
					AlertCode:               "NEW_SERVER",
					AlertSeverity:           "INFO",
					AlertStatus:             "ACK",
					Date:                    utils.P("2020-04-10T08:46:58.38Z"),
					Description:             "The server 'test-virt' was added to ercole",
					OtherInfo: map[string]interface{}{
						"hostname": "test-virt",
					},
					ID: utils.Str2oid("5e96ade270c184faca93fe1b"),
				},
			},
			Archived: false,
			Cluster:  "Puzzait",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				HACMP:                false,
				OracleClusterware:    false,
				SunCluster:           false,
				VeritasClusterServer: false,
			},
			Clusters: []model.ClusterInfo{
				{
					CPU:           140,
					FetchEndpoint: "???",
					Name:          "Puzzait",
					Sockets:       10,
					Type:          "vmware",
					VMs: []model.VMInfo{
						{
							CappedCPU:          false,
							Hostname:           "test-virt",
							Name:               "test-virt",
							VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402",
						},
						{
							CappedCPU:          false,
							Hostname:           "test-db",
							Name:               "test-db",
							VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402",
						},
					},
				},
				{
					CPU:           130,
					FetchEndpoint: "???",
					Name:          "Puzzait2",
					Sockets:       13,
					Type:          "vmware",
					VMs: []model.VMInfo{
						{
							CappedCPU:          false,
							Hostname:           "test-virt2",
							Name:               "test-virt2",
							VirtualizationNode: "s157-cb32c10a56c256746c337e21b3ffffff",
						},
						{
							CappedCPU:          false,
							Hostname:           "test-db2",
							Name:               "test-db2",
							VirtualizationNode: "s157-cb32c10a56c256746c337e21b3fffeua",
						},
					},
				},
			},
			CreatedAt:   utils.P("2020-05-04T14:09:46.608Z"),
			Environment: "PROD",
			Filesystems: []model.Filesystem{
				{
					AvailableSpace: 4.93921239e+09,
					Filesystem:     "/dev/mapper/vg_os-lv_root",
					MountedOn:      "/",
					Size:           8.589934592e+09,
					Type:           "xfs",
					UsedSpace:      3.758096384e+09,
				},
			},
			History: []model.History{
				{
					CreatedAt:          utils.P("2020-05-04T14:09:46.608Z"),
					TotalDailyCPUUsage: 0,
					ID:                 utils.Str2oid("5eb0222a45d85f4193704944"),
				},
			},
			Info: model.Host{
				CPUCores:                      1,
				CPUFrequency:                  "2.50GHz",
				CPUModel:                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				CPUSockets:                    2,
				CPUThreads:                    2,
				CoresPerSocket:                1,
				HardwareAbstraction:           "VIRT",
				HardwareAbstractionTechnology: "VMWARE",
				Hostname:                      "test-virt",
				Kernel:                        "Linux",
				KernelVersion:                 "3.10.0-862.9.1.el7.x86_64",
				MemoryTotal:                   3,
				OS:                            "Red Hat Enterprise Linux",
				OSVersion:                     "7.5",
				SwapTotal:                     4,
				ThreadsPerCore:                2,
			},
			Location:            "Italy",
			VirtualizationNode:  "s157-cb32c10a56c256746c337e21b3f82402",
			SchemaVersion:       1,
			ServerSchemaVersion: 1,
			ServerVersion:       "latest",
			Hostname:            "test-virt",
			Tags:                []string{},
			ID:                  utils.Str2oid("5eb0222a45d85f4193704944"),
		}
		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_detect_history", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.MAX_TIME, false)
		require.NoError(t, err)

		dailyCPUUsage := 3.4
		dbTime := 184.81
		elapsed := 12059.18
		work := 1.0

		expectedResult := dto.HostData{
			AgentVersion: "latest",
			Alerts:       []model.Alert{},
			Archived:     false,
			Cluster:      "",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				HACMP:                false,
				OracleClusterware:    false,
				SunCluster:           false,
				VeritasClusterServer: false,
			},
			Clusters:    nil,
			CreatedAt:   utils.P("2020-05-21T09:32:54.83Z"),
			Environment: "TST",
			Features: model.Features{
				Oracle: &model.OracleFeature{
					Database: &model.OracleDatabaseFeature{
						UnlistedRunningDatabases: []string{},
						Databases: []model.OracleDatabase{
							{
								ADDMs:        []model.OracleDatabaseAddm{},
								ASM:          false,
								Allocable:    129,
								Archivelog:   false,
								Backups:      []model.OracleDatabaseBackup{},
								BlockSize:    8192,
								CPUCount:     2,
								SegmentsSize: 50,
								DatafileSize: 8,
								Changes: []model.Changes{
									{
										DailyCPUUsage: 3.4,
										SegmentsSize:  50,
										Updated:       utils.P("2020-05-21T09:32:54.83Z"),
										DatafileSize:  8,
										Allocable:     129,
									},
									{
										DailyCPUUsage: 5.3,
										SegmentsSize:  100,
										Updated:       utils.P("2020-05-21T09:32:09.288Z"),
										DatafileSize:  10,
										Allocable:     129,
									},
									{
										DailyCPUUsage: 0.7,
										SegmentsSize:  3,
										Updated:       utils.P("2020-05-21T09:30:55.061Z"),
										DatafileSize:  6,
										Allocable:     129,
									},
								},
								Charset:           "AL32UTF8",
								DBTime:            &dbTime,
								DailyCPUUsage:     &dailyCPUUsage,
								Dataguard:         false,
								Elapsed:           &elapsed,
								FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{},
								InstanceNumber:    1,
								InstanceName:      "pippodb1",
								IsCDB:             false,
								Licenses:          []model.OracleDatabaseLicense{},
								MemoryTarget:      1.484,
								NCharset:          "AL16UTF16",
								Name:              "pippodb",
								PDBs:              []model.OracleDatabasePluggableDatabase{},
								PGATarget:         0,
								PSUs:              []model.OracleDatabasePSU{},
								Patches:           []model.OracleDatabasePatch{},
								Platform:          "Linux x86 64-bit",
								SGAMaxSize:        1.484,
								SGATarget:         0,
								Schemas:           []model.OracleDatabaseSchema{},
								SegmentAdvisors:   []model.OracleDatabaseSegmentAdvisor{},
								Services:          []model.OracleDatabaseService{},
								Status:            "OPEN",
								Tablespaces:       []model.OracleDatabaseTablespace{},
								UniqueName:        "pippodb",
								Version:           "12.2.0.1.0 Enterprise Edition",
								Work:              &work,
							},
						},
					},
				},
			},
			Filesystems: []model.Filesystem{
				{
					AvailableSpace: 5.798205849e+09,
					Filesystem:     "/dev/mapper/cl_itl--csllab--112-root",
					MountedOn:      "/",
					Size:           1.3958643712e+10,
					Type:           "ext4",
					UsedSpace:      7.19407022e+09,
				},
				{
					AvailableSpace: 3.3554432e+08,
					Filesystem:     "/dev/sda1",
					MountedOn:      "/boot",
					Size:           5.11705088e+08,
					Type:           "ext4",
					UsedSpace:      1.39460608e+08,
				},
			},
			History: []model.History{
				{
					CreatedAt:          utils.P("2020-05-21T09:32:54.83Z"),
					TotalDailyCPUUsage: 3.4,
					ID:                 utils.Str2oid("5ec64ac640c089c5aff44e9d"),
				},
				{
					CreatedAt:          utils.P("2020-05-21T09:32:09.288Z"),
					TotalDailyCPUUsage: 5.3,
					ID:                 utils.Str2oid("5ec64a9940c089c5aff44e9c"),
				},
				{
					CreatedAt:          utils.P("2020-05-21T09:30:55.061Z"),
					TotalDailyCPUUsage: 0.7,
					ID:                 utils.Str2oid("5ec64a4f40c089c5aff44e99"),
				},
			},
			Hostname: "newdb",
			Info: model.Host{
				CPUCores:                      1,
				CPUFrequency:                  "2.53GHz",
				CPUModel:                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				CPUSockets:                    2,
				CPUThreads:                    2,
				CoresPerSocket:                1,
				HardwareAbstraction:           "VIRT",
				HardwareAbstractionTechnology: "VMWARE",
				Hostname:                      "newdb",
				Kernel:                        "Linux",
				KernelVersion:                 "3.10.0-514.el7.x86_64",
				MemoryTotal:                   3,
				OS:                            "Red Hat Enterprise Linux",
				OSVersion:                     "7.6",
				SwapTotal:                     1,
				ThreadsPerCore:                2,
			},
			Location:            "Germany",
			SchemaVersion:       1,
			ServerSchemaVersion: 1,
			ServerVersion:       "latest",
			Tags:                []string{},
			VirtualizationNode:  "",
			ID:                  utils.Str2oid("5ec64ac640c089c5aff44e9d"),
		}

		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_detect_partial_history", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.P("2020-05-21T11:31:00.061+02:00"), false)
		require.NoError(t, err)

		dailyCPUUsage := 0.7
		dbTime := 184.81
		elapsed := 12059.18
		work := 1.0

		expectedResult := dto.HostData{
			AgentVersion: "latest",
			Alerts:       []model.Alert{},
			Archived:     true,
			Cluster:      "",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				HACMP:                false,
				OracleClusterware:    false,
				SunCluster:           false,
				VeritasClusterServer: false,
			},
			Clusters:    nil,
			CreatedAt:   utils.P("2020-05-21T09:30:55.061Z"),
			Environment: "TST",
			Features: model.Features{
				Oracle: &model.OracleFeature{
					Database: &model.OracleDatabaseFeature{
						UnlistedRunningDatabases: []string{},
						Databases: []model.OracleDatabase{
							{
								ADDMs:      []model.OracleDatabaseAddm{},
								ASM:        false,
								Allocable:  129,
								Archivelog: false,
								Backups:    []model.OracleDatabaseBackup{},
								BlockSize:  8192,
								CPUCount:   2,
								Changes: []model.Changes{
									{
										DailyCPUUsage: 0.7,
										SegmentsSize:  3,
										Updated:       utils.P("2020-05-21T09:30:55.061Z"),
										DatafileSize:  6,
										Allocable:     129,
									},
								},
								Charset:           "AL32UTF8",
								DBTime:            &dbTime,
								DailyCPUUsage:     &dailyCPUUsage,
								DatafileSize:      6,
								Dataguard:         false,
								Elapsed:           &elapsed,
								FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{},
								InstanceNumber:    1,
								InstanceName:      "pippodb1",
								IsCDB:             false,
								Licenses:          []model.OracleDatabaseLicense{},
								MemoryTarget:      1.484,
								NCharset:          "AL16UTF16",
								Name:              "pippodb",
								PDBs:              []model.OracleDatabasePluggableDatabase{},
								PGATarget:         0,
								PSUs:              []model.OracleDatabasePSU{},
								Patches:           []model.OracleDatabasePatch{},
								Platform:          "Linux x86 64-bit",
								SGAMaxSize:        1.484,
								SGATarget:         0,
								Schemas:           []model.OracleDatabaseSchema{},
								SegmentAdvisors:   []model.OracleDatabaseSegmentAdvisor{},
								SegmentsSize:      3,
								Services:          []model.OracleDatabaseService{},
								Status:            "OPEN",
								Tablespaces:       []model.OracleDatabaseTablespace{},
								UniqueName:        "pippodb",
								Version:           "12.2.0.1.0 Enterprise Edition",
								Work:              &work,
							},
						},
					},
				},
			},
			Filesystems: []model.Filesystem{
				{
					AvailableSpace: 5.798205849e+09,
					Filesystem:     "/dev/mapper/cl_itl--csllab--112-root",
					MountedOn:      "/",
					Size:           1.3958643712e+10,
					Type:           "ext4",
					UsedSpace:      7.19407022e+09,
				},
				{
					AvailableSpace: 3.3554432e+08,
					Filesystem:     "/dev/sda1",
					MountedOn:      "/boot",
					Size:           5.11705088e+08,
					Type:           "ext4",
					UsedSpace:      1.39460608e+08,
				},
			},
			History: []model.History{
				{
					CreatedAt:          utils.P("2020-05-21T09:30:55.061Z"),
					TotalDailyCPUUsage: 0.7,
					ID:                 utils.Str2oid("5ec64a4f40c089c5aff44e99"),
				},
			},
			Hostname: "newdb",
			Info: model.Host{
				CPUCores:                      1,
				CPUFrequency:                  "2.53GHz",
				CPUModel:                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				CPUSockets:                    2,
				CPUThreads:                    2,
				CoresPerSocket:                1,
				HardwareAbstraction:           "VIRT",
				HardwareAbstractionTechnology: "VMWARE",
				Hostname:                      "newdb",
				Kernel:                        "Linux",
				KernelVersion:                 "3.10.0-514.el7.x86_64",
				MemoryTotal:                   3,
				OS:                            "Red Hat Enterprise Linux",
				OSVersion:                     "7.6",
				SwapTotal:                     1,
				ThreadsPerCore:                2,
			},
			Location:            "Germany",
			SchemaVersion:       1,
			ServerSchemaVersion: 1,
			ServerVersion:       "latest",
			Tags:                []string{},
			ID:                  utils.Str2oid("5ec64a4f40c089c5aff44e99"),
		}

		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_return_raw_result", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.MAX_TIME, true)
		require.NoError(t, err)

		dailyCPUUsage := 3.4
		dbTime := 184.81
		elapsed := 12059.18
		work := 1.0

		expectedResult := dto.HostData{
			AgentVersion: "latest",
			Archived:     false,
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				HACMP:                false,
				OracleClusterware:    false,
				SunCluster:           false,
				VeritasClusterServer: false,
			},
			Clusters:    nil,
			CreatedAt:   utils.P("2020-05-21T09:32:54.83Z"),
			Environment: "TST",
			Features: model.Features{
				Oracle: &model.OracleFeature{
					Database: &model.OracleDatabaseFeature{
						UnlistedRunningDatabases: []string{},
						Databases: []model.OracleDatabase{
							{
								ADDMs:             []model.OracleDatabaseAddm{},
								ASM:               false,
								Allocable:         129,
								Archivelog:        false,
								Backups:           []model.OracleDatabaseBackup{},
								BlockSize:         8192,
								CPUCount:          2,
								Charset:           "AL32UTF8",
								DBTime:            &dbTime,
								DailyCPUUsage:     &dailyCPUUsage,
								DatafileSize:      8,
								Dataguard:         false,
								Elapsed:           &elapsed,
								FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{},
								InstanceNumber:    1,
								InstanceName:      "pippodb1",
								IsCDB:             false,
								Licenses:          []model.OracleDatabaseLicense{},
								MemoryTarget:      1.484,
								NCharset:          "AL16UTF16",
								Name:              "pippodb",
								PDBs:              []model.OracleDatabasePluggableDatabase{},
								PGATarget:         0,
								PSUs:              []model.OracleDatabasePSU{},
								Patches:           []model.OracleDatabasePatch{},
								Platform:          "Linux x86 64-bit",
								SGAMaxSize:        1.484,
								SGATarget:         0,
								Schemas:           []model.OracleDatabaseSchema{},
								SegmentAdvisors:   []model.OracleDatabaseSegmentAdvisor{},
								SegmentsSize:      50,
								Services:          []model.OracleDatabaseService{},
								Status:            "OPEN",
								Tablespaces:       []model.OracleDatabaseTablespace{},
								UniqueName:        "pippodb",
								Version:           "12.2.0.1.0 Enterprise Edition",
								Work:              &work,
							},
						},
					},
				},
			},
			Filesystems: []model.Filesystem{
				{
					AvailableSpace: 5.798205849e+09,
					Filesystem:     "/dev/mapper/cl_itl--csllab--112-root",
					MountedOn:      "/",
					Size:           1.3958643712e+10,
					Type:           "ext4",
					UsedSpace:      7.19407022e+09,
				},
				{
					AvailableSpace: 3.3554432e+08,
					Filesystem:     "/dev/sda1",
					MountedOn:      "/boot",
					Size:           5.11705088e+08,
					Type:           "ext4",
					UsedSpace:      1.39460608e+08,
				},
			},
			Hostname: "newdb",
			Info: model.Host{
				CPUCores:                      1,
				CPUFrequency:                  "2.53GHz",
				CPUModel:                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				CPUSockets:                    2,
				CPUThreads:                    2,
				CoresPerSocket:                1,
				HardwareAbstraction:           "VIRT",
				HardwareAbstractionTechnology: "VMWARE",
				Hostname:                      "newdb",
				Kernel:                        "Linux",
				KernelVersion:                 "3.10.0-514.el7.x86_64",
				MemoryTotal:                   3,
				OS:                            "Red Hat Enterprise Linux",
				OSVersion:                     "7.6",
				SwapTotal:                     1,
				ThreadsPerCore:                2,
			},
			Location:            "Germany",
			SchemaVersion:       1,
			ServerSchemaVersion: 1,
			ServerVersion:       "latest",
			Tags:                []string{},
			ID:                  utils.Str2oid("5ec64ac640c089c5aff44e9d"),
		}

		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_14.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_15.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_16.json"))

	actual, err := m.db.GetHostData("newdb", utils.MAX_TIME)
	require.NoError(m.T(), err)

	expected := mongoutils.LoadFixtureMongoHostDataMapAsHostData(m.T(), "../../fixture/test_apiservice_mongohostdata_16.json")

	assert.Equal(m.T(), expected, *actual)
}

func (m *MongodbSuite) TestGetHostDatas() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_14.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_15.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_16.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

	actual, err := m.db.GetHostDatas(utils.MAX_TIME)
	require.NoError(m.T(), err)

	expected := []model.HostDataBE{
		mongoutils.LoadFixtureMongoHostDataMapAsHostData(m.T(), "../../fixture/test_apiservice_mongohostdata_16.json"),
		mongoutils.LoadFixtureMongoHostDataMapAsHostData(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"),
	}

	assert.Equal(m.T(), expected, actual)
}

func (m *MongodbSuite) TestListLocations() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.ListLocations("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{}, out)
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.ListLocations("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{}, out)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.ListLocations("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{}, out)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.ListLocations("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{"Italy", "Germany"}, out)
	})
}

func (m *MongodbSuite) TestListEnvironments() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.ListEnvironments("France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{}, out)
	})

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.ListEnvironments("", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{}, out)
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.ListEnvironments("", "", utils.MIN_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{}, out)
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.ListEnvironments("", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.ElementsMatch(t, []string{"PROD", "DEV", "TST"}, out)
	})
}

func (m *MongodbSuite) TestFindHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	testSmall := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json")
	m.InsertHostData(testSmall)
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_05.json"))

	m.T().Run("should_find_test_small", func(t *testing.T) {
		out, err := m.db.FindHostData("test-small")
		m.Require().NoError(err)
		assert.Equal(t, utils.Str2oid("5ea2d26d20d55cbdc35022b4"), out.ID)
		assert.False(t, out.Archived)
		assert.Equal(t, "test-small", out.Hostname)
		assert.Equal(t, utils.P("2020-04-24T11:50:05.46Z"), out.CreatedAt)
	})

	m.T().Run("should_not_find_anything", func(t *testing.T) {
		_, err := m.db.FindHostData("foobar")
		assert.Equal(t, utils.ErrHostNotFound, err)
	})

	m.T().Run("should_not_find_archived_host", func(t *testing.T) {
		_, err := m.db.FindHostData("test-small3")
		assert.Equal(t, utils.ErrHostNotFound, err)
	})
}

func (m *MongodbSuite) TestReplaceHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	newHostdata := mongoutils.LoadFixtureMongoHostDataMapAsHostData(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json")
	newHostdata.OtherInfo["foo"] = "bar"
	newHostdata.CreatedAt = utils.P("2020-04-28T13:50:05.46Z").Local()
	err := m.db.ReplaceHostData(newHostdata)
	m.Require().NoError(err)

	hs, err := m.db.FindHostData("test-small")
	m.Require().NoError(err)
	m.Require().NotNil(hs)

	m.Assert().Equal("bar", hs.OtherInfo["foo"])
}

func (m *MongodbSuite) TestExistHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_05.json"))

	m.T().Run("should_find_test_small", func(t *testing.T) {
		out, err := m.db.ExistHostdata("test-small")
		require.NoError(t, err)

		assert.True(t, out)
	})

	m.T().Run("should_not_find_anything", func(t *testing.T) {
		out, err := m.db.ExistHostdata("foobar")
		require.NoError(t, err)

		assert.False(t, out)
	})

	m.T().Run("should_not_find_archived_host", func(t *testing.T) {
		out, err := m.db.ExistHostdata("test-small3")
		require.NoError(t, err)

		assert.False(t, out)
	})
}

func (m *MongodbSuite) TestDismissHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	val, err := m.db.ExistHostdata("test-small")
	m.Require().NoError(err)
	m.Assert().True(val)

	err = m.db.DismissHost("test-small")
	m.Require().NoError(err)
	val, err = m.db.ExistHostdata("test-small")
	m.Require().NoError(err)
	m.Assert().False(val)
}

func (m *MongodbSuite) TestGetHostMinValidCreatedAtDate() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))

	createdAt := utils.P("2020-04-24T13:50:05.460+02:00").UTC()
	createdDateTemp, err := m.db.GetHostMinValidCreatedAtDate("test-small")
	m.Require().NoError(err)
	createdDate := createdDateTemp["createdAt"].(primitive.DateTime).Time().UTC()
	m.Assert().EqualValues(createdAt, createdDate)
}

func (m *MongodbSuite) TestGetListValidHostsByRangeDates() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))

	hosts, err := m.db.GetListValidHostsByRangeDates(utils.P("2020-04-23T13:50:05.460+02:00").UTC(), utils.P("2020-04-25T13:50:05.460+02:00").UTC())
	m.Require().NoError(err)
	actual := []string{"test-small"}
	m.Assert().EqualValues(hosts, actual)
}

func (m *MongodbSuite) TestGetListDismissedHostsByRangeDates() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_25.json"))

	dismissedHosts, err := m.db.GetListDismissedHostsByRangeDates(utils.P("2020-04-23T13:50:05.460+02:00").UTC(), utils.P("2020-04-25T13:50:05.460+02:00").UTC())
	m.Require().NoError(err)
	actual := []string{"test-small"}
	m.Assert().EqualValues(dismissedHosts, actual)
}

func (m *MongodbSuite) TestExistNotInClusterHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

	m.T().Run("foobar_not_exist", func(t *testing.T) {
		out, err := m.db.ExistNotInClusterHost("foobar")
		require.NoError(t, err)
		assert.False(t, out)
	})
	m.T().Run("test_db_in_cluster", func(t *testing.T) {
		out, err := m.db.ExistNotInClusterHost("test-db")
		require.NoError(t, err)
		assert.False(t, out)
	})
	m.T().Run("test_db3_not_in_cluster", func(t *testing.T) {
		out, err := m.db.ExistNotInClusterHost("test-db3")
		require.NoError(t, err)
		assert.True(t, out)
	})
}
