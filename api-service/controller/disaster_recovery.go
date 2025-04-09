// Copyright (c) 2025 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package controller

import (
	"net/http"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) CreateDr(w http.ResponseWriter, r *http.Request) {
	hostname := mux.Vars(r)["hostname"]

	type data struct {
		ClusterVeritasHostnames []string `json:"clusterVeritasHostnames"`
	}

	d := data{}

	if err := utils.Decode(r.Body, &d); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	drName, err := ctrl.Service.CreateDR(hostname, d.ClusterVeritasHostnames)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{
		"hostname": drName,
	})
}
