// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"go.mongodb.org/mongo-driver/bson/primitive"

	alertServiceClient "github.com/ercole-io/ercole/v2/alert-service/client"
	"github.com/ercole-io/ercole/v2/api-service/database"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
)

// APIServiceInterface is a interface that wrap methods used to querying data
type APIServiceInterface interface {
	// Init initialize the service
	Init()
	// SearchHosts search hosts
	SearchHosts(mode string, filters dto.SearchHostsFilters) ([]map[string]interface{}, error)
	// SearchHostsAsLMS return LMS template file with the hosts filtered
	SearchHostsAsLMS(filters dto.SearchHostsAsLMS) (*excelize.File, error)
	SearchHostsAsXLSX(filters dto.SearchHostsFilters) (*excelize.File, error)
	GetHostDataSummaries(filters dto.SearchHostsFilters) ([]dto.HostDataSummary, error)
	// GetHost return the host specified in the hostname param
	GetHost(hostname string, olderThan time.Time, raw bool) (*dto.HostData, error)
	// ListManagedTechnologies returns the list of technologies with some stats
	ListManagedTechnologies(sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]model.TechnologyStatus, error)
	// SearchAlerts search alerts
	SearchAlerts(mode string, search string, sortBy string, sortDesc bool, page, pageSize int, location, environment, severity, status string, from, to time.Time) ([]map[string]interface{}, error)
	SearchAlertsAsXLSX(from time.Time, to time.Time, filter dto.GlobalFilter) (*excelize.File, error)
	// SearchClusters search clusters
	SearchClusters(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	SearchClustersAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	// GetCluster return the cluster specified in the clusterName param
	GetCluster(clusterName string, olderThan time.Time) (*dto.Cluster, error)
	// GetClusterXLSX return  cluster vms as xlxs file
	GetClusterXLSX(clusterName string, olderThan time.Time) (*excelize.File, error)
	// SearchOracleDatabaseAddms search addm
	SearchOracleDatabaseAddms(search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// SearchOracleDatabaseSegmentAdvisors search segment advisors
	SearchOracleDatabaseSegmentAdvisors(search string, sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]dto.OracleDatabaseSegmentAdvisor, error)
	SearchOracleDatabaseSegmentAdvisorsAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	// SearchOracleDatabasePatchAdvisors search patch advisors
	SearchOracleDatabasePatchAdvisors(search string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) (*dto.PatchAdvisorResponse, error)
	SearchOracleDatabasePatchAdvisorsAsXLSX(windowTime time.Time, filter dto.GlobalFilter) (*excelize.File, error)
	// SearchOracleDatabases search databases
	SearchOracleDatabases(filter dto.SearchOracleDatabasesFilter) (*dto.OracleDatabaseResponse, error)
	// SearchOracleDatabases search databases
	SearchOracleDatabasesAsXLSX(filter dto.SearchOracleDatabasesFilter) (*excelize.File, error)
	// SearchOracleExadata search exadata
	SearchOracleExadata(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleExadataResponse, error)
	SearchOracleExadataAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	// SearchOracleDatabaseUsedLicenses return the list of consumed licenses
	SearchOracleDatabaseUsedLicenses(hostname string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleDatabaseUsedLicenseSearchResponse, error)

	// ListLocations list locations
	ListLocations(location string, environment string, olderThan time.Time) ([]string, error)
	// ListEnvironments list environments
	ListEnvironments(location string, environment string, olderThan time.Time) ([]string, error)

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
	GetOracleDatabasesStatistics(filter dto.GlobalFilter) (*dto.OracleDatabasesStatistics, error)

	// ORACLE DATABASE CONTRACTS

	AddOracleDatabaseContract(contract model.OracleDatabaseContract) (*dto.OracleDatabaseContractFE, error)
	UpdateOracleDatabaseContract(contract model.OracleDatabaseContract) (*dto.OracleDatabaseContractFE, error)
	GetOracleDatabaseContracts(filter dto.GetOracleDatabaseContractsFilter) ([]dto.OracleDatabaseContractFE, error)
	GetOracleDatabaseContractsAsXLSX(filter dto.GetOracleDatabaseContractsFilter) (*excelize.File, error)
	DeleteOracleDatabaseContract(id primitive.ObjectID) error
	AddHostToOracleDatabaseContract(id primitive.ObjectID, hostname string) error
	DeleteHostFromOracleDatabaseContract(id primitive.ObjectID, hostname string) error
	DeleteHostFromOracleDatabaseContracts(hostname string) error

	// ORACLE DATABASE LICENSES

	GetOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error)
	GetOracleDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error)
	DeleteOracleDatabaseLicenseType(id string) error
	AddOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) (*model.OracleDatabaseLicenseType, error)
	UpdateOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) (*model.OracleDatabaseLicenseType, error)

	ListOracleGrantDbaByHostname(hostname string, filter dto.GlobalFilter) ([]dto.OracleGrantDbaDto, error)

	// SQL SERVER DATABASE LICENSES
	GetSqlServerDatabaseLicenseTypes() ([]model.SqlServerDatabaseLicenseType, error)

	// MySQL DATABASE LICENSES
	GetMySqlLicenseTypes() ([]model.MySqlLicenseType, error)

	// SearchSqlServerInstances search databases
	SearchSqlServerInstances(filter dto.SearchSqlServerInstancesFilter) (*dto.SqlServerInstanceResponse, error)
	// SearchOracleDatabases search databases
	SearchSqlServerInstancesAsXLSX(filter dto.SearchSqlServerInstancesFilter) (*excelize.File, error)

	// SQL SERVER DATABASE CONTRACTS
	AddSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) (*model.SqlServerDatabaseContract, error)
	GetSqlServerDatabaseContracts() ([]model.SqlServerDatabaseContract, error)
	GetSqlServerDatabaseContractsAsXLSX() (*excelize.File, error)
	DeleteSqlServerDatabaseContract(id primitive.ObjectID) error
	UpdateSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) (*model.SqlServerDatabaseContract, error)

	// AckAlerts ack the specified alerts
	AckAlerts(alertsFilter dto.AlertsFilter) error
	// DismissHost dismiss the specified host
	DismissHost(hostname string) error
	// UpdateAlertsStatus update alerts status
	UpdateAlertsStatus(alertsFilter dto.AlertsFilter, newStatus string) error

	// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
	GetInfoForFrontendDashboard(location string, environment string, olderThan time.Time) (map[string]interface{}, error)

	// UpdateLicenseIgnoredField update license ignored field (true/false)
	UpdateLicenseIgnoredField(hostname string, dbname string, licensetypeid string, ignored bool, ignoredComment string) error

	// UpdateSqlServerLicenseIgnoredField update license ignored field (true/false)
	UpdateSqlServerLicenseIgnoredField(hostname string, instancename string, ignored bool, ignoredComment string) error

	// UpdateMySqlLicenseIgnoredField update license ignored field (true/false)
	UpdateMySqlLicenseIgnoredField(hostname string, instancename string, ignored bool, ignoredComment string) error

	// ALL

	SearchDatabases(filter dto.GlobalFilter) ([]dto.Database, error)
	SearchDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	GetDatabasesStatistics(filter dto.GlobalFilter) (*dto.DatabasesStatistics, error)
	GetUsedLicensesPerDatabases(hostname string, filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error)
	GetUsedLicensesPerHost(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicensePerHost, error)
	GetUsedLicensesPerHostAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	GetUsedLicensesPerCluster(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicensePerCluster, error)
	GetUsedLicensesPerClusterAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	GetDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error)
	GetDatabaseLicensesComplianceAsXLSX() (*excelize.File, error)

	// MYSQL

	SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error)
	SearchMySQLInstancesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	GetMySQLUsedLicenses(hostname string, filter dto.GlobalFilter) ([]dto.MySQLUsedLicense, error)
	GetUsedLicensesPerDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	// MYSQL CONTRACTS

	AddMySQLContract(contract model.MySQLContract) (*model.MySQLContract, error)
	UpdateMySQLContract(contract model.MySQLContract) (*model.MySQLContract, error)
	GetMySQLContracts() ([]model.MySQLContract, error)
	GetMySQLContractsAsXLSX() (*excelize.File, error)
	DeleteMySQLContract(id primitive.ObjectID) error
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
	Log logger.Logger
	// TechnologyInfos contains the list of technologies with their informations
	TechnologyInfos []model.TechnologyInfo
	// NewObjectID return a new ObjectID
	NewObjectID func() primitive.ObjectID

	mockGetOracleDatabaseContracts func(filters dto.GetOracleDatabaseContractsFilter) ([]dto.OracleDatabaseContractFE, error)

	AlertSvcClient alertServiceClient.AlertSvcClientInterface
}

// Init initializes the service and database
func (as *APIService) Init() {
	as.loadManagedTechnologiesList()

	as.NewObjectID = func() primitive.ObjectID {
		return primitive.NewObjectIDFromTimestamp(as.TimeNow())
	}
}
