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

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	var sampleListOracleDatabaseAgreements []dto.OracleDatabaseAgreementFE = []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			AgreementID:              "",
			CSI:                      "",
			LicenseTypeID:            "PID001",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "pippo"}, {Hostname: "pluto"}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 50,
			AvailableLicensesPerUser: 0,
		},
		{
			ID:                       utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
			AgreementID:              "",
			CSI:                      "",
			LicenseTypeID:            "PID002",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "topolino"}, {Hostname: "minnie"}},
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

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	gomock.InOrder(
		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
			Return(map[string]float64{
				model.TechnologyOracleDatabase: 42,
				model.TechnologyOracleExadata:  43,
				model.TechnologyOracleMySQL:    44,
			}, nil),
		db.EXPECT().
			ListOracleDatabaseAgreements().
			Return(sampleListOracleDatabaseAgreements, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetHostDatas(utils.MAX_TIME).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),

		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(sampleLicenseTypes, nil),
	)

	actual, err := as.ListManagedTechnologies(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)
	require.NoError(t, err)

	expected := []model.TechnologyStatus{
		{Product: "Oracle/Database", ConsumedByHosts: 32, CoveredByAgreements: 18, TotalCost: 0, PaidCost: 0, Compliance: 0.5625, UnpaidDues: 0, HostsCount: 42},
		{Product: "Oracle/MySQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 44},
		{Product: "MariaDBFoundation/MariaDB", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
		{Product: "PostgreSQL/PostgreSQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
		{Product: "Microsoft/SQLServer", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
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

	returnedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
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
	gomock.InOrder(
		db.EXPECT().
			GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
			Return(map[string]float64{
				model.TechnologyOracleDatabase: 42,
				model.TechnologyOracleExadata:  43,
			}, nil),
		db.EXPECT().
			ListOracleDatabaseAgreements().
			Return(returnedAgreements, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetHostDatas(utils.MAX_TIME).
			Return(hostdatas, nil),
		db.EXPECT().GetClusters(globalFilterAny).
			Return(clusters, nil),

		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(sampleLicenseTypes, nil),
	)

	actual, err := as.ListManagedTechnologies(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	expected := []model.TechnologyStatus{
		{Product: "Oracle/Database", ConsumedByHosts: 100, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 42},
		{Product: "Oracle/MySQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
		{Product: "MariaDBFoundation/MariaDB", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
		{Product: "PostgreSQL/PostgreSQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
		{Product: "Microsoft/SQLServer", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 0},
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

	t.Run("Fail ListOracleDatabaseAgreements", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().
				GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
				Return(map[string]float64{
					model.TechnologyMariaDBFoundationMariaDB: 42,
					model.TechnologyMicrosoftSQLServer:       43,
				}, nil),
			db.EXPECT().
				ListOracleDatabaseAgreements().
				Return(nil, aerrMock),
		)

		_, err := as.ListManagedTechnologies(
			"Count", true,
			"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
		)

		require.Equal(t, aerrMock, err)
	})
	t.Run("Fail ListHostUsingOracleDatabaseLicenses", func(t *testing.T) {
		var sampleListOracleDatabaseAgreements []dto.OracleDatabaseAgreementFE = []dto.OracleDatabaseAgreementFE{
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				AgreementID:              "",
				CSI:                      "",
				LicenseTypeID:            "PID001",
				ItemDescription:          "",
				Metric:                   "",
				ReferenceNumber:          "",
				Unlimited:                false,
				Basket:                   false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "pippo"}, {Hostname: "pluto"}},
				LicensesPerCore:          0,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 50,
				AvailableLicensesPerUser: 0,
			},
			{
				ID:                       utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				AgreementID:              "",
				CSI:                      "",
				LicenseTypeID:            "PID002",
				ItemDescription:          "",
				Metric:                   "",
				ReferenceNumber:          "",
				Unlimited:                false,
				Basket:                   false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "topolino"}, {Hostname: "minnie"}},
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
				ListOracleDatabaseAgreements().
				Return(sampleListOracleDatabaseAgreements, nil),
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
