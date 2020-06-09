// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestProcessMsg_HostDataInsertion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:      utils.NewLogger("TEST"),
		Queue:    hub.New(),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", utils.P("2019-11-05T14:02:03Z")).Return(emptyHostData, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, "The server 'superhost1' was added to ercole", alert.Description)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
	}).Times(1)

	msg := hub.Message{
		Name: "hostdata.insertion",
		Fields: hub.Fields{
			"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
		},
	}
	as.ProcessMsg(msg)
}

func TestProcessMsg_AlertInsertion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	emailer := NewMockEmailer(mockCtrl)

	as := AlertService{
		Emailer: emailer,
		TimeNow: utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:     utils.NewLogger("TEST"),
		Queue:   hub.New(),
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					To: []string{"test@ercole.test"},
				},
			},
		},
	}

	emailer.EXPECT().SendEmail(
		"MAJOR This is just an alert test to a mocked emailer. on TestHostname",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: MAJOR
Host: TestHostname
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To)

	fields := make(hub.Fields, 1)
	fields["alert"] = model.Alert{
		AlertCategory:      model.AlertCategoryLicense,
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertStatus:        model.AlertStatusNew,
		OtherInfo:          map[string]interface{}{"Hostname": "TestHostname"},
		AlertSeverity:      model.AlertSeverityMajor,
		Description:        "This is just an alert test to a mocked emailer.",
		Date:               utils.P("2019-09-02T10:25:28Z"),
		AlertCode:          model.AlertCodeNewLicense,
	}

	msg := hub.Message{
		Name:   "alert.insertion",
		Fields: fields,
	}
	as.ProcessMsg(msg)
}

func TestProcessMsg_WrongInsertion(t *testing.T) {
	as := AlertService{
		Config: config.Configuration{
			AlertService: config.AlertService{LogMessages: true}},
		Log: utils.NewLogger("TEST"),
	}

	msg := hub.Message{
		Name: "",
	}

	as.ProcessMsg(msg)
}

func TestProcessHostDataInsertion_SuccessNewHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:      utils.NewLogger("TEST"),
		Queue:    hub.New(),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", utils.P("2019-11-05T14:02:03Z")).Return(emptyHostData, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, "The server 'superhost1' was added to ercole", alert.Description)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, model.AlertCategorySystem, alert.AlertCategory)
		assert.Nil(t, alert.AlertAffectedAsset)
	}).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:      utils.NewLogger("TEST"),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(emptyHostData, aerrMock).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:      utils.NewLogger("TEST"),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", utils.P("2019-11-05T14:02:03Z")).Return(model.HostData{}, aerrMock).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DiffHostError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:      utils.NewLogger("TEST"),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", utils.P("2019-11-05T14:02:03Z")).Return(emptyHostData, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessAlertInsertion_WithHostname(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	emailer := NewMockEmailer(mockCtrl)

	as := AlertService{
		Emailer: emailer,
		TimeNow: utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:     utils.NewLogger("TEST"),
		Queue:   hub.New(),
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					To: []string{"test@ercole.test"},
				},
			},
		},
	}

	emailer.EXPECT().SendEmail(
		"MAJOR This is just an alert test to a mocked emailer. on TestHostname",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: MAJOR
Host: TestHostname
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To)

	params := make(hub.Fields, 1)
	params["alert"] = model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		OtherInfo:          map[string]interface{}{"Hostname": "TestHostname"},
		AlertSeverity:      model.AlertSeverityMajor,
		Description:        "This is just an alert test to a mocked emailer.",
		Date:               utils.P("2019-09-02T10:25:28Z"),
		AlertCode:          model.AlertCodeNewLicense,
	}

	as.ProcessAlertInsertion(params)
}

