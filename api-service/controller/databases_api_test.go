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

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/plandem/xlsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAddms_JSONPaged(t *testing.T) {
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
		"Content": []interface{}{
			map[string]interface{}{
				"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
				"Benefit":        83.34,
				"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "SQL Tuning",
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
			},
			map[string]interface{}{
				"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
				"Benefit":        68.24,
				"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "Segment Tuning",
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
			},
		},
		"Metadata": map[string]interface{}{
			"Empty":         false,
			"First":         true,
			"Last":          true,
			"Number":        0,
			"Size":          20,
			"TotalElements": 25,
			"TotalPages":    1,
		},
	}

	resFromService := []map[string]interface{}{
		expectedRes,
	}

	as.EXPECT().
		SearchAddms("foobar", "Benefit", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?search=foobar&sort-by=Benefit&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchAddms_JSONUnpaged(t *testing.T) {
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
		map[string]interface{}{
			"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"Benefit":        83.34,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
		map[string]interface{}{
			"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"Benefit":        68.24,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchAddms("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchAddms_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?sort-desc=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?page=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?size=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?older-than=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_JSONInternalServerError(t *testing.T) {
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
		SearchAddms("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchAddms_XLSXSuccess(t *testing.T) {
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

	expectedRes := []map[string]interface{}{
		map[string]interface{}{
			"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"Benefit":        83.34,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
		map[string]interface{}{
			"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"Benefit":        68.24,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchAddms("foobar", "Benefit", true, -1, -1, "Germany", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?search=foobar&location=Germany&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Addm")
	require.NotNil(t, sh)
	assert.Equal(t, "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".", sh.Cell(0, 1).String())
	AssertXLSXFloat(t, 83.34, sh.Cell(1, 1))
	assert.Equal(t, "ERCOLE", sh.Cell(2, 1).String())
	assert.Equal(t, "TST", sh.Cell(3, 1).String())
	assert.Equal(t, "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.", sh.Cell(4, 1).String())
	assert.Equal(t, "test-db", sh.Cell(5, 1).String())
	assert.Equal(t, "SQL Tuning", sh.Cell(6, 1).String())

	assert.Equal(t, "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.", sh.Cell(0, 2).String())
	AssertXLSXFloat(t, 68.24, sh.Cell(1, 2))
	assert.Equal(t, "ERCOLE", sh.Cell(2, 2).String())
	assert.Equal(t, "TST", sh.Cell(3, 2).String())
	assert.Equal(t, "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.", sh.Cell(4, 2).String())
	assert.Equal(t, "test-db", sh.Cell(5, 2).String())
	assert.Equal(t, "Segment Tuning", sh.Cell(6, 2).String())
}

func TestSearchAddms_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms?older-than=aasdasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAddms_XLSXInternalServerError1(t *testing.T) {
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
		SearchAddms("", "Benefit", true, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchAddms_XLSXInternalServerError2(t *testing.T) {
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
		map[string]interface{}{
			"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"Benefit":        83.34,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
		map[string]interface{}{
			"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"Benefit":        68.24,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchAddms("", "Benefit", true, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
