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

// Package controller contains structs and methods used to provide endpoints for storing hostdata informations
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

//GetAwsRecommendations get recommendation related to cloud from Ercole
func (ctrl *ThunderController) GetAwsRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := ctrl.Service.GetAwsRecommendations()

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

func (ctrl *ThunderController) GetAwsRecommendationsErrors(w http.ResponseWriter, r *http.Request) {
	if seqValue, ok := mux.Vars(r)["seqnum"]; ok {
		seqValue, err := strconv.ParseUint(seqValue, 10, 64)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, err)
			return
		}

		recommendations, err := ctrl.Service.GetAwsRecommendationsBySeqValue(seqValue)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, dto.ToAwsRecommendationsErrorsDto(recommendations))

		return
	}

	recommendations, err := ctrl.Service.GetLastAwsRecommendations()
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dto.ToAwsRecommendationsErrorsDto(recommendations))
}

func (ctrl *ThunderController) ForceGetAwsRecommendations(w http.ResponseWriter, r *http.Request) {
	err := ctrl.Service.ForceGetAwsRecommendations()

	if errors.Is(err, utils.ErrClusterNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
