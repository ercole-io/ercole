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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/context"
)

// SearchMongoDBInstances search instances data using the filters in the request
func (ctrl *APIController) SearchMongoDBInstances(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	filter, err := dto.GetMongoDBInstancesFilter(r)
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

	switch choice {
	case "application/json":
		ctrl.SearchMongoDBInstancesJSON(w, r, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchMongoDBInstancesXLSX(w, r, *filter)
	}
}

//SearchMongoDBInstancesJSON search instances data using the filters in the request returning it in JSON
func (ctrl *APIController) SearchMongoDBInstancesJSON(w http.ResponseWriter, r *http.Request, filter dto.SearchMongoDBInstancesFilter) {
	instances, err := ctrl.Service.SearchMongoDBInstances(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if filter.PageNumber == -1 || filter.PageSize == -1 {
		utils.WriteJSONResponse(w, http.StatusOK, instances.Content)
	} else {
		utils.WriteJSONResponse(w, http.StatusOK, instances)
	}
}

// SearchMongoDBInstancesXLSX search instances data using the filters in the request returning it in XLSX
func (ctrl *APIController) SearchMongoDBInstancesXLSX(w http.ResponseWriter, r *http.Request, filter dto.SearchMongoDBInstancesFilter) {
	file, err := ctrl.Service.SearchMongoDBInstancesAsXLSX(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, file)
}
