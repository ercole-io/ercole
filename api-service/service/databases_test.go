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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchDatabases_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	oracleDbs := []map[string]interface{}{
		{
			"name":        "pippo",
			"version":     "",
			"hostname":    "",
			"environment": "",
			"charset":     "",

			"memory":       42.42,
			"datafileSize": 75.42,
			"segmentsSize": 99.42,
			"archivelog":   true,
			"ha":           false,
			"dataguard":    true,
		},
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	db.EXPECT().SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(oracleDbs, nil)

	mysqlInstances := []dto.MySQLInstance{
		{
			Hostname:    "pluto",
			Location:    "Cuba",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "mysql",
				Version:            "",
				Edition:            "",
				Platform:           "",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           1,
				ThreadsConcurrency: 2,
				BufferPoolSize:     43008,
				LogBufferSize:      4,
				SortBufferSize:     5,
				ReadOnly:           false,
				LogBin:             true,
				HighAvailability:   false,
				UUID:               "000000000000000000000000",
				IsMaster:           true,
				SlaveUUIDs:         []string{"111111111111111111111111"},
				IsSlave:            false,
				MasterUUID:         new(string),
				Databases:          []model.MySQLDatabase{{Name: "", Charset: "", Collation: "", Encrypted: false}},
				TableSchemas:       []model.MySQLTableSchema{{Name: "", Engine: "", Allocation: 24576}},
				SegmentAdvisors:    []model.MySQLSegmentAdvisor{{TableSchema: "", TableName: "", Engine: "", Allocation: 76, Data: 0, Index: 0, Free: 0}},
			},
		},
	}

	globalFilter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	actual, err := as.SearchDatabases(globalFilter)
	require.NoError(t, err)

	expected := []dto.Database{
		{
			Name:             "pippo",
			Type:             "Oracle/Database",
			Version:          "",
			Hostname:         "",
			Environment:      "",
			Charset:          "",
			Memory:           42.42,
			DatafileSize:     75.42,
			SegmentsSize:     99.42,
			Archivelog:       true,
			HighAvailability: false,
			DisasterRecovery: true,
		},
		{
			Name:             "mysql",
			Type:             "Oracle/MySQL",
			Version:          "",
			Hostname:         "pluto",
			Environment:      "TST",
			Charset:          "",
			Memory:           42.0,
			DatafileSize:     0,
			SegmentsSize:     24.0,
			Archivelog:       true,
			HighAvailability: false,
			DisasterRecovery: true,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestSearchDatabasesAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}

	oracleDbs := []map[string]interface{}{
		{
			"name":        "pippo",
			"version":     "",
			"hostname":    "",
			"environment": "",
			"charset":     "",

			"memory":       42.42,
			"datafileSize": 75.42,
			"segmentsSize": 99.42,
			"archivelog":   true,
			"ha":           false,
			"dataguard":    true,
		},
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	db.EXPECT().SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(oracleDbs, nil)

	mysqlInstances := []dto.MySQLInstance{
		{
			Hostname:    "pluto",
			Location:    "Cuba",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "mysql",
				Version:            "",
				Edition:            "",
				Platform:           "",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           1,
				ThreadsConcurrency: 2,
				BufferPoolSize:     43008,
				LogBufferSize:      4,
				SortBufferSize:     5,
				ReadOnly:           false,
				LogBin:             false,
				HighAvailability:   false,
				UUID:               "",
				IsMaster:           false,
				SlaveUUIDs:         []string{},
				IsSlave:            false,
				MasterUUID:         new(string),
				Databases:          []model.MySQLDatabase{{Name: "", Charset: "", Collation: "", Encrypted: false}},
				TableSchemas:       []model.MySQLTableSchema{{Name: "", Engine: "", Allocation: 24576}},
				SegmentAdvisors:    []model.MySQLSegmentAdvisor{{TableSchema: "", TableName: "", Engine: "", Allocation: 76, Data: 0, Index: 0, Free: 0}},
			},
		},
	}

	globalFilter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	actual, err := as.SearchDatabasesAsXLSX(globalFilter)
	require.NoError(t, err)

	expected := []dto.Database{
		{
			Name:         "pippo",
			Type:         "Oracle/Database",
			Version:      "",
			Hostname:     "",
			Environment:  "",
			Charset:      "",
			Memory:       42.42,
			DatafileSize: 75.42,
			SegmentsSize: 99.42,
		},
		{
			Name:         "mysql",
			Type:         "Oracle/MySQL",
			Version:      "",
			Hostname:     "pluto",
			Environment:  "TST",
			Charset:      "",
			Memory:       42.0,
			DatafileSize: 0,
			SegmentsSize: 24.0,
		},
	}

	assert.Equal(t, "Name", actual.GetCellValue("Databases", "A1"))
	assert.Equal(t, expected[0].Name, actual.GetCellValue("Databases", "A2"))
	assert.Equal(t, expected[1].Name, actual.GetCellValue("Databases", "A3"))

	assert.Equal(t, "Type", actual.GetCellValue("Databases", "B1"))
	assert.Equal(t, expected[0].Type, actual.GetCellValue("Databases", "B2"))
	assert.Equal(t, expected[1].Type, actual.GetCellValue("Databases", "B3"))

	assert.Equal(t, "Memory", actual.GetCellValue("Databases", "G1"))
	assert.Equal(t, "42.42", actual.GetCellValue("Databases", "G2"))
	assert.Equal(t, "42", actual.GetCellValue("Databases", "G3"))
}

