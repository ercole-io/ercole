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
	"go.mongodb.org/mongo-driver/mongo"
)

// SearchHosts search hosts
func (md *MongoDatabase) SearchHosts(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			FilterByOldnessSteps(olderThan),
			mu.APSearchFilterStage([]string{
				"Hostname",
				"Extra.Databases.Name",
				"Extra.Databases.UniqueName",
				"Extra.Clusters.Name",
			}, keywords),
			mu.APLookupPipeline("hosts", bson.M{"hn": "$Hostname"}, "VM", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APUnwind("$Extra.Clusters"),
				mu.APReplaceWith("$Extra.Clusters"),
				mu.APUnwind("$VMs"),
				mu.APReplaceWith("$VMs"),
				mu.APMatch(bson.M{
					"$expr": mu.APOEqual("$Hostname", "$$hn"),
				}),
				mu.APLimit(1),
			)),
			mu.APSet(bson.M{
				"VM": mu.APOArrayElemAt("$VM", 0),
			}),
			mu.APAddFields(bson.M{
				"Cluster":      mu.APOIfNull("$vm.ClusterName", nil),
				"PhysicalHost": mu.APOIfNull("$vm.PhysicalHost", nil),
			}),
			mu.APUnset("VM"),
			mu.APOptionalStage(!full, mu.APProject(bson.M{
				"Hostname":       true,
				"Location":       true,
				"Environment":    true,
				"HostType":       true,
				"Cluster":        true,
				"PhysicalHost":   true,
				"CreatedAt":      true,
				"Databases":      true,
				"OS":             "$Info.OS",
				"Kernel":         "$Info.Kernel",
				"OracleCluster":  "$Info.OracleCluster",
				"SunCluster":     "$Info.SunCluster",
				"VeritasCluster": "$Info.VeritasCluster",
				"Virtual":        "$Info.Virtual",
				"Type":           "$Info.Type",
				"CPUThreads":     "$Info.CPUThreads",
				"CPUCores":       "$Info.CPUCores",
				"Socket":         "$Info.Socket",
				"MemTotal":       "$Info.MemTotal",
				"SwapTotal":      "$Info.SwapTotal",
				"CPUModel":       "$Info.CPUModel",
			})),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
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

// GetHost fetch all informations about a host in the database
func (md *MongoDatabase) GetHost(hostname string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			mu.APMatch(bson.M{
				"Hostname": hostname,
			}),
			mu.APLookupPipeline("alerts", bson.M{"hn": "$Hostname"}, "Alerts", mu.MAPipeline(
				mu.APMatch(bson.M{
					"$expr": mu.APOEqual("$OtherInfo.Hostname", "$$hn"),
				}),
			)),
			mu.APLookupPipeline("hosts", bson.M{"hn": "$Hostname"}, "VM", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APUnwind("$Extra.Clusters"),
				mu.APReplaceWith("$Extra.Clusters"),
				mu.APUnwind("$VMs"),
				mu.APReplaceWith("$VMs"),
				mu.APMatch(bson.M{
					"$expr": mu.APOEqual("$Hostname", "$$hn"),
				}),
				mu.APLimit(1),
			)),
			mu.APSet(bson.M{
				"VM": mu.APOArrayElemAt("$VM", 0),
			}),
			mu.APAddFields(bson.M{
				"Cluster":      mu.APOIfNull("$VM.ClusterName", nil),
				"PhysicalHost": mu.APOIfNull("$VM.PhysicalHost", nil),
			}),
			mu.APUnset("VM"),
			mu.APLookupPipeline(
				"hosts",
				bson.M{
					"hn": "$Hostname",
					"ca": "$CreatedAt",
				},
				"History",
				mu.MAPipeline(
					mu.APMatch(bson.M{
						"$expr": mu.APOAnd(mu.APOEqual("$Hostname", "$$hn"), mu.APOGreaterOrEqual("$$ca", "$CreatedAt")),
					}),
					mu.APProject(bson.M{
						"CreatedAt":                    1,
						"Extra.Databases.Name":         1,
						"Extra.Databases.Used":         1,
						"Extra.Databases.SegmentsSize": 1,
					}),
				),
			),
			mu.APSet(bson.M{
				"Extra.Databases": mu.APOMap(
					"$Extra.Databases",
					"db",
					mu.APOMergeObjects(
						"$$db",
						bson.M{
							"Changes": mu.APOFilter(
								mu.APOMap("$History", "hh", mu.APOMergeObjects(
									bson.M{"Updated": "$$hh.CreatedAt"},
									mu.APOArrayElemAt(mu.APOFilter("$$hh.Extra.Databases", "hdb", mu.APOEqual("$$hdb.Name", "$$db.Name")), 0),
								)),
								"time_frame",
								"$$time_frame.SegmentsSize",
							),
						},
					),
				),
			}),
			mu.APUnset(
				"Extra.Databases.Changes.Name",
				"History.Extra",
			),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return nil, utils.AerrHostNotFound
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}

// FindHostData find the current hostdata with a certain hostname
func (md *MongoDatabase) FindHostData(hostname string) (map[string]interface{}, utils.AdvancedErrorInterface) {
	//Find the hostdata
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").FindOne(context.TODO(), bson.M{
		"Hostname": hostname,
		"Archived": false,
	})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, nil
	} else if res.Err() != nil {
		return nil, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
	}

	//Decode the data
	var out map[string]interface{}
	if err := res.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
	}

	//Return it!
	return out, nil
}

// ReplaceHostData adds a new hostdata to the database
func (md *MongoDatabase) ReplaceHostData(hostData map[string]interface{}) utils.AdvancedErrorInterface {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").ReplaceOne(context.TODO(),
		bson.M{
			"_id": hostData["_id"],
		},
		hostData,
	)
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return nil
}
