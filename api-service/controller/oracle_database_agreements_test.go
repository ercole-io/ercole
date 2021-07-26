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
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOracleDatabaseAgreement_Success(t *testing.T) {
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

	request := model.OracleDatabaseAgreement{
		AgreementID:     "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		CatchAll:        false,
		Restricted:      false,
		Hosts:           []string{"foobar"},
	}

	result := dto.OracleDatabaseAgreementFE{
		ID:              utils.Str2oid(fmt.Sprintf("%024d", 1)),
		AgreementID:     request.AgreementID,
		CSI:             request.CSI,
		LicenseTypeID:   request.LicenseTypeID,
		ReferenceNumber: request.ReferenceNumber,
		Unlimited:       request.Unlimited,
		CatchAll:        request.CatchAll,
		Restricted:      request.Restricted,
	}

	as.EXPECT().AddOracleDatabaseAgreement(request).Return(&result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(result), rr.Body.String())
}

func TestAddOracleDatabaseAgreement_ReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddOracleDatabaseAgreement_BadRequests(t *testing.T) {
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
		handler := http.HandlerFunc(ac.AddOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("id must be empty", func(t *testing.T) {
		request := model.OracleDatabaseAgreement{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			AgreementID:     "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "REF001",
			Unlimited:       false,
			Count:           42,
			CatchAll:        false,
			Hosts:           []string{"foobar"},
		}

		handler := http.HandlerFunc(ac.AddOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("if restricted can't be catchAll", func(t *testing.T) {
		request := model.OracleDatabaseAgreement{
			ID:              utils.Str2oid(""),
			AgreementID:     "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "REF001",
			Unlimited:       false,
			Count:           42,
			CatchAll:        true,
			Restricted:      true,
			Hosts:           []string{"foobar"},
		}

		handler := http.HandlerFunc(ac.AddOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestAddOracleDatabaseAgreement_InternalServerError(t *testing.T) {
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

	request := model.OracleDatabaseAgreement{
		AgreementID: "foobar",
		Count:       20,
	}

	as.EXPECT().AddOracleDatabaseAgreement(request).Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateOracleDatabaseAgreement_Success(t *testing.T) {
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

	request := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		AgreementID:     "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		CatchAll:        false,
		Hosts:           []string{"foobar"},
	}

	agr := new(dto.OracleDatabaseAgreementFE)
	as.EXPECT().UpdateOracleDatabaseAgreement(request).Return(agr, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.UpdateOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(agr), rr.Body.String())
}

func TestUpdateOracleDatabaseAgreement_BadRequests(t *testing.T) {
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
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", &FailingReader{})
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("if restricted can't be catchAll", func(t *testing.T) {
		request := model.OracleDatabaseAgreement{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			AgreementID:     "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "REF001",
			Unlimited:       false,
			Count:           42,
			CatchAll:        true,
			Restricted:      true,
			Hosts:           []string{"foobar"},
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestUpdateOracleDatabaseAgreement_InternalServerError(t *testing.T) {
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

	request := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		AgreementID:     "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		CatchAll:        false,
		Hosts:           []string{"foobar"},
	}

	t.Run("Unknown error", func(t *testing.T) {
		as.EXPECT().UpdateOracleDatabaseAgreement(request).Return(nil, aerrMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Agreement not found", func(t *testing.T) {
		as.EXPECT().UpdateOracleDatabaseAgreement(request).
			Return(nil, utils.ErrOracleDatabaseAgreementNotFound)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
	t.Run("Invalid PartID", func(t *testing.T) {
		as.EXPECT().UpdateOracleDatabaseAgreement(request).
			Return(nil, utils.ErrOracleDatabaseLicenseTypeIDNotFound)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateOracleDatabaseAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
}

func TestGetOracleDatabaseAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	agreements := []dto.OracleDatabaseAgreementFE{
		{
			ItemDescription: "foobar",
		},
	}

	as.EXPECT().
		GetOracleDatabaseAgreements(dto.GetOracleDatabaseAgreementsFilter{
			AgreementID:                 "",
			LicenseTypeID:               "",
			ItemDescription:             "",
			CSI:                         "",
			Metric:                      "",
			ReferenceNumber:             "",
			Unlimited:                   "true",
			CatchAll:                    "",
			LicensesPerCoreLTE:          -1,
			LicensesPerCoreGTE:          -1,
			LicensesPerUserLTE:          -1,
			LicensesPerUserGTE:          -1,
			AvailableLicensesPerCoreLTE: -1,
			AvailableLicensesPerCoreGTE: -1,
			AvailableLicensesPerUserLTE: -1,
			AvailableLicensesPerUserGTE: -1,
		}).
		Return(agreements, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=true", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := map[string]interface{}{
		"agreements": agreements,
	}
	assert.JSONEq(t, utils.ToJSON(expectedResponse), rr.Body.String())
}

func TestGetOracleDatabaseAgreements_FailedUnprocessableEntity(t *testing.T) {
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
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=sasasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetOracleDatabaseAgreements_FailedInternalServerError(t *testing.T) {
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
		GetOracleDatabaseAgreements(dto.NewGetOracleDatabaseAgreementsFilter()).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestParseGetOracleDatabaseAgreementsFilters_SuccessEmpty(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	filters, err := parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.NoError(t, err)
	assert.Equal(t, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}, filters)
}

func TestParseSearchOracleDatabaseAgreementsFilters_SuccessFull(t *testing.T) {
	r, err := http.NewRequest("GET",
		"/?agreement-id=foo&license-type-id=bar&item-description=boz&csi=pippo&metrics=pluto&reference-number=foobar&"+
			"unlimited=false&catch-all=true&"+
			"licenses-per-core-lte=1&licenses-per-core-gte=2&"+
			"licenses-per-user-lte=3&licenses-per-user-gte=4&"+
			"available-licenses-per-core-lte=3&available-licenses-per-core-gte=13&"+
			"available-licenses-per-user-lte=4&available-licenses-per-user-gte=2",
		nil)
	require.NoError(t, err)
	filters, err := parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.NoError(t, err)
	assert.Equal(t, dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "foo",
		LicenseTypeID:               "bar",
		ItemDescription:             "boz",
		CSI:                         "pippo",
		Metric:                      "pluto",
		ReferenceNumber:             "foobar",
		Unlimited:                   "false",
		CatchAll:                    "true",
		LicensesPerCoreLTE:          1,
		LicensesPerCoreGTE:          2,
		LicensesPerUserLTE:          3,
		LicensesPerUserGTE:          4,
		AvailableLicensesPerCoreLTE: 3,
		AvailableLicensesPerCoreGTE: 13,
		AvailableLicensesPerUserLTE: 4,
		AvailableLicensesPerUserGTE: 2,
	}, filters)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail1(t *testing.T) {
	r, err := http.NewRequest("GET", "/?unlimited=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail2(t *testing.T) {
	r, err := http.NewRequest("GET", "/?catch-all=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail3(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-core-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail4(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-core-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail5(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-user-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail6(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-per-user-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail7(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-licenses-per-core-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail8(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-licenses-per-core-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
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

	as.EXPECT().AddHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
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
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
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
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
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
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", &FailingReader{})
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddHostToAssociatedPart_FailedAgreementNotFound(t *testing.T) {
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

	as.EXPECT().AddHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
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

	as.EXPECT().AddHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrNotInClusterHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddHostToAssociatedPart_FailedAerrOracleDatabaseAgreementNotFound(t *testing.T) {
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

	as.EXPECT().AddHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
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

	as.EXPECT().AddHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToOracleDatabaseAgreement)
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

	as.EXPECT().DeleteHostFromOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseAgreement)
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
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseAgreement)
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
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "sdsdfaasdf",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestRemoveHostFromAssociatedPart_FailedAgreementNotFound(t *testing.T) {
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

	as.EXPECT().DeleteHostFromOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.ErrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseAgreement)
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

	as.EXPECT().DeleteHostFromOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteHostFromOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseAgreement_Success(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseAgreement_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseAgreement_FailedInvalidID(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "sdasdasdf",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseAgreement_FailedAgreementNotFound(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(utils.ErrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteAssociatedPartFromOracleDatabaseAgreement_FailedInternalServerError(t *testing.T) {
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

	as.EXPECT().DeleteOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetOracleDatabaseAgreementsXLSX_Success(t *testing.T) {
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

	filter := dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "true",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}

	xlsx := excelize.File{}

	as.EXPECT().
		GetOracleDatabaseAgreementsAsXLSX(filter).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=true", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
}

func TestGetOracleDatabaseAgreementsXLSX_UnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=sadasddasasd", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetOracleDatabaseAgreementsXLSX_InternalServerError1(t *testing.T) {
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

	filter := dto.GetOracleDatabaseAgreementsFilter{
		AgreementID:                 "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		CatchAll:                    "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}

	as.EXPECT().
		GetOracleDatabaseAgreementsAsXLSX(filter).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
