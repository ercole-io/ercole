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

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchOracleDatabases() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "", "PROD", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{""}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"archivelog":   false,
						"blockSize":    8192,
						"cpuCount":     2,
						"charset":      "AL32UTF8",
						"createdAt":    utils.P("2020-04-15T08:46:58.471+02:00").Local(),
						"datafileSize": 6,
						"dataguard":    false,
						"environment":  "TST",
						"ha":           false,
						"hostname":     "test-db",
						"location":     "Germany",
						"memory":       1.484,
						"name":         "ERCOLE",
						"rac":          false,
						"segmentsSize": 3,
						"status":       "OPEN",
						"uniqueName":   "ERCOLE",
						"version":      "12.2.0.1.0 Enterprise Edition",
						"work":         1,
						"_id":          utils.Str2oid("5e96ade270c184faca93fe36"),
					},
				},
				"metadata": map[string]interface{}{
					"empty":         false,
					"first":         true,
					"last":          false,
					"number":        0,
					"size":          1,
					"totalElements": 2,
					"totalPages":    2,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{""}, "memory", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"archivelog":   true,
				"blockSize":    8192,
				"cpuCount":     2,
				"charset":      "AL32UTF8",
				"createdAt":    utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"datafileSize": 6,
				"dataguard":    true,
				"environment":  "TST",
				"ha":           true,
				"hostname":     "test-db2",
				"location":     "Germany",
				"memory":       90.254,
				"name":         "pokemons",
				"rac":          true,
				"segmentsSize": 3,
				"status":       "OPEN",
				"uniqueName":   "pokemons",
				"version":      "12.2.0.1.0 Enterprise Edition",
				"work":         1,
				"_id":          utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
			map[string]interface{}{
				"archivelog":   false,
				"blockSize":    8192,
				"cpuCount":     2,
				"charset":      "AL32UTF8",
				"createdAt":    utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"datafileSize": 6,
				"dataguard":    false,
				"environment":  "TST",
				"ha":           false,
				"hostname":     "test-db",
				"location":     "Germany",
				"memory":       1.484,
				"name":         "ERCOLE",
				"rac":          false,
				"segmentsSize": 3,
				"status":       "OPEN",
				"uniqueName":   "ERCOLE",
				"version":      "12.2.0.1.0 Enterprise Edition",
				"work":         1,
				"_id":          utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{"pokemon", "test-db2"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"archivelog":   true,
				"blockSize":    8192,
				"cpuCount":     2,
				"charset":      "AL32UTF8",
				"createdAt":    utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"datafileSize": 6,
				"dataguard":    true,
				"environment":  "TST",
				"ha":           true,
				"hostname":     "test-db2",
				"location":     "Germany",
				"memory":       90.254,
				"name":         "pokemons",
				"rac":          true,
				"segmentsSize": 3,
				"status":       "OPEN",
				"uniqueName":   "pokemons",
				"version":      "12.2.0.1.0 Enterprise Edition",
				"work":         1,
				"_id":          utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("notfullmode", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"archivelog":   false,
				"blockSize":    8192,
				"cpuCount":     2,
				"charset":      "AL32UTF8",
				"createdAt":    utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"datafileSize": 6,
				"dataguard":    false,
				"environment":  "TST",
				"ha":           false,
				"hostname":     "test-db",
				"location":     "Germany",
				"memory":       1.484,
				"name":         "ERCOLE",
				"rac":          false,
				"segmentsSize": 3,
				"status":       "OPEN",
				"uniqueName":   "ERCOLE",
				"version":      "12.2.0.1.0 Enterprise Edition",
				"work":         1,
				"_id":          utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			map[string]interface{}{
				"archivelog":   true,
				"blockSize":    8192,
				"cpuCount":     2,
				"charset":      "AL32UTF8",
				"createdAt":    utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"datafileSize": 6,
				"dataguard":    true,
				"environment":  "TST",
				"ha":           true,
				"hostname":     "test-db2",
				"location":     "Germany",
				"memory":       90.254,
				"name":         "pokemons",
				"rac":          true,
				"segmentsSize": 3,
				"status":       "OPEN",
				"uniqueName":   "pokemons",
				"version":      "12.2.0.1.0 Enterprise Edition",
				"work":         1,
				"_id":          utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases(true, []string{""}, "memory", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"addms": []interface{}{
					map[string]interface{}{
						"action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						"benefit":        83.34,
						"finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						"recommendation": "SQL Tuning",
					},
					map[string]interface{}{
						"action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						"benefit":        68.24,
						"finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						"recommendation": "Segment Tuning",
					},
				},
				"asm":        false,
				"allocable":  129,
				"archivelog": false,
				"backups": []interface{}{
					map[string]interface{}{
						"avgBckSize": 13,
						"backupType": "Archivelog",
						"hour":       "01:30",
						"retention":  "1 NUMBERS",
						"weekDays":   []interface{}{"Wednesday"},
					},
					map[string]interface{}{
						"avgBckSize": 45,
						"backupType": "Archivelog",
						"hour":       "03:00",
						"retention":  "1 NUMBERS",
						"weekDays": []interface{}{"Tuesday",
							"Sunday",
							"Monday",
							"Saturday",
							"Wednesday",
						},
					},
				},
				"blockSize":     8192,
				"cpuCount":      2,
				"charset":       "AL32UTF8",
				"createdAt":     utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"dbTime":        184.81,
				"dailyCPUUsage": 0.7,
				"datafileSize":  6,
				"dataguard":     false,
				"elapsed":       12059.18,
				"environment":   "TST",
				"featureUsageStats": []interface{}{
					map[string]interface{}{
						"currentlyUsed":    false,
						"detectedUsages":   91,
						"extraFeatureInfo": "",
						"feature":          "ADDM",
						"firstUsageDate":   utils.P("2019-06-24T19:34:20+02:00").Local(),
						"lastUsageDate":    utils.P("2019-11-09T05:48:23+01:00").Local(),
						"product":          "Diagnostics Pack",
					},
					map[string]interface{}{
						"currentlyUsed":    false,
						"detectedUsages":   90,
						"extraFeatureInfo": "",
						"feature":          "AWR Report",
						"firstUsageDate":   utils.P("2019-06-27T17:15:44+02:00").Local(),
						"lastUsageDate":    utils.P("2019-11-09T05:48:23+01:00").Local(),
						"product":          "Diagnostics Pack",
					},
					map[string]interface{}{
						"currentlyUsed":    false,
						"detectedUsages":   7,
						"extraFeatureInfo": "",
						"feature":          "Automatic Workload Repository",
						"firstUsageDate":   utils.P("2019-06-27T19:01:09+02:00").Local(),
						"lastUsageDate":    utils.P("2019-07-02T07:35:05+02:00").Local(),
						"product":          "Diagnostics Pack",
					},
				},
				"ha":             false,
				"hostname":       "test-db",
				"instanceNumber": 1,
				"instanceName":   "ERCOLE1",
				"isCDB":          false,
				"licenses": []interface{}{
					map[string]interface{}{
						"count": 0,
						"name":  "Oracle EXE",
					},
					map[string]interface{}{
						"count": 0.5,
						"name":  "Oracle ENT",
					},
					map[string]interface{}{
						"count": 0,
						"name":  "Oracle STD",
					},
					map[string]interface{}{
						"count": 0,
						"name":  "WebLogic Server Management Pack Enterprise Edition",
					},
					map[string]interface{}{
						"count": 0.5,
						"name":  "Diagnostics Pack",
					},
				},
				"location":     "Germany",
				"memory":       1.484,
				"memoryTarget": 1.484,
				"nCharset":     "AL16UTF16",
				"name":         "ERCOLE",
				"pdbs":         []interface{}{},
				"pgaTarget":    0,
				"psus": []interface{}{
					map[string]interface{}{
						"date":        "2012-04-16",
						"description": "PSU 11.2.0.3.2",
					},
				},
				"patches":    []interface{}{},
				"platform":   "Linux x86 64-bit",
				"rac":        false,
				"sgaMaxSize": 1.484,
				"sgaTarget":  0,
				"schemas": []interface{}{
					map[string]interface{}{
						"indexes": 0,
						"lob":     0,
						"tables":  192,
						"total":   192,
						"user":    "RAF",
					},
					map[string]interface{}{
						"indexes": 0,
						"lob":     0,
						"tables":  0,
						"total":   0,
						"user":    "REMOTE_SCHEDULER_AGENT",
					},
				},
				"segmentAdvisors": []interface{}{
					map[string]interface{}{
						"partitionName":  "iyyiuyyoy",
						"reclaimable":    0.5,
						"recommendation": "32b36a77e7481343ef175483c086859e",
						"segmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						"segmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						"segmentType":    "TABLE",
					},
				},
				"segmentsSize": 3,
				"services":     []interface{}{},
				"status":       "OPEN",
				"tablespaces": []interface{}{
					map[string]interface{}{
						"maxSize":  32767.9844,
						"name":     "SYSTEM",
						"status":   "ONLINE",
						"total":    850,
						"used":     842.875,
						"usedPerc": 2.57,
					},
					map[string]interface{}{
						"maxSize":  32767.9844,
						"name":     "USERS",
						"status":   "ONLINE",
						"total":    1024,
						"used":     576,
						"usedPerc": 1.76,
					},
				},
				"tags":       []interface{}{"foobar"},
				"uniqueName": "ERCOLE",
				"version":    "12.2.0.1.0 Enterprise Edition",
				"work":       1,
				"_id":        utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			map[string]interface{}{
				"addms": []interface{}{
					map[string]interface{}{
						"action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						"benefit":        83.34,
						"finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						"recommendation": "SQL Tuning",
					},
					map[string]interface{}{
						"action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						"benefit":        68.24,
						"finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						"recommendation": "Segment Tuning",
					},
				},
				"asm":        false,
				"allocable":  129,
				"archivelog": true,
				"backups": []interface{}{
					map[string]interface{}{
						"avgBckSize": 13,
						"backupType": "Archivelog",
						"hour":       "01:30",
						"retention":  "1 NUMBERS",
						"weekDays":   []interface{}{"Wednesday"},
					},
					map[string]interface{}{
						"avgBckSize": 45,
						"backupType": "Archivelog",
						"hour":       "03:00",
						"retention":  "1 NUMBERS",
						"weekDays": []interface{}{"Tuesday",
							"Sunday",
							"Monday",
							"Saturday",
							"Wednesday",
						},
					},
				},
				"blockSize":     8192,
				"cpuCount":      2,
				"charset":       "AL32UTF8",
				"createdAt":     utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"dbTime":        184.81,
				"dailyCPUUsage": 0.7,
				"datafileSize":  6,
				"dataguard":     true,
				"elapsed":       12059.18,
				"environment":   "TST",
				"featureUsageStats": []interface{}{
					map[string]interface{}{
						"currentlyUsed":    false,
						"detectedUsages":   91,
						"extraFeatureInfo": "",
						"feature":          "ADDM",
						"firstUsageDate":   utils.P("2019-06-24T19:34:20+02:00").Local(),
						"lastUsageDate":    utils.P("2019-11-09T05:48:23+01:00").Local(),
						"product":          "Diagnostics Pack",
					},
					map[string]interface{}{
						"currentlyUsed":    false,
						"detectedUsages":   90,
						"extraFeatureInfo": "",
						"feature":          "AWR Report",
						"firstUsageDate":   utils.P("2019-06-27T17:15:44+02:00").Local(),
						"lastUsageDate":    utils.P("2019-11-09T05:48:23+01:00").Local(),
						"product":          "Diagnostics Pack",
					},
					map[string]interface{}{
						"currentlyUsed":    false,
						"detectedUsages":   7,
						"extraFeatureInfo": "",
						"feature":          "Automatic Workload Repository",
						"firstUsageDate":   utils.P("2019-06-27T19:01:09+02:00").Local(),
						"lastUsageDate":    utils.P("2019-07-02T07:35:05+02:00").Local(),
						"product":          "Diagnostics Pack",
					},
				},
				"ha":             true,
				"hostname":       "test-db2",
				"instanceNumber": 1,
				"instanceName":   "pokemons1",
				"isCDB":          false,
				"licenses": []interface{}{
					map[string]interface{}{
						"count": 0,
						"name":  "Oracle EXE",
					},
					map[string]interface{}{
						"count": 0.5,
						"name":  "Oracle ENT",
					},
					map[string]interface{}{
						"count": 0,
						"name":  "Oracle STD",
					},
					map[string]interface{}{
						"count": 0,
						"name":  "WebLogic Server Management Pack Enterprise Edition",
					},
					map[string]interface{}{
						"count": 0.5,
						"name":  "Diagnostics Pack",
					},
					map[string]interface{}{
						"count": 0.5,
						"name":  "Real Application Clusters",
					},
				},
				"location":     "Germany",
				"memory":       90.254,
				"memoryTarget": 1.484,
				"nCharset":     "AL16UTF16",
				"name":         "pokemons",
				"pdbs":         []interface{}{},
				"pgaTarget":    53.45,
				"psus": []interface{}{
					map[string]interface{}{
						"date":        "2012-04-16",
						"description": "PSU 11.2.0.3.2",
					},
				},
				"patches":    []interface{}{},
				"platform":   "Linux x86 64-bit",
				"rac":        true,
				"sgaMaxSize": 1.484,
				"sgaTarget":  35.32,
				"schemas": []interface{}{
					map[string]interface{}{
						"indexes": 0,
						"lob":     0,
						"tables":  192,
						"total":   192,
						"user":    "RAF",
					},
					map[string]interface{}{
						"indexes": 0,
						"lob":     0,
						"tables":  0,
						"total":   0,
						"user":    "REMOTE_SCHEDULER_AGENT",
					},
				},
				"segmentAdvisors": []interface{}{
					map[string]interface{}{
						"partitionName":  "iyyiuyyoy",
						"reclaimable":    0.5,
						"recommendation": "32b36a77e7481343ef175483c086859e",
						"segmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						"segmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						"segmentType":    "TABLE",
					},
				},
				"segmentsSize": 3,
				"services":     []interface{}{},
				"status":       "OPEN",
				"tablespaces": []interface{}{
					map[string]interface{}{
						"maxSize":  32767.9844,
						"name":     "SYSTEM",
						"status":   "ONLINE",
						"total":    850,
						"used":     842.875,
						"usedPerc": 2.57,
					},
					map[string]interface{}{
						"maxSize":  32767.9844,
						"name":     "USERS",
						"status":   "ONLINE",
						"total":    1024,
						"used":     576,
						"usedPerc": 1.76,
					},
				},
				"uniqueName": "pokemons",
				"version":    "12.2.0.1.0 Enterprise Edition",
				"work":       1,
				"_id":        utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
