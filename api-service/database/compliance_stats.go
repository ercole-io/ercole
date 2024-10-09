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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (md *MongoDatabase) CountAllHost() (int64, error) {
	filter := bson.D{
		{Key: "archived", Value: false},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return cur, nil
}

func (md *MongoDatabase) CountOracleInstance() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountOracleHosts() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$hostname"}}}},
		bson.D{{Key: "$count", Value: "count"}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountMySqlInstance() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.mysql.instances"}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountMySqlHosts() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.mysql.instances"}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$hostname"}}}},
		bson.D{{Key: "$count", Value: "count"}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountSqlServerlInstance() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.microsoft.sqlServer.instances"}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountSqlServerHosts() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.microsoft.sqlServer.instances"}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$hostname"}}}},
		bson.D{{Key: "$count", Value: "count"}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountPostgreSqlInstance() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.postgresql.instances"}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountPostgreSqlHosts() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.postgresql.instances"}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$hostname"}}}},
		bson.D{{Key: "$count", Value: "count"}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountMongoDbInstance() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.mongodb.instances"}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}

func (md *MongoDatabase) CountMongoDbHosts() (int64, error) {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.mongodb.instances"}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$hostname"}}}},
		bson.D{{Key: "$count", Value: "count"}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var out int64

	for cur.Next(ctx) {
		var item map[string]int64
		if cur.Decode(&item) != nil {
			return 0, err
		}

		out = item["count"]
	}

	return out, nil
}
