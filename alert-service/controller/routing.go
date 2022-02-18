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

// GetAlertControllerHandler setup the routes of the router using the handler in the controller as http handler
func (ctrl *AlertQueueController) GetAlertControllerHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Pong")); err != nil {
			ctrl.Log.Error(err)
			return
		}
	})

	subrouter := router.NewRoute().Subrouter()
	subrouter.Use(ctrl.AuthenticateMiddleware())
	ctrl.setupProtectedRoutes(subrouter)

	return router
}

func (ctrl *AlertQueueController) setupProtectedRoutes(router *mux.Router) {
	router.HandleFunc("/alerts", ctrl.ThrowNewAlert).Methods("POST")
}
