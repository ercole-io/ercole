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

// Package controller contains structs and methods used to provide endpoints for storing hostdata informations
package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *AlertQueueController) ThrowNewAlert(w http.ResponseWriter, r *http.Request) {
	var alert model.Alert
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&alert); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	if !alert.IsValid() {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(errors.New("Invalid alert"), "INVALID_ALERT"))
		return
	}

	aerr := ctrl.Service.ThrowNewAlert(alert)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteNoContentResponse(w)
}
