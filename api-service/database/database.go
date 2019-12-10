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

// Database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"log"
	"time"

	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// SearchCurrentHosts search current hosts
	SearchCurrentHosts(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface)
	// GetCurrentHost fetch all informations about a current host in the database
	GetCurrentHost(hostname string) (interface{}, utils.AdvancedErrorInterface)
	// SearchAlerts search alerts
	SearchAlerts(keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentClusters search current clusters
	SearchCurrentClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentAddms search current addms
	SearchCurrentAddms(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentSegmentAdvisors search current segment advisors
	SearchCurrentSegmentAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentPatchAdvisors search current patch advisors
	SearchCurrentPatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface)
	// SearchCurrentDatabases search current databases
	SearchCurrentDatabases(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface)

	// GetEnvironmentStats return a array containing the number of hosts per environment
	GetEnvironmentStats(location string) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per operating system
	GetOperatingSystemStats(location string) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTypeStats return a array containing the number of hosts per type
	GetTypeStats(location string) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseEnvironmentStats return a array containing the number of databases per environment
	GetDatabaseEnvironmentStats(location string) ([]interface{}, utils.AdvancedErrorInterface)
	// GetDatabaseVersionStats return a array containing the number of databases per version
	GetDatabaseVersionStats(location string) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopReclaimableDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
	GetTopReclaimableDatabaseStats(location string, limit int) ([]interface{}, utils.AdvancedErrorInterface)
	// GetTopWorkloadDatabaseStats return a array containing top databases by workload
	GetTopWorkloadDatabaseStats(location string, limit int) ([]interface{}, utils.AdvancedErrorInterface)
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
}

// Init initializes the connection to the database
func (md *MongoDatabase) Init() {
	//Connect to mongodb
	md.ConnectToMongodb()
	log.Println("MongoDatabase is connected to MongoDB!", md.Config.Mongodb.URI)
}

// ConnectToMongodb connects to the MongoDB and return the connection
func (md *MongoDatabase) ConnectToMongodb() {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(md.Config.Mongodb.URI)

	//Connect to MongoDB
	md.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//Check the connection
	err = md.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func optionalStep(optional bool, step bson.M) bson.M {
	if optional {
		return step
	}
	return bson.M{"$skip": 0}
}

func optionalSortingStep(sortBy string, sortDesc bool) bson.M {
	if sortBy == "" {
		return bson.M{"$skip": 0}
	}

	sortOrder := 0
	if sortDesc {
		sortOrder = -1
	} else {
		sortOrder = 1
	}

	return bson.M{"$sort": bson.M{
		sortBy: sortOrder,
	}}
}

func optionalPagingStep(page int, size int) bson.M {
	if page == -1 || size == -1 {
		return bson.M{"$skip": 0}
	}

	return bson.M{"$facet": bson.M{
		"content": bson.A{
			bson.M{"$skip": page * size},
			bson.M{"$limit": size},
		},
		"metadata": bson.A{
			bson.M{"$count": "total_elements"},
			bson.M{"$addFields": bson.M{
				"total_pages": bson.M{
					"$floor": bson.M{
						"$divide": bson.A{
							"$total_elements",
							size,
						},
					},
				},
				"size": bson.M{
					"$min": bson.A{
						size,
						bson.M{"$subtract": bson.A{
							"$total_elements",
							size * page,
						}},
					},
				},
				"number": page,
			}},
			bson.M{"$addFields": bson.M{
				"empty": bson.M{
					"$eq": bson.A{
						"$size",
						0,
					},
				},
				"first": page == 0,
				"last": bson.M{
					"$eq": bson.A{
						page,
						bson.M{"$subtract": bson.A{
							"$total_pages",
							1,
						}},
					},
				},
			}},
		},
	}}
}
