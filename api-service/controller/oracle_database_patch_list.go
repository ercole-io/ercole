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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/context"
)

func (ctrl *APIController) GetOraclePatchList(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
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

	switch choice {
	case "application/json":
		result, err := ctrl.GetOraclePatchListJSON(filter)
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, result)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		result, err := ctrl.GetOraclePatchListXLSX(filter)
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
			return
		}

		utils.WriteXLSXResponse(w, result)
	}
}

func (ctrl *APIController) GetOraclePatchListJSON(filters *dto.GlobalFilter) ([]dto.OracleDatabasePatchDto, error) {
	return ctrl.Service.GetOraclePatchList(*filters)
}

func (ctrl *APIController) GetOraclePatchListXLSX(filters *dto.GlobalFilter) (*excelize.File, error) {
	return ctrl.Service.CreateGetOraclePatchListXLSX(*filters)
}
