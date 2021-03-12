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

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// FindHostData find a host data
	FindHostData(id primitive.ObjectID) (model.HostDataBE, error)
	// FindMostRecentHostDataOlderThan return the most recest hostdata that is older than t
	FindMostRecentHostDataOlderThan(hostname string, t time.Time) (model.HostDataBE, error)
	// InsertAlert inserr the alert in the database
	InsertAlert(alert model.Alert) (*mongo.InsertOneResult, error)
	// ExistNoDataAlertByHost return true if the host has associated a new NO_DATA alert
	ExistNoDataAlertByHost(hostname string) (bool, error)
}

// MongoDatabase is a implementation
type MongoDatabase struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Client contain the mongodb client
	Client *mongo.Client
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
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

// FindHostData find a host data
func (md *MongoDatabase) FindHostData(id primitive.ObjectID) (model.HostDataBE, error) {
	//Find the hostdata
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").FindOne(context.TODO(), bson.M{
		"_id": id,
	})
	if res.Err() != nil {
		return model.HostDataBE{}, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
	}

	//Decode the data

	var out model.HostDataBE
	if err := res.Decode(&out); err != nil {
		return model.HostDataBE{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Return it!
	return out, nil
}

//TODO RM?
// FindMostRecentHostDataOlderThan return the most recest hostdata that is older than t
func (md *MongoDatabase) FindMostRecentHostDataOlderThan(hostname string, t time.Time) (model.HostDataBE, error) {
	var out model.HostDataBE

	//Find the most recent HostData older than t
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"hostname":  hostname,
				"createdAt": mu.QOLessThan(t),
			}),
			mu.APSort(bson.M{
				"createdAt": -1,
			}),
			mu.APLimit(1),
		),
	)
	if err != nil {
		return model.HostDataBE{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return model.HostDataBE{}, nil
	}

	if err := cur.Decode(&out); err != nil {
		return model.HostDataBE{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}

// InsertAlert insert the alert in the database
func (md *MongoDatabase) InsertAlert(alert model.Alert) (*mongo.InsertOneResult, error) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").InsertOne(context.TODO(), alert)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return res, nil
}

// ExistNoDataAlertByHost return true if the host has associated a new NO_DATA alert
func (md *MongoDatabase) ExistNoDataAlertByHost(hostname string) (bool, error) {
	//Count the number of new NO_DATA alerts associated to the host
	val, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").CountDocuments(context.TODO(), bson.M{
		"alertCode":          model.AlertCodeNoData,
		"alertStatus":        model.AlertStatusNew,
		"otherInfo.hostname": hostname,
	}, &options.CountOptions{
		Limit: utils.Intptr(1),
	})
	if err != nil {
		return false, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Return true if the count > 0
	return val > 0, nil
}
