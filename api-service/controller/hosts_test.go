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
	"os"
	"testing"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	gomock "go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

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
		From:               utils.MIN_TIME,
		To:                 utils.MAX_TIME,
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
		From:               utils.MIN_TIME,
		To:                 utils.MAX_TIME,
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

	expectedRes := dto.HostData{
		Archived:    false,
		Cluster:     "Puzzait",
		CreatedAt:   utils.P("2020-04-15T08:46:58.466Z"),
		Environment: "PROD",
		Clusters: []model.ClusterInfo{
			{
				CPU:     140,
				Name:    "Puzzait",
				Sockets: 10,
				Type:    "vmware",
				VMs: []model.VMInfo{
					{
						CappedCPU:               false,
						Hostname:                "test-virt",
						Name:                    "test-virt",
						VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
						PhysicalServerModelName: "test physical server model name",
					},
					{
						CappedCPU:               false,
						Hostname:                "test-db",
						Name:                    "test-db",
						VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
						PhysicalServerModelName: "test physical server model name",
					},
				},
			},
		},
		Filesystems: []model.Filesystem{
			{
				AvailableSpace: 4.60000000e+09,
				Filesystem:     "/dev/mapper/vg_os-lv_root",
				MountedOn:      "/",
				Size:           8.00000000e+09,
				Type:           "xfs",
				UsedSpace:      3.50000000e+09,
			},
		},
		SchemaVersion: 3,
		Hostname:      "test-virt",
		Info: model.Host{
			CPUCores:                      1,
			CPUFrequency:                  "2.50GHz",
			CPUModel:                      "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			CPUSockets:                    2,
			CPUThreads:                    2,
			CoresPerSocket:                1,
			HardwareAbstraction:           "VIRT",
			HardwareAbstractionTechnology: "VMWARE",
			Hostname:                      "test-virt",
			Kernel:                        "Linux",
			KernelVersion:                 "3.10.0-862.9.1.el7.x86_64",
			MemoryTotal:                   3,
			OS:                            "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			OSVersion:                     "7.5",
			SwapTotal:                     4,
			ThreadsPerCore:                2,
		},
		Location:            "Italy",
		VirtualizationNode:  "s157-cb32c10a56c256746c337e21b3f82402",
		ServerSchemaVersion: 1,
		ServerVersion:       "latest",
		AgentVersion:        "1.6.1",
		ID:                  utils.Str2oid("5e8c234b24f648a08585bd41"),
	}

	var user interface{}
	locations := []string{"Italy"}

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().
		GetHost("test-virt", utils.P("2020-06-10T11:54:59Z"), false).
		Return(&expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/test-virt?older-than=2020-06-10T11%3A54%3A59Z", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-virt",
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

	var res dto.HostData
	raw, err := os.ReadFile("../../fixture/test_dataservice_mongohostdata_05.json")
	require.NoError(t, err)
	err = bson.UnmarshalExtJSON(raw, true, &res)
	require.NoError(t, err)
	expectedRes, err := os.ReadFile("../../fixture/test_dataservice_mongohostdata_05.json")
	require.NoError(t, err)

	var user interface{}
	as.EXPECT().
		ListLocations(user).
		Return([]string{"Italy", "Germany", "France"}, nil)

	as.EXPECT().
		GetHost("foobar", utils.P("2020-06-10T11:54:59Z"), true).
		Return(&res, nil)

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

func TestDismissHost_Success(t *testing.T) {
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
		Return(&dto.HostData{Location: "Italy"}, nil)

	var user interface{}
	as.EXPECT().
		ListLocations(user).
		Return([]string{"Italy", "Germany", "France"}, nil)

	as.EXPECT().DismissHost("foobar").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DismissHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDismissHost_FailReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DismissHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDismissHost_FailNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetHost("foobar", utils.MAX_TIME, true).Return(nil, utils.ErrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DismissHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDismissHost_FailInternalServerError(t *testing.T) {
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
		Return(&dto.HostData{Location: "Italy"}, nil)

	var user interface{}
	as.EXPECT().
		ListLocations(user).
		Return([]string{"Italy", "Germany", "France"}, nil)

	as.EXPECT().DismissHost("foobar").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DismissHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
