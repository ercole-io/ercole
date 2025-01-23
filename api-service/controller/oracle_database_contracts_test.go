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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestAddOracleDatabaseContract_Success(t *testing.T) {
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

	request := model.OracleDatabaseContract{
		ContractID:      "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		Basket:          false,
		Restricted:      false,
		Hosts:           []string{"foobar"},
	}

	result := dto.OracleDatabaseContractFE{
		ID:              utils.Str2oid(fmt.Sprintf("%024d", 1)),
		ContractID:      request.ContractID,
		CSI:             request.CSI,
		LicenseTypeID:   request.LicenseTypeID,
		ReferenceNumber: request.ReferenceNumber,
		Unlimited:       request.Unlimited,
		Basket:          request.Basket,
		Restricted:      request.Restricted,
	}

	as.EXPECT().AddOracleDatabaseContract(request).Return(&result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(result), rr.Body.String())
}

func TestAddOracleDatabaseContract_ReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddOracleDatabaseContract_BadRequests(t *testing.T) {
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
		handler := http.HandlerFunc(ac.AddOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("id must be empty", func(t *testing.T) {
		request := model.OracleDatabaseContract{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			ContractID:      "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "REF001",
			Unlimited:       false,
			Count:           42,
			Basket:          false,
			Hosts:           []string{"foobar"},
		}

		handler := http.HandlerFunc(ac.AddOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("if restricted can't be basket", func(t *testing.T) {
		request := model.OracleDatabaseContract{
			ID:              utils.Str2oid(""),
			ContractID:      "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "REF001",
			Unlimited:       false,
			Count:           42,
			Basket:          true,
			Restricted:      true,
			Hosts:           []string{"foobar"},
		}

		handler := http.HandlerFunc(ac.AddOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestAddOracleDatabaseContract_InternalServerError(t *testing.T) {
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

	request := model.OracleDatabaseContract{
		ContractID: "foobar",
		Count:      20,
	}

	as.EXPECT().AddOracleDatabaseContract(request).Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateOracleDatabaseContract_Success(t *testing.T) {
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

	request := model.OracleDatabaseContract{
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		ContractID:      "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		Basket:          false,
		Hosts:           []string{"foobar"},
	}

	agr := new(dto.OracleDatabaseContractFE)
	as.EXPECT().UpdateOracleDatabaseContract(request).Return(agr, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.UpdateOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(agr), rr.Body.String())
}

func TestUpdateOracleDatabaseContract_BadRequests(t *testing.T) {
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
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("if restricted can't be basket", func(t *testing.T) {
		request := model.OracleDatabaseContract{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			ContractID:      "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "REF001",
			Unlimited:       false,
			Count:           42,
			Basket:          true,
			Restricted:      true,
			Hosts:           []string{"foobar"},
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestUpdateOracleDatabaseContract_InternalServerError(t *testing.T) {
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

	request := model.OracleDatabaseContract{
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		ContractID:      "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		Basket:          false,
		Hosts:           []string{"foobar"},
	}

	t.Run("Unknown error", func(t *testing.T) {
		as.EXPECT().UpdateOracleDatabaseContract(request).Return(nil, aerrMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Contract not found", func(t *testing.T) {
		as.EXPECT().UpdateOracleDatabaseContract(request).
			Return(nil, utils.ErrContractNotFound)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
	t.Run("Invalid PartID", func(t *testing.T) {
		as.EXPECT().UpdateOracleDatabaseContract(request).
			Return(nil, utils.ErrOracleDatabaseLicenseTypeIDNotFound)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseContract)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
}

func TestGetOracleDatabaseContracts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contracts := []dto.OracleDatabaseContractFE{
		{
			ItemDescription: "foobar",
		},
	}

	as.EXPECT().ListLocations(gomock.Any()).Return([]string{}, nil)
	as.EXPECT().
		GetOracleDatabaseContracts(dto.GetOracleDatabaseContractsFilter{
			ContractID:                  "",
			LicenseTypeID:               "",
			ItemDescription:             "",
			CSI:                         "",
			Metric:                      "",
			ReferenceNumber:             "",
			Unlimited:                   "true",
			Basket:                      "",
			LicensesPerCoreLTE:          -1,
			LicensesPerCoreGTE:          -1,
			LicensesPerUserLTE:          -1,
			LicensesPerUserGTE:          -1,
			AvailableLicensesPerCoreLTE: -1,
			AvailableLicensesPerCoreGTE: -1,
			AvailableLicensesPerUserLTE: -1,
			AvailableLicensesPerUserGTE: -1,
			Locations:                   []string{},
		}).
		Return(contracts, nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ac.GetOracleDatabaseContracts)
	req, err := http.NewRequest("GET", "/?unlimited=true", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := map[string]interface{}{
		"contracts": contracts,
	}
	assert.JSONEq(t, utils.ToJSON(expectedResponse), rr.Body.String())
}

func TestGetOracleDatabaseContracts_FailedUnprocessableEntity(t *testing.T) {
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
	as.EXPECT().ListLocations(gomock.Any()).Return([]string{}, nil)

	handler := http.HandlerFunc(ac.GetOracleDatabaseContracts)
	req, err := http.NewRequest("GET", "/?unlimited=sasasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetOracleDatabaseContracts_FailedInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().ListLocations(gomock.Any()).Return([]string{}, nil)
	as.EXPECT().
		GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter()).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseContracts)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestParseGetOracleDatabaseContractsFilters_SuccessEmpty(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	filters, err := parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.NoError(t, err)
	assert.Equal(t, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
		Locations:                   []string{},
	}, filters)
}

func TestParseSearchOracleDatabaseContractsFilters_SuccessFull(t *testing.T) {
	r, err := http.NewRequest("GET",
		"/?contract-id=foo&license-type-id=bar&item-description=boz&csi=pippo&metrics=pluto&reference-number=foobar&"+
			"unlimited=false&basket=true&"+
			"licenses-per-core-lte=1&licenses-per-core-gte=2&"+
			"licenses-per-user-lte=3&licenses-per-user-gte=4&"+
			"available-licenses-per-core-lte=3&available-licenses-per-core-gte=13&"+
			"available-licenses-per-user-lte=4&available-licenses-per-user-gte=2",
		nil)
	require.NoError(t, err)
	filters, err := parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.NoError(t, err)
	assert.Equal(t, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "foo",
		LicenseTypeID:               "bar",
		ItemDescription:             "boz",
		CSI:                         "pippo",
		Metric:                      "pluto",
		ReferenceNumber:             "foobar",
		Unlimited:                   "false",
		Basket:                      "true",
		LicensesPerCoreLTE:          1,
		LicensesPerCoreGTE:          2,
		LicensesPerUserLTE:          3,
		LicensesPerUserGTE:          4,
		AvailableLicensesPerCoreLTE: 3,
		AvailableLicensesPerCoreGTE: 13,
		AvailableLicensesPerUserLTE: 4,
		AvailableLicensesPerUserGTE: 2,
		Locations:                   []string{},
	}, filters)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail1(t *testing.T) {
	r, err := http.NewRequest("GET", "/?unlimited=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail2(t *testing.T) {
	r, err := http.NewRequest("GET", "/?basket=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail3(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-core-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail4(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-core-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail5(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-user-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail6(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-user-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail7(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-licenses-per-core-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseContractsFilters_Fail8(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-licenses-per-core-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseContractsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestAddHostToAssociatedPart_Success(t *testing.T) {
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

	as.EXPECT().AddHostToOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestAddHostToAssociatedPart_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddHostToAssociatedPart_FailedInvalidID(t *testing.T) {
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "saddsfadasf",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestAddHostToAssociatedPart_FailedBrokenBody(t *testing.T) {
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", &FailingReader{})
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddHostToAssociatedPart_FailedContractNotFound(t *testing.T) {
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

	as.EXPECT().AddHostToOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrContractNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddHostToAssociatedPart_FailedNotInClusterHostNotFound(t *testing.T) {
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

	as.EXPECT().AddHostToOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrNotInClusterHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddHostToAssociatedPart_FailedAerrOracleDatabaseContractNotFound(t *testing.T) {
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

	as.EXPECT().AddHostToOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrContractNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddHostToAssociatedPart_FailedInternalServerError(t *testing.T) {
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

	as.EXPECT().AddHostToOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseContract)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRemoveHostFromAssociatedPart_Success(t *testing.T) {
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

	as.EXPECT().DeleteHostFromOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestRemoveHostFromAssociatedPart_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRemoveHostFromAssociatedPart_FailedInvalidID(t *testing.T) {
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "sdsdfaasdf",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestRemoveHostFromAssociatedPart_FailedContractNotFound(t *testing.T) {
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

	as.EXPECT().DeleteHostFromOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrContractNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRemoveHostFromAssociatedPart_FailedInternalServerError(t *testing.T) {
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

	as.EXPECT().DeleteHostFromOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseContract_Success(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e")).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseContract_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseContract_FailedInvalidID(t *testing.T) {
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

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "sdasdasdf",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseContract_FailedContractNotFound(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e")).Return(utils.ErrContractNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseContract_FailedInternalServerError(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseContract(utils.Str2oid("5f50a98611959b1baa17525e")).Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseContract)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetOracleDatabaseContractsXLSX_Success(t *testing.T) {
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

	filter := dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "true",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
		Locations:                   []string{},
	}

	xlsx := excelize.File{}

	as.EXPECT().ListLocations(gomock.Any()).Return([]string{}, nil)
	as.EXPECT().
		GetOracleDatabaseContractsAsXLSX(filter).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseContracts)
	req, err := http.NewRequest("GET", "/?unlimited=true", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
}

func TestGetOracleDatabaseContractsXLSX_UnprocessableEntity1(t *testing.T) {
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
	as.EXPECT().ListLocations(gomock.Any()).Return([]string{}, nil)

	handler := http.HandlerFunc(ac.GetOracleDatabaseContracts)
	req, err := http.NewRequest("GET", "/?unlimited=sadasddasasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetOracleDatabaseContractsXLSX_InternalServerError1(t *testing.T) {
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

	filter := dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
		Locations:                   []string{},
	}

	as.EXPECT().ListLocations(gomock.Any()).Return([]string{}, nil)
	as.EXPECT().
		GetOracleDatabaseContractsAsXLSX(filter).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseContracts)
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
