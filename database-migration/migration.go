package migration

import (
	"context"

	"github.com/amreo/ercole-services/model"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

// Migrate migrate the client database
func Migrate(client *mongo.Database) {
	//NB: ALL OPERATIONS SHOULD BE IDEMPOTENT
	//THE RESULT OF
	//	Migrate(&db)
	//	Migrate(&db)
	//SHOULD BE EQUAL TO THE RESULT OF
	//	Migrate(&db)
	//AND POSSIBLY AVOID DESTRUCTIVE CHANGES
}

// MigrateHostsSchema create or update the hosts schema
func MigrateHostsSchema(client *mongo.Database) {
	//TODO: add data updater

	client.RunCommand(context.TODO(), bson.D{
		{"collMod", "hosts"},
		{"validator", bson.D{
			{"$jsonSchema", model.HostDataBsonValidatorRules},
		}},
	})
}
