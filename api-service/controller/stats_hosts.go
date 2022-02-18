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
	"time"

	"github.com/ercole-io/ercole/v2/utils"
)

// GetHostsCountStats return the number of the hosts using the filters in the request
func (ctrl *APIController) GetHostsCountStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location, environment string

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetHostsCountStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetEnvironmentStats return all statistics about the environments of the hosts using the filters in the request
func (ctrl *APIController) GetEnvironmentStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location string

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetEnvironmentStats(location, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTypeStats return all statistics about the types of the hosts using the filters in the request
func (ctrl *APIController) GetTypeStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location string

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTypeStats(location, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetOperatingSystemStats return all statistics about the operating systems of the hosts using the filters in the request
func (ctrl *APIController) GetOperatingSystemStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location string

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetOperatingSystemStats(location, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
