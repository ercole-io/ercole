// Copyright (c) 2020 Sorint.lab S.p.A.
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

// Package database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (md *MongoDatabase) HistoricizeOracleDbsLicenses(licenses []dto.OracleDatabaseLicenseUsage) error {
	now := md.TimeNow()

	updateOptions := options.Update()
	updateOptions.SetUpsert(true)

	for _, license := range licenses {
		filter := bson.D{{"licenseTypeID", license.LicenseTypeID}}
		update := bson.D{
			{"$push", bson.D{{
				"history", bson.D{{"date", now}, {"consumed", license.Consumed}, {"covered", license.Covered}},
			}}}}

		_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("oracle_database_licenses_history").
			UpdateMany(context.TODO(),
				filter,
				update,
				updateOptions)
		if err != nil {
			return utils.NewAdvancedErrorPtr(err, "DB ERROR")
		}
	}

	return nil
}
