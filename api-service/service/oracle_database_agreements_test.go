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

var licenseTypesSample = []model.OracleDatabaseLicenseType{
	{
		ID:              "PID001",
		ItemDescription: "itemDesc1",
		Aliases:         []string{"alias1"},
		Metric:          "metric1",
	},
	{
		ID:              "PID002",
		ItemDescription: "itemDesc2",
		Aliases:         []string{"alias2"},
		Metric:          "metric2",
	},
	{
		ID:              "PID003",
		ItemDescription: "itemDesc3",
		Aliases:         []string{"alias3"},
		Metric:          "metric3",
	},
}

func TestAddOracleDatabaseAgreement_Success_InsertNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreement := model.OracleDatabaseAgreement{
		AgreementID:     "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Restricted:      false,
		Hosts:           []string{"test-db", "ercsoldbx"},
	}

	expectedAgr := agreement
	expectedAgr.ID = utils.Str2oid("000000000000000000000001")

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames",
			dto.SearchHostsFilters{
				Search:         []string{""},
				OlderThan:      utils.MAX_TIME,
				PageNumber:     -1,
				PageSize:       -1,
				LTEMemoryTotal: -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
				GTESwapTotal:   -1,
				LTECPUCores:    -1,
				GTECPUCores:    -1,
				LTECPUThreads:  -1,
				GTECPUThreads:  -1,
			}).Return([]map[string]interface{}{
			{"hostname": "test-db"},
			{"hostname": "foobar"},
			{"hostname": "ercsoldbx"},
		}, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypesSample, nil),
		db.EXPECT().InsertOracleDatabaseAgreement(expectedAgr).
			Return(nil),
	)

	searchedAgreementItem := dto.OracleDatabaseAgreementFE{
		ID:                       expectedAgr.ID,
		AgreementID:              agreement.AgreementID,
		CSI:                      agreement.CSI,
		LicenseTypeID:            agreement.LicenseTypeID,
		ItemDescription:          "",
		Metric:                   "",
		ReferenceNumber:          "",
		Unlimited:                false,
		CatchAll:                 false,
		Restricted:               false,
		Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
		LicensesPerCore:          0,
		LicensesPerUser:          0,
		AvailableLicensesPerCore: 0,
		AvailableLicensesPerUser: 0,
	}
	as.mockGetOracleDatabaseAgreements = func(filters dto.GetOracleDatabaseAgreementsFilter) ([]dto.OracleDatabaseAgreementFE, error) {
		return []dto.OracleDatabaseAgreementFE{searchedAgreementItem}, nil
	}

	res, err := as.AddOracleDatabaseAgreement(agreement)
	require.NoError(t, err)
	assert.Equal(t,
		searchedAgreementItem,
		*res)
}

func TestAddOracleDatabaseAgreements_Fail(t *testing.T) {
	t.Run("Fail: can't find host", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		db := NewMockMongoDatabaseInterface(mockCtrl)

		as := APIService{
			Database: db,
			Config: config.Configuration{
				ResourceFilePath: "../../resources",
			},
			TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			NewObjectID: utils.NewObjectIDForTests(),
		}

		addRequest := model.OracleDatabaseAgreement{
			AgreementID:     "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Restricted:      false,
			Hosts: []string{
				"test-db",
				"ercsoldbx",
			},
		}

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				dto.SearchHostsFilters{
					Search:         []string{""},
					OlderThan:      utils.MAX_TIME,
					PageNumber:     -1,
					PageSize:       -1,
					LTEMemoryTotal: -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
					GTESwapTotal:   -1,
					LTECPUCores:    -1,
					GTECPUCores:    -1,
					LTECPUThreads:  -1,
					GTECPUThreads:  -1,
				}).
				Return([]map[string]interface{}{
					{"hostname": "paperino"},
					{"hostname": "pippo"},
					{"hostname": "pluto"},
				}, nil),
		)

		res, err := as.AddOracleDatabaseAgreement(addRequest)
		assert.EqualError(t, err, utils.ErrHostNotFound.Error())
		assert.Nil(t, res)
	})

	t.Run("Fail: can't find licenseType", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		db := NewMockMongoDatabaseInterface(mockCtrl)

		as := APIService{
			Database: db,
			Config: config.Configuration{
				ResourceFilePath: "../../resources",
			},
			TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			NewObjectID: utils.NewObjectIDForTests(),
		}

		agreementWrongLicenseType := model.OracleDatabaseAgreement{
			AgreementID:     "AID001",
			LicenseTypeID:   "xxxxxx",
			CSI:             "CSI001",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Restricted:      false,
			Hosts: []string{
				"test-db",
				"ercsoldbx",
			},
		}
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				dto.SearchHostsFilters{
					Search:         []string{""},
					OlderThan:      utils.MAX_TIME,
					PageNumber:     -1,
					PageSize:       -1,
					LTEMemoryTotal: -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
					GTESwapTotal:   -1,
					LTECPUCores:    -1,
					GTECPUCores:    -1,
					LTECPUThreads:  -1,
					GTECPUThreads:  -1,
				}).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseLicenseTypes().
				Return(licenseTypesSample, nil),
		)

		res, err := as.AddOracleDatabaseAgreement(agreementWrongLicenseType)

		assert.EqualError(t, err, utils.ErrOracleDatabaseLicenseTypeIDNotFound.Error())
		assert.Nil(t, res)
	})
}

