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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

//TODO: add SearchHostsFilters tests for SearchHosts!
func TestSearchHosts_JSONPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	resFromService := []dto.HostDataSummary{
		{
			CreatedAt:               time.Time{},
			Hostname:                "sample",
			Location:                "",
			Environment:             "",
			AgentVersion:            "",
			Info:                    model.Host{},
			ClusterMembershipStatus: model.ClusterMembershipStatus{},
			Databases:               map[string][]string{},
		},
	}
	filters := dto.SearchHostsFilters{
		Search:         []string{"foobar"},
		SortBy:         "Hostname",
		SortDesc:       true,
		Location:       "Italy",
		Environment:    "TST",
		OlderThan:      utils.P("2020-06-10T11:54:59Z"),
		PageNumber:     2,
		PageSize:       3,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}
	as.EXPECT().
		GetHostDataSummaries(gomock.Any()).
		DoAndReturn(func(actual dto.SearchHostsFilters) ([]dto.HostDataSummary, error) {
			assert.EqualValues(t, filters, actual)

			return resFromService, nil
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?mode=summary&search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	expectedRes := map[string]interface{}{
		"hosts": resFromService,
	}
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchHosts_JSONUnpaged(t *testing.T) {
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
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Cluster":                       "Angola-1dac9f7418db9b52c259ce4ba087cdb6",
			"CreatedAt":                     utils.P("2020-04-07T08:52:59.844+02:00"),
			"Databases":                     "8888888-d41d8cd98f00b204e9800998ecf8427e",
			"Environment":                   "PROD",
			"Hostname":                      "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemTotal":                      3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"VirtualizationNode":            "suspended-290dce22a939f3868f8f23a6e1f57dd8",
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"Version":                       "1.6.1",
			"HardwareAbstraction":           "VIRT",
			"_id":                           utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Cluster":                       "Puzzait",
			"CreatedAt":                     utils.P("2020-04-07T08:52:59.869+02:00"),
			"Databases":                     "",
			"Environment":                   "PROD",
			"Hostname":                      "test-virt",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemTotal":                      3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"VirtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"Version":                       "1.6.1",
			"HardwareAbstraction":           "VIRT",
			"_id":                           utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	filters := dto.SearchHostsFilters{
		Search:         []string{""},
		SortBy:         "",
		SortDesc:       false,
		Location:       "",
		Environment:    "",
		OlderThan:      utils.MAX_TIME,
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}
	as.EXPECT().
		SearchHosts("full", gomock.Any()).
		DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filters, actual)

			return expectedRes, nil
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchHosts_JSONHostnames(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	returnedRes := []map[string]interface{}{
		{
			"hostname": "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
		},
		{
			"hostname": "test-virt",
		},
	}

	expectedRes := []string{
		"fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
		"test-virt",
	}

	filters := dto.SearchHostsFilters{
		Search:         []string{""},
		SortBy:         "",
		SortDesc:       false,
		Location:       "",
		Environment:    "",
		OlderThan:      utils.MAX_TIME,
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}
	as.EXPECT().
		SearchHosts("hostnames", gomock.Any()).
		DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filters, actual)

			return returnedRes, nil
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?mode=hostnames", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchHosts_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?mode=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?sort-desc=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?page=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?size=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity5(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?older-than=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	filters := dto.SearchHostsFilters{
		Search:         []string{""},
		SortBy:         "",
		SortDesc:       false,
		Location:       "",
		Environment:    "",
		OlderThan:      utils.MAX_TIME,
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}
	as.EXPECT().
		SearchHosts("full", gomock.Any()).
		DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filters, actual)

			return nil, aerrMock
		})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchHosts_LMSSuccess(t *testing.T) {
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

	expected := excelize.NewFile()
	filters := dto.SearchHostsFilters{
		Search:         []string{"foobar"},
		SortBy:         "Processors",
		SortDesc:       true,
		Location:       "Italy",
		Environment:    "TST",
		OlderThan:      utils.P("2020-06-10T11:54:59Z"),
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}
	filterlsm := dto.SearchHostsAsLMS{
		SearchHostsFilters: filters,
		GlobalFilter:       dto.GlobalFilter{"Italy", "TST", utils.P("2020-06-10T11:54:59Z")},
		NewerThan:          utils.P("2021-06-10T11:54:59Z"),
	}
	as.EXPECT().
		SearchHostsAsLMS(gomock.Any()).
		DoAndReturn(func(actual dto.SearchHostsAsLMS) (*excelize.File, error) {
			assert.EqualValues(t, filterlsm, actual)

			return expected, nil
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z&newer-than=2021-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.ms-excel.sheet.macroEnabled.12")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
}

func TestSearchHosts_LMSUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?sort-desc=sdfsdf", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.ms-excel.sheet.macroEnabled.12")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_LMSUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?older-than=sdfsdf", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.ms-excel.sheet.macroEnabled.12")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_LMSSuccessInternalServerError1(t *testing.T) {
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

	filters := dto.SearchHostsFilters{
		Search:         []string{"foobar"},
		SortBy:         "Processors",
		SortDesc:       true,
		Location:       "Italy",
		Environment:    "TST",
		OlderThan:      utils.P("2020-06-10T11:54:59Z"),
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}
	filterlsm := dto.SearchHostsAsLMS{
		SearchHostsFilters: filters,
		GlobalFilter:       dto.GlobalFilter{"Italy", "TST", utils.P("2020-06-10T11:54:59Z")},
		NewerThan:          utils.P("2021-06-10T11:54:59Z"),
	}
	as.EXPECT().
		SearchHostsAsLMS(gomock.Any()).
		DoAndReturn(func(actual dto.SearchHostsAsLMS) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filterlsm, actual)

			return nil, aerrMock
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z&newer-than=2021-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.ms-excel.sheet.macroEnabled.12")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchHosts_XLSXSuccess(t *testing.T) {
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

	filters := dto.SearchHostsFilters{
		Search:         []string{"foobar"},
		SortBy:         "Processors",
		SortDesc:       true,
		Location:       "Italy",
		Environment:    "TST",
		OlderThan:      utils.P("2020-06-10T11:54:59Z"),
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}

	xlsx := excelize.NewFile()

	as.EXPECT().
		SearchHostsAsXLSX(filters).
		Return(xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	_, err = excelize.OpenReader(rr.Body)
	assert.NoError(t, err)
}

func TestSearchHosts_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?sort-desc=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_XLSXUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?older-than=asasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_XLSXInternalServerError1(t *testing.T) {
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

	filters := dto.SearchHostsFilters{
		Search:         []string{""},
		SortBy:         "",
		SortDesc:       false,
		Location:       "",
		Environment:    "",
		OlderThan:      utils.MAX_TIME,
		PageNumber:     -1,
		PageSize:       -1,
		Cluster:        new(string),
		LTEMemoryTotal: -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
		GTESwapTotal:   -1,
		LTECPUCores:    -1,
		GTECPUCores:    -1,
		LTECPUThreads:  -1,
		GTECPUThreads:  -1,
	}
	as.EXPECT().
		SearchHostsAsXLSX(filters).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHost_JSONSuccess(t *testing.T) {
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
		"Archived":    false,
		"Cluster":     "Puzzait",
		"CreatedAt":   utils.P("2020-04-15T08:46:58.466+02:00"),
		"Databases":   "",
		"Environment": "PROD",
		"Extra": map[string]interface{}{
			"Clusters": []interface{}{
				map[string]interface{}{
					"CPU":     140,
					"Name":    "Puzzait",
					"Sockets": 10,
					"Type":    "vmware",
					"VMs": []interface{}{
						map[string]interface{}{
							"CappedCPU":          false,
							"ClusterName":        "Puzzait",
							"Hostname":           "test-virt",
							"Name":               "test-virt",
							"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
						},
						map[string]interface{}{
							"CappedCPU":          false,
							"ClusterName":        "Puzzait",
							"Hostname":           "test-db",
							"Name":               "test-db",
							"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
						},
					},
				},
			},
			"Databases": []interface{}{},
			"Filesystems": []interface{}{
				map[string]interface{}{
					"Available":  "4.6G",
					"Filesystem": "/dev/mapper/vg_os-lv_root",
					"FsType":     "xfs",
					"MountedOn":  "/",
					"Size":       "8.0G",
					"Used":       "3.5G",
					"UsedPerc":   "43%",
				},
			},
		},
		"HostDataSchemaVersion": 3,
		"Hostname":              "test-virt",
		"Info": map[string]interface{}{
			"AixCluster":                    false,
			"CPUCores":                      1,
			"CPUModel":                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":                    2,
			"Environment":                   "PROD",
			"Hostname":                      "test-virt",
			"Kernel":                        "3.10.0-862.9.1.el7.x86_64",
			"Location":                      "Italy",
			"MemoryTotal":                   3,
			"OS":                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":                 false,
			"Socket":                        2,
			"SunCluster":                    false,
			"SwapTotal":                     4,
			"HardwareAbstractionTechnology": "VMWARE",
			"VeritasCluster":                false,
			"HardwareAbstraction":           "VIRT",
		},
		"Location":           "Italy",
		"VirtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
		"SchemaVersion":      1,
		"Schemas":            "",
		"ServerVersion":      "latest",
		"Version":            "1.6.1",
		"_id":                utils.Str2oid("5e96ade270c184faca93fe34"),
	}

	as.EXPECT().
		GetHost("foobar", utils.P("2020-06-10T11:54:59Z"), false).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=2020-06-10T11%3A54%3A59Z", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetHost_JSONFailUnprocessableEntity(t *testing.T) {
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
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=fgfggf", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHost_JSONFailInternalServerError(t *testing.T) {
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
		GetHost("foobar", utils.MAX_TIME, false).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHost_JSONFailNotFound(t *testing.T) {
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
		GetHost("foobar", utils.MAX_TIME, false).
		Return(nil, utils.ErrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetHost_MongoJSONSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	res := mongoutils.LoadFixtureMongoHostDataMap(t, "../../fixture/test_dataservice_mongohostdata_02.json")
	expectedRes, err := ioutil.ReadFile("../../fixture/test_dataservice_mongohostdata_02.json")
	require.NoError(t, err)

	as.EXPECT().
		GetHost("foobar", utils.P("2020-06-10T11:54:59Z"), true).
		Return(res, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=2020-06-10T11%3A54%3A59Z", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, string(expectedRes), rr.Body.String())
}

func TestGetHost_MongoJSONFailUnprocessableEntity(t *testing.T) {
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
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=fgfggf", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHost_MongoJSONFailInternalServerError(t *testing.T) {
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
		GetHost("foobar", utils.MAX_TIME, true).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHost_MongoJSONFailNotFound(t *testing.T) {
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
		GetHost("foobar", utils.MAX_TIME, true).
		Return(nil, utils.ErrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestListLocations_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	expectedRes := []string{"Italy", "German", "France"}

	as.EXPECT().
		ListLocations("Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLocations)
	req, err := http.NewRequest("GET", "/locations?environment=TST&location=Italy&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestListLocations_FailUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ListLocations)
	req, err := http.NewRequest("GET", "/locations?older-than=dfsgdfsg", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListLocations_FailInternalServerError(t *testing.T) {
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
		ListLocations("", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLocations)
	req, err := http.NewRequest("GET", "/locations", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestListEnvironments_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	expectedRes := []string{"TST", "PRD", "DEV"}

	as.EXPECT().
		ListEnvironments("Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListEnvironments)
	req, err := http.NewRequest("GET", "/environments?environment=TST&location=Italy&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestListEnvironments_FailUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ListEnvironments)
	req, err := http.NewRequest("GET", "/environments?older-than=dfsgdfsg", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListEnvironments_FailInternalServerError(t *testing.T) {
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
		ListEnvironments("", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListEnvironments)
	req, err := http.NewRequest("GET", "/environments", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestArchiveHost_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().ArchiveHost("foobar").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestArchiveHost_FailReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestArchiveHost_FailNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().ArchiveHost("foobar").Return(utils.ErrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestArchiveHost_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().ArchiveHost("foobar").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
