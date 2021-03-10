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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//TODO TestSearchDatabases

func TestGetDatabasesStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	expected := &dto.DatabasesStatistics{
		TotalMemorySize:   42.42,
		TotalSegmentsSize: 53.53,
	}
	as.EXPECT().GetDatabasesStatistics(filter).
		Return(expected, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetDatabasesStatistics)
	req, err := http.NewRequest("GET", "/stats?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(expected), rr.Body.String())
}

func TestGetDatabasesStats_GlobalFilter_Error(t *testing.T) {
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
	handler := http.HandlerFunc(ac.GetDatabasesStatistics)
	req, err := http.NewRequest("GET", "/stats?older-than=xxxx-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expected := utils.ErrorResponseFE{
		Error:            "Unable to parse string to time.Time",
		ErrorDescription: "parsing time \"xxxx-06-10T11:54:59Z\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"xxxx-06-10T11:54:59Z\" as \"2006\"",
	}

	var actual utils.ErrorResponseFE
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected.Error, actual.Error)
	assert.Equal(t, expected.ErrorDescription, actual.ErrorDescription)
}

func TestGetDatabasesStats_Service_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	as.EXPECT().GetDatabasesStatistics(filter).
		Return(nil, errMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetDatabasesStatistics)
	req, err := http.NewRequest("GET", "/stats?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expected := utils.ErrorResponseFE{
		ErrorDescription: "MockError",
	}
	assert.Equal(t, utils.ToJSON(expected), rr.Body.String())
}
