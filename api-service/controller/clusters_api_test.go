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
	"testing"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	dto "github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestSearchCluster_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	expectedRes := []dto.Cluster{
		{
			CPU:                         0,
			Environment:                 "PROD",
			Hostname:                    "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			HostnameAgentVirtualization: "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			Location:                    "Italy",
			Name:                        "not_in_cluster",
			VirtualizationNodes:         []string{"aspera-b1fe49e8501c9ef031e5acff4b5e69a9"},
			PhysicalServerModelNames:    []string{"model name"},
			Sockets:                     0,
			Type:                        "unknown",
			ID:                          utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			CPU:                         140,
			Environment:                 "PROD",
			Hostname:                    "test-virt",
			HostnameAgentVirtualization: "test-virt",
			Location:                    "Italy",
			Name:                        "Puzzait",
			VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
			PhysicalServerModelNames:    []string{"new model name"},
			Sockets:                     10,
			Type:                        "vmware",
			ID:                          utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	as.EXPECT().
		SearchClusters("full", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchClusterNames_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	returnedRes := []dto.Cluster{
		{
			Name: "not_in_cluster",
		},
		{
			Name: "Puzzait",
		},
	}

	expectedRes := []string{
		"not_in_cluster",
		"Puzzait",
	}

	as.EXPECT().
		SearchClusters("clusternames", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(returnedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?mode=clusternames", nil)
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
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?mode=pppccc", nil)
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
		Log:     logger.NewLogger("TEST"),
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
		Log:     logger.NewLogger("TEST"),
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
		Log:     logger.NewLogger("TEST"),
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
		Log:     logger.NewLogger("TEST"),
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
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchClusters("full", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchClustersAsXLSX_Success(t *testing.T) {
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
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	xlsx := excelize.File{}

	as.EXPECT().
		SearchClustersAsXLSX(filter).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClustersXLSX)
	req, err := http.NewRequest("GET", "/clusters?location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
}

func TestSearchClustersAS_XLSXUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchClustersXLSX)
	req, err := http.NewRequest("GET", "/clusters?older-than=dsasdasd", nil)
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
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters?older-than=sdsdsdf", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchClustersAs_XLSXInternalServerError1(t *testing.T) {
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

	as.EXPECT().
		SearchClustersAsXLSX(filter).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchClusters)
	req, err := http.NewRequest("GET", "/clusters", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetCluster(t *testing.T) {
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

	t.Run("json", func(t *testing.T) {
		cluster := &dto.Cluster{
			ID:                          [12]byte{},
			CPU:                         0,
			CreatedAt:                   time.Time{},
			Environment:                 "",
			FetchEndpoint:               "",
			Hostname:                    "",
			HostnameAgentVirtualization: "",
			Location:                    "Italy",
			Name:                        "Pippo",
			Sockets:                     0,
			Type:                        "",
			VirtualizationNodes:         []string{},
			VirtualizationNodesCount:    0,
			VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
			VMs:                         []dto.VM{},
			VMsCount:                    0,
			VMsErcoleAgentCount:         0,
		}

		as.EXPECT().
			GetCluster("Pippo", utils.P("2020-06-10T11:54:59Z")).
			Return(cluster, nil)

		var user interface{}
		locations := []string{"Italy"}

		as.EXPECT().
			ListLocations(user).
			Return(locations, nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetCluster)
		req, err := http.NewRequest("GET", "/hosts/cluster/Pippo?location=Italy&older-than=2020-06-10T11%3A54%3A59Z", nil)
		require.NoError(t, err)

		req = mux.SetURLVars(req, map[string]string{
			"name": "Pippo",
		})
		req.Header.Add("Accept", "application/json")

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, utils.ToJSON(cluster), rr.Body.String())
	})

	t.Run("xlsx", func(t *testing.T) {
		cluster := &dto.Cluster{
			ID:                          [12]byte{},
			CPU:                         0,
			CreatedAt:                   time.Time{},
			Environment:                 "",
			FetchEndpoint:               "",
			Hostname:                    "",
			HostnameAgentVirtualization: "",
			Location:                    "Italy",
			Name:                        "Pippo",
			Sockets:                     0,
			Type:                        "",
			VirtualizationNodes:         []string{},
			VirtualizationNodesCount:    0,
			VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
			VMs:                         []dto.VM{},
			VMsCount:                    0,
			VMsErcoleAgentCount:         0,
		}

		as.EXPECT().
			GetCluster("Pippo", utils.P("2020-06-10T11:54:59Z")).
			Return(cluster, nil)

		var user interface{}
		locations := []string{"Italy"}

		as.EXPECT().
			ListLocations(user).
			Return(locations, nil)

		xlsx := &excelize.File{}

		as.EXPECT().
			GetClusterXLSX("Pippo", utils.P("2020-06-10T11:54:59Z")).
			Return(xlsx, nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetCluster)
		req, err := http.NewRequest("GET", "/hosts/cluster/Pippo?older-than=2020-06-10T11%3A54%3A59Z", nil)
		require.NoError(t, err)

		req = mux.SetURLVars(req, map[string]string{
			"name": "Pippo",
		})
		req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		_, err = excelize.OpenReader(rr.Body)
		require.NoError(t, err)
	})
}
