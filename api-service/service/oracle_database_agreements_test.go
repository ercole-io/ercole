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

	"github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/api-service/dto"
	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO add tests for UpdateOracleDatabaseAgreement

var partsSample = []model.OracleDatabasePart{
	{
		PartID:          "PID001",
		ItemDescription: "itemDesc1",
		Aliases:         []string{"alias1"},
		Metric:          "metric1",
	},
	{
		PartID:          "PID002",
		ItemDescription: "itemDesc2",
		Aliases:         []string{"alias2"},
		Metric:          "metric2",
	},
	{
		PartID:          "PID003",
		ItemDescription: "itemDesc3",
		Aliases:         []string{"alias3"},
		Metric:          "metric3",
	},
}

func TestAddOracleDatabaseAgreements_Success_InsertNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: partsSample,
		TimeNow:                      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID:                  utils.NewObjectIDForTests(),
	}
	addRequest := dto.AssociatedPartInOracleDbAgreementRequest{
		AgreementID:     "AID001",
		PartID:          "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Hosts: []string{
			"test-db",
			"ercsoldbx",
		},
	}

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
			GTECPUCores:    -1,
			LTECPUCores:    -1,
			LTECPUThreads:  -1,
			LTEMemoryTotal: -1,
			GTECPUThreads:  -1,
			GTESwapTotal:   -1,
			GTEMemoryTotal: -1,
			LTESwapTotal:   -1,
		}, "", false, -1, -1, "", "", utils.MAX_TIME).Return([]map[string]interface{}{
			{"hostname": "test-db"},
			{"hostname": "foobar"},
			{"hostname": "ercsoldbx"},
		}, nil),
		db.EXPECT().GetOracleDatabaseAgreement(addRequest.AgreementID).Return(nil, utils.AerrOracleDatabaseAgreementNotFound),
		db.EXPECT().InsertOracleDatabaseAgreement(gomock.Any()).Do(func(actual model.OracleDatabaseAgreement) {
			actual.ID = primitive.NilObjectID

			expected := model.OracleDatabaseAgreement{
				AgreementID: "AID001",
				CSI:         "CSI001",
				Parts: []model.AssociatedPart{
					{
						ID:                 utils.Str2oid("000000000000000000000001"),
						OracleDatabasePart: partsSample[0],
						ReferenceNumber:    "RF0001",
						Unlimited:          true,
						Count:              30,
						CatchAll:           true,
						Hosts:              []string{"test-db", "ercsoldbx"},
					},
				},
			}

			assert.Equal(t, expected, actual)
		}).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5f4d0a2b27fe53da8a4aec45")}, nil),
	)

	res, err := as.AddAssociatedPartToOracleDbAgreement(addRequest)
	require.NoError(t, err)
	assert.Equal(t,
		"5f4d0a2b27fe53da8a4aec45",
		res)
}

