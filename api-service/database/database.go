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
	"regexp"
	"strings"
	"time"

	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// SearchCurrentHosts search current hosts
	SearchCurrentHosts(full bool, keywords []string) ([]interface{}, utils.AdvancedErrorInterface)
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

// SearchCurrentHosts search current hosts
func (md *MongoDatabase) SearchCurrentHosts(full bool, keywords []string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	var quotedKeywords []string
	for _, k := range keywords {
		quotedKeywords = append(quotedKeywords, regexp.QuoteMeta(k))
	}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		bson.A{
			bson.D{{"$match", bson.D{
				{"archived", false},
				{"$or", bson.A{
					bson.D{{"hostname", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"extra.databases.name", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"extra.databases.unique_name", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"extra.clusters.name", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
				}},
			}}},
			optionalStep(!full, bson.D{{"$project", bson.D{
				{"hostname", true},
				{"environment", true},
				{"host_type", true},
				{"cluster", ""},
				{"physical_host", ""},
				{"created_at", true},
				{"databases", true},
				{"os", "$info.os"},
				{"kernel", "$info.kernel"},
				{"oracle_cluster", "$info.oracle_cluster"},
				{"sun_cluster", "$info.sun_cluster"},
				{"veritas_cluster", "$info.veritas_cluster"},
				{"virtual", "$info.virtual"},
				{"type", "$info.type"},
				{"cpu_threads", "$info.cpu_threads"},
				{"cpu_cores", "$info.cpu_cores"},
				{"socket", "$info.socket"},
				{"mem_total", "$info.memory_total"},
				{"swap_total", "$info.swap_total"},
				{"cpu_model", "$info.cpu_model"},
			}}}),
		},
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

func optionalStep(optional bool, step bson.D) bson.D {
	if optional {
		return step
	}
	return bson.D{{"$skip", 0}}
}