func TestUpdateOracleDatabaseAgreement_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreement := model.OracleDatabaseAgreement{
		AgreementID:     "AID001",
		CSI:             "CSI001",
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		LicenseTypeID:   licenseTypesSample[0].ID,
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Restricted:      false,
		Hosts:           []string{"test-db", "ercsoldbx"},
	}

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames",
			dto.SearchHostsFilters{
				Search:         []string{""},
				OlderThan:      utils.MAX_TIME,
				PageNumber:     -1,
				PageSize:       -1,
				LTEMemoryTotal: -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
				GTESwapTotal:   -1,
				LTECPUCores:    -1,
				GTECPUCores:    -1,
				LTECPUThreads:  -1,
				GTECPUThreads:  -1,
			}).
			Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypesSample, nil),
		db.EXPECT().UpdateOracleDatabaseAgreement(agreement).Return(nil),
	)

	searchedAgreementItem := dto.OracleDatabaseAgreementFE{
		ID:                       agreement.ID,
		AgreementID:              agreement.AgreementID,
		CSI:                      agreement.CSI,
		LicenseTypeID:            agreement.LicenseTypeID,
		ItemDescription:          "",
		Metric:                   "",
		ReferenceNumber:          "",
		Unlimited:                false,
		CatchAll:                 false,
		Restricted:               false,
		Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
		LicensesPerCore:          0,
		LicensesPerUser:          0,
		AvailableLicensesPerCore: 0,
		AvailableLicensesPerUser: 0,
	}
	as.mockGetOracleDatabaseAgreements = func(filters dto.GetOracleDatabaseAgreementsFilter) ([]dto.OracleDatabaseAgreementFE, error) {
		return []dto.OracleDatabaseAgreementFE{searchedAgreementItem}, nil
	}

	actualAgreement, err := as.UpdateOracleDatabaseAgreement(agreement)
	require.NoError(t, err)
	assert.Equal(t, searchedAgreementItem, *actualAgreement)
}

func TestUpdateOracleDatabaseAgreement_Fail_LicenseTypeIdNotValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreement := model.OracleDatabaseAgreement{
		AgreementID:     "AID001",
		CSI:             "CSI001",
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		LicenseTypeID:   "invalidLicenseTypeID",
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Restricted:      false,
		Hosts:           []string{"test-db", "ercsoldbx"},
	}

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames",
			dto.SearchHostsFilters{
				Search:         []string{""},
				OlderThan:      utils.MAX_TIME,
				PageNumber:     -1,
				PageSize:       -1,
				LTEMemoryTotal: -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
				GTESwapTotal:   -1,
				LTECPUCores:    -1,
				GTECPUCores:    -1,
				LTECPUThreads:  -1,
				GTECPUThreads:  -1,
			}).
			Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypesSample, nil),
	)

	actual, err := as.UpdateOracleDatabaseAgreement(agreement)

	assert.EqualError(t, err, utils.ErrOracleDatabaseLicenseTypeIDNotFound.Error())
	assert.Nil(t, actual)
}

func TestGetOracleDatabaseAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
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
			Unlimited:                true,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	returnedHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseAgreements().
			Return(returnedAgreements, nil),
		db.EXPECT().ListHostUsingOracleDatabaseLicenses().
			Return(returnedHosts, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
	)

	res, err := as.GetOracleDatabaseAgreements(dto.NewGetOracleDatabaseAgreementsFilter())
	require.NoError(t, err)
	assert.Equal(t, expectedAgreements, res)
}

