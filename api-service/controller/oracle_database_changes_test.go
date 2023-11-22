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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetOracleChanges_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	result := []dto.OracleChangesDto{
		{
			Hostname: "newdb",
			Databasenames: []dto.OracleChangesDBs{
				{
					Databasename: "pippodb",
					OracleChanges: []dto.OracleChangesGrowth{
						{
							DailyCPUUsage: 3.4,
							SegmentsSize:  50,
							Updated:       utils.P("2020-05-21T09:32:54.83Z"),
							DatafileSize:  8,
							Allocable:     129,
						},
						{
							DailyCPUUsage: 5.3,
							SegmentsSize:  100,
							Updated:       utils.P("2020-05-21T09:32:09.288Z"),
							DatafileSize:  10,
							Allocable:     129,
						},
						{
							DailyCPUUsage: 0.7,
							SegmentsSize:  3,
							Updated:       utils.P("2020-05-21T09:30:55.061Z"),
							DatafileSize:  6,
							Allocable:     129,
						},
					},
				},
			},
		},
	}

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().
		GetOracleChanges(dto.GlobalFilter{OlderThan: utils.MAX_TIME}, "newdb").
		Return(result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleChanges)
	req, err := http.NewRequest("GET", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "newdb",
	})

	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	expectedRes := result
	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}
