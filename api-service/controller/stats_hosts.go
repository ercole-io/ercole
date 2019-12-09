// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/amreo/ercole-services/utils"
)

// GetEnvironmentStats return all statistics about the environments using the filters in the request
func (ctrl *APIController) GetEnvironmentStats(w http.ResponseWriter, r *http.Request) {
	var location string
	var err utils.AdvancedErrorInterface

	//parse the query params
	location = r.URL.Query().Get("location")

	//get the data
	stats, err := ctrl.Service.GetEnvironmentStats(location)
	if err != nil {
		utils.WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
