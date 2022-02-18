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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package migrations

import (
	"context"
	"fmt"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/utils"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := migrateHostsSchema(db); err != nil {
			return err
		}
		if err := migrateAlertsSchema(db); err != nil {
			return err
		}
		if err := migrateOracleDatabaseAgreementsSchema(db); err != nil {
			return err
		}
		if err := migrateOracleDatabaseLicenseTypes(db); err != nil {
			return err
		}

		return nil

	}, func(db *mongo.Database) error {
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

// migrateHostsSchema create or update the hosts schema
func migrateHostsSchema(client *mongo.Database) error {
	//Create the collection
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		return err
	} else if !utils.Contains(cols, "hosts") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{Key: "create", Value: "hosts"},
		}).Err(); err != nil {
			return err
		}
	}

	//Set the collection validator
	if err := client.RunCommand(context.TODO(), bson.D{
		{Key: "collMod", Value: "hosts"},
		{Key: "validator", Value: bson.D{
			{Key: "$jsonSchema", Value: bson.M{}},
		}},
		{Key: "validationAction", Value: "error"},
	}).Err(); err != nil {
		return err
	}

	//index creations
	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "archived", Value: 1},
			{Key: "hostname", Value: 1},
		},
		Options: (&options.IndexOptions{
			PartialFilterExpression: bson.D{{Key: "archived", Value: false}},
		}).SetUnique(true),
	}); err != nil {
		return err
	}

	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "hostname", Value: 1},
			{Key: "createdAt", Value: -1},
		},
	}); err != nil {
		return err
	}

	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "archived", Value: 1},
			{Key: "clusters.name", Value: 1},
		},
	}); err != nil {
		return err
	}

	if _, err := client.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "archived", Value: 1},
			{Key: "clusters.vms.hostname", Value: 1},
		},
	}); err != nil {
		return err
	}

	return nil
}

// migrateAlertsSchema create or update the alerts schema
func migrateAlertsSchema(client *mongo.Database) error {
	//Create the collection
	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		return err
	} else if !utils.Contains(cols, "alerts") {
		if err := client.RunCommand(context.TODO(), bson.D{
			{Key: "create", Value: "alerts"},
		}).Err(); err != nil {
			return err
		}
	}

	return nil
}

// migrateOracleDatabaseAgreementsSchema create or update the oracle_database_agreements schema
func migrateOracleDatabaseAgreementsSchema(client *mongo.Database) error {
	collection := "oracle_database_agreements"

	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		return err
	} else if !utils.Contains(cols, collection) {
		if err := client.RunCommand(context.TODO(), bson.D{
			{Key: "create", Value: collection},
		}).Err(); err != nil {
			return err
		}
	}

	if _, err := client.Collection(collection).
		Indexes().
		CreateMany(context.TODO(),
			[]mongo.IndexModel{
				{

					Keys: bson.D{
						{Key: "agreementID", Value: 1},
					},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{
						{Key: "licenseTypes._id", Value: 1},
					},
					Options: options.Index().SetUnique(true),
				}},
		); err != nil {
		return err
	}

	return nil
}

// MigrateOracleDatabaseAgreementsSchema create or update the oracle_database_agreements schema
func migrateOracleDatabaseLicenseTypes(client *mongo.Database) error {
	collection := "oracle_database_license_types"

	if cols, err := client.ListCollectionNames(context.TODO(), bson.D{}); err != nil {
		return err
	} else if !utils.Contains(cols, collection) {
		if err := client.RunCommand(context.TODO(), bson.D{
			{Key: "create", Value: collection},
		}).Err(); err != nil {
			return err
		}
	}

	return nil
}
