// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"net/http"

	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

func (ctrl *ThunderController) GetGcpProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := ctrl.Service.GetGcpProfiles()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, profiles)
}

func (ctrl *ThunderController) AddGcpProfile(w http.ResponseWriter, r *http.Request) {
	var profile dto.GcpProfileRequest

	if err := utils.Decode(r.Body, &profile); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.AddGcpProfile(profile); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *ThunderController) SelectGcpProfile(w http.ResponseWriter, r *http.Request) {
	profileId := mux.Vars(r)["profileid"]

	if err := ctrl.Service.SelectGcpProfile(profileId); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *ThunderController) UpdateGcpProfile(w http.ResponseWriter, r *http.Request) {
	profileID := mux.Vars(r)["profileid"]

	profileReq := dto.GcpProfileRequest{}
	if err := utils.Decode(r.Body, &profileReq); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.UpdateGcpProfile(profileID, profileReq); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *ThunderController) RemoveGcpProfile(w http.ResponseWriter, r *http.Request) {
	profileID := mux.Vars(r)["profileid"]

	if err := ctrl.Service.RemoveGcpProfile(profileID); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
