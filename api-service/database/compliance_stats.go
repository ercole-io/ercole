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

const (
	oraclePipelinePathMatch     = "$features.oracle.database.databases"
	mysqlPipelinePathMatch      = "$features.mysql.instances"
	sqlServerPipelinePathMatch  = "$features.microsoft.sqlServer.instances"
	postgresqlPipelinePathMatch = "$features.postgresql.instances"
	mongoPipelinePathMatch      = "$features.mongodb.instances"
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
	pipeline := md.getCountInstancePipeline(oraclePipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountOracleInstanceByLocations(locations []string) (int64, error) {
	pipeline := md.getCountInstancePipeline(oraclePipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountOracleHosts() (int64, error) {
	pipeline := md.getCountHostPipeline(oraclePipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountOracleHostsByLocations(locations []string) (int64, error) {
	pipeline := md.getCountHostPipeline(oraclePipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMySqlInstance() (int64, error) {
	pipeline := md.getCountInstancePipeline(mysqlPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMySqlInstanceByLocations(locations []string) (int64, error) {
	pipeline := md.getCountInstancePipeline(mysqlPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMySqlHosts() (int64, error) {
	pipeline := md.getCountHostPipeline(mysqlPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMySqlHostsByLocations(locations []string) (int64, error) {
	pipeline := md.getCountHostPipeline(mysqlPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountSqlServerlInstance() (int64, error) {
	pipeline := md.getCountInstancePipeline(sqlServerPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountSqlServerlInstanceByLocations(locations []string) (int64, error) {
	pipeline := md.getCountInstancePipeline(sqlServerPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountSqlServerHosts() (int64, error) {
	pipeline := md.getCountHostPipeline(sqlServerPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountSqlServerHostsByLocations(locations []string) (int64, error) {
	pipeline := md.getCountHostPipeline(sqlServerPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountPostgreSqlInstance() (int64, error) {
	pipeline := md.getCountInstancePipeline(postgresqlPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountPostgreSqlInstanceByLocations(locations []string) (int64, error) {
	pipeline := md.getCountInstancePipeline(postgresqlPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountPostgreSqlHosts() (int64, error) {
	pipeline := md.getCountHostPipeline(postgresqlPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountPostgreSqlHostsByLocations(locations []string) (int64, error) {
	pipeline := md.getCountHostPipeline(postgresqlPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMongoDbInstance() (int64, error) {
	pipeline := md.getCountInstancePipeline(mongoPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMongoDbInstanceByLocations(locations []string) (int64, error) {
	pipeline := md.getCountInstancePipeline(mongoPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMongoDbHosts() (int64, error) {
	pipeline := md.getCountHostPipeline(mongoPipelinePathMatch)
	return md.count(pipeline)
}

func (md *MongoDatabase) CountMongoDbHostsByLocations(locations []string) (int64, error) {
	pipeline := md.getCountHostPipeline(mongoPipelinePathMatch, locations...)
	return md.count(pipeline)
}

func (md *MongoDatabase) getCountInstancePipeline(path string, locations ...string) bson.A {
	match := bson.D{{Key: "archived", Value: false}}
	if len(locations) > 0 {
		match = append(match, bson.E{Key: "location", Value: bson.M{"$in": locations}})
	}

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: path}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: primitive.Null{}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	return pipeline
}

func (md *MongoDatabase) getCountHostPipeline(path string, locations ...string) bson.A {
	match := bson.D{{Key: "archived", Value: false}}
	if len(locations) > 0 {
		match = append(match, bson.E{Key: "location", Value: bson.M{"$in": locations}})
	}

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: path}}}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$hostname"}}}},
		bson.D{{Key: "$count", Value: "count"}},
	}

	return pipeline
}

func (md *MongoDatabase) count(pipeline bson.A) (int64, error) {
	ctx := context.TODO()
	
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).Aggregate(ctx, pipeline)
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
