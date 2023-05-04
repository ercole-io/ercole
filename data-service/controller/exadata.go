// Copyright (c) 2023 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/schema"
	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *DataController) InsertExadata(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	if raw, err = ctrl.sanitizeJson(raw); err != nil {
		if errors.Is(err, utils.ErrInvalidExadata) {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)

			return
		}

		ctrl.Log.Error(err)
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)

		return
	}

	var exadata model.OracleExadataInstance

	if validationErr := schema.ValidateExadata(raw); validationErr != nil {
		if errors.Is(validationErr, utils.ErrInvalidExadata) {
			ctrl.Log.Info(validationErr)
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, validationErr)

			return
		}

		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, validationErr)

		return
	}

	err = json.Unmarshal(raw, &exadata)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	err = ctrl.Service.SaveExadata(&exadata)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
