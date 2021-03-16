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
)

// GetOracleDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
func (as *APIService) GetOracleDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleDatabaseArchivelogStatusStats(location, environment, olderThan)
}

// GetOracleDatabaseEnvironmentStats return a array containing the number of databases per environment
func (as *APIService) GetOracleDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleDatabaseEnvironmentStats(location, olderThan)
}

// GetOracleDatabaseHighReliabilityStats return a array containing the number of databases per high-reliability status
func (as *APIService) GetOracleDatabaseHighReliabilityStats(location string, environment string, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleDatabaseHighReliabilityStats(location, environment, olderThan)
}

// GetOracleDatabaseVersionStats return a array containing the number of databases per version
func (as *APIService) GetOracleDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleDatabaseVersionStats(location, olderThan)
}

// GetTopReclaimableOracleDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
func (as *APIService) GetTopReclaimableOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetTopReclaimableOracleDatabaseStats(location, limit, olderThan)
}

// GetOracleDatabasePatchStatusStats return a array containing the number of databases per patch status
func (as *APIService) GetOracleDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleDatabasePatchStatusStats(location, windowTime, olderThan)
}

// GetTopWorkloadOracleDatabaseStats return a array containing top databases by workload
func (as *APIService) GetTopWorkloadOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetTopWorkloadOracleDatabaseStats(location, limit, olderThan)
}

// GetOracleDatabaseRACStatusStats return a array containing the number of databases per RAC status
func (as *APIService) GetOracleDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleDatabaseRACStatusStats(location, environment, olderThan)
}

// GetOracleDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
func (as *APIService) GetOracleDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleDatabaseDataguardStatusStats(location, environment, olderThan)
}

// GetTotalOracleDatabaseWorkStats return the total work of databases
func (as *APIService) GetTotalOracleDatabaseWorkStats(location string, environment string, olderThan time.Time) (float64, error) {
	return as.Database.GetTotalOracleDatabaseWorkStats(location, environment, olderThan)
}

// GetTotalOracleDatabaseMemorySizeStats return the total of memory size of databases
func (as *APIService) GetTotalOracleDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float64, error) {
	return as.Database.GetTotalOracleDatabaseMemorySizeStats(location, environment, olderThan)
}

// GetTotalOracleDatabaseDatafileSizeStats return the total size of datafiles of databases
func (as *APIService) GetTotalOracleDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float64, error) {
	return as.Database.GetTotalOracleDatabaseDatafileSizeStats(location, environment, olderThan)
}

// GetTotalOracleDatabaseSegmentSizeStats return the total size of segments of databases
func (as *APIService) GetTotalOracleDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float64, error) {
	return as.Database.GetTotalOracleDatabaseSegmentSizeStats(location, environment, olderThan)
}

// GetTopUnusedOracleDatabaseInstanceResourceStats return a array containing top unused instance resource by workload
func (as *APIService) GetTopUnusedOracleDatabaseInstanceResourceStats(location string, environment string, limit int, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetTopUnusedOracleDatabaseInstanceResourceStats(location, environment, limit, olderThan)
}
