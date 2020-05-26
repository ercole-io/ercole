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

	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchDatabases() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{""}, "", false, -1, -1, "", "PROD", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{""}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Content": []interface{}{
					map[string]interface{}{
						"ArchiveLogStatus": false,
						"BlockSize":        "8192",
						"CPUCount":         "2",
						"Charset":          "AL32UTF8",
						"CreatedAt":        utils.P("2020-04-15T08:46:58.471+02:00").Local(),
						"DatafileSize":     "6",
						"Dataguard":        false,
						"Environment":      "TST",
						"HA":               false,
						"Hostname":         "test-db",
						"Location":         "Germany",
						"Memory":           1.484,
						"Name":             "ERCOLE",
						"RAC":              false,
						"SegmentsSize":     "3",
						"Status":           "OPEN",
						"UniqueName":       "ERCOLE",
						"Version":          "12.2.0.1.0 Enterprise Edition",
						"Work":             "1",
						"_id":              "5e96ade270c184faca93fe36",
					},
				},
				"Metadata": map[string]interface{}{
					"Empty":         false,
					"First":         true,
					"Last":          false,
					"Number":        0,
					"Size":          1,
					"TotalElements": 2,
					"TotalPages":    2,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{""}, "Memory", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"ArchiveLogStatus": true,
				"BlockSize":        "8192",
				"CPUCount":         "2",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"DatafileSize":     "6",
				"Dataguard":        true,
				"Environment":      "TST",
				"HA":               true,
				"Hostname":         "test-db2",
				"Location":         "Germany",
				"Memory":           90.254,
				"Name":             "pokemons",
				"RAC":              true,
				"SegmentsSize":     "3",
				"Status":           "OPEN",
				"UniqueName":       "pokemons",
				"Version":          "12.2.0.1.0 Enterprise Edition",
				"Work":             "1",
				"_id":              utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
			map[string]interface{}{
				"ArchiveLogStatus": false,
				"BlockSize":        "8192",
				"CPUCount":         "2",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"DatafileSize":     "6",
				"Dataguard":        false,
				"Environment":      "TST",
				"HA":               false,
				"Hostname":         "test-db",
				"Location":         "Germany",
				"Memory":           1.484,
				"Name":             "ERCOLE",
				"RAC":              false,
				"SegmentsSize":     "3",
				"Status":           "OPEN",
				"UniqueName":       "ERCOLE",
				"Version":          "12.2.0.1.0 Enterprise Edition",
				"Work":             "1",
				"_id":              utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{"pokemon", "test-db2"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"ArchiveLogStatus": true,
				"BlockSize":        "8192",
				"CPUCount":         "2",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"DatafileSize":     "6",
				"Dataguard":        true,
				"Environment":      "TST",
				"HA":               true,
				"Hostname":         "test-db2",
				"Location":         "Germany",
				"Memory":           90.254,
				"Name":             "pokemons",
				"RAC":              true,
				"SegmentsSize":     "3",
				"Status":           "OPEN",
				"UniqueName":       "pokemons",
				"Version":          "12.2.0.1.0 Enterprise Edition",
				"Work":             "1",
				"_id":              utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("notfullmode", func(t *testing.T) {
		out, err := m.db.SearchDatabases(false, []string{""}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"ArchiveLogStatus": false,
				"BlockSize":        "8192",
				"CPUCount":         "2",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"DatafileSize":     "6",
				"Dataguard":        false,
				"Environment":      "TST",
				"HA":               false,
				"Hostname":         "test-db",
				"Location":         "Germany",
				"Memory":           1.484,
				"Name":             "ERCOLE",
				"RAC":              false,
				"SegmentsSize":     "3",
				"Status":           "OPEN",
				"UniqueName":       "ERCOLE",
				"Version":          "12.2.0.1.0 Enterprise Edition",
				"Work":             "1",
				"_id":              utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			map[string]interface{}{
				"ArchiveLogStatus": true,
				"BlockSize":        "8192",
				"CPUCount":         "2",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"DatafileSize":     "6",
				"Dataguard":        true,
				"Environment":      "TST",
				"HA":               true,
				"Hostname":         "test-db2",
				"Location":         "Germany",
				"Memory":           90.254,
				"Name":             "pokemons",
				"RAC":              true,
				"SegmentsSize":     "3",
				"Status":           "OPEN",
				"UniqueName":       "pokemons",
				"Version":          "12.2.0.1.0 Enterprise Edition",
				"Work":             "1",
				"_id":              utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchDatabases(true, []string{""}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"ADDMs": []interface{}{
					map[string]interface{}{
						"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						"Benefit":        "83.34",
						"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						"Recommendation": "SQL Tuning",
					},
					map[string]interface{}{
						"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						"Benefit":        "68.24",
						"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						"Recommendation": "Segment Tuning",
					},
				},
				"ASM":              false,
				"Allocated":        "129",
				"ArchiveLogStatus": false,
				"Archivelog":       "NOARCHIVELOG",
				"Backups": []interface{}{
					map[string]interface{}{
						"AvgBckSize": "13",
						"BackupType": "Archivelog",
						"Hour":       "01:30",
						"Retention":  "1 NUMBERS",
						"WeekDays":   "Wednesday"},
					map[string]interface{}{
						"AvgBckSize": "45",
						"BackupType": "Archivelog",
						"Hour":       "03:00",
						"Retention":  "1 NUMBERS",
						"WeekDays":   "Tuesday,Sunday,Monday,Saturday,Wednesday",
					},
				},
				"BlockSize":     "8192",
				"CPUCount":      "2",
				"Charset":       "AL32UTF8",
				"CreatedAt":     utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"DBTime":        "184.81",
				"DailyCPUUsage": "0.7",
				"DatafileSize":  "6",
				"Dataguard":     false,
				"Elapsed":       "12059.18",
				"Environment":   "TST",
				"Features": []interface{}{
					map[string]interface{}{
						"Name":   "Exadata",
						"Status": false,
					},
					map[string]interface{}{"Name": "Diagnostics Pack",
						"Status": true,
					},
				},
				"Features2": []interface{}{
					map[string]interface{}{
						"CurrentlyUsed":    false,
						"DetectedUsages":   91,
						"ExtraFeatureInfo": "",
						"Feature":          "ADDM",
						"FirstUsageDate":   "2019-06-24 17:34:20",
						"LastUsageDate":    "2019-11-09 04:48:23",
						"Product":          "Diagnostics Pack"},
					map[string]interface{}{
						"CurrentlyUsed":    false,
						"DetectedUsages":   90,
						"ExtraFeatureInfo": "",
						"Feature":          "AWR Report",
						"FirstUsageDate":   "2019-06-27 15:15:44",
						"LastUsageDate":    "2019-11-09 04:48:23",
						"Product":          "Diagnostics Pack",
					},
					map[string]interface{}{
						"CurrentlyUsed":    false,
						"DetectedUsages":   7,
						"ExtraFeatureInfo": "",
						"Feature":          "Automatic Workload Repository",
						"FirstUsageDate":   "2019-06-27 17:01:09",
						"LastUsageDate":    "2019-07-02 05:35:05",
						"Product":          "Diagnostics Pack",
					},
				},
				"HA":             false,
				"Hostname":       "test-db",
				"InstanceNumber": "1",
				"LastPSUs": []interface{}{
					map[string]interface{}{
						"Date":        "2012-04-16",
						"Description": "PSU 11.2.0.3.2",
					},
				},
				"Licenses": []interface{}{map[string]interface{}{
					"Count": 0,
					"Name":  "Oracle EXE",
				},
					map[string]interface{}{
						"Count": 0.5,
						"Name":  "Oracle ENT",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "Oracle STD",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "WebLogic Server Management Pack Enterprise Edition",
					},
					map[string]interface{}{"Count": 0.5,
						"Name": "Diagnostics Pack",
					},
				},
				"Location":     "Germany",
				"Memory":       1.484,
				"MemoryTarget": "1.484",
				"NCharset":     "AL16UTF16",
				"Name":         "ERCOLE",
				"PGATarget":    "0.",
				"Patches":      []interface{}{},
				"Platform":     "Linux x86 64-bit",
				"RAC":          false,
				"SGAMaxSize":   "1.484",
				"SGATarget":    "0.",
				"Schemas": []interface{}{
					map[string]interface{}{
						"Database": "ERCOLE",
						"Indexes":  0,
						"LOB":      0,
						"Tables":   192,
						"Total":    192,
						"User":     "RAF",
					},
					map[string]interface{}{
						"Database": "ERCOLE",
						"Indexes":  0,
						"LOB":      0,
						"Tables":   0,
						"Total":    0,
						"User":     "REMOTE_SCHEDULER_AGENT",
					},
				},
				"SegmentAdvisors": []interface{}{
					map[string]interface{}{
						"PartitionName":  "iyyiuyyoy",
						"Reclaimable":    "<1",
						"Recommendation": "32b36a77e7481343ef175483c086859e",
						"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						"SegmentType":    "TABLE",
					},
				},
				"SegmentsSize": "3",
				"Status":       "OPEN",
				"Tablespaces": []interface{}{
					map[string]interface{}{
						"Database": "ERCOLE",
						"MaxSize":  "32767.9844",
						"Name":     "SYSTEM",
						"Status":   "ONLINE",
						"Total":    "850",
						"Used":     "842.875",
						"UsedPerc": "2.57",
					},
					map[string]interface{}{
						"Database": "ERCOLE",
						"MaxSize":  "32767.9844",
						"Name":     "USERS",
						"Status":   "ONLINE",
						"Total":    "1024",
						"Used":     "576",
						"UsedPerc": "1.76",
					},
				},
				"Tags":       []interface{}{"foobar"},
				"UniqueName": "ERCOLE",
				"Used":       "6",
				"Version":    "12.2.0.1.0 Enterprise Edition",
				"Work":       "1",
				"_id":        "5e96ade270c184faca93fe36",
			},
			map[string]interface{}{
				"ADDMs":            []interface{}{},
				"ASM":              false,
				"Allocated":        "129",
				"ArchiveLogStatus": true,
				"Archivelog":       "ARCHIVELOG",
				"Backups":          []interface{}{},
				"BlockSize":        "8192",
				"CPUCount":         "2",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-05-06T13:39:23.259+02:00").Local(),
				"DBTime":           "184.81",
				"DailyCPUUsage":    "1.3",
				"DatafileSize":     "6",
				"Dataguard":        true,
				"Elapsed":          "12059.18",
				"Environment":      "TST",
				"Features": []interface{}{
					map[string]interface{}{
						"Name":   "WebLogic Server Management Pack Enterprise Edition",
						"Status": false,
					},
					map[string]interface{}{
						"Name":   "Tuning Pack",
						"Status": false,
					},
					map[string]interface{}{
						"Name":   "Provisioning and Patch Automation Pack for Database",
						"Status": false,
					},
					map[string]interface{}{
						"Name":   "Label Security",
						"Status": false,
					},
					map[string]interface{}{
						"Name":   "HW",
						"Status": false,
					},
					map[string]interface{}{
						"Name":   "GoldenGate",
						"Status": false,
					},
					map[string]interface{}{
						"Name":   "Exadata",
						"Status": false,
					},
					map[string]interface{}{
						"Name":   "Diagnostics Pack",
						"Status": true,
					},
					map[string]interface{}{
						"Name":   "Real Application Clusters",
						"Status": true,
					},
				},
				"Features2": []interface{}{
					map[string]interface{}{
						"CurrentlyUsed":    false,
						"DetectedUsages":   91,
						"ExtraFeatureInfo": "",
						"Feature":          "ADDM",
						"FirstUsageDate":   "2019-06-24 17:34:20",
						"LastUsageDate":    "2019-11-09 04:48:23",
						"Product":          "Diagnostics Pack",
					},
					map[string]interface{}{
						"CurrentlyUsed":    false,
						"DetectedUsages":   90,
						"ExtraFeatureInfo": "",
						"Feature":          "AWR Report",
						"FirstUsageDate":   "2019-06-27 15:15:44",
						"LastUsageDate":    "2019-11-09 04:48:23",
						"Product":          "Diagnostics Pack",
					},
					map[string]interface{}{
						"CurrentlyUsed":    false,
						"DetectedUsages":   7,
						"ExtraFeatureInfo": "",
						"Feature":          "Automatic Workload Repository",
						"FirstUsageDate":   "2019-06-27 17:01:09",
						"LastUsageDate":    "2019-07-02 05:35:05",
						"Product":          "Diagnostics Pack",
					},
				},
				"HA":             true,
				"Hostname":       "test-db2",
				"InstanceNumber": "1",
				"LastPSUs":       []interface{}{},
				"Licenses": []interface{}{
					map[string]interface{}{
						"Count": 0,
						"Name":  "Oracle EXE",
					},
					map[string]interface{}{
						"Count": 0.5,
						"Name":  "Oracle ENT",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "Oracle STD",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "WebLogic Server Management Pack Enterprise Edition",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "RAC or RAC One Node",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "Provisioning and Patch Automation Pack",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "Provisioning and Patch Automation Pack for Database",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "Pillar Storage",
					},
					map[string]interface{}{
						"Count": 0,
						"Name":  "Partitioning",
					},
					map[string]interface{}{
						"Count": 0.5,
						"Name":  "Diagnostics Pack",
					},
				},
				"Location":     "Germany",
				"Memory":       90.254,
				"MemoryTarget": "1.484",
				"NCharset":     "AL16UTF16",
				"Name":         "pokemons",
				"PGATarget":    "53.45",
				"Patches":      []interface{}{},
				"Platform":     "Linux x86 64-bit",
				"RAC":          true,
				"SGAMaxSize":   "1.484",
				"SGATarget":    "35.32",
				"Schemas": []interface{}{
					map[string]interface{}{
						"Database": "pokemons",
						"Indexes":  0,
						"LOB":      0,
						"Tables":   192,
						"Total":    192,
						"User":     "RAF",
					},
					map[string]interface{}{
						"Database": "pokemons",
						"Indexes":  0,
						"LOB":      0,
						"Tables":   0,
						"Total":    0,
						"User":     "REMOTE_SCHEDULER_AGENT",
					},
					map[string]interface{}{
						"Database": "pokemons",
						"Indexes":  0,
						"LOB":      0,
						"Tables":   0,
						"Total":    0,
						"User":     "SYS$UMF",
					},
				},
				"SegmentAdvisors": []interface{}{},
				"SegmentsSize":    "3",
				"Status":          "OPEN",
				"Tablespaces": []interface{}{
					map[string]interface{}{
						"Database": "pokemons",
						"MaxSize":  "32767.9844",
						"Name":     "SYSTEM",
						"Status":   "ONLINE",
						"Total":    "850",
						"Used":     "842.875",
						"UsedPerc": "2.57",
					},
					map[string]interface{}{
						"Database": "pokemons",
						"MaxSize":  "32767.9844",
						"Name":     "USERS",
						"Status":   "ONLINE",
						"Total":    "1024",
						"Used":     "576",
						"UsedPerc": "1.76",
					},
					map[string]interface{}{
						"Database": "pokemons",
						"MaxSize":  "32767.9844",
						"Name":     "UNDOTBS1",
						"Status":   "ONLINE",
						"Total":    "895",
						"Used":     "41.25",
						"UsedPerc": "0.13",
					},
					map[string]interface{}{
						"Database": "pokemons",
						"MaxSize":  "32767.9844",
						"Name":     "SYSAUX",
						"Status":   "ONLINE",
						"Total":    "1870",
						"Used":     "1749.625",
						"UsedPerc": "5.34",
					},
				},
				"UniqueName": "pokemons",
				"Used":       "6",
				"Version":    "12.2.0.1.0 Enterprise Edition",
				"Work":       "1",
				"_id":        "5eb2a1eba77f5e4badf8a2cc",
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
