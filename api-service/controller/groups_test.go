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
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestInsertGroup_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.GroupType{
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	returnAgr := group
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("aaaaaaaaaaaaaaaaaaaaaaaa")
	require.Nil(t, err)

	as.EXPECT().InsertGroup(group).
		Return(&returnAgr, nil)

	agrBytes, err := json.Marshal(group)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.InsertGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, utils.ToJSON(returnAgr), rr.Body.String())
}

func TestInsertGroup_BadRequest_CantDecode(t *testing.T) {
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

	handler := http.HandlerFunc(ac.InsertGroup)
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

func TestInsertGroup_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.GroupType{
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	as.EXPECT().InsertGroup(group).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(group)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.InsertGroup)
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

func TestUpdateGroup_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.GroupType{
		ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	as.EXPECT().UpdateGroup(group).
		Return(&group, nil)

	agrBytes, err := json.Marshal(group)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(group), rr.Body.String())
}

func TestUpdateGroup_BadRequest_CantDecode(t *testing.T) {
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

	handler := http.HandlerFunc(ac.UpdateGroup)
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

func TestUpdateGroup_BadRequest_HasWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongAgr := model.GroupType{
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	agrBytes, err := json.Marshal(wrongAgr)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateGroup)
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

func TestUpdateGroup_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.GroupType{
		ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	aerr := utils.NewError(utils.ErrGroupNotFound, "test")
	as.EXPECT().UpdateGroup(group).
		Return(nil, aerr)

	agrBytes, err := json.Marshal(group)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Group not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestUpdateGroup_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.GroupType{
		ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	as.EXPECT().UpdateGroup(group).
		Return(nil, errMock)

	agrBytes, err := json.Marshal(group)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateGroup)
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

func TestGetGroups_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	groups := []model.GroupType{
		{
			ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			Name:  "Test",
			Roles: []string{"role1", "role2"},
		},
	}

	as.EXPECT().GetGroups().
		Return(groups, nil)

	expBytes, err := json.Marshal(groups)
	require.NoError(t, err)

	reader := bytes.NewReader(expBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetGroups)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expected := map[string]interface{}{
		"groups": groups,
	}
	assert.JSONEq(t, utils.ToJSON(expected), rr.Body.String())
}

func TestGetGroups_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetGroups().
		Return(nil, errMock)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetGroups)
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

func TestDeleteGroup_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteGroup(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(nil)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteGroup_BadRequest_HasWrongID(t *testing.T) {
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

	handler := http.HandlerFunc(ac.DeleteGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "the provided hex string is not a valid ObjectID", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
}

func TestDeleteGroup_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	aerr := utils.NewError(utils.ErrGroupNotFound, "test")
	as.EXPECT().DeleteGroup(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(aerr)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Group not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestDeleteGroup_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteGroup(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(errMock)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteGroup)
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

func TestGetGroup_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.GroupType{
		ID:    utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	as.EXPECT().GetGroup(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(&group, nil)

	expBytes, err := json.Marshal(group)
	require.NoError(t, err)

	reader := bytes.NewReader(expBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.GetGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(&group), rr.Body.String())
}

func TestGetGroup_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetGroup(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(nil, errMock)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.GetGroup)
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
