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
	"github.com/ercole-io/ercole/v2/utils"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := createGroupsCollection(db); err != nil {
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

func createGroupsCollection(client *mongo.Database) error {
	collectionName := "groups"
	ctx := context.TODO()

	if cols, err := client.ListCollectionNames(ctx, bson.D{}); err != nil {
		return err
	} else if !utils.Contains(cols, collectionName) {
		if err := client.RunCommand(ctx, bson.D{
			{Key: "create", Value: collectionName},
		}).Err(); err != nil {
			return err
		}

		_, err := client.Collection(collectionName).Indexes().
			CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "name", Value: 1},
				},
				Options: (&options.IndexOptions{}).SetUnique(true),
			})
		if err != nil {
			return err
		}

		if _, err := client.Collection(collectionName).InsertMany(ctx, getDefaultGroups()); err != nil {
			return err
		}
	}

	return nil
}

func getDefaultGroups() []interface{} {
	return []interface{}{
		model.Group{
			Name:        "admin",
			Description: "Admin group",
			Roles:       []string{"admin"},
		},
		model.Group{
			Name:        "procurement",
			Description: "procurement",
			Roles: []string{
				"read_dashboard",
				"read_license_contract",
				"write_license_contract",
				"read_license_compliance",
				"write_license_compliance",
				"read_license_used",
				"read_alert_type_license",
			},
		},
		model.Group{
			Name:        "dba",
			Description: "dba",
			Roles: []string{
				"read_dashboard",
				"write_dashboard",
				"read_host",
				"write_host",
				"read_databases",
				"write_databases",
				"read_hypervisor",
				"write_hypervisor",
				"read_exadata",
				"write_exadata",
				"read_alert",
				"write_alert",
				"read_cloud",
			},
		},
	}
}
