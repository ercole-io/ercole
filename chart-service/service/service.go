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

// Package service is a package that provides methods for querying data
package service

import (
	"math/rand"
	"time"

	"github.com/ercole-io/ercole/chart-service/chartmodel"
	"github.com/ercole-io/ercole/chart-service/database"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/config"
)

// ChartServiceInterface is a interface that wrap methods used to querying data
type ChartServiceInterface interface {
	// Init initialize the service
	Init()

	// GetOracleDatabaseChart return a chart associated to teh
	GetOracleDatabaseChart(metric string, location string, environment string, olderThan time.Time) (chartmodel.Chart, utils.AdvancedErrorInterface)

	// GetErcoleFeatures return the list of technologies
	GetTechnologyList() (map[string]model.TechnologySupportedMetrics, utils.AdvancedErrorInterface)
}

// ChartService is the concrete implementation of APIServiceInterface.
type ChartService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log *logrus.Logger
	// Random contains the generator used to generate colors
	Random *rand.Rand
}

// Init initializes the service and database
func (as *ChartService) Init() {
	as.Random = rand.New(rand.NewSource(as.TimeNow().UnixNano()))
}

// GetTechnologyList return the list of technologies
func (as *ChartService) GetTechnologyList() (map[string]model.TechnologySupportedMetrics, utils.AdvancedErrorInterface) {
	// at the moment, the list of technologies is hardcoded here
	return model.TechnologySupportedMetricsMap, nil
}
