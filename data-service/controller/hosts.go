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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"

	"github.com/xeipuuv/gojsonschema"
)

// InsertHostData update the informations about a host using the HostData in the request
func (ctrl *DataController) InsertHostData(w http.ResponseWriter, r *http.Request) {
	var hostdata model.HostDataBE

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	documentLoader := gojsonschema.NewBytesLoader(raw)
	schemaLoader := gojsonschema.NewStringLoader(model.FrontendHostdataSchemaValidator)
	if result, err := gojsonschema.Validate(schemaLoader, documentLoader); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, "HOSTDATA_VALIDATION"))
		return
	} else if !result.Valid() {
		var errorMsg strings.Builder
		errorMsg.WriteString("Invalid schema! The input hostdata is not valid!\n")

		for _, err := range result.Errors() {

			value := fmt.Sprintf("%v", err.Value())
			if len(value) > 80 {
				value = value[:78] + ".."
			}
			errorMsg.WriteString(fmt.Sprintf("- %s. Value: [%v]\n", err, value))
		}

		errorMsg.WriteString(fmt.Sprintf("hostdata:\n%v", string(raw)))

		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity,
			utils.NewAdvancedErrorPtr(errors.New(errorMsg.String()), "HOSTDATA_VALIDATION"))
		return
	}

	err = json.Unmarshal(raw, &hostdata)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	err = ctrl.Service.InsertHostData(hostdata)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
