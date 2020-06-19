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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source ../service/service.go -destination=fake_service.go -package=controller

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

func TestHostDataInsertion_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAlertServiceInterface(mockCtrl)
	aqc := AlertQueueController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}
	as.EXPECT().HostDataInsertion(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(nil).Times(1)
	as.EXPECT().HostDataInsertion(gomock.Any()).Times(0)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aqc.HostDataInsertion)
	req, err := http.NewRequest("GET", "/queue/host-data-insertion/5dc3f534db7e81a98b726a52", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5dc3f534db7e81a98b726a52",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHostDataInsertion_RequestError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAlertServiceInterface(mockCtrl)
	aqc := AlertQueueController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}
	as.EXPECT().HostDataInsertion(gomock.Any()).Times(0)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aqc.HostDataInsertion)
	req, err := http.NewRequest("GET", "/queue/host-data-insertion/pippo", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "pippo",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestHostDataInsertion_ServiceError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAlertServiceInterface(mockCtrl)
	aqc := AlertQueueController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
		Log: utils.NewLogger("TEST"),
	}
	as.EXPECT().HostDataInsertion(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(aerrMock).Times(1)
	as.EXPECT().HostDataInsertion(gomock.Any()).Times(0)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aqc.HostDataInsertion)
	req, err := http.NewRequest("GET", "/queue/host-data-insertion/5dc3f534db7e81a98b726a52", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5dc3f534db7e81a98b726a52",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
