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

package service

import (
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// GetSqlServerDatabaseLicenseTypes return the list of SqlServerDatabaseLicenseType
func (as *APIService) GetSqlServerDatabaseLicenseTypes() ([]model.SqlServerDatabaseLicenseType, error) {
	parts, err := as.Database.GetSqlServerDatabaseLicenseTypes()
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return parts, nil
}

// GetSqlServerDatabaseLicenseTypesAsMap return the list of SqlServerDatabaseLicenseType as map by ID
func (as *APIService) GetSqlServerDatabaseLicenseTypesAsMap() (map[string]model.SqlServerDatabaseLicenseType, error) {
	parts, err := as.GetSqlServerDatabaseLicenseTypes()
	if err != nil {
		return nil, err
	}

	partsMap := make(map[string]model.SqlServerDatabaseLicenseType)
	for _, part := range parts {
		partsMap[part.ID] = part
	}

	return partsMap, nil
}

func (as *APIService) GetSqlServerDatabaseLicenseType(id string) (*model.SqlServerDatabaseLicenseType, error) {
	licenseType, err := as.Database.GetSqlServerDatabaseLicenseType(id)
	if err != nil {
		return nil, err
	}

	return licenseType, nil
}
