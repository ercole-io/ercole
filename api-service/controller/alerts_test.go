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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSearchAlerts_JSONSuccessPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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

func TestSearchAlerts_XLSXSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	res := []map[string]interface{}{
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef82"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_LICENSE",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "A new Enterprise license has been enabled to ercsoldbx",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"hostname": "ercsoldbx",
			},
		},
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef83"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_OPTION",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "The database ERCSOL19 on ercsoldbx has enabled new features (Diagnostics Pack) on server",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"dbname": "ERCSOL19",
				"features": []string{
					"Diagnostics Pack",
				},
				"hostname": "ercsoldbx",
			},
		},
	}
	as.EXPECT().
		SearchAlerts("all",
			"foo", "CreatedAt", true, -1, -1, "", "", model.AlertSeverityCritical, model.AlertStatusAck, utils.P("2020-06-10T11:54:59Z"), utils.P("2020-06-17T11:54:59Z")).
		Return(res, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?search=foo&sort-by=CreatedAt&sort-desc=true&severity=CRITICAL&status=ACK&from=2020-06-10T11%3A54%3A59Z&to=2020-06-17T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := excelize.OpenReader(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "LICENSE", sp.GetCellValue("Alerts", "A2"))
	assert.Equal(t, "2020-07-23 08:01:13.746 +0000 UTC", sp.GetCellValue("Alerts", "B2"))
	assert.Equal(t, "CRITICAL", sp.GetCellValue("Alerts", "C2"))
	assert.Equal(t, "ercsoldbx", sp.GetCellValue("Alerts", "D2"))
	assert.Equal(t, "NEW_LICENSE", sp.GetCellValue("Alerts", "E2"))
	assert.Equal(t, "A new Enterprise license has been enabled to ercsoldbx", sp.GetCellValue("Alerts", "F2"))

	assert.Equal(t, "LICENSE", sp.GetCellValue("Alerts", "A3"))
	assert.Equal(t, "2020-07-23 08:01:13.746 +0000 UTC", sp.GetCellValue("Alerts", "B3"))
	assert.Equal(t, "CRITICAL", sp.GetCellValue("Alerts", "C3"))
	assert.Equal(t, "ercsoldbx", sp.GetCellValue("Alerts", "D3"))
	assert.Equal(t, "NEW_OPTION", sp.GetCellValue("Alerts", "E3"))
	assert.Equal(t, "The database ERCSOL19 on ercsoldbx has enabled new features (Diagnostics Pack) on server", sp.GetCellValue("Alerts", "F3"))
}

func TestSearchAlerts_XLSXUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?sort-desc=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_XLSXUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?page=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_XLSXUnprocessableEntity3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?size=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_XLSXUnprocessableEntity4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?severity=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_XLSXUnprocessableEntity5(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?status=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_XLSXUnprocessableEntity6(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?from=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_XLSXUnprocessableEntity7(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?to=sasa", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAlerts_XLSXInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchAlerts("all",
			"foo", "CreatedAt", true, -1, -1, "", "", model.AlertSeverityCritical, model.AlertStatusAck, utils.P("2020-06-10T11:54:59Z"), utils.P("2020-06-17T11:54:59Z")).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?search=foo&sort-by=CreatedAt&sort-desc=true&severity=CRITICAL&status=ACK&from=2020-06-10T11%3A54%3A59Z&to=2020-06-17T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchAlerts_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "asdsad",
		},
		Log: utils.NewLogger("TEST"),
	}

	res := []map[string]interface{}{
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef82"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_LICENSE",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "A new Enterprise license has been enabled to ercsoldbx",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"hostname": "ercsoldbx",
			},
		},
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef83"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_OPTION",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "The database ERCSOL19 on ercsoldbx has enabled new features (Diagnostics Pack) on server",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"dbname": "ERCSOL19",
				"features": []string{
					"Diagnostics Pack",
				},
				"hostname": "ercsoldbx",
			},
		},
	}
	as.EXPECT().
		SearchAlerts("all",
			"foo", "CreatedAt", true, -1, -1, "", "", model.AlertSeverityCritical, model.AlertStatusAck, utils.P("2020-06-10T11:54:59Z"), utils.P("2020-06-17T11:54:59Z")).
		Return(res, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAlerts)
	req, err := http.NewRequest("GET", "/alerts?search=foo&sort-by=CreatedAt&sort-desc=true&severity=CRITICAL&status=ACK&from=2020-06-10T11%3A54%3A59Z&to=2020-06-17T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAckAlerts_Success(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AckAlerts([]primitive.ObjectID{utils.Str2oid("5dc3f534db7e81a98b726a52")}).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AckAlerts)
	body := map[string]interface{}{
		"ids": []string{"5dc3f534db7e81a98b726a52"},
	}
	req, err := http.NewRequest("POST", "/alerts/acks", bytes.NewReader([]byte(utils.ToJSON(body))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
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
		Log: utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AckAlerts([]primitive.ObjectID{utils.Str2oid("5dc3f534db7e81a98b726a52")}).
		Return(utils.ErrAlertNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AckAlerts)
	body := map[string]interface{}{
		"ids": []string{"5dc3f534db7e81a98b726a52"},
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AckAlerts([]primitive.ObjectID{utils.Str2oid("5dc3f534db7e81a98b726a52")}).
		Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AckAlerts)
	body := map[string]interface{}{
		"ids": []string{"5dc3f534db7e81a98b726a52"},
	}
	req, err := http.NewRequest("POST", "/alerts/acks", bytes.NewReader([]byte(utils.ToJSON(body))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAckAlerts_ByFilter(t *testing.T) {
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
			Log: utils.NewLogger("TEST"),
		}

		s := model.AlertStatusNew
		a := dto.AlertsFilter{
			ID:          utils.Str2oid("000000000000"),
			AlertStatus: &s,
			Date:        time.Time{},
			OtherInfo: map[string]interface{}{
				"host": "pippo",
			},
		}
		as.EXPECT().AckAlertsByFilter(a).Return(nil)

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
