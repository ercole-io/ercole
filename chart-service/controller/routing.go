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

	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/api-service/auth"
)

// GetChartControllerHandler setup the routes of the router using the handler in the controller as http handler
func (ctrl *ChartController) GetChartControllerHandler(auth auth.AuthenticationProvider) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Pong")); err != nil {
			ctrl.Log.Error(err)
			return
		}
	})

	router.HandleFunc("/user/login", auth.GetToken).Methods("POST")

	subrouter := router.NewRoute().Subrouter()
	subrouter.Use(auth.AuthenticateMiddleware)
	ctrl.setupProtectedRoutes(subrouter)

	return router
}

func (ctrl *ChartController) setupProtectedRoutes(router *mux.Router) {
	router.HandleFunc("/settings/technologies-metrics", ctrl.GetTechnologiesMetrics).Methods("GET")

	router.HandleFunc("/technologies/all/license-history", ctrl.GetLicenseComplianceHistory).Methods("GET")
	router.HandleFunc("/technologies/oracle/database", ctrl.GetOracleDatabaseChart).Methods("GET")

	router.HandleFunc("/technologies/changes", ctrl.GetChangeChart).Methods("GET")
	router.HandleFunc("/technologies/types", ctrl.GetTechnologyTypes).Methods("GET")

	router.HandleFunc("/hosts/cores", ctrl.GetHostCores).Methods("GET")
}
