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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/database"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/config"
)

//TODO Order as in routing?

// APIServiceInterface is a interface that wrap methods used to querying data
type APIServiceInterface interface {
	// Init initialize the service
	Init()
	// SearchHosts search hosts
	SearchHosts(mode string, filters dto.SearchHostsFilters) ([]map[string]interface{}, error)
	// SearchHostsAsLMS return LMS template file with the hosts filtered
	SearchHostsAsLMS(filters dto.SearchHostsFilters) (*excelize.File, error)
	// GetHost return the host specified in the hostname param
	GetHost(hostname string, olderThan time.Time, raw bool) (interface{}, error)
	// ListManagedTechnologies returns the list of technologies with some stats
	ListManagedTechnologies(sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]model.TechnologyStatus, error)
	// SearchAlerts search alerts
	SearchAlerts(mode string, search string, sortBy string, sortDesc bool, page, pageSize int, location, environment, severity, status string, from, to time.Time) ([]map[string]interface{}, error)
	// SearchClusters search clusters
	SearchClusters(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// GetCluster return the cluster specified in the clusterName param
	GetCluster(clusterName string, olderThan time.Time) (*dto.Cluster, error)
	// GetClusterXLSX return  cluster vms as xlxs file
	GetClusterXLSX(clusterName string, olderThan time.Time) (*excelize.File, error)
	// SearchOracleDatabaseAddms search addm
	SearchOracleDatabaseAddms(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// SearchOracleDatabaseSegmentAdvisors search segment advisors
	SearchOracleDatabaseSegmentAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// SearchOracleDatabasePatchAdvisors search patch advisors
	SearchOracleDatabasePatchAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) ([]map[string]interface{}, error)
	// SearchOracleDatabases search databases
	SearchOracleDatabases(filter dto.SearchOracleDatabasesFilter) ([]map[string]interface{}, error)
	// SearchOracleDatabases search databases
	SearchOracleDatabasesAsXLSX(filter dto.SearchOracleDatabasesFilter) (*excelize.File, error)
	// SearchOracleExadata search exadata
	SearchOracleExadata(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, error)
	// SearchOracleDatabaseUsedLicenses return the list of consumed licenses
	SearchOracleDatabaseUsedLicenses(sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleDatabaseUsedLicenseSearchResponse, error)

	// ListLocations list locations
	ListLocations(location string, environment string, olderThan time.Time) ([]string, error)
	// ListEnvironments list environments
	ListEnvironments(location string, environment string, olderThan time.Time) ([]string, error)
	// SearchOracleDatabaseLicenseModifiers search license modifiers
	SearchOracleDatabaseLicenseModifiers(search string, sortBy string, sortDesc bool, page int, pageSize int) ([]map[string]interface{}, error)

	// GetPatchingFunction return the patching function specified in the hostname param
	GetPatchingFunction(hostname string) (interface{}, error)

	// GetHostsCountStats return the number of the non-archived hosts
	GetHostsCountStats(location string, environment string, olderThan time.Time) (int, error)
	// GetEnvironmentStats return a array containing the number of hosts per environment
	GetEnvironmentStats(location string, olderThan time.Time) ([]interface{}, error)
	// GetOperatingSystemStats return a array containing the number of hosts per operating system
	GetOperatingSystemStats(location string, olderThan time.Time) ([]interface{}, error)
	// GetTypeStats return a array containing the number of hosts per type
	GetTypeStats(location string, olderThan time.Time) ([]interface{}, error)
	// GetTopUnusedOracleDatabaseInstanceResourceStats return a array containing top unused instance resource by workload
	GetTopUnusedOracleDatabaseInstanceResourceStats(location string, environment string, limit int, olderThan time.Time) ([]interface{}, error)
	// GetOracleDatabaseEnvironmentStats return a array containing the number of databases per environment
	GetOracleDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, error)
	// GetOracleDatabaseHighReliabilityStats return a array containing the number of databases per high-reliability status
	GetOracleDatabaseHighReliabilityStats(location string, environment string, olderThan time.Time) ([]interface{}, error)
	// GetOracleDatabaseVersionStats return a array containing the number of databases per version
	GetOracleDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, error)
	// GetTopReclaimableOracleDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
	GetTopReclaimableOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, error)
	// GetTotalTechnologiesComplianceStats return the total compliance of all technologie
	GetTotalTechnologiesComplianceStats(location string, environment string, olderThan time.Time) (map[string]interface{}, error)
	// GetOracleDatabasePatchStatusStats return a array containing the number of databases per patch status
	GetOracleDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, error)
	// GetTopWorkloadOracleDatabaseStats return a array containing top databases by workload
	GetTopWorkloadOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, error)
	// GetOracleDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
	GetOracleDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error)
	// GetOracleDatabaseRACStatusStats return a array containing the number of databases per RAC status
	GetOracleDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error)
	// GetOracleDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
	GetOracleDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error)
	// GetTotalOracleDatabaseWorkStats return the total work of databases
	GetTotalOracleDatabaseWorkStats(location string, environment string, olderThan time.Time) (float64, error)
	// GetTotalOracleDatabaseMemorySizeStats return the total of memory size of databases
	GetTotalOracleDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float64, error)
	// GetTotalOracleDatabaseDatafileSizeStats return the total size of datafiles of databases
	GetTotalOracleDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float64, error)
	// GetTotalOracleDatabaseSegmentSizeStats return the total size of segments of databases
	GetTotalOracleDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float64, error)
	// GetTotalOracleExadataMemorySizeStats return the total size of memory of exadata
	GetTotalOracleExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float64, error)
	// GetTotalOracleExadataCPUStats return the total cpu of exadata
	GetTotalOracleExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, error)
	// GetAverageOracleExadataStorageUsageStats return the average usage of cell disks of exadata
	GetAverageOracleExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float64, error)
	// GetOracleExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
	GetOracleExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error)
	// GetOracleExadataPatchStatusStats return a array containing the number of exadata per patch status
	GetOracleExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, error)
	// GetDefaultDatabaseTags return the default list of database tags from configuration
	GetDefaultDatabaseTags() ([]string, error)
	// GetErcoleFeatures return a map of active/inactive features
	GetErcoleFeatures() (map[string]bool, error)
	// GetErcoleFeatures return the list of technologies
	GetTechnologyList() ([]model.TechnologyInfo, error)

	// ORACLE DATABASE AGREEMENTS

	// Add associated part to OracleDatabaseAgreement or create a new one
	AddAssociatedLicenseTypeToOracleDbAgreement(request dto.AssociatedLicenseTypeInOracleDbAgreementRequest) (string, error)
	// Update associated part in OracleDatabaseAgreement
	UpdateAssociatedLicenseTypeOfOracleDbAgreement(request dto.AssociatedLicenseTypeInOracleDbAgreementRequest) error
	// Search OracleDatabase associated parts agreements
	SearchAssociatedLicenseTypesInOracleDatabaseAgreements(filters dto.SearchOracleDatabaseAgreementsFilter) ([]dto.OracleDatabaseAgreementFE, error)
	// Delete associated part from OracleDatabaseAgreement
	DeleteAssociatedLicenseTypeFromOracleDatabaseAgreement(associateLicenseTypeID primitive.ObjectID) error
	// Add an host to AssociatedLicenseType
	AddHostToAssociatedLicenseType(associateLicenseTypeID primitive.ObjectID, hostname string) error
	// Remove host from AssociatedLicenseType
	RemoveHostFromAssociatedLicenseType(associateLicenseTypeID primitive.ObjectID, hostname string) error

	// PARTS

	GetOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error)
	GetOracleDatabaseLicensesCompliance() ([]dto.OracleDatabaseLicenseUsage, error)

	// PATCHING FUNCTIONS
	// SetPatchingFunction set the patching function of a host
	SetPatchingFunction(hostname string, pf model.PatchingFunction) (interface{}, error)
	// DeletePatchingFunction delete the patching function of a host
	DeletePatchingFunction(hostname string) error

	// AddTagToOracleDatabase add the tag to the database if it hasn't the tag
	AddTagToOracleDatabase(hostname string, dbname string, tagname string) error
	// DeleteTagOfOracleDatabase delete the tag from the database if it hasn't the tag
	DeleteTagOfOracleDatabase(hostname string, dbname string, tagname string) error
	// SetOracleDatabaseLicenseModifier set the value of certain license to newValue
	SetOracleDatabaseLicenseModifier(hostname string, dbname string, licenseName string, newValue int) error
	// DeleteOracleDatabaseLicenseModifier delete the modifier of a certain license
	DeleteOracleDatabaseLicenseModifier(hostname string, dbname string, licenseName string) error
	// AckAlerts ack the specified alerts
	AckAlerts(ids []primitive.ObjectID) error
	// ArchiveHost archive the specified host
	ArchiveHost(hostname string) error

	// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
	GetInfoForFrontendDashboard(location string, environment string, olderThan time.Time) (map[string]interface{}, error)

	SearchDatabases(filter dto.GlobalFilter) ([]dto.Database, error)
	SearchDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	GetDatabasesStatistics(filter dto.GlobalFilter) (*dto.DatabasesStatistics, error)

	SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error)
	SearchMySQLInstancesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
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
	// Log contains logger formatted
	Log *logrus.Logger
	// TechnologyInfos contains the list of technologies with their informations
	TechnologyInfos []model.TechnologyInfo
	// NewObjectID return a new ObjectID
	NewObjectID func() primitive.ObjectID
}

// Init initializes the service and database
func (as *APIService) Init() {
	as.loadManagedTechnologiesList()

	as.NewObjectID = func() primitive.ObjectID {
		return primitive.NewObjectIDFromTimestamp(as.TimeNow())
	}
}