func TestGetDatabasesStatistics_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	oracleDbs := []map[string]interface{}{
		{
			"name":        "pippo",
			"version":     "",
			"hostname":    "",
			"environment": "",
			"charset":     "",

			"memory":       42.42,
			"datafileSize": 75.42,
			"segmentsSize": 99.42,
			"archivelog":   true,
			"ha":           false,
			"dataguard":    true,
		},
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	db.EXPECT().SearchOracleDatabases(false, []string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(oracleDbs, nil)

	mysqlInstances := []dto.MySQLInstance{
		{
			Hostname:    "pluto",
			Location:    "Cuba",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "mysql",
				Version:            "",
				Edition:            "",
				Platform:           "",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           1,
				ThreadsConcurrency: 2,
				BufferPoolSize:     43008,
				LogBufferSize:      4,
				SortBufferSize:     5,
				ReadOnly:           false,
				Databases: []model.MySQLDatabase{
					{
						Name:      "",
						Charset:   "",
						Collation: "",
						Encrypted: false,
					},
				},
				TableSchemas: []model.MySQLTableSchema{
					{
						Name:       "",
						Engine:     "",
						Allocation: 24576,
					},
				},
				SegmentAdvisors: []model.MySQLSegmentAdvisor{
					{
						TableSchema: "",
						TableName:   "",
						Engine:      "",
						Allocation:  76,
						Data:        0,
						Index:       0,
						Free:        0,
					},
				},
			},
		},
	}

	globalFilter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	actual, err := as.GetDatabasesStatistics(globalFilter)
	require.NoError(t, err)

	expected := dto.DatabasesStatistics{
		TotalMemorySize:   84.42 * 1024 * 1024 * 1024,
		TotalSegmentsSize: 123.42 * 1024 * 1024 * 1024,
	}

	assert.Equal(t, expected, *actual)
}

