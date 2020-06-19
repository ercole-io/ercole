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
				"_id":   "$Environment",
				"Count": mu.APOSum(mu.APOCond("$Features.Oracle.Database.Databases", mu.APOSize("$Features.Oracle.Database.Databases"), 0)),
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

// GetOracleDatabaseHighReliabilityStats return a array containing the number of databases per high-reliability status
func (md *MongoDatabase) GetOracleDatabaseHighReliabilityStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)
	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			AddHardwareAbstraction("HA"),
			mu.APGroup(bson.M{
				"_id":   "$HA",
				"Count": mu.APOSum(mu.APOCond("$Features.Oracle.Database.Databases", mu.APOSize("$Features.Oracle.Database.Databases"), 0)),
			}),
			mu.APProject(bson.M{
				"_id":   false,
				"HA":    "$_id",
				"Count": true,
			}),
			mu.APSort(bson.M{
				"HA": 1,
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
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
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

// GetTopReclaimableOracleDatabaseStats return a array containing the total sum of reclaimable of segments advisors of the top reclaimable databases
func (md *MongoDatabase) GetTopReclaimableOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
				"Hostname": true,
			}),
			mu.APProject(bson.M{
				"Hostname":                   true,
				"Dbname":                     "$Database.Name",
				"ReclaimableSegmentAdvisors": mu.APOSumReducer("$Database.SegmentAdvisors", "$$this.Reclaimable"),
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

// GetTopWorkloadOracleDatabaseStats return a array containing top databases by workload
func (md *MongoDatabase) GetTopWorkloadOracleDatabaseStats(location string, limit int, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
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

// GetOracleDatabasePatchStatusStats return a array containing the number of databases per patch status
func (md *MongoDatabase) GetOracleDatabasePatchStatusStats(location string, windowTime time.Time, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, ""),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
			}),
			//TODO: we can map directly PSU to date instead of mapping indirectly using a object that contain the Date field
			mu.APProject(bson.M{
				"Database.PSUs": mu.APOReduce(
					mu.APOMap("$Database.PSUs", "psu", mu.APOMergeObjects(
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
				mu.APOCond(mu.APOGreater("$Database.PSUs.Date", windowTime), "OK", "KO"),
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

// GetOracleDatabaseDataguardStatusStats return a array containing the number of databases per dataguard status
func (md *MongoDatabase) GetOracleDatabaseDataguardStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
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

// GetOracleDatabaseRACStatusStats return a array containing the number of databases per RAC status
func (md *MongoDatabase) GetOracleDatabaseRACStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
			}),
			mu.APGroupAndCountStages("RAC", "Count", mu.APOAny("$Database.Licenses", "lic",
				mu.APOAnd(
					mu.APOEqual("$$lic.Name", "Real Application Clusters"),
					mu.APOGreater("$$lic.Count", 0),
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

// GetOracleDatabaseArchivelogStatusStats return a array containing the number of databases per archivelog status
func (md *MongoDatabase) GetOracleDatabaseArchivelogStatusStats(location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
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

// GetTotalOracleDatabaseWorkStats return the total work of databases
func (md *MongoDatabase) GetTotalOracleDatabaseWorkStats(location string, environment string, olderThan time.Time) (float64, utils.AdvancedErrorInterface) {
	var out map[string]float64

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
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

	return float64(out["Value"]), nil
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
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
			}),
			mu.APGroup(bson.M{
				"_id": 0,
				"Value": mu.APOSum(mu.APOAdd(
					"$Database.PGATarget",
					"$Database.SGATarget",
					"$Database.MemoryTarget",
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

	return float64(out["Value"]), nil
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
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
			}),
			mu.APGroup(bson.M{
				"_id":   0,
				"Value": mu.APOSum(mu.APOConvertToDoubleOrZero("$Database.DatafileSize")),
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

	return float64(out["Value"]), nil
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
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Database": "$Features.Oracle.Database.Databases",
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

	return float64(out["Value"]), nil
}

// GetOracleDatabaseLicenseComplianceStatusStats return the status of the compliance of licenses of databases
func (md *MongoDatabase) GetOracleDatabaseLicenseComplianceStatusStats(location string, environment string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
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
							mu.APOMap("$Features.Oracle.Database.Databases", "db", bson.M{
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
					mu.APUnwind("$Clusters"),
					mu.APReplaceWith("$Clusters"),
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
				"Compliance": mu.APOGreaterOrEqual(mu.APOCond("$Unlimited", "$Used", "$Count"), "$Used"),
			}),
			mu.APGroup(bson.M{
				"_id":                     0,
				"LicensesNumber":          mu.APOSum(1),
				"Count":                   mu.APOSum(mu.APOCond("$Unlimited", "$Used", "$Count")),
				"Used":                    mu.APOSum("$Used"),
				"CompliantLicensesNumber": mu.APOSum(mu.APOCond("$Compliance", 1, 0)),
			}),
			mu.APProject(bson.M{
				"_id":       0,
				"Unlimited": 1,
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
				"Hostname": 1,
				"Works": mu.APOReduce(
					mu.APOFilter("$Features.Oracle.Database.Databases", "db", mu.APONotEqual("$$db.Work", nil)),
					bson.M{"TotalWork": 0, "TotalCPUCount": 0},
					bson.M{
						"TotalWork":     mu.APOAdd("$$value.TotalWork", "$$this.Work"),
						"TotalCPUCount": mu.APOAdd("$$value.TotalCPUCount", "$$this.CPUCount"),
					},
				),
			}),
			mu.APProject(bson.M{
				"Hostname": 1,
				"Unused":   mu.APOSubtract("$Works.TotalCPUCount", "$Works.TotalWork"),
			}),
			mu.APSort(bson.M{
				"Unused": -1,
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
