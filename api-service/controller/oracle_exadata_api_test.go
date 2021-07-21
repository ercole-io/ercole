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
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearchOracleExadata_SuccessPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	var resFromService2 = []dto.OracleExadata{
		{
			Id:        "5eba60d00b606515fdc2c554",
			CreatedAt: utils.P("2020-05-12T08:39:44.831Z"),
			DbServers: []dto.DbServers{
				{
					Hostname:           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
					Memory:             376,
					Model:              "X7-2",
					RunningCPUCount:    48,
					RunningPowerSupply: 2,
					SwVersion:          "19.2.4.0.0.190709",
					TempActual:         24,
					TotalCPUCount:      48,
					TotalPowerSupply:   2,
				},
				{
					Hostname:           "kantoor-43a6cdc54bb211eb127bca5c6651950c",
					Memory:             376,
					Model:              "X7-2",
					RunningCPUCount:    48,
					RunningPowerSupply: 2,
					SwVersion:          "19.2.4.0.0.190709",
					TempActual:         24,
					TotalCPUCount:      48,
					TotalPowerSupply:   2,
				},
			},
			Environment: "PROD",
			Hostname:    "test-exadata",
			IbSwitches: []dto.IbSwitches{
				{
					Hostname:  "off-df8b95a01746a464e69203c840a6a46a",
					Model:     "SUN_DCS_36p",
					SwVersion: "2.2.13-2.190326",
				},
				{
					Hostname:  "aspen-8d1d1b210625b1f1024b686135f889a1",
					Model:     "SUN_DCS_36p",
					SwVersion: "2.2.13-2.190326",
				},
			},
			Location: "Italy",
			StorageServers: []dto.StorageServers{
				{
					Hostname:           "s75-c2449b0e89e5a0b38401636eaa07abd5",
					Memory:             188,
					Model:              "X7-2L_High_Capacity",
					RunningCPUCount:    20,
					RunningPowerSupply: 2,
					SwVersion:          "19.2.4.0.0.190709",
					TempActual:         23,
					TotalCPUCount:      40,
					TotalPowerSupply:   2,
				},
				{
					Hostname:           "itl-b22fa37cad1326aba990cdec7facace2",
					Memory:             188,
					Model:              "X7-2L_High_Capacity",
					RunningCPUCount:    20,
					RunningPowerSupply: 2,
					SwVersion:          "19.2.4.0.0.190709",
					TempActual:         24,
					TotalCPUCount:      40,
					TotalPowerSupply:   2,
				},
			},
		},
	}

	var resFromService = dto.OracleExadataResponse{
		Content: resFromService2,
		Metadata: dto.PagingMetadata{
			Empty:         false,
			First:         true,
			Last:          true,
			Number:        0,
			Size:          1,
			TotalElements: 1,
			TotalPages:    0,
		},
	}

	as.EXPECT().
		SearchOracleExadata(true, "foobar", "Hostname", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(&resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?full=true&search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(&resFromService), rr.Body.String())
}

func TestSearchOracleExadata_SuccessUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	var exo = dto.OracleExadataResponse{
		Content: []dto.OracleExadata{
			{
				Id:             "",
				CreatedAt:      time.Time{},
				DbServers:      nil,
				Environment:    "",
				Hostname:       "",
				IbSwitches:     nil,
				Location:       "",
				StorageServers: nil,
			},
		},
		Metadata: dto.PagingMetadata{},
	}

	as.EXPECT().
		SearchOracleExadata(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(&exo, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(&exo.Content), rr.Body.String())
}

func TestSearchOracleExadata_FailUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?full=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?sort-desc=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?page=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?size=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity5(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?older-than=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailInternalServerError(t *testing.T) {
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
		SearchOracleExadata(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleExadataXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	xlsx := excelize.File{}

	as.EXPECT().
		SearchOracleExadataAsXLSX(filter).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
}

func TestSearchOracleExadataXLSX_InternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	as.EXPECT().
		SearchOracleExadataAsXLSX(filter).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadataXLSX)
	req, err := http.NewRequest("GET", "/?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleExadataXLSX_StatusBadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadataXLSX)
	req, err := http.NewRequest("GET", "/?location=Italy&environment=TST&older-than=sdaaadsd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