func TestGetOracleDatabaseAgreements_SuccessFilter1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
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
			Unlimited:                true,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 7,
			AvailableLicensesPerUser: 0,
		},
	}
	returnedHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "ID Partioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseAgreements().
			Return(returnedAgreements, nil),
		db.EXPECT().ListHostUsingOracleDatabaseLicenses().
			Return(returnedHosts, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(parts, nil),
	)

	res, err := as.GetOracleDatabaseAgreements(dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "asddfa",
		Unlimited:                   "",
		CatchAll:                    "",
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerCoreLTE: -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerCoreLTE:          -1,
		LicensesPerUserGTE:          -1,
		LicensesPerUserLTE:          -1,
	})
	require.NoError(t, err)
	assert.Empty(t, res)

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseAgreements().
			Return(returnedAgreements, nil),
		db.EXPECT().ListHostUsingOracleDatabaseLicenses().
			Return(returnedHosts, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(parts, nil),
	)

	res, err = as.GetOracleDatabaseAgreements(dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "asddfa",
		Unlimited:                   "",
		CatchAll:                    "",
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerCoreLTE: -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerCoreLTE:          -1,
		LicensesPerUserGTE:          -1,
		LicensesPerUserLTE:          -1,
	})
	require.NoError(t, err)
	assert.Empty(t, res)

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseAgreements().
			Return(returnedAgreements, nil),
		db.EXPECT().ListHostUsingOracleDatabaseLicenses().
			Return(returnedHosts, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(parts, nil),
	)

	res, err = as.GetOracleDatabaseAgreements(dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "asddfa",
		Unlimited:                   "",
		CatchAll:                    "",
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerCoreLTE: -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerCoreLTE:          -1,
		LicensesPerUserGTE:          -1,
		LicensesPerUserLTE:          -1,
	})

	require.NoError(t, err)
	assert.Empty(t, res)
}

func TestGetOracleDatabaseAgreements_Failed2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
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
			Unlimited:                true,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 7,
			AvailableLicensesPerUser: 0,
		},
	}

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseAgreements().
			Return(returnedAgreements, nil),
		db.EXPECT().ListHostUsingOracleDatabaseLicenses().
			Return(nil, aerrMock),
	)

	_, err := as.GetOracleDatabaseAgreements(dto.NewGetOracleDatabaseAgreementsFilter())
	require.Equal(t, aerrMock, err)
}

