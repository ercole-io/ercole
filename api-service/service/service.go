// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"encoding/csv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"go.mongodb.org/mongo-driver/bson/primitive"

	alertServiceClient "github.com/ercole-io/ercole/v2/alert-service/client"
	"github.com/ercole-io/ercole/v2/api-service/database"
	"github.com/ercole-io/ercole/v2/api-service/domain"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	alert_filter "github.com/ercole-io/ercole/v2/api-service/dto/filter"
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
	SearchAlerts(alertFilter alert_filter.Alert) (*dto.Pagination, error)
	SearchAlertsAsXLSX(status string, from, to time.Time, filter dto.GlobalFilter) (*excelize.File, error)
	GetAlerts(status string, from, to time.Time, filter dto.GlobalFilter) ([]map[string]interface{}, error)
	// SearchClusters search clusters
	SearchClusters(mode string, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]dto.Cluster, error)
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
	// SearchOracleDatabaseUsedLicenses return the list of consumed licenses
	SearchOracleDatabaseUsedLicenses(hostname string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleDatabaseUsedLicenseSearchResponse, error)

	GetOraclePsqlMigrabilities(hostname, dbname string) ([]model.PgsqlMigrability, error)
	GetOraclePsqlMigrabilitiesSemaphore(hostname, dbname string) (string, error)
	ListOracleDatabasePsqlMigrabilities() ([]dto.OracleDatabasePgsqlMigrability, error)
	ListOracleDatabasePdbPsqlMigrabilities() ([]dto.OracleDatabasePdbPgsqlMigrability, error)
	CreateOraclePsqlMigrabilitiesXlsx(dbs []dto.OracleDatabasePgsqlMigrability, pdbs []dto.OracleDatabasePdbPgsqlMigrability) (*excelize.File, error)

	GetOraclePdbPsqlMigrabilities(hostname, dbname, pdbname string) ([]model.PgsqlMigrability, error)
	GetOraclePdbPsqlMigrabilitiesSemaphore(hostname, dbname, pdbname string) (string, error)

	// ListAllLocations list locations
	ListAllLocations(location string, environment string, olderThan time.Time) ([]string, error)
	// ListLocations list locations
	ListLocations(user interface{}) ([]string, error)
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

	ImportOracleDatabaseContracts(reader *csv.Reader) error
	GetLicenseContractSample(dbtype string) ([]byte, error)

	// ORACLE DATABASE LICENSES

	GetOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error)
	GetOracleDatabaseLicensesCompliance(locations []string) ([]dto.LicenseCompliance, error)
	DeleteOracleDatabaseLicenseType(id string) error
	AddOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) (*model.OracleDatabaseLicenseType, error)
	UpdateOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) (*model.OracleDatabaseLicenseType, error)

	ListOracleGrantDbaByHostname(hostname string, filter dto.GlobalFilter) ([]dto.OracleGrantDbaDto, error)
	CreateOracleGrantDbaXlsx(hostname string, filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE PATCH
	GetOraclePatchList(filter dto.GlobalFilter) ([]dto.OracleDatabasePatchDto, error)
	CreateGetOraclePatchListXLSX(filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE OPTION
	GetOracleOptionList(filter dto.GlobalFilter) ([]dto.OracleDatabaseFeatureUsageStatDto, error)
	CreateGetOracleOptionListXLSX(filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE CHANGES
	GetOracleChanges(filter dto.GlobalFilter, hostname string) ([]dto.OracleChangesDto, error)

	// ORACLE DATABASE TABLESPACE
	ListOracleDatabaseTablespaces(filter dto.GlobalFilter) ([]dto.OracleDatabaseTablespace, error)
	CreateOracleDatabaseTablespacesXlsx(filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE SCHEMA
	ListOracleDatabaseSchemas(filter dto.GlobalFilter) ([]dto.OracleDatabaseSchema, error)
	CreateOracleDatabaseSchemasXlsx(filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE PARTITIONING
	ListOracleDatabasePartitionings(filter dto.GlobalFilter) ([]dto.OracleDatabasePartitioning, error)
	CreateOracleDatabasePartitioningsXlsx(filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE DISK GROUPS
	GetOracleDiskGroups(hostname, dbname string) ([]dto.OracleDatabaseDiskGroupDto, error)
	ListOracleDiskGroups(filter dto.GlobalFilter) ([]dto.OracleDatabaseDiskGroupDto, error)
	CreateOracleDiskGroupsXLSX(filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE PLUGGABLE DBS
	ListOracleDatabasePdbs(filter dto.GlobalFilter) ([]dto.OracleDatabasePluggableDatabase, error)
	CreateOracleDatabasePdbsXlsx(filter dto.GlobalFilter) (*excelize.File, error)
	GetOraclePDBChanges(filter dto.GlobalFilter, hostname string, start time.Time, end time.Time) ([]dto.OraclePdbChange, error)

	// ORACLE DATABASE BACKUP
	GetOracleBackupList(filter dto.GlobalFilter) ([]dto.OracleDatabaseBackupDto, error)
	CreateGetOracleBackupListXLSX(filter dto.GlobalFilter) (*excelize.File, error)

	// ORACLE DATABASE SERVICE
	GetOracleServiceList(filter dto.GlobalFilter) ([]dto.OracleDatabaseServiceDto, error)
	CreateGetOracleServiceListXLSX(filter dto.GlobalFilter) (*excelize.File, error)

	ListOracleDatabasePoliciesAudit() ([]dto.OraclePoliciesAuditListResponse, error)
	ListOracleDatabasePdbPoliciesAudit() ([]dto.OraclePdbPoliciesAuditListResponse, error)

	CreateOraclePoliciesAuditXlsx(dbs []dto.OraclePoliciesAuditListResponse, pdbs []dto.OraclePdbPoliciesAuditListResponse) (*excelize.File, error)

	GetOracleDatabasePoliciesAuditFlag(hostname, dbname string) (map[string][]string, error)
	GetOracleDatabasePdbPoliciesAuditFlag(hostname, dbname, pdbname string) (map[string][]string, error)

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
	GetSqlServerDatabaseContracts(locations []string) ([]model.SqlServerDatabaseContract, error)
	GetSqlServerDatabaseContractsAsXLSX(locations []string) (*excelize.File, error)
	DeleteSqlServerDatabaseContract(id primitive.ObjectID) error
	UpdateSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) (*model.SqlServerDatabaseContract, error)

	ImportSQLServerDatabaseContracts(reader *csv.Reader) error

	// AckAlerts ack the specified alerts
	AckAlerts(alertsFilter dto.AlertsFilter) error
	// DismissHost dismiss the specified host
	DismissHost(hostname string) error

	GetMissingDatabases() ([]dto.OracleDatabaseMissingDbs, error)
	GetMissingDatabasesByHostname(hostname string) ([]model.MissingDatabase, error)
	UpdateMissingDatabaseIgnoredField(hostname string, dbname string, ignored bool, ignoredComment string) error

	GetVirtualHostWithoutCluster() ([]dto.VirtualHostWithoutCluster, error)

	// UpdateAlertsStatus update alerts status
	UpdateAlertsStatus(alertsFilter dto.AlertsFilter, newStatus string) error

	// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
	GetInfoForFrontendDashboard(location string, environment string, olderThan time.Time) (map[string]interface{}, error)
	GetComplianceStats(user model.User) (*dto.ComplianceStats, error)

	// UpdateLicenseIgnoredField update license ignored field (true/false)
	UpdateLicenseIgnoredField(hostname string, dbname string, licensetypeid string, ignored bool, ignoredComment string) error

	CanMigrateLicense(hostname string, dbname string, filter dto.GlobalFilter) (bool, error)

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
	GetDatabaseLicensesCompliance(locations []string) ([]dto.LicenseCompliance, error)
	GetDatabaseLicensesComplianceAsXLSX(locations []string) (*excelize.File, error)

	GetClusterVeritasLicenses(filter dto.GlobalFilter) ([]dto.ClusterVeritasLicense, error)
	GetClusterVeritasLicensesXlsx(filter dto.GlobalFilter) (*excelize.File, error)

	// MYSQL

	SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error)
	SearchMySQLInstancesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)
	GetMySQLUsedLicenses(hostname string, filter dto.GlobalFilter) ([]dto.MySQLUsedLicense, error)
	GetUsedLicensesPerDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error)

	GetHostsMysqlAsLMS(filters dto.SearchHostsAsLMS) (*excelize.File, error)

	// MYSQL CONTRACTS

	AddMySQLContract(contract model.MySQLContract) (*model.MySQLContract, error)
	UpdateMySQLContract(contract model.MySQLContract) (*model.MySQLContract, error)
	GetMySQLContracts(locations []string) ([]model.MySQLContract, error)
	GetMySQLContractsAsXLSX(locations []string) (*excelize.File, error)
	DeleteMySQLContract(id primitive.ObjectID) error

	ImportMySQLDatabaseContracts(reader *csv.Reader) error

	// POSTGRESQL
	// SearchSqlServerInstances search databases
	SearchPostgreSqlInstances(filter dto.SearchPostgreSqlInstancesFilter) (*dto.PostgreSqlInstanceResponse, error)
	// SearchOracleDatabases search databases
	SearchPostgreSqlInstancesAsXLSX(filter dto.SearchPostgreSqlInstancesFilter) (*excelize.File, error)

	// MONGODB
	// SearchMongoDBInstances search databases
	SearchMongoDBInstances(filter dto.SearchMongoDBInstancesFilter) (*dto.MongoDBInstanceResponse, error)
	// SearchOracleDatabases search databases
	SearchMongoDBInstancesAsXLSX(filter dto.SearchMongoDBInstancesFilter) (*excelize.File, error)

	// ROLES
	GetRole(name string) (*model.Role, error)
	GetRoles() ([]model.Role, error)
	AddRole(role model.Role) error
	UpdateRole(role model.Role) error
	RemoveRole(roleName string) error

	// GROUPS
	InsertGroup(group model.Group) (*model.Group, error)
	UpdateGroup(group model.Group) (*model.Group, error)
	GetGroup(name string) (*model.Group, error)
	GetGroups() ([]model.Group, error)
	DeleteGroup(name string) error
	GetMatchedGroupsName(tags []string) []string

	GetDatabaseConnectionStatus() bool
	GetConfig() (*config.Configuration, error)
	ChangeConfig(config config.Configuration) error

	ListUsers() ([]model.User, error)
	GetUser(username string) (*model.User, error)
	AddUser(user model.User) error
	UpdateUserGroups(username string, groups []string) error
	UpdateUserLastLogin(updatedUser model.User) error
	RemoveLimitedGroup(updatedUser model.User) error
	AddLimitedGroup(updatedUser model.User) error
	RemoveUser(username string) error
	NewPassword(username string) (string, error)
	UpdatePassword(username string, oldPass string, newPass string) error
	MatchPassword(user *model.User, password string) bool
	GetUserLocations(username string) ([]string, error)

	GetNodes(groups []string) ([]model.Node, error)
	GetNode(name string) (*model.Node, error)
	AddNode(node model.Node) error
	UpdateNode(node model.Node) error
	RemoveNode(name string) error

	// EXADATA
	ListExadataInstances(filter dto.GlobalFilter, hidden bool) ([]dto.ExadataInstanceResponse, error)
	GetExadataInstance(rackid string, hidden bool) (*domain.OracleExadataInstance, error)
	UpdateExadataVmClusterName(rackID, hostID, vmname, clustername string) error
	UpdateExadataComponentClusterName(RackID, hostID string, clusternames []string) error
	UpdateExadataRdma(rackID string, rdma model.OracleExadataRdma) error
	GetAllExadataInstanceAsXlsx() (*excelize.File, error)
	HideExadataInstance(rackID string) error
	ShowExadataInstance(rackID string) error
	GetExadataPatchAdvisors() ([]dto.OracleExadataPatchAdvisor, error)
	GetAllExadataPatchAdvisorsAsXlsx() (*excelize.File, error)
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
