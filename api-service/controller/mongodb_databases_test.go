// Copyright (c) 2023 Sorint.lab S.p.A.
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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestSearchMongoDBInstances_JSONPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	resContent := []dto.MongoDBInstance{
		{
			Hostname:     "test-db",
			Environment:  "TST",
			Location:     "Italy",
			InstanceName: "",
			DBName:       "",
			Charset:      "UTF8",
			Version:      "6.0.1",
		},
	}

	var resFromService = dto.MongoDBInstanceResponse{
		Content: resContent,
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
		SearchMongoDBInstances(
			dto.SearchMongoDBInstancesFilter{
				GlobalFilter: dto.GlobalFilter{
					Location: "Italy", Environment: "TST", OlderThan: utils.P("2020-06-10T11:54:59Z"),
				},
				Search: "foobar", SortBy: "Hostname", SortDesc: true, PageNumber: 2, PageSize: 3,
			}).
		Return(&resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?full=true&search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(&resFromService), rr.Body.String())
}

func TestSearchMongoDBInstances_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	resContent := []dto.MongoDBInstance{
		{
			Hostname:     "test-db",
			Environment:  "TST",
			Location:     "Italy",
			InstanceName: "",
			DBName:       "",
			Charset:      "UTF8",
			Version:      "6.0.1",
		},
		{
			Hostname:     "test-db2",
			Environment:  "TST",
			Location:     "Italy",
			InstanceName: "",
			DBName:       "",
			Charset:      "UTF8",
			Version:      "6.0.1",
		},
	}

	var resFromService = dto.MongoDBInstanceResponse{
		Content:  resContent,
		Metadata: dto.PagingMetadata{},
	}

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().
		SearchMongoDBInstances(
			dto.SearchMongoDBInstancesFilter{
				GlobalFilter: dto.GlobalFilter{
					Location: "", Environment: "", OlderThan: utils.MAX_TIME,
				},
				Search: "", SortBy: "", SortDesc: false, PageNumber: -1, PageSize: -1,
			},
		).
		Return(&resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(&resFromService.Content), rr.Body.String())
}

func TestSearchMongoDBInstances_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?sort-desc=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchMongoDBInstances_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?page=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchMongoDBInstances_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?size=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchMongoDBInstances_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?older-than=sasdasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchMongoDBInstances_JSONInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().
		SearchMongoDBInstances(
			dto.SearchMongoDBInstancesFilter{
				GlobalFilter: dto.GlobalFilter{
					Location: "", Environment: "", OlderThan: utils.MAX_TIME,
				},
				Search: "", SortBy: "", SortDesc: false, PageNumber: -1, PageSize: -1,
			},
		).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchMongoDBInstances_XLSXSuccess(t *testing.T) {
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

	expectedRes := excelize.NewFile()

	as.EXPECT().
		SearchMongoDBInstancesAsXLSX(
			dto.SearchMongoDBInstancesFilter{
				GlobalFilter: dto.GlobalFilter{
					Location: "Italy", Environment: "TST", OlderThan: utils.P("2020-06-10T11:54:59Z"),
				},
				Search: "foobar", SortBy: "Hostname", SortDesc: true, PageNumber: -1, PageSize: -1,
			},
		).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?search=foobar&sort-by=Hostname&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
}

func TestSearchMongoDBInstances_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?sort-desc=sdddaadasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchMongoDBInstances_XLSXUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases?older-than=sdddaadasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchMongoDBInstances_XLSXInternalServerError1(t *testing.T) {
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

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().
		SearchMongoDBInstancesAsXLSX(
			dto.SearchMongoDBInstancesFilter{
				GlobalFilter: dto.GlobalFilter{
					Location: "", Environment: "", OlderThan: utils.MAX_TIME,
				},
				Search: "", SortBy: "", SortDesc: false, PageNumber: -1, PageSize: -1,
			},
		).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchMongoDBInstances)
	req, err := http.NewRequest("GET", "/databases", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
