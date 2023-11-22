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

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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

	group := model.Group{
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	returnAgr := group
	var err error
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

func TestInsertGroup_UnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.Group{
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

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
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

	group := model.Group{
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
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.UpdateGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(group), rr.Body.String())
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

	group := model.Group{
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
		"name": "Test",
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

func TestUpdateGroup_UnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	group := model.Group{
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
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.UpdateGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
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

	groups := []model.Group{
		{
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

func TestGetGroups_UnprocessableEntity(t *testing.T) {
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

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
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

	as.EXPECT().ListUsers().
		Return(nil, nil)

	as.EXPECT().DeleteGroup("Test").
		Return(nil)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.DeleteGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
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

	as.EXPECT().ListUsers().
		Return(nil, nil)

	aerr := utils.NewError(utils.ErrGroupNotFound, "test")
	as.EXPECT().DeleteGroup("Test").
		Return(aerr)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
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

func TestDeleteGroup_UnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().ListUsers().
		Return(nil, nil)

	as.EXPECT().DeleteGroup("Test").
		Return(errMock)

	req, err := http.NewRequest("DELETE", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.DeleteGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
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

	group := model.Group{
		Name:  "Test",
		Roles: []string{"role1", "role2"},
	}

	as.EXPECT().GetGroup("Test").
		Return(&group, nil)

	expBytes, err := json.Marshal(group)
	require.NoError(t, err)

	reader := bytes.NewReader(expBytes)
	req, err := http.NewRequest("GET", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.GetGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(&group), rr.Body.String())
}

func TestGetGroup_UnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetGroup("Test").
		Return(nil, errMock)

	req, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"name": "Test",
	})

	handler := http.HandlerFunc(ac.GetGroup)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
}
