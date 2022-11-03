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
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

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

func (ctrl *APIController) AddRole(w http.ResponseWriter, r *http.Request) {
	role := &model.Role{}

	if err := utils.Decode(r.Body, role); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.AddRole(*role); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *APIController) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleName := mux.Vars(r)["roleName"]
	role := &model.Role{Name: roleName}

	if err := utils.Decode(r.Body, role); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if roleName != role.Name {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.ErrInvalidRole)
		return
	}

	if err := ctrl.Service.UpdateRole(*role); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *APIController) RemoveRole(w http.ResponseWriter, r *http.Request) {
	roleName := mux.Vars(r)["roleName"]

	groups, err := ctrl.Service.GetGroups()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	for _, group := range groups {
		for _, role := range group.Roles {
			if role == roleName {
				utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.ErrRoleCannotBeDeleted)
				return
			}
		}
	}

	if errDel := ctrl.Service.RemoveRole(roleName); errDel != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errDel)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
