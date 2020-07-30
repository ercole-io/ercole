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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	var hostdata model.HostDataBE
	var aerr utils.AdvancedErrorInterface

	//Read all bytes for the request
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	//Validate the data
	documentLoader := gojsonschema.NewBytesLoader(raw)
	schemaLoader := gojsonschema.NewStringLoader(model.FrontendHostdataSchemaValidator)
	if result, err := gojsonschema.Validate(schemaLoader, documentLoader); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, "HOSTDATA_VALIDATION"))
		return
	} else if !result.Valid() {
		var errorMsg strings.Builder
		errorMsg.WriteString("Invalid schema! The input hostdata is not valid!\n")

		for _, desc := range result.Errors() {
			errorMsg.WriteString(fmt.Sprintf("- %s\n", desc))
		}
		errorMsg.WriteString(fmt.Sprintf("hostdata:\n%v", string(raw)))

		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity,
			utils.NewAdvancedErrorPtr(errors.New(errorMsg.String()), "HOSTDATA_VALIDATION"))
		return
	}

	//Unmarshal raw to hostdata. The err checking is not needed
	_ = json.Unmarshal(raw, &hostdata)

	//Save the HostData
	id, aerr := ctrl.Service.UpdateHostInfo(hostdata)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the created id
	utils.WriteJSONResponse(w, http.StatusOK, id)
}
