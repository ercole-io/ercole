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
		if err := createRolesCollection(db); err != nil {
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

func createRolesCollection(client *mongo.Database) error {
	collectionName := "roles"
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

		if _, err := client.Collection(collectionName).InsertMany(ctx, getDefaultRoles()); err != nil {
			return err
		}
	}

	return nil
}

func getDefaultRoles() []interface{} {
	return []interface{}{
		model.Role{
			Name:        "admin",
			Description: "Admin role",
			Locations:   model.AllLocations,
			Permission:  model.AdminPermission,
		},
		model.Role{
			Name:        "read_dashboard",
			Description: "Read dashboard",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_dashboard",
			Description: "Write dashboard",
			Locations:   model.AllLocations,
			Permission:  model.WritePermission,
		},
		model.Role{
			Name:        "read_license_contract",
			Description: "read_license_contract",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_license_contract",
			Description: "write_license_contract",
			Locations:   model.AllLocations,
			Permission:  model.WritePermission,
		},
		model.Role{
			Name:        "read_license_compliance",
			Description: "read_license_compliance",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_license_compliance",
			Description: "write_license_compliance",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_license_used",
			Description: "read_license_used",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_alert_type_license",
			Description: "read_alert_type_license",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_host",
			Description: "read_host",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_host",
			Description: "write_host",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_databases",
			Description: "read_databases",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_databases",
			Description: "write_databases",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_hypervisor",
			Description: "read_hypervisor",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_hypervisor",
			Description: "write_hypervisor",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_exadata",
			Description: "read_exadata",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_exadata",
			Description: "write_exadata",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_alert",
			Description: "read_alert",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "write_alert",
			Description: "write_alert",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
		model.Role{
			Name:        "read_cloud",
			Description: "read_cloud",
			Locations:   model.AllLocations,
			Permission:  model.ReadPermission,
		},
	}
}
