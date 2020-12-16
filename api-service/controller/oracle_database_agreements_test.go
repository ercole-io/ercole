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
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAssociatedPartToOracleDbAgreement_Success(t *testing.T) {
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

	request := dto.AssociatedPartInOracleDbAgreementRequest{
		ID:              "",
		AgreementID:     "AID001",
		PartID:          "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		CatchAll:        false,
		Hosts:           []string{"foobar"},
	}

	newID := fmt.Sprintf("%024d", 1)
	as.EXPECT().AddAssociatedPartToOracleDbAgreement(request).Return(newID, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedPartToOracleDbAgreement)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(newID), rr.Body.String())
}

func TestAddAssociatedPartToOracleDbAgreement_ReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddAssociatedPartToOracleDbAgreement)
	req, err := http.NewRequest("POST", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddAssociatedPartToOracleDbAgreement_FailToDecode(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddAssociatedPartToOracleDbAgreement)
	req, err := http.NewRequest("POST", "/", &FailingReader{})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddAssociatedPartToOracleDbAgreement_IDMustBeEmpty(t *testing.T) {
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

	request := dto.AssociatedPartInOracleDbAgreementRequest{
		ID:              "aaaaaaaaaaaaaaaaaaaaaaaa",
		AgreementID:     "AID001",
		PartID:          "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		CatchAll:        false,
		Hosts:           []string{"foobar"},
	}

	handler := http.HandlerFunc(ac.AddAssociatedPartToOracleDbAgreement)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddAssociatedPartToOracleDbAgreement_InternalServerError(t *testing.T) {
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

	request := dto.AssociatedPartInOracleDbAgreementRequest{
		AgreementID: "foobar",
		Count:       20,
	}

	as.EXPECT().AddAssociatedPartToOracleDbAgreement(request).Return("", aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddAssociatedPartToOracleDbAgreement)
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
		Log: utils.NewLogger("TEST"),
	}

	request := dto.AssociatedPartInOracleDbAgreementRequest{
		ID:              "aaaaaaaaaaaaaaaaaaaaaaaa",
		AgreementID:     "AID001",
		PartID:          "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		CatchAll:        false,
		Hosts:           []string{"foobar"},
	}

	as.EXPECT().UpdateAssociatedPartOfOracleDbAgreement(request).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.UpdateAssociatedPartOfOracleDbAgreement)
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "null", rr.Body.String())
}

