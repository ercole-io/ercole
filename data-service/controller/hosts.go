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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/schema"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/sanitizer"
)

// InsertHostData update the informations about a host using the HostData in the request
func (ctrl *DataController) InsertHostData(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	if raw, err = ctrl.sanitizeJsonHostdata(raw); err != nil {
		if errors.Is(err, utils.ErrInvalidHostdata) {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
			ctrl.Service.AlertInvalidHostData(err, nil)
			return
		}

		ctrl.Log.Error(err)
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, nil)
		return
	}

	var hostdata model.HostDataBE
	if validationErr := schema.ValidateHostdata(raw); validationErr != nil {
		if errors.Is(validationErr, utils.ErrInvalidHostdata) {
			ctrl.Log.Info(validationErr)
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, validationErr)

			if unmarshalErr := json.Unmarshal(raw, &hostdata); unmarshalErr != nil {
				ctrl.Service.AlertInvalidHostData(validationErr, nil)
			} else {
				ctrl.Service.AlertInvalidHostData(validationErr, &hostdata)
			}

			return
		}

		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, validationErr)
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

func (ctrl *DataController) sanitizeJsonHostdata(raw []byte) ([]byte, error) {
	var m map[string]interface{}

	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, utils.ErrInvalidHostdata
	}

	sanitizer := sanitizer.NewSanitizer(ctrl.Log)

	sanitizedInt, err := sanitizer.Sanitize(m)
	if err != nil {
		return nil, fmt.Errorf("Unable to sanitize: %w", err)
	}

	if raw, err = json.Marshal(sanitizedInt); err != nil {
		return nil, fmt.Errorf("Unable to marshal: %w", err)
	}

	return raw, nil
}
