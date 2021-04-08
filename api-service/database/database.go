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

// Package database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/v2/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// SearchHosts search hosts
	SearchHosts(mode string, filters dto.SearchHostsFilters) ([]map[string]interface{}, error)
	// GetHost fetch all informations about a host in the database
	GetHost(hostname string, olderThan time.Time, raw bool) (interface{}, error)
	GetHostData(hostname string, olderThan time.Time) (*model.HostDataBE, error)
	// SearchAlerts search alerts
	SearchAlerts(mode string, keywords []string, sortBy string, sortDesc bool, page, pageSize int, location, environment, severity, status string, from, to time.Time) ([]map[string]interface{}, error)
	// SearchClusters search clusters
	SearchClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	GetClusters(filter dto.GlobalFilter) ([]dto.Cluster, error)
	// GetCluster fetch all information about a cluster in the database
	GetCluster(clusterName string, olderThan time.Time) (*dto.Cluster, error)
	// SearchOracleDatabaseAddms search addms
	SearchOracleDatabaseAddms(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// SearchOracleDatabaseSegmentAdvisors search segment advisors
	SearchOracleDatabaseSegmentAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// SearchOracleDatabasePatchAdvisors search patch advisors
	SearchOracleDatabasePatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) ([]map[string]interface{}, error)
	// SearchOracleDatabases search databases
	SearchOracleDatabases(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error)
	// SearchOracleExadata search exadata
	SearchOracleExadata(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, error)
	// SearchOracleDatabaseUsedLicenses search consumed licenses
	SearchOracleDatabaseUsedLicenses(sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleDatabaseUsedLicenseSearchResponse, error)
	// SearchOracleDatabaseLicenseModifiers search license modifiers
	SearchOracleDatabaseLicenseModifiers(keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]map[string]interface{}, error)

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

	GetOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error)
	// InsertOracleDatabaseAgreement insert an Oracle/Database agreement into the database
	InsertOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) (*mongo.InsertOneResult, error)
	// GetOracleDatabaseAgreement return the agreement specified by id
	GetOracleDatabaseAgreement(agreementID string) (*model.OracleDatabaseAgreement, error)
	// GetOracleDatabaseAgreementByAssociatedLicenseType return the agreement specified by an associated part id
	GetOracleDatabaseAgreementByAssociatedLicenseType(associateLicenseTypeID primitive.ObjectID) (*model.OracleDatabaseAgreement, error)
	// UpdateOracleDatabaseAgreement update an Oracle/Database agreement in the database
	UpdateOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) error
	// RemoveOracleDatabaseAgreement remove an Oracle/Database agreement from the database
	RemoveOracleDatabaseAgreement(id primitive.ObjectID) error

	// ListOracleDatabaseAgreements lists the Oracle/Database agreements
	ListOracleDatabaseAgreements() ([]dto.OracleDatabaseAgreementFE, error)
	// ListHostUsingOracleDatabaseLicenses lists the hosts/clusters that need to be licensed by Oracle/Database agreements
	ListHostUsingOracleDatabaseLicenses() ([]dto.HostUsingOracleDatabaseLicenses, error)

	// SavePatchingFunction saves the patching function
	SavePatchingFunction(pf model.PatchingFunction) error
	// ReplaceHostData adds a new hostdata to the database
	ReplaceHostData(hostData model.HostDataBE) error
	// UpdateAlertsStatus change the status of the specified alerts
	UpdateAlertsStatus(id []primitive.ObjectID, newStatus string) error
	// ArchiveHost archive the specified host
	ArchiveHost(hostname string) error
	// DeletePatchingFunction delete the patching function
	DeletePatchingFunction(hostname string) error

	// FindPatchingFunction find the the patching function associated to the hostname in the database
	FindPatchingFunction(hostname string) (model.PatchingFunction, error)
	// FindHostData find the current hostdata with a certain hostname
	FindHostData(hostname string) (model.HostDataBE, error)
	// ExistHostdata return true if the host specified by hostname exist, otherwise false
	ExistHostdata(hostname string) (bool, error)
	// GetHostsCountUsingTechnologies return a map that contains the number of usages for every features
	GetHostsCountUsingTechnologies(location string, environment string, olderThan time.Time) (map[string]float64, error)
	// ExistNotInClusterHost return true if the host specified by hostname exist and it is not in cluster, otherwise false
	ExistNotInClusterHost(hostname string) (bool, error)

	// MYSQL

	SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error)
	//GetMySQLUsedLicenses return MySQL used licenses.
	// Only ENTERPRISE MySQL db are considered as licenses
	GetMySQLUsedLicenses(filter dto.GlobalFilter) ([]dto.MySQLUsedLicense, error)

	// MYSQL AGREEMENTS

	AddMySQLAgreement(agreement model.MySQLAgreement) (primitive.ObjectID, error)
	UpdateMySQLAgreement(agreement model.MySQLAgreement) error
	GetMySQLAgreements() ([]model.MySQLAgreement, error)
	DeleteMySQLAgreement(id primitive.ObjectID) error
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
	Log *logrus.Logger
}

// Init initializes the connection to the database
func (md *MongoDatabase) Init() {
	//Connect to mongodb
	md.ConnectToMongodb()
	md.Log.Info("MongoDatabase is connected to MongoDB! ", md.Config.Mongodb.URI)
}

// ConnectToMongodb connects to the MongoDB and return the connection
func (md *MongoDatabase) ConnectToMongodb() {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(md.Config.Mongodb.URI)

	//Connect to MongoDB
	md.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		md.Log.Fatal(err)
	}

	//Check the connection
	err = md.Client.Ping(context.TODO(), nil)
	if err != nil {
		md.Log.Fatal(err)
	}
}
