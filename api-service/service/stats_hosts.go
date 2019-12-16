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
	"github.com/amreo/ercole-services/utils"
)

// GetEnvironmentStats return a array containing the number of hosts per environment
func (as *APIService) GetEnvironmentStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetEnvironmentStats(location)
}

// GetTypeStats return a array containing the number of hosts per type
func (as *APIService) GetTypeStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetTypeStats(location)
}

// GetOperatingSystemStats return a array containing the number of hosts per operating system
func (as *APIService) GetOperatingSystemStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetOperatingSystemStats(location)
}

// GetDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
func (as *APIService) GetDatabaseArchivelogStatusStats(location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.GetDatabaseArchivelogStatusStats(location, environment)
}
