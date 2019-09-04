package migration

import (
	"context"

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
	client.RunCommand(context.TODO(), bson.D{
		{"collMod", "hosts"},
		{"validator", bson.D{
			{"$jsonSchema", model.HostDataBsonValidatorRules},
		}},
	})
	client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"hostname", 1},
		},
	})
	client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"hostname", 1},
		},
		Options: &options.IndexOptions{
			Unique:                  &TRUE,
			PartialFilterExpression: bson.D{{"archived", false}},
		},
	})
	client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"host_type", 1},
			{"hostname", 1},
		},
	})
	client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"hostname", 1},
		},
	})
	client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"archived", 1},
			{"hostname", "text"},
			{"extra_info.databases.name", "text"},
			{"extra_info.databases.lunique_name", "text"},
			{"extra_info.clusters.name", "text"},
		},
		Options: &options.IndexOptions{
			Collation: &options.Collation{
				Locale:          "simple",
				CaseLevel:       false,
				NumericOrdering: true,
				Strength:        1,
			},
		},
	})
}

func MigrateClustersSchema(client *mongo.Database) {
	// client.RunCommand(context.TODO(), bson.D{
	// 	{"collMod", "clusters"},
	// 	{"validator", bson.D{
	// 		{"$jsonSchema", model.ClusterInfoBsonValidatorRules},
	// 	}},
	// })
	client.RunCommand(context.TODO(), bson.D{
		{"create", "CurrentClusters"},
		{"viewOn", "hosts"},
		{"pipeline", bson.D{
			//TODO: Write pipeline!
		}},
	})
}

func MigrateLicensesSchema(client *mongo.Database) {
	client.RunCommand(context.TODO(), bson.D{
		{"collMod", "licenses"},
		{"validator", bson.D{
			{"$jsonSchema", model.LicenseCountBsonValidatorRules},
		}},
	})
}

func MigrateAlertsSchema(client *mongo.Database) {
	client.RunCommand(context.TODO(), bson.D{
		{"collMod", "alerts"},
		{"validator", bson.D{
			{"$jsonSchema", model.AlertBsonValidatorRules},
		}},
	})
	client.Collection("licenses").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"hostname", 1},
		},
	})
}

func MigrateCurrentDatabasesSchema(client *mongo.Database) {
	client.RunCommand(context.TODO(), bson.D{
		{"create", "CurrentDatabases"},
		{"viewOn", "hosts"},
		{"pipeline", bson.D{
			//TODO: Write pipeline!
		}},
	})
}