func TestGetDatabasesUsedLicenses_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	filter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{{
			LicenseTypeID: "A12345",
			DbName:        "topolino-dbname",
			Hostname:      "topolino-hostname",
			UsedLicenses:  0,
		}},
	}

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "A12345",
			ItemDescription: "ThisDesc",
			Metric:          "ThisMetric",
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
	}
	usedLicenses := []dto.MySQLUsedLicense{
		{
			Hostname:        "pluto",
			InstanceName:    "pluto-instance",
			InstanceEdition: model.MySQLEditionEnterprise,
			AgreementType:   "",
		},
	}
	clusters := []dto.Cluster{
		{
			Hostname: "pluto-cluster",
			VMs: []dto.VM{
				{
					Hostname: "pluto",
				},
			},
		},
	}
	agreements := []model.MySQLAgreement{
		{
			ID:               [12]byte{},
			Type:             model.MySQLAgreementTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{"pippo-cluster", "pluto-cluster"},
			Hosts:            []string{},
		},
	}
	any := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses(filter).
			Return(usedLicenses, nil),
		db.EXPECT().GetClusters(any).
			Return(clusters, nil),
		db.EXPECT().GetMySQLAgreements().
			Return(agreements, nil),
	)

	actual, err := as.GetDatabasesUsedLicenses(filter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:      "topolino-hostname",
			DbName:        "topolino-dbname",
			LicenseTypeID: "A12345",
			Description:   "ThisDesc",
			Metric:        "ThisMetric",
			UsedLicenses:  0,
		},
		{
			Hostname:      "pluto",
			DbName:        "pluto-instance",
			LicenseTypeID: "",
			Description:   "MySQL ENTERPRISE",
			Metric:        "CLUSTER",
			UsedLicenses:  1,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestGetDatabasesUsedLicensesAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses := dto.DatabaseUsedLicense{
		Hostname:      "topolino-hostname",
		DbName:        "topolino-dbname",
		LicenseTypeID: "A12345",
		Description:   "ThisDesc",
		Metric:        "ThisMetric",
		UsedLicenses:  0,
	}

	as.mockGetDatabasesUsedLicenses = func(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
		return []dto.DatabaseUsedLicense{usedLicenses}, nil
	}

	actual, err := as.GetDatabasesUsedLicensesAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "topolino-hostname", actual.GetCellValue("Licenses Used", "A2"))
	assert.Equal(t, "topolino-dbname", actual.GetCellValue("Licenses Used", "B2"))
	assert.Equal(t, "A12345", actual.GetCellValue("Licenses Used", "C2"))
	assert.Equal(t, "ThisDesc", actual.GetCellValue("Licenses Used", "D2"))
	assert.Equal(t, "ThisMetric", actual.GetCellValue("Licenses Used", "E2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Used", "F2"))
}

func TestGetDatabaseLicensesComplianceAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	licenses := dto.LicenseCompliance{
			LicenseTypeID:   "L47247",
			ItemDescription: "Oracle Real Application Testing",
			Metric:          "Processor Perpetual",
			Consumed:        0,
			Covered:         0,
			Compliance:      1,
			Unlimited:       false,
	}
	as.mockGetDatabaseLicensesCompliance = func() ([]dto.LicenseCompliance, error) {
		return []dto.LicenseCompliance{licenses}, nil
	}

	actual, err := as.GetDatabaseLicensesComplianceAsXLSX()
	require.NoError(t, err)

	assert.Equal(t, "L47247", actual.GetCellValue("Licenses Compliance", "A2"))
	assert.Equal(t, "Oracle Real Application Testing", actual.GetCellValue("Licenses Compliance", "B2"))
	assert.Equal(t, "Processor Perpetual", actual.GetCellValue("Licenses Compliance", "C2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "D2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "E2"))
	assert.Equal(t, "1", actual.GetCellValue("Licenses Compliance", "F2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "G2"))
}

func TestGetDatabasesUsedLicensesPerHostAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	filter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   utils.MAX_TIME,
	}

	licenses := dto.DatabaseUsedLicensePerHost{
		Hostname:      "ercsoldbx",
		Databases:     2,
		LicenseTypeID: "A90611",
		Description:   "Oracle Database Enterprise Edition",
		Metric:        "Processor Perpetual",
		UsedLicenses:  2,
	}

	as.mockGetDatabasesUsedLicensesPerHost = func(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicensePerHost, error) {
		return []dto.DatabaseUsedLicensePerHost{licenses}, nil
	}

	actual, err := as.GetDatabasesUsedLicensesPerHostAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Licenses Used", "A2"))
	assert.Equal(t, "2", actual.GetCellValue("Licenses Used", "B2"))
	assert.Equal(t, "A90611", actual.GetCellValue("Licenses Used", "C2"))
	assert.Equal(t, "Oracle Database Enterprise Edition", actual.GetCellValue("Licenses Used", "D2"))
	assert.Equal(t, "Processor Perpetual", actual.GetCellValue("Licenses Used", "E2"))
	assert.Equal(t, "2", actual.GetCellValue("Licenses Used", "F2"))
}