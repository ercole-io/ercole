// Copyright (c) 2023 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	Init()

	// Oracle Cloud Configuration
	GetOciProfiles(hidePrivateKey bool) ([]model.OciProfile, error)
	GetMapOciProfiles() (map[primitive.ObjectID]model.OciProfile, error)
	AddOciProfile(profile model.OciProfile) error
	DeleteOciProfile(id primitive.ObjectID) error
	UpdateOciProfile(profile model.OciProfile) error
	GetErcoleDatabases() ([]dto.ErcoleDatabase, error)
	GetErcoleActiveDatabases() ([]dto.ErcoleDatabase, error)
	AddOciObjects(objects model.OciObjects) error
	GetOciObjects() ([]model.OciObjects, error)
	DeleteOldOciObjects(dateFrom time.Time) error
	GetOciRecommendationsByProfiles(profileIDs []string) ([]model.OciRecommendation, error)
	GetOciRecommendations(profileIDs []string) ([]model.OciRecommendation, error)
	AddOciRecommendations(ercoleRecommendations []model.OciRecommendation) error
	AddOciRecommendationErrors(ociRecommendationErrors []model.OciRecommendationError) error
	GetOciRecommendationErrors(seqNum uint64) ([]model.OciRecommendationError, error)
	GetLastOciSeqValue() (uint64, error)
	DeleteOldOciRecommendations(dateFrom time.Time) error
	DeleteOldOciRecommendationErrors(dateFrom time.Time) error
	SelectOciProfile(profileId string, selected bool) error
	GetSelectedOciProfiles() ([]string, error)
	GetAwsProfiles(hidePrivateKey bool) ([]model.AwsProfile, error)
	GetMapAwsProfiles() (map[primitive.ObjectID]model.AwsProfile, error)
	AddAwsProfile(profile model.AwsProfile) error
	DeleteAwsProfile(id primitive.ObjectID) error
	UpdateAwsProfile(profile model.AwsProfile) error
	SelectAwsProfile(profileId string, selected bool) error
	GetSelectedAwsProfiles() ([]primitive.ObjectID, error)
	AddAwsObject(m interface{}, collection string) error
	AddAwsObjects(m []interface{}, collection string) error
	GetLastAwsSeqValue() (uint64, error)
	GetAwsRecommendationsByProfiles(profileIDs []primitive.ObjectID) ([]model.AwsRecommendation, error)
	GetAwsRecommendationsBySeqValue(seqValue uint64) ([]model.AwsRecommendation, error)
	GetAwsObjectsBySeqValue(seqValue uint64) ([]model.AwsObject, error)
	GetAzureProfiles(hidePrivateKey bool) ([]model.AzureProfile, error)
	GetMapAzureProfiles() (map[primitive.ObjectID]model.AzureProfile, error)
	AddAzureProfile(profile model.AzureProfile) error
	DeleteAzureProfile(id primitive.ObjectID) error
	UpdateAzureProfile(profile model.AzureProfile) error
	SelectAzureProfile(profileId string, selected bool) error
	GetSelectedAzureProfiles() ([]primitive.ObjectID, error)
	GetLastAwsRDSSeqValue() (uint64, error)
	GetAwsRDS() ([]model.AwsRDS, error)

	AddGcpProfile(profile model.GcpProfile) error
	GetActiveGcpProfiles() ([]model.GcpProfile, error)
	ListGcpProfiles() ([]model.GcpProfile, error)
	SelectGcpProfile(id primitive.ObjectID, selected bool) error
	UpdateGcpProfile(id primitive.ObjectID, profile model.GcpProfile) error
	RemoveGcpProfile(id primitive.ObjectID) error

	GetLastGcpSeqValue() (uint64, error)
	ListGcpRecommendationsByProfiles(profileIDs []primitive.ObjectID) ([]model.GcpRecommendation, error)
}

// MongoDatabase is a implementation
type MongoDatabase struct {
	Config  config.Configuration
	Client  *mongo.Client
	TimeNow func() time.Time
	Log     logger.Logger
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
