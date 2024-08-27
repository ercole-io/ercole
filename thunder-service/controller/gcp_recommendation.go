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

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
)

func (ctrl *ThunderController) GetGcpRecommendations(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		recommendations, err := ctrl.Service.ListGcpRecommendations()
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		}

		utils.WriteJSONResponse(w, http.StatusOK, recommendations)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		result, err := ctrl.Service.CreateGcpRecommendationsXlsx()
		if err != nil {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
			return
		}

		utils.WriteXLSXResponse(w, result)
	}
}

func (ctrl *ThunderController) ForceGetGcpRecommendations(w http.ResponseWriter, r *http.Request) {
	go ctrl.Service.ForceGetGcpRecommendations()

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *ThunderController) GetGcpErrors(w http.ResponseWriter, r *http.Request) {
	gcpErrors, err := ctrl.Service.ListGcpError()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
	}

	utils.WriteJSONResponse(w, http.StatusOK, gcpErrors)
}
