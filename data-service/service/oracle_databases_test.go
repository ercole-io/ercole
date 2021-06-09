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
	"fmt"
	reflect "reflect"
	"sort"
	"testing"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type alertSimilarTo struct{ al model.Alert }

func (sa *alertSimilarTo) Matches(x interface{}) bool {
	if val, ok := x.(model.Alert); !ok {
		return false
	} else if sa.al.AlertSeverity != "" && val.AlertSeverity != sa.al.AlertSeverity {
		return false
	} else if val.AlertCode != sa.al.AlertCode {
		return false
	} else if (sa.al.AlertAffectedTechnology == nil && val.AlertAffectedTechnology != sa.al.AlertAffectedTechnology) || (sa.al.AlertAffectedTechnology != nil && *val.AlertAffectedTechnology != *sa.al.AlertAffectedTechnology) {
		return false
	} else if val.AlertCategory != sa.al.AlertCategory {
		return false
	} else if !reflect.DeepEqual(sa.al.OtherInfo, val.OtherInfo) {
		return false
	}

	return true
}

func (sa *alertSimilarTo) String() string {
	return fmt.Sprintf("is similar to %v", sa.al)
}

func TestAddLicensesToSecondaryDbs(t *testing.T) {
	//TODO Add test
	//	mockCtrl := gomock.NewController(t)
	//	defer mockCtrl.Finish()
	//	mongoDb := NewMockMongoDatabaseInterface(mockCtrl)
	//
	//	hds := HostDataService{
	//		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	//		Database: mongoDb,
	//		Config: config.Configuration{
	//			AlertService: config.AlertService{
	//				PublisherUsername: "publ1sh3r",
	//				PublisherPassword: "M0stS3cretP4ssw0rd",
	//				RemoteEndpoint:    "http://ercole.example.org",
	//			},
	//			DataService: config.DataService{
	//				EnablePatching:       true,
	//				LogInsertingHostdata: true,
	//			},
	//		},
	//		Version: "1.6.6",
	//		Log:     utils.NewLogger("TEST"),
	//	}
	//
	//	hdPrimary := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_22.json")
	//	hdSecondary := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_23.json")
	//
	//	db := hdSecondary.Features.Oracle.Database.Databases[0]
	//	assert.True(t, db.Status == model.OracleDatabaseStatusMounted && db.Role != model.OracleDatabaseRolePrimary)
	//
	//	hds.addLicensesToSecondaryDbs(hdSecondary.Info, &db)
}

var hostData1 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dc3f534db7e81a98b726a52"),
	Hostname:  "superhost1",
	Archived:  false,
	CreatedAt: utils.P("2019-11-05T14:02:03Z"),
	Features: model.Features{
		Oracle: &model.OracleFeature{
			Database: &model.OracleDatabaseFeature{
				UnlistedRunningDatabases: []string{"FOOBAR"},
				Databases:                []model.OracleDatabase{},
			},
		},
	},
	Info: model.Host{
		CPUCores: 0,
	},
}

var hostData2 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T12:02:03Z"),
	Features: model.Features{
		Oracle: &model.OracleFeature{
			Database: &model.OracleDatabaseFeature{
				UnlistedRunningDatabases: []string{},
				Databases:                []model.OracleDatabase{},
			},
		},
	},
	Info: model.Host{
		CPUCores: 0,
	},
}

var hostData3 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T16:02:03Z"),
	Features: model.Features{
		Oracle: &model.OracleFeature{
			Database: &model.OracleDatabaseFeature{
				UnlistedRunningDatabases: []string{},
				Databases: []model.OracleDatabase{
					{
						Name:     "acd",
						Licenses: []model.OracleDatabaseLicense{},
					},
				},
			},
		},
	},
	Info: model.Host{
		CPUCores: 2,
	},
}

