// Copyright (c) 2022 Sorint.lab S.p.A.
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

	"github.com/golang/gddo/httputil"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *APIController) AddMySQLContract(w http.ResponseWriter, r *http.Request) {
	var contract model.MySQLContract

	if err := utils.Decode(r.Body, &contract); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if contract.ID != primitive.NilObjectID {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("ID must be empty"))
		return
	}

	if !contract.IsValid() {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Contract isn't valid"))
		return
	}

	contractAdded, err := ctrl.Service.AddMySQLContract(contract)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, contractAdded)
}

func (ctrl *APIController) UpdateMySQLContract(w http.ResponseWriter, r *http.Request) {
	var contract model.MySQLContract

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := utils.Decode(r.Body, &contract); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if contract.ID != id {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Object ID does not correspond"))
		return
	}

	if !contract.IsValid() {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Contract isn't valid"))
		return
	}

	contractUpdated, err := ctrl.Service.UpdateMySQLContract(contract)
	if errors.Is(err, utils.ErrNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, contractUpdated)
}

func (ctrl *APIController) GetMySQLContracts(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")

	locations, err := ctrl.Service.ListLocations(user)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.GetMySQLContractsJSON(w, r, locations)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.GetMySQLContractsXLSX(w, r, locations)
	}
}

func (ctrl *APIController) GetMySQLContractsJSON(w http.ResponseWriter, r *http.Request, locations []string) {
	contracts, err := ctrl.Service.GetMySQLContracts(locations)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"contracts": contracts,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetMySQLContractsXLSX(w http.ResponseWriter, r *http.Request, locations []string) {
	xlsx, err := ctrl.Service.GetMySQLContractsAsXLSX(locations)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

func (ctrl *APIController) DeleteMySQLContract(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, fmt.Errorf("Can't decode id: %w", err))
		return
	}

	err = ctrl.Service.DeleteMySQLContract(id)
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
