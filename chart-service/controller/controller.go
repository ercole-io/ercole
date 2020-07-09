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

	"github.com/ercole-io/ercole/api-service/auth"
	"github.com/ercole-io/ercole/chart-service/service"
	"github.com/ercole-io/ercole/utils"

	"github.com/ercole-io/ercole/config"
	"github.com/sirupsen/logrus"
)

// ChartControllerInterface is a interface that wrap methods used to querying data
type ChartControllerInterface interface {
	// GetOracleDatabaseChart return the chart data related to oracle databases
	GetOracleDatabaseChart(w http.ResponseWriter, r *http.Request)
	// GetChangeChart return the chart data related to changes
	GetChangeChart(w http.ResponseWriter, r *http.Request)
	// GetTechnologyTypes return the types of techonlogies
	GetTechnologyTypes(w http.ResponseWriter, r *http.Request)

	// GetTechnologyList return the list of techonlogies
	GetTechnologyList(w http.ResponseWriter, r *http.Request)
}

// ChartController is the struct used to handle the requests from agents and contains the concrete implementation of ChartControllerInterface
type ChartController struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Service contains the underlying service used to perform various logical and store operations
	Service service.ChartServiceInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log *logrus.Logger
	// Authenticator contains the authenticator
	Authenticator auth.AuthenticationProvider
}

// GetTechnologyList return the list of techonlogies
func (ctrl *ChartController) GetTechnologyList(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetTechnologyList()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}
