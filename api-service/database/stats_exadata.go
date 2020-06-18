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
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTotalExadataMemorySizeStats return the total size of memory of exadata
func (md *MongoDatabase) GetTotalExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id":   0,
				"Value": mu.APOSum(mu.APOSumReducer("$Features.Oracle.Exadata.Components", "$$this.Memory")),
			}),
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

	return float32(out["Value"]), nil
}

// GetTotalExadataCPUStats return the total cpu of exadata
func (md *MongoDatabase) GetTotalExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Exadata.Components"),
			mu.APGroup(bson.M{
				"_id":     0,
				"Running": mu.APOSum("$Features.Oracle.Exadata.Components.RunningCPUCount"),
				"Total":   mu.APOSum("$Features.Oracle.Exadata.Components.TotalCPUCount"),
			}),
			mu.APUnset("_id"),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return map[string]interface{}{
			"Running": 0,
			"Total":   0,
		}, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}

// GetAverageExadataStorageUsageStats return the average usage of cell disks of exadata
func (md *MongoDatabase) GetAverageExadataStorageUsageStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"Value": mu.APOReduce("$Features.Oracle.Exadata.Components", bson.M{"Count": 0, "Sum": 0},
					mu.APOLet(
						bson.M{
							"part": mu.APOIfNull(
								mu.APOReduce("$$this.CellDisks", bson.M{"Count": 0, "Sum": 0},
									bson.M{
										"Count": mu.APOAdd("$$value.Count", 1),
										"Sum": mu.APOAdd(
											"$$value.Sum",
											"$$this.UsedPerc",
										),
									},
								),
								bson.M{"Count": 0, "Sum": 0},
							),
						},
						bson.M{
							"Count": mu.APOAdd("$$value.Count", "$$part.Count"),
							"Sum":   mu.APOAdd("$$value.Sum", "$$part.Sum"),
						},
					),
				),
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"Count": mu.APOSum("$Value.Count"),
				"Sum":   mu.APOSum("$Value.Sum"),
			}),
			mu.APProject(bson.M{
				"_id":   0,
				"Value": mu.APODivide("$Sum", "$Count"),
			}),
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

	return float32(out["Value"]), nil
}

// GetExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
func (md *MongoDatabase) GetExadataStorageErrorCountStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"Devs": mu.APOMap(
					mu.APOFilter("$Features.Oracle.Exadata.Components", "dev",
						mu.APOEqual("$$dev.ServerType", "StorageServer"),
					),
					"dev",
					mu.APOMap("$$dev.CellDisks", "cd", mu.APOGreater("$$cd.ErrCount", 0)),
				),
			}),
			mu.APUnwind("$Devs"),
			mu.APUnwind("$Devs"),
			mu.APGroupAndCountStages("Failing", "Count", "$Devs"),
			mu.APSort(bson.M{
				"Failing": 1,
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

// GetExadataPatchStatusStats return a array containing the number of exadata per patch status
func (md *MongoDatabase) GetExadataPatchStatusStats(location string, environment string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"Status": mu.APOMap("$Features.Oracle.Exadata.Components", "dev",
					mu.APOGreater(
						mu.APODateFromString("$$dev.ExaSwVersion", "%Y%m%d"),
						windowTime,
					),
				),
			}),
			mu.APUnwind("$Status"),
			mu.APGroupAndCountStages("Status", "Count", "$Status"),
			mu.APSort(bson.M{
				"Status": 1,
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
