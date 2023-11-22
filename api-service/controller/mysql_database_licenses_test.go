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
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestUpdateMySqlLicenseIgnoredField_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	as.EXPECT().UpdateMySqlLicenseIgnoredField("serv123", "TEST123", false, "").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.UpdateMySqlLicenseIgnoredField)
	body := model.MySQLLicense{IgnoredComment: ""}
	req, err := http.NewRequest("PUT", "/hosts/serv123/technologies/mysql/databases/TEST123/false", bytes.NewReader([]byte(utils.ToJSON(body))))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "serv123",
		"dbname":   "TEST123",
		"ignored":  "false",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

}