func TestAddOracleDatabaseAgreements_Success_AlreadyExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: partsSample,
		TimeNow:                      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID:                  utils.NewObjectIDForTests(),
	}
	addRequest := dto.AssociatedPartInOracleDbAgreementRequest{
		AgreementID:     "AID001",
		PartID:          "PID002",
		CSI:             "CSI001",
		ReferenceNumber: "RF0002",
		Unlimited:       false,
		Count:           33,
		CatchAll:        true,
		Hosts:           []string{"pippo", "pluto"},
	}

	alreadyExistsAgreement := model.OracleDatabaseAgreement{
		ID:          utils.Str2oid("5f4d0a2b27fe53da8a4aec45"),
		AgreementID: "AID001",
		CSI:         "CSI001",
		Parts: []model.AssociatedPart{
			{
				ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				OracleDatabasePart: partsSample[0],
				ReferenceNumber:    "RF0001",
				Unlimited:          true,
				Count:              30,
				CatchAll:           false,
				Hosts:              []string{"test-db", "ercsoldbx"},
			},
		},
	}

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
			GTECPUCores:    -1,
			LTECPUCores:    -1,
			LTECPUThreads:  -1,
			LTEMemoryTotal: -1,
			GTECPUThreads:  -1,
			GTESwapTotal:   -1,
			GTEMemoryTotal: -1,
			LTESwapTotal:   -1,
		}, "", false, -1, -1, "", "", utils.MAX_TIME).Return([]map[string]interface{}{
			{"hostname": "pippo"},
			{"hostname": "pluto"},
		}, nil),
		db.EXPECT().GetOracleDatabaseAgreement(addRequest.AgreementID).Return(&alreadyExistsAgreement, nil),
		db.EXPECT().UpdateOracleDatabaseAgreement(gomock.Any()).Do(func(actual model.OracleDatabaseAgreement) {
			expected := model.OracleDatabaseAgreement{
				ID:          utils.Str2oid("5f4d0a2b27fe53da8a4aec45"),
				AgreementID: "AID001",
				CSI:         "CSI001",
				Parts: []model.AssociatedPart{
					{
						ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
						OracleDatabasePart: partsSample[0],
						ReferenceNumber:    "RF0001",
						Unlimited:          true,
						Count:              30,
						CatchAll:           false,
						Hosts:              []string{"test-db", "ercsoldbx"},
					},
					{
						ID:                 utils.Str2oid("000000000000000000000001"),
						OracleDatabasePart: partsSample[1],
						ReferenceNumber:    "RF0002",
						Unlimited:          false,
						Count:              33,
						CatchAll:           true,
						Hosts:              []string{"pippo", "pluto"},
					},
				},
			}

			assert.Equal(t, expected, actual)
		}).Return(nil),
	)

	res, err := as.AddAssociatedPartToOracleDbAgreement(addRequest)
	require.NoError(t, err)
	assert.Equal(t,
		"5f4d0a2b27fe53da8a4aec45",
		res)
}

func TestAddOracleDatabaseAgreements_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: partsSample,
		TimeNow:                      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID:                  utils.NewObjectIDForTests(),
	}

	addRequest := dto.AssociatedPartInOracleDbAgreementRequest{
		AgreementID:     "AID001",
		PartID:          "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Hosts: []string{
			"test-db",
			"ercsoldbx",
		},
	}

	t.Run("Fail: can't find host", func(t *testing.T) {

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
				GTECPUCores:    -1,
				LTECPUCores:    -1,
				LTECPUThreads:  -1,
				LTEMemoryTotal: -1,
				GTECPUThreads:  -1,
				GTESwapTotal:   -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
			}, "", false, -1, -1, "", "", utils.MAX_TIME).Return([]map[string]interface{}{
				{"hostname": "paperino"},
				{"hostname": "pippo"},
				{"hostname": "pluto"},
			}, nil),
		)

		res, err := as.AddAssociatedPartToOracleDbAgreement(addRequest)
		require.EqualError(t, err, utils.AerrHostNotFound.Error())

		assert.Equal(t, "", res)

	})

	t.Run("Fail: can't find part", func(t *testing.T) {

		addRequestWrongPart := dto.AssociatedPartInOracleDbAgreementRequest{
			AgreementID:     "AID001",
			PartID:          "xxxxxx",
			CSI:             "CSI001",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			CatchAll:        true,
			Hosts: []string{
				"test-db",
				"ercsoldbx",
			},
		}
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
				GTECPUCores:    -1,
				LTECPUCores:    -1,
				LTECPUThreads:  -1,
				LTEMemoryTotal: -1,
				GTECPUThreads:  -1,
				GTESwapTotal:   -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
			}, "", false, -1, -1, "", "", utils.MAX_TIME).Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "ercsoldbx"},
			}, nil),
			db.EXPECT().GetOracleDatabaseAgreement(addRequest.AgreementID).Return(nil, utils.AerrOracleDatabaseAgreementNotFound),
		)

		res, err := as.AddAssociatedPartToOracleDbAgreement(addRequestWrongPart)

		require.EqualError(t, err, utils.AerrOracleDatabaseAgreementInvalidPartID.Error())

		assert.Equal(t, "", res)

	})
}

