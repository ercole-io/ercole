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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
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

	AckOldAlerts(dueDays int) (*mongo.UpdateResult, error)
	RemoveOldAlerts(dueDays int) (*mongo.DeleteResult, error)
	FindAlertsByDate(startDate, endDate time.Time) ([]model.Alert, error) 
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

func (md *MongoDatabase) ReadConfig() (*config.Configuration, error) {
	ctx := context.TODO()

	conf := config.Configuration{}
	if err := md.Client.Database(md.Config.Mongodb.DBName).Collection("config").FindOne(ctx, bson.D{}).Decode(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
