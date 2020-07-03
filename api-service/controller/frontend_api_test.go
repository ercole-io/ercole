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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInfoForFrontendDashboard_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	res := map[string]interface{}{
		"Alerts": []map[string]interface{}{
			{
				"AffectedHosts": 4,
				"Code":          "NEW_DATABASE",
				"Count":         9,
				"OldestAlert":   "2020-05-11T11:38:11.992+02:00",
				"Severity":      "INFO",
			},
			{
				"AffectedHosts": 12,
				"Code":          "NEW_SERVER",
				"Count":         12,
				"OldestAlert":   "2020-05-11T11:38:11.988+02:00",
				"Severity":      "INFO",
			},
			{
				"AffectedHosts": 4,
				"Code":          "NEW_OPTION",
				"Count":         7,
				"OldestAlert":   "2020-05-11T11:38:11.992+02:00",
				"Severity":      "CRITICAL",
			},
		},
		"Technologies": map[string]interface{}{
			"Technologies": []map[string]interface{}{
				{
					"Compliance": false,
					"Cost":       0,
					"Count":      0,
					"Name":       model.TechnologyOracleDatabase,
					"Used":       8,
				},
				{
					"Compliance": true,
					"Cost":       0,
					"Count":      2,
					"Name":       model.TechnologyOracleExadata,
					"Used":       2,
				},
			},
			"Total": map[string]interface{}{
				"Compliant": false,
				"Cost":      0,
				"Count":     2,
				"Used":      10,
			},
		},
		"Features": map[string]interface{}{
			"Oracle/Database": true,
			"Oracle/Exadata":  true,
		},
	}

	as.EXPECT().
		GetInfoForFrontendDashboard("Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(res, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetInfoForFrontendDashboard)
	req, err := http.NewRequest("GET", "/settings/features?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(res), rr.Body.String())
}

func TestGetInfoForFrontendDashboard_UnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetInfoForFrontendDashboard)
	req, err := http.NewRequest("GET", "/settings/features?older-than=sdfsdfsdf", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetInfoForFrontendDashboard_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		GetInfoForFrontendDashboard("", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetInfoForFrontendDashboard)
	req, err := http.NewRequest("GET", "/settings/features", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
