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
	"fmt"
	reflect "reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"

	dto "github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	alertsc := NewMockAlertSvcClientInterface(mockCtrl)
	apisc := NewMockApiSvcClientInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}},
			},
			DataService: config.DataService{
				LogInsertingHostdata: true,
			},
		},
		AlertSvcClient: alertsc,
		ApiSvcClient:   apisc,
		ServerVersion:  "",
		Log:            logger.NewLogger("TEST"),
	}

	hdPrimary := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_22.json")
	hdSecondary := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_23.json")

	primaryDB := hdSecondary.Features.Oracle.Database.Databases[0]
	assert.True(t, utils.Contains(model.OracleDatabaseStatusMounted, primaryDB.Status) && primaryDB.Role != model.OracleDatabaseRolePrimary)

	apisc.EXPECT().GetOracleDatabases()
	apisc.EXPECT().AckAlerts(dto.AlertsFilter{
		AlertCategory: utils.Str2ptr(model.AlertCategoryEngine),
		AlertCode:     utils.Str2ptr(model.AlertCodeMissingPrimaryDatabase),
		AlertSeverity: utils.Str2ptr(model.AlertSeverityWarning),
		OtherInfo: map[string]interface{}{
			"hostname": "itl-csllab-112.sorint.localpippo",
			"dbname":   "ERCOLE",
		},
	})
	alertsc.EXPECT().ThrowNewAlert(&alertSimilarTo{
		al: model.Alert{
			AlertCategory: model.AlertCategoryEngine,
			AlertCode:     model.AlertCodeMissingPrimaryDatabase,
			AlertSeverity: model.AlertSeverityWarning,
			OtherInfo: map[string]interface{}{
				"hostname": "itl-csllab-112.sorint.localpippo",
				"dbname":   "ERCOLE",
			},
		}}).Return(nil)

	hds.addLicensesToSecondaryDb(hdPrimary.Info, 2, &primaryDB)
}

