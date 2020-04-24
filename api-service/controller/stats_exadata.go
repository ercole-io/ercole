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

// GetTotalExadataMemorySizeStats return the total size of memory of exadata using the filters in the request
func (ctrl *APIController) GetTotalExadataMemorySizeStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var location string
	var environment string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTotalExadataMemorySizeStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTotalExadataCPUStats return the total cpu of exadata using the filters in the request
func (ctrl *APIController) GetTotalExadataCPUStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var location string
	var environment string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTotalExadataCPUStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetAverageExadataStorageUsageStats return the average usage of cell disks of exadata using the filters in the request
func (ctrl *APIController) GetAverageExadataStorageUsageStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var location string
	var environment string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetAverageExadataStorageUsageStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetExadataStorageErrorCountStatusStats return all statistics about the ErrorCount status of the storage of the exadata using the filters in the request
func (ctrl *APIController) GetExadataStorageErrorCountStatusStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var location string
	var environment string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetExadataStorageErrorCountStatusStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetExadataPatchStatusStats return all statistics about the patch status of the exadata using the filters in the request
func (ctrl *APIController) GetExadataPatchStatusStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var location string
	var environment string
	var windowTime int

	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")
	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetExadataPatchStatusStats(location, environment, ctrl.TimeNow().AddDate(0, -windowTime, 0), olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
