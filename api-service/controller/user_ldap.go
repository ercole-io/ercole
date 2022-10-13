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

func (ctrl *APIController) GetLDAPUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	user := mux.Vars(r)["user"]

	userLDAP, err := ctrl.Service.GetLDAPUsers(user)
	if errors.Is(err, utils.ErrConnectLDAPServer) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.ErrConnectLDAPServer)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, userLDAP)
}
