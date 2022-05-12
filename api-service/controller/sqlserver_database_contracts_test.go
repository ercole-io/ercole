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
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOSqlServerDatabaseContract_Success(t *testing.T) {
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
		Log: logger.NewLogger("TEST"),
	}

	request := model.SqlServerDatabaseContract{
		Type:           "359-06320",
		ContractID:     "contractID test",
		LicensesNumber: 9999,
		Hosts: []string{
			"ERCWIN2016",
			"sflnxdb104",
		},
		Clusters: []string{},
	}

	result := request

	as.EXPECT().AddSqlServerDatabaseContract(request).Return(&result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddSqlServerDatabaseContract)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(result), rr.Body.String())
}

func TestAddSqlServerDatabaseContract_ReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddSqlServerDatabaseContract)
	req, err := http.NewRequest("POST", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddSqlServerDatabaseContract_BadRequests(t *testing.T) {
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
		Log: logger.NewLogger("TEST"),
	}

	t.Run("fail to decode", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.AddSqlServerDatabaseContract)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("id must be empty", func(t *testing.T) {
		request := model.SqlServerDatabaseContract{
			ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			Hosts: []string{"foobar"},
		}

		handler := http.HandlerFunc(ac.AddSqlServerDatabaseContract)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestAddSqlServerDatabaseContract_InternalServerError(t *testing.T) {
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
		Log: logger.NewLogger("TEST"),
	}

	request := model.SqlServerDatabaseContract{
		ContractID:     "foobar",
		LicensesNumber: 20,
	}

	as.EXPECT().AddSqlServerDatabaseContract(request).Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddSqlServerDatabaseContract)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateSqlServerDatabaseContract_Success(t *testing.T) {
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
		Log: logger.NewLogger("TEST"),
	}

	request := model.SqlServerDatabaseContract{
		ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Hosts: []string{"foobar"},
	}

	agr := new(model.SqlServerDatabaseContract)
	as.EXPECT().UpdateSqlServerDatabaseContract(request).Return(agr, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.UpdateSqlServerDatabaseContract)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(agr), rr.Body.String())
}

func TestUpdateSqlServerDatabaseContract_BadRequests(t *testing.T) {
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
		Log: logger.NewLogger("TEST"),
	}

	t.Run("fail to decode", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateSqlServerDatabaseContract)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestGetSqlServerDatabaseContracts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contracts := []model.SqlServerDatabaseContract{
		{
			Type: "test",
		},
	}

	as.EXPECT().
		GetSqlServerDatabaseContracts().
		Return(contracts, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetSqlServerDatabaseContracts)
	req, err := http.NewRequest("GET", "/?unlimited=true", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := map[string]interface{}{
		"contracts": contracts,
	}
	assert.JSONEq(t, utils.ToJSON(expectedResponse), rr.Body.String())
}
