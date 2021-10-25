// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestSearchAlerts_JSONSuccessPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"content": []interface{}{
			map[string]interface{}{
				"affectedHosts": 12,
				"code":          "NEW_SERVER",
				"count":         12,
				"oldestAlert":   utils.P("2020-05-06T15:40:04.543+02:00"),
				"severity":      "INFO",
			},
			map[string]interface{}{
				"affectedHosts": 1,
				"code":          "NEW_LICENSE",
				"count":         1,
				"oldestAlert":   utils.P("2020-05-06T15:40:04.62+02:00"),
				"severity":      "CRITICAL",
			},
		},
		"metadata": map[string]interface{}{
			"empty":         false,
			"first":         true,
			"last":          true,
			"number":        0,
			"ize":           20,
			"totalElements": 25,
			"totalPages":    1,
		},
	}

	resFromService := []map[string]interface{}{
		expectedRes,
	}

	as.EXPECT().
		SearchAlerts("aggregated-code-severity", "foo", "CreatedAt", true, 10, 2, "", "", model.AlertSeverityCritical, model.AlertStatusAck, utils.P("2020-06-10T11:54:59Z"), utils.P("2020-06-17T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?mode=aggregated-code-severity&search=foo&sort-by=CreatedAt&sort-desc=true&page=10&size=2&severity=CRITICAL&status=ACK&from=2020-06-10T11%3A54%3A59Z&to=2020-06-17T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchAlerts_JSONSuccessUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"affectedHosts": 12,
			"code":          "NEW_SERVER",
			"count":         12,
			"oldestAlert":   utils.P("2020-05-06T15:40:04.543+02:00"),
			"severity":      "INFO",
		},
		{
			"affectedHosts": 1,
			"code":          "NEW_LICENSE",
			"count":         1,
			"oldestAlert":   utils.P("2020-05-06T15:40:04.62+02:00"),
			"severity":      "CRITICAL",
		},
	}
	as.EXPECT().
		SearchAlerts("aggregated-code-severity",
			"foo", "CreatedAt", true, -1, -1, "", "", model.AlertSeverityCritical, model.AlertStatusAck, utils.P("2020-06-10T11:54:59Z"), utils.P("2020-06-17T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?mode=aggregated-code-severity&search=foo&sort-by=CreatedAt&sort-desc=true&severity=CRITICAL&status=ACK&from=2020-06-10T11%3A54%3A59Z&to=2020-06-17T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchAlerts_JSONFailUnprocessable1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?mode=sdfgsdfg", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailUnprocessable2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?sort-desc=maybe", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailUnprocessable3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?page=ssdfds", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailUnprocessable4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?size=asasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailUnprocessable5(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?severity=asasdsd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailUnprocessable6(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?status=asasdsd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailUnprocessable7(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?from=asasdsd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailUnprocessable8(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?to=asasdsd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_JSONFailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchAlerts("all", "", "", false, -1, -1, "", "", "", "", utils.MIN_TIME, utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchAlertsAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.MAX_TIME,
	}

	from := utils.P("2020-06-10T11:54:59Z")
	to := utils.P("2020-06-17T11:54:59Z")
	xlsx := excelize.File{}

	as.EXPECT().
		SearchAlertsAsXLSX(from, to, filter).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?location=Italy&environment=TST&from=2020-06-10T11%3A54%3A59Z&to=2020-06-17T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	_, err = excelize.OpenReader(rr.Body)

	require.NoError(t, err)
}

func TestSearchAlertsXLSX_UnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?sort-desc=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlertsXLSX_UnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?page=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlertsXLSX_UnprocessableEntity3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?size=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlertsXLSX_UnprocessableEntity4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?severity=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlertsXLSX_UnprocessableEntity5(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?status=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlertsXLSX_UnprocessableEntity6(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?from=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlertsXLSX_UnprocessableEntity7(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?to=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlertsXLSX_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	from := utils.P("2020-06-10T11:54:59Z")
	to := utils.P("2020-06-17T11:54:59Z")

	as.EXPECT().
		SearchAlertsAsXLSX(from, to, filter).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?from=2020-06-10T11%3A54%3A59Z&to=2020-06-17T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAckAlerts_FailForbidden(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: true,
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AckAlerts)
	body := map[string]interface{}{
		"ids": []string{"5dc3f534db7e81a98b726a52"},
	}
	req, err := http.NewRequest("POST", "/alerts/acks", bytes.NewReader([]byte(utils.ToJSON(body))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAckAlerts_FailBadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AckAlerts)
	body := map[string]interface{}{
		"ids": []string{"asdasd"},
	}
	req, err := http.NewRequest("POST", "/alerts/acks", bytes.NewReader([]byte(utils.ToJSON(body))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAckAlerts_FailNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	a := dto.AlertsFilter{
		IDs: []primitive.ObjectID{utils.Str2oid("5dc3f534db7e81a98b726a52")},
	}

	as.EXPECT().AckAlerts(a).
		Return(utils.ErrAlertNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AckAlerts)
	body := struct {
		Filter dto.AlertsFilter
	}{
		Filter: a,
	}
	req, err := http.NewRequest("POST", "/alerts/acks", bytes.NewReader([]byte(utils.ToJSON(body))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAckAlerts_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	a := dto.AlertsFilter{
		IDs: []primitive.ObjectID{utils.Str2oid("5dc3f534db7e81a98b726a52")},
	}

	as.EXPECT().AckAlerts(a).
		Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AckAlerts)
	body := struct {
		Filter dto.AlertsFilter
	}{
		Filter: a,
	}

	req, err := http.NewRequest("POST", "/alerts/acks", bytes.NewReader([]byte(utils.ToJSON(body))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAckAlerts_Success(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		as := NewMockAPIServiceInterface(mockCtrl)
		ac := APIController{
			TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			Service: as,
			Config: config.Configuration{
				APIService: config.APIService{
					ReadOnly: false,
				},
			},
			Log: logger.NewLogger("TEST"),
		}

		s := model.AlertStatusNew
		a := dto.AlertsFilter{
			IDs:         []primitive.ObjectID{utils.Str2oid("000000000000")},
			AlertStatus: &s,
			Date:        time.Time{},
			OtherInfo: map[string]interface{}{
				"host": "pippo",
			},
		}
		as.EXPECT().AckAlerts(a).Return(nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.AckAlerts)
		body := struct {
			Filter dto.AlertsFilter
		}{
			Filter: a,
		}
		req, err := http.NewRequest("POST", "/alerts/acks", bytes.NewReader([]byte(utils.ToJSON(body))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)
	})
}
