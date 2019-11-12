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
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/model"
	"github.com/goji/httpauth"
)

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (ctrl *HostDataController) AuthenticateMiddleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(ctrl.Config.DataService.AgentUsername, ctrl.Config.DataService.AgentPassword)
}

// UpdateHostInfo update the informations about a host using the HostData in the request
func (ctrl *HostDataController) UpdateHostInfo(w http.ResponseWriter, r *http.Request) {
	//Parse the hostdata from the request
	var originalHostData map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&originalHostData); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Update and decode originalHostData
	hostData, err := updateAndDecodeData(originalHostData)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// fixes
	setHostnameAgentVirtualizationToClusters(&hostData)

	//Save the HostData
	id, err := ctrl.Service.UpdateHostInfo(hostData)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the created id
	utils.WriteJSONResponse(w, http.StatusOK, id)
}

// updateAndDecodeData return a decoded and updated hostdata from raw data
func updateAndDecodeData(data map[string]interface{}) (model.HostData, utils.AdvancedErrorInterface) {
	var hostDataSchemaVersion int

	//get correct hostDataSchemaVersion and fix the version
	if val, ok := data["HostDataSchemaVersion"]; !ok {
		hostDataSchemaVersion = 0
	} else if val, ok := val.(float64); !ok {
		return model.HostData{}, utils.NewAdvancedErrorPtr(
			errors.New("Invalid type for $hostDataSchemaVersion property"),
			http.StatusText(http.StatusUnprocessableEntity))
	} else {
		hostDataSchemaVersion = int(val)
	}

	//fix the version
	if val, ok := data["Version"]; !ok {
		data["Version"] = "pre1.5.6"
	} else if val, ok := val.(string); !ok {
		return model.HostData{}, utils.NewAdvancedErrorPtr(
			errors.New("Invalid type for $version property"),
			http.StatusText(http.StatusUnprocessableEntity))
	} else if val == "${VERSION}" {
		data["Version"] = "pre1.5.11"
	}

	//Update the hostData to the version 1
	if hostDataSchemaVersion < 1 {
		data = updateHostDataSchemaTo1(data)
	}

	//Decode the raw data
	var hostData model.HostData
	raw, _ := json.Marshal(data)
	if err := json.Unmarshal(raw, &hostData); err != nil {
		return model.HostData{}, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity))
	}

	//Return the decodec hostData
	return hostData, nil
}

// setHostnameAgentVirtualizationToClusters set the hostname of itself to all clusters inside the hostdata
func setHostnameAgentVirtualizationToClusters(orig *model.HostData) {
	for i := range orig.Extra.Clusters {
		orig.Extra.Clusters[i].HostnameAgentVirtualization = orig.Hostname
	}
}

// updateHostDataSchemaTo1 update the schema in the data to the version one
func updateHostDataSchemaTo1(data map[string]interface{}) map[string]interface{} {
	if _, ok := data["HostType"]; !ok {
		data["HostType"] = "oracledb"
	}

	data["HostDataSchemaVersion"] = 1
	return data
}
