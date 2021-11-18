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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	model "github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var errMock error = errors.New("MockError")
var aerrMock error = utils.NewError(errMock, "mock")

func TestGetOciRecommendations_StatusNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("StatusNotFound", func(t *testing.T) {
		as := NewMockThunderServiceInterface(mockCtrl)
		ac := ThunderController{
			TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     logger.NewLogger("TEST"),
		}

		var strProfiles = []string{"6140c473413cf9de756f9848"}
		as.EXPECT().GetOciRecommendations(strProfiles).Return(nil, utils.ErrClusterNotFound)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetOciRecommendations)

		req, err := http.NewRequest("GET", "/oracle-cloud/recommendations", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"ids": "6140c473413cf9de756f9848"})

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

		var feErr utils.ErrorResponseFE
		decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&feErr)
		require.NoError(t, err)

		fmt.Println("Error = ", feErr.Error, " -- Message = ", feErr.Message)
		assert.Equal(t, "Cluster not found", feErr.Error)
		assert.Equal(t, "Not Found", feErr.Message)

	})
}

func TestGetOciRecommendations_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("StatusNotFound", func(t *testing.T) {
		as := NewMockThunderServiceInterface(mockCtrl)
		ac := ThunderController{
			TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     logger.NewLogger("TEST"),
		}

		var strProfiles = []string{"6140c473413cf9de756f9848"}
		as.EXPECT().GetOciRecommendations(strProfiles).Return(nil, errMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetOciRecommendations)

		req, err := http.NewRequest("GET", "/oracle-cloud/recommendations", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"ids": "6140c473413cf9de756f9848"})

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)

		var feErr utils.ErrorResponseFE
		decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&feErr)
		require.NoError(t, err)

		fmt.Println("Error = ", feErr.Error, " -- Message = ", feErr.Message)
		assert.Equal(t, "MockError", feErr.Error)
		assert.Equal(t, "Internal Server Error", feErr.Message)

	})
}

func TestGetOciRecommendations_BadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("BadRequest", func(t *testing.T) {
		as := NewMockThunderServiceInterface(mockCtrl)
		ac := ThunderController{
			TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     logger.NewLogger("TEST"),
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetOciRecommendations)

		req, err := http.NewRequest("GET", "/oracle-cloud/recommendations", nil)
		require.NoError(t, err)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var feErr utils.ErrorResponseFE
		decoder := json.NewDecoder(bytes.NewReader(rr.Body.Bytes()))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&feErr)
		require.NoError(t, err)

		assert.Equal(t, "Ids not present or malformed", feErr.Error)
		assert.Equal(t, "Bad Request", feErr.Message)

	})
}

func TestGetOciRecommendations_InvalidProfileId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("BadRequest", func(t *testing.T) {
		as := NewMockThunderServiceInterface(mockCtrl)
		ac := ThunderController{
			TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     logger.NewLogger("TEST"),
		}

		expectedRes := map[string]interface{}{
			"error":          "invalid profile id",
			"sourceFilename": "c:/Workspace/GitHub/ercole/thunder-service/controller/oci_recommendations.go",
			"lineNumber":     43,
			"message":        "Bad Request",
		}

		var strProfiles = []string{"aaa", "bbb", "ccc"}
		as.EXPECT().GetOciRecommendations(strProfiles).Return(nil, utils.ErrInvalidProfileId)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetOciRecommendations)

		req, err := http.NewRequest("GET", "/oracle-cloud/recommendations", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"ids": "aaa,bbb,ccc"})

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
	})
}

func TestGetOciRecommendations_PartialContent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("BadRequest", func(t *testing.T) {
		as := NewMockThunderServiceInterface(mockCtrl)
		ac := ThunderController{
			TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     logger.NewLogger("TEST"),
		}

		var strError = "1 error occurred: 'invalid profile id aaa'"
		var mockError error = errors.New("1 error occurred: 'invalid profile id aaa'")

		var recommendations []model.Recommendation
		recommendation := model.Recommendation{
			TenancyOCID:         "ocid1.tenancy.oc1..aaaaaaaazizzbqqbjv2se3y3fvm5osfumnorh32nznanirqoju3uks4buh4q",
			Name:                "cost-management-load-balancer-underutilized-name",
			NumPending:          "1",
			EstimatedCostSaving: "1.92",
			Status:              "PENDING",
			Importance:          "MODERATE",
			RecommendationId:    "ocid1.optimizerrecommendation.oc1..aaaaaaaanptucivsbo24dkml7exd2fldalykkz52gld3nk3jt5c2e4ydgzjq",
		}

		recommendations = append(recommendations, recommendation)

		expectedRes := map[string]interface{}{
			"recommendations": recommendations,
			"error":           strError,
		}

		var strProfiles = []string{"6140c473413cf9de756f9848", "bbb", "ccc"}
		as.EXPECT().GetOciRecommendations(strProfiles).Return(recommendations, mockError)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetOciRecommendations)

		req, err := http.NewRequest("GET", "/oracle-cloud/recommendations", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"ids": "6140c473413cf9de756f9848,bbb,ccc"})

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusPartialContent, rr.Code)
		assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
	})
}

func TestGetOciRecommendations_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Success_NoContent", func(t *testing.T) {
		as := NewMockThunderServiceInterface(mockCtrl)
		ac := ThunderController{
			TimeNow: utils.Btc(utils.P("2021-11-08T12:02:03Z")),
			Service: as,
			Config:  config.Configuration{},
			Log:     logger.NewLogger("TEST"),
		}

		recommendation := model.Recommendation{
			TenancyOCID:         "ocid1.tenancy.oc1..aaaaaaaazizzbqqbjv2se3y3fvm5osfumnorh32nznanirqoju3uks4buh4q",
			Name:                "cost-management-load-balancer-underutilized-name",
			NumPending:          "1",
			EstimatedCostSaving: "1.92",
			Status:              "PENDING",
			Importance:          "MODERATE",
			RecommendationId:    "ocid1.optimizerrecommendation.oc1..aaaaaaaanptucivsbo24dkml7exd2fldalykkz52gld3nk3jt5c2e4ydgzjq",
		}

		var expectedRes []model.Recommendation
		var strProfiles = []string{"6140c473413cf9de756f9848"}
		expectedRes = append(expectedRes, recommendation)
		as.EXPECT().GetOciRecommendations(strProfiles).Return(expectedRes, nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ac.GetOciRecommendations)

		req, err := http.NewRequest("GET", "/oracle-cloud/recommendations", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"ids": "6140c473413cf9de756f9848"})

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, utils.ToJSON(map[string]interface{}{"recommendations": expectedRes}), rr.Body.String())
	})
}
