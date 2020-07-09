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
	"errors"
	"time"

	"github.com/ercole-io/ercole/chart-service/chartmodel"
	"github.com/ercole-io/ercole/utils"
)

// GetOracleDatabaseChart return a chart associated to teh
func (as *ChartService) GetOracleDatabaseChart(metric string, location string, environment string, olderThan time.Time) (chartmodel.Chart, utils.AdvancedErrorInterface) {
	switch metric {
	case "version":
		data, err := as.Database.GetOracleDatabaseChartByVersion(location, environment, olderThan)
		if err != nil {
			return chartmodel.Chart{}, err
		}

		// colorize the data
		for i := range data {
			data[i].Color = chartmodel.RandomColorize(*as.Random)
		}

		// return the data
		return chartmodel.Chart{
			Data: data,
			Legend: map[string]string{
				"size": "Number of occurrences",
			},
		}, nil
	case "work":
		data, err := as.Database.GetOracleDatabaseChartByWork(location, environment, olderThan)
		if err != nil {
			return chartmodel.Chart{}, err
		}

		// colorize the data
		for i := range data {
			data[i].Color = chartmodel.RandomColorize(*as.Random)
		}

		// return the data
		return chartmodel.Chart{
			Data: data,
			Legend: map[string]string{
				"size": "Value of work",
			},
		}, nil
	default:
		return chartmodel.Chart{}, utils.NewAdvancedErrorPtr(errors.New("Unsupported metric"), "UNSUPPORTED_METRIC")
	}
}
