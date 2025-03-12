// Copyright (c) 2023 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data

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
		if err := addMissingDatabases(db); err != nil {
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

func addMissingDatabases(client *mongo.Database) error {
	collectionName := "hosts"

	filter := bson.M{
		"$or": bson.A{
			bson.D{{Key: "features.oracle.database.unlistedRunningDatabases", Value: bson.D{{Key: "$exists", Value: true}}}},
			bson.D{{Key: "features.oracle.database.unretrievedDatabases", Value: bson.D{{Key: "$exists", Value: true}}}},
		},
	}

	pipeline := mongo.Pipeline{
		{{
			Key: "$addFields",
			Value: bson.M{
				"features.oracle.database.missingDatabases": bson.M{
					"$map": bson.M{
						"input": bson.M{
							"$filter": bson.M{
								"input": bson.M{
									"$setDifference": bson.A{
										bson.M{
											"$concatArrays": bson.A{
												bson.M{"$ifNull": bson.A{"$features.oracle.database.unretrievedDatabases", bson.A{}}},
												bson.M{"$ifNull": bson.A{"$features.oracle.database.unlistedRunningDatabases", bson.A{}}},
											},
										},
										[]string{},
									},
								},
								"cond": bson.M{
									"$ne": bson.A{"$$this", ""},
								},
							},
						},
						"as": "e",
						"in": bson.M{
							"name":           "$$e",
							"ignored":        false,
							"ignoredComment": nil,
						},
					},
				},
			},
		}},
		{{Key: "$unset", Value: "features.oracle.database.unretrievedDatabases"}},
		{{Key: "$unset", Value: "features.oracle.database.unlistedRunningDatabases"}},
	}

	if _, err := client.Collection(collectionName).UpdateMany(context.TODO(), filter, pipeline); err != nil {
		return err
	}

	return nil
}
