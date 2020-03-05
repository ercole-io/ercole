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
func (md *MongoDatabase) GetTotalExadataMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id": 0,
				"Value": mu.APOSum(mu.APOSumReducer("$Extra.Exadata.Devices",
					mu.APOConvertToDoubleOrZero(mu.APOGetCaptureFromRegexMatch("$$this.Memory", "^(\\d+)GB$", "i", 0)),
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
func (md *MongoDatabase) GetTotalExadataCPUStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"Value": mu.APOReduce("$Extra.Exadata.Devices", bson.M{"Enabled": 0, "Total": 0},
					mu.APOLet(
						bson.M{
							"match": mu.APORegexFind("$$this.CPUEnabled", "^(\\d+)/(\\d+)$", "i"),
						},
						bson.M{
							"Enabled": mu.APOAdd(
								"$$value.Enabled",
								mu.APOConvertToDoubleOrZero(mu.APOArrayElemAt("$$match.captures", 0)),
							),
							"Total": mu.APOAdd(
								"$$value.Total",
								mu.APOConvertToDoubleOrZero(mu.APOArrayElemAt("$$match.captures", 1)),
							),
						},
					),
				),
			}),
			mu.APGroup(bson.M{
				"_id":     0,
				"Enabled": mu.APOSum("$Value.Enabled"),
				"Total":   mu.APOSum("$Value.Total"),
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
				"Value": mu.APOReduce("$Extra.Exadata.Devices", bson.M{"Count": 0, "Sum": 0},
					mu.APOLet(
						bson.M{
							"part": mu.APOIfNull(
								mu.APOReduce("$$this.CellDisks", bson.M{"Count": 0, "Sum": 0},
									bson.M{
										"Count": mu.APOAdd("$$value.Count", 1),
										"Sum": mu.APOAdd(
											"$$value.Sum",
											mu.APOToDouble("$$this.UsedPerc"),
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

	return float32(out["value"]), nil
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
					mu.APOFilter("$Extra.Exadata.Devices", "Dev",
						mu.APOEqual("$$Dev.ServerType", "StorageServer"),
					),
					"dev",
					mu.APOMap("$$Dev.CellDisks", "cd", mu.APOGreater(mu.APOToDouble("$$cd.ErrCount"), 0)),
				),
			}),
			mu.APUnwind("$Devs"),
			mu.APUnwind("$Devs"),
			mu.APGroupAndCountStages("failing", "count", "$Devs"),
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
				"Status": mu.APOMap("$Extra.Exadata.Devices", "dev",
					mu.APOGreater(
						mu.APODateFromString(
							mu.APOConcat("20", mu.APOGetCaptureFromRegexMatch("$$dev.ExaSwVersion", "^.*\\.(\\d+)$", "i", 0)),
							"%Y%m%d",
						),
						windowTime,
					),
				),
			}),
			mu.APUnwind("$Status"),
			mu.APGroupAndCountStages("Status", "Count", "$Status"),
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
