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

	"github.com/ercole-io/ercole/v2/api-service/dto"
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

func (as *APIService) GetOracleDatabasesStatistics(filter dto.GlobalFilter) (*dto.OracleDatabasesStatistics, error) {
	stats := new(dto.OracleDatabasesStatistics)
	var err error

	stats.TotalMemorySize, err = as.Database.GetTotalOracleDatabaseMemorySizeStats(filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}
	stats.TotalSegmentsSize, err = as.Database.GetTotalOracleDatabaseSegmentSizeStats(filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}
	stats.TotalDatafileSize, err = as.Database.GetTotalOracleDatabaseDatafileSizeStats(filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}
	stats.TotalWork, err = as.Database.GetTotalOracleDatabaseWorkStats(filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetTopUnusedOracleDatabaseInstanceResourceStats return a array containing top unused instance resource by workload
func (as *APIService) GetTopUnusedOracleDatabaseInstanceResourceStats(location string, environment string, limit int, olderThan time.Time) ([]interface{}, error) {
	return as.Database.GetTopUnusedOracleDatabaseInstanceResourceStats(location, environment, limit, olderThan)
}
