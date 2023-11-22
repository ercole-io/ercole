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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetOciObjects_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	objects := []model.OciObjects{
		{
			ID:        utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			ProfileID: "TestProfile1",
			Objects: []model.OciObject{
				{
					ObjectName:   "",
					ObjectNumber: 0,
				},
			},
			CreatedAt: time.Date(2022, 5, 27, 0, 0, 1, 0, time.UTC),
			Error:     "TestError1",
		},
	}

	as.EXPECT().GetOciObjects().
		Return(objects, nil)

	proBytes, err := json.Marshal(objects)
	require.NoError(t, err)

	reader := bytes.NewReader(proBytes)
	req, err := http.NewRequest("GET", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetOciObjects)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(objects), rr.Body.String())
}

func TestGetOciObjects_ClusterNotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	objects := make([]model.OciObjects, 0)

	aerr := utils.NewError(utils.ErrClusterNotFound, "test")
	as.EXPECT().GetOciObjects().
		Return(objects, aerr)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetOciObjects)
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

func TestGetOciObjects_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	objects := make([]model.OciObjects, 0)

	as.EXPECT().GetOciObjects().
		Return(objects, errMock)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.GetOciObjects)
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
