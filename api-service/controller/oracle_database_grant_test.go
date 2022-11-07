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
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
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

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

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

func TestGetOracleGrantDbaXLSX_Success(t *testing.T) {
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

	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	sheet := "DBA Roles"
	headers := []string{
		"Hostname",
		"DB Name",
		"Grantee",
		"Admin Option",
		"Default Role",
	}

	expectedRes, _ := exutils.NewXLSX(ac.Config, sheet, headers...)

	var user interface{}
	var locations []string

	as.EXPECT().
		ListLocations(user).
		Return(locations, nil)

	as.EXPECT().CreateOracleGrantDbaXlsx("", filter).Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListOracleGrantDbaByHostname)
	req, err := http.NewRequest("GET", "/grant-dba", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := excelize.OpenReader(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "Hostname", sp.GetCellValue("DBA Roles", "A1"))
	assert.Equal(t, "DB Name", sp.GetCellValue("DBA Roles", "B1"))
	assert.Equal(t, "Grantee", sp.GetCellValue("DBA Roles", "C1"))
	assert.Equal(t, "Admin Option", sp.GetCellValue("DBA Roles", "D1"))
	assert.Equal(t, "Default Role", sp.GetCellValue("DBA Roles", "E1"))
}
