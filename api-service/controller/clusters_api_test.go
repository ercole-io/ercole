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

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/plandem/xlsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchCluster_JSONPaged(t *testing.T) {
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
				"CPU":                         0,
				"Environment":                 "PROD",
				"Hostname":                    "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
				"HostnameAgentVirtualization": "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
				"Location":                    "Italy",
				"Name":                        "not_in_cluster",
				"PhysicalHosts":               "aspera-b1fe49e8501c9ef031e5acff4b5e69a9",
				"Sockets":                     0,
				"Type":                        "unknown",
				"_id":                         utils.Str2oid("5e8c234b24f648a08585bd3d"),
			},
			map[string]interface{}{
				"CPU":                         140,
				"Environment":                 "PROD",
				"Hostname":                    "test-virt",
				"HostnameAgentVirtualization": "test-virt",
				"Location":                    "Italy",
				"Name":                        "Puzzait",
				"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3f82402",
				"Sockets":                     10,
				"Type":                        "vmware",
				"_id":                         utils.Str2oid("5e8c234b24f648a08585bd41"),
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
		SearchClusters(true, "foobar", "CPU", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?full=true&search=foobar&sort-by=CPU&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchCluster_JSONUnpaged(t *testing.T) {
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
			"CPU":                         0,
			"Environment":                 "PROD",
			"Hostname":                    "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"HostnameAgentVirtualization": "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Location":                    "Italy",
			"Name":                        "not_in_cluster",
			"PhysicalHosts":               "aspera-b1fe49e8501c9ef031e5acff4b5e69a9",
			"Sockets":                     0,
			"Type":                        "unknown",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPU":                         140,
			"Environment":                 "PROD",
			"Hostname":                    "test-virt",
			"HostnameAgentVirtualization": "test-virt",
			"Location":                    "Italy",
			"Name":                        "Puzzait",
			"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3f82402",
			"Sockets":                     10,
			"Type":                        "vmware",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	as.EXPECT().
		SearchClusters(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchCluster_JSONUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?full=ddfssdf", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchCluster_JSONUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?sort-desc=ddfssdf", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchCluster_JSONUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?page=ddfssdf", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchCluster_JSONUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?size=ddfssdf", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchCluster_JSONUnprocessableEntity5(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?older-than=ddfssdf", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchCluster_JSONInternalServerError(t *testing.T) {
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
		SearchClusters(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchCluster_XLSXSuccess(t *testing.T) {
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
			"CPU":                         0.0,
			"Environment":                 "PROD",
			"Hostname":                    "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"HostnameAgentVirtualization": "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Location":                    "Italy",
			"Name":                        "not_in_cluster",
			"PhysicalHosts":               "aspera-b1fe49e8501c9ef031e5acff4b5e69a9",
			"Sockets":                     0.0,
			"Type":                        "unknown",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPU":                         140.0,
			"Environment":                 "PROD",
			"Hostname":                    "test-virt",
			"HostnameAgentVirtualization": "test-virt",
			"Location":                    "Italy",
			"Name":                        "Puzzait",
			"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3f82402",
			"Sockets":                     10.0,
			"Type":                        "vmware",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	as.EXPECT().
		SearchClusters(false, "foobar", "CPU", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(res, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?search=foobar&sort-by=CPU&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Hypervisor")
	require.NotNil(t, sh)
	assert.Equal(t, "not_in_cluster", sh.Cell(0, 1).String())
	assert.Equal(t, "unknown", sh.Cell(1, 1).String())
	AssertXLSXInt(t, 0, sh.Cell(2, 1))
	AssertXLSXInt(t, 0, sh.Cell(3, 1))
	assert.Equal(t, "aspera-b1fe49e8501c9ef031e5acff4b5e69a9", sh.Cell(4, 1).String())
	assert.Equal(t, "Puzzait", sh.Cell(0, 2).String())
	assert.Equal(t, "vmware", sh.Cell(1, 2).String())
	AssertXLSXInt(t, 140, sh.Cell(2, 2))
	AssertXLSXInt(t, 10, sh.Cell(3, 2))
	assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", sh.Cell(4, 2).String())
}

func TestSearchCluster_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?sort-desc=sdsdsdf", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchCluster_XLSXUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?older-than=sdsdsdf", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchCluster_XLSXInternalServerError1(t *testing.T) {
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
		SearchClusters(false, "foobar", "CPU", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?search=foobar&sort-by=CPU&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchCluster_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	res := []map[string]interface{}{
		{
			"CPU":                         0.0,
			"Environment":                 "PROD",
			"Hostname":                    "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"HostnameAgentVirtualization": "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Location":                    "Italy",
			"Name":                        "not_in_cluster",
			"PhysicalHosts":               "aspera-b1fe49e8501c9ef031e5acff4b5e69a9",
			"Sockets":                     0.0,
			"Type":                        "unknown",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPU":                         140.0,
			"Environment":                 "PROD",
			"Hostname":                    "test-virt",
			"HostnameAgentVirtualization": "test-virt",
			"Location":                    "Italy",
			"Name":                        "Puzzait",
			"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3f82402",
			"Sockets":                     10.0,
			"Type":                        "vmware",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	as.EXPECT().
		SearchClusters(false, "foobar", "CPU", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(res, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?search=foobar&sort-by=CPU&sort-desc=true&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
