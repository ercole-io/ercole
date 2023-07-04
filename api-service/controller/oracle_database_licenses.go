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
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/model"
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

	ignored, err := strconv.ParseBool(mux.Vars(r)["ignored"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(err, "BAD_REQUEST"))
		return
	}

	req := model.OracleDatabaseLicense{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	//set the value
	err = ctrl.Service.UpdateLicenseIgnoredField(hostname, dbname, licensetypeid, ignored, req.IgnoredComment)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

func (ctrl *APIController) CanMigrateLicense(w http.ResponseWriter, r *http.Request) {
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]

	res, err := ctrl.Service.CanMigrateLicense(hostname, dbname)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	resJSON := struct {
		Canbemigrate bool
	}{
		Canbemigrate: res,
	}

	utils.WriteJSONResponse(w, http.StatusOK, resJSON)
}
