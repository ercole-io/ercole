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
	"net/http"
	"strings"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

//GetOciRecommendations get recommendation from Oracle Cloud
func (ctrl *ThunderController) GetOciRecommendations(w http.ResponseWriter, r *http.Request) {
	profileList := mux.Vars(r)["ids"]
	if profileList == "" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("Ids not present or malformed"))
		return
	}

	var profiles []string = strings.Split(profileList, ",")

	recommendations, err := ctrl.Service.GetOciRecommendations(profiles)

	if recommendations == nil {
		if errors.Is(err, utils.ErrInvalidProfileId) {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
			return
		}

		if errors.Is(err, utils.ErrClusterNotFound) {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
			return
		}

		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)

		return
	}

	if err == nil {
		response := map[string]interface{}{
			"recommendations": recommendations,
		}
		utils.WriteJSONResponse(w, http.StatusOK, response)

		return
	}

	response := map[string]interface{}{
		"recommendations": recommendations,
		"error":           err.Error(),
	}

	utils.WriteJSONResponse(w, http.StatusPartialContent, response)
}
