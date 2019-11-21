// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"time"

	"github.com/amreo/ercole-services/api-service/service"
	"github.com/amreo/ercole-services/config"
	"github.com/goji/httpauth"
)

// APIControllerInterface is a interface that wrap methods used to querying data
type APIControllerInterface interface {
	// AuthenticateMiddleware return the middleware used to authenticate users
	AuthenticateMiddleware() func(http.Handler) http.Handler

	// GetCurrentHosts return all current hosts data using the filters in the request
	GetCurrentHosts(w http.ResponseWriter, r *http.Request)
}

// APIController is the struct used to handle the requests from agents and contains the concrete implementation of APIControllerInterface
type APIController struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Service contains the underlying service used to perform various logical and store operations
	Service service.APIServiceInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
}

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (ctrl *APIController) AuthenticateMiddleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(ctrl.Config.APIService.UserUsername, ctrl.Config.APIService.UserPassword)
}
