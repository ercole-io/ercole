package migration

import (
	"context"
	"log"

	"github.com/amreo/ercole-services/utils"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/amreo/ercole-services/model"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

var TRUE bool = true

// Migrate migrate the client database
func Migrate(client *mongo.Database) {
	//NB: ALL OPERATIONS SHOULD BE IDEMPOTENT
	//THE RESULT OF
	//	Migrate(&db)
	//	Migrate(&db)
	//SHOULD BE EQUAL TO THE RESULT OF
	//	Migrate(&db)
	//AND POSSIBLY AVOID DESTRUCTIVE CHANGES
	UpdateData(client)

	MigrateHostsSchema(client)
	MigrateClustersSchema(client)
	MigrateLicensesSchema(client)
	MigrateAlertsSchema(client)
	MigrateCurrentDatabasesSchema(client)
}

func UpdateData(client *mongo.Database) {

}

// MigrateHostsSchema create or update the hosts schema
func MigrateHostsSchema(client *mongo.Database) {
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "hosts") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "hosts"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}

	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "hosts"},
		{"validator", bson.D{
			{"$jsonSchema", model.HostDataBsonValidatorRules},
		}},
		{"validationAction", "error"},
	}).Err(); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"hostname", 1},
		},
		Options: &options.IndexOptions{
			Unique:                  &TRUE,
			PartialFilterExpression: bson.D{{"archived", false}},
		},
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
}

func MigrateClustersSchema(client *mongo.Database) {
	// client.RunCommand(context.TODO(), bson.D{
	// 	{"collMod", "clusters"},
	// 	{"validator", bson.D{
	// 		{"$jsonSchema", model.ClusterInfoBsonValidatorRules},
	// 	}},
	// })
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

func MigrateLicensesSchema(client *mongo.Database) {
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "licenses") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "licenses"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "licenses"},
		{"validator", bson.D{
			{"$jsonSchema", model.LicenseCountBsonValidatorRules},
		}},
	}).Err(); err != nil {
		log.Panicln(err)
	}
}

func MigrateAlertsSchema(client *mongo.Database) {
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		log.Panicln(err)
	} else if !utils.Contains(cols, "alerts") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{"create", "alerts"},
		}).Err(); err != nil {
			log.Panicln(err)
		}
	}
	if err := client.RunCommand(context.TODO(), bson.D{
		{"collMod", "alerts"},
		{"validator", bson.D{
			{"$jsonSchema", model.AlertBsonValidatorRules},
		}},
	}).Err(); err != nil {
		log.Panicln(err)
	}
	if _, err := client.Collection("licenses").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"hostname", 1},
		},
	}); err != nil {
		log.Panicln(err)
	}
}

func MigrateCurrentDatabasesSchema(client *mongo.Database) {
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
