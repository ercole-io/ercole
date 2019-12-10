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
	// SearchCurrentHosts search current hosts data using the filters in the request
	SearchCurrentHosts(w http.ResponseWriter, r *http.Request)
	// SearchCurrentDatabases search current databases data using the filters in the request
	SearchCurrentDatabases(w http.ResponseWriter, r *http.Request)
	// SearchCurrentClusters search current clusters data using the filters in the request
	SearchCurrentClusters(w http.ResponseWriter, r *http.Request)
	// SearchCurrentAddms search current addms data using the filters in the request
	SearchCurrentAddms(w http.ResponseWriter, r *http.Request)
	// SearchCurrentSegmentAdvisors search current segment advisors data using the filters in the request
	SearchCurrentSegmentAdvisors(w http.ResponseWriter, r *http.Request)
	// SearchCurrentPatchAdvisors search current patch advisors data using the filters in the request
	SearchCurrentPatchAdvisors(w http.ResponseWriter, r *http.Request)
	// GetCurrentHost return all'informations about the current host requested in the id path variable
	GetCurrentHost(w http.ResponseWriter, r *http.Request)
	// SearchAlerts search alerts using the filters in the request
	SearchAlerts(w http.ResponseWriter, r *http.Request)

	// GetEnvironmentStats return all statistics about the environments of the hosts using the filters in the request
	GetEnvironmentStats(w http.ResponseWriter, r *http.Request)
	// GetTypeStats return all statistics about the types of the hosts using the filters in the request
	GetTypeStats(w http.ResponseWriter, r *http.Request)
	// GetOperatingSystemStats return all statistics about the operating systems of the hosts using the filters in the request
	GetOperatingSystemStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseEnvironmentStats return all statistics about the environments of the databases using the filters in the request
	GetDatabaseEnvironmentStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseVersionStats return all statistics about the versions of the databases using the filters in the request
	GetDatabaseVersionStats(w http.ResponseWriter, r *http.Request)
	// GetTopReclaimableDatabaseStats return all the top database by reclaimable segment advisors using the filters in the request
	GetTopReclaimableDatabaseStats(w http.ResponseWriter, r *http.Request)
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
