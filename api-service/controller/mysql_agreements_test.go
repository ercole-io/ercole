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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddMySQLAgreement_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	agreement := model.MySQLAgreement{
		Type:             model.MySQLAgreementTypeCluster,
		AgreementID:      "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{"pippo", "pluto"},
		Hosts:            []string{"topolino", "minnie"},
	}

	returnAgr := agreement
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("aaaaaaaaaaaaaaaaaaaaaaaa")
	require.Nil(t, err)

	as.EXPECT().AddMySQLAgreement(agreement).
		Return(&returnAgr, nil)

	agrBytes, err := json.Marshal(agreement)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, utils.ToJSON(returnAgr), rr.Body.String())
}

func TestAddMySQLAgreement_BadRequest_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	wrongAgr := struct {
		Pippo string
		Pluto int
	}{
		Pippo: "pippo",
		Pluto: 42,
	}

	agrBytes, err := json.Marshal(wrongAgr)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "json: unknown field \"Pippo\"", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestAddMySQLAgreement_BadRequest_HasID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	wrongAgr := model.MySQLAgreement{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             "",
		AgreementID:      "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 0,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	agrBytes, err := json.Marshal(wrongAgr)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "ID must be empty", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestAddMySQLAgreement_BadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	wrongAgr := model.MySQLAgreement{
		Type: "",
		// AgreementID:      "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 0,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	agrBytes, err := json.Marshal(wrongAgr)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Agreement isn't valid", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestAddMySQLAgreement_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	agreement := model.MySQLAgreement{
		Type:             model.MySQLAgreementTypeCluster,
		AgreementID:      "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	as.EXPECT().AddMySQLAgreement(agreement).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(agreement)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Internal Server Error", feErr.Message)
}

func TestUpdateMySQLAgreement_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	agreement := model.MySQLAgreement{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             model.MySQLAgreementTypeCluster,
		AgreementID:      "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	as.EXPECT().UpdateMySQLAgreement(agreement).
		Return(&agreement, nil)

	agrBytes, err := json.Marshal(agreement)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(agreement), rr.Body.String())
}

func TestUpdateMySQLAgreement_BadRequest_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	wrongAgr := struct {
		Pippo string
		Pluto int
	}{
		Pippo: "pippo",
		Pluto: 42,
	}

	agrBytes, err := json.Marshal(wrongAgr)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "json: unknown field \"Pippo\"", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestUpdateMySQLAgreement_BadRequest_HasWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	wrongAgr := model.MySQLAgreement{
		Type:             "",
		NumberOfLicenses: 0,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	agrBytes, err := json.Marshal(wrongAgr)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Object ID does not correspond", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestUpdateMySQLAgreement_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	agreement := model.MySQLAgreement{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             model.MySQLAgreementTypeCluster,
		AgreementID:      "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().UpdateMySQLAgreement(agreement).
		Return(nil, aerr)

	agrBytes, err := json.Marshal(agreement)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestUpdateMySQLAgreement_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	agreement := model.MySQLAgreement{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             model.MySQLAgreementTypeCluster,
		AgreementID:      "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	as.EXPECT().UpdateMySQLAgreement(agreement).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(agreement)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Internal Server Error", feErr.Message)
}

func TestGetMySQLAgreements_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	agreements := []model.MySQLAgreement{
		{
			ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			Type:             "hosts",
			NumberOfLicenses: 7,
			Clusters:         []string{},
			Hosts:            []string{"pippo", "pluto"},
		},
	}

	as.EXPECT().GetMySQLAgreements().
		Return(agreements, nil)

	expBytes, err := json.Marshal(agreements)
	require.NoError(t, err)

	reader := bytes.NewReader(expBytes)
	req, err := http.NewRequest("GET", "/?location=Italy", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetMySQLAgreements)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expected := map[string]interface{}{
		"agreements": agreements,
	}
	assert.JSONEq(t, utils.ToJSON(expected), rr.Body.String())
}

func TestGetMySQLAgreements_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().GetMySQLAgreements().
		Return(nil, errMock)

	req, err := http.NewRequest("GET", "/?environment=TEST", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetMySQLAgreements)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Internal Server Error", feErr.Message)
}

func TestDeleteMySQLAgreement_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteMySQLAgreement(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(nil)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteMySQLAgreement_BadRequest_HasWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "asdf",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Can't decode id: encoding/hex: invalid byte: U+0073 's'", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
}

func TestDeleteMySQLAgreement_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().DeleteMySQLAgreement(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(aerr)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestDeleteMySQLAgreement_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteMySQLAgreement(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(errMock)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLAgreement)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Internal Server Error", feErr.Message)
}
