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

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) FindPsqlMigrabilities(hostname, dbname string) ([]model.PgsqlMigrability, error) {
	ctx := context.TODO()

	result := make([]model.PgsqlMigrability, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx,
			bson.A{
				bson.D{
					{Key: "$match",
						Value: bson.D{
							{Key: "archived", Value: false},
							{Key: "hostname", Value: hostname},
						},
					},
				},
				bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
				bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.name", Value: dbname}}}},
				bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pgsqlMigrability"}}}},
				bson.D{
					{Key: "$project",
						Value: bson.D{
							{Key: "metric", Value: "$features.oracle.database.databases.pgsqlMigrability.metric"},
							{Key: "count", Value: "$features.oracle.database.databases.pgsqlMigrability.count"},
							{Key: "schema", Value: "$features.oracle.database.databases.pgsqlMigrability.schema"},
							{Key: "objectType", Value: "$features.oracle.database.databases.pgsqlMigrability.objectType"},
						},
					},
				},
			})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) FindPdbPsqlMigrabilities(hostname, dbname, pdbname string) ([]model.PgsqlMigrability, error) {
	ctx := context.TODO()

	result := make([]model.PgsqlMigrability, 0)

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "archived", Value: false},
					{Key: "hostname", Value: hostname},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.name", Value: dbname}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.pdbs.name", Value: pdbname}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "metric", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.metric"},
					{Key: "count", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.count"},
					{Key: "schema", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.schema"},
					{Key: "objectType", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.objectType"},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
