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
	"strings"
	"testing"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/plandem/xlsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchOracleDatabaseAddms_JSONPaged(t *testing.T) {
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
		SearchOracleDatabaseAddms("foobar", "Benefit", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms?search=foobar&sort-by=Benefit&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabaseAddms_JSONUnpaged(t *testing.T) {
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
		{
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
		SearchOracleDatabaseAddms("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabaseAddms_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms?sort-desc=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseAddms_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms?page=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseAddms_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms?size=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseAddms_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms?older-than=sdfdfsdfs", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseAddms_JSONInternalServerError(t *testing.T) {
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
		SearchOracleDatabaseAddms("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabaseAddms_XLSXSuccess(t *testing.T) {
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
		{
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
		{
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
		SearchOracleDatabaseAddms("foobar", "Benefit", true, -1, -1, "Germany", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
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

func TestSearchOracleDatabaseAddms_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms?older-than=aasdasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseAddms_XLSXInternalServerError1(t *testing.T) {
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
		SearchOracleDatabaseAddms("", "Benefit", true, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabaseAddms_XLSXInternalServerError2(t *testing.T) {
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
		{
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
		SearchOracleDatabaseAddms("", "Benefit", true, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAddms)
	req, err := http.NewRequest("GET", "/addms", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_JSONPaged(t *testing.T) {
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
				"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
				"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				"Environment":    "SVIL",
				"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
				"Location":       "Germany",
				"PartitionName":  "",
				"Reclaimable":    "\u003c1",
				"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
				"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
				"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				"SegmentType":    "TABLE",
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
			},
			map[string]interface{}{
				"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"PartitionName":  "iyyiuyyoy",
				"Reclaimable":    "\u003c1",
				"Recommendation": "32b36a77e7481343ef175483c086859e",
				"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				"SegmentType":    "TABLE",
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
		SearchOracleDatabaseSegmentAdvisors("foobar", "Reclaimable", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?search=foobar&sort-by=Reclaimable&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabaseSegmentAdvisors_JSONUnpaged(t *testing.T) {
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
			"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
			"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":    "SVIL",
			"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":       "Germany",
			"PartitionName":  "",
			"Reclaimable":    "\u003c1",
			"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
			"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"PartitionName":  "iyyiuyyoy",
			"Reclaimable":    "\u003c1",
			"Recommendation": "32b36a77e7481343ef175483c086859e",
			"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabaseSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabaseSegmentAdvisors_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?sort-desc=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?page=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?size=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?older-than=asasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_JSONInternalServerError(t *testing.T) {
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
		SearchOracleDatabaseSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_XLSXSuccess(t *testing.T) {
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
		{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
			"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":    "SVIL",
			"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":       "Germany",
			"PartitionName":  "",
			"Reclaimable":    "\u003c1",
			"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
			"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"PartitionName":  "iyyiuyyoy",
			"Reclaimable":    "\u003c1",
			"Recommendation": "32b36a77e7481343ef175483c086859e",
			"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabaseSegmentAdvisors("foobar", "Reclaimable", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?search=foobar&sort-by=Reclaimable&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Segment_Advisor")
	require.NotNil(t, sh)
	assert.Equal(t, "4wcqjn-ecf040bdfab7695ab332aef7401f185c", sh.Cell(0, 1).String())
	assert.Equal(t, "SVIL", sh.Cell(1, 1).String())
	assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sh.Cell(2, 1).String())
	assert.Equal(t, "", sh.Cell(3, 1).String())
	assert.Equal(t, "\u003c1", sh.Cell(4, 1).String())
	assert.Equal(t, "3d7e603f515ed171fc99bdb908f38fb2", sh.Cell(5, 1).String())
	assert.Equal(t, "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3", sh.Cell(6, 1).String())
	assert.Equal(t, "Brittany-424f6a749eef846fa40a1ad1ee3d3674", sh.Cell(7, 1).String())
	assert.Equal(t, "TABLE", sh.Cell(8, 1).String())

	assert.Equal(t, "ERCOLE", sh.Cell(0, 2).String())
	assert.Equal(t, "TST", sh.Cell(1, 2).String())
	assert.Equal(t, "test-db", sh.Cell(2, 2).String())
	assert.Equal(t, "iyyiuyyoy", sh.Cell(3, 2).String())
	assert.Equal(t, "\u003c1", sh.Cell(4, 2).String())
	assert.Equal(t, "32b36a77e7481343ef175483c086859e", sh.Cell(5, 2).String())
	assert.Equal(t, "pasta-973e4d1f937da4d9bc1b092f934ab0ec", sh.Cell(6, 2).String())
	assert.Equal(t, "Brittany-424f6a749eef846fa40a1ad1ee3d3674", sh.Cell(7, 2).String())
	assert.Equal(t, "TABLE", sh.Cell(8, 2).String())
}

func TestSearchOracleDatabaseSegmentAdvisors_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?sort-desc=sadasddasasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_XLSXUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors?older-than=sadasddasasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_XLSXInternalServerError1(t *testing.T) {
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
		SearchOracleDatabaseSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabaseSegmentAdvisors_XLSXInternalServerError2(t *testing.T) {
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
			"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
			"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":    "SVIL",
			"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":       "Germany",
			"PartitionName":  "",
			"Reclaimable":    "\u003c1",
			"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
			"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"PartitionName":  "iyyiuyyoy",
			"Reclaimable":    "\u003c1",
			"Recommendation": "32b36a77e7481343ef175483c086859e",
			"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabaseSegmentAdvisors("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_JSONPaged(t *testing.T) {
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
				"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
				"Date":        utils.P("2012-04-16T02:00:00+02:00"),
				"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				"Dbver":       "11.2.0.3.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Environment": "SVIL",
				"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
				"Location":    "Germany",
				"Status":      "KO",
				"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
			},
			map[string]interface{}{
				"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
				"Date":        utils.P("2012-04-16T02:00:00+02:00"),
				"Dbname":      "ERCOLE",
				"Dbver":       "12.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Environment": "TST",
				"Hostname":    "test-db",
				"Location":    "Germany",
				"Status":      "KO",
				"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
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
		SearchOracleDatabasePatchAdvisors("foobar", "Hostname", true, 2, 3, utils.P("2019-03-05T14:02:03Z"), "Italy", "TST", utils.P("2020-06-10T11:54:59Z"), "KO").
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&window-time=8&status=KO&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabasePatchAdvisors_JSONUnpaged(t *testing.T) {
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
			"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
			"Date":        utils.P("2012-04-16T02:00:00+02:00"),
			"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Dbver":       "11.2.0.3.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "SVIL",
			"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
			"Date":        utils.P("2012-04-16T02:00:00+02:00"),
			"Dbname":      "ERCOLE",
			"Dbver":       "12.2.0.1.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "TST",
			"Hostname":    "test-db",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabasePatchAdvisors("", "", false, -1, -1, utils.P("2019-05-05T14:02:03Z"), "", "", utils.MAX_TIME, "").
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabasePatchAdvisors_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?sort-desc=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?page=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?size=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?window-time=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_JSONUnprocessableEntity5(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?status=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_JSONUnprocessableEntity6(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?older-than=sdasdasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_JSONInternalServerError(t *testing.T) {
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
		SearchOracleDatabasePatchAdvisors("", "", false, -1, -1, utils.P("2019-05-05T14:02:03Z"), "", "", utils.MAX_TIME, "").
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_XLSXSuccess(t *testing.T) {
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
		{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Dbver":       "11.2.0.3.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "SVIL",
			"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "ERCOLE",
			"Dbver":       "12.2.0.1.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "TST",
			"Hostname":    "test-db",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabasePatchAdvisors("foobar", "Hostname", true, -1, -1, utils.P("2019-03-05T14:02:03Z"), "Italy", "TST", utils.P("2020-06-10T11:54:59Z"), "KO").
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?search=foobar&sort-by=Hostname&sort-desc=true&window-time=8&status=KO&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Patch_Advisor")
	require.NotNil(t, sh)
	assert.Equal(t, "PSU 11.2.0.3.2", sh.Cell(0, 1).String())
	assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sh.Cell(1, 1).String())
	assert.Equal(t, "4wcqjn-ecf040bdfab7695ab332aef7401f185c", sh.Cell(2, 1).String())
	assert.Equal(t, "11.2.0.3.0 Enterprise Edition", sh.Cell(3, 1).String())
	assert.Equal(t, utils.P("2012-04-16T00:00:00Z").String(), sh.Cell(4, 1).String())
	assert.Equal(t, "KO", sh.Cell(5, 1).String())

	assert.Equal(t, "PSU 11.2.0.3.2", sh.Cell(0, 2).String())
	assert.Equal(t, "test-db", sh.Cell(1, 2).String())
	assert.Equal(t, "ERCOLE", sh.Cell(2, 2).String())
	assert.Equal(t, "12.2.0.1.0 Enterprise Edition", sh.Cell(3, 2).String())
	assert.Equal(t, utils.P("2012-04-16T00:00:00Z").String(), sh.Cell(4, 2).String())
	assert.Equal(t, "KO", sh.Cell(5, 2).String())
}

func TestSearchOracleDatabasePatchAdvisors_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?sort-desc=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_XLSXUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?window-time=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_XLSXUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?older-than=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_XLSXUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?status=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabasePatchAdvisors_XLSXInternalServerError1(t *testing.T) {
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
		SearchOracleDatabasePatchAdvisors("", "", false, -1, -1, utils.P("2019-05-05T14:02:03Z"), "", "", utils.MAX_TIME, "").
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
func TestSearchOracleDatabasePatchAdvisors_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Dbver":       "11.2.0.3.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "SVIL",
			"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
			"Date":        utils.PDT("2012-04-16T02:00:00+02:00"),
			"Dbname":      "ERCOLE",
			"Dbver":       "12.2.0.1.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "TST",
			"Hostname":    "test-db",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabasePatchAdvisors("foobar", "Hostname", true, -1, -1, utils.P("2019-03-05T14:02:03Z"), "Italy", "TST", utils.P("2020-06-10T11:54:59Z"), "KO").
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?search=foobar&sort-by=Hostname&sort-desc=true&window-time=8&status=KO&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabases_JSONPaged(t *testing.T) {
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
				"ArchiveLogStatus": false,
				"BlockSize":        "8192",
				"CPUCount":         "16",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-04-07T08:52:59.82+02:00"),
				"DatafileSize":     "61",
				"Dataguard":        false,
				"Environment":      "SVIL",
				"HA":               false,
				"Hostname":         "publicitate-36d06ca83eafa454423d2097f4965517",
				"Location":         "Germany",
				"Memory":           4.199,
				"Name":             "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				"RAC":              false,
				"SegmentsSize":     "41",
				"Status":           "OPEN",
				"UniqueName":       "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				"Version":          "11.2.0.3.0 Enterprise Edition",
				"Work":             "1",
				"_id":              utils.Str2oid("5e8c234b24f648a08585bd32"),
			},
			map[string]interface{}{
				"ArchiveLogStatus": false,
				"BlockSize":        "8192",
				"CPUCount":         "2",
				"Charset":          "AL32UTF8",
				"CreatedAt":        utils.P("2020-04-07T08:52:59.872+02:00"),
				"DatafileSize":     "6",
				"Dataguard":        false,
				"Environment":      "TST",
				"HA":               false,
				"Hostname":         "test-db",
				"Location":         "Germany",
				"Memory":           1.484,
				"Name":             "ERCOLE",
				"RAC":              false,
				"SegmentsSize":     "3",
				"Status":           "OPEN",
				"UniqueName":       "ERCOLE",
				"Version":          "12.2.0.1.0 Enterprise Edition",
				"Work":             "1",
				"_id":              utils.Str2oid("5e8c234b24f648a08585bd43"),
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
		SearchOracleDatabases(true, "foobar", "Hostname", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?full=true&search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabases_JSONUnpaged(t *testing.T) {
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
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "16",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.82+02:00"),
			"DatafileSize":     "61",
			"Dataguard":        false,
			"Environment":      "SVIL",
			"HA":               false,
			"Hostname":         "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":         "Germany",
			"Memory":           4.199,
			"Name":             "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"RAC":              false,
			"SegmentsSize":     "41",
			"Status":           "OPEN",
			"UniqueName":       "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Version":          "11.2.0.3.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "2",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.872+02:00"),
			"DatafileSize":     "6",
			"Dataguard":        false,
			"Environment":      "TST",
			"HA":               false,
			"Hostname":         "test-db",
			"Location":         "Germany",
			"Memory":           1.484,
			"Name":             "ERCOLE",
			"RAC":              false,
			"SegmentsSize":     "3",
			"Status":           "OPEN",
			"UniqueName":       "ERCOLE",
			"Version":          "12.2.0.1.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabases(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabases_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?full=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabases_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?sort-desc=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabases_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?page=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabases_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?size=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabases_JSONUnprocessableEntity5(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?older-than=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabases_JSONInternalServerError1(t *testing.T) {
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
		SearchOracleDatabases(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabases_XLSXSuccess(t *testing.T) {
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
		{
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "16",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.82+02:00"),
			"DatafileSize":     "61",
			"Dataguard":        false,
			"Environment":      "SVIL",
			"HA":               false,
			"Hostname":         "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":         "Germany",
			"Memory":           4.199,
			"Name":             "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"RAC":              false,
			"SegmentsSize":     "41",
			"Status":           "OPEN",
			"UniqueName":       "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Version":          "11.2.0.3.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "2",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.872+02:00"),
			"DatafileSize":     "6",
			"Dataguard":        false,
			"Environment":      "TST",
			"HA":               false,
			"Hostname":         "test-db",
			"Location":         "Germany",
			"Memory":           1.484,
			"Name":             "ERCOLE",
			"RAC":              false,
			"SegmentsSize":     "3",
			"Status":           "OPEN",
			"UniqueName":       "ERCOLE",
			"Version":          "12.2.0.1.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabases(false, "foobar", "Hostname", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?search=foobar&sort-by=Hostname&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Databases")
	require.NotNil(t, sh)
	assert.Equal(t, "4wcqjn-ecf040bdfab7695ab332aef7401f185c", sh.Cell(0, 1).String())
	assert.Equal(t, "4wcqjn-ecf040bdfab7695ab332aef7401f185c", sh.Cell(1, 1).String())
	assert.Equal(t, "11.2.0.3.0 Enterprise Edition", sh.Cell(2, 1).String())
	assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sh.Cell(3, 1).String())
	assert.Equal(t, "OPEN", sh.Cell(4, 1).String())
	assert.Equal(t, "SVIL", sh.Cell(5, 1).String())
	assert.Equal(t, "Germany", sh.Cell(6, 1).String())
	assert.Equal(t, "AL32UTF8", sh.Cell(7, 1).String())
	assert.Equal(t, "8192", sh.Cell(8, 1).String())
	assert.Equal(t, "16", sh.Cell(9, 1).String())
	assert.Equal(t, "1", sh.Cell(10, 1).String())
	AssertXLSXFloat(t, 4.199, sh.Cell(11, 1))
	assert.Equal(t, "61", sh.Cell(12, 1).String())
	assert.Equal(t, "41", sh.Cell(13, 1).String())
	AssertXLSXBool(t, false, sh.Cell(14, 1))
	AssertXLSXBool(t, false, sh.Cell(15, 1))
	AssertXLSXBool(t, false, sh.Cell(16, 1))
	AssertXLSXBool(t, false, sh.Cell(17, 1))

	assert.Equal(t, "ERCOLE", sh.Cell(0, 2).String())
	assert.Equal(t, "ERCOLE", sh.Cell(1, 2).String())
	assert.Equal(t, "12.2.0.1.0 Enterprise Edition", sh.Cell(2, 2).String())
	assert.Equal(t, "test-db", sh.Cell(3, 2).String())
	assert.Equal(t, "OPEN", sh.Cell(4, 2).String())
	assert.Equal(t, "TST", sh.Cell(5, 2).String())
	assert.Equal(t, "Germany", sh.Cell(6, 2).String())
	assert.Equal(t, "AL32UTF8", sh.Cell(7, 2).String())
	assert.Equal(t, "8192", sh.Cell(8, 2).String())
	assert.Equal(t, "2", sh.Cell(9, 2).String())
	assert.Equal(t, "1", sh.Cell(10, 2).String())
	AssertXLSXFloat(t, 1.484, sh.Cell(11, 2))
	assert.Equal(t, "6", sh.Cell(12, 2).String())
	assert.Equal(t, "3", sh.Cell(13, 2).String())
	AssertXLSXBool(t, false, sh.Cell(14, 2))
	AssertXLSXBool(t, false, sh.Cell(15, 2))
	AssertXLSXBool(t, false, sh.Cell(16, 2))
	AssertXLSXBool(t, false, sh.Cell(17, 2))
}

func TestSearchOracleDatabases_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?sort-desc=sdddaadasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabases_XLSXUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?older-than=sdddaadasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabases_XLSXInternalServerError1(t *testing.T) {
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
		SearchOracleDatabases(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabases_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "16",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.82+02:00"),
			"DatafileSize":     "61",
			"Dataguard":        false,
			"Environment":      "SVIL",
			"HA":               false,
			"Hostname":         "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":         "Germany",
			"Memory":           4.199,
			"Name":             "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"RAC":              false,
			"SegmentsSize":     "41",
			"Status":           "OPEN",
			"UniqueName":       "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Version":          "11.2.0.3.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "2",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.872+02:00"),
			"DatafileSize":     "6",
			"Dataguard":        false,
			"Environment":      "TST",
			"HA":               false,
			"Hostname":         "test-db",
			"Location":         "Germany",
			"Memory":           1.484,
			"Name":             "ERCOLE",
			"RAC":              false,
			"SegmentsSize":     "3",
			"Status":           "OPEN",
			"UniqueName":       "ERCOLE",
			"Version":          "12.2.0.1.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	as.EXPECT().
		SearchOracleDatabases(false, "foobar", "Hostname", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?search=foobar&sort-by=Hostname&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestListLicenses_JSONPaged(t *testing.T) {
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
				"Compliance": false,
				"Count":      0,
				"Used":       5,
				"_id":        "Oracle ENT",
			},
			map[string]interface{}{
				"Compliance": true,
				"Count":      0,
				"Used":       0,
				"_id":        "Oracle STD",
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

	resFromService := []interface{}{
		expectedRes,
	}

	as.EXPECT().
		ListLicenses(true, "Benefit", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses?full=true&sort-by=Benefit&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestListLicenses_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Compliance": false,
			"Count":      0,
			"Used":       5,
			"_id":        "Oracle ENT",
		},
		map[string]interface{}{
			"Compliance": true,
			"Count":      0,
			"Used":       0,
			"_id":        "Oracle STD",
		},
	}

	as.EXPECT().
		ListLicenses(false, "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestListLicenses_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses?full=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListLicenses_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses?sort-desc=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListLicenses_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses?page=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListLicenses_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses?size=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListLicenses_JSONUnprocessableEntity5(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses?older-than=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListLicenses_JSONInternalServerError(t *testing.T) {
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
		ListLicenses(false, "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLicenses)
	req, err := http.NewRequest("GET", "/licenses", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetLicense_Success(t *testing.T) {
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
		"Compliance": false,
		"Count":      0,
		"Hosts": []interface{}{
			map[string]interface{}{
				"Databases": []interface{}{
					"ERCOLE",
					"urcole",
				},
				"Hostname": "itl-csllab-112.sorint.localpippo",
			},
			map[string]interface{}{
				"Databases": []interface{}{
					"ERCOLE",
				},
				"Hostname": "test-db",
			},
			map[string]interface{}{
				"Databases": []interface{}{
					"rudeboy-fb3160a04ffea22b55555bbb58137f77",
					"007bond-f260462ca34bbd17deeda88f042e42a1",
					"jacket-d4a157354d91bfc68fce6f45546d8f3d",
					"allstate-9a6a2a820a3f61aeb345a834abf40fba",
					"4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				},
				"Hostname": "publicitate-36d06ca83eafa454423d2097f4965517",
			},
		},
		"Used": 5,
		"_id":  "Oracle ENT",
	}

	as.EXPECT().
		GetLicense("Oracle ENT", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetLicense)
	req, err := http.NewRequest("GET", "/licenses/Oracle%20ENT?older-than=2020-06-10T11%3A54%3A59Z", nil)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetLicense_UnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.GetLicense)
	req, err := http.NewRequest("GET", "/licenses/Oracle%20ENT?older-than=asasdas", nil)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetLicense_InternalServerError1(t *testing.T) {
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
		GetLicense("Oracle ENT", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetLicense)
	req, err := http.NewRequest("GET", "/licenses/Oracle%20ENT", nil)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetLicense_NotFound(t *testing.T) {
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
		GetLicense("Oracle ENT", utils.MAX_TIME).
		Return(nil, utils.AerrLicenseNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetLicense)
	req, err := http.NewRequest("GET", "/licenses/Oracle%20ENT", nil)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestSetLicenseCount_Success(t *testing.T) {
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

	payload := []map[string]interface{}{
		{
			"_id":   "Oracle ENT",
			"Count": 10,
		},
		{
			"_id":   "Oracle STD",
			"Count": 20,
		},
	}

	as.EXPECT().SetLicensesCount(payload).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestSetLicensesCount_Success(t *testing.T) {
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

	as.EXPECT().SetLicenseCount("Oracle ENT", 10).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicenseCount)
	req, err := http.NewRequest("PUT", "/licenses/Oracle%20ENT/count", strings.NewReader("10"))
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestSetLicensesCount_FailReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SetLicenseCount)
	req, err := http.NewRequest("PUT", "/licenses/Oracle%20ENT/count", strings.NewReader("10"))
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestSetLicensesCount_FailNotFound(t *testing.T) {
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

	as.EXPECT().SetLicenseCount("Oracle ENT", 10).Return(utils.AerrLicenseNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicenseCount)
	req, err := http.NewRequest("PUT", "/licenses/Oracle%20ENT/count", strings.NewReader("10"))
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestSetLicensesCount_FailInternalServerError(t *testing.T) {
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

	as.EXPECT().SetLicenseCount("Oracle ENT", 10).Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicenseCount)
	req, err := http.NewRequest("PUT", "/licenses/Oracle%20ENT/count", strings.NewReader("10"))
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSetLicensesCount_FailUnprocessableEntity(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SetLicenseCount)
	req, err := http.NewRequest("PUT", "/licenses/Oracle%20ENT/count", strings.NewReader("sdfsdf"))
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSetLicensesCount_FailBadRequest(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SetLicenseCount)
	req, err := http.NewRequest("PUT", "/licenses/Oracle%20ENT/count", &FailingReader{})
	req = mux.SetURLVars(req, map[string]string{
		"name": "Oracle ENT",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
func TestSetLicensesCount_FailedReadOnly(t *testing.T) {
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

	payload := []map[string]interface{}{
		{
			"_id":   "Oracle ENT",
			"Count": 10,
		},
		{
			"_id":   "Oracle STD",
			"Count": 20,
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestSetLicensesCount_FailedInternalServerError(t *testing.T) {
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

	payload := []map[string]interface{}{
		{
			"_id":   "Oracle ENT",
			"Count": 10,
		},
		{
			"_id":   "Oracle STD",
			"Count": 20,
		},
	}

	as.EXPECT().SetLicensesCount(payload).Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSetLicensesCount_FailUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte{100, 200}))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSetLicensesCount_FailUnprocessableEntity2(t *testing.T) {
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

	payload := []map[string]interface{}{
		{
			"Count": 10,
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSetLicensesCount_FailUnprocessableEntity3(t *testing.T) {
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

	payload := []map[string]interface{}{
		{
			"_id":   456546,
			"Count": 10,
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSetLicensesCount_FailUnprocessableEntity4(t *testing.T) {
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

	payload := []map[string]interface{}{
		{
			"_id": "Oracle ENT",
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSetLicensesCount_FailUnprocessableEntity5(t *testing.T) {
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

	payload := []map[string]interface{}{
		{
			"_id":   "Oracle ENT",
			"Count": "ssadsad",
		},
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicensesCount)
	req, err := http.NewRequest("PUT", "/licenses", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}
