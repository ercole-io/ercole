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

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestListOracleGrantDbaByHostname_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	gdRes := []dto.OracleGrantDbaDto{
		{
			Hostname:     "hostname",
			Databasename: "databasename",
			OracleGrantDba: model.OracleGrantDba{
				Grantee:     "test#001",
				AdminOption: "yes",
				DefaultRole: "no",
			},
		},
	}

	as.EXPECT().
		ListOracleGrantDbaByHostname("hostname", filter).
		Return(gdRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListOracleGrantDbaByHostname)
	req, err := http.NewRequest("GET", "/", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "hostname",
	})

	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	expectedRes := gdRes
	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}
