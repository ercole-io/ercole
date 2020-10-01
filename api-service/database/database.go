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

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// SearchHosts search hosts
	SearchHosts(mode string, keywords []string, otherFilters SearchHostsFilters, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// GetHost fetch all informations about a host in the database
	GetHost(hostname string, olderThan time.Time, raw bool) (interface{}, utils.AdvancedErrorInterface)
	// SearchAlerts search alerts
	SearchAlerts(mode string, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, severity string, status string, from time.Time, to time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchClusters search clusters
	SearchClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// GetCluster fetch all information about a cluster in the database
	GetCluster(clusterName string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabaseAddms search addms
	SearchOracleDatabaseAddms(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabaseSegmentAdvisors search segment advisors
	SearchOracleDatabaseSegmentAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabasePatchAdvisors search patch advisors
	SearchOracleDatabasePatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabases search databases
	SearchOracleDatabases(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface)
	// SearchOracleExadata search exadata
	SearchOracleExadata(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchLicenses search licenses
	SearchLicenses(location string, environment string, olderThan time.Time) ([]apimodel.OracleDatabaseLicenseUsageInfo, utils.AdvancedErrorInterface)
	// ListLicenses list licenses
	ListLicenses(sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetLicense get a certain license
	GetLicense(name string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// SearchOracleDatabaseLicenseModifiers search license modifiers
	SearchOracleDatabaseLicenseModifiers(keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]map[string]interface{}, utils.AdvancedErrorInterface)

	// ListLocations list locations
	ListLocations(location string, environment string, olderThan time.Time) ([]string, utils.AdvancedErrorInterface)
	// ListEnvironments list environments
	ListEnvironments(location string, environment string, olderThan time.Time) ([]string, utils.AdvancedErrorInterface)
	// GetHostsCountStats return the number of the non-archived hosts
	GetHostsCountStats(location string, environment string, olderThan time.Time) (int, utils.AdvancedErrorInterface)
	// GetEnvironmentStats return a array containing the number of hosts per environment
	GetEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per operating system
	GetOperatingSystemStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per type
	GetTypeStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopUnusedOracleDatabaseInstanceResourceStats return a array containing top unused instance resource by workload
	GetTopUnusedOracleDatabaseInstanceResourceStats(location string, environment string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseEnvironmentStats return a array containing the number of databases per environment
	GetOracleDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseHighReliabilityStats return a array containing the number of databases per high-reliability status
	GetOracleDatabaseHighReliabilityStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseVersionStats return a array containing the number of databases per version
	GetOracleDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopReclaimableOracleDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
	GetTopReclaimableOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabasePatchStatusStats return a array containing the number of databases per patch status
	GetOracleDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopWorkloadOracleDatabaseStats return a array containing top databases by workload
	GetTopWorkloadOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
	GetOracleDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseRACStatusStats return a array containing the number of databases per RAC status
	GetOracleDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
	GetOracleDatabaseLicenseComplianceStatusStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetOracleDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
	GetOracleDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseWorkStats return the total work of databases
	GetTotalOracleDatabaseWorkStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseMemorySizeStats return the total of memory size of databases
	GetTotalOracleDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseDatafileSizeStats return the total size of datafiles of databases
	GetTotalOracleDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleDatabaseSegmentSizeStats return the total size of segments of databases
	GetTotalOracleDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleExadataMemorySizeStats return the total size of memory of exadata
	GetTotalOracleExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetTotalOracleExadataCPUStats return the total cpu of exadata
	GetTotalOracleExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface)
	// GetAverageOracleExadataStorageUsageStats return the average usage of cell disks of exadata
	GetAverageOracleExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface)
	// GetOracleExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
	GetOracleExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)
	// GetOracleExadataPatchStatusStats return a array containing the number of exadata per patch status
	GetOracleExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface)

	// InsertOracleDatabaseAgreement insert an Oracle/Database agreement into the database
	InsertOracleDatabaseAgreement(aggreement model.OracleDatabaseAgreement) (*mongo.InsertOneResult, utils.AdvancedErrorInterface)
	// UpdateOracleDatabaseAgreement update an Oracle/Database agreement in the database
	UpdateOracleDatabaseAgreement(newAgreement model.OracleDatabaseAgreement) utils.AdvancedErrorInterface
	// RemoveOracleDatabaseAgreement remove an Oracle/Database agreement from the database
	RemoveOracleDatabaseAgreement(id primitive.ObjectID) utils.AdvancedErrorInterface

	// ListOracleDatabaseAgreements lists the Oracle/Database agreements
	ListOracleDatabaseAgreements() ([]apimodel.OracleDatabaseAgreementFE, utils.AdvancedErrorInterface)
	// ListHostUsingOracleDatabaseLicenses lists the hosts/clusters that need to be licensed by Oracle/Database agreements
	ListHostUsingOracleDatabaseLicenses() ([]apimodel.HostUsingOracleDatabaseLicenses, utils.AdvancedErrorInterface)
	// FindOracleDatabaseAgreement return the agreement specified by id
	FindOracleDatabaseAgreement(id primitive.ObjectID) (model.OracleDatabaseAgreement, utils.AdvancedErrorInterface)

	// SetLicenseCostPerProcessor set the cost per processor of a certain license
	SetLicenseCostPerProcessor(name string, costPerProcessor float64) utils.AdvancedErrorInterface

	// SavePatchingFunction saves the patching function
	SavePatchingFunction(pf model.PatchingFunction) utils.AdvancedErrorInterface
	// ReplaceHostData adds a new hostdata to the database
	ReplaceHostData(hostData model.HostDataBE) utils.AdvancedErrorInterface
	// UpdateAlertsStatus change the status of the specified alerts
	UpdateAlertsStatus(id []primitive.ObjectID, newStatus string) utils.AdvancedErrorInterface
	// ArchiveHost archive the specified host
	ArchiveHost(hostname string) utils.AdvancedErrorInterface
	// DeletePatchingFunction delete the patching function
	DeletePatchingFunction(hostname string) utils.AdvancedErrorInterface

	// FindPatchingFunction find the the patching function associated to the hostname in the database
	FindPatchingFunction(hostname string) (model.PatchingFunction, utils.AdvancedErrorInterface)
	// FindHostData find the current hostdata with a certain hostname
	FindHostData(hostname string) (model.HostDataBE, utils.AdvancedErrorInterface)
	// ExistHostdata return true if the host specified by hostname exist, otherwise false
	ExistHostdata(hostname string) (bool, utils.AdvancedErrorInterface)
	// GetTechnologiesUsage return a map that contains the number of usages for every features
	GetTechnologiesUsage(location string, environment string, olderThan time.Time) (map[string]float64, utils.AdvancedErrorInterface)
	// ExistNotInClusterHost return true if the host specified by hostname exist and it is not in cluster, otherwise false
	ExistNotInClusterHost(hostname string) (bool, utils.AdvancedErrorInterface)
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

// SearchHostsFilters contains all filters for the SearchHosts API
type SearchHostsFilters struct {
	Hostname                      string
	Database                      string
	Technology                    string
	HardwareAbstractionTechnology string
	Cluster                       *string
	VirtualizationNode            string
	OperatingSystem               string
	Kernel                        string
	LTEMemoryTotal                float64
	GTEMemoryTotal                float64
	LTESwapTotal                  float64
	GTESwapTotal                  float64
	IsMemberOfCluster             *bool
	CPUModel                      string
	LTECPUCores                   int
	GTECPUCores                   int
	LTECPUThreads                 int
	GTECPUThreads                 int
}
