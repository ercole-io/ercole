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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package migrations

import (
	"context"
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {

	newNode := model.Node{
		Name: "Capacity",
		Roles: []string{
			"admin",
			"read_databases",
			"write_databases",
		},
		Parent: "Oracle",
	}

	err := migrate.Register(func(db *mongo.Database) error {
		if err := insertNode(newNode, db); err != nil {
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

func insertNode(node model.Node, client *mongo.Database) error {
	collectionName := "nodes"

	if collectionNames, err := client.ListCollectionNames(context.TODO(), bson.D{{Key: "name", Value: collectionName}}); len(collectionNames) > 0 {
		if err != nil {
			return err
		}
	}

	if _, err := client.Collection(collectionName).InsertOne(context.TODO(), node); err != nil {
		return err
	}

	return nil
}