func TestUpdateAssociatedPartOfOracleDbAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: partsSample,
		TimeNow:                      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID:                  utils.NewObjectIDForTests(),
	}

	agreement := model.OracleDatabaseAgreement{
		AgreementID: "AID001",
		CSI:         "CSI001",
		Parts: []model.AssociatedPart{
			{
				ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				OracleDatabasePart: partsSample[0],
				ReferenceNumber:    "RF0001",
				Unlimited:          true,
				Count:              30,
				CatchAll:           true,
				Hosts:              []string{"test-db", "ercsoldbx"},
			},
		},
	}

	req := dto.AssociatedPartInOracleDbAgreementRequest{
		ID:              "aaaaaaaaaaaaaaaaaaaaaaaa",
		AgreementID:     "AID999",
		PartID:          "PID002",
		CSI:             "CSI999",
		ReferenceNumber: "REFREF",
		Unlimited:       true,
		Count:           42,
		CatchAll:        false,
		Hosts:           []string{"foobar"},
	}

	t.Run("Update successfully", func(t *testing.T) {
		agrForGet := agreement

		agrForUpdate := agreement
		agrForUpdate.AgreementID = req.AgreementID
		agrForUpdate.CSI = req.CSI

		agrPart := &agrForUpdate.Parts[0]
		agrPart.ReferenceNumber = req.AgreementID
		agrPart.Unlimited = req.Unlimited
		agrPart.Count = req.Count
		agrPart.CatchAll = req.CatchAll
		agrPart.Hosts = req.Hosts

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
				GTECPUCores:    -1,
				LTECPUCores:    -1,
				LTECPUThreads:  -1,
				LTEMemoryTotal: -1,
				GTECPUThreads:  -1,
				GTESwapTotal:   -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
			}, "", false, -1, -1, "", "", utils.MAX_TIME).Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(utils.Str2oid(req.ID)).
				Return(&agrForGet, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreement).Return(nil),
		)

		err := as.UpdateAssociatedPartOfOracleDbAgreement(req)
		require.NoError(t, err)
	})

	t.Run("Fail: agreement not found", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames", []string{""},
				database.SearchHostsFilters{
					GTECPUCores:    -1,
					LTECPUCores:    -1,
					LTECPUThreads:  -1,
					LTEMemoryTotal: -1,
					GTECPUThreads:  -1,
					GTESwapTotal:   -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
				}, "", false, -1, -1, "", "", utils.MAX_TIME).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "foobar"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(utils.Str2oid(req.ID)).
				Return(nil, utils.AerrOracleDatabaseAgreementNotFound),
		)

		err := as.UpdateAssociatedPartOfOracleDbAgreement(req)
		assert.EqualError(t, err, utils.AerrOracleDatabaseAgreementNotFound.Error())
	})

	t.Run("Fail: partID not valid", func(t *testing.T) {
		agrForGet := agreement

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames", []string{""},
				database.SearchHostsFilters{
					GTECPUCores:    -1,
					LTECPUCores:    -1,
					LTECPUThreads:  -1,
					LTEMemoryTotal: -1,
					GTECPUThreads:  -1,
					GTESwapTotal:   -1,
					GTEMemoryTotal: -1,
					LTESwapTotal:   -1,
				}, "", false, -1, -1, "", "", utils.MAX_TIME).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "foobar"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(utils.Str2oid(req.ID)).
				Return(&agrForGet, nil),
		)

		req.PartID = "this is a wrong partID"

		err := as.UpdateAssociatedPartOfOracleDbAgreement(req)
		assert.EqualError(t, err, utils.AerrOracleDatabaseAgreementInvalidPartID.Error())
	})
}

func TestSearchOracleDatabaseAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
				// TODO Cost: ,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
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
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	returnedHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 0,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      3,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 3,
					ConsumedLicensesCount:     3,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListHostUsingOracleDatabaseLicenses().Return(returnedHosts, nil)

	res, err := as.SearchAssociatedPartsInOracleDatabaseAgreements(dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	})
	require.NoError(t, err)
	assert.Equal(t, expectedAgreements, res)
}

