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

func (md *MongoDatabase) FindAllOracleDatabasePartitionings(filter dto.GlobalFilter) ([]dto.OracleDatabasePartitioning, error) {
	ctx := context.TODO()
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.partitionings"}},
			bson.M{"$project": bson.M{
				"hostname":     1,
				"databaseName": "$features.oracle.database.databases.uniqueName",
				"owner":        "$features.oracle.database.databases.partitionings.owner",
				"segmentName":  "$features.oracle.database.databases.partitionings.segmentName",
				"count":        "$features.oracle.database.databases.partitionings.count",
				"mb":           "$features.oracle.database.databases.partitionings.mb",
			}},
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	out := make([]dto.OracleDatabasePartitioning, 0)
	if err = cur.All(ctx, &out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}

func (md *MongoDatabase) FindAllOraclePDBPartitionings(filter dto.GlobalFilter) ([]dto.OracleDatabasePartitioning, error) {
	ctx := context.TODO()
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.pdbs"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.pdbs.partitionings"}},
			bson.M{"$project": bson.M{
				"hostname":     1,
				"databaseName": "$features.oracle.database.databases.uniqueName",
				"pdb":          "$features.oracle.database.databases.pdbs.name",
				"owner":        "$features.oracle.database.databases.pdbs.partitionings.owner",
				"segmentName":  "$features.oracle.database.databases.pdbs.partitionings.segmentName",
				"count":        "$features.oracle.database.databases.pdbs.partitionings.count",
				"mb":           "$features.oracle.database.databases.pdbs.partitionings.mb",
			}},
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	out := make([]dto.OracleDatabasePartitioning, 0)
	if err = cur.All(ctx, &out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}
