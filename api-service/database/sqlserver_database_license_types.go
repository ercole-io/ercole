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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

const sqlServerDbLicenseTypesCollection = "ms_sqlserver_database_license_types"

func (md *MongoDatabase) GetSqlServerDatabaseLicenseTypes() ([]model.SqlServerDatabaseLicenseType, error) {
	ctx := context.TODO()
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection(sqlServerDbLicenseTypesCollection).
		Find(ctx, bson.M{})

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	licenseTypes := make([]model.SqlServerDatabaseLicenseType, 0)

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var licenseType model.SqlServerDatabaseLicenseType

		err := cur.Decode(&licenseType)
		if err != nil {
			log.Fatal(err)
		}

		licenseTypes = append(licenseTypes, licenseType)
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return licenseTypes, nil
}

// InsertSqlServerDatabaseLicenseType insert an SqlServer/Database license type into the database
func (md *MongoDatabase) InsertSqlServerDatabaseLicenseType(licenseType model.SqlServerDatabaseLicenseType) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(sqlServerDbLicenseTypesCollection).
		InsertOne(context.TODO(), licenseType)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}
