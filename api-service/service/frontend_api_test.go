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

package service

import (
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInfoForFrontendDashboard_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"technologies": map[string]interface{}{
			"technologies": []map[string]interface{}{
				{
					"compliance": 0,
					"product":    model.TechnologyOracleDatabase,
					"hostsCount": 8,
					"unpaidDues": 0,
				},
				{
					"compliance": 1,
					"unpaidDues": 0,
					"product":    model.TechnologyOracleMySQL,
					"hostsCount": 0,
				},
				{
					"compliance": 1,
					"unpaidDues": 0,
					"product":    model.TechnologyMicrosoftSQLServer,
					"hostsCount": 8,
				},
				{
					"compliance": 0,
					"unpaidDues": 0,
					"product":    model.TechnologyMariaDBFoundationMariaDB,
					"hostsCount": 0,
				},
				{
					"compliance": 0,
					"unpaidDues": 0,
					"product":    model.TechnologyPostgreSQLPostgreSQL,
					"hostsCount": 0,
				},
			},
			"total": map[string]interface{}{
				"compliance": 0.10256410256410256,
				"unpaidDues": 0,
				"hostsCount": 20,
			},
		},
		"features": map[string]interface{}{
			"Oracle/Database": true,
			"Oracle/Exadata":  false,
		},
	}

	getTechnologiesUsageRes := map[string]float64{
		"Oracle/Database":     8,
		"Oracle/Exadata":      0,
		"Microsoft/SQLServer": 8,
	}

	contracts := []dto.OracleDatabaseContractFE{
		{
			LicenseTypeID:   "PID002",
			ItemDescription: "foobar",
		},
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID002",
				Hostname:      "test-db",
				DbName:        "",
				UsedLicenses:  70,
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

	getTechnologiesUsageRes2 := map[string]float64{
		"Oracle/Database":     8,
		"Oracle/Exadata":      2,
		"Microsoft/SQLServer": 8,
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

	usedLicenses := []dto.MySQLUsedLicense{
		{
			Hostname:        "pluto",
			InstanceName:    "pluto-instance",
			InstanceEdition: model.MySQLEditionEnterprise,
			ContractType:    "",
		},
	}

	mySqlcontracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{},
			Hosts:            []string{},
		},
	}

	db.EXPECT().GetHostDatas(utils.MAX_TIME).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	gomock.InOrder(

		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(getTechnologiesUsageRes, nil),

		db.EXPECT().
			ListOracleDatabaseContracts().
			Return(contracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts().
			Return(mySqlcontracts, nil),
		db.EXPECT().GetMySQLContracts().
			Return(mySqlcontracts, nil),

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

		db.EXPECT().
			GetHostsCountStats("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(20, nil),
		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(getTechnologiesUsageRes2, nil),

		db.EXPECT().
			ListOracleDatabaseContracts().
			Return(contracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts().
			Return(mySqlcontracts, nil),
		db.EXPECT().GetMySQLContracts().
			Return(mySqlcontracts, nil),

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

		db.EXPECT().
			GetHostsCountUsingTechnologies("", "", utils.MAX_TIME).
			Return(getTechnologiesUsageRes, nil),
	)

	res, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetInfoForFrontendDashboard_Fail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetHostsCountUsingTechnologies("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(nil, aerrMock).AnyTimes().MinTimes(1)

	_, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.Equal(t, aerrMock, err)
}

func TestGetInfoForFrontendDashboard_Fail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	getTechnologiesUsageRes := map[string]float64{
		"Oracle/Database": 8,
		"Oracle/Exadata":  0,
	}

	contracts := []dto.OracleDatabaseContractFE{
		{
			ItemDescription: "foobar",
		},
	}

	ltRes := []model.OracleDatabaseLicenseType{
		{
			ItemDescription: "foobar",
		},
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "A90649",
				Hostname:      "test-db",
				DbName:        "",
				UsedLicenses:  70,
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

	usedLicenses := []dto.MySQLUsedLicense{
		{
			Hostname:        "pluto",
			InstanceName:    "pluto-instance",
			InstanceEdition: model.MySQLEditionEnterprise,
			ContractType:    "",
		},
	}

	mySqlcontracts := []model.MySQLContract{
		{
			ID:               [12]byte{},
			Type:             model.MySQLContractTypeCluster,
			NumberOfLicenses: 12,
			Clusters:         []string{},
			Hosts:            []string{},
		},
	}

	db.EXPECT().GetHostDatas(utils.MAX_TIME).
		Times(1).
		Return(hostdatas, nil).AnyTimes()

	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	gomock.InOrder(

		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(getTechnologiesUsageRes, nil).AnyTimes().MinTimes(1),

		db.EXPECT().
			ListOracleDatabaseContracts().
			Return(contracts, nil).AnyTimes().MinTimes(1),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Return(ltRes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Return(ltRes, nil),

		db.EXPECT().GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().GetMySQLContracts().
			Return(mySqlcontracts, nil),
		db.EXPECT().GetMySQLContracts().
			Return(mySqlcontracts, nil),

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

		db.EXPECT().
			GetHostsCountStats("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(20, nil).AnyTimes().MinTimes(1),

		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(nil, aerrMock),
	)

	_, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.Equal(t, aerrMock, err)
}
