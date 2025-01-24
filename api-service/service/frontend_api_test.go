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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetInfoForFrontendDashboard_Success(t *testing.T) {
	t.Skip("writing new code on this API")

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
					"compliance": 1,
					"unpaidDues": 0,
					"product":    model.TechnologyPostgreSQLPostgreSQL,
					"hostsCount": 0,
				},
				{
					"compliance": 1,
					"unpaidDues": 0,
					"product":    model.TechnologyMongoDBMongoDB,
					"hostsCount": 0,
				},
				{
					"compliance": 0,
					"unpaidDues": 0,
					"product":    model.TechnologyMariaDBFoundationMariaDB,
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
			ListOracleDatabaseContracts(gomock.Any()).
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
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(mySqlcontracts, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(mySqlcontracts, nil),

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

		db.EXPECT().
			GetHostsCountStats("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(20, nil),
		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
			Return(getTechnologiesUsageRes2, nil),

		db.EXPECT().
			ListOracleDatabaseContracts(gomock.Any()).
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
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(mySqlcontracts, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(mySqlcontracts, nil),

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

		db.EXPECT().
			GetHostsCountUsingTechnologies("", "", utils.MAX_TIME).
			Return(getTechnologiesUsageRes, nil),
	)

	res, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetInfoForFrontendDashboard_Fail1(t *testing.T) {
	t.Skip("writing new code on this API")

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
	t.Skip("writing new code on this API")

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
			ListOracleDatabaseContracts(gomock.Any()).
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
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(mySqlcontracts, nil),
		db.EXPECT().GetMySQLContracts(gomock.Any()).
			Return(mySqlcontracts, nil),

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

func TestGetComplianceStatsAsAdmin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
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

	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
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

	var instancesCount, hostsCount int64 = 1, 1

	oracleContracts := []dto.OracleDatabaseContractFE{
		{
			ItemDescription: "foobar",
		},
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
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: licenseTypesID,
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
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
			ID:              licenseTypesID,
			ItemDescription: "SQL Server Enterprise Edition",
			Edition:         "ENT",
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

	db.EXPECT().
		GetOracleDatabaseLicenseTypes().
		Return(licenseTypes, nil).
		AnyTimes()

	db.EXPECT().
		GetClusters(globalFilterAny).
		Return(clusters, nil).
		AnyTimes()

	db.EXPECT().
		GetHostDatas(utils.MAX_TIME).
		Times(1).
		Return(hostdatas, nil).
		AnyTimes()

	gomock.InOrder(
		db.EXPECT().
			CountOracleInstance().
			Return(instancesCount, nil),
		db.EXPECT().
			CountOracleHosts().
			Return(hostsCount, nil),
		db.EXPECT().
			ListOracleDatabaseContracts(gomock.Any()).
			Return(oracleContracts, nil),
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),

		db.EXPECT().
			CountSqlServerlInstance().
			Return(instancesCount, nil),
		db.EXPECT().
			CountSqlServerHosts().
			Return(hostsCount, nil),
		db.EXPECT().
			SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().
			ListSqlServerDatabaseContracts(gomock.Any()).
			Times(2).
			Return(sqlServerContracts, nil),
		db.EXPECT().
			GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),

		db.EXPECT().
			CountMySqlInstance().
			Return(instancesCount, nil),
		db.EXPECT().
			CountMySqlHosts().
			Return(hostsCount, nil),
		db.EXPECT().
			GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().
			GetMySQLContracts(gomock.Any()).
			Times(2).
			Return(mySqlcontracts, nil),

		db.EXPECT().
			CountPostgreSqlInstance().
			Return(instancesCount, nil),
		db.EXPECT().
			CountPostgreSqlHosts().
			Return(hostsCount, nil),

		db.EXPECT().
			CountMongoDbInstance().
			Return(instancesCount, nil),
		db.EXPECT().
			CountMongoDbHosts().
			Return(hostsCount, nil),
	)

	user := model.User{
		Groups: []string{model.GroupAdmin},
	}
	expectedRes := map[string]interface{}{
		"ercole": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   5,
			"hostCount":               5,
		},
		"mariaDb": map[string]interface{}{
			"compliancePercentageStr": "100%",
			"compliancePercentageVal": 100,
			"count":                   0,
			"hostCount":               0,
		},
		"mongoDb": map[string]interface{}{
			"compliancePercentageStr": "100%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"mySql": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"oracle": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"postgreSql": map[string]interface{}{
			"compliancePercentageStr": "100%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"sqlServer": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
	}

	res, err := as.GetComplianceStats(user)

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetComplianceStatsAsUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
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

	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
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

	var instancesCount, hostsCount int64 = 1, 1

	oracleContracts := []dto.OracleDatabaseContractFE{
		{
			ItemDescription: "foobar",
		},
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
		},
	}

	sqlServerLics := dto.SqlServerDatabaseUsedLicenseSearchResponse{
		Content: []dto.SqlServerDatabaseUsedLicense{
			{
				LicenseTypeID: licenseTypesID,
				DbName:        "topolino-dbname",
				Hostname:      "topolino-hostname",
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
			ID:              licenseTypesID,
			ItemDescription: "SQL Server Enterprise Edition",
			Edition:         "ENT",
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

	user := model.User{
		Username: "foobar",
	}

	locations := []string{"wonderland"}

	db.EXPECT().
		GetUserLocations(user.Username).
		Return(locations, nil)

	db.EXPECT().
		GetOracleDatabaseLicenseTypes().
		Return(licenseTypes, nil).
		AnyTimes()

	db.EXPECT().
		GetClusters(globalFilterAny).
		Return(clusters, nil).
		AnyTimes()

	db.EXPECT().
		GetHostDatas(utils.MAX_TIME).
		Times(1).
		Return(hostdatas, nil).
		AnyTimes()

	gomock.InOrder(
		db.EXPECT().
			CountOracleInstanceByLocations(locations).
			Return(instancesCount, nil),
		db.EXPECT().
			CountOracleHostsByLocations(locations).
			Return(hostsCount, nil),
		db.EXPECT().
			ListOracleDatabaseContracts(gomock.Any()).
			Return(oracleContracts, nil),
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),

		db.EXPECT().
			CountSqlServerlInstanceByLocations(locations).
			Return(instancesCount, nil),
		db.EXPECT().
			CountSqlServerHostsByLocations(locations).
			Return(hostsCount, nil),
		db.EXPECT().
			SearchSqlServerDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&sqlServerLics, nil),
		db.EXPECT().
			ListSqlServerDatabaseContracts(gomock.Any()).
			Times(2).
			Return(sqlServerContracts, nil),
		db.EXPECT().
			GetSqlServerDatabaseLicenseTypes().
			Return(sqlServerLicenseTypes, nil),

		db.EXPECT().
			CountMySqlInstanceByLocations(locations).
			Return(instancesCount, nil),
		db.EXPECT().
			CountMySqlHostsByLocations(locations).
			Return(hostsCount, nil),
		db.EXPECT().
			GetMySQLUsedLicenses("", globalFilterAny).
			Return(usedLicenses, nil),
		db.EXPECT().
			GetMySQLContracts(gomock.Any()).
			Times(2).
			Return(mySqlcontracts, nil),

		db.EXPECT().
			CountPostgreSqlInstanceByLocations(locations).
			Return(instancesCount, nil),
		db.EXPECT().
			CountPostgreSqlHostsByLocations(locations).
			Return(hostsCount, nil),

		db.EXPECT().
			CountMongoDbInstanceByLocations(locations).
			Return(instancesCount, nil),
		db.EXPECT().
			CountMongoDbHostsByLocations(locations).
			Return(hostsCount, nil),
	)

	expectedRes := map[string]interface{}{
		"ercole": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   5,
			"hostCount":               5,
		},
		"mariaDb": map[string]interface{}{
			"compliancePercentageStr": "100%",
			"compliancePercentageVal": 100,
			"count":                   0,
			"hostCount":               0,
		},
		"mongoDb": map[string]interface{}{
			"compliancePercentageStr": "100%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"mySql": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"oracle": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"postgreSql": map[string]interface{}{
			"compliancePercentageStr": "100%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
		"sqlServer": map[string]interface{}{
			"compliancePercentageStr": "100.00%",
			"compliancePercentageVal": 100,
			"count":                   1,
			"hostCount":               1,
		},
	}

	res, err := as.GetComplianceStats(user)

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}
