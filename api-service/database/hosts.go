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

// SearchCurrentHosts search current hosts
func (md *MongoDatabase) SearchCurrentHosts(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"archived": false,
			}),
			mu.APSearchFilterStage([]string{
				"hostname",
				"extra.databases.name",
				"extra.databases.unique_name",
				"extra.clusters.name",
			}, keywords),
			mu.APLookupSimple("currentClusters", "hostname", "cluster.vms.hostname", "cluster"),
			mu.APSet(bson.M{
				"cluster": mu.APOArrayElemAt("$cluster", 0),
			}),
			mu.APSet(bson.M{
				"cluster": mu.APOArrayElemAt(mu.APOFilter("$cluster.cluster.vms", "vm", mu.APOEqual("$$vm.hostname", "$hostname")), 0),
			}),
			mu.APAddFields(bson.M{
				"cluster":       mu.APOIfNull("$cluster.cluster_name", nil),
				"physical_host": mu.APOIfNull("$cluster.physical_host", nil),
			}),
			mu.APOptionalStage(!full, mu.APProject(bson.M{
				"hostname":        true,
				"location":        true,
				"environment":     true,
				"host_type":       true,
				"cluster":         true,
				"physical_host":   true,
				"created_at":      true,
				"databases":       true,
				"os":              "$info.os",
				"kernel":          "$info.kernel",
				"oracle_cluster":  "$info.oracle_cluster",
				"sun_cluster":     "$info.sun_cluster",
				"veritas_cluster": "$info.veritas_cluster",
				"virtual":         "$info.virtual",
				"type":            "$info.type",
				"cpu_threads":     "$info.cpu_threads",
				"cpu_cores":       "$info.cpu_cores",
				"socket":          "$info.socket",
				"mem_total":       "$info.memory_total",
				"swap_total":      "$info.swap_total",
				"cpu_model":       "$info.cpu_model",
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

// GetCurrentHost fetch all informations about a current host in the database
func (md *MongoDatabase) GetCurrentHost(hostname string) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"archived": false,
				"hostname": hostname,
			}),
			mu.APLookupSimple("alerts", "hostname", "other_info.hostname", "alerts"),
			mu.APLookupSimple("currentClusters", "hostname", "cluster.vms.hostname", "cluster"),
			mu.APSet(bson.M{
				"cluster": mu.APOArrayElemAt("$cluster", 0),
			}),
			mu.APSet(bson.M{
				"cluster": mu.APOArrayElemAt(
					mu.APOFilter("$cluster.cluster.vms", "vm", mu.APOEqual("$$vm.hostname", "$hostname")),
					0,
				),
			}),
			mu.APAddFields(bson.M{
				"cluster":       mu.APOIfNull("$cluster.cluster_name", nil),
				"physical_host": mu.APOIfNull("$cluster.physical_host", nil),
			}),
			mu.APLookupPipeline(
				"hosts",
				bson.M{
					"hn": "$hostname",
				},
				"history",
				mu.MAPipeline(
					mu.APMatch(bson.M{
						"$expr": mu.APOEqual("$hostname", "$$hn"),
					}),
					mu.APProject(bson.M{
						"created_at":                    1,
						"extra.databases.name":          1,
						"extra.databases.used":          1,
						"extra.databases.segments_size": 1,
					}),
				),
			),
			mu.APSet(bson.M{
				"extra.databases": mu.APOMap(
					"$extra.databases",
					"db",
					mu.APOMergeObjects(
						"$$db",
						bson.M{
							"changes": mu.APOFilter(
								mu.APOMap("$history", "hh", mu.APOMergeObjects(
									bson.M{"updated": "$$hh.created_at"},
									mu.APOArrayElemAt(mu.APOFilter("$$hh.extra.databases", "hdb", mu.APOEqual("$$hdb.name", "$$db.name")), 0),
								)),
								"time_frame",
								"$$time_frame.segments_size",
							),
						},
					),
				),
			}),
			mu.APUnset(
				"extra.databases.changes.name",
				"history.extra",
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
