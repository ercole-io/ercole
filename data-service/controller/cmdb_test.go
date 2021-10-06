// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/data-service/dto"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestCompareCmdbsInfo_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	cmdbInfo := dto.CmdbInfo{}
	as.EXPECT().CompareCmdbInfo(cmdbInfo).Return(nil)

	handler := http.HandlerFunc(ac.CompareCmdbInfo)

	cmdbInfoBytes, err := json.Marshal(cmdbInfo)
	require.NoError(t, err)

	reader := bytes.NewReader(cmdbInfoBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestCompareCmdbsInfo_BadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	req, err := http.NewRequest("POST", "/", strings.NewReader("asdf"))
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.CompareCmdbInfo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	actual := utils.ErrorResponseFE{}
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.Nil(t, err)

	assert.Equal(t, "Bad Request", actual.Message)
}

func TestCompareCmdbsInfo_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	cmdbInfo := dto.CmdbInfo{}
	as.EXPECT().CompareCmdbInfo(cmdbInfo).Return(aerrMock)

	cmdbInfoBytes, err := json.Marshal(cmdbInfo)
	require.NoError(t, err)

	reader := bytes.NewReader(cmdbInfoBytes)
	req, err := http.NewRequest("POST", "/", reader)
	require.NoError(t, err)

	handler := http.HandlerFunc(ac.CompareCmdbInfo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	actual := utils.ErrorResponseFE{}
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.Nil(t, err)

	assert.Equal(t, "mock", actual.Message)
}
