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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestAddMySQLContract_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contract := model.MySQLContract{
		Type:             model.MySQLContractTypeCluster,
		ContractID:       "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{"pippo", "pluto"},
		Hosts:            []string{"topolino", "minnie"},
	}

	returnAgr := contract
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("aaaaaaaaaaaaaaaaaaaaaaaa")
	require.Nil(t, err)

	as.EXPECT().AddMySQLContract(contract).
		Return(&returnAgr, nil)

	agrBytes, err := json.Marshal(contract)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddMySQLContract)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, utils.ToJSON(returnAgr), rr.Body.String())
}

func TestAddMySQLContract_BadRequest_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
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

	handler := http.HandlerFunc(ac.AddMySQLContract)
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

func TestAddMySQLContract_BadRequest_HasID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongAgr := model.MySQLContract{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             "",
		ContractID:       "agr01",
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

	handler := http.HandlerFunc(ac.AddMySQLContract)
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

func TestAddMySQLContract_BadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongAgr := model.MySQLContract{
		Type: "",
		// ContractID:      "agr01",
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

	handler := http.HandlerFunc(ac.AddMySQLContract)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Contract isn't valid", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestAddMySQLContract_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contract := model.MySQLContract{
		Type:             model.MySQLContractTypeCluster,
		ContractID:       "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	as.EXPECT().AddMySQLContract(contract).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(contract)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddMySQLContract)
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

func TestUpdateMySQLContract_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contract := model.MySQLContract{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             model.MySQLContractTypeCluster,
		ContractID:       "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	as.EXPECT().UpdateMySQLContract(contract).
		Return(&contract, nil)

	agrBytes, err := json.Marshal(contract)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLContract)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(contract), rr.Body.String())
}

func TestUpdateMySQLContract_BadRequest_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
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

	handler := http.HandlerFunc(ac.UpdateMySQLContract)
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

func TestUpdateMySQLContract_BadRequest_HasWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongAgr := model.MySQLContract{
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

	handler := http.HandlerFunc(ac.UpdateMySQLContract)
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

func TestUpdateMySQLContract_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contract := model.MySQLContract{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             model.MySQLContractTypeCluster,
		ContractID:       "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().UpdateMySQLContract(contract).
		Return(nil, aerr)

	agrBytes, err := json.Marshal(contract)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLContract)
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

func TestUpdateMySQLContract_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contract := model.MySQLContract{
		ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Type:             model.MySQLContractTypeCluster,
		ContractID:       "agr01",
		CSI:              "csi01",
		NumberOfLicenses: 42,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	as.EXPECT().UpdateMySQLContract(contract).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(contract)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateMySQLContract)
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

func TestGetMySQLContracts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	contracts := []model.MySQLContract{
		{
			ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			Type:             "hosts",
			NumberOfLicenses: 7,
			Clusters:         []string{},
			Hosts:            []string{"pippo", "pluto"},
		},
	}

	as.EXPECT().GetMySQLContracts(gomock.Any()).
		Return(contracts, nil)

	expBytes, err := json.Marshal(contracts)
	require.NoError(t, err)

	reader := bytes.NewReader(expBytes)
	req, err := http.NewRequest("GET", "/?location=Italy", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetMySQLContracts)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expected := map[string]interface{}{
		"contracts": contracts,
	}
	assert.JSONEq(t, utils.ToJSON(expected), rr.Body.String())
}

func TestGetMySQLContracts_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetMySQLContracts(gomock.Any()).
		Return(nil, errMock)

	req, err := http.NewRequest("GET", "/?environment=TEST", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetMySQLContracts)
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

func TestGetMySQLContractsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	xlsx := excelize.File{}

	as.EXPECT().
		GetMySQLContractsAsXLSX(gomock.Any()).
		Return(&xlsx, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetMySQLContracts)
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	_, err = excelize.OpenReader(rr.Body)
	require.NoError(t, err)
}

func TestGetMySQLContractsXLSX_InternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	as.EXPECT().
		GetMySQLContractsAsXLSX(gomock.Any()).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetMySQLContracts)
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteMySQLContract_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteMySQLContract(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(nil)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLContract)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteMySQLContract_BadRequest_HasWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "asdf",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLContract)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Can't decode id: "+primitive.ErrInvalidHex.Error(), feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
}

func TestDeleteMySQLContract_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().DeleteMySQLContract(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(aerr)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLContract)
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

func TestDeleteMySQLContract_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteMySQLContract(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(errMock)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteMySQLContract)
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
