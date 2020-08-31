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
	"io/ioutil"
	"net/http"

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/utils"
	"gopkg.in/square/go-jose.v2/json"
)

// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
func (ctrl *APIController) GetOracleDatabaseAgreementPartsList(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetOracleDatabaseAgreementPartsList()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

// AddOracleDatabaseAgreements add some agreements
func (ctrl *APIController) AddOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var aerr utils.AdvancedErrorInterface
	var req apimodel.OracleDatabaseAgreementsAddRequest

	//Read all bytes for the request
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	//Unmarshal it to req
	if err := json.Unmarshal(raw, &req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Add it!
	var ids interface{}
	if ids, aerr = ctrl.Service.AddOracleDatabaseAgreements(req); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the created id
	utils.WriteJSONResponse(w, http.StatusOK, ids)
}
