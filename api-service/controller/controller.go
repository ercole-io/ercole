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

package controller

import (
	"net/http"
	"time"

	"github.com/amreo/ercole-services/api-service/auth"
	"github.com/amreo/ercole-services/api-service/service"
	"github.com/amreo/ercole-services/config"
	"github.com/goji/httpauth"
	"github.com/sirupsen/logrus"
)

// APIControllerInterface is a interface that wrap methods used to querying data
type APIControllerInterface interface {
	// AuthenticateMiddleware return the middleware used to authenticate users
	AuthenticateMiddleware() func(http.Handler) http.Handler
	// SearchHosts search hosts data using the filters in the request
	SearchHosts(w http.ResponseWriter, r *http.Request)
	// SearchDatabases search databases data using the filters in the request
	SearchDatabases(w http.ResponseWriter, r *http.Request)
	// SearchClusters search clusters data using the filters in the request
	SearchClusters(w http.ResponseWriter, r *http.Request)
	// SearchAddms search addms data using the filters in the request
	SearchAddms(w http.ResponseWriter, r *http.Request)
	// SearchSegmentAdvisors search segment advisors data using the filters in the request
	SearchSegmentAdvisors(w http.ResponseWriter, r *http.Request)
	// SearchPatchAdvisors search patch advisors data using the filters in the request
	SearchPatchAdvisors(w http.ResponseWriter, r *http.Request)
	// GetHost return all'informations about the host requested in the id path variable
	GetHost(w http.ResponseWriter, r *http.Request)
	// SearchAlerts search alerts using the filters in the request
	SearchAlerts(w http.ResponseWriter, r *http.Request)
	// SearchExadata search exadata data using the filters in the request
	SearchExadata(w http.ResponseWriter, r *http.Request)
	// ListLicenses list licenses using the filters in the request
	ListLicenses(w http.ResponseWriter, r *http.Request)
	// GetPatchingFunction return all'informations about the patching function of the host requested in the hostnmae path variable
	GetPatchingFunction(w http.ResponseWriter, r *http.Request)
	// ListLocations list locations using the filters in the request
	ListLocations(w http.ResponseWriter, r *http.Request)
	// ListEnvironments list environments using the filters in the request
	ListEnvironments(w http.ResponseWriter, r *http.Request)

	// GetEnvironmentStats return all statistics about the environments of the hosts using the filters in the request
	GetEnvironmentStats(w http.ResponseWriter, r *http.Request)
	// GetTypeStats return all statistics about the types of the hosts using the filters in the request
	GetTypeStats(w http.ResponseWriter, r *http.Request)
	// GetOperatingSystemStats return all statistics about the operating systems of the hosts using the filters in the request
	GetOperatingSystemStats(w http.ResponseWriter, r *http.Request)
	// GetTopUnusedInstanceResourceStats return top unused instance resource by databases work using the filters in the request
	GetTopUnusedInstanceResourceStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseEnvironmentStats return all statistics about the environments of the databases using the filters in the request
	GetDatabaseEnvironmentStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseVersionStats return all statistics about the versions of the databases using the filters in the request
	GetDatabaseVersionStats(w http.ResponseWriter, r *http.Request)
	// GetTopReclaimableDatabaseStats return top databases by reclaimable segment advisors using the filters in the request
	GetTopReclaimableDatabaseStats(w http.ResponseWriter, r *http.Request)
	// GetDatabasePatchStatusStats return all statistics about the patch status of the databases using the filters in the request
	GetDatabasePatchStatusStats(w http.ResponseWriter, r *http.Request)
	// GetTopWorkloadDatabaseStats return top databases by workload advisors using the filters in the request
	GetTopWorkloadDatabaseStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseDataguardStatusStats return all statistics about the dataguard status of the databases using the filters in the request
	GetDatabaseDataguardStatusStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseRACStatusStats return all statistics about the RAC status of the databases using the filters in the request
	GetDatabaseRACStatusStats(w http.ResponseWriter, r *http.Request)
	// GetDatabasArchivelogStatusStats return all statistics about the archivelog status of the databases using the filters in the request
	GetDatabaseArchivelogStatusStats(w http.ResponseWriter, r *http.Request)
	// GetTotalDatabaseWorkStats return the total work of databases using the filters in the request
	GetTotalDatabaseWorkStats(w http.ResponseWriter, r *http.Request)
	// GetTotalDatabaseMemorySizeStats return the total size of memory of databases using the filters in the request
	GetTotalDatabaseMemorySizeStats(w http.ResponseWriter, r *http.Request)
	// GetTotalDatabaseDatafileSizeStats return the total size of datafiles of databases using the filters in the request
	GetTotalDatabaseDatafileSizeStats(w http.ResponseWriter, r *http.Request)
	// GetTotalDatabaseSegmentSizeStats return the total size of segments of databases using the filters in the request
	GetTotalDatabaseSegmentSizeStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases using the filters in the request
	GetDatabaseLicenseComplianceStatusStats(w http.ResponseWriter, r *http.Request)
	// GetTotalExadataMemorySizeStats return the total size of memory of exadata using the filters in the request
	GetTotalExadataMemorySizeStats(w http.ResponseWriter, r *http.Request)
	// GetTotalExadataCPUStats return the total cpu of exadata using the filters in the request
	GetTotalExadataCPUStats(w http.ResponseWriter, r *http.Request)
	// GetAverageExadataStorageUsageStats return the average usage of cell disks of exadata using the filters in the request
	GetAverageExadataStorageUsageStats(w http.ResponseWriter, r *http.Request)
	// GetExadataStorageErrorCountStatusStats return all statistics about the ErrorCount status of the storage of the exadata using the filters in the request
	GetExadataStorageErrorCountStatusStats(w http.ResponseWriter, r *http.Request)
	// GetExadataPatchStatusStats return all statistics about the patch status of the exadata using the filters in the request
	GetExadataPatchStatusStats(w http.ResponseWriter, r *http.Request)
	//GetDefaultDatabaseTags return the default list of database tags from configuration
	GetDefaultDatabaseTags(w http.ResponseWriter, r *http.Request)

	// SetLicenseCount set the count of a certain license
	SetLicenseCount(w http.ResponseWriter, r *http.Request)
	// SetPatchingFunction set the patching function of a host specified in the hostname path variable to the content of the request body
	SetPatchingFunction(w http.ResponseWriter, r *http.Request)
	// AddTagToDatabase add a tag to the database if it hasn't the tag
	AddTagToDatabase(w http.ResponseWriter, r *http.Request)
	// DeleteTagOfDatabase remove a certain tag from a database if it has the tag
	DeleteTagOfDatabase(w http.ResponseWriter, r *http.Request)
	// SetLicenseModifier set the license modifier of specified license/db/host in the request to the value in the body
	SetLicenseModifier(w http.ResponseWriter, r *http.Request)
	// DeleteLicenseModifier delete the license modifier of specified license/db/host in the request
	DeleteLicenseModifier(w http.ResponseWriter, r *http.Request)
	// AckAlert ack the specified alert in the request
	AckAlert(w http.ResponseWriter, r *http.Request)
	// ArchiveHost archive the specified host in the request
	ArchiveHost(w http.ResponseWriter, r *http.Request)
}

// APIController is the struct used to handle the requests from agents and contains the concrete implementation of APIControllerInterface
type APIController struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Service contains the underlying service used to perform various logical and store operations
	Service service.APIServiceInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log *logrus.Logger
	// Authenticator contains the authenticator
	Authenticator auth.AuthenticationProvider
}

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (ctrl *APIController) AuthenticateMiddleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(ctrl.Config.APIService.AuthenticationProvider.Username, ctrl.Config.APIService.AuthenticationProvider.Password)
}
