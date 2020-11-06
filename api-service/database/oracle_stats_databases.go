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
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetOracleDatabaseEnvironmentStats return a array containing the number of databases per environment
func (md *MongoDatabase) GetOracleDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APGroup(bson.M{
				"_id":   "$environment",
				"count": mu.APOSum(mu.APOCond("$features.oracle.database.databases", mu.APOSize("$features.oracle.database.databases"), 0)),
			}),
			mu.APProject(bson.M{
				"_id":         false,
				"environment": "$_id",
				"count":       true,
			}),
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

// GetOracleDatabaseHighReliabilityStats return a array containing the number of databases per high-reliability status
func (md *MongoDatabase) GetOracleDatabaseHighReliabilityStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			AddHardwareAbstraction("ha"),
			mu.APGroup(bson.M{
				"_id":   "$ha",
				"count": mu.APOSum(mu.APOCond("$features.oracle.database.databases", mu.APOSize("$features.oracle.database.databases"), 0)),
			}),
			mu.APProject(bson.M{
				"_id":   false,
				"ha":    "$_id",
				"count": true,
			}),
			mu.APSort(bson.M{
				"ha": 1,
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

// GetOracleDatabaseVersionStats return a array containing the number of databases per version
func (md *MongoDatabase) GetOracleDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroupAndCountStages("version", "count", "$database.version"),
			mu.APSort(bson.M{
				"version": 1,
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

// GetTopReclaimableOracleDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
func (md *MongoDatabase) GetTopReclaimableOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
				"hostname": true,
			}),
			mu.APProject(bson.M{
				"hostname":                   true,
				"dbname":                     "$database.name",
				"reclaimableSegmentAdvisors": mu.APOSumReducer("$database.segmentAdvisors", "$$this.reclaimable"),
			}),
			mu.APSort(bson.M{
				"reclaimableSegmentAdvisors": -1,
			}),
			mu.APLimit(limit),
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

// GetTopWorkloadOracleDatabaseStats return a array containing top databases by workload
func (md *MongoDatabase) GetTopWorkloadOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
				"hostname": true,
			}),
			mu.APProject(bson.M{
				"hostname": true,
				"dbname":   "$database.name",
				"workload": mu.APOConvertToDoubleOrZero("$database.work"),
			}),
			mu.APSort(bson.M{
				"workload": -1,
			}),
			mu.APLimit(limit),
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

// GetOracleDatabasePatchStatusStats return a array containing the number of databases per patch status
func (md *MongoDatabase) GetOracleDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			//TODO: we can map directly PSU to date instead of mapping indirectly using a object that contain the Date field
			mu.APProject(bson.M{
				"database.psus": mu.APOReduce(
					mu.APOMap("$database.psus", "psu", mu.APOMergeObjects(
						"$$psu",
						bson.M{
							"date": mu.APODateFromString("$$psu.date", "%Y-%m-%d"),
						},
					)),
					nil,
					mu.APOCond(
						mu.APOEqual("$$value", nil),
						"$$this",
						mu.APOMaxWithCmpExpr("$$value.date", "$$this.date", "$$value", "$$this"),
					),
				),
			}),
			mu.APGroupAndCountStages("status", "count",
				mu.APOCond(mu.APOGreater("$database.psus.date", windowTime), "OK", "KO"),
			),
			mu.APSort(bson.M{
				"status": 1,
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

// GetOracleDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
func (md *MongoDatabase) GetOracleDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroupAndCountStages("dataguard", "count", "$database.dataguard"),
			mu.APSort(bson.M{
				"dataguard": 1,
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

// GetOracleDatabaseRACStatusStats return a array containing the number of databases per RAC status
func (md *MongoDatabase) GetOracleDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroupAndCountStages("rac", "count", mu.APOAny("$database.licenses", "lic",
				mu.APOAnd(
					mu.APOEqual("$$lic.name", "Real Application Clusters"),
					mu.APOGreater("$$lic.count", 0),
				),
			)),
			mu.APSort(bson.M{
				"rac": 1,
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

// GetOracleDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
func (md *MongoDatabase) GetOracleDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroupAndCountStages("archivelog", "count", "$database.archivelog"),
			mu.APSort(bson.M{
				"archivelog": 1,
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

// GetTotalOracleDatabaseWorkStats return the total work of databases
func (md *MongoDatabase) GetTotalOracleDatabaseWorkStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"value": mu.APOSum(mu.APOConvertToDoubleOrZero("$database.work")),
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

	return float64(out["value"]), nil
}

// GetTotalOracleDatabaseMemorySizeStats return the total of memory size of databases
func (md *MongoDatabase) GetTotalOracleDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroup(bson.M{
				"_id": 0,
				"value": mu.APOSum(mu.APOAdd(
					"$database.pgaTarget",
					"$database.sgaTarget",
					"$database.memoryTarget",
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

	return float64(out["value"]), nil
}

// GetTotalOracleDatabaseDatafileSizeStats return the total size of datafiles of databases
func (md *MongoDatabase) GetTotalOracleDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"value": mu.APOSum("$database.datafileSize"),
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

	return float64(out["value"]), nil
}

// GetTotalOracleDatabaseSegmentSizeStats return the total size of segments of databases
func (md *MongoDatabase) GetTotalOracleDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"database": "$features.oracle.database.databases",
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"value": mu.APOSum(mu.APOConvertToDoubleOrZero("$database.segmentsSize")),
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

	return float64(out["value"]), nil
}

// GetTopUnusedOracleDatabaseInstanceResourceStats return a array containing top unused instance resource by workload
func (md *MongoDatabase) GetTopUnusedOracleDatabaseInstanceResourceStats(location string, environment string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(bson.M{
				"hostname": 1,
				"works": mu.APOReduce(
					mu.APOFilter("$features.oracle.database.databases", "db", mu.APONotEqual("$$db.work", nil)),
					bson.M{"totalWork": 0, "totalCPUCount": 0},
					bson.M{
						"totalWork":     mu.APOAdd("$$value.totalWork", "$$this.work"),
						"totalCPUCount": mu.APOAdd("$$value.totalCPUCount", "$$this.cpuCount"),
					},
				),
			}),
			mu.APProject(bson.M{
				"hostname": 1,
				"unused":   mu.APOSubtract("$works.totalCPUCount", "$works.totalWork"),
			}),
			mu.APSort(bson.M{
				"unused": -1,
			}),
			mu.APLimit(limit),
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
