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
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"

	"github.com/goji/httpauth"
	"github.com/xeipuuv/gojsonschema"
)

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (ctrl *HostDataController) AuthenticateMiddleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(ctrl.Config.DataService.AgentUsername, ctrl.Config.DataService.AgentPassword)
}

// UpdateHostInfo update the informations about a host using the HostData in the request
func (ctrl *HostDataController) UpdateHostInfo(w http.ResponseWriter, r *http.Request) {
	//Parse the hostdata from the request
	var originalHostData model.RawObject
	var aerr utils.AdvancedErrorInterface

	if err := json.NewDecoder(r.Body).Decode(&originalHostData); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Update and decode originalHostData
	if err := updateData(originalHostData); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//Validate the data
	documentLoader := gojsonschema.NewGoLoader(originalHostData)
	schemaLoader := gojsonschema.NewStringLoader(model.FrontendHostdataSchemaValidator)

	if result, err := gojsonschema.Validate(schemaLoader, documentLoader); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, "HOSTDATA_VALIDATION"))
		return
	} else if !result.Valid() {
		ctrl.Log.Printf("The input hostdata is not valid:\n")
		for _, desc := range result.Errors() {
			ctrl.Log.Printf("- %s\n", desc)
		}
		ctrl.Log.Println(utils.ToJSON(originalHostData))
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("Invalid schema. See the log"), "HOSTDATA_VALIDATION"))
		return
	}

	//Convert the originalHostData to a hostdata
	tempUpdatedRawJSON, err := json.Marshal(originalHostData)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	var updatedHostData model.HostDataBE
	err = json.Unmarshal(tempUpdatedRawJSON, &updatedHostData)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Save the HostData
	id, aerr := ctrl.Service.UpdateHostInfo(updatedHostData)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the created id
	utils.WriteJSONResponse(w, http.StatusOK, id)
}

// updateAndDecodeData return a decoded and updated hostdata from raw data
func updateData(data map[string]interface{}) utils.AdvancedErrorInterface {

	return nil
}
