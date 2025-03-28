// Copyright (c) 2023 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) FindOraclePDBChangesByHostname(filter dto.GlobalFilter, hostname string, start time.Time, end time.Time) ([]dto.OraclePdbChange, error) {
	ctx := context.TODO()

	result := make([]dto.OraclePdbChange, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(ctx, bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "isDR", Value: false},
					{Key: "hostname", Value: hostname},
					{Key: "createdAt",
						Value: bson.D{
							{Key: "$gte", Value: start},
							{Key: "$lt", Value: end},
						},
					},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "dbname", Value: "$features.oracle.database.databases.name"},
					{Key: "pdbname", Value: "$features.oracle.database.databases.pdbs.name"},
					{Key: "updated", Value: "$createdAt"},
					{Key: "datafileSize", Value: "$features.oracle.database.databases.pdbs.datafileSize"},
					{Key: "segmentsSize", Value: "$features.oracle.database.databases.pdbs.segmentsSize"},
					{Key: "allocable", Value: "$features.oracle.database.databases.pdbs.allocable"},
				},
			},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{{Key: "updated", Value: 1}}},
		},
	})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}
