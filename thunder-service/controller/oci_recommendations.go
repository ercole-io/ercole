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

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
)

//GetOciRecommendations get recommendation related to cloud from Ercole
func (ctrl *ThunderController) GetOciRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := ctrl.Service.GetOciRecommendations()

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

	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		response := map[string]interface{}{
			"recommendations": recommendations,
			"error":           err,
		}

		utils.WriteJSONResponse(w, http.StatusOK, response)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		response, err := ctrl.Service.WriteOciRecommendationsXlsx(recommendations)
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
			return
		}

		utils.WriteXLSXResponse(w, response)
	}
}

func (ctrl *ThunderController) ForceGetOciRecommendations(w http.ResponseWriter, r *http.Request) {
	err := ctrl.Service.ForceGetOciRecommendations()

	if errors.Is(err, utils.ErrClusterNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
