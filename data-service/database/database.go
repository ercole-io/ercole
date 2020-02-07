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

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// ArchiveHost archives tho host with hostname as hostname
	ArchiveHost(hostname string) (*mongo.UpdateResult, utils.AdvancedErrorInterface)
	// InsertHostData adds a new hostdata to the database
	InsertHostData(hostData interface{}) (*mongo.InsertOneResult, utils.AdvancedErrorInterface)
	// FindOldCurrentHost return the list of current hosts that haven't sent hostdata after time t
	FindOldCurrentHosts(t time.Time) ([]string, utils.AdvancedErrorInterface)
	// FindOldArchivedHosts return the list of archived hosts older than t
	FindOldArchivedHosts(t time.Time) ([]primitive.ObjectID, utils.AdvancedErrorInterface)
	// DeleteHostData delete the hostdata
	DeleteHostData(id primitive.ObjectID) utils.AdvancedErrorInterface
	// FindPatchingFunction find the the patching function associated to the hostname in the database
	FindPatchingFunction(hostname string) (model.PatchingFunction, utils.AdvancedErrorInterface)
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

// ArchiveHost archives tho host with hostname as hostname
func (md *MongoDatabase) ArchiveHost(hostname string) (*mongo.UpdateResult, utils.AdvancedErrorInterface) {
	if res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").UpdateOne(context.TODO(), bson.M{
		"Hostname": hostname,
		"Archived": false,
	}, bson.M{
		"$set": bson.M{
			"Archived": true,
		},
	}); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	} else {
		return res, nil
	}
}

// InsertHostData adds a new hostdata to the database
func (md *MongoDatabase) InsertHostData(hostData interface{}) (*mongo.InsertOneResult, utils.AdvancedErrorInterface) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").InsertOne(context.TODO(), hostData)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return res, nil
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

// FindOldCurrentHosts return the list of current hosts that haven't sent hostdata after time t
func (md *MongoDatabase) FindOldCurrentHosts(t time.Time) ([]string, utils.AdvancedErrorInterface) {
	//Get the list of old current hosts
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"Hostname",
		bson.M{
			"Archived": false,
			"CreatedAt": bson.M{
				"$lt": t,
			},
		})
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Convert the slice of interface{} to []string
	var hosts []string
	for _, val := range values {
		hosts = append(hosts, val.(string))
	}

	//Return it
	return hosts, nil
}

// FindOldArchivedHosts return the list of archived hosts older than t
func (md *MongoDatabase) FindOldArchivedHosts(t time.Time) ([]primitive.ObjectID, utils.AdvancedErrorInterface) {
	//Get the list of old archived hosts
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"_id",
		bson.M{
			"Archived": true,
			"CreatedAt": bson.M{
				"$lt": t,
			},
		})
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Convert the slice of interface{} to []primitive.ObjectID
	var ids []primitive.ObjectID
	for _, val := range values {
		ids = append(ids, val.(primitive.ObjectID))
	}

	//Return it
	return ids, nil
}

// DeleteHostData delete the hostdata
func (md *MongoDatabase) DeleteHostData(id primitive.ObjectID) utils.AdvancedErrorInterface {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").DeleteOne(
		context.TODO(),
		bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return nil
}

// FindPatchingFunction find the the patching function associated to the hostname in the database
func (md *MongoDatabase) FindPatchingFunction(hostname string) (model.PatchingFunction, utils.AdvancedErrorInterface) {
	var out model.PatchingFunction

	//Find the hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("patching_functions").Find(context.TODO(), bson.M{
		"Hostname": hostname,
	})
	if err != nil {
		return model.PatchingFunction{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return model.PatchingFunction{}, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return model.PatchingFunction{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}
