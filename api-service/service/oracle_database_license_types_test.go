// Copyright (c) 2021 Sorint.lab S.p.A.
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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetOracleDatabaseLicenseTypes_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []model.OracleDatabaseLicenseType{
		{
			ID:              "Pippo",
			ItemDescription: "Pluto",
			Metric:          "Topolino",
			Cost:            12,
			Aliases:         []string{"Minny"},
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().Return(expected, nil)

	res, err := as.GetOracleDatabaseLicenseTypes()
	require.NoError(t, err)
	assert.Equal(t, []model.OracleDatabaseLicenseType{
		{
			ID:              "Pippo",
			ItemDescription: "Pluto",
			Metric:          "Topolino",
			Cost:            12,
			Aliases:         []string{"Minny"},
		},
	}, res)
}

func TestGetLicensesCompliance(t *testing.T) {
	var sampleAgreements []dto.OracleDatabaseAgreementFE = []dto.OracleDatabaseAgreementFE{
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
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 10,
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
			LicensesPerUser:          500,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 500,
		},
		{
			ID:                       utils.Str2oid("cccccccccccccccccccccccc"),
			AgreementID:              "",
			CSI:                      "",
			LicenseTypeID:            "PID003",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                true,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
		{
			ID:                       utils.Str2oid("dddddddddddddddddddddddd"),
			AgreementID:              "",
			CSI:                      "",
			LicenseTypeID:            "PID004",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          40,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 40,
			AvailableLicensesPerUser: 0,
		},
	}

	var sampleHosts []dto.HostUsingOracleDatabaseLicenses = []dto.HostUsingOracleDatabaseLicenses{
		{LicenseTypeID: "PID001", Name: "test1", Type: "host", LicenseCount: 3, OriginalCount: 3},
		{LicenseTypeID: "PID001", Name: "pluto", Type: "host", LicenseCount: 1.5, OriginalCount: 1.5},
		{LicenseTypeID: "PID001", Name: "pippo", Type: "host", LicenseCount: 5.5, OriginalCount: 5.5},

		{LicenseTypeID: "PID002", Name: "topolino", Type: "host", LicenseCount: 7, OriginalCount: 7},
		{LicenseTypeID: "PID002", Name: "minnie", Type: "host", LicenseCount: 3, OriginalCount: 3},

		{LicenseTypeID: "PID003", Name: "minnie", Type: "host", LicenseCount: 0.5, OriginalCount: 0.5},
		{LicenseTypeID: "PID003", Name: "pippo", Type: "host", LicenseCount: 0.5, OriginalCount: 0.5},
		{LicenseTypeID: "PID003", Name: "test2", Type: "host", LicenseCount: 4, OriginalCount: 4},
		{LicenseTypeID: "PID003", Name: "test3", Type: "host", LicenseCount: 6, OriginalCount: 6},

		{LicenseTypeID: "PID005", Name: "test5", Type: "host", LicenseCount: 12, OriginalCount: 12},
	}

	var expectedLicenseTypes = []model.OracleDatabaseLicenseType{
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
		{
			ID:              "PID004",
			ItemDescription: "itemDesc4",
			Aliases:         []string{"alias4"},
			Metric:          model.LicenseTypeMetricComputerPerpetual,
		},
		{
			ID:              "PID005",
			ItemDescription: "itemDesc5",
			Aliases:         []string{"alias5"},
			Metric:          model.LicenseTypeMetricComputerPerpetual,
		},
	}
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

	gomock.InOrder(
		db.EXPECT().
			ListOracleDatabaseAgreements().
			Return(sampleAgreements, nil),
		db.EXPECT().
			ListHostUsingOracleDatabaseLicenses().
			Return(sampleHosts, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(expectedLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(expectedLicenseTypes, nil),
	)

	db.EXPECT().GetHostDatas(utils.MAX_TIME).
		Return(hostdatas, nil)

	actual, err := as.GetOracleDatabaseLicensesCompliance()
	require.NoError(t, err)

	expected := []dto.LicenseCompliance{
		{LicenseTypeID: "PID001", ItemDescription: "itemDesc1", Metric: "Processor Perpetual", Consumed: 10, Covered: 0, Purchased: 10, Compliance: 0, Available: 3, Unlimited: false},
		{LicenseTypeID: "PID002", ItemDescription: "itemDesc2", Metric: "Named User Plus Perpetual", Consumed: 250, Covered: 0, Purchased: 500, Compliance: 0, Available: 250, Unlimited: false},
		{LicenseTypeID: "PID003", ItemDescription: "itemDesc3", Metric: "Computer Perpetual", Consumed: 11, Covered: 0, Purchased: 0, Compliance: 1, Available: 0, Unlimited: true},
		{LicenseTypeID: "PID004", ItemDescription: "itemDesc4", Metric: "Computer Perpetual", Consumed: 0.0, Covered: 0.0, Purchased: 40, Compliance: 1, Available: 40, Unlimited: false},
		{LicenseTypeID: "PID005", ItemDescription: "itemDesc5", Metric: "Computer Perpetual", Consumed: 12, Covered: 0.0, Purchased: 0, Compliance: 0, Available: 0, Unlimited: false},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestGetLicensesCompliance_Veritas(t *testing.T) {
	var sampleAgreements []dto.OracleDatabaseAgreementFE = []dto.OracleDatabaseAgreementFE{
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
			LicensesPerUser:          5,
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
			LicensesPerUser:          100,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 75,
		},
		{
			ID:                       utils.Str2oid("cccccccccccccccccccccccc"),
			AgreementID:              "",
			CSI:                      "",
			LicenseTypeID:            "PID003",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "topolino"}, {Hostname: "minnie"}},
			LicensesPerCore:          0,
			LicensesPerUser:          10,
			AvailableLicensesPerCore: 75,
			AvailableLicensesPerUser: 0,
		},
		{
			ID:                       utils.Str2oid("dddddddddddddddddddddddd"),
			AgreementID:              "",
			CSI:                      "",
			LicenseTypeID:            "PID004",
			ItemDescription:          "",
			Metric:                   "",
			ReferenceNumber:          "",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          5,
			AvailableLicensesPerCore: 50,
			AvailableLicensesPerUser: 0,
		},
	}

	var sampleHosts []dto.HostUsingOracleDatabaseLicenses = []dto.HostUsingOracleDatabaseLicenses{
		{LicenseTypeID: "PID001", Name: "test1", Type: "host", LicenseCount: 2, OriginalCount: 2},
		{LicenseTypeID: "PID001", Name: "test2", Type: "host", LicenseCount: 2, OriginalCount: 2},
		{LicenseTypeID: "PID001", Name: "test3", Type: "host", LicenseCount: 2, OriginalCount: 2},

		{LicenseTypeID: "PID002", Name: "test1", Type: "host", LicenseCount: 2, OriginalCount: 2},

		{LicenseTypeID: "PID003", Name: "test1", Type: "host", LicenseCount: 2, OriginalCount: 2},

		{LicenseTypeID: "PID004", Name: "test3", Type: "host", LicenseCount: 2, OriginalCount: 2},
	}

	var expectedLicenseTypes = []model.OracleDatabaseLicenseType{
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
		{
			ID:              "PID004",
			ItemDescription: "itemDesc4",
			Aliases:         []string{"alias4"},
			Metric:          model.LicenseTypeMetricComputerPerpetual,
		},
		{
			ID:              "PID005",
			ItemDescription: "itemDesc5",
			Aliases:         []string{"alias5"},
			Metric:          model.LicenseTypeMetricComputerPerpetual,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      logger.NewLogger("TEST"),
	}

	hostdatas := []model.HostDataBE{
		{
			Hostname: "test1",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    true,
				VeritasClusterHostnames: []string{"test1", "test2", "test3"},
			},
			Info: model.Host{
				CPUCores: 2,
			},
		},
		{
			Hostname: "test2",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    true,
				VeritasClusterHostnames: []string{"test1", "test2", "test3"},
			},
			Info: model.Host{
				CPUCores: 2,
			},
		},
		{
			Hostname: "test3",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    true,
				VeritasClusterHostnames: []string{"test1", "test2", "test3"},
			},
			Info: model.Host{
				CPUCores: 2,
			},
		},
	}

	gomock.InOrder(
		db.EXPECT().
			ListOracleDatabaseAgreements().
			Return(sampleAgreements, nil),
		db.EXPECT().
			ListHostUsingOracleDatabaseLicenses().
			Return(sampleHosts, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(expectedLicenseTypes, nil),
		db.EXPECT().
			GetOracleDatabaseLicenseTypes().
			Times(1).
			Return(expectedLicenseTypes, nil),
	)

	db.EXPECT().GetHostDatas(utils.MAX_TIME).
		Return(hostdatas, nil)

	actual, err := as.GetOracleDatabaseLicensesCompliance()
	require.NoError(t, err)

	expected := []dto.LicenseCompliance{
		{LicenseTypeID: "PID001", ItemDescription: "itemDesc1", Metric: "Processor Perpetual", Consumed: 3, Covered: 0, Purchased: 5, Compliance: 0, Unlimited: false, Available: 50},
		{LicenseTypeID: "PID002", ItemDescription: "itemDesc2", Metric: "Named User Plus Perpetual", Consumed: 75, Covered: 0, Purchased: 100, Compliance: 0, Unlimited: false, Available: 75},
		{LicenseTypeID: "PID003", ItemDescription: "itemDesc3", Metric: "Computer Perpetual", Consumed: 3, Covered: 0, Purchased: 10, Compliance: 1, Unlimited: true, Available: 75},
		{LicenseTypeID: "PID004", ItemDescription: "itemDesc4", Metric: "Computer Perpetual", Consumed: 3, Covered: 0, Purchased: 5, Compliance: 0, Unlimited: false, Available: 50},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestAddOracleDatabaseLicenseTypes_Success_InsertNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	licenseType := model.OracleDatabaseLicenseType{
		ID:              "Test",
		ItemDescription: "Oracle Database Enterprise Edition",
		Metric:          "Processor Perpetual",
		Cost:            500,
		Aliases:         []string{"Tuning Pack"},
		Option:          false,
	}

	expectedLT := licenseType

	db.EXPECT().InsertOracleDatabaseLicenseType(expectedLT).Return(nil)

	searchedLT := model.OracleDatabaseLicenseType{
		ID:              expectedLT.ID,
		ItemDescription: licenseType.ItemDescription,
		Metric:          licenseType.Metric,
		Cost:            licenseType.Cost,
		Aliases:         licenseType.Aliases,
		Option:          licenseType.Option,
	}

	res, err := as.AddOracleDatabaseLicenseType(licenseType)
	require.NoError(t, err)
	assert.Equal(t,
		searchedLT,
		*res)

}

func TestUpdateOracleDatabaseLicenseTypes_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	licenseType := model.OracleDatabaseLicenseType{
		ID:              "Test",
		ItemDescription: "Oracle Database Enterprise Edition",
		Metric:          "Processor Perpetual",
		Cost:            500,
		Aliases:         []string{"Tuning Pack"},
		Option:          false,
	}

	db.EXPECT().UpdateOracleDatabaseLicenseType(licenseType).Return(nil)

	searchedLTItem := model.OracleDatabaseLicenseType{
		ID:              licenseType.ID,
		ItemDescription: licenseType.ItemDescription,
		Metric:          licenseType.Metric,
		Cost:            licenseType.Cost,
		Aliases:         licenseType.Aliases,
		Option:          licenseType.Option,
	}

	actualLT, err := as.UpdateOracleDatabaseLicenseType(licenseType)
	require.NoError(t, err)
	assert.Equal(t, searchedLTItem, *actualLT)
}

func TestDeleteOracleDatabaseLicenseTypes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	t.Run("Success", func(t *testing.T) {
		LT := model.OracleDatabaseLicenseType{
			ID:              "Test",
			ItemDescription: "Oracle Database Enterprise Edition",
			Metric:          "Processor Perpetual",
			Cost:            500,
			Aliases:         []string{"Tuning Pack"},
			Option:          false,
		}

		gomock.InOrder(
			db.EXPECT().RemoveOracleDatabaseLicenseType(LT.ID).
				Return(nil),
		)

		err := as.DeleteOracleDatabaseLicenseType(LT.ID)
		assert.Nil(t, err)
	})
}
