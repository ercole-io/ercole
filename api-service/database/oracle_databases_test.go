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

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
)

func (m *MongodbSuite) TestSearchOracleDatabases() {
	var work float64 = 1
	enabled := false
	name := "ECXSERVER"
	creationdate := utils.P("2019-06-24T17:34:20Z")
	dbtime := 184.81
	dailycpuusage := 0.7
	elapsed := 12059.18

	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, -1, -1, "", "PROD", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := dto.OracleDatabaseResponse{
			Content: []dto.OracleDatabase{},
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := dto.OracleDatabaseResponse{
			Content: []dto.OracleDatabase{},
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)

		expectedOut := dto.OracleDatabaseResponse{
			Content: []dto.OracleDatabase{},
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				Archivelog:     false,
				BlockSize:      8192,
				CPUCount:       2,
				Charset:        "AL32UTF8",
				CreatedAt:      utils.P("2020-04-15T06:46:58.471Z"),
				DatafileSize:   6,
				Dataguard:      false,
				Environment:    "TST",
				Ha:             false,
				Hostname:       "test-db",
				IsCDB:          false,
				Location:       "Germany",
				Memory:         1.484,
				Name:           "ERCOLE",
				PDBs:           []model.OracleDatabasePluggableDatabase{},
				Rac:            false,
				SegmentsSize:   3,
				Services:       []model.OracleDatabaseService{},
				Status:         "OPEN",
				UniqueName:     "ERCOLE",
				Version:        "12.2.0.1.0 Enterprise Edition",
				Work:           &work,
				ID:             utils.Str2oid("5e96ade270c184faca93fe36"),
				InstanceNumber: 1,
				InstanceName:   "ERCOLE1",
				DbID:           0,
				Role:           "",
				Platform:       "Linux x86 64-bit",
				NCharset:       "AL16UTF16",
				SGATarget:      0,
				PGATarget:      0,
				MemoryTarget:   1.484,
				SGAMaxSize:     1.484,
				Allocable:      129,
				Elapsed:        &elapsed,
				DBTime:         &dbtime,
				DailyCPUUsage:  &dailycpuusage,
				ASM:            false,
				Tags:           []string{"foobar"},
				Patches:        []model.OracleDatabasePatch{},
				Tablespaces: []model.OracleDatabaseTablespace{
					{
						MaxSize:  32767.9844,
						Name:     "SYSTEM",
						Status:   "ONLINE",
						Total:    850,
						Used:     842.875,
						UsedPerc: 2.57,
					},
					{
						MaxSize:  32767.9844,
						Name:     "USERS",
						Status:   "ONLINE",
						Total:    1024,
						Used:     576,
						UsedPerc: 1.76,
					},
				},
				Schemas: []model.OracleDatabaseSchema{
					{
						Indexes: 0,
						LOB:     0,
						Tables:  192,
						Total:   192,
						User:    "RAF",
					},
					{
						Indexes: 0,
						LOB:     0,
						Tables:  0,
						Total:   0,
						User:    "REMOTE_SCHEDULER_AGENT",
					},
				},
				Licenses: []model.OracleDatabaseLicense{
					{
						Count:         0.5,
						Name:          "Oracle ENT",
						LicenseTypeID: "A90611",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "Oracle STD",
						LicenseTypeID: "L103399",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "WebLogic Server Management Pack Enterprise Edition",
						LicenseTypeID: "L104095",
						Ignored:       false,
					},
					{
						Count:         0.5,
						Name:          "Diagnostics Pack",
						LicenseTypeID: "A90649",
						Ignored:       false,
					},
				},
				ADDMs: []model.OracleDatabaseAddm{
					{
						Action:         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						Benefit:        83.34,
						Finding:        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						Recommendation: "SQL Tuning",
					},
					{
						Action:         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						Benefit:        68.24,
						Finding:        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						Recommendation: "Segment Tuning",
					},
				},
				SegmentAdvisors: []model.OracleDatabaseSegmentAdvisor{
					{
						PartitionName:  "iyyiuyyoy",
						Reclaimable:    0.5,
						Recommendation: "32b36a77e7481343ef175483c086859e",
						SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						SegmentType:    "TABLE",
					},
				},
				PSUs: []model.OracleDatabasePSU{
					{
						Date:        "2012-04-16",
						Description: "PSU 11.2.0.3.2",
					},
				},
				Backups: []model.OracleDatabaseBackup{
					{
						AvgBckSize: 13,
						BackupType: "Archivelog",
						Hour:       "01:30",
						Retention:  "1 NUMBERS",
						WeekDays:   []string{"Wednesday"},
					},
					{
						AvgBckSize: 45,
						BackupType: "Archivelog",
						Hour:       "03:00",
						Retention:  "1 NUMBERS",
						WeekDays: []string{"Tuesday",
							"Sunday",
							"Monday",
							"Saturday",
							"Wednesday",
						},
					},
				},
				FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{
					{
						CurrentlyUsed:    false,
						DetectedUsages:   91,
						ExtraFeatureInfo: "",
						Feature:          "ADDM",
						FirstUsageDate:   utils.P("2019-06-24T17:34:20Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   90,
						ExtraFeatureInfo: "",
						Feature:          "AWR Report",
						FirstUsageDate:   utils.P("2019-06-27T15:15:44Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   7,
						ExtraFeatureInfo: "",
						Feature:          "Automatic Workload Repository",
						FirstUsageDate:   utils.P("2019-06-27T17:01:09Z"),
						LastUsageDate:    utils.P("2019-07-02T05:35:05Z"),
						Product:          "Diagnostics Pack",
					},
				},
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          false,
				Number:        0,
				Size:          1,
				TotalElements: 2,
				TotalPages:    2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{""}, "memory", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				ADDMs: []model.OracleDatabaseAddm{
					{
						Action:         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						Benefit:        83.34,
						Finding:        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						Recommendation: "SQL Tuning",
					},
					{
						Action:         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						Benefit:        68.24,
						Finding:        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						Recommendation: "Segment Tuning",
					},
				},
				Allocable:  129,
				Archivelog: true,
				ASM:        false,
				Backups: []model.OracleDatabaseBackup{
					{
						AvgBckSize: 13,
						BackupType: "Archivelog",
						Hour:       "01:30",
						Retention:  "1 NUMBERS",
						WeekDays:   []string{"Wednesday"},
					},
					{
						AvgBckSize: 45,
						BackupType: "Archivelog",
						Hour:       "03:00",
						Retention:  "1 NUMBERS",
						WeekDays:   []string{"Tuesday", "Sunday", "Monday", "Saturday", "Wednesday"},
					},
				},
				BlockSize:     8192,
				Charset:       "AL32UTF8",
				CPUCount:      2,
				CreatedAt:     utils.P("2020-05-06T11:39:23.259Z"),
				DailyCPUUsage: &dailycpuusage,
				DatafileSize:  6,
				Dataguard:     true,
				DbID:          0,
				DBTime:        &dbtime,
				Elapsed:       &elapsed,
				Environment:   "TST",
				FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{
					{
						CurrentlyUsed:    false,
						DetectedUsages:   91,
						ExtraFeatureInfo: "",
						Feature:          "ADDM",
						FirstUsageDate:   utils.P("2019-06-24T17:34:20Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   90,
						ExtraFeatureInfo: "",
						Feature:          "AWR Report",
						FirstUsageDate:   utils.P("2019-06-27T15:15:44Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   7,
						ExtraFeatureInfo: "",
						Feature:          "Automatic Workload Repository",
						FirstUsageDate:   utils.P("2019-06-27T17:01:09Z"),
						LastUsageDate:    utils.P("2019-07-02T05:35:05Z"),
						Product:          "Diagnostics Pack",
					},
				},
				Ha:             true,
				Hostname:       "test-db2",
				ID:             utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
				InstanceName:   "pokemons1",
				InstanceNumber: 1,
				IsCDB:          true,
				Licenses: []model.OracleDatabaseLicense{
					{
						Count:         0.5,
						Ignored:       false,
						LicenseTypeID: "A90611",
						Name:          "Oracle ENT",
					},
					{
						Count:         0,
						Ignored:       false,
						LicenseTypeID: "L103399",
						Name:          "Oracle STD",
					},
					{
						Count:         0,
						Ignored:       false,
						LicenseTypeID: "L104095",
						Name:          "WebLogic Server Management Pack Enterprise Edition",
					},
					{
						Count:         0.5,
						Ignored:       false,
						LicenseTypeID: "A90649",
						Name:          "Diagnostics Pack",
					},
					{
						Count:         0.5,
						Ignored:       false,
						LicenseTypeID: "A90619",
						Name:          "Real Application Clusters",
					},
				},
				Location:     "Germany",
				Memory:       90.254,
				MemoryTarget: 1.484,
				NCharset:     "AL16UTF16",
				Name:         "pokemons",
				Patches:      []model.OracleDatabasePatch{},
				PDBs: []model.OracleDatabasePluggableDatabase{
					{
						Name: "PDB1",
					},
				},
				PGATarget: 53.45,
				Platform:  "Linux x86 64-bit",
				PSUs: []model.OracleDatabasePSU{
					{
						Date:        "2012-04-16",
						Description: "PSU 11.2.0.3.2",
					},
				},
				Rac:  true,
				Role: "",
				Schemas: []model.OracleDatabaseSchema{
					{
						Indexes: 0,
						LOB:     0,
						Tables:  192,
						Total:   192,
						User:    "RAF",
					},
					{
						Indexes: 0,
						LOB:     0,
						Tables:  0,
						Total:   0,
						User:    "REMOTE_SCHEDULER_AGENT",
					},
				},
				SegmentAdvisors: []model.OracleDatabaseSegmentAdvisor{
					{
						PartitionName:  "iyyiuyyoy",
						Reclaimable:    0.5,
						Recommendation: "32b36a77e7481343ef175483c086859e",
						SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						SegmentType:    "TABLE",
					},
				},
				SegmentsSize: 3,
				Services: []model.OracleDatabaseService{
					{
						CreationDate: &creationdate,
						Enabled:      &enabled,
						Name:         &name,
					},
				},
				SGAMaxSize: 1.484,
				SGATarget:  35.32,
				Status:     "OPEN",
				Tablespaces: []model.OracleDatabaseTablespace{
					{
						MaxSize:  32767.9844,
						Name:     "SYSTEM",
						Status:   "ONLINE",
						Total:    850,
						Used:     842.875,
						UsedPerc: 2.57,
					},
					{
						MaxSize:  32767.9844,
						Name:     "USERS",
						Status:   "ONLINE",
						Total:    1024,
						Used:     576,
						UsedPerc: 1.76,
					},
				},
				UniqueName: "pokemons",
				Version:    "12.2.0.1.0 Enterprise Edition",
				Work:       &work,
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          2,
				TotalElements: 2,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut.Content[0]), utils.ToJSON(out.Content[0]))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{}

		expectedOut := dto.OracleDatabaseResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{"pokemon", "test-db2"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				ADDMs: []model.OracleDatabaseAddm{
					{
						Action:         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						Benefit:        83.34,
						Finding:        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						Recommendation: "SQL Tuning",
					},
					{
						Action:         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						Benefit:        68.24,
						Finding:        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						Recommendation: "Segment Tuning",
					},
				},
				ASM:        false,
				Allocable:  129,
				Archivelog: true,
				Backups: []model.OracleDatabaseBackup{
					{
						AvgBckSize: 13,
						BackupType: "Archivelog",
						Hour:       "01:30",
						Retention:  "1 NUMBERS",
						WeekDays:   []string{"Wednesday"},
					},
					{
						AvgBckSize: 45,
						BackupType: "Archivelog",
						Hour:       "03:00",
						Retention:  "1 NUMBERS",
						WeekDays: []string{"Tuesday",
							"Sunday",
							"Monday",
							"Saturday",
							"Wednesday",
						},
					},
				},
				BlockSize:     8192,
				CPUCount:      2,
				Charset:       "AL32UTF8",
				CreatedAt:     utils.P("2020-05-06T11:39:23.259Z"),
				DBTime:        &dbtime,
				DailyCPUUsage: &dailycpuusage,
				DatafileSize:  6,
				Dataguard:     true,
				Elapsed:       &elapsed,
				Environment:   "TST",
				FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{
					{
						CurrentlyUsed:    false,
						DetectedUsages:   91,
						ExtraFeatureInfo: "",
						Feature:          "ADDM",
						FirstUsageDate:   utils.P("2019-06-24T17:34:20Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   90,
						ExtraFeatureInfo: "",
						Feature:          "AWR Report",
						FirstUsageDate:   utils.P("2019-06-27T15:15:44Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   7,
						ExtraFeatureInfo: "",
						Feature:          "Automatic Workload Repository",
						FirstUsageDate:   utils.P("2019-06-27T17:01:09Z"),
						LastUsageDate:    utils.P("2019-07-02T05:35:05Z"),
						Product:          "Diagnostics Pack",
					},
				},
				Ha:             true,
				Hostname:       "test-db2",
				InstanceNumber: 1,
				InstanceName:   "pokemons1",
				IsCDB:          true,
				Licenses: []model.OracleDatabaseLicense{
					{
						Count:         0.5,
						Name:          "Oracle ENT",
						LicenseTypeID: "A90611",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "Oracle STD",
						LicenseTypeID: "L103399",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "WebLogic Server Management Pack Enterprise Edition",
						LicenseTypeID: "L104095",
						Ignored:       false,
					},
					{
						Count:         0.5,
						Name:          "Diagnostics Pack",
						LicenseTypeID: "A90649",
						Ignored:       false,
					},
					{
						Count:         0.5,
						Name:          "Real Application Clusters",
						LicenseTypeID: "A90619",
						Ignored:       false,
					},
				},
				Location:     "Germany",
				Memory:       90.254,
				MemoryTarget: 1.484,
				NCharset:     "AL16UTF16",
				Name:         "pokemons",
				PDBs: []model.OracleDatabasePluggableDatabase{
					{
						Name: "PDB1",
					},
				},
				PGATarget: 53.45,
				PSUs: []model.OracleDatabasePSU{
					{
						Date:        "2012-04-16",
						Description: "PSU 11.2.0.3.2",
					},
				},
				Patches:    []model.OracleDatabasePatch{},
				Platform:   "Linux x86 64-bit",
				Rac:        true,
				SGAMaxSize: 1.484,
				SGATarget:  35.32,
				Schemas: []model.OracleDatabaseSchema{
					{
						Indexes: 0,
						LOB:     0,
						Tables:  192,
						Total:   192,
						User:    "RAF",
					},
					{
						Indexes: 0,
						LOB:     0,
						Tables:  0,
						Total:   0,
						User:    "REMOTE_SCHEDULER_AGENT",
					},
				},
				SegmentAdvisors: []model.OracleDatabaseSegmentAdvisor{
					{
						PartitionName:  "iyyiuyyoy",
						Reclaimable:    0.5,
						Recommendation: "32b36a77e7481343ef175483c086859e",
						SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						SegmentType:    "TABLE",
					},
				},
				SegmentsSize: 3,
				Services: []model.OracleDatabaseService{
					{
						CreationDate: &creationdate,
						Enabled:      &enabled,
						Name:         &name,
					},
				},
				Status: "OPEN",
				Tablespaces: []model.OracleDatabaseTablespace{
					{
						MaxSize:  32767.9844,
						Name:     "SYSTEM",
						Status:   "ONLINE",
						Total:    850,
						Used:     842.875,
						UsedPerc: 2.57,
					},
					{
						MaxSize:  32767.9844,
						Name:     "USERS",
						Status:   "ONLINE",
						Total:    1024,
						Used:     576,
						UsedPerc: 1.76,
					},
				},
				UniqueName: "pokemons",
				Version:    "12.2.0.1.0 Enterprise Edition",
				Work:       &work,
				ID:         utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          1,
				TotalElements: 1,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{""}, "memory", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				ADDMs: []model.OracleDatabaseAddm{
					{
						Action:         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						Benefit:        83.34,
						Finding:        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						Recommendation: "SQL Tuning",
					},
					{
						Action:         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						Benefit:        68.24,
						Finding:        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						Recommendation: "Segment Tuning",
					},
				},
				ASM:        false,
				Allocable:  129,
				Archivelog: false,
				Backups: []model.OracleDatabaseBackup{
					{
						AvgBckSize: 13,
						BackupType: "Archivelog",
						Hour:       "01:30",
						Retention:  "1 NUMBERS",
						WeekDays:   []string{"Wednesday"},
					},
					{
						AvgBckSize: 45,
						BackupType: "Archivelog",
						Hour:       "03:00",
						Retention:  "1 NUMBERS",
						WeekDays: []string{"Tuesday",
							"Sunday",
							"Monday",
							"Saturday",
							"Wednesday",
						},
					},
				},
				BlockSize:     8192,
				CPUCount:      2,
				Charset:       "AL32UTF8",
				CreatedAt:     utils.P("2020-04-15T06:46:58.471Z"),
				DBTime:        &dbtime,
				DailyCPUUsage: &dailycpuusage,
				DatafileSize:  6,
				Dataguard:     false,
				Elapsed:       &elapsed,
				Environment:   "TST",
				FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{
					{
						CurrentlyUsed:    false,
						DetectedUsages:   91,
						ExtraFeatureInfo: "",
						Feature:          "ADDM",
						FirstUsageDate:   utils.P("2019-06-24T17:34:20Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   90,
						ExtraFeatureInfo: "",
						Feature:          "AWR Report",
						FirstUsageDate:   utils.P("2019-06-27T15:15:44Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   7,
						ExtraFeatureInfo: "",
						Feature:          "Automatic Workload Repository",
						FirstUsageDate:   utils.P("2019-06-27T17:01:09Z"),
						LastUsageDate:    utils.P("2019-07-02T05:35:05Z"),
						Product:          "Diagnostics Pack",
					},
				},
				Ha:             false,
				Hostname:       "test-db",
				InstanceNumber: 1,
				InstanceName:   "ERCOLE1",
				IsCDB:          false,
				Licenses: []model.OracleDatabaseLicense{
					{
						Count:         0.5,
						Name:          "Oracle ENT",
						LicenseTypeID: "A90611",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "Oracle STD",
						LicenseTypeID: "L103399",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "WebLogic Server Management Pack Enterprise Edition",
						LicenseTypeID: "L104095",
						Ignored:       false,
					},
					{
						Count:         0.5,
						Name:          "Diagnostics Pack",
						LicenseTypeID: "A90649",
						Ignored:       false,
					},
				},
				Location:     "Germany",
				Memory:       1.484,
				MemoryTarget: 1.484,
				NCharset:     "AL16UTF16",
				Name:         "ERCOLE",
				PDBs:         []model.OracleDatabasePluggableDatabase{},
				PGATarget:    0,
				PSUs: []model.OracleDatabasePSU{
					{
						Date:        "2012-04-16",
						Description: "PSU 11.2.0.3.2",
					},
				},
				Patches:    []model.OracleDatabasePatch{},
				Platform:   "Linux x86 64-bit",
				Rac:        false,
				SGAMaxSize: 1.484,
				SGATarget:  0,
				Schemas: []model.OracleDatabaseSchema{
					{
						Indexes: 0,
						LOB:     0,
						Tables:  192,
						Total:   192,
						User:    "RAF",
					},
					{
						Indexes: 0,
						LOB:     0,
						Tables:  0,
						Total:   0,
						User:    "REMOTE_SCHEDULER_AGENT",
					},
				},
				SegmentAdvisors: []model.OracleDatabaseSegmentAdvisor{
					{
						PartitionName:  "iyyiuyyoy",
						Reclaimable:    0.5,
						Recommendation: "32b36a77e7481343ef175483c086859e",
						SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						SegmentType:    "TABLE",
					},
				},
				SegmentsSize: 3,
				Services:     []model.OracleDatabaseService{},
				Status:       "OPEN",
				Tablespaces: []model.OracleDatabaseTablespace{
					{
						MaxSize:  32767.9844,
						Name:     "SYSTEM",
						Status:   "ONLINE",
						Total:    850,
						Used:     842.875,
						UsedPerc: 2.57,
					},
					{
						MaxSize:  32767.9844,
						Name:     "USERS",
						Status:   "ONLINE",
						Total:    1024,
						Used:     576,
						UsedPerc: 1.76,
					},
				},
				Tags:       []string{"foobar"},
				UniqueName: "ERCOLE",
				Version:    "12.2.0.1.0 Enterprise Edition",
				Work:       &work,
				ID:         utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			{
				ADDMs: []model.OracleDatabaseAddm{
					{
						Action:         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						Benefit:        83.34,
						Finding:        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						Recommendation: "SQL Tuning",
					},
					{
						Action:         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
						Benefit:        68.24,
						Finding:        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
						Recommendation: "Segment Tuning",
					},
				},
				ASM:        false,
				Allocable:  129,
				Archivelog: true,
				Backups: []model.OracleDatabaseBackup{
					{
						AvgBckSize: 13,
						BackupType: "Archivelog",
						Hour:       "01:30",
						Retention:  "1 NUMBERS",
						WeekDays:   []string{"Wednesday"},
					},
					{
						AvgBckSize: 45,
						BackupType: "Archivelog",
						Hour:       "03:00",
						Retention:  "1 NUMBERS",
						WeekDays: []string{"Tuesday",
							"Sunday",
							"Monday",
							"Saturday",
							"Wednesday",
						},
					},
				},
				BlockSize:     8192,
				CPUCount:      2,
				Charset:       "AL32UTF8",
				CreatedAt:     utils.P("2020-05-06T11:39:23.259Z"),
				DBTime:        &dbtime,
				DailyCPUUsage: &dailycpuusage,
				DatafileSize:  6,
				Dataguard:     true,
				Elapsed:       &elapsed,
				Environment:   "TST",
				FeatureUsageStats: []model.OracleDatabaseFeatureUsageStat{
					{
						CurrentlyUsed:    false,
						DetectedUsages:   91,
						ExtraFeatureInfo: "",
						Feature:          "ADDM",
						FirstUsageDate:   utils.P("2019-06-24T17:34:20Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   90,
						ExtraFeatureInfo: "",
						Feature:          "AWR Report",
						FirstUsageDate:   utils.P("2019-06-27T15:15:44Z"),
						LastUsageDate:    utils.P("2019-11-09T04:48:23Z"),
						Product:          "Diagnostics Pack",
					},
					{
						CurrentlyUsed:    false,
						DetectedUsages:   7,
						ExtraFeatureInfo: "",
						Feature:          "Automatic Workload Repository",
						FirstUsageDate:   utils.P("2019-06-27T17:01:09Z"),
						LastUsageDate:    utils.P("2019-07-02T05:35:05Z"),
						Product:          "Diagnostics Pack",
					},
				},
				Ha:             true,
				Hostname:       "test-db2",
				InstanceNumber: 1,
				InstanceName:   "pokemons1",
				IsCDB:          true,
				Licenses: []model.OracleDatabaseLicense{
					{
						Count:         0.5,
						Name:          "Oracle ENT",
						LicenseTypeID: "A90611",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "Oracle STD",
						LicenseTypeID: "L103399",
						Ignored:       false,
					},
					{
						Count:         0,
						Name:          "WebLogic Server Management Pack Enterprise Edition",
						LicenseTypeID: "L104095",
						Ignored:       false,
					},
					{
						Count:         0.5,
						Name:          "Diagnostics Pack",
						LicenseTypeID: "A90649",
						Ignored:       false,
					},
					{
						Count:         0.5,
						Name:          "Real Application Clusters",
						LicenseTypeID: "A90619",
						Ignored:       false,
					},
				},
				Location:     "Germany",
				Memory:       90.254,
				MemoryTarget: 1.484,
				NCharset:     "AL16UTF16",
				Name:         "pokemons",
				PDBs: []model.OracleDatabasePluggableDatabase{
					{
						Name: "PDB1",
					},
				},
				PGATarget: 53.45,
				PSUs: []model.OracleDatabasePSU{
					{
						Date:        "2012-04-16",
						Description: "PSU 11.2.0.3.2",
					},
				},
				Patches:    []model.OracleDatabasePatch{},
				Platform:   "Linux x86 64-bit",
				Rac:        true,
				SGAMaxSize: 1.484,
				SGATarget:  35.32,
				Schemas: []model.OracleDatabaseSchema{
					{
						Indexes: 0,
						LOB:     0,
						Tables:  192,
						Total:   192,
						User:    "RAF",
					},
					{
						Indexes: 0,
						LOB:     0,
						Tables:  0,
						Total:   0,
						User:    "REMOTE_SCHEDULER_AGENT",
					},
				},
				SegmentAdvisors: []model.OracleDatabaseSegmentAdvisor{
					{
						PartitionName:  "iyyiuyyoy",
						Reclaimable:    0.5,
						Recommendation: "32b36a77e7481343ef175483c086859e",
						SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
						SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
						SegmentType:    "TABLE",
					},
				},
				SegmentsSize: 3,
				Services: []model.OracleDatabaseService{
					{
						CreationDate: &creationdate,
						Enabled:      &enabled,
						Name:         &name,
					},
				},
				Status: "OPEN",
				Tablespaces: []model.OracleDatabaseTablespace{
					{
						MaxSize:  32767.9844,
						Name:     "SYSTEM",
						Status:   "ONLINE",
						Total:    850,
						Used:     842.875,
						UsedPerc: 2.57,
					},
					{
						MaxSize:  32767.9844,
						Name:     "USERS",
						Status:   "ONLINE",
						Total:    1024,
						Used:     576,
						UsedPerc: 1.76,
					},
				},
				UniqueName: "pokemons",
				Version:    "12.2.0.1.0 Enterprise Edition",
				Work:       &work,
				ID:         utils.Str2oid("5eb2a1eba77f5e4badf8a2cc"),
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          2,
				TotalElements: 2,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
