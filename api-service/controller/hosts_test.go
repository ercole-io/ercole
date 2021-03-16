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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"content": []interface{}{
			map[string]interface{}{
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
			map[string]interface{}{
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
		SearchHosts("summary", gomock.Any()).
		DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filters, actual)

			return resFromService, nil
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?mode=summary&search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
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
	as.EXPECT().
		SearchHostsAsLMS(gomock.Any()).
		DoAndReturn(func(actual dto.SearchHostsFilters) (*excelize.File, error) {
			assert.EqualValues(t, filters, actual)

			return expected, nil
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&&older-than=2020-06-10T11%3A54%3A59Z", nil)
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
		Log: utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
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
	as.EXPECT().
		SearchHostsAsLMS(gomock.Any()).
		DoAndReturn(func(actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filters, actual)

			return nil, aerrMock
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&&older-than=2020-06-10T11%3A54%3A59Z", nil)
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
		Log: utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"agentVersion":                  "latest",
			"cpuCores":                      24,
			"cpuModel":                      "Intel(R) Xeon(R) Platinum 8160 CPU @ 2.10GHz",
			"cpuSockets":                    1,
			"cpuThreads":                    48,
			"cluster":                       nil,
			"createdAt":                     utils.PDT("2020-07-01T09:18:03.715+02:00"),
			"environment":                   "PROD",
			"hacmp":                         false,
			"hardwareAbstraction":           "PH",
			"hardwareAbstractionTechnology": "PH",
			"hostname":                      "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
			"kernel":                        "Linux 4.1.12-124.26.12.el7uek.x86_64",
			"location":                      "Italy",
			"memTotal":                      376,
			"os":                            "Red Hat Enterprise Linux 7.6",
			"oracleClusterware":             true,
			"sunCluster":                    false,
			"swapTotal":                     23,
			"veritasClusterServer":          false,
			"virtualizationNode":            nil,
			"_id":                           utils.Str2oid("5efc38ab79f92e4cbf283b0b"),
		},
		{
			"agentVersion":                  "latest",
			"cpuCores":                      1,
			"cpuModel":                      "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
			"cpuSockets":                    2,
			"cpuThreads":                    2,
			"cluster":                       "Puzzait",
			"createdAt":                     utils.PDT("2020-07-01T09:18:03.726+02:00"),
			"environment":                   "TST",
			"hacmp":                         false,
			"hardwareAbstraction":           "VIRT",
			"hardwareAbstractionTechnology": "VMWARE",
			"hostname":                      "test-db",
			"kernel":                        "Linux 3.10.0-514.el7.x86_64",
			"location":                      "Germany",
			"memTotal":                      3,
			"os":                            "Red Hat Enterprise Linux 7.6",
			"oracleClusterware":             false,
			"sunCluster":                    false,
			"swapTotal":                     1,
			"veritasClusterServer":          false,
			"virtualizationNode":            "s157-cb32c10a56c256746c337e21b3f82402",
			"_id":                           utils.Str2oid("5efc38ab79f92e4cbf283b13"),
		},
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
	as.EXPECT().
		SearchHosts("summary", gomock.Any()).
		DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filters, actual)

			return expectedRes, nil
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := excelize.OpenReader(rr.Body)
	require.NoError(t, err)
	assert.Equal(t, "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", sp.GetCellValue("Hosts", "A2"))
	assert.Equal(t, "PROD", sp.GetCellValue("Hosts", "B2"))
	assert.Equal(t, "", sp.GetCellValue("Hosts", "C2"))
	assert.Equal(t, "", sp.GetCellValue("Hosts", "D2"))
	assert.Equal(t, "", sp.GetCellValue("Hosts", "E2"))
	assert.Equal(t, "latest", sp.GetCellValue("Hosts", "F2"))
	assert.Equal(t, utils.P("2020-07-01T09:18:03.715+02:00").UTC().String(), sp.GetCellValue("Hosts", "G2"))
	assert.Equal(t, "", sp.GetCellValue("Hosts", "H2"))
	assert.Equal(t, "Red Hat Enterprise Linux 7.6", sp.GetCellValue("Hosts", "I2"))
	assert.Equal(t, "Linux 4.1.12-124.26.12.el7uek.x86_64", sp.GetCellValue("Hosts", "J2"))
	assert.Equal(t, "1", sp.GetCellValue("Hosts", "K2"))
	assert.Equal(t, "0", sp.GetCellValue("Hosts", "L2"))
	assert.Equal(t, "0", sp.GetCellValue("Hosts", "M2"))
	assert.Equal(t, "PH", sp.GetCellValue("Hosts", "N2"))
	assert.Equal(t, "PH", sp.GetCellValue("Hosts", "O2"))
	assert.Equal(t, "48", sp.GetCellValue("Hosts", "P2"))
	assert.Equal(t, "24", sp.GetCellValue("Hosts", "Q2"))
	assert.Equal(t, "1", sp.GetCellValue("Hosts", "R2"))
	assert.Equal(t, "376", sp.GetCellValue("Hosts", "S2"))
	assert.Equal(t, "23", sp.GetCellValue("Hosts", "T2"))
	assert.Equal(t, "Intel(R) Xeon(R) Platinum 8160 CPU @ 2.10GHz", sp.GetCellValue("Hosts", "U2"))

	assert.Equal(t, "test-db", sp.GetCellValue("Hosts", "A3"))
	assert.Equal(t, "TST", sp.GetCellValue("Hosts", "B3"))
	assert.Equal(t, "", sp.GetCellValue("Hosts", "C3"))
	assert.Equal(t, "Puzzait", sp.GetCellValue("Hosts", "D3"))
	assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", sp.GetCellValue("Hosts", "E3"))
	assert.Equal(t, "latest", sp.GetCellValue("Hosts", "F3"))
	assert.Equal(t, utils.P("2020-07-01T09:18:03.726+02:00").UTC().String(), sp.GetCellValue("Hosts", "G3"))
	assert.Equal(t, "", sp.GetCellValue("Hosts", "H3"))
	assert.Equal(t, "Red Hat Enterprise Linux 7.6", sp.GetCellValue("Hosts", "I3"))
	assert.Equal(t, "Linux 3.10.0-514.el7.x86_64", sp.GetCellValue("Hosts", "J3"))
	assert.Equal(t, "0", sp.GetCellValue("Hosts", "K3"))
	assert.Equal(t, "0", sp.GetCellValue("Hosts", "L3"))
	assert.Equal(t, "0", sp.GetCellValue("Hosts", "M3"))
	assert.Equal(t, "VIRT", sp.GetCellValue("Hosts", "N3"))
	assert.Equal(t, "VMWARE", sp.GetCellValue("Hosts", "O3"))
	assert.Equal(t, "2", sp.GetCellValue("Hosts", "P3"))
	assert.Equal(t, "1", sp.GetCellValue("Hosts", "Q3"))
	assert.Equal(t, "2", sp.GetCellValue("Hosts", "R3"))
	assert.Equal(t, "3", sp.GetCellValue("Hosts", "S3"))
	assert.Equal(t, "1", sp.GetCellValue("Hosts", "T3"))
	assert.Equal(t, "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz", sp.GetCellValue("Hosts", "U3"))
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
		Log: utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
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
		SearchHosts("summary", filters).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchHosts_XLSXInternalServerError2(t *testing.T) {
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
		SearchHosts("summary", gomock.Any()).
		DoAndReturn(func(_ string, actual dto.SearchHostsFilters) ([]map[string]interface{}, error) {
			assert.EqualValues(t, filters, actual)

			return expectedRes, nil
		})

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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
	}

	res := utils.LoadFixtureMongoHostDataMap(t, "../../fixture/test_dataservice_mongohostdata_02.json")
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log: utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
		Log:     utils.NewLogger("TEST"),
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
