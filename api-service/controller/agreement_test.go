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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOracleDatabaseAgreementPartsList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []model.OracleDatabaseAgreementPart{
		{
			ItemDescription: "foobar",
		},
	}

	as.EXPECT().
		GetOracleDatabaseAgreementPartsList().
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreementPartsList)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetOracleDatabaseAgreementPartsList_InternalServerError(t *testing.T) {
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
		GetOracleDatabaseAgreementPartsList().
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleDatabaseAgreementPartsList)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAddOracleDatabaseAgreements_Success(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	payload := apimodel.OracleDatabaseAgreementsAddRequest{
		AgreementID: "foobar",
		Count:       20,
	}

	as.EXPECT().AddOracleDatabaseAgreements(payload).Return([]string{"foo", "bar"}, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreements)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON([]string{"foo", "bar"}), rr.Body.String())
}

func TestAddOracleDatabaseAgreements_ReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreements)
	req, err := http.NewRequest("POST", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddOracleDatabaseAgreements_BadRequest(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreements)
	req, err := http.NewRequest("POST", "/", &FailingReader{})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddOracleDatabaseAgreements_UnprocessableEntity1(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreements)
	req, err := http.NewRequest("POST", "/", strings.NewReader("{\"sadasdsad"))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestAddOracleDatabaseAgreements_InternalServerError(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	payload := apimodel.OracleDatabaseAgreementsAddRequest{
		AgreementID: "foobar",
		Count:       20,
	}

	as.EXPECT().AddOracleDatabaseAgreements(payload).Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddOracleDatabaseAgreements)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(payload))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchOracleDatabaseAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []apimodel.OracleDatabaseAgreementsFE{
		{
			ItemDescription: "foobar",
		},
	}

	as.EXPECT().
		SearchOracleDatabaseAgreements("", apimodel.SearchOracleDatabaseAgreementsFilters{
			Unlimited:         "true",
			CatchAll:          "NULL",
			AvailableCountGTE: -1,
			AvailableCountLTE: -1,
			LicensesCountGTE:  -1,
			LicensesCountLTE:  -1,
			UsersCountGTE:     -1,
			UsersCountLTE:     -1,
		}).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=true", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleDatabaseAgreements_FailedUnprocessableEntity(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=sasasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleDatabaseAgreements_FailedInternalServerError(t *testing.T) {
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
		SearchOracleDatabaseAgreements("", apimodel.SearchOracleDatabaseAgreementsFilters{
			Unlimited:         "NULL",
			CatchAll:          "NULL",
			AvailableCountGTE: -1,
			AvailableCountLTE: -1,
			LicensesCountGTE:  -1,
			LicensesCountLTE:  -1,
			UsersCountGTE:     -1,
			UsersCountLTE:     -1,
		}).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetSearchOracleDatabaseAgreementsFilters_SuccessEmpty(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	filters, err := GetSearchOracleDatabaseAgreementsFilters(r)
	require.NoError(t, err)
	assert.Equal(t, apimodel.SearchOracleDatabaseAgreementsFilters{
		CatchAll:          "NULL",
		Unlimited:         "NULL",
		AvailableCountGTE: -1,
		AvailableCountLTE: -1,
		LicensesCountGTE:  -1,
		LicensesCountLTE:  -1,
		UsersCountGTE:     -1,
		UsersCountLTE:     -1,
	}, filters)
}

func TestGetSearchOracleDatabaseAgreementsFilters_SuccessFull(t *testing.T) {
	r, err := http.NewRequest("GET", "/?agreement-id=foo&part-id=bar&item-description=boz&csi=pippo&metrics=pluto&reference-number=foobar&unlimited=false&catch-all=true&licenses-count-lte=10&licenses-count-gte=20&users-count-lte=5&users-count-gte=15&available-count-lte=3&available-count-gte=13", nil)
	require.NoError(t, err)
	filters, err := GetSearchOracleDatabaseAgreementsFilters(r)
	require.NoError(t, err)
	assert.Equal(t, apimodel.SearchOracleDatabaseAgreementsFilters{
		AgreementID:       "foo",
		PartID:            "bar",
		ItemDescription:   "boz",
		CSI:               "pippo",
		Metrics:           "pluto",
		ReferenceNumber:   "foobar",
		Unlimited:         "false",
		CatchAll:          "true",
		LicensesCountLTE:  10,
		LicensesCountGTE:  20,
		UsersCountLTE:     5,
		UsersCountGTE:     15,
		AvailableCountLTE: 3,
		AvailableCountGTE: 13,
	}, filters)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail1(t *testing.T) {
	r, err := http.NewRequest("GET", "/?unlimited=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail2(t *testing.T) {
	r, err := http.NewRequest("GET", "/?catch-all=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail3(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-count-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail4(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-count-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail5(t *testing.T) {
	r, err := http.NewRequest("GET", "/?users-count-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail6(t *testing.T) {
	r, err := http.NewRequest("GET", "/?users-count-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail7(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-count-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestGetSearchOracleDatabaseAgreementsFilters_Fail8(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-count-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = GetSearchOracleDatabaseAgreementsFilters(r)
	require.Error(t, err)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_Success(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedInvalidID(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "saddsfadasf",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedBrokenBody(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", &FailingReader{})
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedAgreementNotFound(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedNotInClusterHostNotFound(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrNotInClusterHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedAerrOracleDatabaseAgreementNotFound(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddAssociatedHostToOracleDatabaseAgreement_FailedInternalServerError(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("POST", "/", strings.NewReader("foohost"))
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_Success(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().RemoveAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.RemoveAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_FailedInvalidID(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "sdsdfaasdf",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_FailedAgreementNotFound(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().RemoveAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRemoveAssociatedHostToOracleDatabaseAgreement_FailedInternalServerError(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().RemoveAssociatedHostToOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveAssociatedHostToOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":       "5f50a98611959b1baa17525e",
		"hostname": "foohost",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteOracleDatabaseAgreement_Success(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
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

func TestDeleteOracleDatabaseAgreement_FailedReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DeleteOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDeleteOracleDatabaseAgreement_FailedInvalidID(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
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

func TestDeleteOracleDatabaseAgreement_FailedAgreementNotFound(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(utils.AerrOracleDatabaseAgreementNotFound)

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

func TestDeleteOracleDatabaseAgreement_FailedInternalServerError(t *testing.T) {
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
		Log: utils.NewLogger("TEST"),
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
