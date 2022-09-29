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

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestInsertRole_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	role := model.Role{
		Name: "Test",
	}

	returnAgr := role
	var err error
	require.Nil(t, err)

	as.EXPECT().InsertRole(role).
		Return(&returnAgr, nil)

	agrBytes, err := json.Marshal(role)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.InsertRole)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, utils.ToJSON(returnAgr), rr.Body.String())
}

func TestInsertRole_BadRequest_CantDecode(t *testing.T) {
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

	handler := http.HandlerFunc(ac.InsertRole)
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

func TestInsertRole_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	role := model.Role{
		Name: "Test",
	}

	as.EXPECT().InsertRole(role).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(role)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.InsertRole)
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

func TestUpdateRole_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	role := model.Role{
		Name: "Test",
	}

	as.EXPECT().UpdateRole(role).
		Return(&role, nil)

	agrBytes, err := json.Marshal(role)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.UpdateRole)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(role), rr.Body.String())
}

func TestUpdateRole_BadRequest_CantDecode(t *testing.T) {
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
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.UpdateRole)
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

func TestUpdateRole_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	role := model.Role{
		Name: "Test",
	}

	aerr := utils.NewError(utils.ErrRoleNotFound, "test")
	as.EXPECT().UpdateRole(role).
		Return(nil, aerr)

	agrBytes, err := json.Marshal(role)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.UpdateRole)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Role not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestUpdateRole_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	role := model.Role{
		Name: "Test",
	}

	as.EXPECT().UpdateRole(role).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(role)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.UpdateRole)
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

func TestGetRoles_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	roles := []model.Role{
		{
			Name: "Test",
		},
	}

	as.EXPECT().GetRoles().
		Return(roles, nil)

	expBytes, err := json.Marshal(roles)
	require.NoError(t, err)

	reader := bytes.NewReader(expBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetRoles)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expected := map[string]interface{}{
		"roles": roles,
	}
	assert.JSONEq(t, utils.ToJSON(expected), rr.Body.String())
}

func TestGetRoles_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetRoles().
		Return(nil, errMock)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetRoles)
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

func TestDeleteRole_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteRole("Test").
		Return(nil)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.DeleteRole)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteRole_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	aerr := utils.NewError(utils.ErrRoleNotFound, "test")
	as.EXPECT().DeleteRole("Test").
		Return(aerr)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.DeleteRole)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Role not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestDeleteRole_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteRole("Test").
		Return(errMock)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.DeleteRole)
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

func TestGetRole_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	role := model.Role{
		Name: "Test",
	}

	as.EXPECT().GetRole("Test").
		Return(&role, nil)

	expBytes, err := json.Marshal(role)
	require.NoError(t, err)

	reader := bytes.NewReader(expBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.GetRole)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(&role), rr.Body.String())
}

func TestGetRole_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetRole("Test").
		Return(nil, errMock)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.GetRole)
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
