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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestListManagedTechnologies_Success(t *testing.T) {
	var sampleLicenseTypes = []model.OracleDatabaseLicenseType{
		{
			ID:              "PID001",
			ItemDescription: "itemDesc1",
			Aliases:         []string{"alias1"},
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
		{
			ID:              "PID002",
			ItemDescription: "itemDesc2",
			Aliases:         []string{"alias2"},
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		},
		{
			ID:              "PID003",
			ItemDescription: "itemDesc3",
			Aliases:         []string{"alias3"},
			Metric:          model.LicenseTypeMetricComputerPerpetual,
		},
	}

	var sampleListOracleDatabaseContracts []dto.OracleDatabaseContractFE = []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			ContractID:               "",
			CSI:                      "",
			LicenseTypeID:            "PID001",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "pippo"}, {Hostname: "pluto"}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 50,
			AvailableLicensesPerUser: 0,
		},
		{
			ID:                       utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
			ContractID:               "",
			CSI:                      "",
			LicenseTypeID:            "PID002",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "topolino"}, {Hostname: "minnie"}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 75,
		},
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID001",
				DbName:        "",
				Hostname:      "test1",
				UsedLicenses:  3,
			},
			{
				LicenseTypeID: "PID001",
				DbName:        "",
				Hostname:      "pluto",
				UsedLicenses:  1.5,
			},
			{
				LicenseTypeID: "PID001",
				DbName:        "",
				Hostname:      "pippo",
				UsedLicenses:  5.5,
			},

			{
				LicenseTypeID: "PID002",
				DbName:        "",
				Hostname:      "topolino",
				UsedLicenses:  7,
			},
			{
				LicenseTypeID: "PID002",
				DbName:        "",
				Hostname:      "minnie",
				UsedLicenses:  4,
			},
			{
				LicenseTypeID: "PID003",
				DbName:        "",
				Hostname:      "minnie",
				UsedLicenses:  0.5,
			},
			{
				LicenseTypeID: "PID003",
				DbName:        "",
				Hostname:      "pippo",
				UsedLicenses:  0.5,
			},
			{
				LicenseTypeID: "PID003",
				DbName:        "",
				Hostname:      "test2",
				UsedLicenses:  4,
			},
			{
				LicenseTypeID: "PID003",
				DbName:        "",
				Hostname:      "test3",
				UsedLicenses:  6,
			},
		},
	}
	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "test-db",
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
	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}
	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: "359-06320",
				DbName:        "topolino-dbname",
				Hostname:      "plutohost",
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
			LicenseTypeID:  "359-06320",
			Clusters:       []string{},
			Hosts:          []string{},
		},
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeHost,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  "359-06320",
			Clusters:       []string{},
			Hosts:          []string{},
		},
	}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "359-06320",
			ItemDescription: "SQL Server Standard Edition",
			Edition:         "STD",
			Version:         "2019",
		},
	}
	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{},
			Hosts:            []string{},
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

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	db.EXPECT().GetHostDatas(utils.MAX_TIME).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
			Return(map[string]float64{
				model.TechnologyOracleDatabase:     42,
				model.TechnologyOracleExadata:      43,
				model.TechnologyOracleMySQL:        44,
				model.TechnologyMicrosoftSQLServer: 42,
			}, nil),
		db.EXPECT().
			ListOracleDatabaseContracts().
			Return(sampleListOracleDatabaseContracts, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(sampleLicenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts().
			Return(contracts, nil),
		db.EXPECT().GetMySQLContracts().
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts().
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().ListSqlServerDatabaseContracts().
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Times(1).
			Return(sqlServerLicenseTypes, nil),
	)

	actual, err := as.ListManagedTechnologies(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)
	require.NoError(t, err)

	expected := []model.TechnologyStatus{
		{Product: "Oracle/Database", ConsumedByHosts: 32, CoveredByContracts: 18, TotalCost: 0, PaidCost: 0, Compliance: 0.5625, UnpaidDues: 0, HostsCount: 42},
		{Product: "Oracle/MySQL", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 44},
		{Product: "Microsoft/SQLServer", ConsumedByHosts: 8, CoveredByContracts: 8, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 42},
		{Product: "PostgreSQL/PostgreSQL", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "MongoDB/MongoDB", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "MariaDBFoundation/MariaDB", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
	}

	assert.Equal(t, expected, actual)
}

func TestListManagedTechnologies_Success2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	returnedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          55,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID002",
				DbName:        "test-dbname",
				Hostname:      "test-db",
				UsedLicenses:  100,
			},
		},
	}
	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "test-db",
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
	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}
	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}

	var sampleLicenseTypes = []model.OracleDatabaseLicenseType{
		{
			ID:              "PID001",
			ItemDescription: "itemDesc1",
			Aliases:         []string{"alias1"},
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
		{
			ID:              "PID002",
			ItemDescription: "itemDesc2",
			Aliases:         []string{"alias2"},
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		},
		{
			ID:              "PID003",
			ItemDescription: "itemDesc3",
			Aliases:         []string{"alias3"},
			Metric:          model.LicenseTypeMetricComputerPerpetual,
		},
	}
	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: "359-06320",
				DbName:        "topolino-dbname",
				Hostname:      "plutohost",
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
			LicenseTypeID:  "359-06320",
			Clusters:       []string{},
			Hosts:          []string{},
		},
		{
			ID:             [12]byte{},
			Type:           model.SqlServerContractTypeHost,
			LicensesNumber: 12,
			ContractID:     "abc",
			LicenseTypeID:  "359-06320",
			Clusters:       []string{},
			Hosts:          []string{},
		},
	}

	sqlServerLicenseTypes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "359-06320",
			ItemDescription: "SQL Server Standard Edition",
			Edition:         "STD",
			Version:         "2019",
		},
	}

	contracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{},
			Hosts:            []string{},
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

	db.EXPECT().GetHostDatas(utils.MAX_TIME).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
			Return(map[string]float64{
				model.TechnologyOracleDatabase:     42,
				model.TechnologyOracleExadata:      43,
				model.TechnologyOracleMySQL:        44,
				model.TechnologyMicrosoftSQLServer: 42,
			}, nil),
		db.EXPECT().
			ListOracleDatabaseContracts().
			Return(returnedContracts, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(sampleLicenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts().
			Return(contracts, nil),
		db.EXPECT().GetMySQLContracts().
			Return(contracts, nil),

		db.EXPECT().SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().ListSqlServerDatabaseContracts().
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().ListSqlServerDatabaseContracts().
			Times(1).
			Return(sqlServerContracts, nil),
		db.EXPECT().GetSqlServerDatabaseLicenseTypes().
			Times(1).
			Return(sqlServerLicenseTypes, nil),
	)

	actual, err := as.ListManagedTechnologies(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	expected := []model.TechnologyStatus{
		{Product: "Oracle/Database", ConsumedByHosts: 100, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 42},
		{Product: "Oracle/MySQL", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 44},
		{Product: "Microsoft/SQLServer", ConsumedByHosts: 8, CoveredByContracts: 8, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 42},
		{Product: "PostgreSQL/PostgreSQL", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "MongoDB/MongoDB", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "MariaDBFoundation/MariaDB", ConsumedByHosts: 0, CoveredByContracts: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
	}

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestListManagedTechnologies_FailInternalServerErrors(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Fail GetHostsCountUsingTechnologies", func(t *testing.T) {
		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
			Return(nil, aerrMock)

		_, err := as.ListManagedTechnologies(
			"Count", true,
			"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
		)

		require.Equal(t, aerrMock, err)
	})

	t.Run("Fail ListOracleDatabaseContracts", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().
				GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
				Return(map[string]float64{
					model.TechnologyMariaDBFoundationMariaDB: 42,
					model.TechnologyMicrosoftSQLServer:       43,
				}, nil),
			db.EXPECT().
				ListOracleDatabaseContracts().
				Return(nil, aerrMock),
		)

		_, err := as.ListManagedTechnologies(
			"Count", true,
			"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
		)

		require.Equal(t, aerrMock, err)
	})
	t.Run("Fail ListHostUsingOracleDatabaseLicenses", func(t *testing.T) {
		var sampleListOracleDatabaseContracts []dto.OracleDatabaseContractFE = []dto.OracleDatabaseContractFE{
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				ContractID:               "",
				CSI:                      "",
				LicenseTypeID:            "PID001",
				ItemDescription:          "",
				Metric:                   "",
				ReferenceNumber:          "",
				Unlimited:                false,
				Basket:                   false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "pippo"}, {Hostname: "pluto"}},
				LicensesPerCore:          0,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 50,
				AvailableLicensesPerUser: 0,
			},
			{
				ID:                       utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				ContractID:               "",
				CSI:                      "",
				LicenseTypeID:            "PID002",
				ItemDescription:          "",
				Metric:                   "",
				ReferenceNumber:          "",
				Unlimited:                false,
				Basket:                   false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "topolino"}, {Hostname: "minnie"}},
				LicensesPerCore:          0,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 0,
				AvailableLicensesPerUser: 75,
			},
		}

		gomock.InOrder(
			db.EXPECT().
				GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
				Return(map[string]float64{
					model.TechnologyMariaDBFoundationMariaDB: 42,
					model.TechnologyMicrosoftSQLServer:       43,
				}, nil),
			db.EXPECT().
				ListOracleDatabaseContracts().
				Return(sampleListOracleDatabaseContracts, nil),
			db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
				Return(nil, aerrMock),
		)

		_, err := as.ListManagedTechnologies(
			"Count", true,
			"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
		)

		require.Equal(t, aerrMock, err)
	})
}
