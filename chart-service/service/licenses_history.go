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

import "github.com/ercole-io/ercole/v2/chart-service/dto"

//TODO Add tests
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

		license.History = keepOnlyLastEntryOfEachDay(license.History)
	}

	return licenses, nil
}
