// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := inserPartitioningsSettingstNodes(db); err != nil {
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

func inserPartitioningsSettingstNodes(client *mongo.Database) error {
	collectionName := "nodes"

	nodes := getPartitioningsSettingsNodes()

	if _, err := client.Collection(collectionName).InsertMany(context.TODO(), nodes); err != nil {
		return err
	}

	return nil
}

func getPartitioningsSettingsNodes() []interface{} {
	return []interface{}{
		model.Node{
			Name: "Partitionings",
			Roles: []string{
				"admin",
				"read_databases",
				"write_databases",
			},
			Parent: "Oracle",
		},

		model.Node{
			Name: "Settings",
			Roles: []string{
				"admin",
			},
			Parent: "",
		},
	}
}
