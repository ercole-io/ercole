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

package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestThrowNewAlert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("BadRequest", func(t *testing.T) {
		as := NewMockAlertServiceInterface(mockCtrl)
		ac := AlertQueueController{
			TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     utils.NewLogger("TEST"),
		}

		alert := model.Alert{
			AlertCategory:           "pippo",
			AlertAffectedTechnology: new(string),
			AlertCode:               "",
			AlertSeverity:           "",
			AlertStatus:             "",
			Description:             "",
			Date:                    time.Time{},
			OtherInfo:               map[string]interface{}{},
		}
		alertBytes, _ := json.Marshal(alert)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.ThrowNewAlert)

		req, err := http.NewRequest("POST", "/alerts", bytes.NewReader(alertBytes))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Success_NoContent", func(t *testing.T) {
		as := NewMockAlertServiceInterface(mockCtrl)
		ac := AlertQueueController{
			TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     utils.NewLogger("TEST"),
		}

		alert := model.Alert{
			AlertCategory:           model.AlertCategoryEngine,
			AlertAffectedTechnology: nil,
			AlertCode:               model.AlertCodeMissingPrimaryDatabase,
			AlertSeverity:           model.AlertSeverityWarning,
			AlertStatus:             model.AlertStatusNew,
			Description:             "",
			Date:                    time.Time{},
			OtherInfo: map[string]interface{}{
				"hostname": "pippo",
				"database": "pluto",
			},
		}

		as.EXPECT().ThrowNewAlert(alert)

		alertBytes, _ := json.Marshal(alert)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.ThrowNewAlert)

		req, err := http.NewRequest("POST", "/alerts", bytes.NewReader(alertBytes))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)
	})
}
