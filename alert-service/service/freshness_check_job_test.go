package service

import (
	"testing"
	"time"

	"github.com/amreo/ercole-services/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFreshnessCheckJobRun_SuccessNoOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      btc(p("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}

	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return([]string{}, nil).Do(func(tm time.Time) {
		assert.Equal(t, p("2019-10-26T14:02:03Z"), tm)
	}).Times(1)
	db.EXPECT().ExistNoDataAlertByHost(gomock.Any()).Times(0)
	as.EXPECT().ThrowNoDataAlert(gomock.Any(), gomock.Any()).Times(0)

	fcj.Run()
}

func TestFreshnessCheckJobRun_SuccessTwoOldCurrentHostsWithoutNoDataAlert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      btc(p("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}

	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return([]string{"pippohost", "plutohost"}, nil).Do(func(tm time.Time) {
		assert.Equal(t, p("2019-10-26T14:02:03Z"), tm)
	})
	db.EXPECT().ExistNoDataAlertByHost("pippohost").Return(false, nil).Times(1)
	db.EXPECT().ExistNoDataAlertByHost("plutohost").Return(false, nil).Times(1)
	db.EXPECT().ExistNoDataAlertByHost(gomock.Any()).Times(0)

	as.EXPECT().ThrowNoDataAlert("pippohost", 10).Return(nil).Times(1)
	as.EXPECT().ThrowNoDataAlert("plutohost", 10).Return(nil).Times(1)
	as.EXPECT().ThrowNoDataAlert(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	fcj.Run()
}

func TestFreshnessCheckJobRun_SuccessTwoOldCurrentHostsWithNoDataAlert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      btc(p("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}

	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return([]string{"pippohost", "plutohost"}, nil).Do(func(tm time.Time) {
		assert.Equal(t, p("2019-10-26T14:02:03Z"), tm)
	})
	db.EXPECT().ExistNoDataAlertByHost("pippohost").Return(true, nil).Times(1)
	db.EXPECT().ExistNoDataAlertByHost("plutohost").Return(true, nil).Times(1)
	db.EXPECT().ExistNoDataAlertByHost(gomock.Any()).Times(0)

	as.EXPECT().ThrowNoDataAlert(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	fcj.Run()
}
func TestFreshnessCheckJobRun_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      btc(p("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}

	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return(nil, aerrMock).Do(func(tm time.Time) {
		assert.Equal(t, p("2019-10-26T14:02:03Z"), tm)
	}).Times(1)
	db.EXPECT().ExistNoDataAlertByHost(gomock.Any()).Times(0)
	as.EXPECT().ThrowNoDataAlert(gomock.Any(), gomock.Any()).Times(0)

	fcj.Run()
}

func TestFreshnessCheckJobRun_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      btc(p("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}

	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return([]string{"pippohost", "plutohost"}, nil).Do(func(tm time.Time) {
		assert.Equal(t, p("2019-10-26T14:02:03Z"), tm)
	}).Times(1)

	db.EXPECT().ExistNoDataAlertByHost("pippohost").Return(false, aerrMock).Times(1)
	db.EXPECT().ExistNoDataAlertByHost(gomock.Any()).Times(0)

	fcj.Run()
}

func TestFreshnessCheckJobRun_AlertServiceError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      btc(p("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}

	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return([]string{"pippohost", "plutohost"}, nil).Do(func(tm time.Time) {
		assert.Equal(t, p("2019-10-26T14:02:03Z"), tm)
	}).Times(1)

	db.EXPECT().ExistNoDataAlertByHost("pippohost").Return(false, nil).Times(1)
	db.EXPECT().ExistNoDataAlertByHost(gomock.Any()).Times(0)
	as.EXPECT().ThrowNoDataAlert("pippohost", 10).Return(aerrMock).Times(1)
	as.EXPECT().ThrowNoDataAlert(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	fcj.Run()
}
