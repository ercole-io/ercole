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
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/context"
)

// GetOracleDatabaseEnvironmentStats return all statistics about the environments of the databases using the filters in the request
func (ctrl *APIController) GetOracleDatabaseEnvironmentStats(w http.ResponseWriter, r *http.Request) {
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
	stats, err := ctrl.Service.GetOracleDatabaseEnvironmentStats(location, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetOracleDatabaseHighReliabilityStats return all statistics about the high-reliability status of the databases using the filters in the request
func (ctrl *APIController) GetOracleDatabaseHighReliabilityStats(w http.ResponseWriter, r *http.Request) {
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
	stats, err := ctrl.Service.GetOracleDatabaseHighReliabilityStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetOracleDatabaseVersionStats return all statistics about the versions of the databases using the filters in the request
func (ctrl *APIController) GetOracleDatabaseVersionStats(w http.ResponseWriter, r *http.Request) {
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
	stats, err := ctrl.Service.GetOracleDatabaseVersionStats(location, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTopReclaimableOracleDatabaseStats return top databases by reclaimable segment advisors using the filters in the request
func (ctrl *APIController) GetTopReclaimableOracleDatabaseStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location string

	var limit int

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")

	if limit, err = utils.Str2int(r.URL.Query().Get("limit"), 15); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTopReclaimableOracleDatabaseStats(location, limit, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetOracleDatabasePatchStatusStats return all statistics about the patch status of the databases using the filters in the request
func (ctrl *APIController) GetOracleDatabasePatchStatusStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location string

	var windowTime int

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")

	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetOracleDatabasePatchStatusStats(location, ctrl.TimeNow().AddDate(0, -windowTime, 0), olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTopWorkloadOracleDatabaseStats return top databases by workload advisors using the filters in the request
func (ctrl *APIController) GetTopWorkloadOracleDatabaseStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location string

	var limit int

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")
	if location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		location = strings.Join(locations, ",")
	}

	if limit, err = utils.Str2int(r.URL.Query().Get("limit"), 10); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTopWorkloadOracleDatabaseStats(location, limit, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetOracleDatabaseDataguardStatusStats return all statistics about the dataguard status of the databases using the filters in the request
func (ctrl *APIController) GetOracleDatabaseDataguardStatusStats(w http.ResponseWriter, r *http.Request) {
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
	stats, err := ctrl.Service.GetOracleDatabaseDataguardStatusStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetOracleDatabaseRACStatusStats return all statistics about the RAC status of the databases using the filters in the request
func (ctrl *APIController) GetOracleDatabaseRACStatusStats(w http.ResponseWriter, r *http.Request) {
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
	stats, err := ctrl.Service.GetOracleDatabaseRACStatusStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetOracleDatabaseArchivelogStatusStats return all statistics about the archivelog status of the databases using the filters in the request
func (ctrl *APIController) GetOracleDatabaseArchivelogStatusStats(w http.ResponseWriter, r *http.Request) {
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
	stats, err := ctrl.Service.GetOracleDatabaseArchivelogStatusStats(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

func (ctrl *APIController) GetOracleDatabasesStatistics(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if filter.Location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		filter.Location = strings.Join(locations, ",")
	}

	stats, err := ctrl.Service.GetOracleDatabasesStatistics(*filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

// GetTopUnusedOracleDatabaseInstanceResourceStats return top unused instance resource by databases work using the filters in the request
func (ctrl *APIController) GetTopUnusedOracleDatabaseInstanceResourceStats(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time

	var location, environment string

	var limit int

	var err error

	//parse the query params
	location = r.URL.Query().Get("location")

	if location == "" {
		user := context.Get(r, "user")
		locations, errLocation := ctrl.Service.ListLocations(user)

		if errLocation != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errLocation)
			return
		}

		location = strings.Join(locations, ",")
	}

	environment = r.URL.Query().Get("environment")

	if limit, err = utils.Str2int(r.URL.Query().Get("limit"), 15); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	stats, err := ctrl.Service.GetTopUnusedOracleDatabaseInstanceResourceStats(location, environment, limit, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
