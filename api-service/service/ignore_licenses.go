// Copyright (c) 2025 Sorint.lab S.p.A.
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
package service

import (
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
)

func (as *APIService) IgnoreLicenses(licenses []dto.IgnoreLicenseRequest) *dto.IgnoreLicenseResponse {
	result := &dto.IgnoreLicenseResponse{}

	for _, license := range licenses {
		switch license.Technology {
		case model.TechnologyOracleDatabase:
			err := as.Database.UpdateLicenseIgnoredField(license.Hostname, license.DatabaseName, license.LicenseTypeID, license.Ignored, license.IgnoredComment)
			if err != nil {
				result.Error = append(result.Error, license)
				continue
			}

			result.Updated = append(result.Updated, license)
		case model.TechnologyOracleMySQL:
			err := as.Database.UpdateMySqlLicenseIgnoredField(license.Hostname, license.DatabaseName, license.Ignored, license.IgnoredComment)
			if err != nil {
				result.Error = append(result.Error, license)
				continue
			}

			result.Updated = append(result.Updated, license)
		case model.TechnologyMicrosoftSQLServer:
			err := as.Database.UpdateSqlServerLicenseIgnoredField(license.Hostname, license.DatabaseName, license.Ignored, license.IgnoredComment)
			if err != nil {
				result.Error = append(result.Error, license)
				continue
			}

			result.Updated = append(result.Updated, license)
		default:
			result.Error = append(result.Error, license)
		}
	}

	return result
}
