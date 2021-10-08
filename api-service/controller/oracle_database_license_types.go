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
	"errors"
	"net/http"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

// GetOracleDatabaseLicenseTypes return the list of OracleDatabaseLicenseTypes
func (ctrl *APIController) GetOracleDatabaseLicenseTypes(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetOracleDatabaseLicenseTypes()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

// GetOracleDatabaseLicensesCompliance return list of licenses with usage and compliance
func (ctrl *APIController) GetOracleDatabaseLicensesCompliance(w http.ResponseWriter, r *http.Request) {
	licenses, err := ctrl.Service.GetOracleDatabaseLicensesCompliance()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, licenses)
}

// DeleteOracleDatabaseLicenseTypes remove a licence type - Oracle/Database agreement part
func (ctrl *APIController) DeleteOracleDatabaseLicenseTypes(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var err error

	id := mux.Vars(r)["id"]

	if err = ctrl.Service.DeleteOracleDatabaseLicenseTypes(id); errors.Is(err, utils.ErrOracleDatabaseLicenseTypeNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
