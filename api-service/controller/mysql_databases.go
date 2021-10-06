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

	"github.com/golang/gddo/httputil"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *APIController) SearchMySQLInstances(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r,
		[]string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		"application/json")

	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	switch choiche {
	case "application/json":
		ctrl.SearchMySQLInstancesJSON(w, r, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchMySQLInstancesXLSX(w, r, *filter)
	}
}

func (ctrl *APIController) SearchMySQLInstancesJSON(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	databases, err := ctrl.Service.SearchMySQLInstances(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"databases": databases,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) SearchMySQLInstancesXLSX(w http.ResponseWriter, r *http.Request, filter dto.GlobalFilter) {
	file, err := ctrl.Service.SearchMySQLInstancesAsXLSX(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, file)
}
