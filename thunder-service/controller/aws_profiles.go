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
	"strconv"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ctrl *ThunderController) AddAwsProfile(w http.ResponseWriter, r *http.Request) {
	var profile model.AwsProfile

	if err := utils.Decode(r.Body, &profile); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if profile.ID != primitive.NilObjectID {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("ID must be empty"))
		return
	}

	if profile.SecretAccessKey == nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("SecretAccessKey must not be null"))
		return
	}

	profileAdded, err := ctrl.Service.AddAwsProfile(profile)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, profileAdded)
}

func (ctrl *ThunderController) UpdateAwsProfile(w http.ResponseWriter, r *http.Request) {
	var profile model.AwsProfile

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

	profileUpdated, err := ctrl.Service.UpdateAwsProfile(profile)
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

func (ctrl *ThunderController) GetAwsProfiles(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetAwsProfiles()

	if errors.Is(err, utils.ErrClusterNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

func (ctrl *ThunderController) DeleteAwsProfile(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, fmt.Errorf("Can't decode id: %w", err))
		return
	}

	err = ctrl.Service.DeleteAwsProfile(id)
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

func (ctrl *ThunderController) SelectAwsProfile(w http.ResponseWriter, r *http.Request) {
	profileId := mux.Vars(r)["profileid"]
	selected, err := strconv.ParseBool(mux.Vars(r)["selected"])

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(err, "BAD_REQUEST"))
		return
	}

	err = ctrl.Service.SelectAwsProfile(profileId, selected)
	if errors.Is(err, utils.ErrNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
