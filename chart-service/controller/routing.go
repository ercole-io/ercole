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

	"github.com/ercole-io/ercole/api-service/auth"
	"github.com/gorilla/mux"
)

// SetupRoutesForChartController setup the routes of the router using the handler in the controller as http handler
func SetupRoutesForChartController(router *mux.Router, ctrl ChartControllerInterface, auth auth.AuthenticationProvider) {

	//Add the routes
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	router.HandleFunc("/user/login", auth.GetToken).Methods("POST")
	//Enable authentication using the ctrl
	router = router.NewRoute().Subrouter()
	router.Use(auth.AuthenticateMiddleware)
	setupProtectedRoutes(router, ctrl)
}

func setupProtectedRoutes(router *mux.Router, ctrl ChartControllerInterface) {
	router.HandleFunc("/settings/technologiy-metrics", ctrl.GetTechnologyList).Methods("GET")
	router.HandleFunc("/technologies/oracle/database", ctrl.GetOracleDatabaseChart).Methods("GET")
	router.HandleFunc("/technologies/changes", ctrl.GetChangeChart).Methods("GET")
	router.HandleFunc("/technologies/types", ctrl.GetTechnologyTypes).Methods("GET")
}
