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

	apiservice_client "github.com/ercole-io/ercole/v2/api-service/client"
	"github.com/ercole-io/ercole/v2/chart-service/database"
	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
)

// ChartServiceInterface is a interface that wrap methods used to querying data
type ChartServiceInterface interface {
	// Init initialize the service
	Init()

	// GetChangeChart return the chart data related to changes to databases
	GetChangeChart(from time.Time, location string, environment string, olderThan time.Time) (dto.ChangeChart, error)

	// GetOracleDatabaseChart return a chart associated to teh
	GetOracleDatabaseChart(metric string, location string, environment string, olderThan time.Time) (dto.Chart, error)
	GetLicenseComplianceHistory() ([]dto.LicenseComplianceHistory, error)

	// GetTechnologiesMetrics return metrics of all technologies
	GetTechnologiesMetrics() (map[string]model.TechnologySupportedMetrics, error)
	// GetTechnologyTypes return the types of techonlogies
	GetTechnologyTypesChart(location string, environment string, olderThan time.Time) (dto.TechnologyTypesChart, error)

	GetHostCores(location string, environment string, olderThan time.Time, newerThan time.Time) ([]dto.HostCores, error)
}

type ChartService struct {
	Config       config.Configuration
	Database     database.MongoDatabaseInterface
	ApiSvcClient apiservice_client.ApiSvcClientInterface
	TimeNow      func() time.Time
	Log          logger.Logger
	// Random contains the generator used to generate colors
	Random *rand.Rand
}

func (as *ChartService) Init() {
	as.Random = rand.New(rand.NewSource(as.TimeNow().UnixNano()))
}

// GetTechnologiesMetrics return the list of technologies
func (as *ChartService) GetTechnologiesMetrics() (map[string]model.TechnologySupportedMetrics, error) {
	// at the moment, the list of technologies is hardcoded here
	return model.TechnologiesSupportedMetricsMap, nil
}
