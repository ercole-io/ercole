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

	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *APIController) SearchDatabases(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r,
		[]string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		"application/json")

	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	switch choice {
	case "application/json":
		ctrl.SearchDatabasesJSON(w, r, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchDatabasesXLSX(w, r, *filter)
	}
}

func (ctrl *APIController) SearchDatabasesJSON(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	databases, err := ctrl.Service.SearchDatabases(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"databases": databases,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) SearchDatabasesXLSX(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	file, err := ctrl.Service.SearchDatabasesAsXLSX(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, file)
}

func (ctrl *APIController) GetDatabasesStatistics(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	stats, err := ctrl.Service.GetDatabasesStatistics(*filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

func (ctrl *APIController) GetUsedLicensesPerDatabases(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	hostname := ""

	vars := mux.Vars(r)
	if v, ok := vars["hostname"]; ok {
		hostname = v
	}

	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.GetUsedLicensesPerDatabasesJSON(w, r, hostname, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.GetUsedLicensesPerDatabasesAsXLSX(w, r, *filter)
	}
}

func (ctrl *APIController) GetUsedLicensesPerDatabasesJSON(w http.ResponseWriter, r *http.Request, hostname string, filter dto.GlobalFilter) {
	usedLicenses, err := ctrl.Service.GetUsedLicensesPerDatabases(hostname, filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"usedLicenses": usedLicenses,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetUsedLicensesPerDatabasesAsXLSX(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	xlsx, err := ctrl.Service.GetUsedLicensesPerDatabasesAsXLSX(filter)

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

func (ctrl *APIController) GetDatabaseLicensesCompliance(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.GetDatabaseLicensesComplianceJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.GetDatabaseLicensesComplianceXLSX(w, r)
	}
}

func (ctrl *APIController) GetDatabaseLicensesComplianceJSON(w http.ResponseWriter, r *http.Request) {
	licenses, err := ctrl.Service.GetDatabaseLicensesCompliance()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"licensesCompliance": licenses,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetDatabaseLicensesComplianceXLSX(w http.ResponseWriter, r *http.Request) {
	xlsx, err := ctrl.Service.GetDatabaseLicensesComplianceAsXLSX()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

func (ctrl *APIController) GetUsedLicensesPerHost(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.GetUsedLicensesPerHostJSON(w, r, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.GetUsedLicensesPerHostAsXLSX(w, r, *filter)
	}
}

func (ctrl *APIController) GetUsedLicensesPerHostJSON(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	usedLicenses, err := ctrl.Service.GetUsedLicensesPerHost(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"usedLicenses": usedLicenses,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetUsedLicensesPerHostAsXLSX(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	xlsx, err := ctrl.Service.GetUsedLicensesPerHostAsXLSX(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

func (ctrl *APIController) GetUsedLicensesPerCluster(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	choice := httputil.NegotiateContentType(r, []string{"application/json",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.GetUsedLicensesPerClusterJSON(w, r, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.GetUsedLicensesPerClusterAsXLSX(w, r, *filter)
	}
}

func (ctrl *APIController) GetUsedLicensesPerClusterJSON(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	usedLicenses, err := ctrl.Service.GetUsedLicensesPerCluster(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"usedLicensesPerCluster": usedLicenses,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetUsedLicensesPerClusterAsXLSX(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	xlsx, err := ctrl.Service.GetUsedLicensesPerClusterAsXLSX(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}
