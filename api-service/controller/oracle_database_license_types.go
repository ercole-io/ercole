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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// GetOracleDatabaseLicensesCompliance return list of licenses with usage and compliance
func (ctrl *APIController) GetOracleDatabaseLicensesCompliance(w http.ResponseWriter, r *http.Request) {
	licenses, err := ctrl.Service.GetOracleDatabaseLicensesCompliance()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, licenses)
}

// GetOracleDatabaseLicenseTypes return the list of OracleDatabaseLicenseTypes
func (ctrl *APIController) GetOracleDatabaseLicenseTypes(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetOracleDatabaseLicenseTypes()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"license-types": data,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// DeleteOracleDatabaseLicenseType remove a licence type - Oracle/Database agreement part
func (ctrl *APIController) DeleteOracleDatabaseLicenseType(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	id := mux.Vars(r)["id"]

	err := ctrl.Service.DeleteOracleDatabaseLicenseType(id)

	if errors.Is(err, utils.ErrOracleDatabaseLicenseTypeIDNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// AddOracleDatabaseLicenseType add a licence type - Oracle/Database agreement part to the database if it hasn't a licence type
func (ctrl *APIController) AddOracleDatabaseLicenseType(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req model.OracleDatabaseLicenseType

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	lt, err := ctrl.Service.AddOracleDatabaseLicenseType(req)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, lt)
}

func (ctrl *APIController) UpdateOracleDatabaseLicenseType(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req model.OracleDatabaseLicenseType

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	lt, err := ctrl.Service.UpdateOracleDatabaseLicenseType(req)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, lt)
}

// UpdateLicenseIgnoredField update license ignored field (true/false)
func (ctrl *APIController) UpdateLicenseIgnoredField(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]
	licensetypeid := mux.Vars(r)["licenseTypeID"]

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewError(err, "BAD_REQUEST"))
		return
	}

	ignored, err := strconv.ParseBool(string(raw))
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(err, "BAD_REQUEST"))
		return
	}

	//set the value
	err = ctrl.Service.UpdateLicenseIgnoredField(hostname, dbname, licensetypeid, ignored)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
