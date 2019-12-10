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

// GetEnvironmentStats return a array containing the number of hosts per environment
func (md *MongoDatabase) GetEnvironmentStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Find the matching hostdata
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
					"$sum": 1,
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

// GetTypeStats return a array containing the number of hosts per type
func (md *MongoDatabase) GetTypeStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Find the matching hostdata
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
				"_id": "$info.type",
				"count": bson.M{
					"$sum": 1,
				},
			}},
			bson.M{"$project": bson.M{
				"_id":   false,
				"type":  "$_id",
				"count": true,
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
