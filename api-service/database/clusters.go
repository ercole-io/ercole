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

// SearchClusters search clusters
func (md *MongoDatabase) SearchClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	var out []map[string]interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Clusters"),
			mu.APProject(bson.M{
				"Hostname":    1,
				"Environment": 1,
				"Location":    1,
				"CreatedAt":   1,
				"Cluster":     "$Extra.Clusters",
			}),
			mu.APSearchFilterStage([]string{"Cluster.Name"}, keywords),
			mu.APProject(bson.M{
				"_id":                         true,
				"Environment":                 true,
				"Location":                    true,
				"CreatedAt":                   1,
				"HostnameAgentVirtualization": "$Hostname",
				"Hostname":                    true,
				"Name":                        "$Cluster.Name",
				"Type":                        "$Cluster.Type",
				"CPU":                         "$Cluster.CPU",
				"Sockets":                     "$Cluster.Sockets",
				"VMs":                         "$Cluster.VMs",
				"PhysicalHosts":               mu.APOSetUnion(mu.APOMap("$Cluster.VMs", "vm", "$$vm.PhysicalHost")),
			}),
			mu.APUnset("VMs.ClusterName"),
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
				"PhysicalHosts":               mu.APOJoin("$PhysicalHosts", " "),
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
