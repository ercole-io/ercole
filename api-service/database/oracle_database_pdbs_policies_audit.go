// Copyright (c) 2024 Sorint.lab S.p.A.
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

	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) FindOracleDatabasePdbPoliciesAudit(hostname, dbname, pdbname string) ([]string, error) {
	ctx := context.TODO()

	res := make([]string, 0)

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
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs.policiesAudit"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "policiesAudit", Value: "$features.oracle.database.databases.pdbs.policiesAudit"},
				},
			},
		},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$policiesAudit"}}}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var item map[string]string
		if cur.Decode(&item) != nil {
			return nil, err
		}

		res = append(res, item["_id"])
	}

	return res, nil
}
