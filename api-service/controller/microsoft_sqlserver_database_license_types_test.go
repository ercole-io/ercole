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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetSqlServerDatabaseLicenseTypes_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	ltRes := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "DG7GMGF0FLR2-0002",
			ItemDescription: "SQL Server 2019 Standard Core - 2 Core License Pack",
			Edition:         "STD",
			Version:         "2019",
		},
		{
			ID:              "DG7GMGF0FKZV-0001",
			ItemDescription: "SQL Server 2019 Enterprise Core - 2 Core License Pack",
			Edition:         "ENT",
			Version:         "2019",
		},
	}

	as.EXPECT().GetSqlServerDatabaseLicenseTypes().Return(ltRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetSqlServerDatabaseLicenseTypes)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	expectedRes := map[string]interface{}{
		"license-types": ltRes,
	}
	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetSqlServerDatabaseLicenseTypes_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	as.EXPECT().
		GetSqlServerDatabaseLicenseTypes().
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetSqlServerDatabaseLicenseTypes)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
