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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
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
		Log:      logger.NewLogger("TEST"),
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
			UsedLicenses:  2,
		},
			{
				LicenseTypeID: "A98765",
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
				UsedLicenses:  2,
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
		{
			ID:              "A98765",
			ItemDescription: "ThisDesc",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
	}
	usedLicensesMySQL := []dto.MySQLUsedLicense{
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
	hostdatas := []model.HostDataBE{
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
			Return(usedLicensesMySQL, nil),
		db.EXPECT().GetClusters(any).
			Return(clusters, nil),
		db.EXPECT().GetMySQLAgreements().
			Return(agreements, nil),

		db.EXPECT().GetHostDatas(utils.MAX_TIME).
			Return(hostdatas, nil),
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
			UsedLicenses:  2,
		},
		{
			Hostname:      "topolino-hostname",
			DbName:        "topolino-dbname",
			LicenseTypeID: "A98765",
			Description:   "ThisDesc",
			Metric:        model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:  50,
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
		Log:      logger.NewLogger("TEST"),
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
	hostdatas := []model.HostDataBE{
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

		db.EXPECT().GetHostDatas(utils.MAX_TIME).
			Return(hostdatas, nil),
	)

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

	objID, _ := primitive.ObjectIDFromHex("609ce4782eff5d5540ec4a30")

	oracleAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:              objID,
			AgreementID:     "5051863",
			CSI:             "18000742",
			LicenseTypeID:   "L47225",
			ItemDescription: "Oracle Advanced Compression",
			Metric:          "Named User Plus Perpetual",
			ReferenceNumber: "66880702",
			Unlimited:       false,
			Basket:          false,
			Restricted:      false,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
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

	oracleHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "L47225",
			Name:          "sdlsts103",
			Type:          "host",
			LicenseCount:  6,
			OriginalCount: 6,
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
		db.EXPECT().ListOracleDatabaseAgreements().
			Return(oracleAgreements, nil),
		db.EXPECT().ListHostUsingOracleDatabaseLicenses().
			Return(oracleHosts, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(oracleLicenseTypes, nil).
			Times(2),
		db.EXPECT().GetHostDatas(utils.MAX_TIME).
			Return([]model.HostDataBE{{
				Hostname: "sdlsts103",
			}}, nil),

		db.EXPECT().GetMySQLUsedLicenses(any).
			Return(usedLicenses, nil),
		db.EXPECT().GetClusters(any).
			Return(clusters, nil),
		db.EXPECT().GetMySQLAgreements().
			Return(agreements, nil),
	)

	actual, err := as.GetDatabaseLicensesComplianceAsXLSX()
	require.NoError(t, err)

	assert.Equal(t, "L47225", actual.GetCellValue("Licenses Compliance", "A2"))
	assert.Equal(t, "Oracle Advanced Compression", actual.GetCellValue("Licenses Compliance", "B2"))
	assert.Equal(t, "Named User Plus Perpetual", actual.GetCellValue("Licenses Compliance", "C2"))
	assert.Equal(t, "150", actual.GetCellValue("Licenses Compliance", "D2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "E2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "F2"))
	assert.Equal(t, "0", actual.GetCellValue("Licenses Compliance", "G2"))
	assert.Equal(t, "", actual.GetCellValue("Licenses Compliance", "H2"))
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

	objID, _ := primitive.ObjectIDFromHex("609ce4782eff5d5540ec4a30")

	oracleAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       objID,
			AgreementID:              "5051863",
			CSI:                      "19338486",
			LicenseTypeID:            "L10006",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   "Named User Plus Perpetual",
			ReferenceNumber:          "96661555",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          450,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 450,
			CoveredLicenses:          0,
		},
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "TEST",
			CSI:                      "99999999",
			LicenseTypeID:            "M10080",
			ItemDescription:          "Oracle Database Enterprise Edition",
			Metric:                   "Processor Perpetual",
			ReferenceNumber:          "666666666",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          50,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 50,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          0,
		},
	}

	oracleHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "L10006",
			Name:          "sdlmoc100.syssede.systest.sanpaoloimi.com",
			Type:          "host",
			LicenseCount:  10,
			OriginalCount: 20,
		},
		{
			LicenseTypeID: "M10080",
			Name:          "sdlmoc100.syssede.systest.sanpaoloimi.com",
			Type:          "host",
			LicenseCount:  20,
			OriginalCount: 50,
		},
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
		db.EXPECT().ListOracleDatabaseAgreements().
			Return(oracleAgreements, nil),
		db.EXPECT().ListHostUsingOracleDatabaseLicenses().
			Return(oracleHosts, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(oracleLicenseTypes, nil).
			Times(2),
		db.EXPECT().GetHostDatas(utils.MAX_TIME).
			Return([]model.HostDataBE{{
				Hostname: "sdlmoc100.syssede.systest.sanpaoloimi.com",
			}}, nil),

		db.EXPECT().GetMySQLUsedLicenses(any).
			Return(usedLicenses, nil),
		db.EXPECT().GetClusters(any).
			Return(clusters, nil),
		db.EXPECT().GetMySQLAgreements().
			Return(agreements, nil),
	)

	actual, err := as.GetDatabaseLicensesCompliance()
	require.NoError(t, err)

	expected := []dto.LicenseCompliance{
		{
			LicenseTypeID:   "M10080",
			ItemDescription: "Application Testing",
			Metric:          "Processor Perpetual",
			Consumed:        50,
			Covered:         20,
			Purchased:       50,
			Compliance:      0.4,
			Unlimited:       false,
			Available:       30,
		},
		{
			LicenseTypeID:   "L10006",
			ItemDescription: "Oracle Partitioning",
			Metric:          "Named User Plus Perpetual",
			Consumed:        500,
			Covered:         250,
			Purchased:       450,
			Compliance:      0.5,
			Unlimited:       false,
			Available:       200,
		},
		{
			LicenseTypeID:   "",
			ItemDescription: "MySQL Enterprise per cluster",
			Metric:          "",
			Consumed:        1,
			Covered:         1,
			Purchased:       0,
			Compliance:      1,
			Unlimited:       false,
			Available:       0,
		},
	}
	assert.ElementsMatch(t, expected, actual)
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
	hostdatas := []model.HostDataBE{
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

		db.EXPECT().GetHostDatas(utils.MAX_TIME).
			Return(hostdatas, nil),
	)

	actual, err := as.GetDatabasesUsedLicensesPerHostAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Licenses Used Per Host", "A2"))
	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Licenses Used Per Host", "B2"))
	assert.Equal(t, "A90611", actual.GetCellValue("Licenses Used Per Host", "C2"))
	assert.Equal(t, "Oracle Database Enterprise Edition", actual.GetCellValue("Licenses Used Per Host", "D2"))
	assert.Equal(t, "Processor Perpetual", actual.GetCellValue("Licenses Used Per Host", "E2"))
	assert.Equal(t, "2", actual.GetCellValue("Licenses Used Per Host", "F2"))
	assert.Equal(t, "", actual.GetCellValue("Licenses Used Per Host", "H2"))
}

