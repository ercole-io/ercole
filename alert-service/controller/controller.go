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

// Package controller contains structs and methods used to provide endpoints for storing hostdata informations
package controller

import (
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/alert-service/service"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/goji/httpauth"
)

// AlertQueueControllerInterface is a interface that wrap methods used to inserting events in the queue
type AlertQueueControllerInterface interface {
	ThrowNewAlert(w http.ResponseWriter, r *http.Request)
	// AuthenticateMiddleware return the middleware used to authenticate users
	AuthenticateMiddleware() func(http.Handler) http.Handler
}

// AlertQueueController is the struct used to handle the requests from agents and contains the concrete implementation of AlertQueueControllerInterface
type AlertQueueController struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Service contains the underlying service used to perform various logical and store operations
	Service service.AlertServiceInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log logger.Logger
}

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (ctrl *AlertQueueController) AuthenticateMiddleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(ctrl.Config.AlertService.PublisherUsername, ctrl.Config.AlertService.PublisherPassword)
}
