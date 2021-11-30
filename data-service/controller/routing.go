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

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
)

// GetDataControllerHandler setup the routes of the router using the handler in the controller as http handler
func (ctrl *DataController) GetDataControllerHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Pong"))
	})

	router.StrictSlash(true)
	router.Use(ctrl.AuthenticateMiddleware)

	ctrl.setupProtectedRoutes(router)

	return router
}

func (ctrl *DataController) setupProtectedRoutes(router *mux.Router) {
	router.HandleFunc("/hosts", ctrl.InsertHostData).Methods("POST")
	router.HandleFunc("/cmdbs", ctrl.CompareCmdbInfo).Methods("POST")
}

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (ctrl *DataController) AuthenticateMiddleware(h http.Handler) http.Handler {
	basicAuthHandler := httpauth.SimpleBasicAuth(ctrl.Config.DataService.AgentUsername, ctrl.Config.DataService.AgentPassword)
	return basicAuthHandler(h)
}
