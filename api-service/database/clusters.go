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
			mu.APUnwind("$Clusters"),
			mu.APProject(bson.M{
				"Hostname":    1,
				"Environment": 1,
				"Location":    1,
				"CreatedAt":   1,
				"Cluster":     "$Clusters",
			}),
			mu.APSearchFilterStage([]interface{}{"$Cluster.Name"}, keywords),
			mu.APProject(bson.M{
				"_id":                         true,
				"Environment":                 true,
				"Location":                    true,
				"CreatedAt":                   1,
				"HostnameAgentVirtualization": "$Hostname",
				"Hostname":                    true,
				"FetchEndpoint":               "$Cluster.FetchEndpoint",
				"Name":                        "$Cluster.Name",
				"Type":                        "$Cluster.Type",
				"CPU":                         "$Cluster.CPU",
				"Sockets":                     "$Cluster.Sockets",
				"VMs":                         "$Cluster.VMs",
				"VirtualizationNodes":         mu.APOSetUnion(mu.APOMap("$Cluster.VMs", "vm", "$$vm.VirtualizationNode")),
				"VMsCount":                    mu.APOSize("$Cluster.VMs"),
			}),
			mu.APLookupPipeline("hosts", bson.M{
				"vms": "$VMs",
			}, "VMsErcoleAgentCount", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APProject(bson.M{
					"Hostname": 1,
				}),
				mu.APMatch(mu.QOExpr(mu.APOAny("$$vms", "vm", mu.APOEqual("$$vm.Hostname", "$Hostname")))),
			)),
			mu.APSet(bson.M{
				"VMsErcoleAgentCount": mu.APOSize("$VMsErcoleAgentCount"),
			}),
			mu.APOptionalStage(!full, mu.APProject(bson.M{
				"_id":                         true,
				"Environment":                 true,
				"Location":                    true,
				"HostnameAgentVirtualization": true,
				"Hostname":                    true,
				"Name":                        true,
				"Type":                        true,
				"CPU":                         true,
				"Sockets":                     true,
				"VirtualizationNodes":         mu.APOJoin("$VirtualizationNodes", " "),
				"VMsCount":                    true,
				"VMsErcoleAgentCount":         true,
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
			mu.APUnwind("$Clusters"),
			mu.APProject(bson.M{
				"Hostname":    1,
				"Environment": 1,
				"Location":    1,
				"CreatedAt":   1,
				"Cluster":     "$Clusters",
			}),
			mu.APMatch(bson.M{
				"Cluster.Name": clusterName,
			}),
			mu.APProject(bson.M{
				"_id":                         true,
				"Environment":                 true,
				"Location":                    true,
				"CreatedAt":                   1,
				"HostnameAgentVirtualization": "$Hostname",
				"Hostname":                    true,
				"FetchEndpoint":               "$Cluster.FetchEndpoint",
				"Name":                        "$Cluster.Name",
				"Type":                        "$Cluster.Type",
				"CPU":                         "$Cluster.CPU",
				"Sockets":                     "$Cluster.Sockets",
				"VMs":                         "$Cluster.VMs",
				"VirtualizationNodes":         mu.APOSetUnion(mu.APOMap("$Cluster.VMs", "vm", "$$vm.VirtualizationNode")),
				"VMsCount":                    mu.APOSize("$Cluster.VMs"),
			}),
			mu.APLookupPipeline("hosts", bson.M{
				"vms": "$VMs",
			}, "VMsErcoleAgentCount", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APProject(bson.M{
					"Hostname": 1,
				}),
				mu.APSet(bson.M{
					"VirtualizationNode": mu.APOArrayElemAt(mu.APOFilter("$$vms", "vm", mu.APOEqual("$$vm.Hostname", "$Hostname")), 0),
				}),
				mu.APMatch(bson.M{
					"VirtualizationNode": mu.QONotEqual(nil),
				}),
				mu.APSet(bson.M{
					"VirtualizationNode": "$VirtualizationNode.VirtualizationNode",
				}),
			)),
			mu.APSet(bson.M{
				"VirtualizationNodesCount": mu.APOSize("$VirtualizationNodes"),
				"VirtualizationNodesStats": mu.APOMap("$VirtualizationNodes", "ph", mu.APOLet(bson.M{
					"vmCount":                mu.APOSize(mu.APOFilter("$VMs", "vm", mu.APOEqual("$$vm.VirtualizationNode", "$$ph"))),
					"vmWithErcoleAgentCount": mu.APOSize(mu.APOFilter("$VMsErcoleAgentCount", "vmea", mu.APOEqual("$$vmea.VirtualizationNode", "$$ph"))),
				}, bson.M{
					"VirtualizationNode":              "$$ph",
					"TotalVMsCount":                   "$$vmCount",
					"TotalVMsWithErcoleAgentCount":    "$$vmWithErcoleAgentCount",
					"TotalVMsWithoutErcoleAgentCount": mu.APOSubtract("$$vmCount", "$$vmWithErcoleAgentCount"),
				})),
				"VMsErcoleAgentCount": mu.APOSize("$VMsErcoleAgentCount"),
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