func TestCheckOracleDatabaseAgreementMatchFilter(t *testing.T) {
	agg1 := dto.OracleDatabaseAgreementFE{
		ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
		AgreementID:              "5051863",
		CSI:                      "6871235",
		LicenseTypeID:            "A90620",
		ItemDescription:          "Oracle Partitioning",
		Metric:                   model.LicenseTypeMetricProcessorPerpetual,
		ReferenceNumber:          "10032246681",
		Unlimited:                false,
		CatchAll:                 true,
		Restricted:               false,
		Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: -1, Hostname: "test-db", TotalCoveredLicensesCount: -1}, {CoveredLicensesCount: -1, Hostname: "ercsoldbx", TotalCoveredLicensesCount: -1}},
		LicensesPerCore:          30,
		LicensesPerUser:          5,
		AvailableLicensesPerCore: 7,
		AvailableLicensesPerUser: 0,
	}

	assert.True(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.NewGetOracleDatabaseAgreementsFilter()))

	assert.True(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "5051",
		LicenseTypeID:               "A9062",
		ItemDescription:             "Partitioning",
		CSI:                         "6871",
		Metric:                      model.LicenseTypeMetricProcessorPerpetual,
		ReferenceNumber:             "100322",
		Unlimited:                   "false",
		CatchAll:                    "true",
		LicensesPerCoreLTE:          35,
		LicensesPerCoreGTE:          25,
		LicensesPerUserLTE:          10,
		LicensesPerUserGTE:          0,
		AvailableLicensesPerCoreLTE: 8,
		AvailableLicensesPerCoreGTE: 6,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.True(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          30,
		LicensesPerCoreGTE:          30,
		LicensesPerUserLTE:          5,
		LicensesPerUserGTE:          5,
		AvailableLicensesPerCoreLTE: 7,
		AvailableLicensesPerCoreGTE: 7,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))

	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "fdgdfgsdsfg",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "fdgdfgsdsfg",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "fdgdfgsdsfg",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "fdgdfgsdsfg",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "fdgdfgsdsfg",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "fdgdfgsdsfg",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "true",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "false",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          35,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          25,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          0,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          10,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: 3,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: 8,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleUnlimitedCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "test-db", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 7,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleProcessorPerpetualCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 2,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleNamedUserPlusCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          250,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 250,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  5,
			LicenseTypeID: "PID002",
			OriginalCount: 5,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 125, Hostname: "test-db", TotalCoveredLicensesCount: 125, ConsumedLicensesCount: 125}},
			LicensesPerCore:          0,
			LicensesPerUser:          250,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 125,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SharedAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}, {CoveredLicensesCount: 0, Hostname: "test-db2", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "test-db2",
			LicenseCount:  4,
			LicenseTypeID: "PID002",
			OriginalCount: 4,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 4, Hostname: "test-db2", TotalCoveredLicensesCount: 4, ConsumedLicensesCount: 4}, {CoveredLicensesCount: 1, Hostname: "test-db", TotalCoveredLicensesCount: 1, ConsumedLicensesCount: 3}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: -2,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SharedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 10,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  20,
			LicenseTypeID: "PID002",
			OriginalCount: 20,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 10, Hostname: "test-db", TotalCoveredLicensesCount: 15, ConsumedLicensesCount: 20}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: -5,
			AvailableLicensesPerUser: 0,
		},
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 5, Hostname: "test-db", TotalCoveredLicensesCount: 15, ConsumedLicensesCount: 20}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: -5,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleUnlimitedCaseNoAssociatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "ID Partioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleProcessorPerpetualCaseNoAssociatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 2,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleNamedUserPlusCaseNoAssociatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          200,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 200,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  5,
			LicenseTypeID: "PID002",
			OriginalCount: 5,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          200,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 75,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_CompleCase1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "test-db"}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 10,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "dbclust",
			LicenseCount:  20,
			LicenseTypeID: "PID002",
			OriginalCount: 20,
			Type:          "cluster",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			CatchAll:                 true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "test-db", CoveredLicensesCount: 3, TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: -13,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestSortHostsUsingLicenses(t *testing.T) {
	list := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "Real Application Clusters",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
		},
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "test-db4",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
	}

	expected := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "Real Application Clusters",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
		},
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "test-db4",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
		},
	}

	sortHostsByLicenses(list)

	assert.Equal(t, expected, list)
}

func TestSortOracleDatabaseAgreements(t *testing.T) {
	list := []dto.OracleDatabaseAgreementFE{
		{CatchAll: true, Unlimited: false, LicensesPerUser: 10},
		{CatchAll: true, Unlimited: false, LicensesPerCore: 10},
		{CatchAll: true, Unlimited: true, LicensesPerUser: 20},
		{CatchAll: false, Unlimited: false, LicensesPerCore: 20},
		{CatchAll: false, Unlimited: true, LicensesPerUser: 10},
		{CatchAll: false, Unlimited: true, LicensesPerCore: 20},
		{CatchAll: false, Unlimited: false, LicensesPerCore: 10},
		{CatchAll: true, Unlimited: true, LicensesPerCore: 10},
		{CatchAll: false, Unlimited: true, LicensesPerUser: 20},
		{CatchAll: false, Unlimited: false, LicensesPerUser: 10},
		{CatchAll: true, Unlimited: true, LicensesPerUser: 10},
		{CatchAll: true, Unlimited: true, LicensesPerCore: 20},
		{CatchAll: true, Unlimited: false, LicensesPerCore: 20},
		{CatchAll: false, Unlimited: false, LicensesPerUser: 20},
		{CatchAll: false, Unlimited: true, LicensesPerCore: 10},
		{CatchAll: true, Unlimited: false, LicensesPerUser: 20},
	}

	expected := []dto.OracleDatabaseAgreementFE{
		{CatchAll: false, Unlimited: false, LicensesPerUser: 20},
		{CatchAll: false, Unlimited: false, LicensesPerUser: 10},
		{CatchAll: false, Unlimited: false, LicensesPerCore: 20},
		{CatchAll: false, Unlimited: false, LicensesPerCore: 10},
		{CatchAll: false, Unlimited: true, LicensesPerUser: 20},
		{CatchAll: false, Unlimited: true, LicensesPerUser: 10},
		{CatchAll: false, Unlimited: true, LicensesPerCore: 20},
		{CatchAll: false, Unlimited: true, LicensesPerCore: 10},
		{CatchAll: true, Unlimited: false, LicensesPerUser: 20},
		{CatchAll: true, Unlimited: false, LicensesPerUser: 10},
		{CatchAll: true, Unlimited: false, LicensesPerCore: 20},
		{CatchAll: true, Unlimited: false, LicensesPerCore: 10},
		{CatchAll: true, Unlimited: true, LicensesPerUser: 20},
		{CatchAll: true, Unlimited: true, LicensesPerUser: 10},
		{CatchAll: true, Unlimited: true, LicensesPerCore: 20},
		{CatchAll: true, Unlimited: true, LicensesPerCore: 10},
	}

	sortOracleDatabaseAgreements(list)

	assert.Equal(t, expected, list)
}

