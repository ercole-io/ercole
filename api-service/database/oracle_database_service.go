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

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) GetOracleServiceList(filter dto.GlobalFilter) ([]dto.OracleDatabaseServiceDto, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabaseServiceDto, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			ExcludeDR(),
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.services"}},
			bson.M{"$project": bson.M{
				"hostname":                              1,
				"location":                              1,
				"environment":                           1,
				"createdAt":                             1,
				"databasename":                          "$features.oracle.database.databases.name",
				"oracleDatabaseService.name":            "$features.oracle.database.databases.services.name",
				"oracleDatabaseService.failoverMethod":  "$features.oracle.database.databases.services.failoverMethod",
				"oracleDatabaseService.failoverType":    "$features.oracle.database.databases.services.failoverType",
				"oracleDatabaseService.failoverRetries": "$features.oracle.database.databases.services.failoverRetries",
				"oracleDatabaseService.failoverDelay":   "$features.oracle.database.databases.services.failoverDelay",
				"oracleDatabaseService.enabled":         "$features.oracle.database.databases.services.enabled",
				"oracleDatabaseService.containerName":   "$features.oracle.database.databases.services.containerName",
			}},
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}
