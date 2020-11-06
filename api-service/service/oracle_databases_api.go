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
	"strings"
	"time"

	"github.com/ercole-io/ercole/api-service/dto"
	"github.com/ercole-io/ercole/utils"
)

// SearchOracleDatabaseAddms search addms
func (as *APIService) SearchOracleDatabaseAddms(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchOracleDatabaseAddms(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// SearchOracleDatabaseSegmentAdvisors search segment advisors
func (as *APIService) SearchOracleDatabaseSegmentAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchOracleDatabaseSegmentAdvisors(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// SearchOracleDatabasePatchAdvisors search patch advisors
func (as *APIService) SearchOracleDatabasePatchAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchOracleDatabasePatchAdvisors(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, windowTime, location, environment, olderThan, status)
}

// SearchOracleDatabases search databases
func (as *APIService) SearchOracleDatabases(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchOracleDatabases(full, strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// SearchOracleDatabaseUsedLicenses return the list of used licenses
func (as *APIService) SearchOracleDatabaseUsedLicenses(sortBy string, sortDesc bool, page int, pageSize int,
	location string, environment string, olderThan time.Time,
) (*dto.OracleDatabaseUsedLicenseSearchResponse, utils.AdvancedErrorInterface) {
	return as.Database.SearchOracleDatabaseUsedLicenses(sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}
