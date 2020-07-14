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
	"errors"
	"net/http"
	"time"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
)

// GetOracleDatabaseChart return the list of techonlogies
func (ctrl *ChartController) GetOracleDatabaseChart(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface
	var metric string
	var location string
	var environment string
	var olderThan time.Time

	metric = r.URL.Query().Get("metric")
	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")
	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if !utils.Contains(model.TechnologiesSupportedMetricsMap[model.TechnologyOracleDatabase].Metrics, metric) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("Unrecognized"), http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	data, err := ctrl.Service.GetOracleDatabaseChart(metric, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}
