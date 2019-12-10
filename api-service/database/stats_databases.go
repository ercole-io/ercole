// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetDatabaseEnvironmentStats return a array containing the number of databases per environment
func (md *MongoDatabase) GetDatabaseEnvironmentStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		bson.A{
			optionalStep(location != "", bson.M{"$match": bson.M{
				"location": location,
			}}),
			bson.M{"$match": bson.M{
				"archived": false,
			}},
			bson.M{"$group": bson.M{
				"_id": "$environment",
				"count": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							"if": "$extra.databases",
							"then": bson.M{
								"$size": "$extra.databases",
							},
							"else": 0,
						},
					},
				},
			}},
			bson.M{"$project": bson.M{
				"_id":         false,
				"environment": "$_id",
				"count":       true,
			}},
		},
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetDatabaseVersionStats return a array containing the number of databases per version
func (md *MongoDatabase) GetDatabaseVersionStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		bson.A{
			optionalStep(location != "", bson.M{"$match": bson.M{
				"location": location,
			}}),
			bson.M{"$group": bson.M{
				"_id": "$database.version",
				"count": bson.M{
					"$sum": 1,
				},
			}},
			bson.M{"$project": bson.M{
				"_id":     false,
				"version": "$_id",
				"count":   true,
			}},
		},
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}