func TestSearchOracleDatabaseAgreements_SuccessFilter1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	returnedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 7,
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
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	returnedHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListHostUsingOracleDatabaseLicenses().Return(returnedHosts, nil)

	res, err := as.SearchAssociatedPartsInOracleDatabaseAgreements(dto.SearchOracleDatabaseAgreementsFilter{
		AgreementID:       "asddfa",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	})

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListHostUsingOracleDatabaseLicenses().Return(returnedHosts, nil)

	res, err = as.SearchAssociatedPartsInOracleDatabaseAgreements(dto.SearchOracleDatabaseAgreementsFilter{
		AgreementID:       "asddfa",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	})

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListHostUsingOracleDatabaseLicenses().Return(returnedHosts, nil)

	res, err = as.SearchAssociatedPartsInOracleDatabaseAgreements(dto.SearchOracleDatabaseAgreementsFilter{
		AgreementID:       "asddfa",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	})

	require.NoError(t, err)
	assert.Empty(t, res)
}

func TestSearchOracleDatabaseAgreements_Failed2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	returnedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 7,
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
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListHostUsingOracleDatabaseLicenses().Return(nil, aerrMock)

	_, err := as.SearchAssociatedPartsInOracleDatabaseAgreements(dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	})
	require.Equal(t, aerrMock, err)
}

