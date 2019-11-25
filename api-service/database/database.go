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
	SearchCurrentHosts(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]interface{}, utils.AdvancedErrorInterface)
}

// MongoDatabase is a implementation
type MongoDatabase struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Client contain the mongodb client
	Client *mongo.Client
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
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

func optionalStep(optional bool, step bson.D) bson.D {
	if optional {
		return step
	}
	return bson.D{{"$skip", 0}}
}

func optionalSortingStep(sortBy string, sortDesc bool) bson.D {
	if sortBy == "" {
		return bson.D{{"$skip", 0}}
	}

	sortOrder := 0
	if sortDesc {
		sortOrder = -1
	} else {
		sortOrder = 1
	}

	return bson.D{{"$sort", bson.D{
		{sortBy, sortOrder},
	}}}
}

func optionalPagingStep(page int, size int) bson.D {
	if page == -1 || size == -1 {
		return bson.D{{"$skip", 0}}
	}

	return bson.D{{"$facet", bson.D{
		{"content", bson.A{
			bson.D{{"$skip", page * size}},
			bson.D{{"$limit", size}},
		}},
		{"metadata", bson.A{
			bson.D{{"$count", "total_elements"}},
			bson.D{{"$addFields", bson.D{
				{"total_pages", bson.D{
					{"$floor", bson.D{
						{"$divide", bson.A{
							"$total_elements",
							size,
						}},
					}},
				}},
				{"size", bson.D{
					{"$min", bson.A{
						size,
						bson.D{{"$subtract", bson.A{
							"$total_elements",
							size * page,
						}}},
					}},
				}},
				{"number", page},
			}}},
			bson.D{{"$addFields", bson.D{
				{"empty", bson.D{
					{"$eq", bson.A{
						"$size",
						0,
					}},
				}},
				{"first", page == 0},
				{"last", bson.D{
					{"$eq", bson.A{
						page,
						bson.D{{"$subtract", bson.A{
							"$total_pages",
							1,
						}}},
					}},
				}},
			}}},
		}},
	}}}
}
