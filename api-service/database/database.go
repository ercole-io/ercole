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

// Package database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/config"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	CheckStatusMongodb() error
	// SearchHosts search hosts
	SearchHosts(mode string, filters dto.SearchHostsFilters) ([]map[string]interface{}, error)
	GetHostDataSummaries(filters dto.SearchHostsFilters) ([]dto.HostDataSummary, error)
	// GetHost fetch all informations about a host in the database
	GetHost(hostname string, olderThan time.Time, raw bool) (*dto.HostData, error)
	GetHostData(hostname string, olderThan time.Time) (*model.HostDataBE, error)
	GetHostDatas(olderThan time.Time) ([]model.HostDataBE, error)
	// SearchAlerts search alerts
	SearchAlerts(mode string, keywords []string, sortBy string, sortDesc bool, page, pageSize int, location, environment, severity, status string, from, to time.Time) ([]map[string]interface{}, error)
	// SearchClusters search clusters
	SearchClusters(mode string, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	GetClusters(filter dto.GlobalFilter) ([]dto.Cluster, error)
	// GetCluster fetch all information about a cluster in the database
	GetCluster(clusterName string, olderThan time.Time) (*dto.Cluster, error)
	// SearchOracleDatabaseAddms search addms
	SearchOracleDatabaseAddms(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// SearchOracleDatabaseSegmentAdvisors search segment advisors
	SearchOracleDatabaseSegmentAdvisors(keywords []string, sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]dto.OracleDatabaseSegmentAdvisor, error)
	// SearchOracleDatabasePatchAdvisors search patch advisors
	SearchOracleDatabasePatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) (*dto.PatchAdvisorResponse, error)
	// SearchOracleDatabases search databases
	SearchOracleDatabases(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleDatabaseResponse, error)
	// SearchOracleExadata search exadata
	SearchOracleExadata(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleExadataResponse, error)
	// SearchOracleDatabaseUsedLicenses search consumed licenses
	SearchOracleDatabaseUsedLicenses(hostname string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleDatabaseUsedLicenseSearchResponse, error)

	// ListLocations list locations
	ListLocations(location string, environment string, olderThan time.Time) ([]string, error)
	// ListEnvironments list environments
	ListEnvironments(location string, environment string, olderThan time.Time) ([]string, error)
	// GetHostsCountStats return the number of the non-archived hosts
	GetHostsCountStats(location string, environment string, olderThan time.Time) (int, error)
	// GetEnvironmentStats return a array containing the number of hosts per environment
	GetEnvironmentStats(location string, olderThan time.Time) ([]interface{}, error)
	// GetTypeStats return a array containing the number of hosts per operating system
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
	//GetOracleDatabaseLicenseTypes return an array of OracleDatabaseLicenseType
	GetOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error)
	//GetOracleDatabaseLicenseType return a OracleDatabaseLicenseType
	GetOracleDatabaseLicenseType(id string) (*model.OracleDatabaseLicenseType, error)
	// InsertOracleDatabaseContract insert an Oracle/Database contract into the database
	InsertOracleDatabaseContract(contract model.OracleDatabaseContract) error
	// GetOracleDatabaseContract return the contract specified by id
	GetOracleDatabaseContract(id primitive.ObjectID) (*model.OracleDatabaseContract, error)
	// UpdateOracleDatabaseContract update an Oracle/Database contract in the database
	UpdateOracleDatabaseContract(contract model.OracleDatabaseContract) error
	// RemoveOracleDatabaseContract remove an Oracle/Database contract from the database
	RemoveOracleDatabaseContract(id primitive.ObjectID) error

	// ListOracleDatabasContracts lists the Oracle/Database contracts
	ListOracleDatabaseContracts() ([]dto.OracleDatabaseContractFE, error)
	// UpdateLicenseIgnoredField update license ignored field (true/false)
	UpdateLicenseIgnoredField(hostname string, dbname string, licenseTypeID string, ignored bool, ignoredComment string) error

	// InsertOracleDatabaseLicenseType insert an Oracle/Database license type into the database
	InsertOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) error
	// UpdateOracleDatabaseLicenseType update an Oracle/Database license type in the database
	UpdateOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) error
	// RemoveOracleDatabaseLicenseType remove a licence type - Oracle/Database contract part from the database
	RemoveOracleDatabaseLicenseType(id string) error

	FindGrantDbaByHostname(hostname string, filter dto.GlobalFilter) ([]dto.OracleGrantDbaDto, error)

	GetOraclePatchList() ([]dto.OracleDatabasePatchDto, error)
	GetOracleOptionList() ([]dto.OracleDatabaseFeatureUsageStatDto, error)

	// ReplaceHostData adds a new hostdata to the database
	ReplaceHostData(hostData model.HostDataBE) error
	// UpdateAlertsStatus change the status of the specified alerts
	UpdateAlertsStatus(alertsFilter dto.AlertsFilter, newStatus string) error
	// CountAlertsNODATA gets alert with alertCode equals to "NO_DATA"
	CountAlertsNODATA(alertsFilter dto.AlertsFilter) (int64, error)
	// DismissHost dismiss the specified host
	DismissHost(hostname string) error
	// GetHostMinValidCreatedAtDate get the host's minimun valid CreatedAt date
	GetHostMinValidCreatedAtDate(hostname string) (time.Time, error)
	// GetListValidHostsByRangeDates get list of valid hosts by range dates
	GetListValidHostsByRangeDates(from time.Time, to time.Time) ([]string, error)
	// GetListDismissedHostsByRangeDates get list of dismissed hosts by range dates
	GetListDismissedHostsByRangeDates(from time.Time, to time.Time) ([]string, error)
	// RemoveAlertsNODATA delete all alerts with alertCode equals to "NO_DATA"
	RemoveAlertsNODATA(alertsFilter dto.AlertsFilter) error

	// FindHostData find the current hostdata with a certain hostname
	FindHostData(hostname string) (model.HostDataBE, error)
	// ExistHostdata return true if the host specified by hostname exist, otherwise false
	ExistHostdata(hostname string) (bool, error)
	// GetHostsCountUsingTechnologies return a map that contains the number of usages for every features
	GetHostsCountUsingTechnologies(location string, environment string, olderThan time.Time) (map[string]float64, error)
	// ExistNotInClusterHost return true if the host specified by hostname exist and it is not in cluster, otherwise false
	ExistNotInClusterHost(hostname string) (bool, error)

	FindAllOracleDatabaseTablespaces(filter dto.GlobalFilter) ([]dto.OracleDatabaseTablespace, error)

	// MYSQL

	SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error)
	//GetMySQLUsedLicenses return MySQL used licenses.
	// Only ENTERPRISE MySQL db are considered as licenses
	GetMySQLUsedLicenses(hostname string, filter dto.GlobalFilter) ([]dto.MySQLUsedLicense, error)
	UpdateMySqlLicenseIgnoredField(hostname string, instancename string, ignored bool, ignoredComment string) error

	// MYSQL CONTRACTS

	AddMySQLContract(contract model.MySQLContract) error
	UpdateMySQLContract(contract model.MySQLContract) error
	GetMySQLContracts() ([]model.MySQLContract, error)
	DeleteMySQLContract(id primitive.ObjectID) error

	// SQL SERVER

	GetSqlServerDatabaseLicenseTypes() ([]model.SqlServerDatabaseLicenseType, error)
	InsertSqlServerDatabaseLicenseType(licenseType model.SqlServerDatabaseLicenseType) error
	SearchSqlServerInstances(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.SqlServerInstanceResponse, error)
	SearchSqlServerDatabaseUsedLicenses(hostname string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.SqlServerDatabaseUsedLicenseSearchResponse, error)
	UpdateSqlServerLicenseIgnoredField(hostname string, instancename string, ignored bool, ignoredComment string) error

	InsertSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) error
	GetSqlServerDatabaseLicenseType(id string) (*model.SqlServerDatabaseLicenseType, error)
	ListSqlServerDatabaseContracts() ([]model.SqlServerDatabaseContract, error)
	RemoveSqlServerDatabaseContract(id primitive.ObjectID) error
	UpdateSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) error

	// POSTGRESQL
	SearchPostgreSqlInstances(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.PostgreSqlInstanceResponse, error)
}

// MongoDatabase is a implementation
type MongoDatabase struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Client contain the mongodb client
	Client *mongo.Client
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// OperatingSystemAggregationRules contains rules used to aggregate various operating systems
	OperatingSystemAggregationRules []config.AggregationRule
	// Log contains logger formatted
	Log logger.Logger
}

// Init initializes the connection to the database
func (md *MongoDatabase) Init() {
	md.ConnectToMongodb()

	md.Log.Debug("MongoDatabase is connected to MongoDB! ", utils.HideMongoDBPassword(md.Config.Mongodb.URI))
}

// ConnectToMongodb connects to the MongoDB and return the connection
func (md *MongoDatabase) ConnectToMongodb() {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(md.Config.Mongodb.URI)

	//Connect to MongoDB
	md.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		md.Log.Warn(err)
	}

	//Check the connection
	err = md.Client.Ping(context.TODO(), nil)
	if err != nil {
		md.Log.Warn(err)
	}
}

func (md *MongoDatabase) CheckStatusMongodb() error {
	var err error

	clientOptions := options.Client().ApplyURI(md.Config.Mongodb.URI)

	md.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = md.Client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