func TestCheckOracleDatabaseAgreementMatchFilter(t *testing.T) {
	agg1 := dto.OracleDatabaseAgreementFE{
		AgreementID:    "5051863",
		AvailableCount: 7,
		CatchAll:       true,
		CSI:            "6871235",
		Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
			{
				CoveredLicensesCount:      -1,
				Hostname:                  "test-db",
				TotalCoveredLicensesCount: -1,
			},
			{
				CoveredLicensesCount:      -1,
				Hostname:                  "ercsoldbx",
				TotalCoveredLicensesCount: -1,
			},
		},
		ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
		ItemDescription: "Oracle Partitioning",
		LicensesCount:   30,
		Metric:          model.AgreementPartMetricProcessorPerpetual,
		PartID:          "A90620",
		ReferenceNumber: "10032246681",
		Unlimited:       false,
		UsersCount:      5,
	}

	assert.True(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))

	assert.True(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		AgreementID:       "5051",
		PartID:            "A9062",
		ItemDescription:   "Partitioning",
		CSI:               "6871",
		Metric:            model.AgreementPartMetricProcessorPerpetual,
		ReferenceNumber:   "100322",
		Unlimited:         "false",
		CatchAll:          "true",
		AvailableCountGTE: 6,
		AvailableCountLTE: 8,
		LicensesCountGTE:  25,
		LicensesCountLTE:  35,
		UsersCountGTE:     0,
		UsersCountLTE:     10,
	}))
	assert.True(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: 7,
		AvailableCountLTE: 7,
		LicensesCountGTE:  30,
		LicensesCountLTE:  30,
		UsersCountGTE:     5,
		UsersCountLTE:     5,
	}))

	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		AgreementID:       "fdgdfgsdsfg",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		PartID:            "fdgdfgsdsfg",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		ItemDescription:   "fdgdfgsdsfg",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		CSI:               "fdgdfgsdsfg",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Metric:            "fdgdfgsdsfg",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		ReferenceNumber:   "fdgdfgsdsfg",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "true",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "false",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  35,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  25,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     0,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     10,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: 3,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, checkOracleDatabaseAgreementMatchFilter(agg1, dto.SearchOracleDatabaseAgreementsFilter{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: 8,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleUnlimitedCase(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					Hostname:                  "test-db",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 0,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      3,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 3,
					ConsumedLicensesCount:     3,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleProcessorPerpetualCase(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 5,
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
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "AID001",
			AvailableCount: 2,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      3,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 3,
					ConsumedLicensesCount:     3,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleNamedUserPlusCase(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 250,
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
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      250,
			Count:           250,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  5,
			LicenseName:   "Partitioning",
			OriginalCount: 5,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 125,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      125,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 125,
					ConsumedLicensesCount:     125,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      250,
			Count:           250,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SharedAgreement(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}
	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 5,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db2",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "test-db2",
			LicenseCount:  4,
			LicenseName:   "Partitioning",
			OriginalCount: 4,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: -2,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      4,
					Hostname:                  "test-db2",
					TotalCoveredLicensesCount: 4,
					ConsumedLicensesCount:     4,
				},
				{
					CoveredLicensesCount:      1,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 1,
					ConsumedLicensesCount:     3,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SharedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 5,
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
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
		{
			AgreementID:    "5051863",
			AvailableCount: 10,
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
			LicensesCount:   10,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           10,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  20,
			LicenseName:   "Partitioning",
			OriginalCount: 20,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: -5,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      10,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 15,
					ConsumedLicensesCount:     20,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   10,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           10,
		},
		{
			AgreementID:    "5051863",
			AvailableCount: -5,
			CatchAll:       false,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					CoveredLicensesCount:      5,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 15,
					ConsumedLicensesCount:     20,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleUnlimitedCaseNoAssociatedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  0,
			CatchAll:        true,
			CSI:             "CSI001",
			Hosts:           []dto.OracleDatabaseAgreementAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  0,
			CatchAll:        true,
			CSI:             "CSI001",
			Hosts:           []dto.OracleDatabaseAgreementAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleProcessorPerpetualCaseNoAssociatedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  5,
			CatchAll:        true,
			CSI:             "CSI001",
			Hosts:           []dto.OracleDatabaseAgreementAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  2,
			CatchAll:        true,
			CSI:             "CSI001",
			Hosts:           []dto.OracleDatabaseAgreementAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_SimpleNamedUserPlusCaseNoAssociatedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  200,
			CatchAll:        true,
			CSI:             "CSI001",
			Hosts:           []dto.OracleDatabaseAgreementAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      200,
			Count:           200,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  5,
			LicenseName:   "Partitioning",
			OriginalCount: 5,
			Type:          "host",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  75,
			CatchAll:        true,
			CSI:             "CSI001",
			Hosts:           []dto.OracleDatabaseAgreementAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      200,
			Count:           200,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestAssignOracleDatabaseAgreementsToHosts_CompleCase1(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "PID002",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metric:          model.AgreementPartMetricProcessorPerpetual,
			},
		},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 10,
			CatchAll:       true,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{Hostname: "test-db"},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   10,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           10,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "dbclust",
			LicenseCount:  20,
			LicenseName:   "Partitioning",
			OriginalCount: 20,
			Type:          "cluster",
		},
	}

	expectedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			AgreementID:    "5051863",
			AvailableCount: -13,
			CatchAll:       true,
			CSI:            "CSI001",
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					Hostname:                  "test-db",
					CoveredLicensesCount:      3,
					TotalCoveredLicensesCount: 3,
					ConsumedLicensesCount:     3,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   10,
			Metric:          model.AgreementPartMetricProcessorPerpetual,
			PartID:          "PID002",
			ReferenceNumber: "RF0001",
			Unlimited:       false,
			UsersCount:      0,
			Count:           10,
		},
	}

	as.assignOracleDatabaseAgreementsToHosts(agreements, hosts)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestSortHostsUsingLicenses(t *testing.T) {
	list := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseName:  "Diagnostics Pack",
			Name:         "Puzzait",
			Type:         "cluster",
			LicenseCount: 70,
		},
		{
			LicenseName:  "Real Application Clusters",
			Name:         "test-db3",
			Type:         "host",
			LicenseCount: 1.5,
		},
		{
			LicenseName:  "Diagnostics Pack",
			Name:         "test-db4",
			Type:         "host",
			LicenseCount: 0.5,
		},
		{
			LicenseName:  "Oracle ENT",
			Name:         "test-db3",
			Type:         "host",
			LicenseCount: 0.5,
		},
		{
			LicenseName:  "Oracle ENT",
			Name:         "Puzzait",
			Type:         "cluster",
			LicenseCount: 70,
		},
	}

	expected := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseName:  "Oracle ENT",
			Name:         "Puzzait",
			Type:         "cluster",
			LicenseCount: 70,
		},
		{
			LicenseName:  "Diagnostics Pack",
			Name:         "Puzzait",
			Type:         "cluster",
			LicenseCount: 70,
		},
		{
			LicenseName:  "Real Application Clusters",
			Name:         "test-db3",
			Type:         "host",
			LicenseCount: 1.5,
		},
		{
			LicenseName:  "Diagnostics Pack",
			Name:         "test-db4",
			Type:         "host",
			LicenseCount: 0.5,
		},
		{
			LicenseName:  "Oracle ENT",
			Name:         "test-db3",
			Type:         "host",
			LicenseCount: 0.5,
		},
	}

	sortHostsUsingLicenses(list)

	assert.Equal(t, expected, list)
}

