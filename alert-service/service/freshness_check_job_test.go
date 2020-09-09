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
	"time"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFreshnessCheckJobRun_SuccessNoOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	db.EXPECT().DeleteAllNoDataAlerts().Return(nil).Times(1)

	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return([]model.HostDataBE{}, nil).Do(func(tm time.Time) {
		assert.Equal(t, utils.P("2019-10-26T14:02:03Z"), tm)
	}).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_SuccessTwoOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)

	fcj := FreshnessCheckJob{
		TimeNow:      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	db.EXPECT().DeleteAllNoDataAlerts().Return(nil).Times(1)

	pippo := model.HostDataBE{
		Hostname:  "pippohost",
		CreatedAt: utils.P("2019-10-05T14:02:03Z"),
	}
	pluto := model.HostDataBE{
		Hostname:  "plutohost",
		CreatedAt: utils.P("2019-10-15T14:02:03Z"),
	}

	db.EXPECT().FindOldCurrentHosts(utils.P("2019-10-26T14:02:03Z")).
		Return([]model.HostDataBE{pippo, pluto}, nil)

	as.EXPECT().ThrowNoDataAlert("pippohost", 31).Return(nil).Times(1)
	as.EXPECT().ThrowNoDataAlert("plutohost", 21).Return(nil).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_DeleteAllNoDataAlertsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	db.EXPECT().DeleteAllNoDataAlerts().Return(aerrMock).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_FindOldCurrentHostsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	db.EXPECT().DeleteAllNoDataAlerts().Return(nil).Times(1)
	db.EXPECT().FindOldCurrentHosts(gomock.Any()).Return(nil, aerrMock).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_ThrowNoDataAlertError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)

	fcj := FreshnessCheckJob{
		TimeNow:      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		alertService: as,
		Database:     db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	db.EXPECT().DeleteAllNoDataAlerts().Return(nil).Times(1)

	pippo := model.HostDataBE{
		Hostname:  "pippohost",
		CreatedAt: utils.P("2019-10-05T14:02:03Z"),
	}
	pluto := model.HostDataBE{
		Hostname:  "plutohost",
		CreatedAt: utils.P("2019-10-15T14:02:03Z"),
	}

	db.EXPECT().FindOldCurrentHosts(utils.P("2019-10-26T14:02:03Z")).
		Return([]model.HostDataBE{pippo, pluto}, nil)

	as.EXPECT().ThrowNoDataAlert("pippohost", 31).Return(aerrMock).Times(1)
	as.EXPECT().ThrowNoDataAlert("plutohost", 21).Return(aerrMock).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_InvalidDaysThresholdValue(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := NewMockAlertServiceInterface(mockCtrl)

	t.Run("DaysThreshold = 0", func(t *testing.T) {
		fcj := FreshnessCheckJob{
			TimeNow:      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			alertService: as,
			Database:     db,
			Config: config.Configuration{
				AlertService: config.AlertService{
					FreshnessCheckJob: config.FreshnessCheckJob{
						DaysThreshold: 0,
					},
				},
			},
			Log: utils.NewLogger("TEST"),
		}
		fcj.Run()
	})

	t.Run("DaysThreshold < 0", func(t *testing.T) {
		fcj := FreshnessCheckJob{
			TimeNow:      utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			alertService: as,
			Database:     db,
			Config: config.Configuration{
				AlertService: config.AlertService{
					FreshnessCheckJob: config.FreshnessCheckJob{
						DaysThreshold: -42,
					},
				},
			},
			Log: utils.NewLogger("TEST"),
		}
		fcj.Run()
	})
}
