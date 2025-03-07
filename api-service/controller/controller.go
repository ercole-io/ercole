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

package controller

import (
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/auth"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/api-service/service"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/context"
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
	// SearchOracleDatabaseUsedLicenses search licenses consumed by the hosts using the filters in the request
	SearchOracleDatabaseUsedLicenses(w http.ResponseWriter, r *http.Request)

	GetOraclePsqlMigrabilities(w http.ResponseWriter, r *http.Request)
	GetOraclePsqlMigrabilitiesSemaphore(w http.ResponseWriter, r *http.Request)
	ListOracleDatabasePsqlMigrabilities(w http.ResponseWriter, r *http.Request)
	ListOracleDatabasePdbPsqlMigrabilities(w http.ResponseWriter, r *http.Request)

	ListOraclePoliciesAudit(w http.ResponseWriter, r *http.Request)
	ListOraclePdbPoliciesAudit(w http.ResponseWriter, r *http.Request)

	GetOraclePoliciesAudit(w http.ResponseWriter, r *http.Request)
	GetOraclePdbsPoliciesAudit(w http.ResponseWriter, r *http.Request)

	GetOraclePdbPsqlMigrabilities(w http.ResponseWriter, r *http.Request)
	GetOraclePdbPsqlMigrabilitiesSemaphore(w http.ResponseWriter, r *http.Request)

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
	GetOracleDatabasesStatistics(w http.ResponseWriter, r *http.Request)
	// GetOracleDatabaseLicensesCompliance return licenses usage status and compliance
	GetOracleDatabaseLicensesCompliance(w http.ResponseWriter, r *http.Request)

	// GetDefaultDatabaseTags return the default list of database tags from configuration
	GetDefaultDatabaseTags(w http.ResponseWriter, r *http.Request)
	// GetErcoleFeatures return a map of active/inactive features
	GetErcoleFeatures(w http.ResponseWriter, r *http.Request)
	// GetTechnologyList return the list of techonlogies
	GetTechnologyList(w http.ResponseWriter, r *http.Request)

	AddTagToOracleDatabase(w http.ResponseWriter, r *http.Request)
	// DeleteTagOfOracleDatabase remove a certain tag from a database if it has the tag
	DeleteTagOfOracleDatabase(w http.ResponseWriter, r *http.Request)
	// AckAlerts ack the specified alert in the request
	AckAlerts(w http.ResponseWriter, r *http.Request)
	// DismissHost dismiss the specified host in the request
	DismissHost(w http.ResponseWriter, r *http.Request)

	GetMissingDatabases(w http.ResponseWriter, r *http.Request)
	GetMissingDatabasesByHostname(w http.ResponseWriter, r *http.Request)

	GetVirtualHostWithoutCluster(w http.ResponseWriter, r *http.Request)

	// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
	GetInfoForFrontendDashboard(w http.ResponseWriter, r *http.Request)

	// UpdateLicenseIgnoredField update license ignored field (true/false)
	UpdateLicenseIgnoredField(w http.ResponseWriter, r *http.Request)

	CanMigrateLicense(w http.ResponseWriter, r *http.Request)

	// UpdateSqlServerLicenseIgnoredField update license ignored field (true/false)
	UpdateSqlServerLicenseIgnoredField(w http.ResponseWriter, r *http.Request)
	// UpdateMySqlLicenseIgnoredField update license ignored field (true/false)
	UpdateMySqlLicenseIgnoredField(w http.ResponseWriter, r *http.Request)

	// ALL TECHNOLOGIES

	SearchDatabases(w http.ResponseWriter, r *http.Request)
	GetDatabasesStatistics(w http.ResponseWriter, r *http.Request)
	GetUsedLicensesPerDatabases(w http.ResponseWriter, r *http.Request)
	GetUsedLicensesPerDatabasesByHost(w http.ResponseWriter, r *http.Request)
	GetUsedLicensesPerHost(w http.ResponseWriter, r *http.Request)
	GetUsedLicensesPerCluster(w http.ResponseWriter, r *http.Request)
	GetDatabaseLicensesCompliance(w http.ResponseWriter, r *http.Request)

	// ORACLE DATABASE CONTRACTS

	AddOracleDatabaseContract(w http.ResponseWriter, r *http.Request)
	UpdateOracleDatabaseContract(w http.ResponseWriter, r *http.Request)
	GetOracleDatabaseContracts(w http.ResponseWriter, r *http.Request)
	DeleteOracleDatabaseContract(w http.ResponseWriter, r *http.Request)

	AddHostToOracleDatabaseContract(w http.ResponseWriter, r *http.Request)
	DeleteHostFromOracleDatabaseContract(w http.ResponseWriter, r *http.Request)

	ImportContractFromCSV(w http.ResponseWriter, r *http.Request)
	GetContractSampleCSV(w http.ResponseWriter, r *http.Request)

	// ORACLE DATABASE LICENSE TYPES

	// GetOracleDatabaseLicenseTypes return the list of Oracle/Database contract parts
	GetOracleDatabaseLicenseTypes(w http.ResponseWriter, r *http.Request)
	// DeleteOracleDatabaseLicenseType remove a licence type - Oracle/Database contract part
	DeleteOracleDatabaseLicenseType(w http.ResponseWriter, r *http.Request)
	// AddOracleDatabaseLicenseType add a licence type - Oracle/Database contract part to the database if it hasn't a licence type
	AddOracleDatabaseLicenseType(w http.ResponseWriter, r *http.Request)
	// UpdateOracleDatabaseLicenseType update a licence type - Oracle/Database contract part
	UpdateOracleDatabaseLicenseType(w http.ResponseWriter, r *http.Request)

	ListOracleGrantDbaByHostname(w http.ResponseWriter, r *http.Request)
	GetOracleGrantDbaJSON(hostname string, filters *dto.GlobalFilter) ([]dto.OracleGrantDbaDto, error)
	GetOracleGrantDbaXLSX(hostname string, filters *dto.GlobalFilter) (*excelize.File, error)

	GetOraclePatchList(w http.ResponseWriter, r *http.Request)
	GetOracleOptionList(w http.ResponseWriter, r *http.Request)
	GetOracleChanges(w http.ResponseWriter, r *http.Request)
	GetOraclePDBChanges(w http.ResponseWriter, r *http.Request)
	GetOracleBackupList(w http.ResponseWriter, r *http.Request)
	GetOracleServiceList(w http.ResponseWriter, r *http.Request)
	ListOracleDatabasePartitionings(w http.ResponseWriter, r *http.Request)
	GetOracleDiskGroups(w http.ResponseWriter, r *http.Request)

	// MYSQL

	SearchMySQLInstances(w http.ResponseWriter, r *http.Request)
	GetMySqlLicenseTypes(w http.ResponseWriter, r *http.Request)

	// SQL SERVER
	// SearchSqlServerInstances search instances data using the filters in the request
	SearchSqlServerInstances(w http.ResponseWriter, r *http.Request)

	// POSTGRESQL
	// SearchPostgreSqlInstances search instances data using the filters in the request
	SearchPostgreSqlInstances(w http.ResponseWriter, r *http.Request)

	// MONGODB
	// SearchMongoDBInstances search instances data using the filters in the request
	SearchMongoDBInstances(w http.ResponseWriter, r *http.Request)

	// MYSQL CONTRACTS
	AddMySQLContract(w http.ResponseWriter, r *http.Request)
	UpdateMySQLContract(w http.ResponseWriter, r *http.Request)
	GetMySQLContracts(w http.ResponseWriter, r *http.Request)
	DeleteMySQLContract(w http.ResponseWriter, r *http.Request)

	// ROLES
	GetRole(w http.ResponseWriter, r *http.Request)
	GetRoles(w http.ResponseWriter, r *http.Request)
	AddRole(w http.ResponseWriter, r *http.Request)
	UpdateRole(w http.ResponseWriter, r *http.Request)
	RemoveRole(w http.ResponseWriter, r *http.Request)

	// GROUPS
	InsertGroup(w http.ResponseWriter, r *http.Request)
	UpdateGroup(w http.ResponseWriter, r *http.Request)
	GetGroup(w http.ResponseWriter, r *http.Request)
	GetGroups(w http.ResponseWriter, r *http.Request)
	DeleteGroup(w http.ResponseWriter, r *http.Request)

	GetConfig(w http.ResponseWriter, r *http.Request)
	ChangeConfig(w http.ResponseWriter, r *http.Request)

	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	AddUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	RemoveUser(w http.ResponseWriter, r *http.Request)
	GetInfo(w http.ResponseWriter, r *http.Request)

	GetNodes(w http.ResponseWriter, r *http.Request)
	GetNode(w http.ResponseWriter, r *http.Request)
	AddNode(w http.ResponseWriter, r *http.Request)
	UpdateNode(w http.ResponseWriter, r *http.Request)
	RemoveNode(w http.ResponseWriter, r *http.Request)

	// EXADATA
	ListExadata(w http.ResponseWriter, r *http.Request)
	ListHiddenExadata(w http.ResponseWriter, r *http.Request)
	GetExadata(w http.ResponseWriter, r *http.Request)
	UpdateExadataVmClusterName(w http.ResponseWriter, r *http.Request)
	UpdateExadataComponentClusterName(w http.ResponseWriter, r *http.Request)
	UpdateExadataRdma(w http.ResponseWriter, r *http.Request)
	ExportExadataInstances(w http.ResponseWriter, r *http.Request)
	HideExadataInstance(w http.ResponseWriter, r *http.Request)
	ShowExadataInstance(w http.ResponseWriter, r *http.Request)
	ListExadataPatchAdvisors(w http.ResponseWriter, r *http.Request)
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
	Log logger.Logger
	// Authenticator contains the authenticator
	Authenticator []auth.AuthenticationProvider
}

func (ctrl *APIController) userHasAccessToLocation(r *http.Request, location string) bool {
	locations, err := ctrl.Service.ListLocations(context.Get(r, "user"))
	if err != nil {
		return false
	}

	return utils.ContainsSomeI(locations, location, model.AllLocation)
}
