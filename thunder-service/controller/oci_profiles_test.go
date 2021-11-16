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

func TestAddOciProfile_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	returnAgr := profile
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("aaaaaaaaaaaaaaaaaaaaaaaa")
	require.Nil(t, err)

	as.EXPECT().AddOciProfile(profile).
		Return(&returnAgr, nil)

	agrBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(agrBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.AddOciProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, utils.ToJSON(returnAgr), rr.Body.String())
}

func TestAddOciProfile_BadRequest_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
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

	handler := http.HandlerFunc(ac.AddOciProfile)
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

func TestAddOciProfile_BadRequest_HasID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	wrongProfile := model.OciProfile{
		ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddOciProfile)
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

func TestAddOciProfile_BadRequest_PrivateKeyNull(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	wrongProfile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenamcy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     nil,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddOciProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "PrivateKey must not be null", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestAddOciProfile_BadRequest_ProfileNotValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	wrongProfile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.ten.testAdd",
		UserOCID:       "ocid1.usr.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddOciProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Profile configuration isn't valid", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestAddOciProfile_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	as.EXPECT().AddOciProfile(profile).
		Return(nil, errMock)

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.AddOciProfile)
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

func TestUpdateOciProfile_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	returnAgr := profile
	var err error
	returnAgr.ID, err = primitive.ObjectIDFromHex("000000000000000000000000")
	require.Nil(t, err)

	as.EXPECT().UpdateOciProfile(profile).
		Return(&profile, nil)

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateOciProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, utils.ToJSON(profile), rr.Body.String())
}

func TestUpdateOciProfile_UnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "hhhhhhhhhhhhhhhhhhhhhhhh",
	})

	handler := http.HandlerFunc(ac.UpdateOciProfile)
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

func TestUpdateOciProfile_CantDecode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
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

	handler := http.HandlerFunc(ac.UpdateOciProfile)
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

func TestUpdateOciProfile_ObjectIdNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
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

	handler := http.HandlerFunc(ac.UpdateOciProfile)
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

func TestUpdateOciProfile_BadRequest_ProfileNotValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	wrongProfile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.ten.testAdd",
		UserOCID:       "ocid1.usr.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	proBytes, err := json.Marshal(wrongProfile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateOciProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Some profile fields are not valid", feErr.Error)
	assert.Equal(t, "Bad Request", feErr.Message)
}

func TestUpdateOciProfile_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().UpdateOciProfile(profile).
		Return(nil, aerr)

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateOciProfile)
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

func TestUpdateOciProfile_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profile := model.OciProfile{
		Profile:        "TestProfileAdd",
		TenancyOCID:    "ocid1.tenancy.testAdd",
		UserOCID:       "ocid1.user.testAdd",
		KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
		Region:         "eu-frankfurt-testAdd",
		PrivateKey:     &strPrivateKeyTestAdd,
	}

	as.EXPECT().UpdateOciProfile(profile).
		Return(nil, errMock)

	proBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("PUT", "/", reader)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "000000000000000000000000",
	})

	handler := http.HandlerFunc(ac.UpdateOciProfile)
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

func TestGetOciProfiles_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var strPrivateKeyTestAdd = "PrivateKeyTestAdd"
	profiles := []model.OciProfile{
		{
			ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			Profile:        "TestProfileAdd",
			TenancyOCID:    "ocid1.tenancy.testAdd",
			UserOCID:       "ocid1.user.testAdd",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:66:66:66:66",
			Region:         "eu-frankfurt-testAdd",
			PrivateKey:     &strPrivateKeyTestAdd,
		},
	}

	as.EXPECT().GetOciProfiles().
		Return(profiles, nil)

	proBytes, err := json.Marshal(profiles)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("GET", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetOciProfiles)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(profiles), rr.Body.String())
}

func TestGetOciProfiles_ClusterNotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	profiles := make([]model.OciProfile, 0)

	aerr := utils.NewError(utils.ErrClusterNotFound, "test")
	as.EXPECT().GetOciProfiles().
		Return(profiles, aerr)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetOciProfiles)
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

func TestGetOciProfiles_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	profiles := make([]model.OciProfile, 0)

	as.EXPECT().GetOciProfiles().
		Return(profiles, errMock)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetOciProfiles)
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

func TestDeleteOciProfile_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteOciProfile(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(nil)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteOciProfile)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteOciProfile_BadRequest_HasWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "asdf",
	})

	handler := http.HandlerFunc(ac.DeleteOciProfile)
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

func TestDeleteOciProfile_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	aerr := utils.NewError(utils.ErrNotFound, "test")
	as.EXPECT().DeleteOciProfile(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(aerr)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteOciProfile)
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

func TestDeleteOciProfile_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().DeleteOciProfile(utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")).
		Return(errMock)

	req, err := http.NewRequest("DELETE", "/", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "aaaaaaaaaaaaaaaaaaaaaaaa",
	})

	handler := http.HandlerFunc(ac.DeleteOciProfile)
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
