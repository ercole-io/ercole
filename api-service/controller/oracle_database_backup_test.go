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

func TestGetOracleBackupList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	result := []dto.OracleDatabaseBackupDto{
		{
			Hostname:     "hostname",
			Databasename: "databasename",
			OracleDatabaseBackup: model.OracleDatabaseBackup{
				BackupType: "Archivelog",
				Hour:       "01:30",
				WeekDays:   []string{"Wednesday"},
				AvgBckSize: 13.0,
				Retention:  "1 NUMBERS",
			},
		},
	}

	as.EXPECT().
		GetOracleBackupList().
		Return(result, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleBackupList)
	req, err := http.NewRequest("GET", "/", nil)

	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	expectedRes := result
	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetOracleBackupListXLSX_Success(t *testing.T) {
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

	sheet := "Backups"
	headers := []string{
		"Hostname",
		"DB Name",
		"Days of the Week",
		"Hour",
		"Type",
		"Average",
		"RMAN",
	}

	expectedRes, _ := exutils.NewXLSX(ac.Config, sheet, headers...)

	as.EXPECT().CreateGetOracleBackupListXLSX().Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetOracleBackupList)
	req, err := http.NewRequest("GET", "/backup-list", nil)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := excelize.OpenReader(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "Hostname", sp.GetCellValue("Backups", "A1"))
	assert.Equal(t, "DB Name", sp.GetCellValue("Backups", "B1"))
	assert.Equal(t, "Days of the Week", sp.GetCellValue("Backups", "C1"))
	assert.Equal(t, "Hour", sp.GetCellValue("Backups", "D1"))
	assert.Equal(t, "Type", sp.GetCellValue("Backups", "E1"))
	assert.Equal(t, "Average", sp.GetCellValue("Backups", "F1"))
	assert.Equal(t, "RMAN", sp.GetCellValue("Backups", "G1"))
}
