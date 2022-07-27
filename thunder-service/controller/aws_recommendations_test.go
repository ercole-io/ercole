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
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAwsRecommendations_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-05-30T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	recommendation := model.AwsRecommendation{
		SeqValue:   999,
		ProfileID:  primitive.NewObjectID(),
		Category:   "",
		Suggestion: "",
		Name:       "",
		ResourceID: "",
		ObjectType: "",
		Details: []map[string]interface{}{
			{
				"Name":  "",
				"Value": "",
			},
		},
		CreatedAt: time.Date(2022, 5, 30, 0, 0, 1, 0, time.UTC),
	}

	var expectedRes []model.AwsRecommendation
	expectedRes = append(expectedRes, recommendation)

	as.EXPECT().GetAwsRecommendations().
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetAwsRecommendations)

	req, err := http.NewRequest("GET", "/aws/aws-recommendations", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"ids": "6140c473413cf9de756f9848"})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(map[string]interface{}{"recommendations": dto.ToAwsRecommendationsDto(expectedRes)}), rr.Body.String())
}

func TestGetAwsRecommendations_ClusterNotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-05-30T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetAwsRecommendations().Return(nil, utils.ErrClusterNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetAwsRecommendations)

	req, err := http.NewRequest("GET", "/aws/aws-recommendations", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"ids": "6140c473413cf9de756f9848"})

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

func TestGetAwsRecommendations_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-05-30T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().GetAwsRecommendations().Return(nil, errMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetAwsRecommendations)

	req, err := http.NewRequest("GET", "/aws/aws-recommendations", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"ids": "6140c473413cf9de756f9848"})

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

func TestGetAwsRecommendationsErrors_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockThunderServiceInterface(mockCtrl)
	ac := ThunderController{
		TimeNow: utils.Btc(utils.P("2022-05-30T12:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	recs := []model.AwsRecommendation{
		{
			SeqValue:  998,
			Category:  "",
			Errors:    []map[string]string{{"error": "error details"}},
			CreatedAt: time.Date(2022, 5, 30, 0, 0, 1, 0, time.UTC),
		},
	}

	as.EXPECT().GetAwsRecommendationsBySeqValue(uint64(998)).Return(recs, nil)

	expectedRes := dto.ToAwsRecommendationsErrorsDto(recs)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetAwsRecommendationsErrors)

	req, err := http.NewRequest("GET", "/aws/aws-recommendation-errors/{seqnum}", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"seqnum": "998"})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}
