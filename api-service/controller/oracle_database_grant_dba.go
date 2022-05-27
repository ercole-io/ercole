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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) ListOracleGrantDbaByHostname(w http.ResponseWriter, r *http.Request) {
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
		result, err := ctrl.GetOracleGrantDbaJSON(hostname, filter)
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, result)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		result, err := ctrl.GetOracleGrantDbaXLSX(hostname, filter)
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteXLSXResponse(w, result)
	}
}

func (ctrl *APIController) GetOracleGrantDbaJSON(hostname string, filters *dto.GlobalFilter) ([]dto.OracleGrantDbaDto, error) {
	return ctrl.Service.ListOracleGrantDbaByHostname(hostname, *filters)
}

func (ctrl *APIController) GetOracleGrantDbaXLSX(hostname string, filters *dto.GlobalFilter) (*excelize.File, error) {
	return ctrl.Service.CreateOracleGrantDbaXlsx(hostname, *filters)
}
