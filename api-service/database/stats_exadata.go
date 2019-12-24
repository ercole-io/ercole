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

	"github.com/amreo/ercole-services/utils"
	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTotalExadataMemorySizeStats return the total size of memory of exadata
func (md *MongoDatabase) GetTotalExadataMemorySizeStats(location string, environment string) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APGroup(bson.M{
				"_id": 0,
				"value": mu.APOSum(mu.APOSumReducer("$extra.exadata.devices",
					mu.APOConvertToDoubleOrZero(mu.APOGetCaptureFromRegexMatch("$$this.memory", "^(\\d+)GB$", "i", 0)),
				)),
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

	return float32(out["value"]), nil
}

// GetTotalExadataCPUStats return the total cpu of exadata
func (md *MongoDatabase) GetTotalExadataCPUStats(location string, environment string) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APProject(bson.M{
				"value": mu.APOReduce("$extra.exadata.devices", bson.M{"enabled": 0, "total": 0},
					mu.APOLet(
						bson.M{
							"match": mu.APORegexFind("$$this.cpu_enabled", "^(\\d+)/(\\d+)$", "i"),
						},
						bson.M{
							"enabled": mu.APOAdd(
								"$$value.enabled",
								mu.APOConvertToDoubleOrZero(mu.APOArrayElemAt("$$match.captures", 0)),
							),
							"total": mu.APOAdd(
								"$$value.total",
								mu.APOConvertToDoubleOrZero(mu.APOArrayElemAt("$$match.captures", 1)),
							),
						},
					),
				),
			}),
			mu.APGroup(bson.M{
				"_id":     0,
				"enabled": mu.APOSum("$value.enabled"),
				"total":   mu.APOSum("$value.total"),
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
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}

// GetAvegageExadataStorageUsageStats return the average usage of cell disks of exadata
func (md *MongoDatabase) GetAvegageExadataStorageUsageStats(location string, environment string) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APProject(bson.M{
				"value": mu.APOReduce("$extra.exadata.devices", bson.M{"count": 0, "sum": 0},
					mu.APOLet(
						bson.M{
							"part": mu.APOIfNull(
								mu.APOReduce("$$this.cell_disks", bson.M{"count": 0, "sum": 0},
									bson.M{
										"count": mu.APOAdd("$$value.count", 1),
										"sum": mu.APOAdd(
											"$$value.sum",
											mu.APOToDouble("$$this.used_perc"),
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
				"_id":   0,
				"value": mu.APODivide("$sum", "$count"),
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

	return float32(out["value"]), nil
}

// GetExadataStorageErrorCountStatusStats return a array containing the number of cell disks of exadata per error count status
func (md *MongoDatabase) GetExadataStorageErrorCountStatusStats(location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APProject(bson.M{
				"devs": mu.APOMap(
					mu.APOFilter("$extra.exadata.devices", "dev",
						mu.APOEqual("$$dev.server_type", "StorageServer"),
					),
					"dev",
					mu.APOMap("$$dev.cell_disks", "cd", mu.APOGreater(mu.APOToDouble("$$cd.err_count"), 0)),
				),
			}),
			mu.APUnwind("$devs"),
			mu.APUnwind("$devs"),
			mu.APGroupAndCountStages("failing", "count", "$devs"),
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
func (md *MongoDatabase) GetExadataPatchStatusStats(location string, environment string, windowTime time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APProject(bson.M{
				"status": mu.APOMap("$extra.exadata.devices", "dev",
					mu.APOGreater(
						mu.APODateFromString(
							mu.APOConcat("20", mu.APOGetCaptureFromRegexMatch("$$dev.exa_sw_version", "^.*\\.(\\d+)$", "i", 0)),
							"%Y%m%d",
						),
						windowTime,
					),
				),
			}),
			mu.APUnwind("$status"),
			mu.APGroupAndCountStages("status", "count", "$status"),
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
