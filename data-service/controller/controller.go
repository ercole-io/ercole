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

	"github.com/ercole-io/ercole/v2/data-service/service"
	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/v2/config"
)

type DataControllerInterface interface {
	InsertHostData(w http.ResponseWriter, r *http.Request)
	CompareCmdbInfo(w http.ResponseWriter, r *http.Request)

	AuthenticateMiddleware(h http.Handler) http.Handler
}

type DataController struct {
	Config  config.Configuration
	Service service.HostDataServiceInterface
	TimeNow func() time.Time
	Log     *logrus.Logger
}
