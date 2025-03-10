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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) GetMissingDatabases() ([]dto.OracleDatabaseMissingDbs, error) {
	res := make([]dto.OracleDatabaseMissingDbs, 0)

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "archived", Value: false},
			{Key: "$expr", Value: bson.M{
				"$gt": bson.A{
					bson.M{
						"$size": bson.M{"$ifNull": bson.A{"$features.oracle.database.missingDatabases", bson.A{}}},
					},
					0,
				},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"hostname":         1,
			"missingDatabases": "$features.oracle.database.missingDatabases",
		}}},
	}

	ctx := context.TODO()

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (md *MongoDatabase) GetMissingDatabasesByHostname(hostname string) ([]model.MissingDatabase, error) {
	res := make([]model.MissingDatabase, 0)

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "archived", Value: false},
			{Key: "hostname", Value: hostname},
		}}},
		bson.D{{Key: "$unwind", Value: "$features.oracle.database.missingDatabases"}},
		bson.D{{Key: "$project", Value: bson.M{
			"name":           "$features.oracle.database.missingDatabases.name",
			"ignored":        "$features.oracle.database.missingDatabases.ignored",
			"ignoredComment": "$features.oracle.database.missingDatabases.ignoredComment",
		}}},
	}

	ctx := context.TODO()

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}
