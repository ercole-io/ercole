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
