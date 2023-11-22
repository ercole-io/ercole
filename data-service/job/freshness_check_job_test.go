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

package job

import (
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestFreshnessCheckJobRun_SuccessNoOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database:       db,
		AlertSvcClient: nil,
		Log:            logger.NewLogger("TEST"),
		NewObjectID:    utils.NewObjectIDForTests(),
	}

	pippo := []model.HostDataBE{
		{
			Hostname:  "pippohost",
			CreatedAt: utils.P("2019-10-05T14:02:03Z"),
		},
	}

	db.EXPECT().DeleteAllNoDataAlerts().Return(nil).Times(1)

	db.EXPECT().GetActiveHostdata().Return(pippo, nil)

	db.EXPECT().FindOldCurrentHostdata(pippo[0].Hostname, utils.P("2019-11-04T14:02:03Z")).Return(false, nil)

	fcj.Run()
}

func TestFreshnessCheckJobRun_SuccessTwoOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	now := utils.Btc(utils.P("2019-11-05T14:02:03Z"))

	fcj := FreshnessCheckJob{
		TimeNow:        now,
		Database:       db,
		AlertSvcClient: asc,
		Log:            logger.NewLogger("TEST"),
		NewObjectID:    utils.NewObjectIDForTests(),
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

	hosts := []model.HostDataBE{
		{
			Hostname:  pippo.Hostname,
			CreatedAt: pippo.CreatedAt,
		},
		{
			Hostname:  pluto.Hostname,
			CreatedAt: pluto.CreatedAt,
		},
	}

	db.EXPECT().GetActiveHostdata().Return(hosts, nil)

	db.EXPECT().FindOldCurrentHostdata(pippo.Hostname, utils.P("2019-11-04T14:02:03Z")).
		Return(true, nil)

	db.EXPECT().FindOldCurrentHostdata(pluto.Hostname, utils.P("2019-11-04T14:02:03Z")).
		Return(true, nil)

	alert1 := model.Alert{
		ID:                      utils.Str2oid("000000000000000000000001"),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryAgent,
		AlertCode:               model.AlertCodeNoData,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    now(),
		Description: fmt.Sprintf("No data received from the host %s in the last %d day(s)",
			"pippohost", 31),
		OtherInfo: map[string]interface{}{
			"hostname": "pippohost",
		},
	}
	asc.EXPECT().ThrowNewAlert(alert1).Return(nil).Times(1)

	alert2 := model.Alert{
		ID:                      utils.Str2oid("000000000000000000000002"),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryAgent,
		AlertCode:               model.AlertCodeNoData,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    now(),
		Description: fmt.Sprintf("No data received from the host %s in the last %d day(s)",
			"plutohost", 21),
		OtherInfo: map[string]interface{}{
			"hostname": "plutohost",
		},
	}
	asc.EXPECT().ThrowNewAlert(alert2).Return(nil).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_DeleteAllNoDataAlertsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database:       db,
		AlertSvcClient: nil,
		Log:            logger.NewLogger("TEST"),
		NewObjectID:    utils.NewObjectIDForTests(),
	}

	db.EXPECT().DeleteAllNoDataAlerts().Return(aerrMock).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_FindOldCurrentHostdataError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	fcj := FreshnessCheckJob{
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database:       db,
		AlertSvcClient: nil,
		Log:            logger.NewLogger("TEST"),
		NewObjectID:    utils.NewObjectIDForTests(),
	}

	pippo := []model.HostDataBE{
		{
			Hostname:  "pippohost",
			CreatedAt: utils.P("2019-10-05T14:02:03Z"),
		},
	}

	db.EXPECT().GetActiveHostdata().Return(pippo, nil)
	db.EXPECT().DeleteAllNoDataAlerts().Return(nil).Times(1)
	db.EXPECT().FindOldCurrentHostdata(gomock.Any(), gomock.Any()).Return(true, aerrMock).Times(1)

	fcj.Run()
}

func TestFreshnessCheckJobRun_ThrowNoDataAlertError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	now := utils.Btc(utils.P("2019-11-05T14:02:03Z"))

	fcj := FreshnessCheckJob{
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database:       db,
		AlertSvcClient: asc,
		Log:            logger.NewLogger("TEST"),
		NewObjectID:    utils.NewObjectIDForTests(),
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

	hosts := []model.HostDataBE{
		{
			Hostname:  pippo.Hostname,
			CreatedAt: pippo.CreatedAt,
		},
		{
			Hostname:  pluto.Hostname,
			CreatedAt: pluto.CreatedAt,
		},
	}

	db.EXPECT().GetActiveHostdata().Return(hosts, nil)

	db.EXPECT().FindOldCurrentHostdata(pippo.Hostname, utils.P("2019-11-04T14:02:03Z")).
		Return(true, nil)

	db.EXPECT().FindOldCurrentHostdata(pluto.Hostname, utils.P("2019-11-04T14:02:03Z")).
		Return(true, nil)

	alert1 := model.Alert{
		ID:                      utils.Str2oid("000000000000000000000001"),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryAgent,
		AlertCode:               model.AlertCodeNoData,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    now(),
		Description: fmt.Sprintf("No data received from the host %s in the last %d day(s)",
			"pippohost", 31),
		OtherInfo: map[string]interface{}{
			"hostname": "pippohost",
		},
	}
	asc.EXPECT().ThrowNewAlert(alert1).Return(aerrMock).Times(1)

	alert2 := model.Alert{
		ID:                      utils.Str2oid("000000000000000000000002"),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryAgent,
		AlertCode:               model.AlertCodeNoData,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    now(),
		Description: fmt.Sprintf("No data received from the host %s in the last %d day(s)",
			"plutohost", 21),
		OtherInfo: map[string]interface{}{
			"hostname": "plutohost",
		},
	}
	asc.EXPECT().ThrowNewAlert(alert2).Return(aerrMock).Times(1)

	fcj.Run()
}
