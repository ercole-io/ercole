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

// GetDatabasePatchStatusStats return a array containing the number of databases per patch status
func (as *APIService) GetDatabasePatchStatusStats(location string, windowTime time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetDatabasePatchStatusStats(location, windowTime)
}

// GetTopWorkloadDatabaseStats return a array containing top databases by workload
func (as *APIService) GetTopWorkloadDatabaseStats(location string, limit int) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetTopWorkloadDatabaseStats(location, limit)
}

// GetDatabaseRACStatusStats return a array containing the number of databases per RAC status
func (as *APIService) GetDatabaseRACStatusStats(location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetDatabaseRACStatusStats(location, environment)
}

// GetDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
func (as *APIService) GetDatabaseDataguardStatusStats(location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetDatabaseDataguardStatusStats(location, environment)
}

// GetTotalDatabaseWorkStats return the total work of databases
func (as *APIService) GetTotalDatabaseWorkStats(location string, environment string) (float32, utils.AdvancedErrorInterface) {
	return as.Database.GetTotalDatabaseWorkStats(location, environment)
}

// GetTotalDatabaseMemorySizeStats return the total of memory size of databases
func (as *APIService) GetTotalDatabaseMemorySizeStats(location string, environment string) (float32, utils.AdvancedErrorInterface) {
	return as.Database.GetTotalDatabaseMemorySizeStats(location, environment)
}

// GetTotalDatabaseDatafileSizeStats return the total size of datafiles of databases
func (as *APIService) GetTotalDatabaseDatafileSizeStats(location string, environment string) (float32, utils.AdvancedErrorInterface) {
	return as.Database.GetTotalDatabaseDatafileSizeStats(location, environment)
}

// GetTotalDatabaseSegmentSizeStats return the total size of segments of databases
func (as *APIService) GetTotalDatabaseSegmentSizeStats(location string, environment string) (float32, utils.AdvancedErrorInterface) {
	return as.Database.GetTotalDatabaseSegmentSizeStats(location, environment)
}

// GetDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
func (as *APIService) GetDatabaseLicenseComplianceStatusStats(location string, environment string) (interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetDatabaseLicenseComplianceStatusStats(location, environment)
}
