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

// SearchClusters search clusters
func (md *MongoDatabase) SearchClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	var out []map[string]interface{} = make([]map[string]interface{}, 0)

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$clusters"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"createdAt":   1,
				"cluster":     "$clusters",
			}),
			mu.APSearchFilterStage([]interface{}{"$cluster.name"}, keywords),
			mu.APProject(bson.M{
				"_id":                         true,
				"environment":                 true,
				"location":                    true,
				"createdAt":                   1,
				"hostnameAgentVirtualization": "$hostname",
				"hostname":                    true,
				"fetchEndpoint":               "$cluster.fetchEndpoint",
				"name":                        "$cluster.name",
				"type":                        "$cluster.type",
				"cpu":                         "$cluster.cpu",
				"sockets":                     "$cluster.sockets",
				"vms":                         "$cluster.vms",
				"virtualizationNodes":         mu.APOSetUnion(mu.APOMap("$cluster.vms", "vm", "$$vm.virtualizationNode")),
				"vmsCount":                    mu.APOSize("$cluster.vms"),
			}),
			mu.APLookupPipeline("hosts", bson.M{
				"vms": "$vms",
			}, "vmsErcoleAgentCount", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APProject(bson.M{
					"hostname": 1,
				}),
				mu.APMatch(mu.QOExpr(mu.APOAny("$$vms", "vm", mu.APOEqual("$$vm.hostname", "$hostname")))),
			)),
			mu.APSet(bson.M{
				"vmsErcoleAgentCount": mu.APOSize("$vmsErcoleAgentCount"),
			}),
			mu.APOptionalStage(!full, mu.APProject(bson.M{
				"_id":                         true,
				"environment":                 true,
				"location":                    true,
				"hostnameAgentVirtualization": true,
				"hostname":                    true,
				"name":                        true,
				"type":                        true,
				"cpu":                         true,
				"sockets":                     true,
				"virtualizationNodes":         mu.APOJoin("$virtualizationNodes", " "),
				"vmsCount":                    true,
				"vmsErcoleAgentCount":         true,
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
		out = append(out, item)
	}
	return out, nil
}

// GetCluster fetch all information about a cluster in the database
func (md *MongoDatabase) GetCluster(clusterName string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			mu.APUnwind("$clusters"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"createdAt":   1,
				"cluster":     "$clusters",
			}),
			mu.APMatch(bson.M{
				"cluster.name": clusterName,
			}),
			mu.APProject(bson.M{
				"_id":                         true,
				"environment":                 true,
				"location":                    true,
				"createdAt":                   1,
				"hostnameAgentVirtualization": "$hostname",
				"hostname":                    true,
				"fetchEndpoint":               "$cluster.fetchEndpoint",
				"name":                        "$cluster.name",
				"type":                        "$cluster.type",
				"cpu":                         "$cluster.cpu",
				"sockets":                     "$cluster.sockets",
				"vms":                         "$cluster.vms",
				"virtualizationNodes":         mu.APOSetUnion(mu.APOMap("$cluster.vms", "vm", "$$vm.virtualizationNode")),
				"vmsCount":                    mu.APOSize("$cluster.vms"),
			}),
			mu.APLookupPipeline("hosts", bson.M{
				"vms": "$vms",
			}, "vmsErcoleAgentCount", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APProject(bson.M{
					"hostname": 1,
				}),
				mu.APSet(bson.M{
					"virtualizationNode": mu.APOArrayElemAt(mu.APOFilter("$$vms", "vm", mu.APOEqual("$$vm.hostname", "$hostname")), 0),
				}),
				mu.APMatch(bson.M{
					"virtualizationNode": mu.QONotEqual(nil),
				}),
				mu.APSet(bson.M{
					"virtualizationNode": "$virtualizationNode.virtualizationNode",
				}),
			)),
			mu.APSet(bson.M{
				"virtualizationNodesCount": mu.APOSize("$virtualizationNodes"),
				"virtualizationNodesStats": mu.APOMap("$virtualizationNodes", "ph", mu.APOLet(bson.M{
					"vmCount":                mu.APOSize(mu.APOFilter("$vms", "vm", mu.APOEqual("$$vm.virtualizationNode", "$$ph"))),
					"vmWithErcoleAgentCount": mu.APOSize(mu.APOFilter("$vmsErcoleAgentCount", "vmea", mu.APOEqual("$$vmea.virtualizationNode", "$$ph"))),
				}, bson.M{
					"virtualizationNode":              "$$ph",
					"totalVMsCount":                   "$$vmCount",
					"totalVMsWithErcoleAgentCount":    "$$vmWithErcoleAgentCount",
					"totalVMsWithoutErcoleAgentCount": mu.APOSubtract("$$vmCount", "$$vmWithErcoleAgentCount"),
				})),
				"vmsErcoleAgentCount": mu.APOSize("$vmsErcoleAgentCount"),
			}),
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
