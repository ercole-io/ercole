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

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ercole-io/ercole/v2/utils"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := migrateSqlServerDatabaseLicenseTypes(db); err != nil {
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

func migrateSqlServerDatabaseLicenseTypes(client *mongo.Database) error {
	collection := "ms_sqlserver_database_license_types"

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