func TestUpdateAssociatedPartOfOracleDbAgreement_FailToDecode(t *testing.T) {
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
	handler := http.HandlerFunc(ac.UpdateAssociatedPartOfOracleDbAgreement)
	req, err := http.NewRequest("POST", "/", &FailingReader{})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateAssociatedPartOfOracleDbAgreement_InternalServerError(t *testing.T) {
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

	request := dto.AssociatedPartInOracleDbAgreementRequest{
		ID:              "aaaaaaaaaaaaaaaaaaaaaaaa",
		AgreementID:     "AID001",
		PartID:          "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "REF001",
		Unlimited:       false,
		Count:           42,
		CatchAll:        false,
		Hosts:           []string{"foobar"},
	}

	t.Run("Unknown error", func(t *testing.T) {
		as.EXPECT().UpdateAssociatedPartOfOracleDbAgreement(request).Return(aerrMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateAssociatedPartOfOracleDbAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Agreement not found", func(t *testing.T) {
		as.EXPECT().UpdateAssociatedPartOfOracleDbAgreement(request).
			Return(utils.AerrOracleDatabaseAgreementNotFound)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateAssociatedPartOfOracleDbAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
	t.Run("Invalid PartID", func(t *testing.T) {
		as.EXPECT().UpdateAssociatedPartOfOracleDbAgreement(request).
			Return(utils.AerrOracleDatabaseAgreementInvalidPartID)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.UpdateAssociatedPartOfOracleDbAgreement)
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(utils.ToJSON(request))))
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
}
func TestSearchAssociatedPartsInOracleDatabaseAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []dto.OracleDatabaseAgreementFE{
		{
			ItemDescription: "foobar",
		},
	}

	as.EXPECT().
		SearchAssociatedPartsInOracleDatabaseAgreements(dto.SearchOracleDatabaseAgreementsFilter{
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
	handler := http.HandlerFunc(ac.SearchAssociatedPartsInOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=true", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchAssociatedPartsInOracleDatabaseAgreements_FailedUnprocessableEntity(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchAssociatedPartsInOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/?unlimited=sasasd", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchAssociatedPartsInOracleDatabaseAgreements_FailedInternalServerError(t *testing.T) {
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
		SearchAssociatedPartsInOracleDatabaseAgreements(dto.SearchOracleDatabaseAgreementsFilter{
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
	handler := http.HandlerFunc(ac.SearchAssociatedPartsInOracleDatabaseAgreements)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestParseSearchOracleDatabaseAgreementsFilters_SuccessEmpty(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	filters, err := parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.NoError(t, err)
	assert.Equal(t, dto.SearchOracleDatabaseAgreementsFilter{
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

func TestParseSearchOracleDatabaseAgreementsFilters_SuccessFull(t *testing.T) {
	r, err := http.NewRequest("GET", "/?agreement-id=foo&part-id=bar&item-description=boz&csi=pippo&metrics=pluto&reference-number=foobar&unlimited=false&catch-all=true&licenses-count-lte=10&licenses-count-gte=20&users-count-lte=5&users-count-gte=15&available-count-lte=3&available-count-gte=13", nil)
	require.NoError(t, err)
	filters, err := parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.NoError(t, err)
	assert.Equal(t, dto.SearchOracleDatabaseAgreementsFilter{
		AgreementID:       "foo",
		PartID:            "bar",
		ItemDescription:   "boz",
		CSI:               "pippo",
		Metric:            "pluto",
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

func TestParseSearchOracleDatabaseAgreementsFilters_Fail1(t *testing.T) {
	r, err := http.NewRequest("GET", "/?unlimited=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail2(t *testing.T) {
	r, err := http.NewRequest("GET", "/?catch-all=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail3(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-count-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail4(t *testing.T) {
	r, err := http.NewRequest("GET", "/?licenses-count-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail5(t *testing.T) {
	r, err := http.NewRequest("GET", "/?users-count-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail6(t *testing.T) {
	r, err := http.NewRequest("GET", "/?users-count-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail7(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-count-lte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	require.Error(t, err)
}

func TestParseSearchOracleDatabaseAgreementsFilters_Fail8(t *testing.T) {
	r, err := http.NewRequest("GET", "/?available-count-gte=asdfasdf", nil)
	require.NoError(t, err)
	_, err = parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddHostToAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddHostToAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddHostToAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrNotInClusterHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddHostToAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddHostToAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddHostToAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().RemoveHostFromAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveHostFromAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveHostFromAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveHostFromAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().RemoveHostFromAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(utils.AerrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveHostFromAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().RemoveHostFromAssociatedPart(utils.Str2oid("5f50a98611959b1baa17525e"), "foohost").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.RemoveHostFromAssociatedPart)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteAssociatedPartFromOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteAssociatedPartFromOracleDatabaseAgreement)
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteAssociatedPartFromOracleDatabaseAgreement)
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
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteAssociatedPartFromOracleDatabaseAgreement)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteAssociatedPartFromOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(utils.AerrOracleDatabaseAgreementNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteAssociatedPartFromOracleDatabaseAgreement)
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
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteAssociatedPartFromOracleDatabaseAgreement(utils.Str2oid("5f50a98611959b1baa17525e")).Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteAssociatedPartFromOracleDatabaseAgreement)
	req, err := http.NewRequest("DELETE", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5f50a98611959b1baa17525e",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
