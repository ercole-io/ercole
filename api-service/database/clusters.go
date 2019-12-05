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
	"regexp"
	"strings"

	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchCurrentClusters search current clusters
func (md *MongoDatabase) SearchCurrentClusters(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	var quotedKeywords []string
	for _, k := range keywords {
		quotedKeywords = append(quotedKeywords, regexp.QuoteMeta(k))
	}
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentClusters").Aggregate(
		context.TODO(),
		bson.A{
			bson.M{"$match": bson.M{
				"$or": bson.A{
					bson.M{"cluster.name": bson.M{
						"$regex": primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"},
					}},
				},
			}},
			bson.M{"$project": bson.M{
				"_id":                           true,
				"environment":                   true,
				"location":                      true,
				"hostname_agent_virtualization": "$hostname",
				"hostname":                      true,
				"name":                          "$cluster.name",
				"type":                          "$cluster.type",
				"cpu":                           "$cluster.cpu",
				"sockets":                       "$cluster.sockets",
				"vms":                           "$cluster.vms",
				"physical_hosts": bson.M{
					"$setUnion": bson.A{
						bson.M{
							"$map": bson.M{
								"input": "$cluster.vms",
								"as":    "vm",
								"in":    "$$vm.physical_host",
							},
						},
					},
				},
			}},
			bson.M{"$unset": bson.A{
				"vms.cluster_name",
			}},
			optionalStep(!full, bson.M{"$project": bson.M{
				"_id":                           true,
				"environment":                   true,
				"location":                      true,
				"hostname_agent_virtualization": true,
				"hostname":                      true,
				"name":                          true,
				"type":                          true,
				"cpu":                           true,
				"sockets":                       true,
				"physical_hosts": bson.M{
					"$reduce": bson.M{
						"input":        "$physical_hosts",
						"initialValue": "",
						"in": bson.M{
							"$concat": bson.A{
								"$$value",
								bson.M{
									"$cond": bson.A{
										bson.M{"$eq": bson.A{"$$value", ""}},
										"",
										" ",
									},
								},
								"$$this",
							},
						},
					},
				},
			}}),
			optionalSortingStep(sortBy, sortDesc),
			optionalPagingStep(page, pageSize),
		},
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
