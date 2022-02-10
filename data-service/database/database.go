// Copyright (c) 2021 Sorint.lab S.p.A.
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

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
)

type MongoDatabaseInterface interface {
	Init()
	DismissHost(hostname string) error
	InsertHostData(hostData model.HostDataBE) error
	GetCurrentHostnames() ([]string, error)
	// FindOldCurrentHostnames return the list of current hosts names that haven't sent hostdata after time t
	FindOldCurrentHostnames(t time.Time) ([]string, error)
	// FindOldCurrentHostdata return the list of current hosts that haven't sent hostdata after time t
	FindOldCurrentHostdata(t time.Time) ([]model.HostDataBE, error)
	// FindOldArchivedHosts return the list of archived hosts older than t
	FindOldArchivedHosts(t time.Time) ([]primitive.ObjectID, error)
	DeleteHostData(id primitive.ObjectID) error
	HistoricizeLicensesCompliance(licenses []dto.LicenseCompliance) error

	DeleteNoDataAlertByHost(hostname string) error
	DeleteAllNoDataAlerts() error
	// FindMostRecentHostDataOlderThan return the most recest hostdata that is older than t
	FindMostRecentHostDataOlderThan(hostname string, t time.Time) (*model.HostDataBE, error)
	GetHostnames() ([]string, error)
}

type MongoDatabase struct {
	Config  config.Configuration
	Client  *mongo.Client
	TimeNow func() time.Time
	Log     logger.Logger
}

func (md *MongoDatabase) Init() {
	md.ConnectToMongodb()
	md.Log.Debug("MongoDatabase is connected to MongoDB! ", md.Config.Mongodb.URI)
}

func (md *MongoDatabase) ConnectToMongodb() {
	var err error

	clientOptions := options.Client().ApplyURI(md.Config.Mongodb.URI)

	md.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		md.Log.Fatal(err)
	}

	err = md.Client.Ping(context.TODO(), nil)
	if err != nil {
		md.Log.Fatal(err)
	}
}
