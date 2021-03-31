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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
)

func (ctrl *APIController) SearchDatabases(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r,
		[]string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		"application/json")

	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	switch choiche {
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

func (ctrl *APIController) GetDatabasesUsedLicenses(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	usedLicenses, err := ctrl.Service.GetDatabasesUsedLicenses(*filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"usedLicenses": usedLicenses,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetDatabaseLicensesCompliance(w http.ResponseWriter, r *http.Request) {
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
