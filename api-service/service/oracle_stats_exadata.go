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

// GetTotalOracleExadataMemorySizeStats return the total size of memory of exadata
func (as *APIService) GetTotalOracleExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float64, error) {
	return as.Database.GetTotalOracleExadataMemorySizeStats(location, environment, olderThan)
}

// GetTotalOracleExadataCPUStats return the total cpu of exadata
func (as *APIService) GetTotalOracleExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, error) {
	return as.Database.GetTotalOracleExadataCPUStats(location, environment, olderThan)
}

// GetAverageOracleExadataStorageUsageStats return the average usage of cell disks of exadata
func (as *APIService) GetAverageOracleExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float64, error) {
	return as.Database.GetAverageOracleExadataStorageUsageStats(location, environment, olderThan)
}

// GetOracleExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
func (as *APIService) GetOracleExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleExadataStorageErrorCountStatusStats(location, environment, olderThan)
}

// GetOracleExadataPatchStatusStats return a array containing the number of exadata per patch status
func (as *APIService) GetOracleExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetOracleExadataPatchStatusStats(location, environment, windowTime, olderThan)
}
