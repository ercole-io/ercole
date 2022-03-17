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

package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetHostsHistory_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	host := []dto.HostCores{}
	location := "Italy"
	environment := "TST"
	olderThan := utils.MAX_TIME
	newerThan := utils.MIN_TIME

	as.EXPECT().GetHostCores(location, environment, olderThan, newerThan).
		Return(host, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/?location=Italy&environment=TST", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expected := map[string]interface{}{
		"coresHistory": host,
	}
	assert.JSONEq(t, utils.ToJSON(expected), rr.Body.String())
}

func TestGetHostsHistoryUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/?older-than=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHostsHistoryUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/?newer-than=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHostsHistory_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	location := ""
	environment := ""
	olderThan := utils.MAX_TIME
	newerThan := utils.MIN_TIME

	as.EXPECT().GetHostCores(location, environment, olderThan, newerThan).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
