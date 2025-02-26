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

func (ctrl *APIController) GetOracleDiskGroups(w http.ResponseWriter, r *http.Request) {
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]

	res, err := ctrl.Service.GetOracleDiskGroups(hostname, dbname)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, res)
}

func (ctrl *APIController) ListOracleDiskGroups(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	switch choice {
	case "application/json":
		result, err := ctrl.listOracleDiskGroupsJSON(filter)
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, result)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		result, err := ctrl.listOracleDiskGroupsXLSX(filter)
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
			return
		}

		utils.WriteXLSXResponse(w, result)
	}
}

func (ctrl *APIController) listOracleDiskGroupsJSON(filters *dto.GlobalFilter) ([]dto.OracleDatabaseDiskGroupDto, error) {
	return ctrl.Service.ListOracleDiskGroups(*filters)
}

func (ctrl *APIController) listOracleDiskGroupsXLSX(filters *dto.GlobalFilter) (*excelize.File, error) {
	return ctrl.Service.CreateOracleDiskGroupsXLSX(*filters)
}