var hostData1 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dc3f534db7e81a98b726a52"),
	Hostname:  "superhost1",
	Archived:  false,
	CreatedAt: utils.P("2019-11-05T14:02:03Z"),
	Features: model.Features{
		Oracle: &model.OracleFeature{
			Database: &model.OracleDatabaseFeature{
				MissingDatabases: []model.MissingDatabase{},
				Databases:        []model.OracleDatabase{},
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
				MissingDatabases: []model.MissingDatabase{},
				Databases:        []model.OracleDatabase{},
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
				MissingDatabases: []model.MissingDatabase{},
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
				MissingDatabases: []model.MissingDatabase{},
				Databases: []model.OracleDatabase{
					{
						Name: "acd",
						Licenses: []model.OracleDatabaseLicense{
							{
								LicenseTypeID: "Oracle ENT",
								Name:          "Oracle ENT",
								Count:         10,
								Ignored:       true,
							},
							{
								LicenseTypeID: "Driving",
								Name:          "Driving",
								Count:         100,
								Ignored:       false,
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
				MissingDatabases: []model.MissingDatabase{},
				Databases: []model.OracleDatabase{
					{
						Name: "acd",
						Licenses: []model.OracleDatabaseLicense{
							{
								LicenseTypeID: "Oracle ENT",
								Name:          "Oracle ENT",
								Count:         10,
								Ignored:       true,
							},
							{
								LicenseTypeID: "Driving",
								Name:          "Driving",
								Count:         100,
								Ignored:       false,
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
								Ignored:       true,
							},
							{
								LicenseTypeID: "Driving",
								Name:          "Driving",
								Count:         50,
								Ignored:       false,
							},
							{
								LicenseTypeID: "Waving",
								Name:          "Waving",
								Count:         50,
								Ignored:       false,
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

var hostData6 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		Oracle: &model.OracleFeature{
			Database: &model.OracleDatabaseFeature{
				MissingDatabases: []model.MissingDatabase{},
				Databases: []model.OracleDatabase{
					{
						Name: "acd",
						Licenses: []model.OracleDatabaseLicense{
							{
								LicenseTypeID: "Oracle ENT",
								Name:          "Oracle ENT",
								Count:         10,
								Ignored:       false,
							},
						},
					},
					{
						Name:      "Bart",
						Dataguard: true,
						IsRAC:     false,
						Status:    model.OracleDatabaseStatusMounted[0],
						Role:      model.OracleDatabaseRoleSnapshotStandby,
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

func TestCheckNewLicenses_SuccessNoDifferences(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	hds.checkNewLicenses(&hostData2, &hostData1, licenseTypes)
}

func TestCheckNewLicenses_SuccessNewDatabase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
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

	hds.checkNewLicenses(&hostData1, &hostData3, licenseTypes)
}

func TestCheckNewLicenses_ThrowNewLicenseAndNewOption(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
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

	hds.checkNewLicenses(&hostData3, &hostData4, licenseTypes)
}

func TestCheckNewLicenses_ThrowNewLicenseAndNewOptionAlreadyEnabled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
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

	hds.checkNewLicenses(&hostData4, &hostData5, licenseTypes)
}

func TestCheckNewLicenses_CantThrowNewAlert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
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

		hds.checkNewLicenses(&hostData1, &hostData3, licenseTypes)
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

		hds.checkNewLicenses(&hostData3, &hostData4, licenseTypes)
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

		hds.checkNewLicenses(&hostData3, &hostData4, licenseTypes)
	})
}

func TestCheckNewLicenses_ErrOracleDatabaseLicenseTypeIDNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	hds.checkNewLicenses(&hostData3, &hostData4, []model.OracleDatabaseLicenseType{})
}

func TestIgnorePreviousLicences_SuccessNoPreviousIgnored(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	hds.ignorePreviousLicences(&hostData6, &hostData6)

	for _, db := range hostData6.Features.Oracle.Database.Databases {
		for _, license := range db.Licenses {
			if license.LicenseTypeID == "Oracle ENT" {
				if license.Ignored {
					t.Fatal("unexpected ignored license")
				}
			}
		}
	}
}

func TestIgnorePreviousLicences_SuccessWithPreviousIgnored(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	hds.ignorePreviousLicences(&hostData4, &hostData6)

	for _, db := range hostData6.Features.Oracle.Database.Databases {
		for _, license := range db.Licenses {
			if license.LicenseTypeID == "Oracle ENT" {
				if !license.Ignored {
					t.Fatal("expected ignored license")
				}
			}
		}
	}
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

func TestCheckMissingDatabases_NoneMissing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	alertsc := NewMockAlertSvcClientInterface(mockCtrl)
	apisc := NewMockApiSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: alertsc,
		ApiSvcClient:   apisc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	f := dto.AlertsFilter{
		AlertCategory:           utils.Str2ptr(model.AlertCategoryLicense),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCode:               utils.Str2ptr(model.AlertCodeMissingDatabase),
		AlertStatus:             utils.Str2ptr(model.AlertStatusNew),
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}
	apisc.EXPECT().GetAlertsByFilter(f)

	hds.checkMissingDatabases(&hostData3, &hostData4)
}

func TestCheckMissingDatabases_OneMissing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	alertsc := NewMockAlertSvcClientInterface(mockCtrl)
	apisc := NewMockApiSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: alertsc,
		ApiSvcClient:   apisc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	f := dto.AlertsFilter{
		AlertCategory:           utils.Str2ptr(model.AlertCategoryLicense),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCode:               utils.Str2ptr(model.AlertCodeMissingDatabase),
		AlertStatus:             utils.Str2ptr(model.AlertStatusNew),
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}
	apisc.EXPECT().GetAlertsByFilter(f)

	alertsc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(actualAlert model.Alert) {
		expectedAlert := model.Alert{
			AlertCategory:           model.AlertCategoryLicense,
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCode:               model.AlertCodeMissingDatabase,
			AlertSeverity:           model.AlertSeverityWarning,
			AlertStatus:             model.AlertStatusNew,
			Description:             "The databases \"second\" on \"superhost1\" are missing compared to the previous hostdata",
			Date:                    utils.P("2019-11-05T16:02:03Z"),
			OtherInfo: map[string]interface{}{
				"hostname": "superhost1",
				"dbNames":  []string{"second"},
			},
		}

		assert.Equal(t, expectedAlert, actualAlert)
	})

	hdPrevious := model.HostDataBE{
		ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
		Hostname:  "superhost1",
		Archived:  true,
		CreatedAt: utils.P("2019-11-05T16:02:03Z"),
		Features: model.Features{
			Oracle: &model.OracleFeature{
				Database: &model.OracleDatabaseFeature{
					MissingDatabases: []model.MissingDatabase{},
					Databases: []model.OracleDatabase{
						{
							Name:     "first",
							Licenses: []model.OracleDatabaseLicense{},
						},
						{
							Name:     "second",
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

	hdNew := model.HostDataBE{
		ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
		Hostname:  "superhost1",
		Archived:  true,
		CreatedAt: utils.P("2019-11-05T16:02:03Z"),
		Features: model.Features{
			Oracle: &model.OracleFeature{
				Database: &model.OracleDatabaseFeature{
					MissingDatabases: []model.MissingDatabase{},
					Databases: []model.OracleDatabase{
						{
							Name:     "first",
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

	hds.checkMissingDatabases(&hdPrevious, &hdNew)
}

func TestCheckMissingDatabases_AllMissing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	alertsc := NewMockAlertSvcClientInterface(mockCtrl)
	apisc := NewMockApiSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: alertsc,
		ApiSvcClient:   apisc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	f := dto.AlertsFilter{
		AlertCategory:           utils.Str2ptr(model.AlertCategoryLicense),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCode:               utils.Str2ptr(model.AlertCodeMissingDatabase),
		AlertStatus:             utils.Str2ptr(model.AlertStatusNew),
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}
	apisc.EXPECT().GetAlertsByFilter(f)

	alertsc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(actualAlert model.Alert) {
		expectedAlert := model.Alert{
			AlertCategory:           model.AlertCategoryLicense,
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCode:               model.AlertCodeMissingDatabase,
			AlertSeverity:           model.AlertSeverityCritical,
			AlertStatus:             model.AlertStatusNew,
			Description:             "The databases \"first, second\" on \"superhost1\" are missing compared to the previous hostdata",
			Date:                    utils.P("2019-11-05T16:02:03Z"),
			OtherInfo: map[string]interface{}{
				"hostname": "superhost1",
				"dbNames":  []string{"first", "second"},
			},
		}

		assert.Equal(t, expectedAlert, actualAlert)
	})

	hdPrevious := model.HostDataBE{
		ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
		Hostname:  "superhost1",
		Archived:  true,
		CreatedAt: utils.P("2019-11-05T16:02:03Z"),
		Features: model.Features{
			Oracle: &model.OracleFeature{
				Database: &model.OracleDatabaseFeature{
					MissingDatabases: []model.MissingDatabase{},
					Databases: []model.OracleDatabase{
						{
							Name:     "first",
							Licenses: []model.OracleDatabaseLicense{},
						},
						{
							Name:     "second",
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

	hdNew := model.HostDataBE{
		ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
		Hostname:  "superhost1",
		Archived:  true,
		CreatedAt: utils.P("2019-11-05T16:02:03Z"),
		Features: model.Features{
			Oracle: &model.OracleFeature{
				Database: &model.OracleDatabaseFeature{
					MissingDatabases: []model.MissingDatabase{},
					Databases:        []model.OracleDatabase{},
				},
			},
		},
		Info: model.Host{
			CPUCores: 2,
		},
	}

	hds.checkMissingDatabases(&hdPrevious, &hdNew)
}

func TestCheckMissingDatabases_NoFeature(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	alertsc := NewMockAlertSvcClientInterface(mockCtrl)
	apisc := NewMockApiSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: alertsc,
		ApiSvcClient:   apisc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	hd := model.HostDataBE{}
	hds.checkMissingDatabases(&hd, nil)
}

func TestSearchAndAckOldMissingDatabasesAlerts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	alertsc := NewMockAlertSvcClientInterface(mockCtrl)
	apisc := NewMockApiSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: alertsc,
		ApiSvcClient:   apisc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	newDbs := map[string]bool{
		"db1": true,
		"db2": true,
	}

	ids := utils.NewObjectIDForTests()
	alerts := []model.Alert{
		{ID: ids()},
		{
			ID: ids(),
			OtherInfo: map[string]interface{}{
				"dbNames": []interface{}{"db1"},
			},
		},
		{
			ID: ids(),
			OtherInfo: map[string]interface{}{
				"dbNames": []interface{}{"db1", "db2"},
			},
		},
	}

	apisc.EXPECT().GetAlertsByFilter(gomock.Any()).DoAndReturn(func(fff dto.AlertsFilter) ([]model.Alert, error) {
		expectedFilter := dto.AlertsFilter{
			AlertCategory:           utils.Str2ptr(model.AlertCategoryLicense),
			AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
			AlertCode:               utils.Str2ptr(model.AlertCodeMissingDatabase),
			AlertSeverity:           nil,
			AlertStatus:             utils.Str2ptr(model.AlertStatusNew),
			Description:             nil,
			Date:                    time.Time{},
			OtherInfo: map[string]interface{}{
				"hostname": "pippo",
			},
		}
		assert.Equal(t, expectedFilter, fff)

		return alerts, nil
	})

	apisc.EXPECT().AckAlerts(dto.AlertsFilter{IDs: []primitive.ObjectID{alerts[1].ID}})
	apisc.EXPECT().AckAlerts(dto.AlertsFilter{IDs: []primitive.ObjectID{alerts[2].ID}})

	hds.searchAndAckOldMissingDatabasesAlerts("pippo", newDbs)
}

func Test_ignoreRacLicenses_success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					AlertType: config.AlertType{
						NewHost:                    config.Directive{Enable: true},
						NewDatabase:                config.Directive{Enable: true},
						NewLicense:                 config.Directive{Enable: true},
						NewOption:                  config.Directive{Enable: true},
						NewUnlistedRunningDatabase: config.Directive{Enable: true},
						NewHostCpu:                 config.Directive{Enable: true},
						MissingPrimaryDatabase:     config.Directive{Enable: true},
						MissingDatabase:            config.Directive{Enable: true},
						AgentError:                 config.Directive{Enable: true},
						NoData:                     config.Directive{Enable: true},
					}}}},
		ServerVersion:  "",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	host := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_24.json")

	snapHost := &host

	hds.ignoreRacLicenses(&host)

	for _, db := range host.Features.Oracle.Database.Databases {
		if db.Edition() == model.OracleDatabaseEditionStandard {
			for _, license := range db.Licenses {
				if license.IsRAC() {
					assert.Equal(t, license.Ignored, true)
					assert.Equal(t, license.IgnoredComment, "RAC license ignored by Ercole")
				}
			}
		}
	}

	assert.NotEqual(t, host, snapHost)
}
