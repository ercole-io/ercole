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
		"content": []interface{}{
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
			"action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"benefit":        83.34,
			"createdAt":      utils.P("2020-07-01T09:18:03.726+02:00"),
			"dbname":         "ERCOLE",
			"environment":    "TST",
			"finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"hostname":       "test-db",
			"location":       "Germany",
			"recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5efc38ab79f92e4cbf283b13"),
		},
		{
			"action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"benefit":        68.24,
			"createdAt":      utils.P("2020-07-01T09:18:03.726+02:00"),
			"dbname":         "ERCOLE",
			"environment":    "TST",
			"finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"hostname":       "test-db",
			"location":       "Germany",
			"recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5efc38ab79f92e4cbf283b13"),
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
	sp, err := excelize.OpenReader(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".", sp.GetCellValue("Addm", "A2"))
	assert.Equal(t, "83.34", sp.GetCellValue("Addm", "B2"))
	assert.Equal(t, "ERCOLE", sp.GetCellValue("Addm", "C2"))
	assert.Equal(t, "TST", sp.GetCellValue("Addm", "D2"))
	assert.Equal(t, "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.", sp.GetCellValue("Addm", "E2"))
	assert.Equal(t, "test-db", sp.GetCellValue("Addm", "F2"))
	assert.Equal(t, "SQL Tuning", sp.GetCellValue("Addm", "G2"))

	assert.Equal(t, "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.", sp.GetCellValue("Addm", "A3"))
	assert.Equal(t, "68.24", sp.GetCellValue("Addm", "B3"))
	assert.Equal(t, "ERCOLE", sp.GetCellValue("Addm", "C3"))
	assert.Equal(t, "TST", sp.GetCellValue("Addm", "D3"))
	assert.Equal(t, "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.", sp.GetCellValue("Addm", "E3"))
	assert.Equal(t, "test-db", sp.GetCellValue("Addm", "F3"))
	assert.Equal(t, "Segment Tuning", sp.GetCellValue("Addm", "G3"))
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
			"OK": true,
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

