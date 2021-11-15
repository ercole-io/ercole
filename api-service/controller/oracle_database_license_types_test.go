// Copyright (c) 2021 Sorint.lab S.p.A.
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

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetOracleDatabaseLicenseTypes_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	ltRes := []model.OracleDatabaseLicenseType{
		{
			ItemDescription: "foobar",
		},
	}

	as.EXPECT().
		GetOracleDatabaseLicenseTypes().
		Return(ltRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseLicenseTypes)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	expectedRes := map[string]interface{}{
		"license-types": ltRes,
	}
	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetOracleDatabaseLicenseTypes_InternalServerError(t *testing.T) {
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
		GetOracleDatabaseLicenseTypes().
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseLicenseTypes)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetOracleDatabaseLicensesCompliance_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	compliance := 75.0 / 275.0
	expectedRes := []dto.LicenseCompliance{
		{LicenseTypeID: "PID001", ItemDescription: "itemDesc1", Metric: "Processor Perpetual", Consumed: 7, Covered: 7, Compliance: 1},
		{LicenseTypeID: "PID002", ItemDescription: "itemDesc2", Metric: "Named User Plus Perpetual", Consumed: 275, Covered: 75, Compliance: compliance},
	}
	as.EXPECT().
		GetOracleDatabaseLicensesCompliance().
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseLicensesCompliance)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetOracleDatabaseLicensesCompliance_InternalServerError(t *testing.T) {
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
		GetOracleDatabaseLicensesCompliance().
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseLicensesCompliance)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAddOracleDatabaseLicenseType_Success(t *testing.T) {
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

	request := model.OracleDatabaseLicenseType{
		ID:              "Test",
		ItemDescription: "Oracle Database Enterprise Edition",
		Metric:          "Processor Perpetual",
		Cost:            500,
		Aliases:         []string{"Tuning Pack"},
		Option:          false,
	}

	result := model.OracleDatabaseLicenseType{
		ID:              request.ID,
		ItemDescription: request.ItemDescription,
		Metric:          request.Metric,
		Cost:            request.Cost,
		Aliases:         request.Aliases,
		Option:          request.Option,
	}

	as.EXPECT().AddOracleDatabaseLicenseType(request).Return(&result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseLicenseType)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(result), rr.Body.String())
}

func TestAddOracleDatabaseLicenseType_ReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddOracleDatabaseLicenseType)
	req, err := http.NewRequest("POST", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddOracleDatabaseLicenseType_BadRequests(t *testing.T) {
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
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseLicenseType)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestAddOracleDatabaseLicenseType_InternalServerError(t *testing.T) {
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

	request := model.OracleDatabaseLicenseType{
		ItemDescription: "Oracle Database Enterprise Edition",
		Cost:            500,
	}

	as.EXPECT().AddOracleDatabaseLicenseType(request).Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseLicenseType)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateOracleDatabaseLicenseType_Success(t *testing.T) {
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

	request := model.OracleDatabaseLicenseType{
		ID:              "Test",
		ItemDescription: "Oracle Database Enterprise Edition",
		Metric:          "Processor Perpetual",
		Cost:            500,
		Aliases:         []string{"Tuning Pack"},
		Option:          false,
	}

	agr := new(model.OracleDatabaseLicenseType)
	as.EXPECT().UpdateOracleDatabaseLicenseType(request).Return(agr, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.UpdateOracleDatabaseLicenseType)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(agr), rr.Body.String())
}

func TestUpdateOracleDatabaseLicenseType_BadRequests(t *testing.T) {
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
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseLicenseType)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestUpdateOracleDatabaseLicenseType_InternalServerError(t *testing.T) {
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

	request := model.OracleDatabaseLicenseType{
		ID:              "Test",
		ItemDescription: "Oracle Database Enterprise Edition",
		Metric:          "Processor Perpetual",
		Cost:            500,
		Aliases:         []string{"Tuning Pack"},
		Option:          false,
	}

	t.Run("Unknown error", func(t *testing.T) {
		as.EXPECT().UpdateOracleDatabaseLicenseType(request).Return(nil, aerrMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseLicenseType)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteOracleDatabaseLicenseType_Success(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseLicenseType("Test").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseLicenseType)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "Test",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteOracleDatabaseLicenseType_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseLicenseType)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "Test",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDeleteOracleDatabaseLicenseType_FailedInvalidID(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseLicenseType("Test").Return(utils.ErrOracleDatabaseLicenseTypeIDNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseLicenseType)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "Test",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteOracleDatabaseLicenseType_FailedInternalServerError(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseLicenseType("Test").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseLicenseType)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "Test",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
