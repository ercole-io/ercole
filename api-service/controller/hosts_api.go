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
	"github.com/gorilla/mux"
)

// SearchHosts search hosts data using the filters in the request
func (ctrl *APIController) SearchHosts(w http.ResponseWriter, r *http.Request) {
	var full bool
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var location string
	var environment string
	var olderThan time.Time

	var err utils.AdvancedErrorInterface
	//parse the query params
	if full, err = utils.Str2bool(r.URL.Query().Get("full"), false); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	search = r.URL.Query().Get("search")
	sortBy = r.URL.Query().Get("sort-by")
	if sortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if pageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	hosts, err := ctrl.Service.SearchHosts(full, search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts[0])
	}
}

// GetHost return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHost(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var err utils.AdvancedErrorInterface

	hostname := mux.Vars(r)["hostname"]

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	host, err := ctrl.Service.GetHost(hostname, olderThan)
	if err == utils.AerrHostNotFound {
		utils.WriteAndLogError(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, host)
}