func TestSortOracleDatabaseAgreements(t *testing.T) {
	list := []dto.OracleDatabaseAgreementFE{
		{CatchAll: true, Unlimited: false, UsersCount: 10},
		{CatchAll: true, Unlimited: false, LicensesCount: 10},
		{CatchAll: true, Unlimited: true, UsersCount: 20},
		{CatchAll: false, Unlimited: false, LicensesCount: 20},
		{CatchAll: false, Unlimited: true, UsersCount: 10},
		{CatchAll: false, Unlimited: true, LicensesCount: 20},
		{CatchAll: false, Unlimited: false, LicensesCount: 10},
		{CatchAll: true, Unlimited: true, LicensesCount: 10},
		{CatchAll: false, Unlimited: true, UsersCount: 20},
		{CatchAll: false, Unlimited: false, UsersCount: 10},
		{CatchAll: true, Unlimited: true, UsersCount: 10},
		{CatchAll: true, Unlimited: true, LicensesCount: 20},
		{CatchAll: true, Unlimited: false, LicensesCount: 20},
		{CatchAll: false, Unlimited: false, UsersCount: 20},
		{CatchAll: false, Unlimited: true, LicensesCount: 10},
		{CatchAll: true, Unlimited: false, UsersCount: 20},
	}

	expected := []dto.OracleDatabaseAgreementFE{
		{CatchAll: false, Unlimited: false, UsersCount: 20},
		{CatchAll: false, Unlimited: false, UsersCount: 10},
		{CatchAll: false, Unlimited: false, LicensesCount: 20},
		{CatchAll: false, Unlimited: false, LicensesCount: 10},
		{CatchAll: false, Unlimited: true, UsersCount: 20},
		{CatchAll: false, Unlimited: true, UsersCount: 10},
		{CatchAll: false, Unlimited: true, LicensesCount: 20},
		{CatchAll: false, Unlimited: true, LicensesCount: 10},
		{CatchAll: true, Unlimited: false, UsersCount: 20},
		{CatchAll: true, Unlimited: false, UsersCount: 10},
		{CatchAll: true, Unlimited: false, LicensesCount: 20},
		{CatchAll: true, Unlimited: false, LicensesCount: 10},
		{CatchAll: true, Unlimited: true, UsersCount: 20},
		{CatchAll: true, Unlimited: true, UsersCount: 10},
		{CatchAll: true, Unlimited: true, LicensesCount: 20},
		{CatchAll: true, Unlimited: true, LicensesCount: 10},
	}

	sortOracleDatabaseAgreements(list)

	assert.Equal(t, expected, list)
}

