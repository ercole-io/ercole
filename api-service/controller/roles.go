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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/schema"
	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *APIController) InsertRole(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	var role model.Role

	if validationErr := schema.ValidateRole(raw); validationErr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, validationErr)

		return
	}

	err = json.Unmarshal(raw, &role)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	roleInserted, err := ctrl.Service.InsertRole(role)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, roleInserted)
}

func (ctrl *APIController) UpdateRole(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	role := model.Role{Name: name}
	if err := utils.Decode(r.Body, &role); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	roleUpdated, err := ctrl.Service.UpdateRole(role)
	if errors.Is(err, utils.ErrRoleNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, roleUpdated)
}

func (ctrl *APIController) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := ctrl.Service.GetRoles()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	response := map[string]interface{}{
		"roles": roles,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetRole(w http.ResponseWriter, r *http.Request) {
	var err error

	name := mux.Vars(r)["name"]

	role, err := ctrl.Service.GetRole(name)
	if errors.Is(err, utils.ErrRoleNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, role)
}

func (ctrl *APIController) DeleteRole(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	errDel := ctrl.Service.DeleteRole(name)
	if errors.Is(errDel, utils.ErrRoleNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, errDel)
		return
	}

	if errDel != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errDel)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
