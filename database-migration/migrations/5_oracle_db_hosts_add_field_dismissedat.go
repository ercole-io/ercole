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

	"github.com/ercole-io/ercole/v2/utils"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := migrateHostsAddFieldDismissedAtArchivedFalse(db); err != nil {
			return err
		}
		if err := migrateHostsAddFieldDismissedAtArchivedTrueWithFalse(db); err != nil {
			return err
		}
		if err := migrateHostsAddFieldDismissedAtArchivedTrue(db); err != nil {
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

func migrateHostsAddFieldDismissedAtArchivedFalse(db *mongo.Database) error {
	collection := "hosts"
	ctx := context.TODO()

	filter := bson.M{"archived": false}
	update := bson.M{"$set": bson.M{"dismissedAt": nil}}

	_, err := db.Collection(collection).UpdateMany(ctx, filter, update)
	if err != nil {
		return utils.NewError(err, "Can't add new field 'dismissedAt' to hosts collection", collection)
	}

	return nil
}

func migrateHostsAddFieldDismissedAtArchivedTrueWithFalse(db *mongo.Database) error {
	collection := "hosts"
	ctx := context.TODO()
	filter := bson.M{"archived": false}
	update := bson.M{"$set": bson.M{"dismissedAt": nil}}

	var out map[string]interface{}

	cursor, err := db.Collection(collection).Find(
		ctx,
		filter,
	)
	if err != nil {
		return utils.NewError(err, "Can't find Hosts", collection)
	}

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&out); err != nil {
			return nil
		}

		filterUpd := bson.M{"hostname": out["hostname"]}

		if _, err := db.Collection(collection).UpdateMany(ctx, filterUpd, update); err != nil {
			return utils.NewError(err, "Can't add new field 'dismissedAt' to hosts collection")
		}
	}

	return nil
}

func migrateHostsAddFieldDismissedAtArchivedTrue(db *mongo.Database) error {
	collection := "hosts"
	ctx := context.TODO()
	filter := bson.M{"dismissedAt": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"dismissedAt": utils.MIN_TIME}}

	_, err := db.Collection(collection).UpdateMany(ctx, filter, update)
	if err != nil {
		return utils.NewError(err, "Can't add new field 'dismissedAt' to hosts collection", collection)
	}

	return nil
}
