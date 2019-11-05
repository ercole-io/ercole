package service

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amreo/ercole-services/model"
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
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCodeNewDatabase, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityNotice, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "The database 'bestdb' was created on the server myhost", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
			"dbname":   "bestdb",
		}, alert.OtherInfo)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	})
	require.Nil(t, as.ThrowNewDatabaseAlert("bestdb", "myhost"))
}

func TestThrowNewDatabaseAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock)
	assert.Equal(t, aerrMock, as.ThrowNewDatabaseAlert("bestdb", "myhost"))
}

//ThrowNewServerAlert tests

func TestThrowNewServerAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCodeNewServer, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityNotice, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "The server 'myhost' was added to ercole", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
		}, alert.OtherInfo)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	})
	require.Nil(t, as.ThrowNewServerAlert("myhost"))
}

func TestThrowNewServerAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock)
	assert.Equal(t, aerrMock, as.ThrowNewServerAlert("myhost"))
}

// ThrowNewEnterpriseLicenseAlert tests

func TestThrowNewEnterpriseLicenseAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCodeNewLicense, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityCritical, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "A new Enterprise license has been enabled to myhost", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
		}, alert.OtherInfo)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	})
	require.Nil(t, as.ThrowNewEnterpriseLicenseAlert("myhost"))
}

func TestThrowNewEnterpriseLicenseAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock)
	assert.Equal(t, aerrMock, as.ThrowNewEnterpriseLicenseAlert("myhost"))
}

// ThrowActivatedFeaturesAlert tests

func TestThrowActivatedFeaturesAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCodeNewOption, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityCritical, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "The database mydb on myhost has enabled new features (fastibility, slowibility) on server", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
			"dbname":   "mydb",
			"features": []string{"fastibility", "slowibility"},
		}, alert.OtherInfo)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	})
	require.Nil(t, as.ThrowActivatedFeaturesAlert("mydb", "myhost", []string{"fastibility", "slowibility"}))
}

func TestThrowActivatedFeaturesAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock)
	assert.Equal(t, aerrMock, as.ThrowActivatedFeaturesAlert("mydb", "myhost", []string{"fastibility", "slowibility"}))
}

// ThrowNoDataAlert tests

func TestThrowNoDataAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, model.AlertCodeNoData, alert.AlertCode)
		assert.Equal(t, model.AlertSeverityMajor, alert.AlertSeverity)
		assert.Equal(t, model.AlertStatusNew, alert.AlertStatus)
		assert.Equal(t, "No data received from the host myhost in the last 90 days", alert.Description)
		assert.Equal(t, map[string]interface{}{
			"hostname": "myhost",
		}, alert.OtherInfo)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.Date)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.ID.Timestamp())
	})
	require.Nil(t, as.ThrowNoDataAlert("myhost", 90))
}

func TestThrowNoDataAlert_DatabaseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock)
	assert.Equal(t, aerrMock, as.ThrowNoDataAlert("myhost", 90))
}