func TestGetDatabasesUsedLicensesPerCluster_OneVm_Success(t *testing.T) {
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

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "id1",
			ItemDescription: "desc1",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
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
			},
			VMsCount:            0,
			VMsErcoleAgentCount: 0,
		},
	}
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&usedLicenses, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().Return(licenseTypes, nil),
		db.EXPECT().GetClusters(filter).Return(clusters, nil),
	)

	actual, actErr := as.GetDatabasesUsedLicensesPerCluster(filter)
	require.NoError(t, actErr)

	expected := []dto.DatabaseUsedLicensePerCluster{
		{
			Cluster:       "name1",
			Hostnames:     []string{"vm1"},
			LicenseTypeID: "id1",
			Description:   "desc1",
			Metric:        model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:  6,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetDatabasesUsedLicensesPerCluster_MultipleVms_Success(t *testing.T) {
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

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "id1",
			ItemDescription: "desc1",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
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
	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&usedLicenses, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().Return(licenseTypes, nil),
		db.EXPECT().GetClusters(filter).Return(clusters, nil),
	)

	actual, actErr := as.GetDatabasesUsedLicensesPerCluster(filter)
	require.NoError(t, actErr)

	expected := []dto.DatabaseUsedLicensePerCluster{
		{
			Cluster:       "name1",
			Hostnames:     []string{"vm1", "vm2"},
			LicenseTypeID: "id1",
			Description:   "desc1",
			Metric:        model.LicenseTypeMetricNamedUserPlusPerpetual,
			UsedLicenses:  6,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestGetDatabasesUsedLicensesPerClusterAsXLSX_Success(t *testing.T) {
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

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "id1",
			ItemDescription: "desc1",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Cost:            0,
			Aliases:         []string{},
			Option:          false,
		},
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
			},
			VMsCount:            0,
			VMsErcoleAgentCount: 0,
		},
	}

	gomock.InOrder(
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
			Return(&usedLicenses, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().Return(licenseTypes, nil),
		db.EXPECT().GetClusters(filter).Return(clusters, nil),
	)

	actual, err := as.GetDatabasesUsedLicensesPerClusterAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "name1", actual.GetCellValue("Licenses Used Per Cluster", "A2"))
	assert.Equal(t, "id1", actual.GetCellValue("Licenses Used Per Cluster", "B2"))
	assert.Equal(t, "desc1", actual.GetCellValue("Licenses Used Per Cluster", "C2"))
	assert.Equal(t, "Named User Plus Perpetual", actual.GetCellValue("Licenses Used Per Cluster", "D2"))
	assert.Equal(t, "vm1", actual.GetCellValue("Licenses Used Per Cluster", "E2"))
	assert.Equal(t, "6", actual.GetCellValue("Licenses Used Per Cluster", "F2"))
}
