// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/utils"
)

// UpdateLicenseIgnoredField update license ignored field (true/false)
func (ctrl *APIController) UpdateLicenseIgnoredField(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]
	licensetypeid := mux.Vars(r)["licenseTypeID"]
	ignored := mux.Vars(r)["ignored"]

	isIgnored, err := strconv.ParseBool(ignored)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(err, "BAD_REQUEST"))
		return
	}

	//set the value
	err = ctrl.Service.UpdateLicenseIgnoredField(hostname, dbname, licensetypeid, isIgnored)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
