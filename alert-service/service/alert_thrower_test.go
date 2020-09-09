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

	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//ThrowNewDatabaseAlert tests

func TestThrowNewDatabaseAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Queue:    hub.New(),
		Log:      utils.NewLogger("TEST"),
		Config: config.Configuration{
			AlertService: config.AlertService{
				LogAlertThrows: true},
		},
	}

	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCategoryLicense, alert.AlertCategory)
		assert.Equal(t, model.TechnologyOracleDatabase, *alert.AlertAffectedTechnology)
		assert.Equal(t, model.AlertCodeNewDatabase, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityInfo, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "The database 'bestdb' was created on the server myhost", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
			"dbname":   "bestdb",
		}, alert.OtherInfo)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	}).Times(1)

	require.NoError(t, as.ThrowNewDatabaseAlert("bestdb", "myhost"))
}

func TestThrowNewDatabaseAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)
	assert.Equal(t, aerrMock, as.ThrowNewDatabaseAlert("bestdb", "myhost"))
}

//ThrowNewServerAlert tests

func TestThrowNewServerAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Queue:    hub.New(),
		Log:      utils.NewLogger("TEST"),
		Config: config.Configuration{
			AlertService: config.AlertService{
				LogAlertThrows: true},
		},
	}

	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCategoryEngine, alert.AlertCategory)
		assert.Nil(t, alert.AlertAffectedTechnology)
		assert.Equal(t, model.AlertCodeNewServer, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityInfo, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "The server 'myhost' was added to ercole", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
		}, alert.OtherInfo)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	}).Times(1)

	require.NoError(t, as.ThrowNewServerAlert("myhost"))
}

func TestThrowNewServerAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)
	assert.Equal(t, aerrMock, as.ThrowNewServerAlert("myhost"))
}

// ThrowNewEnterpriseLicenseAlert tests

func TestThrowNewEnterpriseLicenseAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Queue:    hub.New(),
		Log:      utils.NewLogger("TEST"),
		Config: config.Configuration{
			AlertService: config.AlertService{
				LogAlertThrows: true},
		},
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCategoryLicense, alert.AlertCategory)
		assert.Equal(t, model.TechnologyOracleDatabase, *alert.AlertAffectedTechnology)
		assert.Equal(t, model.AlertCodeNewLicense, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityCritical, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "A new Enterprise license has been enabled to myhost", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
		}, alert.OtherInfo)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	}).Times(1)
	require.NoError(t, as.ThrowNewEnterpriseLicenseAlert("myhost"))
}

func TestThrowNewEnterpriseLicenseAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)
	assert.Equal(t, aerrMock, as.ThrowNewEnterpriseLicenseAlert("myhost"))
}

// ThrowActivatedFeaturesAlert tests

func TestThrowActivatedFeaturesAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Queue:    hub.New(),
		Log:      utils.NewLogger("TEST"),
		Config: config.Configuration{
			AlertService: config.AlertService{
				LogAlertThrows: true},
		},
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCategoryLicense, alert.AlertCategory)
		assert.Equal(t, model.TechnologyOracleDatabase, *alert.AlertAffectedTechnology)
		assert.Equal(t, model.AlertCodeNewOption, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityCritical, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "The database mydb on myhost has enabled new features (fastibility, slowibility) on server", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
			"dbname":   "mydb",
			"features": []string{"fastibility", "slowibility"},
		}, alert.OtherInfo)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	})
	require.NoError(t, as.ThrowActivatedFeaturesAlert("mydb", "myhost", []string{"fastibility", "slowibility"}))
}

func TestThrowActivatedFeaturesAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)
	assert.Equal(t, aerrMock, as.ThrowActivatedFeaturesAlert("mydb", "myhost", []string{"fastibility", "slowibility"}))
}

// ThrowNoDataAlert tests

func TestThrowNoDataAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Queue:    hub.New(),
		Log:      utils.NewLogger("TEST"),
		Config: config.Configuration{
			AlertService: config.AlertService{
				LogAlertThrows: true},
		},
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCategoryAgent, alert.AlertCategory)
		assert.Nil(t, alert.AlertAffectedTechnology)
		assert.Equal(t, model.AlertCodeNoData, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityCritical, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "No data received from the host myhost in the last 90 day(s)", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
		}, alert.OtherInfo)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	}).Times(1)
	require.NoError(t, as.ThrowNoDataAlert("myhost", 90))
}

func TestThrowNoDataAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)
	assert.Equal(t, aerrMock, as.ThrowNoDataAlert("myhost", 90))
}

func TestThrowUnlistedRunningDatabasesAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Queue:    hub.New(),
		Log:      utils.NewLogger("TEST"),
		Config: config.Configuration{
			AlertService: config.AlertService{
				LogAlertThrows: true},
		},
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCategoryEngine, alert.AlertCategory)
		assert.Equal(t, model.TechnologyOracleDatabase, *alert.AlertAffectedTechnology)
		assert.Equal(t, model.AlertCodeUnlistedRunningDatabase, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityWarning, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "The database mydb is not listed in the oratab of the host myhost", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
			"dbname":   "mydb",
		}, alert.OtherInfo)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	}).Times(1)
	require.NoError(t, as.ThrowUnlistedRunningDatabasesAlert("mydb", "myhost"))
}

func TestThrowUnlistedRunningDatabasesAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)
	assert.Equal(t, aerrMock, as.ThrowUnlistedRunningDatabasesAlert("mydb", "myhost"))
}
