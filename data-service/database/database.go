// Database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"log"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
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
	InsertHostData(hostData model.HostData) (*mongo.InsertOneResult, utils.AdvancedErrorInterface)
}

// MongoDatabase is a implementation
type MongoDatabase struct {
	// Config contains the dataservice global configuration
	// TODO: Should be removed?
	Config config.Configuration
	// Client contain the mongodb client
	Client *mongo.Client
}

// Init initializes the connection to the database
func (md *MongoDatabase) Init() {
	//Connect to mongodb
	md.ConnectToMongodb()
	log.Println("MongoDatabase is connected to MongoDB!", md.Config.Mongodb.URI)
}

// ArchiveHost archives tho host with hostname as hostname
func (md *MongoDatabase) ArchiveHost(hostname string) (*mongo.UpdateResult, utils.AdvancedErrorInterface) {
	if res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").UpdateOne(context.TODO(), bson.D{
		{"hostname", hostname},
		{"archived", false},
	}, bson.D{
		{"$set", bson.D{
			{"archived", true},
		}},
	}); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	} else {
		return res, nil
	}
}

// InsertHostData adds a new hostdata to the database
func (md *MongoDatabase) InsertHostData(hostData model.HostData) (*mongo.InsertOneResult, utils.AdvancedErrorInterface) {
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
