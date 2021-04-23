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

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

//TODO Add tests
func (md *MongoDatabase) GetLicenseComplianceHistory() ([]dto.LicenseComplianceHistory, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection("oracle_database_licenses_history").
		Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var items []dto.LicenseComplianceHistory
	for cur.Next(context.TODO()) {
		var item dto.LicenseComplianceHistory
		if err := cur.Decode(&item); err != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}
		items = append(items, item)
	}

	return items, nil
}
