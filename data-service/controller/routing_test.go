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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestAuthenticateMiddleware_Success(t *testing.T) {
	var err error

	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Config: config.Configuration{
			DataService: config.DataService{
				AgentUsername: "agent",
				AgentPassword: "p4ssW0rd",
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := ac.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Basic YWdlbnQ6cDRzc1cwcmQ=")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestAuthenticateMiddleware_Unauthorized(t *testing.T) {
	var err error

	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Config: config.Configuration{
			DataService: config.DataService{
				AgentUsername: "agent",
				AgentPassword: "p4ssW0rd",
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := ac.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Basic YWdlbnQ6VDBwb0wxbm8=")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
