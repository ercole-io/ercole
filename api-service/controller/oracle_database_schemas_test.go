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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestListOracleDatabaseSchemas_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	result := []dto.OracleDatabaseSchema{
		{
			Hostname:      "hostname",
			Indexes:       0,
			LOB:           0,
			Tables:        0,
			Total:         0,
			User:          "user",
			AccountStatus: "status",
		},
	}

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().
		ListOracleDatabaseSchemas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListOracleDatabaseSchemas)
	req, err := http.NewRequest("GET", "/", nil)

	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(result), rr.Body.String())
}

func TestGetOracleDatabaseSchemasXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: logger.NewLogger("TEST"),
	}

	sheet := "Schema"
	headers := []string{
		"Hostname",
		"Indexes",
		"LOB",
		"Tables",
		"Total",
		"User",
		"Account Status",
	}

	expectedRes, _ := exutils.NewXLSX(ac.Config, sheet, headers...)

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().CreateOracleDatabaseSchemasXlsx(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListOracleDatabaseSchemas)
	req, err := http.NewRequest("GET", "/schemas", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := excelize.OpenReader(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "Hostname", sp.GetCellValue("Schema", "A1"))
	assert.Equal(t, "Indexes", sp.GetCellValue("Schema", "B1"))
	assert.Equal(t, "LOB", sp.GetCellValue("Schema", "C1"))
	assert.Equal(t, "Tables", sp.GetCellValue("Schema", "D1"))
	assert.Equal(t, "Total", sp.GetCellValue("Schema", "E1"))
	assert.Equal(t, "User", sp.GetCellValue("Schema", "F1"))
	assert.Equal(t, "Account Status", sp.GetCellValue("Schema", "G1"))
}