func TestSortAssociatedHostsInOracleDatabaseAgreement(t *testing.T) {
	partsMap := map[string]*model.OracleDatabaseLicenseType{
		"L10005": {
			ID:              "L10005",
			ItemDescription: "Oracle Real Application Clusters",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
			Aliases:         []string{"Real Application Clusters", "RAC or RAC One Node"},
		},
	}

	hostsMap := map[string]map[string]*dto.HostUsingOracleDatabaseLicenses{
		"L10005": {
			"test-db1": {
				LicenseCount: 30,
			},
			"test-db2": {
				LicenseCount: 30,
			},
			"test-db3": {
				LicenseCount: 15,
			},
			"test-db4": {
				LicenseCount: 35,
			},
		},
	}

	agr := dto.OracleDatabaseAgreementFE{
		LicenseTypeID: "L10005",
		Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
			{Hostname: "test-db2"},
			{Hostname: "test-db1"},
			{Hostname: "test-db4"},
			{Hostname: "test-db3"},
		},
	}

	expected := []dto.OracleDatabaseAgreementAssociatedHostFE{
		{Hostname: "test-db4"},
		{Hostname: "test-db2"},
		{Hostname: "test-db1"},
		{Hostname: "test-db3"},
	}

	sortHostsInAgreementByLicenseCount(&agr, hostsMap, partsMap)

	assert.Equal(t, expected, agr.Hosts)
}

func TestBuildHostUsingLicensesMap(t *testing.T) {
	list := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "LTID01",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "LTID02",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "LTID03",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
		},
		{
			LicenseTypeID: "LTID02",
			Name:          "test-db4",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "LTID01",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
		},
	}

	expected := map[string]map[string]*dto.HostUsingOracleDatabaseLicenses{
		"LTID01": {
			"Puzzait":  &list[0],
			"test-db3": &list[4],
		},
		"LTID02": {
			"Puzzait":  &list[1],
			"test-db4": &list[3],
		},
		"LTID03": {
			"test-db3": &list[2],
		},
	}

	assert.Equal(t, expected, buildHostUsingLicensesMap(list))
}

func TestBuildAgreementPartMap(t *testing.T) {
	list := []model.OracleDatabaseLicenseType{
		{
			ItemDescription: "itemDesc1",
			Aliases:         []string{"alias1"},
			Metric:          "metric1",
			ID:              "PID001",
		},
		{
			ItemDescription: "itemDesc2",
			Aliases:         []string{"alias1"},
			Metric:          "metric2",
			ID:              "PID002",
		},
	}

	expected := map[string]*model.OracleDatabaseLicenseType{
		"PID001": &list[0],
		"PID002": &list[1],
	}

	assert.Equal(t, expected, buildLicenseTypesMap(list))
}

func TestDeleteOracleDatabaseAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreementID := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Fail: can't find agreement", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().RemoveOracleDatabaseAgreement(agreementID).
				Return(utils.ErrOracleDatabaseAgreementNotFound),
		)

		err := as.DeleteOracleDatabaseAgreement(agreementID)
		require.EqualError(t, err, utils.ErrOracleDatabaseAgreementNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			ID:              agreementID,
			AgreementID:     "AID001",
			CSI:             "CSI001",
			LicenseTypeID:   licenseTypesSample[0].ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Hosts:           []string{"test-db", "ercsoldbx"},
		}

		gomock.InOrder(
			db.EXPECT().RemoveOracleDatabaseAgreement(agreement.ID).
				Return(nil),
		)

		err := as.DeleteOracleDatabaseAgreement(agreementID)
		assert.Nil(t, err)
	})

}

func TestAddHostToOracleDatabaseAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	anotherAssociatedPartID := utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb")
	t.Run("Fail: can't find host", func(t *testing.T) {

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				dto.SearchHostsFilters{
					Search:         []string{""},
					OlderThan:      utils.MAX_TIME,
					PageNumber:     -1,
					PageSize:       -1,
					LTEMemoryTotal: -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
					GTESwapTotal:   -1,
					LTECPUCores:    -1,
					GTECPUCores:    -1,
					LTECPUThreads:  -1,
					GTECPUThreads:  -1,
				}).Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
		)

		err := as.AddHostToOracleDatabaseAgreement(anotherAssociatedPartID, "pippo")
		assert.EqualError(t, err, utils.ErrHostNotFound.Error())
	})

	id := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Fail: can't find agreement", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				dto.SearchHostsFilters{
					Search:         []string{""},
					OlderThan:      utils.MAX_TIME,
					PageNumber:     -1,
					PageSize:       -1,
					LTEMemoryTotal: -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
					GTESwapTotal:   -1,
					LTECPUCores:    -1,
					GTECPUCores:    -1,
					LTECPUThreads:  -1,
					GTECPUThreads:  -1,
				}).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "foobar"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseAgreement(id).
				Return(nil, utils.ErrOracleDatabaseAgreementNotFound),
		)

		err := as.AddHostToOracleDatabaseAgreement(id, "foobar")
		assert.EqualError(t, err, utils.ErrOracleDatabaseAgreementNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			AgreementID:     "AID001",
			CSI:             "CSI001",
			ID:              id,
			LicenseTypeID:   licenseTypesSample[0].ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Hosts:           []string{"test-db", "ercsoldbx"},
		}

		agreementPostAdd := model.OracleDatabaseAgreement{
			ID:              id,
			AgreementID:     "AID001",
			CSI:             "CSI001",
			LicenseTypeID:   licenseTypesSample[0].ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Hosts:           []string{"test-db", "ercsoldbx", "foobar"},
		}

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				dto.SearchHostsFilters{
					Search:         []string{""},
					OlderThan:      utils.MAX_TIME,
					PageNumber:     -1,
					PageSize:       -1,
					LTEMemoryTotal: -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
					GTESwapTotal:   -1,
					LTECPUCores:    -1,
					GTECPUCores:    -1,
					LTECPUThreads:  -1,
					GTECPUThreads:  -1,
				}).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "foobar"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseAgreement(id).
				Return(&agreement, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreementPostAdd).
				Return(nil),
		)

		err := as.AddHostToOracleDatabaseAgreement(id, "foobar")
		assert.Nil(t, err)
	})
}

func TestDeleteHostFromOracleDatabaseAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	id := utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb")

	t.Run("Fail: can't get agreement", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				dto.SearchHostsFilters{
					Search:         []string{""},
					OlderThan:      utils.MAX_TIME,
					PageNumber:     -1,
					PageSize:       -1,
					LTEMemoryTotal: -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
					GTESwapTotal:   -1,
					LTECPUCores:    -1,
					GTECPUCores:    -1,
					LTECPUThreads:  -1,
					GTECPUThreads:  -1,
				}).
				Return([]map[string]interface{}{
					{"hostname": "pippo"},
				}, nil),
			db.EXPECT().GetOracleDatabaseAgreement(id).
				Return(nil, utils.ErrOracleDatabaseAgreementNotFound),
		)

		err := as.DeleteHostFromOracleDatabaseAgreement(id, "pippo")
		require.EqualError(t, err, utils.ErrOracleDatabaseAgreementNotFound.Error())
	})

	anotherId := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Success", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			AgreementID:     "AID001",
			CSI:             "CSI001",
			ID:              anotherId,
			LicenseTypeID:   licenseTypesSample[0].ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Hosts:           []string{"test-db", "ercsoldbx"},
		}

		agreementPostAdd := model.OracleDatabaseAgreement{
			AgreementID:     "AID001",
			CSI:             "CSI001",
			ID:              anotherId,
			LicenseTypeID:   licenseTypesSample[0].ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Hosts:           []string{"test-db"},
		}

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				dto.SearchHostsFilters{
					Search:         []string{""},
					OlderThan:      utils.MAX_TIME,
					PageNumber:     -1,
					PageSize:       -1,
					LTEMemoryTotal: -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
					GTESwapTotal:   -1,
					LTECPUCores:    -1,
					GTECPUCores:    -1,
					LTECPUThreads:  -1,
					GTECPUThreads:  -1,
				}).
				Return([]map[string]interface{}{
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseAgreement(anotherId).
				Return(&agreement, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreementPostAdd).
				Return(nil),
		)

		err := as.DeleteHostFromOracleDatabaseAgreement(anotherId, "ercsoldbx")
		assert.Nil(t, err)
	})
}
