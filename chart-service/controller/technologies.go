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

package controller

import (
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/utils"
)

// GetChangeChart return the chart data related to changes
func (ctrl *ChartController) GetChangeChart(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface
	var from time.Time
	var location string
	var environment string
	var olderThan time.Time

	if from, err = utils.Str2time(r.URL.Query().Get("from"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")
	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	data, err := ctrl.Service.GetChangeChart(from, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

// GetTechnologyTypes return the types of techonlogies
func (ctrl *ChartController) GetTechnologyTypes(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface
	var location string
	var environment string
	var olderThan time.Time

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")
	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	data, err := ctrl.Service.GetTechnologyTypesChart(location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}