func TestProcessAlertInsertion_WithoutHostname(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	emailer := NewMockEmailer(mockCtrl)

	as := AlertService{
		Emailer: emailer,
		TimeNow: utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:     utils.NewLogger("TEST"),
		Queue:   hub.New(),
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					To: []string{"test@ercole.test"},
				},
			},
		},
	}

	emailer.EXPECT().SendEmail(
		"MAJOR This is just an alert test to a mocked emailer.",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: MAJOR
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To)

	params := make(hub.Fields, 1)
	params["alert"] = model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		OtherInfo:          map[string]interface{}{},
		AlertSeverity:      model.AlertSeverityMajor,
		Description:        "This is just an alert test to a mocked emailer.",
		Date:               utils.P("2019-09-02T10:25:28Z"),
		AlertCode:          model.AlertCodeNewLicense,
	}

	as.ProcessAlertInsertion(params)
}

func TestProcessAlertInsertion_EmailerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	emailer := NewMockEmailer(mockCtrl)

	as := AlertService{
		Emailer: emailer,
		TimeNow: utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:     utils.NewLogger("TEST"),
		Queue:   hub.New(),
		Config: config.Configuration{
			AlertService: config.AlertService{
				Emailer: config.Emailer{
					To: []string{"test@ercole.test"},
				},
			},
		},
	}

	emailer.EXPECT().SendEmail(
		"MAJOR This is just an alert test to a mocked emailer.",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: MAJOR
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To).
		Return(utils.NewAdvancedErrorPtr(fmt.Errorf("test error from emailer"), "test EMAILER"))

	params := make(hub.Fields, 1)
	params["alert"] = model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		OtherInfo:          map[string]interface{}{},
		AlertSeverity:      model.AlertSeverityMajor,
		Description:        "This is just an alert test to a mocked emailer.",
		Date:               utils.P("2019-09-02T10:25:28Z"),
		AlertCode:          model.AlertCodeNewLicense,
	}

	as.ProcessAlertInsertion(params)
}

func TestDiffHostDataMapAndGenerateAlert_SuccessNoDifferences(t *testing.T) {
	as := AlertService{
		Log: utils.NewLogger("TEST"),
	}

	require.NoError(t, as.DiffHostDataMapAndGenerateAlert(hostData2, hostData1))
}

func TestDiffHostDataMapAndGenerateAlert_SuccessNewHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:      utils.NewLogger("TEST"),
		Queue:    hub.New(),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: nil,
		AlertCategory:      model.AlertCategorySystem,
		AlertCode:          model.AlertCodeNewServer,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.NoError(t, as.DiffHostDataMapAndGenerateAlert(emptyHostData, hostData1))
}

func TestDiffHostDataMapAndGenerateAlert_SuccessNewDatabase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:      utils.NewLogger("TEST"),
		Queue:    hub.New(),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		AlertCode:          model.AlertCodeNewDatabase,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.NoError(t, as.DiffHostDataMapAndGenerateAlert(hostData1, hostData3))
}

func TestDiffHostDataMapAndGenerateAlert_SuccessNewEnterpriseLicense(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:      utils.NewLogger("TEST"),
		Queue:    hub.New(),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		AlertCode:          model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		AlertCode:          model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
			"Features": []string{"Driving"},
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.NoError(t, as.DiffHostDataMapAndGenerateAlert(hostData3, hostData4))
}

func TestDiffHostDataMapAndGenerateAlert_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:      utils.NewLogger("TEST"),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: nil,
		AlertCategory:      model.AlertCategorySystem,
		AlertCode:          model.AlertCodeNewServer,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataMapAndGenerateAlert(emptyHostData, hostData1))
}

func TestDiffHostDataMapAndGenerateAlert_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:      utils.NewLogger("TEST"),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		AlertCode:          model.AlertCodeNewDatabase,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataMapAndGenerateAlert(hostData1, hostData3))
}

func TestDiffHostDataMapAndGenerateAlert_DatabaseError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:      utils.NewLogger("TEST"),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		AlertCode:          model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataMapAndGenerateAlert(hostData3, hostData4))
}

func TestDiffHostDataMapAndGenerateAlert_DatabaseError4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:      utils.NewLogger("TEST"),
		Queue:    hub.New(),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		AlertCode:          model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertAffectedAsset: model.AssetOracleDatabasePtr,
		AlertCategory:      model.AlertCategoryLicense,
		AlertCode:          model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
			"Features": []string{"Driving"},
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataMapAndGenerateAlert(hostData3, hostData4))
}
