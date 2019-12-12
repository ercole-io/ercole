// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/amreo/ercole-services/utils"
)

// GetDatabaseEnvironmentStats return a array containing the number of databases per environment
func (as *APIService) GetDatabaseEnvironmentStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetDatabaseEnvironmentStats(location)
}

// GetDatabaseVersionStats return a array containing the number of databases per version
func (as *APIService) GetDatabaseVersionStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetDatabaseVersionStats(location)
}

// GetTopReclaimableDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
func (as *APIService) GetTopReclaimableDatabaseStats(location string, limit int) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetTopReclaimableDatabaseStats(location, limit)
}

// GetPatchStatusDatabaseStats return a array containing the number of databases per patch status
func (as *APIService) GetPatchStatusDatabaseStats(location string, windowTime time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetPatchStatusDatabaseStats(location, windowTime)
}

// GetTopWorkloadDatabaseStats return a array containing top databases by workload
func (as *APIService) GetTopWorkloadDatabaseStats(location string, limit int) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetTopWorkloadDatabaseStats(location, limit)
}
