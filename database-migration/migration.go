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

package migration

import (
	"context"
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectToMongodb connects to the MongoDB and return the connection
func ConnectToMongodb(log *logrus.Logger, conf config.Mongodb) *mongo.Client {
	var err error

	//Set client options
	clientOptions := options.Client().ApplyURI(conf.URI)

	//Connect to MongoDB
	cl, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = cl.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Can't connect to the database!\n%s\n", err)
	}

	return cl
}

// Migrate migrate the client database
func Migrate(log *logrus.Logger, client *mongo.Database) {
	//NB: ALL OPERATIONS SHOULD BE IDEMPOTENT
	//THE RESULT OF
	//	Migrate(&db)
	//	Migrate(&db)
	//SHOULD BE EQUAL TO THE RESULT OF
	//	Migrate(&db)
	//AND POSSIBLY AVOID DESTRUCTIVE CHANGES
	UpdateDataSchemas(log, client)

	MigrateHostsSchema(log, client)
	// MigrateClustersSchema(log, client)
	MigrateAlertsSchema(log, client)
	MigratePatchingFunctionsSchema(log, client)
	// MigrateCurrentDatabasesSchema(log, client)
	MigrateOracleDatabaseAgreementsSchema(log, client)
	MigrateOracleDatabaseLicenseTypes(log, client)
}

// MigrateHostsSchema create or update the hosts schema
func MigrateHostsSchema(log *logrus.Logger, client *mongo.Database) {
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
			{"$jsonSchema", bson.M{}},
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
			{"hostname", 1},
			{"createdAt", -1},
		},
	}); err != nil {
		log.Panicln(err)
	}
	// if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
	// 	Keys: bson.D{
	// 		{"archived", 1},
	// 		{"hostname", "text"},
	// 		{"extra.databases.name", "text"},
	// 		{"extra.databases.unique_name", "text"},
	// 		{"extra.clusters.name", "text"},
	// 	},
	// 	Options: &options.IndexOptions{
	// 		Collation: &options.Collation{
	// 			Locale: "simple",
	// 		},
	// 		Weights: map[string]interface{}{
	// 			"hostname":                    10,
	// 			"extra.databases.name":        7,
	// 			"extra.databases.unique_name": 6,
	// 			"extra.clusters.name":         7,
	// 		},
	// 	},
	// }); err != nil {
	// 	log.Panicln(err)
	// }
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"clusters.name", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"clusters.vms.hostname", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
}

// // MigrateClustersSchema create or update the currentCluster schema
// func MigrateClustersSchema(client *mongo.Database) {
// 	//Create the view
// 	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
// 		log.Panicln(err)
// 	} else if !utils.Contains(cols, "currentClusters") {
// 		if err := client.RunCommand(context.TODO(), bson.D{
// 			{"create", "currentClusters"},
// 			{"viewOn", "hosts"},
// 		}).Err(); err != nil {
// 			log.Panicln(err)
// 		}
// 	}

// 	//Set the view pipeline
// 	if err := client.RunCommand(context.TODO(), bson.D{
// 		{"collMod", "currentClusters"},
// 		{"viewOn", "hosts"},
// 		{"pipeline", bson.A{
// 			bson.D{{"$match", bson.D{
// 				{"archived", false},
// 				{"host_type", "virtualization"},
// 			}}},
// 			bson.D{{"$unwind", bson.D{
// 				{"path", "$extra.clusters"},
// 			}}},
// 			bson.D{{"$project", bson.D{
// 				{"hostname", 1},
// 				{"environment", 1},
// 				{"location", 1},
// 				{"created_at", 1},
// 				{"cluster", "$extra.clusters"},
// 			}}},
// 		}},
// 	}).Err(); err != nil {
// 		log.Panicln(err)
// 	}
// }

// MigrateAlertsSchema create or update the alerts schema
func MigrateAlertsSchema(log *logrus.Logger, client *mongo.Database) {
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
}

// // MigrateCurrentDatabasesSchema create or update the databases schema
// func MigrateCurrentDatabasesSchema(client *mongo.Database) {
// 	//Create the collection
// 	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
// 		log.Panicln(err)
// 	} else if !utils.Contains(cols, "currentDatabases") {
// 		if err := client.RunCommand(context.TODO(), bson.D{
// 			{"create", "currentDatabases"},
// 			{"viewOn", "hosts"},
// 		}).Err(); err != nil {
// 			log.Panicln(err)
// 		}
// 	}

// 	//Set the view pipeline
// 	if err := client.RunCommand(context.TODO(), bson.D{
// 		{"collMod", "currentDatabases"},
// 		{"viewOn", "hosts"},
// 		{"pipeline", bson.A{
// 			bson.D{{"$match", bson.D{
// 				{"archived", false},
// 				{"host_type", "oracledb"},
// 			}}},
// 			bson.D{{"$unwind", bson.D{
// 				{"path", "$extra.databases"},
// 			}}},
// 			bson.D{{"$addFields", bson.D{
// 				{"extra.databases.ha", bson.D{
// 					{"$or", bson.A{
// 						"$info.sun_cluster",
// 						"$info.veritas_cluster",
// 						"$info.oracle_cluster",
// 						"$info.aix_cluster",
// 					}},
// 				}},
// 			}}},
// 			bson.D{{"$project", bson.D{
// 				{"hostname", 1},
// 				{"environment", 1},
// 				{"location", 1},
// 				{"created_at", 1},
// 				{"database", "$extra.databases"},
// 			}}},
// 		}},
// 	}).Err(); err != nil {
// 		log.Panicln(err)
// 	}
// }

// UpdateDataSchemas updates the schema of the data in the database
func UpdateDataSchemas(log *logrus.Logger, client *mongo.Database) {

}

// MigratePatchingFunctionsSchema create or update the patching_functions schema
func MigratePatchingFunctionsSchema(log *logrus.Logger, client *mongo.Database) {
	//Create the collection
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "patching_functions") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "patching_functions"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	//Set the collection validator
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "patching_functions"},
		{"validator", bson.D{
			{"$jsonSchema", bson.M{}},
		}},
		{"validationAction", "error"},
	}).Err(); err != nil {
		log.Panicln(err)
	}

	//index creations
	if _, err := client.Collection("patching_functions").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"hostname", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
}

// MigrateOracleDatabaseAgreementsSchema create or update the oracle_database_agreements schema
func MigrateOracleDatabaseAgreementsSchema(log *logrus.Logger, client *mongo.Database) {
	collection := "oracle_database_agreements"

	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, collection) {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", collection},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	if _, err := client.Collection(collection).
		Indexes().
		CreateMany(context.TODO(),
			[]mongo.IndexModel{
				{

					Keys: bson.D{
						{"agreementID", 1},
					},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{
						{"licenseTypes._id", 1},
					},
					Options: options.Index().SetUnique(true),
				}},
		); err != nil {
		log.Panicln(err)
	}
}

// MigrateOracleDatabaseAgreementsSchema create or update the oracle_database_agreements schema
func MigrateOracleDatabaseLicenseTypes(log *logrus.Logger, client *mongo.Database) {
	collection := "oracle_database_license_types"

	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, collection) {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", collection},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}
}
