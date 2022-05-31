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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func (md *MongoDatabase) SearchSqlServerDatabaseUsedLicenses(hostname string, sortBy string, sortDesc bool, page int, pageSize int,
	location string, environment string, olderThan time.Time,
) (*dto.SqlServerDatabaseUsedLicenseSearchResponse, error) {
	cursor, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FindByHostname(hostname),
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.microsoft.sqlServer.instances"),
			mu.APUnwind("$features.microsoft.sqlServer.instances.license"),
			mu.APMatch(bson.M{"features.microsoft.sqlServer.instances.license.count": bson.M{"$gt": 0}}),
			mu.APLookupSimple("ms_sqlserver_database_license_types", "features.microsoft.sqlServer.instances.license.licenseTypeID", "_id", "licenseType"),
			bson.M{"$unwind": bson.M{"path": "$licenseType", "preserveNullAndEmptyArrays": true}},
			mu.APProject(
				bson.M{
					"_id":            0,
					"hostname":       1,
					"dbName":         "$features.microsoft.sqlServer.instances.name",
					"licenseTypeID":  "$features.microsoft.sqlServer.instances.license.licenseTypeID",
					"usedLicenses":   "$features.microsoft.sqlServer.instances.license.count",
					"ignored":        "$features.microsoft.sqlServer.instances.license.ignored",
					"ignoredComment": "$features.microsoft.sqlServer.instances.license.ignoredComment",
				},
			),

			mu.APOptionalSortingStage(sortBy, sortDesc),
			PagingMetadataStage(page, pageSize),
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var response dto.SqlServerDatabaseUsedLicenseSearchResponse

	cursor.Next(context.TODO())

	if err := cursor.Decode(&response); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return &response, nil
}

func (md *MongoDatabase) UpdateSqlServerLicenseIgnoredField(hostname string, instancename string, ignored bool, ignoredComment string) error {
	result, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		UpdateOne(context.TODO(),
			bson.M{
				"hostname": hostname,
				"archived": false,
				"features.microsoft.sqlServer.instances.name": instancename,
			},
			bson.M{"$set": bson.M{
				"features.microsoft.sqlServer.instances.$[elemDB].license.ignored":        ignored,
				"features.microsoft.sqlServer.instances.$[elemDB].license.ignoredComment": ignoredComment,
			}},
			options.Update().SetArrayFilters(options.ArrayFilters{Filters: []interface{}{bson.M{"elemDB.name": instancename}}}),
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if result.MatchedCount != 1 {
		return utils.ErrLicenseNotFound
	}

	return nil
}
