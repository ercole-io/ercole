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

func (ctrl *APIController) InsertGroup(w http.ResponseWriter, r *http.Request) {
	var group model.Group

	if err := utils.Decode(r.Body, &group); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	groupInserted, err := ctrl.Service.InsertGroup(group)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, groupInserted)
}

func (ctrl *APIController) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	group := model.Group{Name: name}
	if err := utils.Decode(r.Body, &group); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	groupUpdated, err := ctrl.Service.UpdateGroup(group)
	if errors.Is(err, utils.ErrGroupNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, groupUpdated)
}

func (ctrl *APIController) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := ctrl.Service.GetGroups()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"groups": groups,
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetGroup(w http.ResponseWriter, r *http.Request) {
	var err error

	name := mux.Vars(r)["name"]

	group, err := ctrl.Service.GetGroup(name)
	if errors.Is(err, utils.ErrGroupNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, group)
}

func (ctrl *APIController) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	errDel := ctrl.Service.DeleteGroup(name)
	if errors.Is(errDel, utils.ErrGroupNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, errDel)
		return
	}

	if errDel != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, errDel)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
