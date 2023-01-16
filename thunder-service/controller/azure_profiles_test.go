// Copyright (c) 2023 Sorint.lab S.p.A.
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

func TestAddAzureProfile_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	profile := model.AzureProfile{
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   &strClientSecretTestAdd,
		Selected:       false,
	}

	returnAgr := profile
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("aaaaaaaaaaaaaaaaaaaaaaaa")
	require.Nil(t, err)

	as.EXPECT().AddAzureProfile(profile).Return(&returnAgr, nil)

	agrBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.AddAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestAddAzureProfile_BadRequest_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongProfile := struct {
		Pippo string
		Pluto int
	}{
		Pippo: "pippo",
		Pluto: 42,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddAzureProfile)
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

func TestAddAzureProfile_BadRequest_HasID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	wrongProfile := model.AzureProfile{
		ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   &strClientSecretTestAdd,
		Selected:       false,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddAzureProfile)
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

func TestAddAzureProfile_BadRequest_SecretAccessKeyNull(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongProfile := model.AzureProfile{
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   nil,
		Selected:       false,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "ClientSecret must not be null", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestUpdateAzureProfile_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	profile := model.AzureProfile{
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   &strClientSecretTestAdd,
		Selected:       false,
	}

	returnAgr := profile
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("000000000000000000000000")
	require.Nil(t, err)

	as.EXPECT().UpdateAzureProfile(profile).
		Return(&profile, nil)

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(profile), rr.Body.String())
}

func TestUpdateAzureProfile_UnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	profile := model.AzureProfile{
		ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   &strClientSecretTestAdd,
		Selected:       false,
	}

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "hhhhhhhhhhhhhhhhhhhhhhhh",
	})

	handler := http.HandlerFunc(ac.UpdateAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "encoding/hex: invalid byte: U+0068 'h'", feErr.Error)
	assert.Equal(t, "Unprocessable Entity", feErr.Message)
}

func TestUpdateAzureProfile_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongProfile := struct {
		Pippo string
		Pluto int
	}{
		Pippo: "pippo",
		Pluto: 42,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.UpdateAzureProfile)
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

func TestUpdateAzureProfile_ObjectIdNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	profile := model.AzureProfile{
		ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   &strClientSecretTestAdd,
		Selected:       false,
	}

	returnAgr := profile
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("000000000000000000000001")
	require.Nil(t, err)

	proBytes, err := json.Marshal(returnAgr)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateAzureProfile)
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

func TestUpdateAzureProfile_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	profile := model.AzureProfile{
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   &strClientSecretTestAdd,
		Selected:       false,
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().UpdateAzureProfile(profile).
		Return(nil, aerr)

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestUpdateAzureProfile_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	profile := model.AzureProfile{
		TenantId:       "TestProfileAdd",
		ClientId:       "TestProfileAdd",
		SubscriptionId: "TestProfileAdd",
		Region:         "eu-frankfurt-testAdd",
		ClientSecret:   &strClientSecretTestAdd,
		Selected:       false,
	}

	as.EXPECT().UpdateAzureProfile(profile).
		Return(nil, errMock)

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateAzureProfile)
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

func TestGetAzureProfiles_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strClientSecretTestAdd = "ClientSecretTestAdd"
	profiles := []model.AzureProfile{
		{
			ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			TenantId:       "TestProfileAdd",
			ClientId:       "TestProfileAdd",
			SubscriptionId: "TestProfileAdd",
			Region:         "eu-frankfurt-testAdd",
			ClientSecret:   &strClientSecretTestAdd,
			Selected:       false,
		},
	}

	as.EXPECT().GetAzureProfiles().
		Return(profiles, nil)

	proBytes, err := json.Marshal(profiles)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("GET", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetAzureProfiles)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(profiles), rr.Body.String())
}

func TestGetAzureProfiles_ClusterNotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	profiles := make([]model.AzureProfile, 0)

	aerr := utils.NewError(utils.ErrClusterNotFound, "test")
	as.EXPECT().GetAzureProfiles().
		Return(profiles, aerr)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetAzureProfiles)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Cluster not found", feErr.Error)
	assert.Equal(t, "test", feErr.Message)
}

func TestGetAzureProfiles_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	profiles := make([]model.AzureProfile, 0)

	as.EXPECT().GetAzureProfiles().
		Return(profiles, errMock)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetAzureProfiles)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Internal Server Error", feErr.Message)
}

func TestDeleteAzureProfile_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteAzureProfile(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(nil)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteAzureProfile_BadRequest_HasWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "asdf",
	})

	handler := http.HandlerFunc(ac.DeleteAzureProfile)
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

func TestDeleteAzureProfile_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().DeleteAzureProfile(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(aerr)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteAzureProfile)
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

func TestDeleteAzureProfile_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteAzureProfile(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(errMock)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Internal Server Error", feErr.Message)
}

func TestSelectAzureProfile_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T15:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().SelectAzureProfile("aaaaaaaaaaaaaaaaaaaaaaaa", true).
		Return(nil)

	req, err := http.NewRequest("PUT", "/azure/profile-selection/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"profileid": "aaaaaaaaaaaaaaaaaaaaaaaa",
		"selected":  "true",
	})

	handler := http.HandlerFunc(ac.SelectAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestSelectAzureProfile_ClusterNotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T15:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().SelectAzureProfile("aaaaaaaaaaaaaaaaaaaaaaaa", true).
		Return(aerr)

	req, err := http.NewRequest("PUT", "/azure/profile-selection/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"profileid": "aaaaaaaaaaaaaaaaaaaaaaaa",
		"selected":  "true",
	})

	handler := http.HandlerFunc(ac.SelectAzureProfile)
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

func TestSelectAzureProfile_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-06-28T15:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().SelectAzureProfile("aaaaaaaaaaaaaaaaaaaaaaaa", true).
		Return(errMock)

	req, err := http.NewRequest("PUT", "/azure/profile-selection/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"profileid": "aaaaaaaaaaaaaaaaaaaaaaaa",
		"selected":  "true",
	})

	handler := http.HandlerFunc(ac.SelectAzureProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "MockError", feErr.Error)
	assert.Equal(t, "Internal Server Error", feErr.Message)
}