var hostData4 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		Oracle: &model.OracleFeature{
			Database: &model.OracleDatabaseFeature{
				UnlistedRunningDatabases: []string{},
				Databases: []model.OracleDatabase{
					{
						Name: "acd",
						Licenses: []model.OracleDatabaseLicense{
							{
								LicenseTypeID: "Oracle ENT",
								Name:          "Oracle ENT",
								Count:         10,
							},
							{
								LicenseTypeID: "Driving",
								Name:          "Driving",
								Count:         100,
							},
						},
					},
				},
			},
		},
	},
	Info: model.Host{
		CPUCores: 0,
	},
}

var hostData5 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		Oracle: &model.OracleFeature{
			Database: &model.OracleDatabaseFeature{
				UnlistedRunningDatabases: []string{},
				Databases: []model.OracleDatabase{
					{
						Name: "acd",
						Licenses: []model.OracleDatabaseLicense{
							{
								LicenseTypeID: "Oracle ENT",
								Name:          "Oracle ENT",
								Count:         10,
							},
							{
								LicenseTypeID: "Driving",
								Name:          "Driving",
								Count:         100,
							},
						},
					},
					{
						Name: "acd-two",
						Licenses: []model.OracleDatabaseLicense{
							{
								LicenseTypeID: "Oracle ENT",
								Name:          "Oracle ENT",
								Count:         5,
							},
							{
								LicenseTypeID: "Driving",
								Name:          "Driving",
								Count:         50,
							},
							{
								LicenseTypeID: "Waving",
								Name:          "Waving",
								Count:         50,
							},
						},
					},
				},
			},
		},
	},
	Info: model.Host{
		CPUCores: 0,
	},
}

var licenseTypes = []model.OracleDatabaseLicenseType{
	{
		ID:              "Oracle ENT",
		ItemDescription: "",
		Metric:          "",
		Cost:            0,
		Aliases:         []string{},
		Option:          false,
	},
	{
		ID:              "Driving",
		ItemDescription: "",
		Metric:          "",
		Cost:            0,
		Aliases:         []string{},
		Option:          true,
	},
}

//TODO Create client for api and mock it
//func TestOracleDatabasesChecks(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	db := NewMockMongoDatabaseInterface(mockCtrl)
//	asc := NewMockAlertSvcClientInterface(mockCtrl)
//	hds := HostDataService{
//		Config:         config.Configuration{},
//		ServerVersion:  "",
//		Database:       db,
//		AlertSvcClient: asc,
//		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
//		Log:            utils.NewLogger("TEST"),
//	}
//
//	gomock.InOrder(
//		asc.EXPECT().ThrowNewAlert(&alertSimilarTo{
//			al: model.Alert{
//				AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
//				AlertCategory:           model.AlertCategoryLicense,
//				AlertCode:               model.AlertCodeNewDatabase,
//				OtherInfo: map[string]interface{}{
//					"hostname": "superhost1",
//					"dbname":   "acd",
//				},
//			}}).Return(nil),
//
//		asc.EXPECT().ThrowNewAlert(&alertSimilarTo{
//			al: model.Alert{
//				AlertAffectedTechnology: nil,
//				AlertCategory:           model.AlertCategoryLicense,
//				AlertCode:               model.AlertCodeIncreasedCPUCores,
//				AlertSeverity:           model.AlertSeverityCritical,
//				AlertStatus:             model.AlertStatusNew,
//				OtherInfo: map[string]interface{}{
//					"hostname": "superhost1",
//				},
//			}}).Return(nil),
//	)
//
//	hds.oracleDatabasesChecks(&hostData1, &hostData3)
//}

func TestCheckNewLicenses_SuccessNoDifferences(t *testing.T) {
	hds := HostDataService{
		Log: utils.NewLogger("TEST"),
	}

	require.NoError(t, hds.checkNewLicenses(&hostData2, &hostData1, licenseTypes))
}

func TestCheckNewLicenses_SuccessNewDatabase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	asc.EXPECT().ThrowNewAlert(&alertSimilarTo{
		al: model.Alert{
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCategory:           model.AlertCategoryLicense,
			AlertCode:               model.AlertCodeNewDatabase,
			OtherInfo: map[string]interface{}{
				"hostname": "superhost1",
				"dbname":   "acd",
			},
		}}).Return(nil)

	require.NoError(t, hds.checkNewLicenses(&hostData1, &hostData3, licenseTypes))
}

