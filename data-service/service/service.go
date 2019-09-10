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

// Service is a package that provides methods for manipulating host informations
package service

import (
	"errors"
	"time"

	"github.com/amreo/ercole-services/data-service/database"

	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
)

// ErrHostNotFound cotains "Host not found" error
var ErrHostNotFound = errors.New("Host not found")

// HostDataServiceInterface is a interface that wrap methods used to manipulate and save data
type HostDataServiceInterface interface {
	// Init initialize the service
	Init()

	// UpdateHostInfo update the host informations using the provided hostdata
	UpdateHostInfo(hostdata model.HostData) (interface{}, utils.AdvancedErrorInterface)

	// ArchiveHost archive the host
	// ArchiveHost(hostname string) utils.AdvancedError
}

// HostDataService is the concrete implementation of HostDataServiceInterface. It saves data to a MongoDB database
type HostDataService struct {
	// Config contains the dataservice global configuration
	// TODO: Should be removed?
	Config config.Configuration
	// Version of the saved data
	Version string
	// Database contains the database layer
	Database database.MongoDatabaseInterface
}

// Init initializes the service and database
func (hds *HostDataService) Init() {
}

// UpdateHostInfo saves the hostdata
func (hds *HostDataService) UpdateHostInfo(hostdata model.HostData) (interface{}, utils.AdvancedErrorInterface) {
	hostdata.ServerVersion = hds.Version
	hostdata.Archived = false
	hostdata.CreatedAt = time.Now()
	hostdata.SchemaVersion = model.SchemaVersion

	//Archive the host
	_, err := hds.Database.ArchiveHost(hostdata.Hostname)
	if err != nil {
		return nil, err
	}

	//Insert the host
	res, err := hds.Database.InsertHostData(hostdata)
	if err != nil {
		return nil, err
	} else {
		return res.InsertedID, nil
	}

}
