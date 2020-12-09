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
)

// SetupRoutesForAlertQueueController setup the routes of the router using the handler in the controller as http handler
func SetupRoutesForAlertQueueController(router *mux.Router, ctrl AlertQueueControllerInterface) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	router = router.NewRoute().Subrouter()
	router.Use(ctrl.AuthenticateMiddleware())

	setupProtectedRoutes(router, ctrl)
}

func setupProtectedRoutes(router *mux.Router, ctrl AlertQueueControllerInterface) {
	router.HandleFunc("/alerts", ctrl.ThrowNewAlert).Methods("POST")
	router.HandleFunc("/queue/host-data-insertion/{id}", ctrl.HostDataInsertion).Methods("POST")
}
