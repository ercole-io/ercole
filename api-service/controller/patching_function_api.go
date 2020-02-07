// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/amreo/ercole-services/utils"
	"github.com/gorilla/mux"
)

// GetPatchingFunction return all'informations about the patching function of the host requested in the hostnmae path variable
func (ctrl *APIController) GetPatchingFunction(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface

	hostname := mux.Vars(r)["hostname"]

	//get the data
	pf, err := ctrl.Service.GetPatchingFunction(hostname)
	if err == utils.AerrHostNotFound || err == utils.AerrPatchingFunctionNotFound {
		utils.WriteAndLogError(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, pf)
}

// AddTagToDatabase add a tag to the database if it hasn't the tag
func (ctrl *APIController) AddTagToDatabase(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, "BAD_REQUEST"))
		return
	}

	//set the value
	aerr := ctrl.Service.AddTagToDatabase(hostname, dbname, string(raw))
	if aerr != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