func TestSortAssociatedHostsInOracleDatabaseAgreement(t *testing.T) {
	partsMap := map[string]*model.OracleDatabasePart{
		"L10005": {
			PartID:          "L10005",
			ItemDescription: "Oracle Real Application Clusters",
			Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
			Aliases:         []string{"Real Application Clusters", "RAC or RAC One Node"},
		},
	}

	hostsMap := map[string]map[string]*dto.HostUsingOracleDatabaseLicenses{
		"Real Application Clusters": {
			"test-db1": {
				LicenseCount: 10,
			},
			"test-db2": {
				LicenseCount: 30,
			},
		},
		"RAC or RAC One Node": {
			"test-db1": {
				LicenseCount: 20,
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
		PartID: "L10005",
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
			LicenseName:  "Oracle ENT",
			Name:         "Puzzait",
			Type:         "cluster",
			LicenseCount: 70,
		},
		{
			LicenseName:  "Diagnostics Pack",
			Name:         "Puzzait",
			Type:         "cluster",
			LicenseCount: 70,
		},
		{
			LicenseName:  "Real Application Clusters",
			Name:         "test-db3",
			Type:         "host",
			LicenseCount: 1.5,
		},
		{
			LicenseName:  "Diagnostics Pack",
			Name:         "test-db4",
			Type:         "host",
			LicenseCount: 0.5,
		},
		{
			LicenseName:  "Oracle ENT",
			Name:         "test-db3",
			Type:         "host",
			LicenseCount: 0.5,
		},
	}

	expected := map[string]map[string]*dto.HostUsingOracleDatabaseLicenses{
		"Oracle ENT": {
			"Puzzait":  &list[0],
			"test-db3": &list[4],
		},
		"Diagnostics Pack": {
			"Puzzait":  &list[1],
			"test-db4": &list[3],
		},
		"Real Application Clusters": {
			"test-db3": &list[2],
		},
	}

	assert.Equal(t, expected, buildHostUsingLicensesMap(list))
}

func TestBuildAgreementPartMap(t *testing.T) {
	list := []model.OracleDatabasePart{
		{
			ItemDescription: "itemDesc1",
			Aliases:         []string{"alias1"},
			Metric:          "metric1",
			PartID:          "PID001",
		},
		{
			ItemDescription: "itemDesc2",
			Aliases:         []string{"alias1"},
			Metric:          "metric2",
			PartID:          "PID002",
		},
	}

	expected := map[string]*model.OracleDatabasePart{
		"PID001": &list[0],
		"PID002": &list[1],
	}

	assert.Equal(t, expected, buildAgreementPartMap(list))
}

func TestDeleteAssociatedPartFromOracleDatabaseAgreement(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: partsSample,
		TimeNow:                      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID:                  utils.NewObjectIDForTests(),
	}

	associatedPartID := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Fail: can't find associated part", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(associatedPartID).
				Return(nil, utils.AerrOracleDatabaseAssociatedPartNotFound),
		)

		err := as.DeleteAssociatedPartFromOracleDatabaseAgreement(associatedPartID)
		require.EqualError(t, err, utils.AerrOracleDatabaseAssociatedPartNotFound.Error())
	})

	t.Run("Success with only one associated part", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
			},
		}

		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(associatedPartID).
				Return(&agreement, nil),
			db.EXPECT().RemoveOracleDatabaseAgreement(agreement.ID).
				Return(nil),
		)

		err := as.DeleteAssociatedPartFromOracleDatabaseAgreement(associatedPartID)
		assert.Nil(t, err)
	})

	t.Run("Success with multiple associated parts", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
				{
					ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
					OracleDatabasePart: partsSample[1],
					ReferenceNumber:    "RF0002",
					Unlimited:          false,
					Count:              42,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
			},
		}

		agreementResult := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
					OracleDatabasePart: partsSample[1],
					ReferenceNumber:    "RF0002",
					Unlimited:          false,
					Count:              42,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
			},
		}
		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(associatedPartID).
				Return(&agreement, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreementResult).
				Return(nil),
		)

		err := as.DeleteAssociatedPartFromOracleDatabaseAgreement(associatedPartID)
		assert.Nil(t, err)
	})
}

