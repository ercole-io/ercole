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
	"encoding/json"
	"errors"
	"net/http"
	"sort"
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
		return dto.Chart{}, utils.NewAdvancedErrorPtr(errors.New("Unsupported metric"), "UNSUPPORTED_METRIC")
	}
}

func (as *ChartService) GetOracleDbLicenseHistory() ([]dto.OracleDatabaseLicenseHistory, error) {
	licenses, err := as.Database.GetOracleDbLicenseHistory()
	if err != nil {
		return nil, err
	}

	types, err := as.getOracleDatabaseLicenseTypes()
	if err != nil {
		return nil, err
	}

	for i := range licenses {
		license := &licenses[i]
		licenseType := types[license.LicenseTypeID]
		license.ItemDescription = licenseType.ItemDescription
		license.Metric = licenseType.Metric

		license.History = keepOnlyLastEntryOfEachDay(license.History)
	}

	return licenses, nil
}

func (as *ChartService) getOracleDatabaseLicenseTypes() (map[string]model.OracleDatabaseLicenseType, error) {
	url := utils.NewAPIUrlNoParams(
		as.Config.APIService.RemoteEndpoint,
		as.Config.APIService.AuthenticationProvider.Username,
		as.Config.APIService.AuthenticationProvider.Password,
		"/settings/oracle/database/license-types").String()
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, utils.NewAdvancedErrorPtr(err, "Can't retrieve from databases")
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	var licenseTypes []model.OracleDatabaseLicenseType
	if err := decoder.Decode(&licenseTypes); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Can't decode response body")
	}

	licenseTypesMap := make(map[string]model.OracleDatabaseLicenseType)
	for _, licenseType := range licenseTypes {
		licenseTypesMap[licenseType.ID] = licenseType
	}

	return licenseTypesMap, nil
}

func keepOnlyLastEntryOfEachDay(history []dto.OracleDbHistoricValue) []dto.OracleDbHistoricValue {
	sort.Slice(history, func(i, j int) bool {
		return history[i].Date.After(history[j].Date)
	})

	currentDay := utils.MAX_TIME
	newHistory := make([]dto.OracleDbHistoricValue, 0, len(history))

	for i := range history {
		entry := &history[i]
		entryDate := entry.Date
		entryDay := time.Date(entryDate.Year(), entryDate.Month(), entryDate.Day(), 0, 0, 0, 0, time.UTC)

		if entryDay.Before(currentDay) {
			currentDay = entryDay

			entry.Date = entryDay
			newHistory = append(newHistory, *entry)
		}
	}

	return newHistory
}
