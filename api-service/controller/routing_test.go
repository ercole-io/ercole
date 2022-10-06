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
	"testing"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/auth"
	"github.com/ercole-io/ercole/v2/api-service/service"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/gorilla/mux"
)

func TestAPIController_setupProtectedRoutes(t *testing.T) {
	type fields struct {
		Config        config.Configuration
		Service       service.APIServiceInterface
		TimeNow       func() time.Time
		Log           logger.Logger
		Authenticator []auth.AuthenticationProvider
	}
	type args struct {
		router *mux.Router
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := &APIController{
				Config:        tt.fields.Config,
				Service:       tt.fields.Service,
				TimeNow:       tt.fields.TimeNow,
				Log:           tt.fields.Log,
				Authenticator: tt.fields.Authenticator,
			}
			ctrl.setupProtectedRoutes(tt.args.router)
		})
	}
}
