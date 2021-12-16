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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func TestUpdateHostInfo_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	raw, err := ioutil.ReadFile("../../fixture/test_dataservice_hostdata_v1_00.json")
	require.NoError(t, err)

	expectedHostDataBE := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	as.EXPECT().InsertHostData(expectedHostDataBE).Return(nil)

	handler := http.HandlerFunc(ac.InsertHostData)
	req, err := http.NewRequest("PUT", "/", bytes.NewReader(raw))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateHostInfo_FailBadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.InsertHostData)
	req, err := http.NewRequest("PUT", "/", &failingReader{})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateHostInfo_UnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().
		AlertInvalidHostData(gomock.Any(), nil).
		Do(func(err error, _ interface{}) {
			assert.ErrorIs(t, err, utils.ErrInvalidHostdata)
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.InsertHostData)
	req, err := http.NewRequest("PUT", "/", strings.NewReader("{asasdsad"))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestUpdateHostInfo_UnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().
		AlertInvalidHostData(gomock.Any(), gomock.Any()).
		Do(func(err error, hd interface{}) {
			assert.ErrorIs(t, err, utils.ErrInvalidHostdata)
			assert.NotNil(t, hd)
		})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.InsertHostData)
	req, err := http.NewRequest("PUT", "/", strings.NewReader("{}"))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestUpdateHostInfo_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	raw, err := ioutil.ReadFile("../../fixture/test_dataservice_hostdata_v1_00.json")
	require.NoError(t, err)

	expectedHostDataBE := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	as.EXPECT().InsertHostData(expectedHostDataBE).Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.InsertHostData)
	req, err := http.NewRequest("PUT", "/", bytes.NewReader(raw))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateHostInfo_SuccessAndSanitized(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	actual, err := ioutil.ReadFile("../../fixture/test_dataservice_hostdata_xss_v1_00.json")
	require.NoError(t, err)

	expected := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")
	as.EXPECT().InsertHostData(expected).Return(nil)

	handler := http.HandlerFunc(ac.InsertHostData)
	req, err := http.NewRequest("PUT", "/", bytes.NewReader(actual))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
