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
	"time"

	"github.com/ercole-io/ercole/chart-service/chartmodel"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
)

// GetChangeChart return the chart data related to changes to databases
func (as *ChartService) GetChangeChart(from time.Time, location string, environment string, olderThan time.Time) (chartmodel.ChangeChart, utils.AdvancedErrorInterface) {
	// get the old counts
	oldCounts, err := as.Database.GetTechnologyCount(location, environment, from)
	if err != nil {
		return chartmodel.ChangeChart{}, err
	}

	// get the new counts
	newCounts, err := as.Database.GetTechnologyCount(location, environment, olderThan)
	if err != nil {
		return chartmodel.ChangeChart{}, err
	}

	// build the bubbles
	bubbles := make([]chartmodel.ChangeChartBubble, 0)
	for t, v := range newCounts {
		if v > 0 {
			bubbles = append(bubbles, chartmodel.ChangeChartBubble{
				Name:   t,
				Size:   v,
				Change: v/oldCounts[t] - 1,
			})
			if oldCounts[t] == 0 {
				bubbles[len(bubbles)-1].Change = 0
			}
		}
	}

	return chartmodel.ChangeChart{
		Data: bubbles,
		Legend: map[string]string{
			"size": "Number of occurrences",
		},
	}, nil
}

// GetTechnologyTypesChart return the types of techonlogies
func (as *ChartService) GetTechnologyTypesChart(location string, environment string, olderThan time.Time) (chartmodel.TechnologyTypesChart, utils.AdvancedErrorInterface) {
	// get the counts
	counts, err := as.Database.GetTechnologyCount(location, environment, olderThan)
	if err != nil {
		return chartmodel.TechnologyTypesChart{}, err
	}

	out := chartmodel.TechnologyTypesChart{
		Legend: map[string]string{
			"size": "Number of occurrences",
		},
		OperatingSystems: make([]chartmodel.TechnologyTypeChartBubble, 0),
		Databases:        make([]chartmodel.TechnologyTypeChartBubble, 0),
		Middlewares:      make([]chartmodel.TechnologyTypeChartBubble, 0),
	}

	//databases
	if counts[model.TechnologyOracleDatabase] > 0 {
		out.Databases = append(out.Databases, chartmodel.TechnologyTypeChartBubble{
			Name: model.TechnologyOracleDatabase,
			Size: counts[model.TechnologyOracleDatabase],
		})
	}
	//middlewares
	//operating system
	for _, v := range as.Config.APIService.OperatingSystemAggregationRules {
		if counts[v.Product] > 0 {
			out.OperatingSystems = append(out.OperatingSystems, chartmodel.TechnologyTypeChartBubble{
				Name: v.Product,
				Size: counts[v.Product],
			})
		}
	}
	if counts[model.TechnologyUnknownOperatingSystem] > 0 {
		out.OperatingSystems = append(out.OperatingSystems, chartmodel.TechnologyTypeChartBubble{
			Name: model.TechnologyUnknownOperatingSystem,
			Size: counts[model.TechnologyUnknownOperatingSystem],
		})
	}

	return out, nil
}