func TestSearchOracleDatabaseSegmentAdvisors_JSON(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	segmentAdvisors := []dto.OracleDatabaseSegmentAdvisor{
		{
			// Id:             utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
			Hostname:       "test-db3",
			Location:       "Germany",
			Environment:    "PRD",
			CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
			Dbname:         "foobar4",
			Reclaimable:    534.34,
			SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			SegmentType:    "TABLE",
			PartitionName:  "iyyiuyyoy",
			Recommendation: "32b36a77e7481343ef175483c086859e",
		},
		{
			// _id:            utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
			Hostname:       "test-db3",
			Location:       "Germany",
			Environment:    "PRD",
			CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
			Dbname:         "foobar3",
			Reclaimable:    4.3,
			SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			SegmentType:    "TABLE",
			PartitionName:  "iyyiuyyoy",
			Recommendation: "32b36a77e7481343ef175483c086859e",
		},
	}

	as.EXPECT().
		SearchOracleDatabaseSegmentAdvisors("", "", false, "", "", utils.MAX_TIME).
		Return(segmentAdvisors, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisors)
	req, err := http.NewRequest("GET", "/segment-advisors", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expectedRes := map[string]interface{}{
		"segmentAdvisors": segmentAdvisors,
	}
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
		SearchOracleDatabaseSegmentAdvisors("", "", false, "", "", utils.MAX_TIME).
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

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	xlsx := excelize.File{}

	as.EXPECT().
		SearchOracleDatabaseSegmentAdvisorsAsXLSX(filter).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisorsXLSX)
	req, err := http.NewRequest("GET", "/segment-advisors?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseSegmentAdvisorsXLSX)
	req, err := http.NewRequest("GET", "/segment-advisors?older-than=sadasddasasd", nil)
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

	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	as.EXPECT().
		SearchOracleDatabaseSegmentAdvisorsAsXLSX(filter).
		Return(nil, aerrMock)

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
		"content": []interface{}{
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

	windowTime := utils.P("2019-05-05T14:02:03Z")
	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	xlsx := excelize.File{}

	as.EXPECT().
		SearchOracleDatabasePatchAdvisorsAsXLSX(windowTime, filter).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
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
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisorsXLSX)
	req, err := http.NewRequest("GET", "/patch-advisors?older-than=dsasdasd", nil)
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

	windowTime := utils.P("2019-05-05T14:02:03Z")
	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	as.EXPECT().
		SearchOracleDatabasePatchAdvisorsAsXLSX(windowTime, filter).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabasePatchAdvisors)
	req, err := http.NewRequest("GET", "/patch-advisors", nil)
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
		"content": []interface{}{
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
		SearchOracleDatabases(
			dto.SearchOracleDatabasesFilter{
				dto.GlobalFilter{
					"Italy", "TST", utils.P("2020-06-10T11:54:59Z"),
				},
				true, "foobar", "Hostname", true, 2, 3,
			}).
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
		SearchOracleDatabases(
			dto.SearchOracleDatabasesFilter{
				dto.GlobalFilter{
					"", "", utils.MAX_TIME,
				},
				false, "", "", false, -1, -1,
			},
		).
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
		SearchOracleDatabases(
			dto.SearchOracleDatabasesFilter{
				dto.GlobalFilter{
					"", "", utils.MAX_TIME,
				},
				false, "", "", false, -1, -1,
			},
		).
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

	expectedRes := excelize.NewFile()

	as.EXPECT().
		SearchOracleDatabasesAsXLSX(
			dto.SearchOracleDatabasesFilter{
				dto.GlobalFilter{
					"Italy", "TST", utils.P("2020-06-10T11:54:59Z"),
				},
				false, "foobar", "Hostname", true, -1, -1,
			},
		).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases?search=foobar&sort-by=Hostname&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
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
		SearchOracleDatabasesAsXLSX(
			dto.SearchOracleDatabasesFilter{
				dto.GlobalFilter{
					"", "", utils.MAX_TIME,
				},
				false, "", "", false, -1, -1,
			},
		).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabases)
	req, err := http.NewRequest("GET", "/databases", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabaseUsedLicenses_JSONPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	resFromService := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "LID001",
				DbName:        "erclin5dbx",
				Hostname:      "pippo",
				UsedLicenses:  3,
			},
			{
				LicenseTypeID: "LID002",
				DbName:        "erclin6dbx",
				Hostname:      "pluto",
				UsedLicenses:  42,
			},
		},
		Metadata: dto.PagingMetadata{
			Empty: false, First: true, Last: true, Number: 0, Size: 2, TotalElements: 2, TotalPages: 1,
		},
	}

	t.Run("JSON paged", func(t *testing.T) {
		as.EXPECT().
			SearchOracleDatabaseUsedLicenses("Benefit", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
			Return(&resFromService, nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.SearchOracleDatabaseUsedLicenses)
		req, err := http.NewRequest("GET", "/licenses?sort-by=Benefit&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, utils.ToJSON(resFromService), rr.Body.String())
	})

	t.Run("JSON unpaged", func(t *testing.T) {

		as.EXPECT().
			SearchOracleDatabaseUsedLicenses("", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&resFromService, nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.SearchOracleDatabaseUsedLicenses)
		req, err := http.NewRequest("GET", "/licenses", nil)
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, utils.ToJSON(resFromService.Content), rr.Body.String())
	})
}

func TestSearchOracleDatabaseUsedLicenses_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseUsedLicenses)
	req, err := http.NewRequest("GET", "/licenses?sort-desc=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseUsedLicenses_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseUsedLicenses)
	req, err := http.NewRequest("GET", "/licenses?page=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseUsedLicenses_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseUsedLicenses)
	req, err := http.NewRequest("GET", "/licenses?size=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseUsedLicenses_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseUsedLicenses)
	req, err := http.NewRequest("GET", "/licenses?older-than=sadsas", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseUsedLicenses_JSONInternalServerError(t *testing.T) {
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
		SearchOracleDatabaseUsedLicenses("", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseUsedLicenses)
	req, err := http.NewRequest("GET", "/licenses", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
