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
func (md *MongoDatabase) SearchClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$extra.clusters"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"created_at":  1,
				"cluster":     "$extra.clusters",
			}),
			mu.APSearchFilterStage([]string{"cluster.name"}, keywords),
			mu.APProject(bson.M{
				"_id":                           true,
				"environment":                   true,
				"location":                      true,
				"created_at":                    1,
				"hostname_agent_virtualization": "$hostname",
				"hostname":                      true,
				"name":                          "$cluster.name",
				"type":                          "$cluster.type",
				"cpu":                           "$cluster.cpu",
				"sockets":                       "$cluster.sockets",
				"vms":                           "$cluster.vms",
				"physical_hosts":                mu.APOSetUnion(mu.APOMap("$cluster.vms", "vm", "$$vm.physical_host")),
			}),
			mu.APUnset("vms.cluster_name"),
			mu.APOptionalStage(!full, mu.APProject(bson.M{
				"_id":                           true,
				"environment":                   true,
				"location":                      true,
				"hostname_agent_virtualization": true,
				"hostname":                      true,
				"name":                          true,
				"type":                          true,
				"cpu":                           true,
				"sockets":                       true,
				"physical_hosts":                mu.APOJoin("$physical_hosts", " "),
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
