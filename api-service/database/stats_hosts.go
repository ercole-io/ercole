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
	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
)

// GetEnvironmentStats return a array containing the number of hosts per environment
func (md *MongoDatabase) GetEnvironmentStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APGroupAndCountStages("environment", "count", "$environment"),
		),
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
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APGroupAndCountStages("type", "count", "$info.type"),
		),
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

// GetOperatingSystemStats return a array containing the number of hosts per operanting system
func (md *MongoDatabase) GetOperatingSystemStats(location string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Create the aggregation branches
	aggregationBranches := []bson.M{}
	for _, v := range md.OperatingSystemAggregationRules {
		aggregationBranches = append(aggregationBranches, bson.M{
			"case": bson.M{
				"$regexMatch": bson.M{
					"input": "$info.os",
					"regex": v.Regex,
				},
			},
			"then": v.Group,
		})
	}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APGroupAndCountStages("operating_system", "count", bson.M{
				"$switch": bson.M{
					"branches": aggregationBranches,
					"default":  "$info.os",
				},
			}),
		),
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
