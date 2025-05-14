// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

const (
	licenseTypesID = "359-06320"
)

var oracleDbs = []dto.OracleDatabase{
	{
		Name:                      "pippo",
		Version:                   "",
		Hostname:                  "",
		Environment:               "",
		Charset:                   "",
		Memory:                    42.42,
		DatafileSize:              75.42,
		SegmentsSize:              99.42,
		Archivelog:                true,
		Ha:                        false,
		Dataguard:                 true,
		PgsqlMigrabilitySemaphore: "red",
	},
}

var expectedRes = dto.OracleDatabaseResponse{
	Content: oracleDbs,
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

var thisMoment = utils.P("2019-11-05T14:02:03+01:00")

var globalFilter = dto.GlobalFilter{
	Location:    "Dubai",
	Environment: "TEST",
	OlderThan:   thisMoment,
}

var sqlServerInstances = []dto.SqlServerInstance{
	{
		Hostname:      "test-db",
		Environment:   "TST",
		Location:      "Germany",
		Name:          "MSSQLSERVER",
		Status:        "ONLINE",
		Edition:       "ENT",
		CollationName: "Latin1_General_CI_AS",
		Version:       "2019",
	},
}

var expectedSqlServerRes = dto.SqlServerInstanceResponse{
	Content: sqlServerInstances,
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

var postgreSqlInstances = []dto.PostgreSqlInstance{
	{
		Hostname:    "test-db",
		Environment: "TST",
		Location:    "Germany",
		Name:        "PostgreSQL-example:1010",
		Charset:     "UTF8",
		Version:     "PostgreSQL 10.20",
	},
}

var expectedPostgreSqlRes = dto.PostgreSqlInstanceResponse{
	Content: postgreSqlInstances,
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

var mongoDBInstances = []dto.MongoDBInstance{
	{
		Hostname:     "test-db",
		Environment:  "TST",
		Location:     "Germany",
		InstanceName: "host:27017",
		DBName:       "ercole",
		Charset:      "UTF8",
		Version:      "6.0.1",
	},
}

var expectedMongoDBRes = dto.MongoDBInstanceResponse{
	Content: mongoDBInstances,
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

var oracleLics = dto.OracleDatabaseUsedLicenseSearchResponse{
	Content: []dto.OracleDatabaseUsedLicense{
		{
			LicenseTypeID: "A12345",
			DbName:        "topolino-dbname",
			Hostname:      "topolino-hostname",
			UsedLicenses:  2,
		},
		{
			LicenseTypeID: "A98765",
			DbName:        "topolino-dbname",
			Hostname:      "topolino-hostname",
			UsedLicenses:  2,
		},
	},
}

var objID, _ = primitive.ObjectIDFromHex("609ce4782eff5d5540ec4a30")

var licenseTypes = []model.OracleDatabaseLicenseType{
	{
		ID:              "A12345",
		ItemDescription: "ThisDesc",
		Metric:          "ThisMetric",
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
	{
		ID:              "A98765",
		ItemDescription: "ThisDesc",
		Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
}

var hostdatas = []model.HostDataBE{
	{
		ClusterMembershipStatus: model.ClusterMembershipStatus{
			OracleClusterware:       false,
			SunCluster:              false,
			HACMP:                   false,
			VeritasClusterServer:    false,
			VeritasClusterHostnames: []string{},
		},
	},
}

var licenseTypesRac = []model.OracleDatabaseLicenseType{
	{
		ID:              "A12345",
		ItemDescription: "ThisDesc",
		Metric:          "ThisMetric",
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
	{
		ID:              "A98765",
		ItemDescription: "ThisDesc",
		Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
	{
		ID:              racIds[1],
		ItemDescription: "rac",
		Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
	{
		ID:              racOneNodeIds[0],
		ItemDescription: "rac one node",
		Metric:          model.LicenseTypeMetricProcessorPerpetual,
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
}

var globalFilterAny = dto.GlobalFilter{
	Location:    "",
	Environment: "",
	OlderThan:   utils.MAX_TIME,
}

var sqlServerLicenseTypes = []model.SqlServerDatabaseLicenseType{
	{
		ID:              licenseTypesID,
		ItemDescription: "SQL Server Enterprise Edition",
		Edition:         "ENT",
		Version:         "2019",
	},
}

var sqlServerLics = dto.SqlServerDatabaseUsedLicenseSearchResponse{
	Content: []dto.SqlServerDatabaseUsedLicense{
		{
			LicenseTypeID: licenseTypesID,
			DbName:        "topolino-dbname",
			Hostname:      "topolino-hostname",
			UsedLicenses:  8,
		},
	},
}

var hostdatasVm1 = []model.HostDataBE{
	{
		Hostname: "vm1",
		ClusterMembershipStatus: model.ClusterMembershipStatus{
			OracleClusterware:       false,
			SunCluster:              false,
			HACMP:                   false,
			VeritasClusterServer:    false,
			VeritasClusterHostnames: []string{},
		},
	},
}

var clusters = []dto.Cluster{
	{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "cluster1",
		HostnameAgentVirtualization: "",
		Name:                        "name1",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         12,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "vm1",
				Name:               "",
				VirtualizationNode: "",
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	},
}

var licenseTypes1 = []model.OracleDatabaseLicenseType{
	{
		ID:              "id1",
		ItemDescription: "desc1",
		Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
}

func TestSearchDatabases_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	// db.EXPECT().SearchOracleDatabaseUsedLicenses(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
	// 	Return(&dto.OracleDatabaseUsedLicenseSearchResponse{}, nil)

	// db.EXPECT().GetOracleDatabaseLicenseTypes().Return([]model.OracleDatabaseLicenseType{}, nil)

	// db.EXPECT().GetHostDatas(gomock.Any()).Return([]model.HostDataBE{}, nil)

	// db.EXPECT().GetClusters(gomock.Any()).Return([]dto.Cluster{}, nil)

	// db.EXPECT().FindPsqlMigrabilities(gomock.Any(), gomock.Any()).Return([]model.PgsqlMigrability{}, nil).AnyTimes()

	db.EXPECT().SearchOracleDatabases([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedRes, nil)

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

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	db.EXPECT().SearchSqlServerInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedSqlServerRes, nil)
	db.EXPECT().SearchPostgreSqlInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedPostgreSqlRes, nil)

	db.EXPECT().SearchMongoDBInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedMongoDBRes, nil)

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
			Location:         "Cuba",
			Charset:          "",
			Memory:           42.0,
			DatafileSize:     0,
			SegmentsSize:     24.0,
			Archivelog:       true,
			HighAvailability: false,
			DisasterRecovery: true,
		},
		{
			Name:        "MSSQLSERVER",
			Type:        "Microsoft/SQLServer",
			Version:     "2019",
			Hostname:    "test-db",
			Environment: "TST",
			Location:    "Germany",
			Charset:     "Latin1_General_CI_AS",
		},
		{
			Name:        "PostgreSQL-example:1010",
			Type:        "PostgreSQL/PostgreSQL",
			Version:     "PostgreSQL 10.20",
			Hostname:    "test-db",
			Environment: "TST",
			Location:    "Germany",
			Charset:     "UTF8",
		},
		{
			Hostname:    "test-db",
			Type:        "MongoDB/MongoDB",
			Environment: "TST",
			Location:    "Germany",
			Name:        "host:27017",
			Charset:     "UTF8",
			Version:     "6.0.1",
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

	// db.EXPECT().SearchOracleDatabaseUsedLicenses(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
	// 	Return(&dto.OracleDatabaseUsedLicenseSearchResponse{}, nil)

	// db.EXPECT().GetOracleDatabaseLicenseTypes().Return([]model.OracleDatabaseLicenseType{}, nil)

	// db.EXPECT().GetHostDatas(gomock.Any()).Return([]model.HostDataBE{}, nil)

	// db.EXPECT().GetClusters(gomock.Any()).Return([]dto.Cluster{}, nil)

	db.EXPECT().SearchOracleDatabases([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedRes, nil)

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

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	db.EXPECT().FindPsqlMigrabilities(gomock.Any(), gomock.Any()).Return([]model.PgsqlMigrability{}, nil).AnyTimes()

	db.EXPECT().SearchSqlServerInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedSqlServerRes, nil)

	db.EXPECT().SearchPostgreSqlInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedPostgreSqlRes, nil)

	db.EXPECT().SearchMongoDBInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedMongoDBRes, nil)

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
			Location:     "Cuba",
			Charset:      "",
			Memory:       42.0,
			DatafileSize: 0,
			SegmentsSize: 24.0,
		},
		{
			Name:        "MSSQLSERVER",
			Type:        "Microsoft/SQLServer",
			Version:     "2019",
			Hostname:    "test-db",
			Environment: "TST",
			Location:    "Germany",
			Charset:     "Latin1_General_CI_AS",
		},
		{
			Name:        "PostgreSQL-example:1010",
			Type:        "PostgreSQL/PostgreSQL",
			Version:     "PostgreSQL 10.20",
			Hostname:    "test-db",
			Environment: "TST",
			Location:    "Germany",
			Charset:     "UTF8",
		},
		{
			Hostname:    "test-db",
			Type:        "MongoDB/MongoDB",
			Environment: "TST",
			Location:    "Germany",
			Name:        "test",
			Charset:     "UTF8",
			Version:     "6.0.1",
		},
	}

	assert.Equal(t, "Name", actual.GetCellValue("Databases", "A1"))
	assert.Equal(t, expected[0].Name, actual.GetCellValue("Databases", "A2"))
	assert.Equal(t, expected[1].Name, actual.GetCellValue("Databases", "A3"))
	assert.Equal(t, expected[2].Name, actual.GetCellValue("Databases", "A4"))

	assert.Equal(t, "Type", actual.GetCellValue("Databases", "B1"))
	assert.Equal(t, expected[0].Type, actual.GetCellValue("Databases", "B2"))
	assert.Equal(t, expected[1].Type, actual.GetCellValue("Databases", "B3"))
	assert.Equal(t, expected[2].Type, actual.GetCellValue("Databases", "B4"))

	assert.Equal(t, "Memory", actual.GetCellValue("Databases", "H1"))
	assert.Equal(t, "42.42", actual.GetCellValue("Databases", "H2"))
	assert.Equal(t, "42", actual.GetCellValue("Databases", "H3"))
}

func TestGetDatabasesStatistics_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	// db.EXPECT().SearchOracleDatabaseUsedLicenses(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
	// 	Return(&dto.OracleDatabaseUsedLicenseSearchResponse{}, nil)

	// db.EXPECT().GetOracleDatabaseLicenseTypes().Return([]model.OracleDatabaseLicenseType{}, nil)

	// db.EXPECT().GetHostDatas(gomock.Any()).Return([]model.HostDataBE{}, nil)

	// db.EXPECT().GetClusters(gomock.Any()).Return([]dto.Cluster{}, nil)

	// db.EXPECT().FindPsqlMigrabilities(gomock.Any(), gomock.Any()).Return([]model.PgsqlMigrability{}, nil).AnyTimes()

	db.EXPECT().SearchOracleDatabases([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedRes, nil)

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

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(mysqlInstances, nil)

	db.EXPECT().SearchSqlServerInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedSqlServerRes, nil)

	db.EXPECT().SearchPostgreSqlInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedPostgreSqlRes, nil)

	db.EXPECT().SearchMongoDBInstances([]string{""}, "", false, -1, -1, "Dubai", "TEST", thisMoment).
		Return(&expectedMongoDBRes, nil)

	actual, err := as.GetDatabasesStatistics(globalFilter)
	require.NoError(t, err)

	expected := dto.DatabasesStatistics{
		TotalMemorySize:   84.42 * 1024 * 1024 * 1024,
		TotalSegmentsSize: 123.42 * 1024 * 1024 * 1024,
	}

	assert.Equal(t, expected, *actual)
}

func TestGetUsedLicensesPerDatabases_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	usedLicensesMySQL := []dto.MySQLUsedLicense{
		{
			LicenseTypeID:   model.MySqlPartNumber,
			Hostname:        "pluto",
			InstanceName:    "pluto-instance",
			InstanceEdition: model.MySqlItemDescription,
			Clustername:     "plutocluster",
			UsedLicenses:    1,
			Ignored:         false,
			IgnoredComment:  "",
			ContractType:    model.MySQLContractTypeCluster,
		},
	}
	clusters := []dto.Cluster{
		{
			Name:     "plutocluster",
			Hostname: "pluto-cluster",
			VMs: []dto.VM{
				{
					Hostname: "pluto",
				},
			},
		},
	}
	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "pluto-cluster",
		HostnameAgentVirtualization: "",
		Name:                        "plutocluster",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "pluto",
				Name:               "",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}
	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			ContractID:       "",
			CSI:              "",
			LicenseTypeID:    model.MySqlPartNumber,
			NumberOfLicenses: 12,
			Clusters:         []string{"plutocluster"},
			Hosts:            []string{},
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: "ABCDEF-GH",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  8,
			},
		},
	}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "ABCDEF-GH",
			ItemDescription: "test sql server enterprise",
			Edition:         "ENT",
			Version:         "2019",
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLUsedLicenses("", globalFilter).
			Return(usedLicensesMySQL, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().GetCluster("plutocluster", utils.MAX_TIME).
			Return(&cluster, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
	)
	actual, err := as.GetUsedLicensesPerDatabases("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 0,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 0,
			Ignored:         false,
		},
		{
			Hostname:        "pluto",
			DbName:          "pluto-instance",
			LicenseTypeID:   "B64911",
			Description:     model.MySqlItemDescription,
			Metric:          model.MySQLContractTypeHost,
			UsedLicenses:    1,
			ClusterLicenses: 0,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "ABCDEF-GH",
			Description:     "test sql server enterprise",
			Metric:          "HOST",
			UsedLicenses:    8,
			ClusterLicenses: 0,
			Ignored:         false,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestGetUsedLicensesPerDatabases_VMWareCluster_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	usedLicensesMySQL := []dto.MySQLUsedLicense{}
	clusters := []dto.Cluster{
		{
			Hostname: "topolino-cluster",
			CPU:      16,
			VMs: []dto.VM{
				{
					Hostname: "topolino-hostname",
				},
			},
		},
	}
	contracts := []model.MySQLContract{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "topolino-hostname",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "ABCDEF-GH",
			ItemDescription: "test sql server enterprise",
			Edition:         "ENT",
			Version:         "2019",
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLUsedLicenses("", globalFilter).
			Return(usedLicensesMySQL, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	actual, err := as.GetUsedLicensesPerDatabases("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 8,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 200,
			Ignored:         false,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestGetUsedLicensesPerDatabases_VeritasCluster_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	usedLicensesMySQL := []dto.MySQLUsedLicense{}
	clusters := []dto.Cluster{}
	contracts := []model.MySQLContract{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "topolino-hostname",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    true,
				VeritasClusterHostnames: []string{"topolino-hostname", "qui", "quo", "qua"},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "ABCDEF-GH",
			ItemDescription: "test sql server enterprise",
			Edition:         "ENT",
			Version:         "2019",
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLUsedLicenses("", globalFilter).
			Return(usedLicensesMySQL, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	actual, err := as.GetUsedLicensesPerDatabases("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 84,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 2100,
			Ignored:         false,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestGetUsedLicensesPerDatabasesAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
		Log:      logger.NewLogger("TEST"),
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
			ContractType:    "",
		},
	}
	clusters := []dto.Cluster{
		{
			Name:     "plutocluster",
			Hostname: "pluto-cluster",
			VMs: []dto.VM{
				{
					Hostname: "pluto",
				},
			},
		},
	}
	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "pluto-cluster",
		HostnameAgentVirtualization: "",
		Name:                        "plutocluster",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "pluto",
				Name:               "",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}
	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{"plutocluster"},
			Hosts:            []string{},
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: "ABCDEF-GH",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  8,
			},
		},
	}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "ABCDEF-GH",
			ItemDescription: "test sql server enterprise",
			Edition:         "ENT",
			Version:         "2019",
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLUsedLicenses("", globalFilter).
			Return(usedLicenses, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().GetCluster("plutocluster", utils.MAX_TIME).
			Return(&cluster, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	actual, err := as.GetUsedLicensesPerDatabasesAsXLSX(globalFilter)
	require.NoError(t, err)

	assert.Equal(t, "topolino-hostname", actual.GetCellValue("Licenses Used", "A2"))
	assert.Equal(t, "topolino-dbname", actual.GetCellValue("Licenses Used", "B2"))
	assert.Equal(t, "A12345", actual.GetCellValue("Licenses Used", "C2"))
	assert.Equal(t, "ThisDesc", actual.GetCellValue("Licenses Used", "D2"))
	assert.Equal(t, "ThisMetric", actual.GetCellValue("Licenses Used", "E2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Used", "F2"))
}

func TestGetOracleDatabasesUsedLicenses_Host_WithActiveDataguardAndGoldenGate_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "A12345",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: "A98765",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: goldenGateIds[0],
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: activeDataguardIds[1],
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
		},
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
		{
			ID:              "A98765",
			ItemDescription: "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
		{
			ID:              goldenGateIds[0],
			ItemDescription: "golden gate",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
		{
			ID:              activeDataguardIds[1],
			ItemDescription: "active dataguard",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
	}

	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "topolino-hostname",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}

	host := dto.HostData{
		ID:                      [12]byte{},
		Archived:                false,
		CreatedAt:               time.Time{},
		ServerVersion:           "",
		SchemaVersion:           0,
		ServerSchemaVersion:     0,
		Hostname:                "topolino-hostname",
		Location:                "",
		Environment:             "",
		AgentVersion:            "",
		Cluster:                 "",
		VirtualizationNode:      "",
		Tags:                    []string{},
		Info:                    model.Host{},
		ClusterMembershipStatus: model.ClusterMembershipStatus{},
		Features:                model.Features{},
		Filesystems:             []model.Filesystem{},
		Clusters:                []model.ClusterInfo{},
		Cloud:                   model.Cloud{},
		Errors:                  []model.AgentError{},
		OtherInfo:               map[string]interface{}{},
		Alerts:                  []model.Alert{},
		History:                 []model.History{},
	}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),

		db.EXPECT().GetHost("topolino-hostname", utils.MAX_TIME, false).
			Return(&host, nil).AnyTimes(),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	actual, err := as.getOracleDatabasesUsedLicenses("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 0,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 0,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   goldenGateIds[0],
			Description:     "golden gate",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 0,
			Ignored:         false,
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestGetOracleDatabasesUsedLicenses_VeritasCluster_WithActiveDataguardAndGoldenGate_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "A12345",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: "A98765",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: goldenGateIds[0],
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: activeDataguardIds[1],
				DbName:        "topolino-dbname",
				Hostname:      "qui",
				UsedLicenses:  2,
			},
		},
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
		{
			ID:              "A98765",
			ItemDescription: "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
		{
			ID:              goldenGateIds[0],
			ItemDescription: "golden gate",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
		{
			ID:              activeDataguardIds[1],
			ItemDescription: "active dataguard",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
	}

	host := dto.HostData{
		ID:                      [12]byte{},
		Archived:                false,
		CreatedAt:               time.Time{},
		ServerVersion:           "",
		SchemaVersion:           0,
		ServerSchemaVersion:     0,
		Hostname:                "topolino-hostname",
		Location:                "",
		Environment:             "",
		AgentVersion:            "",
		Cluster:                 "",
		VirtualizationNode:      "",
		Tags:                    []string{},
		Info:                    model.Host{},
		ClusterMembershipStatus: model.ClusterMembershipStatus{},
		Features:                model.Features{},
		Filesystems:             []model.Filesystem{},
		Clusters:                []model.ClusterInfo{},
		Cloud:                   model.Cloud{},
		Errors:                  []model.AgentError{},
		OtherInfo:               map[string]interface{}{},
		Alerts:                  []model.Alert{},
		History:                 []model.History{},
	}

	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "topolino-hostname",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    true,
				VeritasClusterHostnames: []string{"topolino-hostname", "qui", "quo", "qua"},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),

		db.EXPECT().GetHost("topolino-hostname", utils.MAX_TIME, false).
			Return(&host, nil).AnyTimes(),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	actual, err := as.getOracleDatabasesUsedLicenses("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 84,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 2100,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   goldenGateIds[0],
			Description:     "golden gate",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 2100,
			Ignored:         false,
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestGetOracleDatabasesUsedLicenses_Host_WithRacAndRacOneNode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "A12345",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: "A98765",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: racIds[1],
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: racOneNodeIds[0],
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
		},
	}

	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "topolino-hostname",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}

	host := dto.HostData{
		ID:                      [12]byte{},
		Archived:                false,
		CreatedAt:               time.Time{},
		ServerVersion:           "",
		SchemaVersion:           0,
		ServerSchemaVersion:     0,
		Hostname:                "topolino-hostname",
		Location:                "",
		Environment:             "",
		AgentVersion:            "",
		Cluster:                 "",
		VirtualizationNode:      "",
		Tags:                    []string{},
		Info:                    model.Host{},
		ClusterMembershipStatus: model.ClusterMembershipStatus{},
		Features:                model.Features{},
		Filesystems:             []model.Filesystem{},
		Clusters:                []model.ClusterInfo{},
		Cloud:                   model.Cloud{},
		Errors:                  []model.AgentError{},
		OtherInfo:               map[string]interface{}{},
		Alerts:                  []model.Alert{},
		History:                 []model.History{},
	}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypesRac, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),

		db.EXPECT().GetHost("topolino-hostname", utils.MAX_TIME, false).
			Return(&host, nil).AnyTimes(),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	actual, err := as.getOracleDatabasesUsedLicenses("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 0,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 0,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   racIds[1],
			Description:     "rac",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 0,
			Ignored:         false,
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestGetOracleDatabasesUsedLicenses_VeritasCluster_WithRacAndRacOneNode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "A12345",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: "A98765",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: racIds[1],
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: racOneNodeIds[0],
				DbName:        "topolino-dbname",
				Hostname:      "qui",
				UsedLicenses:  2,
			},
		},
	}

	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "topolino-hostname",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    true,
				VeritasClusterHostnames: []string{"topolino-hostname", "qui", "quo", "qua"},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}

	host := dto.HostData{
		ID:                      [12]byte{},
		Archived:                false,
		CreatedAt:               time.Time{},
		ServerVersion:           "",
		SchemaVersion:           0,
		ServerSchemaVersion:     0,
		Hostname:                "topolino-hostname",
		Location:                "",
		Environment:             "",
		AgentVersion:            "",
		Cluster:                 "",
		VirtualizationNode:      "",
		Tags:                    []string{},
		Info:                    model.Host{},
		ClusterMembershipStatus: model.ClusterMembershipStatus{},
		Features:                model.Features{},
		Filesystems:             []model.Filesystem{},
		Clusters:                []model.ClusterInfo{},
		Cloud:                   model.Cloud{},
		Errors:                  []model.AgentError{},
		OtherInfo:               map[string]interface{}{},
		Alerts:                  []model.Alert{},
		History:                 []model.History{},
	}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypesRac, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),

		db.EXPECT().GetHost("topolino-hostname", utils.MAX_TIME, false).
			Return(&host, nil).AnyTimes(),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	actual, err := as.getOracleDatabasesUsedLicenses("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 84,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 2100,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			ClusterName:     "qua,qui,quo,topolino-hostname",
			ClusterType:     "VeritasCluster",
			LicenseTypeID:   racIds[1],
			Description:     "rac",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 2100,
			Ignored:         false,
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestGetOracleDatabasesUsedLicenses_VmwareCluster_WithRacAndRacOneNode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "A12345",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: "A98765",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: racIds[1],
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
			},
			{
				LicenseTypeID: racOneNodeIds[0],
				DbName:        "qui-dbname",
				Hostname:      "qui",
				UsedLicenses:  2,
			},
		},
	}

	clusters := []dto.Cluster{
		{
			ID:                          [12]byte{},
			CreatedAt:                   time.Time{},
			Hostname:                    "cluster1",
			HostnameAgentVirtualization: "",
			Name:                        "",
			Environment:                 "",
			Location:                    "",
			FetchEndpoint:               "",
			CPU:                         64,
			Sockets:                     0,
			Type:                        "",
			VirtualizationNodes:         []string{},
			VirtualizationNodesCount:    0,
			VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
			VMs: []dto.VM{
				{
					CappedCPU:          false,
					Hostname:           "qui",
					Name:               "",
					VirtualizationNode: "",
				},
				{
					CappedCPU:          false,
					Hostname:           "topolino-hostname",
					Name:               "",
					VirtualizationNode: "",
				},
			},
			VMsCount:            1,
			VMsErcoleAgentCount: 1,
		},
	}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "topolino-hostname",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
		{
			Hostname: "qui",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 12,
			},
		},
	}

	host := dto.HostData{
		ID:                      [12]byte{},
		Archived:                false,
		CreatedAt:               time.Time{},
		ServerVersion:           "",
		SchemaVersion:           0,
		ServerSchemaVersion:     0,
		Hostname:                "topolino-hostname",
		Location:                "",
		Environment:             "",
		AgentVersion:            "",
		Cluster:                 "",
		VirtualizationNode:      "",
		Tags:                    []string{},
		Info:                    model.Host{},
		ClusterMembershipStatus: model.ClusterMembershipStatus{},
		Features:                model.Features{},
		Filesystems:             []model.Filesystem{},
		Clusters:                []model.ClusterInfo{},
		Cloud:                   model.Cloud{},
		Errors:                  []model.AgentError{},
		OtherInfo:               map[string]interface{}{},
		Alerts:                  []model.Alert{},
		History:                 []model.History{},
	}

	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, globalFilter.Location, globalFilter.Environment, globalFilter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypesRac, nil),

		db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),

		db.EXPECT().GetHost("topolino-hostname", utils.MAX_TIME, false).
			Return(&host, nil).AnyTimes(),
	)

	db.EXPECT().ExistHostdata("topolino-hostname").Return(true, nil).AnyTimes()

	db.EXPECT().ExistHostdata("qui").Return(true, nil).AnyTimes()

	actual, err := as.getOracleDatabasesUsedLicenses("", globalFilter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicense{
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A12345",
			Description:     "ThisDesc",
			Metric:          "ThisMetric",
			UsedLicenses:    2,
			ClusterLicenses: 32,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   "A98765",
			Description:     "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 800,
			Ignored:         false,
		},
		{
			Hostname:        "topolino-hostname",
			DbName:          "topolino-dbname",
			LicenseTypeID:   racIds[1],
			Description:     "rac",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:    50,
			ClusterLicenses: 800,
			Ignored:         false,
		},
	}

	assert.ElementsMatch(t, expected, actual)
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

	oracleContracts := []dto.OracleDatabaseContractFE{
		{
			ID:              objID,
			ContractID:      "5051863",
			CSI:             "18000742",
			LicenseTypeID:   "L47225",
			ItemDescription: "Oracle Advanced Compression",
			Metric:          "Named User Plus Perpetual",
			ReferenceNumber: "66880702",
			Unlimited:       false,
			Basket:          false,
			Restricted:      false,
			Hosts: []dto.OracleDatabaseContractAssociatedHostFE{
				{
					Hostname:                  "sdlsts101",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
					ConsumedLicensesCount:     0,
				},
			},
			LicensesPerCore:          0,
			LicensesPerUser:          150,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 150,
		},
	}

	oracleLicenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "L47225",
			ItemDescription: "Oracle Advanced Compression",
			Metric:          "Named User Plus Perpetual",
			Cost:            230,
			Aliases:         []string{"Advanced Compression"},
			Option:          true,
		},
	}

	usedLicenses := []dto.MySQLUsedLicense{
		{
			Hostname:        "pluto",
			InstanceName:    "pluto-instance",
			InstanceEdition: model.MySQLEditionEnterprise,
			ContractType:    "",
		},
	}
	clusters := []dto.Cluster{
		{
			Name:     "plutocluster",
			Hostname: "pluto-cluster",
			VMs: []dto.VM{
				{
					Hostname: "pluto",
				},
			},
		},
	}
	cluster := dto.Cluster{
		ID:                          objID,
		CreatedAt:                   time.Time{},
		Hostname:                    "pluto-cluster",
		HostnameAgentVirtualization: "",
		Name:                        "plutocluster",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "pluto",
				Name:               "",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}
	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{"plutocluster"},
			Hosts:            []string{},
		},
	}

	searchResponse := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "L47225",
				DbName:        "pippo",
				Hostname:      "sdlsts101",
				UsedLicenses:  6,
				Ignored:       false,
			},
		},
		Metadata: dto.PagingMetadata{},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeCluster,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  licenseTypesID,
			Clusters:       []string{"plutocluster"},
			Hosts:          []string{},
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return([]model.HostDataBE{{
			Hostname: "sdlsts101",
		}}, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(oracleContracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&searchResponse, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().GetCluster("plutocluster", utils.MAX_TIME).
			Return(&cluster, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetCluster("plutocluster", utils.MAX_TIME).
			Return(&cluster, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Times(1).
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("pluto").Return(true, nil).AnyTimes()

	actual, err := as.GetDatabaseLicensesComplianceAsXLSX([]string{})
	require.NoError(t, err)

	assert.Equal(t, "L47225", actual.GetCellValue("Licenses Compliance", "A2"))
	assert.Equal(t, "Oracle Advanced Compression", actual.GetCellValue("Licenses Compliance", "B2"))
	assert.Equal(t, "Named User Plus Perpetual", actual.GetCellValue("Licenses Compliance", "C2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "D2"))
	assert.Equal(t, "150", actual.GetCellValue("Licenses Compliance", "E2"))
	assert.Equal(t, "150", actual.GetCellValue("Licenses Compliance", "F2"))
	assert.Equal(t, "150", actual.GetCellValue("Licenses Compliance", "G2"))
	assert.Equal(t, "1", actual.GetCellValue("Licenses Compliance", "H2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "I2"))
	assert.Equal(t, "", actual.GetCellValue("Licenses Compliance", "J2"))
}

func TestGetDatabaseLicensesCompliance_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	oracleContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       objID,
			ContractID:               "5051863",
			CSI:                      "19338486",
			LicenseTypeID:            "L10006",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   "Named User Plus Perpetual",
			ReferenceNumber:          "96661555",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          450,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 450,
			CoveredLicenses:          0,
		},
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "TEST",
			CSI:                      "99999999",
			LicenseTypeID:            "M10080",
			ItemDescription:          "Oracle Database Enterprise Edition",
			Metric:                   "Processor Perpetual",
			ReferenceNumber:          "666666666",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          50,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 50,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          0,
		},
	}

	searchResponse := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "L10006",
				DbName:        "pippo",
				Hostname:      "homer",
				UsedLicenses:  10,
				Ignored:       false,
			},
			{
				LicenseTypeID: "M10080",
				DbName:        "pluto",
				Hostname:      "homer",
				UsedLicenses:  20,
				Ignored:       false,
			},
		},
		Metadata: dto.PagingMetadata{},
	}

	var oracleLicenseTypes = []model.OracleDatabaseLicenseType{
		{
			ID:              "L10006",
			ItemDescription: "Oracle Partitioning",
			Metric:          "Named User Plus Perpetual",
			Cost:            250,
			Aliases:         []string{"Partitioning"},
			Option:          true,
		},
		{
			ID:              "M10080",
			ItemDescription: "Application Testing",
			Metric:          "Processor Perpetual",
			Cost:            230,
			Aliases:         []string{"Application Testing"},
			Option:          true,
		},
	}
	usedLicenses := []dto.MySQLUsedLicense{
		{
			Hostname:        "pluto",
			InstanceName:    "pluto-instance",
			InstanceEdition: model.MySqlItemDescription,
			ContractType:    model.MySQLContractTypeHost,
			LicenseTypeID:   model.MySqlPartNumber,
			UsedLicenses:    1,
		},
	}
	clusters := []dto.Cluster{
		{
			Name:     "plutocluster",
			Hostname: "pluto-cluster",
			VMs: []dto.VM{
				{
					Hostname: "pluto",
				},
			},
		},
	}
	cluster := dto.Cluster{
		ID:                          objID,
		CreatedAt:                   time.Time{},
		Hostname:                    "pluto-cluster",
		HostnameAgentVirtualization: "",
		Name:                        "plutocluster",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "pluto",
				Name:               "",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}
	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeHost,
			NumberOfLicenses: 12,
			Clusters:         []string{"plutocluster"},
			Hosts:            []string{},
			LicenseTypeID:    model.MySqlPartNumber,
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeCluster,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  licenseTypesID,
			Clusters:       []string{"plutocluster"},
			Hosts:          []string{},
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return([]model.HostDataBE{{
			Hostname: "homer",
		}}, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(oracleContracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&searchResponse, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil).AnyTimes(),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetCluster("plutocluster", utils.MAX_TIME).
			Return(&cluster, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Times(1).
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("pluto").Return(true, nil).AnyTimes()

	actual, err := as.GetDatabaseLicensesCompliance([]string{})
	require.NoError(t, err)

	expected := []dto.LicenseCompliance{
		{
			LicenseTypeID:   "M10080",
			ItemDescription: "Application Testing",
			Metric:          "Processor Perpetual",
			Cost:            230,
			Consumed:        20,
			Covered:         20,
			Purchased:       50,
			Compliance:      1,
			Unlimited:       false,
			Available:       30,
		},
		{
			LicenseTypeID:   "L10006",
			ItemDescription: "Oracle Partitioning",
			Metric:          "Named User Plus Perpetual",
			Cost:            250,
			Consumed:        250,
			Covered:         250,
			Purchased:       450,
			Compliance:      1,
			Unlimited:       false,
			Available:       200,
		},
		{
			LicenseTypeID:   model.MySqlPartNumber,
			ItemDescription: model.MySqlItemDescription,
			Metric:          model.MySQLContractTypeHost,
			Cost:            0,
			Consumed:        1,
			Covered:         1,
			Purchased:       12,
			Compliance:      1,
			Unlimited:       false,
			Available:       11,
		},
		{
			LicenseTypeID:   licenseTypesID,
			ItemDescription: "SQL Server Enterprise Edition",
			Metric:          "HOST",
			Cost:            0,
			Consumed:        8,
			Covered:         0,
			Purchased:       0,
			Compliance:      0,
			Unlimited:       false,
			Available:       0,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetUsedLicensesPerHost_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   utils.MAX_TIME,
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{{
			LicenseTypeID: "A90611",
			DbName:        "ercsoldbx",
			Hostname:      "ercsoldbx",
			UsedLicenses:  2,
		}},
	}
	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "A90611",
			ItemDescription: "Oracle Database Enterprise Edition",
			Metric:          "Processor Perpetual",
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
	}
	usedLicenses := []dto.MySQLUsedLicense{
		{
			Hostname:        "pluto",
			InstanceName:    "pluto-instance",
			InstanceEdition: model.MySqlItemDescription,
			ContractType:    model.MySQLContractTypeHost,
			UsedLicenses:    1,
			LicenseTypeID:   model.MySqlPartNumber,
		},
	}
	clusters := []dto.Cluster{
		{
			Hostname: "pluto",
			VMs: []dto.VM{
				{
					Hostname: "pluto",
				},
			},
			Name: "PLUTO-CLUSTER-NAME",
			CPU:  45,
		},
	}
	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "plutone",
		HostnameAgentVirtualization: "",
		Name:                        "PLUTO-CLUSTER-NAME",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         45,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "pluto",
				Name:               "",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}
	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{"PLUTO-CLUSTER-NAME"},
			Hosts:            []string{},
		},
	}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "pluto",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: "ABCDEF-GH",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  8,
			},
		},
	}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "ABCDEF-GH",
			ItemDescription: "test sql server enterprise",
			Edition:         "ENT",
			Version:         "2019",
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{
		{
			ID:             [12]byte{},
			Type:           model.MySQLContractTypeHost,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  "ABCDEF-GH",
			Clusters:       []string{"PLUTO-CLUSTER-NAME"},
			Hosts:          []string{},
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", filter).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().GetCluster("PLUTO-CLUSTER-NAME", utils.MAX_TIME).
			Return(&cluster, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("pluto").Return(true, nil).AnyTimes()

	actual, err := as.GetUsedLicensesPerHost(filter)
	require.NoError(t, err)

	expected := []dto.DatabaseUsedLicensePerHost{
		{
			Hostname:        "ercsoldbx",
			DatabaseNames:   []string{"ercsoldbx"},
			LicenseTypeID:   "A90611",
			Description:     "Oracle Database Enterprise Edition",
			Metric:          "Processor Perpetual",
			UsedLicenses:    2,
			ClusterLicenses: 0,
		},
		{
			Hostname:        "pluto",
			DatabaseNames:   []string{"pluto-instance"},
			LicenseTypeID:   model.MySqlPartNumber,
			Description:     model.MySqlItemDescription,
			Metric:          model.MySQLContractTypeHost,
			UsedLicenses:    1,
			ClusterLicenses: 0,
		},
		{
			Hostname:        "topolino-hostname",
			DatabaseNames:   []string{"topolino-dbname"},
			LicenseTypeID:   "ABCDEF-GH",
			Description:     "test sql server enterprise",
			Metric:          "HOST",
			UsedLicenses:    8,
			ClusterLicenses: 0,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetUsedLicensesPerHostAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   utils.MAX_TIME,
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{{
			LicenseTypeID: "A90611",
			DbName:        "ercsoldbx",
			Hostname:      "ercsoldbx",
			UsedLicenses:  2,
		}},
	}
	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "A90611",
			ItemDescription: "Oracle Database Enterprise Edition",
			Metric:          "Processor Perpetual",
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
			ContractType:    "",
		},
	}
	clusters := []dto.Cluster{
		{
			Name:     "plutocluster",
			Hostname: "pluto-cluster",
			VMs: []dto.VM{
				{
					Hostname: "pluto",
				},
			},
		},
	}
	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "pluto-cluster",
		HostnameAgentVirtualization: "",
		Name:                        "plutocluster",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "pluto",
				Name:               "",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}
	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{"plutocluster"},
			Hosts:            []string{},
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: "ABCDEF-GH",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  8,
			},
		},
	}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "ABCDEF-GH",
			ItemDescription: "test sql server enterprise",
			Edition:         "ENT",
			Version:         "2019",
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeCluster,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  "ABCDEF-GH",
			Clusters:       []string{"plutocluster"},
			Hosts:          []string{},
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", filter).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().GetCluster("plutocluster", utils.MAX_TIME).
			Return(&cluster, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetCluster("plutocluster", utils.MAX_TIME).
			Return(&cluster, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("pluto").Return(true, nil).AnyTimes()

	actual, err := as.GetUsedLicensesPerHostAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Licenses Used Per Host", "A2"))
	assert.Equal(t, "1", actual.GetCellValue("Licenses Used Per Host", "B2"))
	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Licenses Used Per Host", "C2"))
	assert.Equal(t, "A90611", actual.GetCellValue("Licenses Used Per Host", "D2"))
	assert.Equal(t, "Oracle Database Enterprise Edition", actual.GetCellValue("Licenses Used Per Host", "E2"))
	assert.Equal(t, "Processor Perpetual", actual.GetCellValue("Licenses Used Per Host", "F2"))
	assert.Equal(t, "2", actual.GetCellValue("Licenses Used Per Host", "G2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Used Per Host", "H2"))
}

func TestGetUsedLicensesPerCluster_OneVm_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "id1",
				DbName:        "pippo",
				Hostname:      "vm1",
				UsedLicenses:  42,
			},
		},
		Metadata: dto.PagingMetadata{},
	}

	usedLicensesMySQL := []dto.MySQLUsedLicense{}
	contracts := []model.MySQLContract{}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatasVm1, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&usedLicenses, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes1, nil),
		db.EXPECT().GetMySQLUsedLicenses("", filter).
			Return(usedLicensesMySQL, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
		db.EXPECT().GetClusters(filter).
			Return(clusters, nil),
	)

	db.EXPECT().ExistHostdata("vm1").Return(true, nil).AnyTimes()

	actual, actErr := as.GetUsedLicensesPerCluster(filter)
	require.NoError(t, actErr)

	expected := []dto.DatabaseUsedLicensePerCluster{
		{
			Cluster:       "name1",
			Hostnames:     []string{"vm1"},
			LicenseTypeID: "id1",
			Description:   "desc1",
			Metric:        model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:  150,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetUsedLicensesPerCluster_MultipleVms_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "id1",
				DbName:        "pippo",
				Hostname:      "vm1",
				UsedLicenses:  42,
			},
			{
				LicenseTypeID: "id1",
				DbName:        "pippo",
				Hostname:      "vm2",
				UsedLicenses:  42,
			},
			{
				LicenseTypeID: "id1",
				DbName:        "pippo",
				Hostname:      "vm3",
				UsedLicenses:  42,
			},
		},
		Metadata: dto.PagingMetadata{},
	}

	clusters := []dto.Cluster{
		{
			ID:                          [12]byte{},
			CreatedAt:                   time.Time{},
			Hostname:                    "cluster1",
			HostnameAgentVirtualization: "",
			Name:                        "name1",
			Environment:                 "",
			Location:                    "",
			FetchEndpoint:               "",
			CPU:                         12,
			Sockets:                     0,
			Type:                        "",
			VirtualizationNodes:         []string{},
			VirtualizationNodesCount:    0,
			VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
			VMs: []dto.VM{
				{
					CappedCPU:          false,
					Hostname:           "vm1",
					Name:               "",
					VirtualizationNode: "",
				},
				{
					CappedCPU:          false,
					Hostname:           "vm2",
					Name:               "",
					VirtualizationNode: "",
				},
			},
			VMsCount:            0,
			VMsErcoleAgentCount: 0,
		},
	}
	usedLicensesMySQL := []dto.MySQLUsedLicense{}
	contracts := []model.MySQLContract{}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatasVm1, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&usedLicenses, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes1, nil),
		db.EXPECT().GetMySQLUsedLicenses("", filter).
			Return(usedLicensesMySQL, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
		db.EXPECT().GetClusters(filter).
			Return(clusters, nil),
	)

	db.EXPECT().ExistHostdata("vm1").Return(true, nil).AnyTimes()
	db.EXPECT().ExistHostdata("vm2").Return(true, nil).AnyTimes()

	actual, actErr := as.GetUsedLicensesPerCluster(filter)
	require.NoError(t, actErr)

	expected := []dto.DatabaseUsedLicensePerCluster{
		{
			Cluster:       "name1",
			Hostnames:     []string{"vm1", "vm2"},
			LicenseTypeID: "id1",
			Description:   "desc1",
			Metric:        model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:  150,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetUsedLicensesPerClusterAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "id1",
				DbName:        "pippo",
				Hostname:      "vm1",
				UsedLicenses:  42,
			},
		},
		Metadata: dto.PagingMetadata{},
	}

	usedLicensesMySQL := []dto.MySQLUsedLicense{}
	contracts := []model.MySQLContract{}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatasVm1, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&usedLicenses, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes1, nil),
		db.EXPECT().GetMySQLUsedLicenses("", filter).
			Return(usedLicensesMySQL, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),
		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),
		db.EXPECT().GetClusters(filter).
			Return(clusters, nil),
	)

	db.EXPECT().ExistHostdata("vm1").Return(true, nil).AnyTimes()

	actual, err := as.GetUsedLicensesPerClusterAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "name1", actual.GetCellValue("Licenses Used Per Cluster", "A2"))
	assert.Equal(t, "id1", actual.GetCellValue("Licenses Used Per Cluster", "B2"))
	assert.Equal(t, "desc1", actual.GetCellValue("Licenses Used Per Cluster", "C2"))
	assert.Equal(t, "Named User Plus Perpetual", actual.GetCellValue("Licenses Used Per Cluster", "D2"))
	assert.Equal(t, "vm1", actual.GetCellValue("Licenses Used Per Cluster", "E2"))
	assert.Equal(t, "150", actual.GetCellValue("Licenses Used Per Cluster", "F2"))
}

