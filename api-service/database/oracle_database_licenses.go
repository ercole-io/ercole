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
	"time"

	"github.com/amreo/mu"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

// SearchOracleDatabaseUsedLicenses search used licenses
func (md *MongoDatabase) SearchOracleDatabaseUsedLicenses(sortBy string, sortDesc bool, page int, pageSize int,
	location string, environment string, olderThan time.Time,
) (*dto.OracleDatabaseUsedLicenseSearchResponse, error) {
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
					"_id":           0,
					"hostname":      1,
					"dbName":        "$features.oracle.database.databases.name",
					"licenseTypeID": "$features.oracle.database.databases.licenses.licenseTypeID",
					"usedLicenses":  "$features.oracle.database.databases.licenses.count",
					"ignored":       "$features.oracle.database.databases.licenses.ignored",
				},
			),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			PagingMetadataStage(page, pageSize),
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var response dto.OracleDatabaseUsedLicenseSearchResponse

	cursor.Next(context.TODO())
	if err := cursor.Decode(&response); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}
	return &response, nil
}

// UpdateLicenseIgnoredField update host ignored field (true/false)
func (md *MongoDatabase) UpdateLicenseIgnoredField(hostname string, dbname string, licenseTypeID string, ignored bool) error {
	result, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		UpdateOne(context.TODO(),
			bson.M{
				"hostname": hostname,
				"archived": false,
				"features.oracle.database.databases.instanceName":           dbname,
				"features.oracle.database.databases.licenses.licenseTypeID": licenseTypeID,
			},
			bson.M{"$set": bson.M{"features.oracle.database.databases.$[elemDB].licenses.$[elemLic].ignored": ignored}},
			options.Update().SetArrayFilters(options.ArrayFilters{Filters: []interface{}{bson.M{"elemDB.instanceName": dbname}, bson.M{"elemLic.licenseTypeID": licenseTypeID}}}),
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}
	if result.MatchedCount != 1 {
		return utils.ErrLicenseNotFound
	}

	return nil
}
