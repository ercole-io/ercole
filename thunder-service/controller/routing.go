// Copyright (c) 2022 Sorint.lab S.p.A.
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

// GetThunderControllerHandler setup the routes of the router using the handler in the controller as http handler
func (ctrl *ThunderController) GetThunderControllerHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Pong")); err != nil {
			ctrl.Log.Error(err)
			return
		}
	})

	subrouter := router.NewRoute().Subrouter()
	ctrl.setupProtectedRoutes(subrouter)

	return router
}

func (ctrl *ThunderController) setupProtectedRoutes(router *mux.Router) {
	router.HandleFunc("/oracle-cloud/recommendations/{ids}", ctrl.GetOciNativeRecommendations).Methods("GET")
	router.HandleFunc("/oracle-cloud/configurations", ctrl.GetOciProfiles).Methods("GET")
	router.HandleFunc("/oracle-cloud/configurations", ctrl.AddOciProfile).Methods("POST")
	router.HandleFunc("/oracle-cloud/configurations/{id}", ctrl.UpdateOciProfile).Methods("PUT")
	router.HandleFunc("/oracle-cloud/configurations/{id}", ctrl.DeleteOciProfile).Methods("DELETE")
	router.HandleFunc("/oracle-cloud/oci-objects", ctrl.GetOciObjects).Methods("GET")
	router.HandleFunc("/oracle-cloud/oci-recommendations", ctrl.GetOciRecommendations).Methods("GET")
	router.HandleFunc("/oracle-cloud/oci-recommendation-errors/{seqnum}", ctrl.GetOciRecommendationErrors).Methods("GET")
	router.HandleFunc("/oracle-cloud/retrieve-last-oci-recommendations", ctrl.ForceGetOciRecommendations).Methods("GET")
	router.HandleFunc("/oracle-cloud/profile-selection/profileid/{profileid}/selected/{selected}", ctrl.SelectOciProfile).Methods("PUT")
	router.HandleFunc("/aws/configurations", ctrl.GetAwsProfiles).Methods("GET")
	router.HandleFunc("/aws/configurations", ctrl.AddAwsProfile).Methods("POST")
	router.HandleFunc("/aws/configurations/{id}", ctrl.UpdateAwsProfile).Methods("PUT")
	router.HandleFunc("/aws/configurations/{id}", ctrl.DeleteAwsProfile).Methods("DELETE")
	router.HandleFunc("/aws/profile-selection/profileid/{profileid}/selected/{selected}", ctrl.SelectAwsProfile).Methods("PUT")
	router.HandleFunc("/aws/aws-recommendations", ctrl.GetAwsRecommendations).Methods("GET")
	router.HandleFunc("/aws/retrieve-last-aws-recommendations", ctrl.ForceGetAwsRecommendations).Methods("GET")
}
