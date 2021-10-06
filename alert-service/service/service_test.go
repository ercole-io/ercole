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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/leandro-lugaresi/hub"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestProcessMsg_AlertInsertion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	emailer := NewMockEmailer(mockCtrl)

	as := AlertService{
		Emailer: emailer,
		TimeNow: utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:     logger.NewLogger("TEST"),
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
		"CRITICAL This is just an alert test to a mocked emailer. on TestHostname",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: CRITICAL
Host: TestHostname
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To)

	fields := make(hub.Fields, 1)
	fields["alert"] = model.Alert{
		AlertCategory:           model.AlertCategoryLicense,
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertStatus:             model.AlertStatusNew,
		OtherInfo:               map[string]interface{}{"hostname": "TestHostname"},
		AlertSeverity:           model.AlertSeverityCritical,
		Description:             "This is just an alert test to a mocked emailer.",
		Date:                    utils.P("2019-09-02T10:25:28Z"),
		AlertCode:               model.AlertCodeNewLicense,
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
		Log: logger.NewLogger("TEST"),
	}

	msg := hub.Message{
		Name: "",
	}

	as.ProcessMsg(msg)
}

func TestProcessAlertInsertion_WithHostname(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	emailer := NewMockEmailer(mockCtrl)

	as := AlertService{
		Emailer: emailer,
		TimeNow: utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:     logger.NewLogger("TEST"),
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
		"CRITICAL This is just an alert test to a mocked emailer. on TestHostname",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: CRITICAL
Host: TestHostname
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To)

	params := make(hub.Fields, 1)
	params["alert"] = model.Alert{
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		OtherInfo:               map[string]interface{}{"hostname": "TestHostname"},
		AlertSeverity:           model.AlertSeverityCritical,
		Description:             "This is just an alert test to a mocked emailer.",
		Date:                    utils.P("2019-09-02T10:25:28Z"),
		AlertCode:               model.AlertCodeNewLicense,
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
		Log:     logger.NewLogger("TEST"),
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
		"CRITICAL This is just an alert test to a mocked emailer.",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: CRITICAL
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To)

	params := make(hub.Fields, 1)
	params["alert"] = model.Alert{
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		OtherInfo:               map[string]interface{}{},
		AlertSeverity:           model.AlertSeverityCritical,
		Description:             "This is just an alert test to a mocked emailer.",
		Date:                    utils.P("2019-09-02T10:25:28Z"),
		AlertCode:               model.AlertCodeNewLicense,
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
		Log:     logger.NewLogger("TEST"),
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
		"CRITICAL This is just an alert test to a mocked emailer.",
		`Date: 2019-09-02 10:25:28 +0000 UTC
Severity: CRITICAL
Code: NEW_LICENSE
This is just an alert test to a mocked emailer.`,
		as.Config.AlertService.Emailer.To).
		Return(utils.NewError(fmt.Errorf("test error from emailer"), "test EMAILER"))

	params := make(hub.Fields, 1)
	params["alert"] = model.Alert{
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		OtherInfo:               map[string]interface{}{},
		AlertSeverity:           model.AlertSeverityCritical,
		Description:             "This is just an alert test to a mocked emailer.",
		Date:                    utils.P("2019-09-02T10:25:28Z"),
		AlertCode:               model.AlertCodeNewLicense,
	}

	as.ProcessAlertInsertion(params)
}
