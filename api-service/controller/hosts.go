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
	"errors"
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
)

// SearchHosts search hosts data using the filters in the request
func (ctrl *APIController) SearchHosts(w http.ResponseWriter, r *http.Request) {
	filters, err := dto.GetSearchHostFilters(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	requestContentType := httputil.NegotiateContentType(r,
		[]string{
			"application/json",
			"application/vnd.oracle.lms+vnd.ms-excel.sheet.macroEnabled.12",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.ercole.mongohostdata+json",
		},
		"application/json")

	switch requestContentType {
	case "application/json":
		ctrl.searchHostsJSON(w, r, filters)
	case "application/vnd.oracle.lms+vnd.ms-excel.sheet.macroEnabled.12":
		ctrl.searchHostsLMS(w, r, filters)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.searchHostsXLSX(w, r, filters)
	}
}

// searchHostsJSON search hosts data using the filters in the request returning it in JSON
func (ctrl *APIController) searchHostsJSON(w http.ResponseWriter, r *http.Request, filters *dto.SearchHostsFilters) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "full"
	}

	if mode != "full" && mode != "summary" && mode != "lms" && mode != "mhd" && mode != "hostnames" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errors.New("Invalid mode value"))
		return
	}

	if mode == "summary" {
		ctrl.getHostDataSummaries(w, r, filters)
		return
	}

	hosts, err := ctrl.Service.SearchHosts(mode, *filters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if mode == "hostnames" {
		hostnames := make([]string, len(hosts))
		for i, h := range hosts {
			hostnames[i] = h["hostname"].(string)
		}
		utils.WriteJSONResponse(w, http.StatusOK, hostnames)
	} else {
		if filters.PageNumber == -1 || filters.PageSize == -1 {
			utils.WriteJSONResponse(w, http.StatusOK, hosts)
		} else {
			utils.WriteJSONResponse(w, http.StatusOK, hosts[0])
		}
	}
}

func (ctrl *APIController) getHostDataSummaries(w http.ResponseWriter, r *http.Request, filters *dto.SearchHostsFilters) {
	hosts, err := ctrl.Service.GetHostDataSummaries(*filters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	resp := map[string]interface{}{
		"hosts": hosts,
	}
	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

// searchHostsLMS search hosts data using the filters in the request returning it in LMS+XLSX
func (ctrl *APIController) searchHostsLMS(w http.ResponseWriter, r *http.Request, filters *dto.SearchHostsFilters) {
	filters.PageNumber, filters.PageSize = -1, -1
	lms, err := ctrl.Service.SearchHostsAsLMS(*filters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSMResponse(w, lms)
}

// searchHostsXLSX search hosts data using the filters in the request returning it in XLSX
func (ctrl *APIController) searchHostsXLSX(w http.ResponseWriter, r *http.Request, filters *dto.SearchHostsFilters) {
	filters.PageNumber, filters.PageSize = -1, -1
	xlsx, err := ctrl.Service.SearchHostsAsXLSX(*filters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

// GetHost return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHost(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.ercole.mongohostdata+json"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.GetHostJSON(w, r)
	case "application/vnd.ercole.mongohostdata+json":
		ctrl.GetHostMongoJSON(w, r)
	}
}

// GetHostJSON return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHostJSON(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var err error

	hostname := mux.Vars(r)["hostname"]

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	host, err := ctrl.Service.GetHost(hostname, olderThan, false)
	if errors.Is(err, utils.ErrHostNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, host)
}

// GetHostMongoJSON return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHostMongoJSON(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var err error

	hostname := mux.Vars(r)["hostname"]

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	host, err := ctrl.Service.GetHost(hostname, olderThan, true)
	if errors.Is(err, utils.ErrHostNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteExtJSONResponse(ctrl.Log, w, http.StatusOK, host)
}

// ListLocations list locations using the filters in the request
func (ctrl *APIController) ListLocations(w http.ResponseWriter, r *http.Request) {
	var location string
	var environment string
	var olderThan time.Time

	var err error
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	locations, err := ctrl.Service.ListLocations(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, locations)
}

// ListEnvironments list the environments using the filters in the request
func (ctrl *APIController) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	var location string
	var environment string
	var olderThan time.Time

	var err error
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	environments, err := ctrl.Service.ListEnvironments(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, environments)
}

// ArchiveHost archive the specified host in the request
func (ctrl *APIController) ArchiveHost(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	hostname := mux.Vars(r)["hostname"]

	err := ctrl.Service.ArchiveHost(hostname)
	if errors.Is(err, utils.ErrHostNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
