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

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	model "github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetOciRecommendationErrors_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-05-30T15:15:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	recError := model.OciRecommendationError{

		SeqValue:  999,
		ProfileID: "",
		Category:  "",
		CreatedAt: time.Date(2022, 5, 30, 0, 0, 1, 0, time.UTC),
		Error:     "",
	}

	var expectedRes []model.OciRecommendationError
	expectedRes = append(expectedRes, recError)
	var seqNum uint64 = 999
	as.EXPECT().GetOciRecommendationErrors(seqNum).Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOciRecommendationErrors)

	req, err := http.NewRequest("GET", "/last-oci-recommendation-errors", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"seqnum": "999"})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetOciRecommendationErrors_ClusterNotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-05-30T15:15:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var seqNum uint64 = 999
	as.EXPECT().GetOciRecommendationErrors(seqNum).Return(nil, utils.ErrClusterNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOciRecommendationErrors)

	req, err := http.NewRequest("GET", "/oracle-cloud/last-oci-recommendation-errors", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"seqnum": "999"})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var feErr utils.ErrorResponseFE
	decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&feErr)
	require.NoError(t, err)

	assert.Equal(t, "Cluster not found", feErr.Error)
	assert.Equal(t, "Not Found", feErr.Message)
}

func TestGetOciRecommendationErrors_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-05-30T15:15:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	var seqNum uint64 = 999
	as.EXPECT().GetOciRecommendationErrors(seqNum).Return(nil, errMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOciRecommendationErrors)

	req, err := http.NewRequest("GET", "/oracle-cloud/last-oci-recommendation-errors", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"seqnum": "999"})

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
