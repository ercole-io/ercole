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
	"sort"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
)

func (as *ChartService) GetLicenseComplianceHistory() ([]dto.LicenseComplianceHistory, error) {
	licenses, err := as.Database.GetLicenseComplianceHistory()
	if err != nil {
		return nil, err
	}

	types, err := as.getOracleDatabaseLicenseTypes()
	if err != nil {
		return nil, err
	}

	for i := range licenses {
		license := &licenses[i]

		if len(license.LicenseTypeID) > 0 {
			if licenseType, ok := types[license.LicenseTypeID]; ok {
				license.ItemDescription = licenseType.ItemDescription
				license.Metric = licenseType.Metric
			}
		}

		license.History = sortAndKeepOnlyLastEntryOfEachDay(license.History)
	}

	licenses = mergeMySqlLicensesCompliance(licenses)
	licenses = removeEmptyLicensesCompliance(licenses)

	return licenses, nil
}

func sortAndKeepOnlyLastEntryOfEachDay(history []dto.LicenseComplianceHistoricValue) []dto.LicenseComplianceHistoricValue {
	sort.Slice(history, func(i, j int) bool {
		return history[i].Date.Before(history[j].Date)
	})

	newHistory := make([]dto.LicenseComplianceHistoricValue, 0, len(history))

	var nextEntry *dto.LicenseComplianceHistoricValue

	for i := range history {
		val := &history[i]
		valDay := time.Date(val.Date.Year(), val.Date.Month(), val.Date.Day(), 0, 0, 0, 0, val.Date.Location())

		if nextEntry == nil || valDay.Equal(nextEntry.Date) {
			nextEntry = val
			nextEntry.Date = valDay

			continue
		}

		if valDay.After(nextEntry.Date) {
			newHistory = append(newHistory, *nextEntry)
			nextEntry.Date = time.Date(val.Date.Year(), val.Date.Month(), val.Date.Day(), 0, 0, 0, 0, val.Date.Location())
			nextEntry = val
		}
	}

	if nextEntry != nil {
		newHistory = append(newHistory, *nextEntry)
	}

	return newHistory
}

func mergeMySqlLicensesCompliance(licenses []dto.LicenseComplianceHistory) []dto.LicenseComplianceHistory {
	var mySql *dto.LicenseComplianceHistory

	mySqlEnterprisePrefix := "MySQL Enterprise"

	for i := len(licenses) - 1; i >= 0; i-- {
		l := licenses[i]

		if len(l.LicenseTypeID) > 0 || !strings.HasPrefix(l.ItemDescription, mySqlEnterprisePrefix) {
			continue
		}

		if mySql == nil {
			mySql = new(dto.LicenseComplianceHistory)
			mySql.ItemDescription = mySqlEnterprisePrefix
		}

		mySql.History = mergeLicenseComplianceHistoricValues(mySql.History, l.History)

		licenses = append(licenses[:i], licenses[i+1:]...)
	}

	if mySql == nil {
		return licenses
	}

	return append(licenses, *mySql)
}

func mergeLicenseComplianceHistoricValues(a, b []dto.LicenseComplianceHistoricValue) []dto.LicenseComplianceHistoricValue {
	merged := make([]dto.LicenseComplianceHistoricValue, 0)

	for i, j := 0, 0; i < len(a) || j < len(b); {
		var valA, valB *dto.LicenseComplianceHistoricValue

		if i < len(a) {
			valA = &a[i]
		}

		if j < len(b) {
			valB = &b[j]
		}

		if valA != nil && (valB == nil || valA.Date.Before(valB.Date)) {
			merged = append(merged, *valA)
			i += 1

			continue
		}

		if valA == nil || valB.Date.Before(valA.Date) {
			merged = append(merged, *valB)
			j += 1

			continue
		}

		newVal := dto.LicenseComplianceHistoricValue{
			Date:      valA.Date,
			Consumed:  valA.Consumed + valB.Consumed,
			Covered:   valA.Covered + valB.Covered,
			Purchased: valA.Purchased + valB.Purchased,
		}

		merged = append(merged, newVal)

		i, j = i+1, j+1
	}

	return merged
}

func removeEmptyLicensesCompliance(licenses []dto.LicenseComplianceHistory) []dto.LicenseComplianceHistory {
	result := make([]dto.LicenseComplianceHistory, 0)

licenses:
	for i := range licenses {
		l := &licenses[i]

		for _, x := range l.History {
			if x.Consumed > 0 || x.Covered > 0 || x.Purchased > 0 {
				result = append(result, *l)
				continue licenses
			}
		}
	}

	return result
}
