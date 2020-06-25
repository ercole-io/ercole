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

package database

import (
	"context"
	"testing"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchHosts() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	commonFilters := SearchHostsFilters{
		Cluster:           new(string),
		LTEMemoryTotal:    -1,
		GTEMemoryTotal:    -1,
		LTESwapTotal:      -1,
		GTESwapTotal:      -1,
		IsMemberOfCluster: nil,
		LTECPUCores:       -1,
		GTECPUCores:       -1,
		LTECPUThreads:     -1,
		GTECPUThreads:     -1,
	}

	//TODO: add search hosts filter tests!

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, commonFilters, "", false, -1, -1, "", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, commonFilters, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, commonFilters, "", false, -1, -1, "", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, commonFilters, "_id", true, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Content": []interface{}{
					map[string]interface{}{
						"CPUCores":                      1,
						"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
						"CPUThreads":                    2,
						"Cluster":                       "Puzzait",
						"CreatedAt":                     utils.P("2020-05-04T16:09:46.608+02:00").Local(),
						"Environment":                   "PROD",
						"Hostname":                      "test-virt",
						"Kernel":                        "Linux 3.10.0-862.9.1.el7.x86_64",
						"Location":                      "Italy",
						"MemTotal":                      3,
						"OS":                            "Red Hat Enterprise Linux 7.5",
						"OracleClusterware":             false,
						"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
						"CPUSockets":                    2,
						"SunCluster":                    false,
						"SwapTotal":                     4,
						"HardwareAbstractionTechnology": "VMWARE",
						"VeritasClusterServer":          false,
						"AgentVersion":                  "1.6.1",
						"HardwareAbstraction":           "VIRT",
						"_id":                           utils.Str2oid("5eb0222a45d85f4193704944"),
						"HACMP":                         false,
					},
				},
				"Metadata": map[string]interface{}{
					"Empty":         false,
					"First":         true,
					"Last":          false,
					"Number":        0,
					"Size":          1,
					"TotalElements": 3,
					"TotalPages":    3,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, commonFilters, "CreatedAt", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"CPUCores":                      1,
				"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"CPUThreads":                    2,
				"Cluster":                       "Puzzait",
				"CreatedAt":                     utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"Environment":                   "PROD",
				"Hostname":                      "test-virt",
				"Kernel":                        "Linux 3.10.0-862.9.1.el7.x86_64",
				"Location":                      "Italy",
				"MemTotal":                      3,
				"OS":                            "Red Hat Enterprise Linux 7.5",
				"OracleClusterware":             false,
				"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"CPUSockets":                    2,
				"SunCluster":                    false,
				"SwapTotal":                     4,
				"HardwareAbstractionTechnology": "VMWARE",
				"VeritasClusterServer":          false,
				"AgentVersion":                  "1.6.1",
				"HardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5eb0222a45d85f4193704944"),
				"HACMP":                         false,
			},
			{
				"CPUCores":                      1,
				"CPUModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUThreads":                    2,
				"Cluster":                       interface{}(nil),
				"CreatedAt":                     utils.P("2020-04-24T13:50:05.46+02:00").Local(),
				"Environment":                   "TST",
				"Hostname":                      "test-small",
				"Kernel":                        "Linux 3.10.0-514.el7.x86_64",
				"Location":                      "Germany",
				"MemTotal":                      3,
				"OS":                            "Red Hat Enterprise Linux 7.6",
				"OracleClusterware":             false,
				"VirtualizationNode":            interface{}(nil),
				"CPUSockets":                    2,
				"SunCluster":                    false,
				"SwapTotal":                     1,
				"HardwareAbstractionTechnology": "VMWARE",
				"VeritasClusterServer":          false,
				"AgentVersion":                  "latest",
				"HardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5ea2d26d20d55cbdc35022b4"),
				"HACMP":                         false,
			},
			{
				"CPUCores":                      1,
				"CPUModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUThreads":                    2,
				"Cluster":                       "Puzzait",
				"CreatedAt":                     utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"Environment":                   "TST",
				"Hostname":                      "test-db",
				"Kernel":                        "Linux 3.10.0-514.el7.x86_64",
				"Location":                      "Germany",
				"MemTotal":                      3,
				"OS":                            "Red Hat Enterprise Linux 7.6",
				"OracleClusterware":             false,
				"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"CPUSockets":                    2,
				"SunCluster":                    false,
				"SwapTotal":                     1,
				"HardwareAbstractionTechnology": "VMWARE",
				"VeritasClusterServer":          false,
				"AgentVersion":                  "latest",
				"HardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5e96ade270c184faca93fe36"),
				"HACMP":                         false,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search1", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{"foobar"}, commonFilters, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search2", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{"test-db", "ERCOLE"}, commonFilters, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"CPUCores":                      1,
				"CPUModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUThreads":                    2,
				"Cluster":                       "Puzzait",
				"CreatedAt":                     utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"Environment":                   "TST",
				"Hostname":                      "test-db",
				"Kernel":                        "Linux 3.10.0-514.el7.x86_64",
				"Location":                      "Germany",
				"MemTotal":                      3,
				"OS":                            "Red Hat Enterprise Linux 7.6",
				"OracleClusterware":             false,
				"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"CPUSockets":                    2,
				"SunCluster":                    false,
				"SwapTotal":                     1,
				"HardwareAbstractionTechnology": "VMWARE",
				"VeritasClusterServer":          false,
				"AgentVersion":                  "latest",
				"HardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5e96ade270c184faca93fe36"),
				"HACMP":                         false,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search3", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{"Puzzait"}, commonFilters, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"CPUCores":                      1,
				"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"CPUThreads":                    2,
				"Cluster":                       "Puzzait",
				"CreatedAt":                     utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"Environment":                   "PROD",
				"Hostname":                      "test-virt",
				"Kernel":                        "Linux 3.10.0-862.9.1.el7.x86_64",
				"Location":                      "Italy",
				"MemTotal":                      3,
				"OS":                            "Red Hat Enterprise Linux 7.5",
				"OracleClusterware":             false,
				"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"CPUSockets":                    2,
				"SunCluster":                    false,
				"SwapTotal":                     4,
				"HardwareAbstractionTechnology": "VMWARE",
				"VeritasClusterServer":          false,
				"AgentVersion":                  "1.6.1",
				"HardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5eb0222a45d85f4193704944"),
				"HACMP":                         false,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("lms_mode", func(t *testing.T) {
		out, err := m.db.SearchHosts("lms", []string{""}, commonFilters, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"ConnectString":            "",
				"CoresPerProcessor":        1,
				"DBInstanceName":           "ERCOLE",
				"Environment":              "TST",
				"Features":                 "Diagnostics Pack",
				"Notes":                    "",
				"OperatingSystem":          "Red Hat Enterprise Linux 7.6",
				"PhysicalCores":            2,
				"PhysicalServerName":       "Puzzait",
				"PluggableDatabaseName":    "",
				"ProcessorModel":           "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"ProcessorSpeed":           "2.53GHz",
				"Processors":               2,
				"ProductEdition":           "Enterprise",
				"ProductVersion":           "12",
				"RacNodeNames":             "",
				"ServerPurchaseDate":       "",
				"ThreadsPerCore":           2,
				"VirtualServerName":        "test-db",
				"VirtualizationTechnology": "VMWARE",
				"_id":                      utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_14.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_15.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_16.json"))
	m.InsertAlert(model.Alert{
		ID:                      utils.Str2oid("5e96ade270c184faca93fe1b"),
		AlertCategory:           model.AlertCategoryEngine,
		AlertAffectedTechnology: nil,
		AlertCode:               model.AlertCodeNewServer,
		AlertSeverity:           model.AlertSeverityNotice,
		AlertStatus:             model.AlertStatusAck,
		Date:                    utils.P("2020-04-10T08:46:58.38+02:00"),
		Description:             "The server 'test-virt' was added to ercole",
		OtherInfo: map[string]interface{}{
			"Hostname": "test-virt",
		},
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		_, err := m.db.GetHost("test-virt", utils.MIN_TIME, false)
		m.Assert().Equal(utils.AerrHostNotFound, err)
	})
	m.T().Run("should_not_find", func(t *testing.T) {
		_, err := m.db.GetHost("foobar", utils.MAX_TIME, false)
		m.Assert().Equal(utils.AerrHostNotFound, err)
	})

	m.T().Run("should_detect_cluster_physical_host_and_alerts", func(t *testing.T) {
		out, err := m.db.GetHost("test-virt", utils.MAX_TIME, false)
		m.Require().NoError(err)

		expectedResult := map[string]interface{}{
			"AgentVersion": "1.6.1",
			"Alerts": []interface{}{
				map[string]interface{}{
					"AlertAffectedTechnology": nil,
					"AlertCategory":           "ENGINE",
					"AlertCode":               "NEW_SERVER",
					"AlertSeverity":           "NOTICE",
					"AlertStatus":             "ACK",
					"Date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
					"Description":             "The server 'test-virt' was added to ercole",
					"OtherInfo": map[string]interface{}{
						"Hostname": "test-virt",
					},
					"_id": utils.Str2oid("5e96ade270c184faca93fe1b"),
				},
			},
			"Archived": false,
			"Cluster":  "Puzzait",
			"ClusterMembershipStatus": map[string]interface{}{
				"HACMP":                false,
				"OracleClusterware":    false,
				"SunCluster":           false,
				"VeritasClusterServer": false,
			},
			"Clusters": []interface{}{
				map[string]interface{}{
					"CPU":           140,
					"FetchEndpoint": "???",
					"Name":          "Puzzait",
					"Sockets":       10,
					"Type":          "vmware",
					"VMs": []interface{}{
						map[string]interface{}{
							"CappedCPU":          false,
							"Hostname":           "test-virt",
							"Name":               "test-virt",
							"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
						},
						map[string]interface{}{
							"CappedCPU":          false,
							"Hostname":           "test-db",
							"Name":               "test-db",
							"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
						},
					},
				},
				map[string]interface{}{
					"CPU":           130,
					"FetchEndpoint": "???",
					"Name":          "Puzzait2",
					"Sockets":       13,
					"Type":          "vmware",
					"VMs": []interface{}{
						map[string]interface{}{
							"CappedCPU":          false,
							"Hostname":           "test-virt2",
							"Name":               "test-virt2",
							"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3ffffff",
						},
						map[string]interface{}{
							"CappedCPU":          false,
							"Hostname":           "test-db2",
							"Name":               "test-db2",
							"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3fffeua",
						},
					},
				},
			},
			"CreatedAt":   utils.P("2020-05-04T16:09:46.608+02:00").Local(),
			"Environment": "PROD",
			"Features": map[string]interface{}{
				"Oracle": map[string]interface{}{
					"Database": map[string]interface{}{
						"Databases": nil,
					},
				},
			},
			"Filesystems": []interface{}{
				map[string]interface{}{
					"AvailableSpace": 4.93921239e+09,
					"Filesystem":     "/dev/mapper/vg_os-lv_root",
					"MountedOn":      "/",
					"Size":           8.589934592e+09,
					"Type":           "xfs",
					"UsedSpace":      3.758096384e+09,
				},
			},
			"History": []interface{}{
				map[string]interface{}{
					"CreatedAt":          utils.P("2020-05-04T16:09:46.608+02:00").Local(),
					"TotalDailyCPUUsage": nil,
					"_id":                utils.Str2oid("5eb0222a45d85f4193704944"),
				},
			},
			"Info": map[string]interface{}{
				"CPUCores":                      1,
				"CPUFrequency":                  "2.50GHz",
				"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"CPUSockets":                    2,
				"CPUThreads":                    2,
				"CoresPerSocket":                1,
				"HardwareAbstraction":           "VIRT",
				"HardwareAbstractionTechnology": "VMWARE",
				"Hostname":                      "test-virt",
				"Kernel":                        "Linux",
				"KernelVersion":                 "3.10.0-862.9.1.el7.x86_64",
				"MemoryTotal":                   3,
				"OS":                            "Red Hat Enterprise Linux",
				"OSVersion":                     "7.5",
				"SwapTotal":                     4,
				"ThreadsPerCore":                2,
			},
			"Location":            "Italy",
			"VirtualizationNode":  "s157-cb32c10a56c256746c337e21b3f82402",
			"SchemaVersion":       1,
			"ServerSchemaVersion": 1,
			"ServerVersion":       "latest",
			"Hostname":            "test-virt",
			"Tags":                []interface{}{},
			"_id":                 utils.Str2oid("5eb0222a45d85f4193704944"),
		}
		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_detect_history", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.MAX_TIME, false)
		require.NoError(t, err)

		expectedResult := map[string]interface{}{
			"AgentVersion": "latest",
			"Alerts":       []interface{}{},
			"Archived":     false,
			"Cluster":      nil,
			"ClusterMembershipStatus": map[string]interface{}{
				"HACMP":                false,
				"OracleClusterware":    false,
				"SunCluster":           false,
				"VeritasClusterServer": false,
			},
			"Clusters":    nil,
			"CreatedAt":   utils.P("2020-05-21T11:32:54.83+02:00").Local(),
			"Environment": "TST",
			"Features": map[string]interface{}{
				"Oracle": map[string]interface{}{
					"Database": map[string]interface{}{
						"Databases": []interface{}{
							map[string]interface{}{
								"ADDMs":      []interface{}{},
								"ASM":        false,
								"Allocated":  129,
								"Archivelog": false,
								"Backups":    []interface{}{},
								"BlockSize":  8192,
								"CPUCount":   2,
								"Changes": []interface{}{
									map[string]interface{}{
										"DailyCPUUsage": 3.4,
										"SegmentsSize":  50,
										"Updated":       utils.P("2020-05-21T11:32:54.83+02:00").Local(),
										"DatafileSize":  8,
									},
									map[string]interface{}{
										"DailyCPUUsage": 5.3,
										"SegmentsSize":  100,
										"Updated":       utils.P("2020-05-21T11:32:09.288+02:00").Local(),
										"DatafileSize":  10,
									},
									map[string]interface{}{
										"DailyCPUUsage": 0.7,
										"SegmentsSize":  3,
										"Updated":       utils.P("2020-05-21T11:30:55.061+02:00").Local(),
										"DatafileSize":  6,
									},
								},
								"Charset":           "AL32UTF8",
								"DBTime":            184.81,
								"DailyCPUUsage":     3.4,
								"DatafileSize":      8,
								"Dataguard":         false,
								"Elapsed":           12059.18,
								"FeatureUsageStats": []interface{}{},
								"InstanceNumber":    1,
								"IsCDB":             false,
								"Licenses":          []interface{}{},
								"MemoryTarget":      1.484,
								"NCharset":          "AL16UTF16",
								"Name":              "pippodb",
								"PDBs":              []interface{}{},
								"PGATarget":         0,
								"PSUs":              []interface{}{},
								"Patches":           []interface{}{},
								"Platform":          "Linux x86 64-bit",
								"SGAMaxSize":        1.484,
								"SGATarget":         0,
								"Schemas":           []interface{}{},
								"SegmentAdvisors":   []interface{}{},
								"SegmentsSize":      50,
								"Services":          []interface{}{},
								"Status":            "OPEN",
								"Tablespaces":       []interface{}{},
								"UniqueName":        "pippodb",
								"Version":           "12.2.0.1.0 Enterprise Edition",
								"Work":              1,
							},
						},
					},
					"Exadata": nil,
				},
			},
			"Filesystems": []interface{}{
				map[string]interface{}{
					"AvailableSpace": 5.798205849e+09,
					"Filesystem":     "/dev/mapper/cl_itl--csllab--112-root",
					"MountedOn":      "/",
					"Size":           1.3958643712e+10,
					"Type":           "ext4",
					"UsedSpace":      7.19407022e+09,
				},
				map[string]interface{}{
					"AvailableSpace": 3.3554432e+08,
					"Filesystem":     "/dev/sda1",
					"MountedOn":      "/boot",
					"Size":           5.11705088e+08,
					"Type":           "ext4",
					"UsedSpace":      1.39460608e+08,
				},
			},
			"History": []interface{}{
				map[string]interface{}{
					"CreatedAt":          utils.P("2020-05-21T11:32:54.83+02:00").Local(),
					"TotalDailyCPUUsage": 3.4,
					"_id":                utils.Str2oid("5ec64ac640c089c5aff44e9d"),
				},
				map[string]interface{}{
					"CreatedAt":          utils.P("2020-05-21T11:32:09.288+02:00").Local(),
					"TotalDailyCPUUsage": 5.3,
					"_id":                utils.Str2oid("5ec64a9940c089c5aff44e9c"),
				},
				map[string]interface{}{
					"CreatedAt":          utils.P("2020-05-21T11:30:55.061+02:00").Local(),
					"TotalDailyCPUUsage": 0.7,
					"_id":                utils.Str2oid("5ec64a4f40c089c5aff44e99"),
				},
			},
			"Hostname": "newdb",
			"Info": map[string]interface{}{
				"CPUCores":                      1,
				"CPUFrequency":                  "2.53GHz",
				"CPUModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUSockets":                    2,
				"CPUThreads":                    2,
				"CoresPerSocket":                1,
				"HardwareAbstraction":           "VIRT",
				"HardwareAbstractionTechnology": "VMWARE",
				"Hostname":                      "newdb",
				"Kernel":                        "Linux",
				"KernelVersion":                 "3.10.0-514.el7.x86_64",
				"MemoryTotal":                   3,
				"OS":                            "Red Hat Enterprise Linux",
				"OSVersion":                     "7.6",
				"SwapTotal":                     1,
				"ThreadsPerCore":                2,
			},
			"Location":            "Germany",
			"SchemaVersion":       1,
			"ServerSchemaVersion": 1,
			"ServerVersion":       "latest",
			"Tags":                []interface{}{},
			"VirtualizationNode":  nil,
			"_id":                 utils.Str2oid("5ec64ac640c089c5aff44e9d"),
		}

		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_return_raw_result", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.MAX_TIME, true)
		require.NoError(t, err)

		expectedResult := map[string]interface{}{
			"AgentVersion": "latest",
			"Archived":     false,
			"ClusterMembershipStatus": map[string]interface{}{
				"HACMP":                false,
				"OracleClusterware":    false,
				"SunCluster":           false,
				"VeritasClusterServer": false,
			},
			"Clusters":    nil,
			"CreatedAt":   utils.P("2020-05-21T11:32:54.83+02:00").Local(),
			"Environment": "TST",
			"Features": map[string]interface{}{
				"Oracle": map[string]interface{}{
					"Database": map[string]interface{}{
						"Databases": []interface{}{
							map[string]interface{}{
								"ADDMs":             []interface{}{},
								"ASM":               false,
								"Allocated":         129,
								"Archivelog":        false,
								"Backups":           []interface{}{},
								"BlockSize":         8192,
								"CPUCount":          2,
								"Charset":           "AL32UTF8",
								"DBTime":            184.81,
								"DailyCPUUsage":     3.4,
								"DatafileSize":      8,
								"Dataguard":         false,
								"Elapsed":           12059.18,
								"FeatureUsageStats": []interface{}{},
								"InstanceNumber":    1,
								"IsCDB":             false,
								"Licenses":          []interface{}{},
								"MemoryTarget":      1.484,
								"NCharset":          "AL16UTF16",
								"Name":              "pippodb",
								"PDBs":              []interface{}{},
								"PGATarget":         0,
								"PSUs":              []interface{}{},
								"Patches":           []interface{}{},
								"Platform":          "Linux x86 64-bit",
								"SGAMaxSize":        1.484,
								"SGATarget":         0,
								"Schemas":           []interface{}{},
								"SegmentAdvisors":   []interface{}{},
								"SegmentsSize":      50,
								"Services":          []interface{}{},
								"Status":            "OPEN",
								"Tablespaces":       []interface{}{},
								"UniqueName":        "pippodb",
								"Version":           "12.2.0.1.0 Enterprise Edition",
								"Work":              1,
							},
						},
					},
					"Exadata": nil,
				},
			},
			"Filesystems": []interface{}{
				map[string]interface{}{
					"AvailableSpace": 5.798205849e+09,
					"Filesystem":     "/dev/mapper/cl_itl--csllab--112-root",
					"MountedOn":      "/",
					"Size":           1.3958643712e+10,
					"Type":           "ext4",
					"UsedSpace":      7.19407022e+09,
				},
				map[string]interface{}{
					"AvailableSpace": 3.3554432e+08,
					"Filesystem":     "/dev/sda1",
					"MountedOn":      "/boot",
					"Size":           5.11705088e+08,
					"Type":           "ext4",
					"UsedSpace":      1.39460608e+08,
				},
			},
			"Hostname": "newdb",
			"Info": map[string]interface{}{
				"CPUCores":                      1,
				"CPUFrequency":                  "2.53GHz",
				"CPUModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"CPUSockets":                    2,
				"CPUThreads":                    2,
				"CoresPerSocket":                1,
				"HardwareAbstraction":           "VIRT",
				"HardwareAbstractionTechnology": "VMWARE",
				"Hostname":                      "newdb",
				"Kernel":                        "Linux",
				"KernelVersion":                 "3.10.0-514.el7.x86_64",
				"MemoryTotal":                   3,
				"OS":                            "Red Hat Enterprise Linux",
				"OSVersion":                     "7.6",
				"SwapTotal":                     1,
				"ThreadsPerCore":                2,
			},
			"Location":            "Germany",
			"SchemaVersion":       1,
			"ServerSchemaVersion": 1,
			"ServerVersion":       "latest",
			"Tags":                []interface{}{},
			"_id":                 utils.Str2oid("5ec64ac640c089c5aff44e9d"),
		}

		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestListLocations() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

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

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

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
	testSmall := utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json")
	m.InsertHostData(testSmall)
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_05.json"))

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
		assert.Equal(t, utils.AerrHostNotFound, err)
	})

	m.T().Run("should_not_find_archived_host", func(t *testing.T) {
		_, err := m.db.FindHostData("test-small3")
		assert.Equal(t, utils.AerrHostNotFound, err)
	})
}

func (m *MongodbSuite) TestReplaceHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	newHostdata := utils.LoadFixtureMongoHostDataMapAsHostData(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json")
	newHostdata.OtherInfo["Foo"] = "Bar"
	newHostdata.CreatedAt = utils.P("2020-04-28T13:50:05.46Z").Local()
	err := m.db.ReplaceHostData(newHostdata)
	m.Require().NoError(err)

	hs, err := m.db.FindHostData("test-small")
	m.Require().NoError(err)
	m.Require().NotNil(hs)

	m.Assert().Equal("Bar", hs.OtherInfo["Foo"])
}

func (m *MongodbSuite) TestExistHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_05.json"))

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

func (m *MongodbSuite) TestArchiveHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))
	val, err := m.db.ExistHostdata("test-small")
	m.Require().NoError(err)
	m.Assert().True(val)

	err = m.db.ArchiveHost("test-small")
	m.Require().NoError(err)
	val, err = m.db.ExistHostdata("test-small")
	m.Require().NoError(err)
	m.Assert().False(val)
}
