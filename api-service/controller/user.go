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

	"github.com/ercole-io/ercole/v2/api-service/auth"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ctrl.Service.ListUsers()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dto.ToUsers(users))
}

func (ctrl *APIController) GetUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	user, err := ctrl.Service.GetUser(username, auth.BasicType)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dto.ToUser(user))
}

func (ctrl *APIController) AddUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}

	if err := utils.Decode(r.Body, user); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if user.Password == "" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Invalid password"))
		return
	}

	user.Provider = auth.BasicType

	if err := ctrl.Service.AddUser(*user); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *APIController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	provider := mux.Vars(r)["provider"]

	user := model.User{Username: username}
	if err := utils.Decode(r.Body, &user); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.UpdateUserGroups(username, provider, user.Groups); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *APIController) RemoveUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	provider := mux.Vars(r)["provider"]

	if username == model.SuperUser {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.ErrSuperUserCannotBeDeleted)
		return
	}

	if err := ctrl.Service.RemoveUser(username, provider); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *APIController) NewPassword(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	user, err := ctrl.Service.GetUser(username, "basic")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	newPassword, err := ctrl.Service.NewPassword(username)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := ctrl.Service.AddLimitedGroup(*user); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, newPassword)
}

func (ctrl *APIController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	user := context.Get(r, "user")
	if user != nil {
		if !user.(*model.User).IsAdmin() && user.(*model.User).Username != username {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnauthorized, utils.ErrInvalidUser)
			return
		}
	}

	type changes struct {
		OldPass       string `json:"oldPassword"`
		NewPass       string `json:"newPassword"`
		ConfirmedPass string `json:"confirmedPassword"`
	}

	reqChanges := changes{}

	if err := utils.Decode(r.Body, &reqChanges); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if reqChanges.NewPass != reqChanges.ConfirmedPass {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Invalid password"))
		return
	}

	if err := ctrl.Service.UpdatePassword(username, reqChanges.OldPass, reqChanges.NewPass); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	u := user.(*model.User)
	if err := ctrl.Service.RemoveLimitedGroup(*u); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
