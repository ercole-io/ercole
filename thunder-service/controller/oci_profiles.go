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

// Package controller contains structs and methods used to provide endpoints for storing hostdata informations
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

func (ctrl *ThunderController) AddOciProfile(w http.ResponseWriter, r *http.Request) {
	var profile model.OciProfile

	if err := utils.Decode(r.Body, &profile); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if profile.ID != primitive.NilObjectID {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("ID must be empty"))
		return
	}

	if profile.PrivateKey == nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("PrivateKey must not be null"))
		return
	}

	if !profile.IsValid() {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Profile configuration isn't valid"))
		return
	}

	profileAdded, err := ctrl.Service.AddOciProfile(profile)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, profileAdded)
}

func (ctrl *ThunderController) UpdateOciProfile(w http.ResponseWriter, r *http.Request) {
	var profile model.OciProfile

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := utils.Decode(r.Body, &profile); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if profile.ID != id {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Object ID does not correspond"))
		return
	}

	if !profile.IsValid() {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Some profile fields are not valid"))
		return
	}

	profileUpdated, err := ctrl.Service.UpdateOciProfile(profile)
	if errors.Is(err, utils.ErrNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, profileUpdated)
}

func (ctrl *ThunderController) GetOciProfiles(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetOciProfiles()

	if errors.Is(err, utils.ErrClusterNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

func (ctrl *ThunderController) DeleteOciProfile(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, fmt.Errorf("Can't decode id: %w", err))
		return
	}

	err = ctrl.Service.DeleteOciProfile(id)
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
