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

	"github.com/gorilla/mux"
)

// SetupRoutesForAPIController setup the routes of the router using the handler in the controller as http handler
func SetupRoutesForAPIController(router *mux.Router, ctrl APIControllerInterface) {
	//Enable authentication using the ctrl
	router.Use(ctrl.AuthenticateMiddleware())

	//Add the routes
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	router.HandleFunc("/hosts", ctrl.SearchCurrentHosts).Methods("GET")
	router.HandleFunc("/clusters", ctrl.SearchCurrentClusters).Methods("GET")
	router.HandleFunc("/addms", ctrl.SearchCurrentAddms).Methods("GET")
	router.HandleFunc("/segment-advisors", ctrl.SearchCurrentSegmentAdvisors).Methods("GET")
	router.HandleFunc("/patch-advisors", ctrl.SearchCurrentPatchAdvisors).Methods("GET")
	router.HandleFunc("/databases", ctrl.SearchCurrentDatabases).Methods("GET")
	router.HandleFunc("/hosts/{hostname}", ctrl.GetCurrentHost).Methods("GET")
	router.HandleFunc("/alerts", ctrl.SearchAlerts).Methods("GET")
	router.HandleFunc("/stats/environments", ctrl.GetEnvironmentStats).Methods("GET")
}
