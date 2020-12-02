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

	"github.com/ercole-io/ercole/api-service/dto"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var sampleParts = []model.OracleDatabasePart{
	{
		PartID:          "PID001",
		ItemDescription: "itemDesc1",
		Aliases:         []string{"alias1"},
		Metric:          model.AgreementPartMetricProcessorPerpetual,
	},
	{
		PartID:          "PID002",
		ItemDescription: "itemDesc2",
		Aliases:         []string{"alias2"},
		Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
	},
	{
		PartID:          "PID003",
		ItemDescription: "itemDesc3",
		Aliases:         []string{"alias3"},
		Metric:          model.AgreementPartMetricComputerPerpetual,
	},
}

var sampleListOracleDatabaseAgreements []dto.OracleDatabaseAgreementFE = []dto.OracleDatabaseAgreementFE{
	{
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		AgreementID:     "",
		CSI:             "",
		PartID:          "PID001",
		ItemDescription: "",
		Metric:          "",
		ReferenceNumber: "",
		Unlimited:       false,
		Count:           50,
		CatchAll:        false,
		Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
			{Hostname: "pippo"},
			{Hostname: "pluto"},
		},
		AvailableCount: 50,
		LicensesCount:  0,
		UsersCount:     0,
	},
	{
		ID:              utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
		AgreementID:     "",
		CSI:             "",
		PartID:          "PID002",
		ItemDescription: "",
		Metric:          "",
		ReferenceNumber: "",
		Unlimited:       false,
		Count:           75,
		CatchAll:        false,
		Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
			{Hostname: "topolino"},
			{Hostname: "minnie"},
		},
		AvailableCount: 75,
		LicensesCount:  0,
		UsersCount:     0,
	},
}

var sampleHostUsingOracleDbLicenses []dto.HostUsingOracleDatabaseLicenses = []dto.HostUsingOracleDatabaseLicenses{
	{LicenseName: "alias1", Name: "test1", Type: "host", LicenseCount: 3, OriginalCount: 3},
	{LicenseName: "alias1", Name: "pluto", Type: "host", LicenseCount: 1.5, OriginalCount: 1.5},
	{LicenseName: "alias1", Name: "pippo", Type: "host", LicenseCount: 5.5, OriginalCount: 5.5},

	{LicenseName: "alias2", Name: "topolino", Type: "cluster", LicenseCount: 7, OriginalCount: 7},
	{LicenseName: "alias2", Name: "minnie", Type: "host", LicenseCount: 4, OriginalCount: 4},
	{LicenseName: "alias2", Name: "minnie", Type: "host", LicenseCount: 8, OriginalCount: 8},

	{LicenseName: "alias3", Name: "minnie", Type: "host", LicenseCount: 0.5, OriginalCount: 0.5},
	{LicenseName: "alias3", Name: "pippo", Type: "host", LicenseCount: 0.5, OriginalCount: 0.5},
	{LicenseName: "alias3", Name: "test2", Type: "host", LicenseCount: 4, OriginalCount: 4},
	{LicenseName: "alias3", Name: "test3", Type: "cluster", LicenseCount: 6, OriginalCount: 6},
}

func TestListManagedTechnologies_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:                     db,
		Log:                          utils.NewLogger("TEST"),
		OracleDatabaseAgreementParts: sampleParts,
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
			Return(sampleListOracleDatabaseAgreements, nil),
		db.EXPECT().
			ListHostUsingOracleDatabaseLicenses().
			Return(sampleHostUsingOracleDbLicenses, nil),
	)

	actual, err := as.ListManagedTechnologies(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	expected := []model.TechnologyStatus{
		{Product: "Oracle/Database", ConsumedByHosts: 40, CoveredByAgreements: 10, TotalCost: 0, PaidCost: 0, Compliance: 0.25, UnpaidDues: 0, HostsCount: 42},
		{Product: "MariaDBFoundation/MariaDB", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "PostgreSQL/PostgreSQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "Oracle/MySQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "Microsoft/SQLServer", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
	}

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestListManagedTechnologies_Success2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:                     db,
		Log:                          utils.NewLogger("TEST"),
		OracleDatabaseAgreementParts: sampleParts,
	}

	returnedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 0,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			LicensesCount:   55,
			UsersCount:      0,
			Count:           55,
		},
	}
	returnedHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  100,
			LicenseName:   "Partitioning",
			OriginalCount: 100,
			Type:          "host",
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
			ListHostUsingOracleDatabaseLicenses().
			Return(returnedHosts, nil),
	)

	actual, err := as.ListManagedTechnologies(
		"Count", true,
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	expected := []model.TechnologyStatus{
		{Product: "Oracle/Database", ConsumedByHosts: 100, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 0, UnpaidDues: 0, HostsCount: 42},
		{Product: "MariaDBFoundation/MariaDB", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "PostgreSQL/PostgreSQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "Oracle/MySQL", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
		{Product: "Microsoft/SQLServer", ConsumedByHosts: 0, CoveredByAgreements: 0, TotalCost: 0, PaidCost: 0, Compliance: 1, UnpaidDues: 0, HostsCount: 0},
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
			db.EXPECT().
				ListHostUsingOracleDatabaseLicenses().
				Return(nil, aerrMock),
		)

		_, err := as.ListManagedTechnologies(
			"Count", true,
			"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
		)

		require.Equal(t, aerrMock, err)
	})
}
