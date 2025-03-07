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
	"errors"
	"net/http"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) GetMissingDatabases(w http.ResponseWriter, r *http.Request) {
	res, err := ctrl.Service.GetMissingDatabases()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, res)
}

func (ctrl *APIController) GetMissingDatabasesByHostname(w http.ResponseWriter, r *http.Request) {
	hostname := mux.Vars(r)["hostname"]

	host, err := ctrl.Service.GetHost(hostname, utils.MAX_TIME, true)
	if errors.Is(err, utils.ErrHostNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if !ctrl.userHasAccessToLocation(r, host.Location) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.ErrPermissionDenied)
		return
	}

	res, err := ctrl.Service.GetMissingDatabasesByHostname(hostname)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, res)
}