func TestCheckNewLicenses_ThrowNewLicenseAndNewOption(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	asc.EXPECT().ThrowNewAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"hostname":      "superhost1",
			"dbname":        "acd",
			"licenseTypeID": "Driving",
		},
	}}).Return(nil)
	asc.EXPECT().ThrowNewAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertCode:               model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"hostname":      "superhost1",
			"dbname":        "acd",
			"licenseTypeID": "Oracle ENT",
		},
	}}).Return(nil)

	require.NoError(t, hds.checkNewLicenses(&hostData3, &hostData4, licenseTypes))
}

func TestCheckNewLicenses_ThrowNewLicenseAndNewOptionAlreadyEnabled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	asc.EXPECT().ThrowNewAlert(&alertSimilarTo{
		model.Alert{
			AlertCategory:           model.AlertCategoryLicense,
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCode:               model.AlertCodeNewDatabase,
			AlertSeverity:           model.AlertSeverityInfo,
			AlertStatus:             model.AlertStatusNew,
			Date:                    utils.P("2019-11-05T16:02:03Z"),
			OtherInfo:               map[string]interface{}{"hostname": "superhost1", "dbname": "acd-two"},
		}}).Return(nil)
	asc.EXPECT().ThrowNewAlert(&alertSimilarTo{
		model.Alert{
			AlertCategory:           model.AlertCategoryLicense,
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCode:               model.AlertCodeNewOption,
			AlertSeverity:           model.AlertSeverityInfo,
			Date:                    utils.P("2019-11-05T16:02:03Z"),
			OtherInfo:               map[string]interface{}{"hostname": "superhost1", "dbname": "acd-two", "licenseTypeID": "Driving"},
		}}).Return(nil)
	asc.EXPECT().ThrowNewAlert(&alertSimilarTo{
		model.Alert{
			AlertCategory:           model.AlertCategoryLicense,
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCode:               model.AlertCodeNewLicense,
			AlertSeverity:           model.AlertSeverityInfo,
			AlertStatus:             "",
			Description:             "",
			Date:                    utils.P("2019-11-05T16:02:03Z"),
			OtherInfo:               map[string]interface{}{"hostname": "superhost1", "dbname": "acd-two", "licenseTypeID": "Oracle ENT"},
		}}).Return(nil)

	require.NoError(t, hds.checkNewLicenses(&hostData4, &hostData5, licenseTypes))
}

func TestCheckNewLicenses_CantThrowNewAlert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	t.Run("Fail throwNewDatabaseAlert", func(t *testing.T) {
		asc.EXPECT().ThrowNewAlert(&alertSimilarTo{al: model.Alert{
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCategory:           model.AlertCategoryLicense,
			AlertCode:               model.AlertCodeNewDatabase,
			OtherInfo: map[string]interface{}{
				"hostname": "superhost1",
				"dbname":   "acd",
			},
		}}).Return(aerrMock)

		//TODO Add check that error has been logged
		require.NoError(t, hds.checkNewLicenses(&hostData1, &hostData3, licenseTypes))
	})

	t.Run("Fail throwNewLicenseAlert", func(t *testing.T) {
		asc.EXPECT().ThrowNewAlert(&alertSimilarTo{al: model.Alert{
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCategory:           model.AlertCategoryLicense,
			AlertCode:               model.AlertCodeNewLicense,
			OtherInfo: map[string]interface{}{
				"hostname":      "superhost1",
				"dbname":        "acd",
				"licenseTypeID": "Oracle ENT",
			},
		}}).Return(aerrMock)
		asc.EXPECT().ThrowNewAlert(&alertSimilarTo{al: model.Alert{
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCategory:           model.AlertCategoryLicense,
			AlertCode:               model.AlertCodeNewOption,
			OtherInfo: map[string]interface{}{
				"hostname":      "superhost1",
				"dbname":        "acd",
				"licenseTypeID": "Driving",
			},
		}}).Return(nil)

		//TODO Add check that error has been logged
		require.NoError(t, hds.checkNewLicenses(&hostData3, &hostData4, licenseTypes))
	})

	t.Run("Fail throwActivatedFeaturesAlert", func(t *testing.T) {
		asc.EXPECT().ThrowNewAlert(&alertSimilarTo{al: model.Alert{
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCategory:           model.AlertCategoryLicense,
			AlertCode:               model.AlertCodeNewLicense,
			OtherInfo: map[string]interface{}{
				"hostname":      "superhost1",
				"dbname":        "acd",
				"licenseTypeID": "Oracle ENT",
			},
		}}).Return(nil)
		asc.EXPECT().ThrowNewAlert(&alertSimilarTo{al: model.Alert{
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCategory:           model.AlertCategoryLicense,
			AlertCode:               model.AlertCodeNewOption,
			OtherInfo: map[string]interface{}{
				"hostname":      "superhost1",
				"dbname":        "acd",
				"licenseTypeID": "Driving",
			},
		}}).Return(aerrMock)

		//TODO Add check that error has been logged
		require.NoError(t, hds.checkNewLicenses(&hostData3, &hostData4, licenseTypes))
	})
}