func TestGetDatabaseLicensesComplianceSqlServerHostWithContractContract_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	oracleContracts := []dto.OracleDatabaseContractFE{}

	searchResponse := dto.OracleDatabaseUsedLicenseSearchResponse{}

	var oracleLicenseTypes = []model.OracleDatabaseLicenseType{}

	usedLicenses := []dto.MySQLUsedLicense{}

	clusters := []dto.Cluster{}

	contracts := []model.MySQLContract{}

	sqlServerContracts := []model.SqlServerDatabaseContract{
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeHost,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  licenseTypesID,
			Clusters:       []string{},
			Hosts:          []string{},
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return([]model.HostDataBE{{
			Hostname: "homer",
		}}, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(oracleContracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&searchResponse, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Times(1).
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("homer").Return(true, nil).AnyTimes()

	actual, err := as.GetDatabaseLicensesCompliance([]string{})
	require.NoError(t, err)

	expected := []dto.LicenseCompliance{
		{
			LicenseTypeID:   licenseTypesID,
			ItemDescription: "SQL Server Enterprise Edition",
			Metric:          "HOST",
			Cost:            0,
			Consumed:        8,
			Covered:         8,
			Purchased:       12,
			Compliance:      1,
			Unlimited:       false,
			Available:       4,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetDatabaseLicensesComplianceSqlServerHostNoContract_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	oracleContracts := []dto.OracleDatabaseContractFE{}

	searchResponse := dto.OracleDatabaseUsedLicenseSearchResponse{}

	var oracleLicenseTypes = []model.OracleDatabaseLicenseType{}

	usedLicenses := []dto.MySQLUsedLicense{}

	clusters := []dto.Cluster{}

	contracts := []model.MySQLContract{}

	sqlServerContracts := []model.SqlServerDatabaseContract{}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return([]model.HostDataBE{{
			Hostname: "homer",
		}}, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(oracleContracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&searchResponse, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Times(1).
			Return(sqlServerLicenseTypes, nil),
	)

	db.EXPECT().ExistHostdata("homer").Return(true, nil).AnyTimes()

	actual, err := as.GetDatabaseLicensesCompliance([]string{})
	require.NoError(t, err)

	expected := []dto.LicenseCompliance{
		{
			LicenseTypeID:   licenseTypesID,
			ItemDescription: "SQL Server Enterprise Edition",
			Metric:          "HOST",
			Cost:            0,
			Consumed:        8,
			Covered:         0,
			Purchased:       0,
			Compliance:      0,
			Unlimited:       false,
			Available:       0,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetDatabaseLicensesComplianceSqlServerHostInClusterWithContract_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	oracleContracts := []dto.OracleDatabaseContractFE{}

	searchResponse := dto.OracleDatabaseUsedLicenseSearchResponse{}

	var oracleLicenseTypes = []model.OracleDatabaseLicenseType{}

	usedLicenses := []dto.MySQLUsedLicense{}

	clusters := []dto.Cluster{
		{
			Hostname: "plutocluster",
			VMs: []dto.VM{
				{
					CappedCPU:          false,
					Hostname:           "plutohost",
					Name:               "plutohost",
					VirtualizationNode: "",
					IsErcoleInstalled:  true,
				},
			},
			Name: "PLUTO-CLUSTER-NAME",
			CPU:  40,
		},
	}

	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "plutocluster",
		HostnameAgentVirtualization: "",
		Name:                        "PLUTO-CLUSTER-NAME",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         40,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "plutohost",
				Name:               "plutohost",
				VirtualizationNode: "",
				IsErcoleInstalled:  true,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}

	contracts := []model.MySQLContract{}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: licenseTypesID,
				DbName:        "topolino-dbname",
				Hostname:      "plutohost",
				UsedLicenses:  8,
			},
			{
				LicenseTypeID: licenseTypesID,
				DbName:        "topolino-dbname",
				Hostname:      "plutocluster",
				UsedLicenses:  8,
			},
		},
	}

	sqlServerContracts := []model.SqlServerDatabaseContract{
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeCluster,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  licenseTypesID,
			Clusters:       []string{"PLUTO-CLUSTER-NAME"},
			Hosts:          []string{},
		},
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeHost,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  licenseTypesID,
			Clusters:       []string{},
			Hosts:          []string{"plutohost"},
		},
	}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              licenseTypesID,
			ItemDescription: "SQL Server Standard Edition",
			Edition:         "STD",
			Version:         "2019",
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return([]model.HostDataBE{{
			Hostname: "plutohost",
		}}, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().FindClusterVeritasLicenses(gomock.Any()).
		Return([]dto.ClusterVeritasLicense{}, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(oracleContracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&searchResponse, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(oracleLicenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetCluster("PLUTO-CLUSTER-NAME", utils.MAX_TIME).
			Return(&cluster, nil),
		db.EXPECT().ListSqlServerDatabaseContracts(gomock.Any()).
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Times(1).
			Return(sqlServerLicenseTypes, nil),
		db.EXPECT().GetCluster("PLUTO-CLUSTER-NAME", utils.MAX_TIME).
			Return(&cluster, nil),
		db.EXPECT().ExistHostdata("plutohost").
			Return(true, nil),
	)

	db.EXPECT().ExistHostdata("plutocluster").Return(true, nil).AnyTimes()

	actual, err := as.GetDatabaseLicensesCompliance([]string{})
	require.NoError(t, err)

	expected := []dto.LicenseCompliance{
		{
			LicenseTypeID:   licenseTypesID,
			ItemDescription: "SQL Server Standard Edition",
			Metric:          "CLUSTER",
			Cost:            0,
			Consumed:        40,
			Covered:         12,
			Purchased:       12,
			Compliance:      0.3,
			Unlimited:       false,
			Available:       0,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}
