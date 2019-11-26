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

	"github.com/amreo/ercole-services/api-service/database"
	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/config"
)

// APIServiceInterface is a interface that wrap methods used to querying data
type APIServiceInterface interface {
	// Init initialize the service
	Init()
	// SearchCurrentHosts search current hosts
	SearchCurrentHosts(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int) ([]interface{}, utils.AdvancedErrorInterface)
	// GetCurrentHost return the current host specified in the hostname param
	GetCurrentHost(hostname string) (interface{}, utils.AdvancedErrorInterface)
	// SearchAlerts search alerts
	SearchAlerts(search string, sortBy string, sortDesc bool, page int, pageSize int) ([]interface{}, utils.AdvancedErrorInterface)
}

// APIService is the concrete implementation of APIServiceInterface.
type APIService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Version of the saved data
	Version string
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
}

// Init initializes the service and database
func (as *APIService) Init() {
}
