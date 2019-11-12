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

package migration

import (
	"context"
	"log"

	"github.com/amreo/ercole-services/utils"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

// ConnectToMongodb connects to the MongoDB and return the connection
func ConnectToMongodb(conf config.Mongodb) *mongo.Client {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(conf.URI)

	//Connect to MongoDB
	cl, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//Check the connection
	err = cl.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return cl
}

// Migrate migrate the client database
func Migrate(client *mongo.Database) {
	//NB: ALL OPERATIONS SHOULD BE IDEMPOTENT
	//THE RESULT OF
	//	Migrate(&db)
	//	Migrate(&db)
	//SHOULD BE EQUAL TO THE RESULT OF
	//	Migrate(&db)
	//AND POSSIBLY AVOID DESTRUCTIVE CHANGES
	UpdateDataSchemas(client)

	MigrateHostsSchema(client)
	MigrateClustersSchema(client)
	MigrateLicensesSchema(client)
	MigrateAlertsSchema(client)
	MigrateCurrentDatabasesSchema(client)
}

// MigrateHostsSchema create or update the hosts schema
func MigrateHostsSchema(client *mongo.Database) {
	//Create the collection
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "hosts") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "hosts"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	//Set the collection validator
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "hosts"},
		{"validator", bson.D{
			{"$jsonSchema", model.HostDataBsonValidatorRules},
		}},
		{"validationAction", "error"},
	}).Err(); err != nil {
		log.Panicln(err)
	}

	//index creations
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"hostname", 1},
		},
		Options: (&options.IndexOptions{
			PartialFilterExpression: bson.D{{"archived", false}},
		}).SetUnique(true),
	}); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"host_type", 1},
			{"hostname", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"hostname", 1},
			{"created_at", -1},
		},
	}); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"hostname", "text"},
			{"extra.databases.name", "text"},
			{"extra.databases.lunique_name", "text"},
			{"extra.clusters.name", "text"},
		},
		Options: &options.IndexOptions{
			Collation: &options.Collation{
				Locale: "simple",
			},
		},
	}); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"host_type", 1},
			{"extra.clusters.name", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"host_type", 1},
			{"extra.clusters.vms.hostname", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
}

// MigrateClustersSchema create or update the currentCluster schema
func MigrateClustersSchema(client *mongo.Database) {
	//Create the view
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "currentClusters") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "currentClusters"},
			{"viewOn", "hosts"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	//Set the view pipeline
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "currentClusters"},
		{"viewOn", "hosts"},
		{"pipeline", bson.A{
			bson.D{{"$match", bson.D{
				{"archived", false},
				{"host_type", "virtualization"},
			}}},
			bson.D{{"$unwind", bson.D{
				{"path", "$extra.clusters"},
			}}},
			bson.D{{"$project", bson.D{
				{"hostname", 1},
				{"environment", 1},
				{"location", 1},
				{"created_at", 1},
				{"cluster", "$extra.clusters"},
			}}},
		}},
	}).Err(); err != nil {
		log.Panicln(err)
	}
}

// MigrateLicensesSchema create or update the licenses schema
func MigrateLicensesSchema(client *mongo.Database) {
	//Create the collection
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "licenses") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "licenses"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	//Set the collection validator
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "licenses"},
		{"validator", bson.D{
			{"$jsonSchema", model.LicenseCountBsonValidatorRules},
		}},
	}).Err(); err != nil {
		log.Panicln(err)
	}
}

// MigrateAlertsSchema create or update the alerts schema
func MigrateAlertsSchema(client *mongo.Database) {
	//Create the collection
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "alerts") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "alerts"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	//Set the collection validator
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "alerts"},
		{"validator", bson.D{
			{"$jsonSchema", model.AlertBsonValidatorRules},
		}},
	}).Err(); err != nil {
		log.Panicln(err)
	}

	//index creations
	if _, err := client.Collection("licenses").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"hostname", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
}

// MigrateCurrentDatabasesSchema create or update the databases schema
func MigrateCurrentDatabasesSchema(client *mongo.Database) {
	//Create the collection
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "currentDatabases") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "currentDatabases"},
			{"viewOn", "hosts"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	//Set the view pipeline
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "currentDatabases"},
		{"viewOn", "hosts"},
		{"pipeline", bson.A{
			bson.D{{"$match", bson.D{
				{"archived", false},
				{"host_type", "oracledb"},
			}}},
			bson.D{{"$unwind", bson.D{
				{"path", "$extra.databases"},
			}}},
			bson.D{{"$project", bson.D{
				{"hostname", 1},
				{"environment", 1},
				{"location", 1},
				{"created_at", 1},
				{"database", "$extra.databases"},
			}}},
		}},
	}).Err(); err != nil {
		log.Panicln(err)
	}
}

// UpdateDataSchemas updates the schema of the data in the database
func UpdateDataSchemas(client *mongo.Database) {

}
