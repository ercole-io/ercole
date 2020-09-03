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
				"content": []interface{}{
					map[string]interface{}{
						"cpuCores":                      1,
						"cpuModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
						"cpuThreads":                    2,
						"cluster":                       "Puzzait",
						"createdAt":                     utils.P("2020-05-04T16:09:46.608+02:00").Local(),
						"environment":                   "PROD",
						"hostname":                      "test-virt",
						"kernel":                        "Linux 3.10.0-862.9.1.el7.x86_64",
						"location":                      "Italy",
						"memTotal":                      3,
						"os":                            "Red Hat Enterprise Linux 7.5",
						"oracleClusterware":             false,
						"virtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
						"cpuSockets":                    2,
						"sunCluster":                    false,
						"swapTotal":                     4,
						"hardwareAbstractionTechnology": "VMWARE",
						"veritasClusterServer":          false,
						"agentVersion":                  "1.6.1",
						"hardwareAbstraction":           "VIRT",
						"_id":                           utils.Str2oid("5eb0222a45d85f4193704944"),
						"hacmp":                         false,
					},
				},
				"metadata": map[string]interface{}{
					"empty":         false,
					"first":         true,
					"last":          false,
					"number":        0,
					"size":          1,
					"totalElements": 3,
					"totalPages":    3,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{""}, commonFilters, "createdAt", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"cpuCores":                      1,
				"cpuModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"cpuThreads":                    2,
				"cluster":                       "Puzzait",
				"createdAt":                     utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"environment":                   "PROD",
				"hostname":                      "test-virt",
				"kernel":                        "Linux 3.10.0-862.9.1.el7.x86_64",
				"location":                      "Italy",
				"memTotal":                      3,
				"os":                            "Red Hat Enterprise Linux 7.5",
				"oracleClusterware":             false,
				"virtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"cpuSockets":                    2,
				"sunCluster":                    false,
				"swapTotal":                     4,
				"hardwareAbstractionTechnology": "VMWARE",
				"veritasClusterServer":          false,
				"agentVersion":                  "1.6.1",
				"hardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5eb0222a45d85f4193704944"),
				"hacmp":                         false,
			},
			{
				"cpuCores":                      1,
				"cpuModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"cpuThreads":                    2,
				"cluster":                       interface{}(nil),
				"createdAt":                     utils.P("2020-04-24T13:50:05.46+02:00").Local(),
				"environment":                   "TST",
				"hostname":                      "test-small",
				"kernel":                        "Linux 3.10.0-514.el7.x86_64",
				"location":                      "Germany",
				"memTotal":                      3,
				"os":                            "Red Hat Enterprise Linux 7.6",
				"oracleClusterware":             false,
				"virtualizationNode":            interface{}(nil),
				"cpuSockets":                    2,
				"sunCluster":                    false,
				"swapTotal":                     1,
				"hardwareAbstractionTechnology": "VMWARE",
				"veritasClusterServer":          false,
				"agentVersion":                  "latest",
				"hardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5ea2d26d20d55cbdc35022b4"),
				"hacmp":                         false,
			},
			{
				"cpuCores":                      1,
				"cpuModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"cpuThreads":                    2,
				"cluster":                       "Puzzait",
				"createdAt":                     utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"environment":                   "TST",
				"hostname":                      "test-db",
				"kernel":                        "Linux 3.10.0-514.el7.x86_64",
				"location":                      "Germany",
				"memTotal":                      3,
				"os":                            "Red Hat Enterprise Linux 7.6",
				"oracleClusterware":             false,
				"virtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"cpuSockets":                    2,
				"sunCluster":                    false,
				"swapTotal":                     1,
				"hardwareAbstractionTechnology": "VMWARE",
				"veritasClusterServer":          false,
				"agentVersion":                  "latest",
				"hardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5e96ade270c184faca93fe36"),
				"hacmp":                         false,
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
				"cpuCores":                      1,
				"cpuModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"cpuThreads":                    2,
				"cluster":                       "Puzzait",
				"createdAt":                     utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"environment":                   "TST",
				"hostname":                      "test-db",
				"kernel":                        "Linux 3.10.0-514.el7.x86_64",
				"location":                      "Germany",
				"memTotal":                      3,
				"os":                            "Red Hat Enterprise Linux 7.6",
				"oracleClusterware":             false,
				"virtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"cpuSockets":                    2,
				"sunCluster":                    false,
				"swapTotal":                     1,
				"hardwareAbstractionTechnology": "VMWARE",
				"veritasClusterServer":          false,
				"agentVersion":                  "latest",
				"hardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5e96ade270c184faca93fe36"),
				"hacmp":                         false,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search3", func(t *testing.T) {
		out, err := m.db.SearchHosts("summary", []string{"Puzzait"}, commonFilters, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
				"cpuCores":                      1,
				"cpuModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"cpuThreads":                    2,
				"cluster":                       "Puzzait",
				"createdAt":                     utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"environment":                   "PROD",
				"hostname":                      "test-virt",
				"kernel":                        "Linux 3.10.0-862.9.1.el7.x86_64",
				"location":                      "Italy",
				"memTotal":                      3,
				"os":                            "Red Hat Enterprise Linux 7.5",
				"oracleClusterware":             false,
				"virtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
				"cpuSockets":                    2,
				"sunCluster":                    false,
				"swapTotal":                     4,
				"hardwareAbstractionTechnology": "VMWARE",
				"veritasClusterServer":          false,
				"agentVersion":                  "1.6.1",
				"hardwareAbstraction":           "VIRT",
				"_id":                           utils.Str2oid("5eb0222a45d85f4193704944"),
				"hacmp":                         false,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("lms_mode", func(t *testing.T) {
		out, err := m.db.SearchHosts("lms", []string{""}, commonFilters, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []map[string]interface{}{
			{
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
		out, err := m.db.SearchHosts("hostnames", []string{""}, commonFilters, "hostname", false, -1, -1, "", "", utils.MAX_TIME)
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
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusAck,
		Date:                    utils.P("2020-04-10T08:46:58.38+02:00"),
		Description:             "The server 'test-virt' was added to ercole",
		OtherInfo: map[string]interface{}{
			"hostname": "test-virt",
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
			"agentVersion": "1.6.1",
			"alerts": []interface{}{
				map[string]interface{}{
					"alertAffectedTechnology": nil,
					"alertCategory":           "ENGINE",
					"alertCode":               "NEW_SERVER",
					"alertSeverity":           "INFO",
					"alertStatus":             "ACK",
					"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
					"description":             "The server 'test-virt' was added to ercole",
					"otherInfo": map[string]interface{}{
						"hostname": "test-virt",
					},
					"_id": utils.Str2oid("5e96ade270c184faca93fe1b"),
				},
			},
			"archived": false,
			"cluster":  "Puzzait",
			"clusterMembershipStatus": map[string]interface{}{
				"hacmp":                false,
				"oracleClusterware":    false,
				"sunCluster":           false,
				"veritasClusterServer": false,
			},
			"clusters": []interface{}{
				map[string]interface{}{
					"cpu":           140,
					"fetchEndpoint": "???",
					"name":          "Puzzait",
					"sockets":       10,
					"type":          "vmware",
					"vms": []interface{}{
						map[string]interface{}{
							"cappedCPU":          false,
							"hostname":           "test-virt",
							"name":               "test-virt",
							"virtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
						},
						map[string]interface{}{
							"cappedCPU":          false,
							"hostname":           "test-db",
							"name":               "test-db",
							"virtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
						},
					},
				},
				map[string]interface{}{
					"cpu":           130,
					"fetchEndpoint": "???",
					"name":          "Puzzait2",
					"sockets":       13,
					"type":          "vmware",
					"vms": []interface{}{
						map[string]interface{}{
							"cappedCPU":          false,
							"hostname":           "test-virt2",
							"name":               "test-virt2",
							"virtualizationNode": "s157-cb32c10a56c256746c337e21b3ffffff",
						},
						map[string]interface{}{
							"cappedCPU":          false,
							"hostname":           "test-db2",
							"name":               "test-db2",
							"virtualizationNode": "s157-cb32c10a56c256746c337e21b3fffeua",
						},
					},
				},
			},
			"createdAt":   utils.P("2020-05-04T16:09:46.608+02:00").Local(),
			"environment": "PROD",
			"features": map[string]interface{}{
				"oracle": map[string]interface{}{
					"database": map[string]interface{}{
						"databases": nil,
					},
				},
			},
			"filesystems": []interface{}{
				map[string]interface{}{
					"availableSpace": 4.93921239e+09,
					"filesystem":     "/dev/mapper/vg_os-lv_root",
					"mountedOn":      "/",
					"size":           8.589934592e+09,
					"type":           "xfs",
					"usedSpace":      3.758096384e+09,
				},
			},
			"history": []interface{}{
				map[string]interface{}{
					"createdAt":          utils.P("2020-05-04T16:09:46.608+02:00").Local(),
					"totalDailyCPUUsage": nil,
					"_id":                utils.Str2oid("5eb0222a45d85f4193704944"),
				},
			},
			"info": map[string]interface{}{
				"cpuCores":                      1,
				"cpuFrequency":                  "2.50GHz",
				"cpuModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"cpuSockets":                    2,
				"cpuThreads":                    2,
				"coresPerSocket":                1,
				"hardwareAbstraction":           "VIRT",
				"hardwareAbstractionTechnology": "VMWARE",
				"hostname":                      "test-virt",
				"kernel":                        "Linux",
				"kernelVersion":                 "3.10.0-862.9.1.el7.x86_64",
				"memoryTotal":                   3,
				"os":                            "Red Hat Enterprise Linux",
				"osVersion":                     "7.5",
				"swapTotal":                     4,
				"threadsPerCore":                2,
			},
			"location":            "Italy",
			"virtualizationNode":  "s157-cb32c10a56c256746c337e21b3f82402",
			"schemaVersion":       1,
			"serverSchemaVersion": 1,
			"serverVersion":       "latest",
			"hostname":            "test-virt",
			"tags":                []interface{}{},
			"_id":                 utils.Str2oid("5eb0222a45d85f4193704944"),
		}
		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_detect_history", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.MAX_TIME, false)
		require.NoError(t, err)

		expectedResult := map[string]interface{}{
			"agentVersion": "latest",
			"alerts":       []interface{}{},
			"archived":     false,
			"cluster":      nil,
			"clusterMembershipStatus": map[string]interface{}{
				"hacmp":                false,
				"oracleClusterware":    false,
				"sunCluster":           false,
				"veritasClusterServer": false,
			},
			"clusters":    nil,
			"createdAt":   utils.P("2020-05-21T11:32:54.83+02:00").Local(),
			"environment": "TST",
			"features": map[string]interface{}{
				"oracle": map[string]interface{}{
					"database": map[string]interface{}{
						"unlistedRunningDatabases": []string{},
						"databases": []interface{}{
							map[string]interface{}{
								"addms":      []interface{}{},
								"asm":        false,
								"allocated":  129,
								"archivelog": false,
								"backups":    []interface{}{},
								"blockSize":  8192,
								"cpuCount":   2,
								"changes": []interface{}{
									map[string]interface{}{
										"dailyCPUUsage": 3.4,
										"segmentsSize":  50,
										"updated":       utils.P("2020-05-21T11:32:54.83+02:00").Local(),
										"datafileSize":  8,
									},
									map[string]interface{}{
										"dailyCPUUsage": 5.3,
										"segmentsSize":  100,
										"updated":       utils.P("2020-05-21T11:32:09.288+02:00").Local(),
										"datafileSize":  10,
									},
									map[string]interface{}{
										"dailyCPUUsage": 0.7,
										"segmentsSize":  3,
										"updated":       utils.P("2020-05-21T11:30:55.061+02:00").Local(),
										"datafileSize":  6,
									},
								},
								"charset":           "AL32UTF8",
								"dbTime":            184.81,
								"dailyCPUUsage":     3.4,
								"datafileSize":      8,
								"dataguard":         false,
								"elapsed":           12059.18,
								"featureUsageStats": []interface{}{},
								"instanceNumber":    1,
								"instanceName":      "pippodb1",
								"isCDB":             false,
								"licenses":          []interface{}{},
								"memoryTarget":      1.484,
								"nCharset":          "AL16UTF16",
								"name":              "pippodb",
								"pdbs":              []interface{}{},
								"pgaTarget":         0,
								"psus":              []interface{}{},
								"patches":           []interface{}{},
								"platform":          "Linux x86 64-bit",
								"sgaMaxSize":        1.484,
								"sgaTarget":         0,
								"schemas":           []interface{}{},
								"segmentAdvisors":   []interface{}{},
								"segmentsSize":      50,
								"services":          []interface{}{},
								"status":            "OPEN",
								"tablespaces":       []interface{}{},
								"uniqueName":        "pippodb",
								"version":           "12.2.0.1.0 Enterprise Edition",
								"work":              1,
							},
						},
					},
					"exadata": nil,
				},
			},
			"filesystems": []interface{}{
				map[string]interface{}{
					"availableSpace": 5.798205849e+09,
					"filesystem":     "/dev/mapper/cl_itl--csllab--112-root",
					"mountedOn":      "/",
					"size":           1.3958643712e+10,
					"type":           "ext4",
					"usedSpace":      7.19407022e+09,
				},
				map[string]interface{}{
					"availableSpace": 3.3554432e+08,
					"filesystem":     "/dev/sda1",
					"mountedOn":      "/boot",
					"size":           5.11705088e+08,
					"type":           "ext4",
					"usedSpace":      1.39460608e+08,
				},
			},
			"history": []interface{}{
				map[string]interface{}{
					"createdAt":          utils.P("2020-05-21T11:32:54.83+02:00").Local(),
					"totalDailyCPUUsage": 3.4,
					"_id":                utils.Str2oid("5ec64ac640c089c5aff44e9d"),
				},
				map[string]interface{}{
					"createdAt":          utils.P("2020-05-21T11:32:09.288+02:00").Local(),
					"totalDailyCPUUsage": 5.3,
					"_id":                utils.Str2oid("5ec64a9940c089c5aff44e9c"),
				},
				map[string]interface{}{
					"createdAt":          utils.P("2020-05-21T11:30:55.061+02:00").Local(),
					"totalDailyCPUUsage": 0.7,
					"_id":                utils.Str2oid("5ec64a4f40c089c5aff44e99"),
				},
			},
			"hostname": "newdb",
			"info": map[string]interface{}{
				"cpuCores":                      1,
				"cpuFrequency":                  "2.53GHz",
				"cpuModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"cpuSockets":                    2,
				"cpuThreads":                    2,
				"coresPerSocket":                1,
				"hardwareAbstraction":           "VIRT",
				"hardwareAbstractionTechnology": "VMWARE",
				"hostname":                      "newdb",
				"kernel":                        "Linux",
				"kernelVersion":                 "3.10.0-514.el7.x86_64",
				"memoryTotal":                   3,
				"os":                            "Red Hat Enterprise Linux",
				"osVersion":                     "7.6",
				"swapTotal":                     1,
				"threadsPerCore":                2,
			},
			"location":            "Germany",
			"schemaVersion":       1,
			"serverSchemaVersion": 1,
			"serverVersion":       "latest",
			"tags":                []interface{}{},
			"virtualizationNode":  nil,
			"_id":                 utils.Str2oid("5ec64ac640c089c5aff44e9d"),
		}

		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_detect_partial_history", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.P("2020-05-21T11:31:00.061+02:00"), false)
		require.NoError(t, err)

		expectedResult := map[string]interface{}{
			"agentVersion": "latest",
			"alerts":       []interface{}{},
			"archived":     true,
			"cluster":      nil,
			"clusterMembershipStatus": map[string]interface{}{
				"hacmp":                false,
				"oracleClusterware":    false,
				"sunCluster":           false,
				"veritasClusterServer": false,
			},
			"clusters":    nil,
			"createdAt":   utils.P("2020-05-21T11:30:55.061+02:00").Local(),
			"environment": "TST",
			"features": map[string]interface{}{
				"oracle": map[string]interface{}{
					"database": map[string]interface{}{
						"unlistedRunningDatabases": []string{},
						"databases": []interface{}{
							map[string]interface{}{
								"addms":      []interface{}{},
								"asm":        false,
								"allocated":  129,
								"archivelog": false,
								"backups":    []interface{}{},
								"blockSize":  8192,
								"cpuCount":   2,
								"changes": []interface{}{
									map[string]interface{}{
										"dailyCPUUsage": 0.7,
										"segmentsSize":  3,
										"updated":       utils.P("2020-05-21T11:30:55.061+02:00").Local(),
										"datafileSize":  6,
									},
								},
								"charset":           "AL32UTF8",
								"dbTime":            184.81,
								"dailyCPUUsage":     0.7,
								"datafileSize":      6,
								"dataguard":         false,
								"elapsed":           12059.18,
								"featureUsageStats": []interface{}{},
								"instanceNumber":    1,
								"instanceName":      "pippodb1",
								"isCDB":             false,
								"licenses":          []interface{}{},
								"memoryTarget":      1.484,
								"nCharset":          "AL16UTF16",
								"name":              "pippodb",
								"pdbs":              []interface{}{},
								"pgaTarget":         0,
								"psus":              []interface{}{},
								"patches":           []interface{}{},
								"platform":          "Linux x86 64-bit",
								"sgaMaxSize":        1.484,
								"sgaTarget":         0,
								"schemas":           []interface{}{},
								"segmentAdvisors":   []interface{}{},
								"segmentsSize":      3,
								"services":          []interface{}{},
								"status":            "OPEN",
								"tablespaces":       []interface{}{},
								"uniqueName":        "pippodb",
								"version":           "12.2.0.1.0 Enterprise Edition",
								"work":              1,
							},
						},
					},
					"exadata": nil,
				},
			},
			"filesystems": []interface{}{
				map[string]interface{}{
					"availableSpace": 5.798205849e+09,
					"filesystem":     "/dev/mapper/cl_itl--csllab--112-root",
					"mountedOn":      "/",
					"size":           1.3958643712e+10,
					"type":           "ext4",
					"usedSpace":      7.19407022e+09,
				},
				map[string]interface{}{
					"availableSpace": 3.3554432e+08,
					"filesystem":     "/dev/sda1",
					"mountedOn":      "/boot",
					"size":           5.11705088e+08,
					"type":           "ext4",
					"usedSpace":      1.39460608e+08,
				},
			},
			"history": []interface{}{
				map[string]interface{}{
					"createdAt":          utils.P("2020-05-21T11:30:55.061+02:00").Local(),
					"totalDailyCPUUsage": 0.7,
					"_id":                utils.Str2oid("5ec64a4f40c089c5aff44e99"),
				},
			},
			"hostname": "newdb",
			"info": map[string]interface{}{
				"cpuCores":                      1,
				"cpuFrequency":                  "2.53GHz",
				"cpuModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"cpuSockets":                    2,
				"cpuThreads":                    2,
				"coresPerSocket":                1,
				"hardwareAbstraction":           "VIRT",
				"hardwareAbstractionTechnology": "VMWARE",
				"hostname":                      "newdb",
				"kernel":                        "Linux",
				"kernelVersion":                 "3.10.0-514.el7.x86_64",
				"memoryTotal":                   3,
				"os":                            "Red Hat Enterprise Linux",
				"osVersion":                     "7.6",
				"swapTotal":                     1,
				"threadsPerCore":                2,
			},
			"location":            "Germany",
			"schemaVersion":       1,
			"serverSchemaVersion": 1,
			"serverVersion":       "latest",
			"tags":                []interface{}{},
			"virtualizationNode":  nil,
			"_id":                 utils.Str2oid("5ec64a4f40c089c5aff44e99"),
		}

		assert.JSONEq(t, utils.ToJSON(expectedResult), utils.ToJSON(out))
	})

	m.T().Run("should_return_raw_result", func(t *testing.T) {
		out, err := m.db.GetHost("newdb", utils.MAX_TIME, true)
		require.NoError(t, err)

		expectedResult := map[string]interface{}{
			"agentVersion": "latest",
			"archived":     false,
			"clusterMembershipStatus": map[string]interface{}{
				"hacmp":                false,
				"oracleClusterware":    false,
				"sunCluster":           false,
				"veritasClusterServer": false,
			},
			"clusters":    nil,
			"createdAt":   utils.P("2020-05-21T11:32:54.83+02:00").Local(),
			"environment": "TST",
			"features": map[string]interface{}{
				"oracle": map[string]interface{}{
					"database": map[string]interface{}{
						"unlistedRunningDatabases": []string{},
						"databases": []interface{}{
							map[string]interface{}{
								"addms":             []interface{}{},
								"asm":               false,
								"allocated":         129,
								"archivelog":        false,
								"backups":           []interface{}{},
								"blockSize":         8192,
								"cpuCount":          2,
								"charset":           "AL32UTF8",
								"dbTime":            184.81,
								"dailyCPUUsage":     3.4,
								"datafileSize":      8,
								"dataguard":         false,
								"elapsed":           12059.18,
								"featureUsageStats": []interface{}{},
								"instanceNumber":    1,
								"instanceName":      "pippodb1",
								"isCDB":             false,
								"licenses":          []interface{}{},
								"memoryTarget":      1.484,
								"nCharset":          "AL16UTF16",
								"name":              "pippodb",
								"pdbs":              []interface{}{},
								"pgaTarget":         0,
								"psus":              []interface{}{},
								"patches":           []interface{}{},
								"platform":          "Linux x86 64-bit",
								"sgaMaxSize":        1.484,
								"sgaTarget":         0,
								"schemas":           []interface{}{},
								"segmentAdvisors":   []interface{}{},
								"segmentsSize":      50,
								"services":          []interface{}{},
								"status":            "OPEN",
								"tablespaces":       []interface{}{},
								"uniqueName":        "pippodb",
								"version":           "12.2.0.1.0 Enterprise Edition",
								"work":              1,
							},
						},
					},
					"exadata": nil,
				},
			},
			"filesystems": []interface{}{
				map[string]interface{}{
					"availableSpace": 5.798205849e+09,
					"filesystem":     "/dev/mapper/cl_itl--csllab--112-root",
					"mountedOn":      "/",
					"size":           1.3958643712e+10,
					"type":           "ext4",
					"usedSpace":      7.19407022e+09,
				},
				map[string]interface{}{
					"availableSpace": 3.3554432e+08,
					"filesystem":     "/dev/sda1",
					"mountedOn":      "/boot",
					"size":           5.11705088e+08,
					"type":           "ext4",
					"usedSpace":      1.39460608e+08,
				},
			},
			"hostname": "newdb",
			"info": map[string]interface{}{
				"cpuCores":                      1,
				"cpuFrequency":                  "2.53GHz",
				"cpuModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
				"cpuSockets":                    2,
				"cpuThreads":                    2,
				"coresPerSocket":                1,
				"hardwareAbstraction":           "VIRT",
				"hardwareAbstractionTechnology": "VMWARE",
				"hostname":                      "newdb",
				"kernel":                        "Linux",
				"kernelVersion":                 "3.10.0-514.el7.x86_64",
				"memoryTotal":                   3,
				"os":                            "Red Hat Enterprise Linux",
				"osVersion":                     "7.6",
				"swapTotal":                     1,
				"threadsPerCore":                2,
			},
			"location":            "Germany",
			"schemaVersion":       1,
			"serverSchemaVersion": 1,
			"serverVersion":       "latest",
			"tags":                []interface{}{},
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

func (m *MongodbSuite) TestExistNotInClusterHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

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
