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

package controller

import (
	"net/http"
	"time"

	"github.com/ercole-io/ercole/api-service/auth"
	"github.com/ercole-io/ercole/api-service/service"
	"github.com/ercole-io/ercole/config"
	"github.com/sirupsen/logrus"
)

// APIControllerInterface is a interface that wrap methods used to querying data
type APIControllerInterface interface {
	// SearchHosts search hosts data using the filters in the request
	SearchHosts(w http.ResponseWriter, r *http.Request)
	// ListTechnologies returns the list of Technologies with some stats using the filters in the request
	ListTechnologies(w http.ResponseWriter, r *http.Request)
	// SearchOracleDatabases search databases data using the filters in the request
	SearchOracleDatabases(w http.ResponseWriter, r *http.Request)
	// SearchClusters search clusters data using the filters in the request
	SearchClusters(w http.ResponseWriter, r *http.Request)
	// GetCluster get cluster data using the filters in the request
	GetCluster(w http.ResponseWriter, r *http.Request)
	// SearchOracleDatabaseAddms search addms data using the filters in the request
	SearchOracleDatabaseAddms(w http.ResponseWriter, r *http.Request)
	// SearchOracleDatabaseSegmentAdvisors search segment advisors data using the filters in the request
	SearchOracleDatabaseSegmentAdvisors(w http.ResponseWriter, r *http.Request)
	// SearchOracleDatabasePatchAdvisors search patch advisors data using the filters in the request
	SearchOracleDatabasePatchAdvisors(w http.ResponseWriter, r *http.Request)
	// GetHost return all informations about the host requested in the id path variable
	GetHost(w http.ResponseWriter, r *http.Request)
	// SearchAlerts search alerts using the filters in the request
	SearchAlerts(w http.ResponseWriter, r *http.Request)
	// SearchOracleExadata search exadata data using the filters in the request
	SearchOracleExadata(w http.ResponseWriter, r *http.Request)
	// SearchLicenses search licenses using the filters in the request
	SearchLicenses(w http.ResponseWriter, r *http.Request)
	// SearchOracleDatabaseConsumedLicenses search licenses consumed by the hosts using the filters in the request
	SearchOracleDatabaseConsumedLicenses(w http.ResponseWriter, r *http.Request)
	// GetLicense return a certain license asked in the request
	GetLicense(w http.ResponseWriter, r *http.Request)
	// SearchOracleDatabaseLicenseModifiers search a license modifier using the filters in the request
	SearchOracleDatabaseLicenseModifiers(w http.ResponseWriter, r *http.Request)
	// SearchOracleDatabaseAgreements search Oracle/Database agreements data using the filters in the request
	SearchOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request)

	// GetPatchingFunction return all'informations about the patching function of the host requested in the hostnmae path variable
	GetPatchingFunction(w http.ResponseWriter, r *http.Request)
	// ListLocations list locations using the filters in the request
	ListLocations(w http.ResponseWriter, r *http.Request)
	// ListEnvironments list environments using the filters in the request
	ListEnvironments(w http.ResponseWriter, r *http.Request)
	// GetHostsCountStats return the number of the hosts using the filters in the request
	GetHostsCountStats(w http.ResponseWriter, r *http.Request)
	// GetEnvironmentStats return all statistics about the environments of the hosts using the filters in the request
	GetEnvironmentStats(w http.ResponseWriter, r *http.Request)
	// GetTypeStats return all statistics about the types of the hosts using the filters in the request
	GetTypeStats(w http.ResponseWriter, r *http.Request)
	// GetOperatingSystemStats return all statistics about the operating systems of the hosts using the filters in the request
	GetOperatingSystemStats(w http.ResponseWriter, r *http.Request)
	// GetTopOracleDatabaseUnusedInstanceResourceStats return top unused instance resource by databases work using the filters in the request
	GetTopUnusedOracleDatabaseInstanceResourceStats(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabaseEnvironmentStats return all statistics about the environments of the databases using the filters in the request
	GetOracleDatabaseEnvironmentStats(w http.ResponseWriter, r *http.Request)
	// GetDatabaseHighReliabilityStats return all statistics about the high-reliability status of the databases using the filters in the request
	GetOracleDatabaseHighReliabilityStats(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabaseVersionStats return all statistics about the versions of the databases using the filters in the request
	GetOracleDatabaseVersionStats(w http.ResponseWriter, r *http.Request)
	// GetTopReclaimableOracleDatabaseStats return top databases by reclaimable segment advisors using the filters in the request
	GetTopReclaimableOracleDatabaseStats(w http.ResponseWriter, r *http.Request)
	// GetTotalTechnologiesComplianceStats return the total compliance of all technologies using the filters in the request
	GetTotalTechnologiesComplianceStats(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabasePatchStatusStats return all statistics about the patch status of the databases using the filters in the request
	GetOracleDatabasePatchStatusStats(w http.ResponseWriter, r *http.Request)
	// GetTopWorkloadOracleDatabaseStats return top databases by workload advisors using the filters in the request
	GetTopWorkloadOracleDatabaseStats(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabaseDataguardStatusStats return all statistics about the dataguard status of the databases using the filters in the request
	GetOracleDatabaseDataguardStatusStats(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabaseRACStatusStats return all statistics about the RAC status of the databases using the filters in the request
	GetOracleDatabaseRACStatusStats(w http.ResponseWriter, r *http.Request)
	// GetDatabasArchivelogStatusStats return all statistics about the archivelog status of the databases using the filters in the request
	GetOracleDatabaseArchivelogStatusStats(w http.ResponseWriter, r *http.Request)
	// GetTotalOracleDatabaseWorkStats return the total work of databases using the filters in the request
	GetTotalOracleDatabaseWorkStats(w http.ResponseWriter, r *http.Request)
	// GetTotalOracleDatabaseMemorySizeStats return the total size of memory of databases using the filters in the request
	GetTotalOracleDatabaseMemorySizeStats(w http.ResponseWriter, r *http.Request)
	// GetTotalOracleDatabaseDatafileSizeStats return the total size of datafiles of databases using the filters in the request
	GetTotalOracleDatabaseDatafileSizeStats(w http.ResponseWriter, r *http.Request)
	// GetTotalOracleDatabaseSegmentSizeStats return the total size of segments of databases using the filters in the request
	GetTotalOracleDatabaseSegmentSizeStats(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases using the filters in the request
	GetOracleDatabaseLicenseComplianceStatusStats(w http.ResponseWriter, r *http.Request)
	// GetTotalOracleExadataMemorySizeStats return the total size of memory of exadata using the filters in the request
	GetTotalOracleExadataMemorySizeStats(w http.ResponseWriter, r *http.Request)
	// GetTotalOracleExadataCPUStats return the total cpu of exadata using the filters in the request
	GetTotalOracleExadataCPUStats(w http.ResponseWriter, r *http.Request)
	// GetAverageOracleExadataStorageUsageStats return the average usage of cell disks of exadata using the filters in the request
	GetAverageOracleExadataStorageUsageStats(w http.ResponseWriter, r *http.Request)
	// GetOracleExadataStorageErrorCountStatusStats return all statistics about the ErrorCount status of the storage of the exadata using the filters in the request
	GetOracleExadataStorageErrorCountStatusStats(w http.ResponseWriter, r *http.Request)
	// GetOracleExadataPatchStatusStats return all statistics about the patch status of the exadata using the filters in the request
	GetOracleExadataPatchStatusStats(w http.ResponseWriter, r *http.Request)

	// GetDefaultDatabaseTags return the default list of database tags from configuration
	GetDefaultDatabaseTags(w http.ResponseWriter, r *http.Request)
	// GetErcoleFeatures return a map of active/inactive features
	GetErcoleFeatures(w http.ResponseWriter, r *http.Request)
	// GetTechnologyList return the list of techonlogies
	GetTechnologyList(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
	GetOracleDatabaseAgreementPartsList(w http.ResponseWriter, r *http.Request)

	// AddOracleDatabaseAgreements add some agreements
	AddOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request)
	// DeleteOracleDatabaseAgreement delete an agreement
	DeleteOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request)
	// AddAssociatedHostToOracleDatabaseAgreement add a associated host to an agreement
	AddAssociatedHostToOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request)
	// RemoveAssociatedHostToOracleDatabaseAgreement remove a associated host of an agreement
	RemoveAssociatedHostToOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request)

	// SetLicenseCostPerProcessor set the cost per processor of a certain license
	SetLicenseCostPerProcessor(w http.ResponseWriter, r *http.Request)

	// SetPatchingFunction set the patching function of a host specified in the hostname path variable to the content of the request body
	SetPatchingFunction(w http.ResponseWriter, r *http.Request)
	// DeletePatchingFunction remove the patching function of a host specified in the hostname path variable
	DeletePatchingFunction(w http.ResponseWriter, r *http.Request)
	// AddTagToOracleDatabase add a tag to the database if it hasn't the tag
	AddTagToOracleDatabase(w http.ResponseWriter, r *http.Request)
	// DeleteTagOfOracleDatabase remove a certain tag from a database if it has the tag
	DeleteTagOfOracleDatabase(w http.ResponseWriter, r *http.Request)
	// SetOracleDatabaseLicenseModifier set the license modifier of specified license/db/host in the request to the value in the body
	SetOracleDatabaseLicenseModifier(w http.ResponseWriter, r *http.Request)
	// DeleteOracleDatabaseLicenseModifier delete the license modifier of specified license/db/host in the request
	DeleteOracleDatabaseLicenseModifier(w http.ResponseWriter, r *http.Request)
	// AckAlerts ack the specified alert in the request
	AckAlerts(w http.ResponseWriter, r *http.Request)
	// ArchiveHost archive the specified host in the request
	ArchiveHost(w http.ResponseWriter, r *http.Request)

	// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
	GetInfoForFrontendDashboard(w http.ResponseWriter, r *http.Request)
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
