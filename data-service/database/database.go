// Database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"log"

	"github.com/amreo/ercole-hostdata-dataservice/config"
	"github.com/amreo/ercole-hostdata-dataservice/model"
	"github.com/amreo/ercole-hostdata-dataservice/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseInterface is a interface that wrap methods used to perform CRUD operations in the mongodb database
type MongoDatabaseInterface interface {
	// Init initializes the connection to the database
	Init()
	// ArchiveHost archives tho host with hostname as hostname
	ArchiveHost(hostname string) (*mongo.UpdateResult, utils.AdvancedError)
	// InsertHostData adds a new hostdata to the database
	InsertHostData(hostData model.HostData) (*mongo.InsertOneResult, utils.AdvancedError)
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
func (this *MongoDatabase) Init() {
	//Connect to mongodb
	this.ConnectToMongodb()
	log.Println("MongoDatabase is connected to MongoDB!", this.Config.Mongodb.URI)
}

// ArchiveHost archives tho host with hostname as hostname
func (this *MongoDatabase) ArchiveHost(hostname string) (*mongo.UpdateResult, utils.AdvancedError) {
	res, err := this.Client.Database(this.Config.Mongodb.DBName).Collection("hosts").UpdateOne(context.TODO(), bson.D{
		{"hostname", hostname},
		{"archived", false},
	}, bson.D{
		{"$set", bson.D{
			{"archived", true},
		}},
	})
	return res, utils.NewAdvancedError(err, "DB ERROR")
}

// InsertHostData adds a new hostdata to the database
func (this *MongoDatabase) InsertHostData(hostData model.HostData) (*mongo.InsertOneResult, utils.AdvancedError) {
	res, err := this.Client.Database(this.Config.Mongodb.DBName).Collection("hosts").InsertOne(context.TODO(), hostData)
	if err != nil {
		return nil, utils.NewAdvancedError(err, "DB ERROR")
	}
	return res, utils.AdvancedError{}
}

// ConnectMongodb connects to the MongoDB and return the connection
func (this *MongoDatabase) ConnectToMongodb() {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(this.Config.Mongodb.URI)

	//Connect to MongoDB
	this.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = this.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
