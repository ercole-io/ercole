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

	"github.com/amreo/ercole-services/utils"
)

// GetDatabaseEnvironmentStats return all statistics about the environments of the databases using the filters in the request
func (ctrl *APIController) GetDatabaseEnvironmentStats(w http.ResponseWriter, r *http.Request) {
	var location string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")

	//get the data
	stats, err := ctrl.Service.GetDatabaseEnvironmentStats(location)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetDatabaseVersionStats return all statistics about the versions of the databases using the filters in the request
func (ctrl *APIController) GetDatabaseVersionStats(w http.ResponseWriter, r *http.Request) {
	var location string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")

	//get the data
	stats, err := ctrl.Service.GetDatabaseVersionStats(location)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTopReclaimableDatabaseStats return top databases by reclaimable segment advisors using the filters in the request
func (ctrl *APIController) GetTopReclaimableDatabaseStats(w http.ResponseWriter, r *http.Request) {
	var location string
	var limit int
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	if limit, err = utils.Str2int(r.URL.Query().Get("limit"), 15); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTopReclaimableDatabaseStats(location, limit)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetDatabasePatchStatusStats return all statistics about the patch status of the databases using the filters in the request
func (ctrl *APIController) GetDatabasePatchStatusStats(w http.ResponseWriter, r *http.Request) {
	var location string
	var windowTime int
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetDatabasePatchStatusStats(location, time.Now().AddDate(0, -windowTime, 0))
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTopWorkloadDatabaseStats return top databases by workload advisors using the filters in the request
func (ctrl *APIController) GetTopWorkloadDatabaseStats(w http.ResponseWriter, r *http.Request) {
	var location string
	var limit int
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	if limit, err = utils.Str2int(r.URL.Query().Get("limit"), 10); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTopWorkloadDatabaseStats(location, limit)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetDatabaseDataguardStatusStats return all statistics about the dataguard status of the databases using the filters in the request
func (ctrl *APIController) GetDatabaseDataguardStatusStats(w http.ResponseWriter, r *http.Request) {
	var location string
	var environment string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	//get the data
	stats, err := ctrl.Service.GetDatabaseDataguardStatusStats(location, environment)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
