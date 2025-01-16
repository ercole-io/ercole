// Copyright (c) 2025 Sorint.lab S.p.A.
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
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := refactorLocationsRolesCollection(db); err != nil {
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

func refactorLocationsRolesCollection(client *mongo.Database) error {
	collectionName := "roles"

	filter := bson.M{"location": bson.M{"$exists": true}}

	pipeline := mongo.Pipeline{
		{{Key: "$set", Value: bson.D{{Key: "locations", Value: bson.A{"$location"}}}}},
		{{Key: "$unset", Value: "location"}},
	}

	if _, err := client.Collection(collectionName).UpdateMany(context.TODO(), filter, pipeline); err != nil {
		return err
	}

	return nil
}
