// Copyright (c) 2022 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// GetOracleDatabaseChart return a chart associated to teh
func (as *ChartService) GetOracleDatabaseChart(metric string, location string, environment string, olderThan time.Time) (dto.Chart, error) {
	switch metric {
	case "version":
		data, err := as.Database.GetOracleDatabaseChartByVersion(location, environment, olderThan)
		if err != nil {
			return dto.Chart{}, err
		}

		// colorize the data
		for i := range data {
			data[i].Color = dto.RandomColorize(*as.Random)
		}

		// return the data
		return dto.Chart{
			Data: data,
			Legend: map[string]string{
				"size": "Number of occurrences",
			},
		}, nil
	case "work":
		data, err := as.Database.GetOracleDatabaseChartByWork(location, environment, olderThan)
		if err != nil {
			return dto.Chart{}, err
		}

		// colorize the data
		for i := range data {
			data[i].Color = dto.RandomColorize(*as.Random)
		}

		// return the data
		return dto.Chart{
			Data: data,
			Legend: map[string]string{
				"size": "Value of work",
			},
		}, nil
	default:
		return dto.Chart{}, utils.NewError(errors.New("Unsupported metric"), "UNSUPPORTED_METRIC")
	}
}

func (as *ChartService) getOracleDatabaseLicenseTypes() (map[string]model.OracleDatabaseLicenseType, error) {
	licenseTypes, err := as.ApiSvcClient.GetOracleDatabaseLicenseTypes()
	if err != nil {
		return nil, utils.NewError(err, "Can't retrieve Oracle licenseTypes")
	}

	licenseTypesMap := make(map[string]model.OracleDatabaseLicenseType)
	for _, licenseType := range licenseTypes {
		licenseTypesMap[licenseType.ID] = licenseType
	}

	return licenseTypesMap, nil
}
