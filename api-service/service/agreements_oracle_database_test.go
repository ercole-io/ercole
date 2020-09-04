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

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestLoadOracleDatabaseAgreementPartsList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}
	as.LoadOracleDatabaseAgreementPartsList()

	assert.Equal(t, "L10001", as.OracleDatabaseAgreementParts[0].PartID)
	assert.Equal(t, "Oracle Database Enterprise Edition", as.OracleDatabaseAgreementParts[0].ItemDescription)
	assert.Equal(t, "Named User Plus Perpetual", as.OracleDatabaseAgreementParts[0].Metrics)
	assert.Equal(t, []string{"Oracle ENT"}, as.OracleDatabaseAgreementParts[0].Aliases)
	assert.Equal(t, "L103405", as.OracleDatabaseAgreementParts[2].PartID)
	assert.Equal(t, []string{"Oracle STD"}, as.OracleDatabaseAgreementParts[2].Aliases)

	//Known list of metrics check!
	for i, part := range as.OracleDatabaseAgreementParts {
		assert.Contains(t,
			[]string{"Processor Perpetual", "Named User Plus Perpetual", "Stream Perpetual", "Computer Perpetual"},
			part.Metrics,
			"There is a Oracle/Database agreement part with unknown metric #", i, part,
		)
	}
}

func TestGetOracleDatabaseAgreementPartsList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{},
		},
	}
	res, err := as.GetOracleDatabaseAgreementPartsList()
	require.NoError(t, err)
	assert.Equal(t, []model.OracleDatabaseAgreementPart{
		{},
	}, res)
}

func TestAddOracleDatabaseAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				ItemDescription: "asdasdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdsdfasasd",
				PartID:          "L10006",
			},
			{
				ItemDescription: "asdasdfdsfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdasjkhasd",
				PartID:          "A90620",
			},
			{
				ItemDescription: "asdsdfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdas2435asd",
				PartID:          "A90650",
			},
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	addRequest := apimodel.OracleDatabaseAgreementsAddRequest{
		AgreementID: "5051863",
		PartsID: []string{
			"L10006",
			"A90620",
		},
		CSI:             "6871235",
		ReferenceNumber: "10032246681",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Hosts: []string{
			"test-db",
			"ercsoldbx",
		},
	}

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
	}, nil)

	db.EXPECT().InsertOracleDatabaseAgreement(gomock.Any()).Do(func(agg model.OracleDatabaseAgreement) {
		assert.Equal(t, "5051863", agg.AgreementID)
		assert.Equal(t, "6871235", agg.CSI)
		assert.True(t, agg.CatchAll)
		assert.Equal(t, 30, agg.Count)
		assert.Equal(t, agg.Hosts, []string{
			"test-db",
			"ercsoldbx",
		})
		assert.Equal(t, "asdasdfdsfsdas", agg.ItemDescription)
		assert.Equal(t, "sdasjkhasd", agg.Metrics)
		assert.Equal(t, "A90620", agg.PartID)
		assert.Equal(t, "10032246681", agg.ReferenceNumber)
		assert.True(t, agg.Unlimited)
	}).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5f4d0a4c9015f713a9c66107")}, nil).After(
		db.EXPECT().InsertOracleDatabaseAgreement(gomock.Any()).Do(func(agg model.OracleDatabaseAgreement) {
			assert.Equal(t, "5051863", agg.AgreementID)
			assert.Equal(t, "6871235", agg.CSI)
			assert.True(t, agg.CatchAll)
			assert.Equal(t, 30, agg.Count)
			assert.Equal(t, agg.Hosts, []string{
				"test-db",
				"ercsoldbx",
			})
			assert.Equal(t, "asdasdas", agg.ItemDescription)
			assert.Equal(t, "sdsdfasasd", agg.Metrics)
			assert.Equal(t, "L10006", agg.PartID)
			assert.Equal(t, "10032246681", agg.ReferenceNumber)
			assert.True(t, agg.Unlimited)
		}).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5f4d0a2b27fe53da8a4aec45")}, nil),
	)
	res, err := as.AddOracleDatabaseAgreements(addRequest)
	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON([]mongo.InsertOneResult{
		{InsertedID: utils.Str2oid("5f4d0a2b27fe53da8a4aec45")},
		{InsertedID: utils.Str2oid("5f4d0a4c9015f713a9c66107")},
	}), utils.ToJSON(res))
}

