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
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/utils"
)

// GetTotalOracleExadataMemorySizeStats return the total size of memory of exadata
func (md *MongoDatabase) GetTotalOracleExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float64, error) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.exadata.components"),
			mu.APReplaceWith("$features.oracle.exadata.components"),
			mu.APGroup(bson.M{
				"_id":   0,
				"value": mu.APOSum("$memory"),
			}),
		),
	)
	if err != nil {
		return 0, utils.NewError(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewError(err, "DB ERROR")
	}

	return float64(out["value"]), nil
}

// GetTotalOracleExadataCPUStats return the total cpu of exadata
func (md *MongoDatabase) GetTotalOracleExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, error) {
	var out map[string]interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.exadata.components"),
			mu.APGroup(bson.M{
				"_id":     0,
				"running": mu.APOSum("$features.oracle.exadata.components.runningCPUCount"),
				"total":   mu.APOSum("$features.oracle.exadata.components.totalCPUCount"),
			}),
			mu.APUnset("_id"),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return map[string]interface{}{
			"running": 0,
			"total":   0,
		}, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return out, nil
}

// GetAverageOracleExadataStorageUsageStats return the average usage of cell disks of exadata
func (md *MongoDatabase) GetAverageOracleExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float64, error) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"value": mu.APOReduce("$features.oracle.exadata.components", bson.M{"count": 0, "sum": 0},
					mu.APOLet(
						bson.M{
							"part": mu.APOIfNull(
								mu.APOReduce("$$this.cellDisks", bson.M{"count": 0, "sum": 0},
									bson.M{
										"count": mu.APOAdd("$$value.count", 1),
										"sum": mu.APOAdd(
											"$$value.sum",
											"$$this.usedPerc",
										),
									},
								),
								bson.M{"count": 0, "sum": 0},
							),
						},
						bson.M{
							"count": mu.APOAdd("$$value.count", "$$part.count"),
							"sum":   mu.APOAdd("$$value.sum", "$$part.sum"),
						},
					),
				),
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"count": mu.APOSum("$value.count"),
				"sum":   mu.APOSum("$value.sum"),
			}),
			mu.APProject(bson.M{
				"_id": 0,
				"value": mu.APOCond(
					bson.D{{Key: "$eq", Value: bson.A{0, "$count"}}},
					0,
					mu.APODivide("$sum", "$count"),
				),
			}),
		),
	)
	if err != nil {
		return 0, utils.NewError(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewError(err, "DB ERROR")
	}

	return float64(out["value"]), nil
}

// GetOracleExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
func (md *MongoDatabase) GetOracleExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, error) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"devs": mu.APOMap(
					mu.APOFilter("$features.oracle.exadata.components", "dev",
						mu.APOEqual("$$dev.serverType", "StorageServer"),
					),
					"dev",
					mu.APOMap("$$dev.cellDisks", "cd", mu.APOGreater("$$cd.errCount", 0)),
				),
			}),
			mu.APUnwind("$devs"),
			mu.APUnwind("$devs"),
			mu.APGroupAndCountStages("failing", "count", "$devs"),
			mu.APSort(bson.M{
				"failing": 1,
			}),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}

		out = append(out, &item)
	}

	return out, nil
}

// GetOracleExadataPatchStatusStats return a array containing the number of exadata per patch status
func (md *MongoDatabase) GetOracleExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, error) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"status": mu.APOMap("$features.oracle.exadata.components", "dev",
					mu.APOGreater(
						mu.APODateFromString("$$dev.swReleaseDate", "%Y-%m-%d"),
						windowTime,
					),
				),
			}),
			mu.APUnwind("$status"),
			mu.APGroupAndCountStages("status", "count", "$status"),
			mu.APSort(bson.M{
				"status": 1,
			}),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}

		out = append(out, &item)
	}

	return out, nil
}
