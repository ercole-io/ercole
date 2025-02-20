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

func (md *MongoDatabase) GetOracleDiskGroups(filter dto.GlobalFilter) ([]dto.OracleDatabaseDiskGroupDto, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabaseDiskGroupDto, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.diskGroups"}},
			bson.M{"$group": bson.M{
				"_id": bson.M{
					"hostname":      "$hostname",
					"diskGroupName": "$features.oracle.database.databases.diskGroups.diskGroupName",
				},
				"hostname":      bson.M{"$first": "$hostname"},
				"databases":     bson.M{"$addToSet": "$features.oracle.database.databases.name"},
				"diskGroupName": bson.M{"$first": "$features.oracle.database.databases.diskGroups.diskGroupName"},
				"totalSpace":    bson.M{"$first": "$features.oracle.database.databases.diskGroups.totalSpace"},
				"freeSpace":     bson.M{"$first": "$features.oracle.database.databases.diskGroups.freeSpace"},
				"usedSpace":     bson.M{"$first": "$features.oracle.database.databases.diskGroups.usedSpace"},
			}},
			bson.M{"$project": bson.M{
				"_id":                                   0,
				"hostname":                              1,
				"databases":                             1,
				"oracleDatabaseDiskGroup.diskGroupName": "$diskGroupName",
				"oracleDatabaseDiskGroup.totalSpace":    "$totalSpace",
				"oracleDatabaseDiskGroup.freeSpace":     "$freeSpace",
				"oracleDatabaseDiskGroup.usedSpace":     "$usedSpace",
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