func TestAddOracleDatabaseAgreements_Fail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				ItemDescription: "asdasdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdsdfasasd",
				PartID:          "L10006",
			},
			{
				ItemDescription: "asdasdfdsfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdasjkhasd",
				PartID:          "A90620",
			},
			{
				ItemDescription: "asdsdfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdas2435asd",
				PartID:          "A90650",
			},
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	addRequest := apimodel.OracleDatabaseAgreementsAddRequest{
		AgreementID: "5051863",
		PartsID: []string{
			"L10006",
			"A90620dsf",
		},
		CSI:             "6871235",
		ReferenceNumber: "10032246681",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Hosts: []string{
			"test-db",
			"ercsoldbx",
		},
	}

	_, err := as.AddOracleDatabaseAgreements(addRequest)
	require.Equal(t, err, utils.AerrOracleDatabaseAgreementInvalidPartID)
}

func TestAddOracleDatabaseAgreements_Fail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				ItemDescription: "asdasdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdsdfasasd",
				PartID:          "L10006",
			},
			{
				ItemDescription: "asdasdfdsfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdasjkhasd",
				PartID:          "A90620",
			},
			{
				ItemDescription: "asdsdfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdas2435asd",
				PartID:          "A90650",
			},
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	addRequest := apimodel.OracleDatabaseAgreementsAddRequest{
		AgreementID: "5051863",
		PartsID: []string{
			"L10006",
			"A90620",
		},
		CSI:             "6871235",
		ReferenceNumber: "10032246681",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Hosts: []string{
			"test-db",
			"ercsoldbx",
		},
	}

	db.EXPECT().SearchHosts("hostnames", []string{""}, database.SearchHostsFilters{
		GTECPUCores:    -1,
		LTECPUCores:    -1,
		LTECPUThreads:  -1,
		LTEMemoryTotal: -1,
		GTECPUThreads:  -1,
		GTESwapTotal:   -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
	}, "", false, -1, -1, "", "", utils.MAX_TIME).Return(nil, aerrMock)

	_, err := as.AddOracleDatabaseAgreements(addRequest)
	assert.Equal(t, aerrMock, err)
}

func TestAddOracleDatabaseAgreements_Fail3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				ItemDescription: "asdasdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdsdfasasd",
				PartID:          "L10006",
			},
			{
				ItemDescription: "asdasdfdsfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdasjkhasd",
				PartID:          "A90620",
			},
			{
				ItemDescription: "asdsdfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdas2435asd",
				PartID:          "A90650",
			},
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	addRequest := apimodel.OracleDatabaseAgreementsAddRequest{
		AgreementID: "5051863",
		PartsID: []string{
			"L10006",
			"A90620",
		},
		CSI:             "6871235",
		ReferenceNumber: "10032246681",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Hosts: []string{
			"test-db",
			"ercsoldbx",
		},
	}

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
	}, nil)

	_, err := as.AddOracleDatabaseAgreements(addRequest)
	assert.Equal(t, utils.AerrHostNotFound, err)
}

func TestAddOracleDatabaseAgreements_Fail4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				ItemDescription: "asdasdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdsdfasasd",
				PartID:          "L10006",
			},
			{
				ItemDescription: "asdasdfdsfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdasjkhasd",
				PartID:          "A90620",
			},
			{
				ItemDescription: "asdsdfsdas",
				Aliases:         []string{"dasasd"},
				Metrics:         "sdas2435asd",
				PartID:          "A90650",
			},
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	addRequest := apimodel.OracleDatabaseAgreementsAddRequest{
		AgreementID: "5051863",
		PartsID: []string{
			"L10006",
			"A90620",
		},
		CSI:             "6871235",
		ReferenceNumber: "10032246681",
		Unlimited:       true,
		Count:           30,
		CatchAll:        true,
		Hosts: []string{
			"test-db",
			"ercsoldbx",
		},
	}

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
	}, nil)

	db.EXPECT().InsertOracleDatabaseAgreement(gomock.Any()).Do(func(agg model.OracleDatabaseAgreement) {
		assert.Equal(t, "5051863", agg.AgreementID)
		assert.Equal(t, "6871235", agg.CSI)
		assert.True(t, agg.CatchAll)
		assert.Equal(t, 30, agg.Count)
		assert.Equal(t, agg.Hosts, []string{
			"test-db",
			"ercsoldbx",
		})
		assert.Equal(t, "asdasdas", agg.ItemDescription)
		assert.Equal(t, "sdsdfasasd", agg.Metrics)
		assert.Equal(t, "L10006", agg.PartID)
		assert.Equal(t, "10032246681", agg.ReferenceNumber)
		assert.True(t, agg.Unlimited)
	}).Return(nil, aerrMock)
	_, err := as.AddOracleDatabaseAgreements(addRequest)
	assert.Equal(t, aerrMock, err)
}

func TestSearchOracleDatabaseAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	returnedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	returnedLicensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 0,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListOracleDatabaseLicensingObjects().Return(returnedLicensingObjects, nil)

	res, err := as.SearchOracleDatabaseAgreements("", apimodel.SearchOracleDatabaseAgreementsFilters{
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
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	returnedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	returnedLicensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListOracleDatabaseLicensingObjects().Return(returnedLicensingObjects, nil)

	res, err := as.SearchOracleDatabaseAgreements("", apimodel.SearchOracleDatabaseAgreementsFilters{
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

func TestSearchOracleDatabaseAgreements_Failed1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	db.EXPECT().ListOracleDatabaseAgreements().Return(nil, aerrMock)

	_, err := as.SearchOracleDatabaseAgreements("", apimodel.SearchOracleDatabaseAgreementsFilters{
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

func TestSearchOracleDatabaseAgreements_Failed2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	returnedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)
	db.EXPECT().ListOracleDatabaseLicensingObjects().Return(nil, aerrMock)

	_, err := as.SearchOracleDatabaseAgreements("", apimodel.SearchOracleDatabaseAgreementsFilters{
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
	agg1 := apimodel.OracleDatabaseAgreementsFE{
		AgreementID:    "5051863",
		AvailableCount: 7,
		CatchAll:       true,
		CSI:            "6871235",
		Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
		Metrics:         "Processor Perpetual",
		PartID:          "A90620",
		ReferenceNumber: "10032246681",
		Unlimited:       false,
		UsersCount:      5,
	}

	assert.True(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))

	assert.True(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		AgreementID:       "5051",
		PartID:            "A9062",
		ItemDescription:   "Partitioning",
		CSI:               "6871",
		Metrics:           "Processor Perpetual",
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
	assert.True(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: 7,
		AvailableCountLTE: 7,
		LicensesCountGTE:  30,
		LicensesCountLTE:  30,
		UsersCountGTE:     5,
		UsersCountLTE:     5,
	}))

	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
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
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
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
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
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
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
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
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Metrics:           "fdgdfgsdsfg",
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
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
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "true",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "false",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  35,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  25,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     0,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     10,
		UsersCountLTE:     -1,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
		Unlimited:         "NULL",
		CatchAll:          "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: 3,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}))
	assert.False(t, CheckOracleDatabaseAgreementMatchFilter(agg1, apimodel.SearchOracleDatabaseAgreementsFilters{
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

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SimpleUnlimitedCase(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 0,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SimpleProcessorPerpetualCase(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 0,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           2,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SimpleNamedUserPlusCase(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Named User Plus Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Named User Plus Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      10,
			Count:           10,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         128,
			LicenseName:   "Partitioning",
			OriginalCount: 128,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: -3,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      125,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 125,
					ConsumedLicensesCount:     128,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Named User Plus Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      10,
			Count:           5,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SharedAgreement(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}
	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "test-db2",
			Count:         4,
			LicenseName:   "Partitioning",
			OriginalCount: 4,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: -2,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           0,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SharedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					CoveredLicensesCount:      0,
					Hostname:                  "test-db",
					TotalCoveredLicensesCount: 0,
				},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   10,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           10,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         20,
			LicenseName:   "Partitioning",
			OriginalCount: 20,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: -5,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           0,
		},
		{
			AgreementID:    "5051863",
			AvailableCount: -5,
			CatchAll:       false,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
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
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           0,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SimpleUnlimitedCaseNoAssociatedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  7,
			CatchAll:        true,
			CSI:             "6871235",
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  -0,
			CatchAll:        true,
			CSI:             "6871235",
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       true,
			UsersCount:      0,
			Count:           0,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SimpleProcessorPerpetualCaseNoAssociatedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  7,
			CatchAll:        true,
			CSI:             "6871235",
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           5,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  0,
			CatchAll:        true,
			CSI:             "6871235",
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   5,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           2,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_SimpleNamedUserPlusCaseNoAssociatedHost(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Named User Plus Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  7,
			CatchAll:        true,
			CSI:             "6871235",
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Named User Plus Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      10,
			Count:           10,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         128,
			LicenseName:   "Partitioning",
			OriginalCount: 128,
			Type:          "host",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:     "5051863",
			AvailableCount:  -3,
			CatchAll:        true,
			CSI:             "6871235",
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   0,
			Metrics:         "Named User Plus Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      10,
			Count:           5,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)
}

func TestGreedilyAssignOracleDatabaseAgreementsToLicensingObjects_CompleCase1(t *testing.T) {
	as := APIService{
		Config: config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{
				PartID:          "A90620",
				Aliases:         []string{"Partitioning"},
				ItemDescription: "Oracle Partitioning",
				Metrics:         "Processor Perpetual",
			},
		},
	}

	agreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: 7,
			CatchAll:       true,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{Hostname: "test-db"},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   10,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           10,
		},
	}
	licensingObjects := []apimodel.OracleDatabaseLicensingObjects{
		{
			Name:          "test-db",
			Count:         3,
			LicenseName:   "Partitioning",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "dbclust",
			Count:         20,
			LicenseName:   "Partitioning",
			OriginalCount: 20,
			Type:          "cluster",
		},
	}

	expectedAgreements := []apimodel.OracleDatabaseAgreementsFE{
		{
			AgreementID:    "5051863",
			AvailableCount: -13,
			CatchAll:       true,
			CSI:            "6871235",
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{Hostname: "test-db", CoveredLicensesCount: 3, TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3},
			},
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ItemDescription: "Oracle Partitioning",
			LicensesCount:   10,
			Metrics:         "Processor Perpetual",
			PartID:          "A90620",
			ReferenceNumber: "10032246681",
			Unlimited:       false,
			UsersCount:      0,
			Count:           0,
		},
	}

	as.GreedilyAssignOracleDatabaseAgreementsToLicensingObjects(agreements, licensingObjects)

	assert.Equal(t, expectedAgreements, agreements)

}

func TestSortOracleDatabaseAgreementLicensingObjects(t *testing.T) {
	list := []apimodel.OracleDatabaseLicensingObjects{
		{
			LicenseName: "Diagnostics Pack",
			Name:        "Puzzait",
			Type:        "cluster",
			Count:       70,
		},
		{
			LicenseName: "Real Application Clusters",
			Name:        "test-db3",
			Type:        "host",
			Count:       1.5,
		},
		{
			LicenseName: "Diagnostics Pack",
			Name:        "test-db4",
			Type:        "host",
			Count:       0.5,
		},
		{
			LicenseName: "Oracle ENT",
			Name:        "test-db3",
			Type:        "host",
			Count:       0.5,
		},
		{
			LicenseName: "Oracle ENT",
			Name:        "Puzzait",
			Type:        "cluster",
			Count:       70,
		},
	}

	expected := []apimodel.OracleDatabaseLicensingObjects{
		{
			LicenseName: "Oracle ENT",
			Name:        "Puzzait",
			Type:        "cluster",
			Count:       70,
		},
		{
			LicenseName: "Diagnostics Pack",
			Name:        "Puzzait",
			Type:        "cluster",
			Count:       70,
		},
		{
			LicenseName: "Real Application Clusters",
			Name:        "test-db3",
			Type:        "host",
			Count:       1.5,
		},
		{
			LicenseName: "Diagnostics Pack",
			Name:        "test-db4",
			Type:        "host",
			Count:       0.5,
		},
		{
			LicenseName: "Oracle ENT",
			Name:        "test-db3",
			Type:        "host",
			Count:       0.5,
		},
	}

	SortOracleDatabaseAgreementLicensingObjects(list)

	assert.Equal(t, expected, list)
}

func TestSortOracleDatabaseAgreements(t *testing.T) {
	list := []apimodel.OracleDatabaseAgreementsFE{
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

	expected := []apimodel.OracleDatabaseAgreementsFE{
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

	SortOracleDatabaseAgreements(list)

	assert.Equal(t, expected, list)
}

func TestSortAssociatedHostsInOracleDatabaseAgreement(t *testing.T) {
	partsMap := map[string]*model.OracleDatabaseAgreementPart{
		"L10005": {
			PartID:          "L10005",
			ItemDescription: "Oracle Real Application Clusters",
			Metrics:         "Named User Plus Perpetual",
			Aliases:         []string{"Real Application Clusters", "RAC or RAC One Node"},
		},
	}

	licensingObjectsMap := map[string]map[string]*apimodel.OracleDatabaseLicensingObjects{
		"Real Application Clusters": {
			"test-db1": {
				Count: 10,
			},
			"test-db2": {
				Count: 30,
			},
		},
		"RAC or RAC One Node": {
			"test-db1": {
				Count: 20,
			},
			"test-db3": {
				Count: 15,
			},
			"test-db4": {
				Count: 35,
			},
		},
	}

	agg := apimodel.OracleDatabaseAgreementsFE{
		PartID: "L10005",
		Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
			{Hostname: "test-db2"},
			{Hostname: "test-db1"},
			{Hostname: "test-db4"},
			{Hostname: "test-db3"},
		},
	}

	expected := []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
		{Hostname: "test-db4"},
		{Hostname: "test-db2"},
		{Hostname: "test-db1"},
		{Hostname: "test-db3"},
	}

	SortAssociatedHostsInOracleDatabaseAgreement(agg, licensingObjectsMap, partsMap)

	assert.Equal(t, expected, agg.Hosts)
}

func TestBuildOracleDatabaseLicensingObjectsMap(t *testing.T) {
	list := []apimodel.OracleDatabaseLicensingObjects{
		{
			LicenseName: "Oracle ENT",
			Name:        "Puzzait",
			Type:        "cluster",
			Count:       70,
		},
		{
			LicenseName: "Diagnostics Pack",
			Name:        "Puzzait",
			Type:        "cluster",
			Count:       70,
		},
		{
			LicenseName: "Real Application Clusters",
			Name:        "test-db3",
			Type:        "host",
			Count:       1.5,
		},
		{
			LicenseName: "Diagnostics Pack",
			Name:        "test-db4",
			Type:        "host",
			Count:       0.5,
		},
		{
			LicenseName: "Oracle ENT",
			Name:        "test-db3",
			Type:        "host",
			Count:       0.5,
		},
	}

	expected := map[string]map[string]*apimodel.OracleDatabaseLicensingObjects{
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

	assert.Equal(t, expected, BuildOracleDatabaseLicensingObjectsMap(list))
}

func TestBuildOracleDatabaseAgreementPartMap(t *testing.T) {
	list := []model.OracleDatabaseAgreementPart{
		{
			ItemDescription: "asdasdas",
			Aliases:         []string{"dasasd"},
			Metrics:         "sdsdfasasd",
			PartID:          "L10006",
		},
		{
			ItemDescription: "asdasdfdsfsdas",
			Aliases:         []string{"dasasd"},
			Metrics:         "sdasjkhasd",
			PartID:          "A90620",
		},
	}

	expected := map[string]*model.OracleDatabaseAgreementPart{
		"L10006": &list[0],
		"A90620": &list[1],
	}

	assert.Equal(t, expected, BuildOracleDatabaseAgreementPartMap(list))
}

func TestAddAssociatedHostToOracleDatabaseAgreement_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	returnedAgg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	updatedAgg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar", "foohost"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	db.EXPECT().ExistNotInClusterHost("foohost").Return(true, nil)
	db.EXPECT().FindOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(returnedAgg, nil)
	db.EXPECT().UpdateOracleDatabaseAgreement(updatedAgg).Return(nil)

	err := as.AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.NoError(t, err)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_SuccessHostIsAlreadyAssociated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	agg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar", "foohost"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	db.EXPECT().ExistNotInClusterHost("foohost").Return(true, nil)
	db.EXPECT().FindOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(agg, nil)

	err := as.AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.NoError(t, err)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedHostNotExist(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().ExistNotInClusterHost("foohost").Return(false, nil)

	err := as.AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.Equal(t, utils.AerrNotInClusterHostNotFound, err)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().ExistNotInClusterHost("foohost").Return(false, aerrMock)

	err := as.AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.Equal(t, aerrMock, err)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().ExistNotInClusterHost("foohost").Return(true, nil)
	db.EXPECT().FindOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(model.OracleDatabaseAgreement{}, aerrMock)

	err := as.AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.Equal(t, aerrMock, err)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedInternalServerError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	returnedAgg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	updatedAgg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar", "foohost"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	db.EXPECT().ExistNotInClusterHost("foohost").Return(true, nil)
	db.EXPECT().FindOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(returnedAgg, nil)
	db.EXPECT().UpdateOracleDatabaseAgreement(updatedAgg).Return(aerrMock)

	err := as.AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.Equal(t, aerrMock, err)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	returnedAgg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar", "foohost"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	updatedAgg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	db.EXPECT().FindOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(returnedAgg, nil)
	db.EXPECT().UpdateOracleDatabaseAgreement(updatedAgg).Return(nil)

	err := as.RemoveAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.NoError(t, err)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_SuccessNoHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	returnedAgg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	db.EXPECT().FindOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(returnedAgg, nil)

	err := as.RemoveAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.NoError(t, err)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_FailedInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(model.OracleDatabaseAgreement{}, aerrMock)

	err := as.RemoveAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost")
	require.Equal(t, aerrMock, err)
}
