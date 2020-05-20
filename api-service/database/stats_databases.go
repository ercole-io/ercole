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

// GetDatabaseEnvironmentStats return a array containing the number of databases per environment
func (md *MongoDatabase) GetDatabaseEnvironmentStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APGroup(bson.M{
				"_id":   "$Environment",
				"Count": mu.APOSum(mu.APOCond("$Extra.Databases", mu.APOSize("$Extra.Databases"), 0)),
			}),
			mu.APProject(bson.M{
				"_id":         false,
				"Environment": "$_id",
				"Count":       true,
			}),
			mu.APSort(bson.M{
				"Environment": 1,
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

// GetDatabaseVersionStats return a array containing the number of databases per version
func (md *MongoDatabase) GetDatabaseVersionStats(location string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroupAndCountStages("Version", "Count", "$Database.Version"),
			mu.APSort(bson.M{
				"Version": 1,
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

// GetTopReclaimableDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
func (md *MongoDatabase) GetTopReclaimableDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
				"Hostname": true,
			}),
			mu.APProject(bson.M{
				"Hostname": true,
				"Dbname":   "$Database.Name",
				"ReclaimableSegmentAdvisors": mu.APOSumReducer("$Database.SegmentAdvisors",
					mu.APOConvertErrorable("$$this.Reclaimable", "double", 0.5),
				),
			}),
			mu.APSort(bson.M{
				"ReclaimableSegmentAdvisors": -1,
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

// GetTopWorkloadDatabaseStats return a array containing top databases by workload
func (md *MongoDatabase) GetTopWorkloadDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
				"Hostname": true,
			}),
			mu.APProject(bson.M{
				"Hostname": true,
				"Dbname":   "$Database.Name",
				"Workload": mu.APOConvertToDoubleOrZero("$Database.Work"),
			}),
			mu.APSort(bson.M{
				"Workload": -1,
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

// GetDatabasePatchStatusStats return a array containing the number of databases per patch status
func (md *MongoDatabase) GetDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			//TODO: we can map directly PSU to date instead of mapping indirectly using a object that contain the Date field
			mu.APProject(bson.M{
				"Database.LastPSUs": mu.APOReduce(
					mu.APOMap("$Database.LastPSUs", "psu", mu.APOMergeObjects(
						"$$psu",
						bson.M{
							"Date": mu.APODateFromString("$$psu.Date", "%Y-%m-%d"),
						},
					)),
					nil,
					mu.APOCond(
						mu.APOEqual("$$value", nil),
						"$$this",
						mu.APOMaxWithCmpExpr("$$value.Date", "$$this.Date", "$$value", "$$this"),
					),
				),
			}),
			mu.APGroupAndCountStages("Status", "Count",
				mu.APOCond(mu.APOGreater("$Database.LastPSUs.Date", windowTime), "OK", "KO"),
			),
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

// GetDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
func (md *MongoDatabase) GetDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroupAndCountStages("Dataguard", "Count", "$Database.Dataguard"),
			mu.APSort(bson.M{
				"Dataguard": 1,
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

// GetDatabaseRACStatusStats return a array containing the number of databases per RAC status
func (md *MongoDatabase) GetDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroupAndCountStages("RAC", "Count", mu.APOAny("$Database.Features", "fe",
				mu.APOAnd(
					mu.APOEqual("$$fe.Name", "Real Application Clusters"),
					mu.APOEqual("$$fe.Status", true),
				),
			)),
			mu.APSort(bson.M{
				"RAC": 1,
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

// GetDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
func (md *MongoDatabase) GetDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroupAndCountStages("Archivelog", "Count",
				mu.APOEqual("$Database.Archivelog", "ARCHIVELOG"),
			),
			mu.APSort(bson.M{
				"Archivelog": 1,
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

// GetTotalDatabaseWorkStats return the total work of databases
func (md *MongoDatabase) GetTotalDatabaseWorkStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"Value": mu.APOSum(mu.APOConvertToDoubleOrZero("$Database.Work")),
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

// GetTotalDatabaseMemorySizeStats return the total of memory size of databases
func (md *MongoDatabase) GetTotalDatabaseMemorySizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroup(bson.M{
				"_id": 0,
				"Value": mu.APOSum(mu.APOAdd(
					mu.APOConvertToDoubleOrZero("$Database.PGATarget"),
					mu.APOConvertToDoubleOrZero("$Database.SGATarget"),
					mu.APOConvertToDoubleOrZero("$Database.MemoryTarget"),
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

	return float32(out["Value"]), nil
}

// GetTotalDatabaseDatafileSizeStats return the total size of datafiles of databases
func (md *MongoDatabase) GetTotalDatabaseDatafileSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"Value": mu.APOSum(mu.APOConvertToDoubleOrZero("$Database.Used")),
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

// GetTotalDatabaseSegmentSizeStats return the total size of segments of databases
func (md *MongoDatabase) GetTotalDatabaseSegmentSizeStats(location string, environment string, olderThan time.Time) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APProject(bson.M{
				"Database": "$Extra.Databases",
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"Value": mu.APOSum(mu.APOConvertToDoubleOrZero("$Database.SegmentsSize")),
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

// GetDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
func (md *MongoDatabase) GetDatabaseLicenseComplianceStatusStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{} = map[string]interface{}{
		"Count":     0,
		"Used":      0,
		"Compliant": true,
	}

	//Find the informations
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APLookupPipeline("hosts", bson.M{
				"ln": "$_id",
			}, "Used", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APProject(bson.M{
					"Hostname": 1,
					"Databases": mu.APOReduce(
						mu.APOFilter(
							mu.APOMap("$Extra.Databases", "db", bson.M{
								"Name": "$$db.Name",
								"Count": mu.APOLet(
									bson.M{
										"val": mu.APOArrayElemAt(mu.APOFilter("$$db.Licenses", "lic", mu.APOEqual("$$lic.Name", "$$ln")), 0),
									},
									"$$val.Count",
								),
							}),
							"db",
							mu.APOGreater("$$db.Count", 0),
						),
						bson.M{"Count": 0, "DBs": bson.A{}},
						bson.M{
							"Count": mu.APOMax("$$value.Count", "$$this.Count"),
							"DBs": bson.M{
								"$concatArrays": bson.A{
									"$$value.DBs",
									bson.A{"$$this.Name"},
								},
							},
						},
					),
				}),
				mu.APMatch(bson.M{
					"Databases.Count": bson.M{
						"$gt": 0,
					},
				}),
				mu.APLookupPipeline("hosts", bson.M{"hn": "$Hostname"}, "VM", mu.MAPipeline(
					FilterByOldnessSteps(olderThan),
					mu.APUnwind("$Extra.Clusters"),
					mu.APReplaceWith("$Extra.Clusters"),
					mu.APUnwind("$VMs"),
					mu.APMatch(mu.QOExpr(mu.APOEqual("$VMs.Hostname", "$$hn"))),
					mu.APLimit(1),
				)),
				mu.APSet(bson.M{
					"VM": mu.APOArrayElemAt("$VM", 0),
				}),
				mu.APAddFields(bson.M{
					"ClusterName":  mu.APOIfNull("$VM.ClusterName", nil),
					"PhysicalHost": mu.APOIfNull("$VM.PhysicalHost", nil),
				}),
				mu.APUnset("VM"),
				mu.APGroup(bson.M{
					"_id": mu.APOCond(
						"$ClusterName",
						mu.APOConcat("cluster_§$#$§_", "$ClusterName"),
						mu.APOConcat("hostname_§$#$§_", "$Hostname"),
					),
					"License":    mu.APOMaxAggr("$Databases.Count"),
					"ClusterCpu": mu.APOMaxAggr("$ClusterCpu"),
				}),
				mu.APSet(bson.M{
					"License": mu.APOCond(
						"$ClusterCpu",
						mu.APODivide("$ClusterCpu", 2),
						"$License",
					),
				}),
				mu.APGroup(bson.M{
					"_id":   0,
					"Value": mu.APOSum("$License"),
				}),
			)),
			mu.APSet(bson.M{
				"Used": mu.APOArrayElemAt("$Used", 0),
			}),
			mu.APSet(bson.M{
				"Used": mu.APOIfNull(mu.APOCeil("$Used.Value"), 0),
			}),
			mu.APSet(bson.M{
				"Compliance": mu.APOGreaterOrEqual("$Count", "$Used"),
			}),
			mu.APGroup(bson.M{
				"_id":                     0,
				"LicensesNumber":          mu.APOSum(1),
				"Count":                   mu.APOSum("$Count"),
				"Used":                    mu.APOSum("$Used"),
				"CompliantLicensesNumber": mu.APOSum(mu.APOCond("$Compliance", 1, 0)),
			}),
			mu.APProject(bson.M{
				"_id":       0,
				"Count":     1,
				"Used":      1,
				"Compliant": mu.APOEqual("$LicensesNumber", "$CompliantLicensesNumber"),
			}),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return out, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}
