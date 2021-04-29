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
	"fmt"
	"net/http"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ctrl *APIController) AddMySQLAgreement(w http.ResponseWriter, r *http.Request) {
	var agreement model.MySQLAgreement

	if err := utils.Decode(r.Body, &agreement); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if agreement.ID != primitive.NilObjectID {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("ID must be empty"))
		return
	}

	id, err := ctrl.Service.AddMySQLAgreement(agreement)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, id)
}

func (ctrl *APIController) UpdateMySQLAgreement(w http.ResponseWriter, r *http.Request) {
	var agreement model.MySQLAgreement

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := utils.Decode(r.Body, &agreement); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if agreement.ID != id {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Object ID does not correspond"))
		return
	}

	err = ctrl.Service.UpdateMySQLAgreement(agreement)
	if errors.Is(err, utils.ErrNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, agreement)
}

func (ctrl *APIController) GetMySQLAgreements(w http.ResponseWriter, r *http.Request) {
	agreements, err := ctrl.Service.GetMySQLAgreements()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"agreements": agreements,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) DeleteMySQLAgreement(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, fmt.Errorf("Can't decode id: %w", err))
		return
	}

	err = ctrl.Service.DeleteMySQLAgreement(id)
	if errors.Is(err, utils.ErrNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
