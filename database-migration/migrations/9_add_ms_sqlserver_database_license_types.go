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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	licenseTypes := make([]interface{}, 0, 2)

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "DG7GMGF0FLR2-0002",
		ItemDescription: "SQL Server 2019 Standard Core - 2 Core License Pack",
		Edition:         "STD",
		Version:         "2019",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "DG7GMGF0FKZV-0001",
		ItemDescription: "SQL Server 2019 Enterprise Core - 2 Core License Pack",
		Edition:         "ENT",
		Version:         "2019",
	})

	err := migrate.Register(func(db *mongo.Database) error {
		if err := addSqlServerDatabaseLicenseTypes(db, licenseTypes); err != nil {
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

func addSqlServerDatabaseLicenseTypes(client *mongo.Database, licenseTypes []interface{}) error {
	collectionName := "ms_sqlserver_database_license_types"

	if collectionNames, err := client.ListCollectionNames(context.TODO(), bson.D{{Key: "name", Value: collectionName}}); len(collectionNames) == 0 {
		if err != nil {
			return err
		}

		return fmt.Errorf("collection [%s] not found", collectionName)
	}

	if _, err := client.Collection(collectionName).InsertMany(context.TODO(), licenseTypes); err != nil {
		return err
	}

	return nil
}
