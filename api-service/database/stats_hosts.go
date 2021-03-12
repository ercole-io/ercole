// Copyright (c) 2020 Sorint.lab S.p.A.
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

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetHostsCountStats return the number of the non-archived hosts
func (md *MongoDatabase) GetHostsCountStats(location string, environment string, olderThan time.Time) (int, error) {
	var out map[string]int

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APCount("value"),
		),
	)
	if err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out["value"], nil
}

// GetEnvironmentStats return a array containing the number of hosts per environment
func (md *MongoDatabase) GetEnvironmentStats(location string, olderThan time.Time) ([]interface{}, error) {
	var out []interface{} = make([]interface{}, 0)
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APGroupAndCountStages("environment", "count", "$environment"),
			mu.APSort(bson.M{
				"environment": 1,
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

// GetTypeStats return a array containing the number of hosts per type
func (md *MongoDatabase) GetTypeStats(location string, olderThan time.Time) ([]interface{}, error) {
	var out []interface{} = make([]interface{}, 0)
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APGroupAndCountStages("hardwareAbstractionTechnology", "count", "$info.hardwareAbstractionTechnology"),
			mu.APSort(bson.M{
				"hardwareAbstractionTechnology": 1,
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

// GetOperatingSystemStats return a array containing the number of hosts per operanting system
func (md *MongoDatabase) GetOperatingSystemStats(location string, olderThan time.Time) ([]interface{}, error) {
	var out []interface{} = make([]interface{}, 0)

	//Create the aggregation branches
	var switchExpr interface{}
	if len(md.OperatingSystemAggregationRules) > 0 {
		aggregationBranches := []bson.M{}
		for _, v := range md.OperatingSystemAggregationRules {
			aggregationBranches = append(aggregationBranches, bson.M{
				"case": bson.M{
					"$regexMatch": bson.M{
						"input": mu.APOConcat("$info.os", " ", "$info.osVersion"),
						"regex": v.Regex,
					},
				},
				"then": v.Group,
			})
		}

		switchExpr = bson.M{
			"$switch": bson.M{
				"branches": aggregationBranches,
				"default":  mu.APOConcat("$info.os", " ", "$info.osVersion"),
			},
		}
	} else {
		switchExpr = mu.APOConcat("$info.os", " ", "$info.osVersion")
	}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APGroupAndCountStages("operatingSystem", "count", switchExpr),
			mu.APSort(bson.M{
				"operatingSystem": 1,
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
