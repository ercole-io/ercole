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
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOracleOptionList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	result := []dto.OracleDatabaseFeatureUsageStatDto{
		{
			Hostname:     "hostname",
			Databasename: "databasename",
			OracleDatabaseFeatureUsageStat: model.OracleDatabaseFeatureUsageStat{
				Product:          "Diagnostics Pack",
				Feature:          "ADDM",
				DetectedUsages:   91,
				CurrentlyUsed:    false,
				FirstUsageDate:   utils.P("2020-05-04T14:09:46.608Z").UTC(),
				LastUsageDate:    utils.P("2020-05-05T09:04:25.000Z").UTC(),
				ExtraFeatureInfo: "",
			},
		},
	}

	as.EXPECT().
		GetOracleOptionList().
		Return(result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleOptionList)
	req, err := http.NewRequest("GET", "/", nil)

	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	expectedRes := result
	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetOracleOptionListXLSX_Success(t *testing.T) {
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

	sheet := "Options"
	headers := []string{
		"Hostname",
		"DB Name",
		"First",
		"Last",
		"Detected",
		"Prod",
		"Currently",
		"Extra",
		"Feature",
	}

	expectedRes, _ := exutils.NewXLSX(ac.Config, sheet, headers...)

	as.EXPECT().CreateGetOracleOptionListXLSX().Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleOptionList)
	req, err := http.NewRequest("GET", "/option-list", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := excelize.OpenReader(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "Hostname", sp.GetCellValue("Options", "A1"))
	assert.Equal(t, "DB Name", sp.GetCellValue("Options", "B1"))
	assert.Equal(t, "First", sp.GetCellValue("Options", "C1"))
	assert.Equal(t, "Last", sp.GetCellValue("Options", "D1"))
	assert.Equal(t, "Detected", sp.GetCellValue("Options", "E1"))
	assert.Equal(t, "Prod", sp.GetCellValue("Options", "F1"))
	assert.Equal(t, "Currently", sp.GetCellValue("Options", "G1"))
	assert.Equal(t, "Extra", sp.GetCellValue("Options", "H1"))
	assert.Equal(t, "Feature", sp.GetCellValue("Options", "I1"))
}
