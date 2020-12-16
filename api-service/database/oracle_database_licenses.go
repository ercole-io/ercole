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

package database

import (
	"context"
	"time"

	"github.com/amreo/mu"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// SearchOracleDatabaseUsedLicenses search used licenses
func (md *MongoDatabase) SearchOracleDatabaseUsedLicenses(sortBy string, sortDesc bool, page int, pageSize int,
	location string, environment string, olderThan time.Time,
) (*dto.OracleDatabaseUsedLicenseSearchResponse, utils.AdvancedErrorInterface) {
	cursor, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APUnwind("$features.oracle.database.databases.licenses"),
			mu.APMatch(bson.M{"features.oracle.database.databases.licenses.count": bson.M{"$gt": 0}}),
			mu.APProject(
				bson.M{
					"_id":          0,
					"hostname":     1,
					"dbName":       "$features.oracle.database.databases.name",
					"licenseName":  "$features.oracle.database.databases.licenses.name",
					"usedLicenses": "$features.oracle.database.databases.licenses.count",
				},
			),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			PagingMetadataStage(page, pageSize),
		),
	)

	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	var response dto.OracleDatabaseUsedLicenseSearchResponse

	cursor.Next(context.TODO())
	if err := cursor.Decode(&response); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
	}
	return &response, nil
}
