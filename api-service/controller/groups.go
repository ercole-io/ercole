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

func (ctrl *APIController) InsertGroup(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	var group model.Group

	if validationErr := schema.ValidateGroup(raw); validationErr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, validationErr)

		return
	}

	err = json.Unmarshal(raw, &group)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	groupInserted, err := ctrl.Service.InsertGroup(group)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, groupInserted)
}

func (ctrl *APIController) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	name := mux.Vars(r)["name"]

	group := model.Group{Name: name}

	if validationErr := schema.ValidateGroup(raw); validationErr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, validationErr)

		return
	}

	err = json.Unmarshal(raw, &group)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	groupUpdated, err := ctrl.Service.UpdateGroup(group)
	if errors.Is(err, utils.ErrGroupNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, groupUpdated)
}

func (ctrl *APIController) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := ctrl.Service.GetGroups()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
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
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, group)
}

func (ctrl *APIController) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	users, err := ctrl.Service.ListUsers()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	for _, user := range users {
		for _, group := range user.Groups {
			if group == name {
				utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.ErrGroupCannotBeDeleted)
				return
			}
		}
	}

	errDel := ctrl.Service.DeleteGroup(name)
	if errors.Is(errDel, utils.ErrGroupNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, errDel)
		return
	}

	if errDel != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, errDel)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
