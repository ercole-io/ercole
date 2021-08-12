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
)

//GetOCRecommendations get recommendation from Oracle Cloud
func (ctrl *ThunderController) GetOCRecommendations(w http.ResponseWriter, r *http.Request) {

	data, err := ctrl.Service.GetOCRecommendations("ocid1.tenancy.oc1..aaaaaaaazizzbqqbjv2se3y3fvm5osfumnorh32nznanirqoju3uks4buh4q")

	if errors.Is(err, utils.ErrClusterNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

//GetOCRecommendations get recommendation with Category from Oracle Cloud
func (ctrl *ThunderController) GetOCRecommendationsWithCategory(w http.ResponseWriter, r *http.Request) {

	data, err := ctrl.Service.GetOCRecommendationsWithCategory("ocid1.tenancy.oc1..aaaaaaaazizzbqqbjv2se3y3fvm5osfumnorh32nznanirqoju3uks4buh4q")

	if errors.Is(err, utils.ErrClusterNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}
