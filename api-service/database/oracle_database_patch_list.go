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

func (md *MongoDatabase) GetOraclePatchList(filter dto.GlobalFilter) ([]dto.OracleDatabasePatchDto, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabasePatchDto, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			ExcludeDR(),
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.patches"}},
			bson.M{"$project": bson.M{
				"hostname":                        1,
				"location":                        1,
				"environment":                     1,
				"createdAt":                       1,
				"databasename":                    "$features.oracle.database.databases.name",
				"oracleDatabasePatch.version":     "$features.oracle.database.databases.patches.version",
				"oracleDatabasePatch.patchID":     "$features.oracle.database.databases.patches.patchID",
				"oracleDatabasePatch.action":      "$features.oracle.database.databases.patches.action",
				"oracleDatabasePatch.description": "$features.oracle.database.databases.patches.description",
				"oracleDatabasePatch.date":        "$features.oracle.database.databases.patches.date",
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