func TestAddHostToAssociatedPart(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: partsSample,
		TimeNow:                      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID:                  utils.NewObjectIDForTests(),
	}

	anotherAssociatedPartID := utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb")
	t.Run("Fail: can't find associated part", func(t *testing.T) {

		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(anotherAssociatedPartID).
				Return(nil, utils.AerrOracleDatabaseAssociatedPartNotFound),
		)

		err := as.AddHostToAssociatedPart(anotherAssociatedPartID, "pippo")
		require.EqualError(t, err, utils.AerrOracleDatabaseAssociatedPartNotFound.Error())
	})

	associatedPartID := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Success with only one associated part", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
			},
		}

		agreementPostAdd := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx", "foobar"},
				},
			},
		}

		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(associatedPartID).
				Return(&agreement, nil),
			db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
				GTECPUCores:    -1,
				LTECPUCores:    -1,
				LTECPUThreads:  -1,
				LTEMemoryTotal: -1,
				GTECPUThreads:  -1,
				GTESwapTotal:   -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
			}, "", false, -1, -1, "", "", utils.MAX_TIME).Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreementPostAdd).
				Return(nil),
		)

		err := as.AddHostToAssociatedPart(associatedPartID, "foobar")
		assert.Nil(t, err)
	})

	t.Run("Success with multiple associated parts", func(t *testing.T) {

		agreement := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
				{
					ID:                 anotherAssociatedPartID,
					OracleDatabasePart: partsSample[1],
					ReferenceNumber:    "RF0002",
					Unlimited:          false,
					Count:              42,
					CatchAll:           false,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
			},
		}

		agreementPostAdd := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
				{
					ID:                 anotherAssociatedPartID,
					OracleDatabasePart: partsSample[1],
					ReferenceNumber:    "RF0002",
					Unlimited:          false,
					Count:              42,
					CatchAll:           false,
					Hosts:              []string{"test-db", "ercsoldbx", "foobar"},
				},
			},
		}

		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(anotherAssociatedPartID).
				Return(&agreement, nil),
			db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
				GTECPUCores:    -1,
				LTECPUCores:    -1,
				LTECPUThreads:  -1,
				LTEMemoryTotal: -1,
				GTECPUThreads:  -1,
				GTESwapTotal:   -1,
				GTEMemoryTotal: -1,
				LTESwapTotal:   -1,
			}, "", false, -1, -1, "", "", utils.MAX_TIME).Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreementPostAdd).
				Return(nil),
		)

		err := as.AddHostToAssociatedPart(anotherAssociatedPartID, "foobar")
		assert.Nil(t, err)
	})
}
func TestRemoveHostFromAssociatedPart(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: partsSample,
		TimeNow:                      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID:                  utils.NewObjectIDForTests(),
	}

	anotherAssociatedPartID := utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb")
	t.Run("Fail: can't find associated part", func(t *testing.T) {

		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(anotherAssociatedPartID).
				Return(nil, utils.AerrOracleDatabaseAssociatedPartNotFound),
		)

		err := as.RemoveHostFromAssociatedPart(anotherAssociatedPartID, "pippo")
		require.EqualError(t, err, utils.AerrOracleDatabaseAssociatedPartNotFound.Error())
	})

	associatedPartID := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Success with only one associated part", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
			},
		}

		agreementPostAdd := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db"},
				},
			},
		}

		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(associatedPartID).
				Return(&agreement, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreementPostAdd).
				Return(nil),
		)

		err := as.RemoveHostFromAssociatedPart(associatedPartID, "ercsoldbx")
		assert.Nil(t, err)
	})

	t.Run("Success with multiple associated parts", func(t *testing.T) {
		agreement := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
				{
					ID:                 anotherAssociatedPartID,
					OracleDatabasePart: partsSample[1],
					ReferenceNumber:    "RF0002",
					Unlimited:          false,
					Count:              42,
					CatchAll:           false,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
			},
		}

		agreementPostAdd := model.OracleDatabaseAgreement{
			AgreementID: "AID001",
			CSI:         "CSI001",
			Parts: []model.AssociatedPart{
				{
					ID:                 associatedPartID,
					OracleDatabasePart: partsSample[0],
					ReferenceNumber:    "RF0001",
					Unlimited:          true,
					Count:              30,
					CatchAll:           true,
					Hosts:              []string{"test-db", "ercsoldbx"},
				},
				{
					ID:                 anotherAssociatedPartID,
					OracleDatabasePart: partsSample[1],
					ReferenceNumber:    "RF0002",
					Unlimited:          false,
					Count:              42,
					CatchAll:           false,
					Hosts:              []string{"test-db"},
				},
			},
		}

		gomock.InOrder(
			db.EXPECT().GetOracleDatabaseAgreementByAssociatedPart(anotherAssociatedPartID).
				Return(&agreement, nil),
			db.EXPECT().UpdateOracleDatabaseAgreement(agreementPostAdd).
				Return(nil),
		)

		err := as.RemoveHostFromAssociatedPart(anotherAssociatedPartID, "ercsoldbx")
		assert.Nil(t, err)
	})
}