func TestCheckNewLicenses_ErrOracleDatabaseLicenseTypeIDNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	//TODO Add check that error has been logged
	require.NoError(t, hds.checkNewLicenses(&hostData3, &hostData4, []model.OracleDatabaseLicenseType{}))
}

func TestLicenseTypesSorter(t *testing.T) {
	testCases := []struct {
		config       config.DataService
		environment  string
		licenseTypes []model.OracleDatabaseLicenseType
		expected     []model.OracleDatabaseLicenseType
	}{
		{
			config: config.DataService{
				LicenseTypeMetricsDefault:       []string{"Processor Perpetual", "Named User Plus Perpetual", "Stream Perpetual", "Computer Perpetual"},
				LicenseTypeMetricsByEnvironment: map[string][]string{},
			},
			environment: "TST",
			licenseTypes: []model.OracleDatabaseLicenseType{
				{
					ID:              "ID02",
					ItemDescription: "bbbbb",
					Metric:          "",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          true,
				},
				{
					ID:              "ID01",
					ItemDescription: "aaaaa",
					Metric:          "",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
			},
			expected: []model.OracleDatabaseLicenseType{
				{
					ID:              "ID01",
					ItemDescription: "aaaaa",
					Metric:          "",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
				{
					ID:              "ID02",
					ItemDescription: "bbbbb",
					Metric:          "",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          true,
				},
			},
		},
		{
			config: config.DataService{
				LicenseTypeMetricsDefault:       []string{"a", "b", "c"},
				LicenseTypeMetricsByEnvironment: map[string][]string{},
			},
			environment: "TST",
			licenseTypes: []model.OracleDatabaseLicenseType{
				{
					ID:              "ID01",
					ItemDescription: "",
					Metric:          "b",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
				{
					ID:              "ID01",
					ItemDescription: "",
					Metric:          "c",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
				{
					ID:              "ID01",
					ItemDescription: "aaaaa",
					Metric:          "a",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
			},
			expected: []model.OracleDatabaseLicenseType{
				{
					ID:              "ID01",
					ItemDescription: "aaaaa",
					Metric:          "a",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
				{
					ID:              "ID01",
					ItemDescription: "",
					Metric:          "b",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
				{
					ID:              "ID01",
					ItemDescription: "",
					Metric:          "c",
					Cost:            0,
					Aliases:         []string{"Pippo"},
					Option:          false,
				},
			},
		},
	}

	for _, tc := range testCases {
		sort.Slice(tc.licenseTypes, licenseTypesSorter(tc.config, tc.environment, tc.licenseTypes))

		assert.Equal(t, tc.licenseTypes, tc.expected)
	}
}
