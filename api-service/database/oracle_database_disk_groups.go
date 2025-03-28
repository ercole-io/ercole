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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) GetOracleDiskGroups(hostname string, dbname string) ([]dto.OracleDatabaseDiskGroupDto, error) {
	return md.load(dto.GlobalFilter{}, hostname, dbname)
}

func (md *MongoDatabase) ListOracleDiskGroups(filter dto.GlobalFilter) ([]dto.OracleDatabaseDiskGroupDto, error) {
	return md.load(filter, "", "")
}

func (md *MongoDatabase) load(filter dto.GlobalFilter, hostname string, dbname string) ([]dto.OracleDatabaseDiskGroupDto, error) {
	ctx := context.TODO()

	percentage := func(amount interface{}, total string) bson.M {
		return bson.M{
			"$cond": bson.A{
				bson.M{"$eq": bson.A{total, 0}},
				0,
				bson.M{"$round": bson.A{
					bson.M{"$multiply": bson.A{
						bson.M{"$divide": bson.A{amount, total}},
						100,
					}},
					2,
				}},
			},
		}
	}

	var pipeline bson.A

	if hostname != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"hostname": hostname}}})
	}

	if filter.OlderThan != utils.MAX_TIME && !filter.OlderThan.IsZero() {
		pipeline = append(pipeline, FilterByOldnessSteps(filter.OlderThan)...)
	}

	if filter.Location != "" || filter.Environment != "" {
		loc := FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment).(bson.A)
		pipeline = append(pipeline, loc...)
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$match", Value: bson.M{"archived": false, "isDR": false}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$features.oracle.database.databases"}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$features.oracle.database.databases.diskGroups"}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"hostname":      "$hostname",
				"diskGroupName": "$features.oracle.database.databases.diskGroups.diskGroupName",
			},
			"hostname":      bson.M{"$first": "$hostname"},
			"databases":     bson.M{"$addToSet": "$features.oracle.database.databases.name"},
			"diskGroupName": bson.M{"$first": "$features.oracle.database.databases.diskGroups.diskGroupName"},
			"totalSpace":    bson.M{"$first": "$features.oracle.database.databases.diskGroups.totalSpace"},
			"freeSpace":     bson.M{"$first": "$features.oracle.database.databases.diskGroups.freeSpace"},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":                                   0,
			"hostname":                              1,
			"databases":                             1,
			"oracleDatabaseDiskGroup.diskGroupName": "$diskGroupName",
			"oracleDatabaseDiskGroup.totalSpace":    "$totalSpace",
			"oracleDatabaseDiskGroup.freeSpace":     "$freeSpace",
			"percentageFreeSpace":                   percentage("$freeSpace", "$totalSpace"),
			"usedSpace":                             bson.M{"$subtract": bson.A{"$totalSpace", "$freeSpace"}},
			"percentageUsedSpace":                   percentage(bson.M{"$subtract": bson.A{"$totalSpace", "$freeSpace"}}, "$totalSpace"),
		}}})

	if dbname != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"databases": dbname}}})
	}

	result := make([]dto.OracleDatabaseDiskGroupDto, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}
