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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	dto "github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetTotalTechnologiesComplianceStats_Success(t *testing.T) {
	t.Skip("writing new code on this API")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	getTechnologiesUsageRes := map[string]float64{
		"Oracle/Exadata": 2,
	}
	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().
		GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(getTechnologiesUsageRes, nil)
	db.EXPECT().
		GetHostsCountStats("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(20, nil).AnyTimes().MinTimes(1)
	db.EXPECT().
		GetOracleDatabaseLicenseTypes().
		Return(licenseTypes, nil)

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
			AvailableLicensesPerCore: 55,
			AvailableLicensesPerUser: 0,
		},
	}
	db.EXPECT().ListOracleDatabaseContracts().Return(returnedContracts, nil)

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
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),

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

	res, err := as.GetTotalTechnologiesComplianceStats(
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.NoError(t, err)

	expectedRes := map[string]interface{}{
		"compliance": 0.5833333333333334,
		"unpaidDues": 0,
		"hostsCount": 20,
	}
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetTotalTechnologiesComplianceStats_FailInternalServerError(t *testing.T) {
	t.Skip("writing new code on this API")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(nil, aerrMock)

	_, err := as.GetTotalTechnologiesComplianceStats(
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.Equal(t, aerrMock, err)
}
