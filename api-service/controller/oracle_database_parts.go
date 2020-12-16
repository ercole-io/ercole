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

	"github.com/ercole-io/ercole/v2/utils"
)

//TODO Rename parts in licenses

// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
func (ctrl *APIController) GetOracleDatabaseAgreementPartsList(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetOracleDatabaseAgreementPartsList()
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
