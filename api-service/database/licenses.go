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

// ListCurrentLicenses list current licenses
func (md *MongoDatabase) ListCurrentLicenses(full bool, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Find the informations
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APLookupPipeline("hosts", bson.M{
				"license_name": "$_id",
			}, "used", mu.MAPipeline(
				mu.APMatch(bson.M{
					"archived": false,
				}),
				mu.APProject(bson.M{
					"hostname": 1,
					"databases": mu.APOReduce(
						mu.APOFilter(
							mu.APOMap("$extra.databases", "db", bson.M{
								"name": "$$db.name",
								"count": mu.APOLet(
									bson.M{
										"val": mu.APOArrayElemAt(mu.APOFilter("$$db.licenses", "lic", mu.APOEqual("$$lic.name", "$$license_name")), 0),
									},
									"$$val.count",
								),
							}),
							"db",
							mu.APOGreater("$$db.count", 0),
						),
						bson.M{"count": 0, "dbs": bson.A{}},
						bson.M{
							"count": mu.APOMax("$$value.count", "$$this.count"),
							"dbs": bson.M{
								"$concatArrays": bson.A{
									"$$value.dbs",
									bson.A{"$$this.name"},
								},
							},
						},
					),
				}),
				mu.APMatch(bson.M{
					"databases.count": bson.M{
						"$gt": 0,
					},
				}),
				mu.APLookupSimple("currentClusters", "hostname", "cluster.vms.hostname", "cluster"),
				mu.APSet(bson.M{
					"cluster": mu.APOArrayElemAt("$cluster", 0),
				}),
				// mu.APSet(bson.M{
				// 	"cluster": mu.APOArrayElemAt(
				// 		mu.APOFilter("$cluster.cluster.vms", "vm", mu.APOEqual("$$vm.hostname", "$hostname")),
				// 		0,
				// 	),
				// }),
				mu.APSet(bson.M{
					"cluster_name": "$cluster.cluster.name",
					"cluster_cpu":  "$cluster.cluster.cpu",
				}),
				mu.APUnset("cluster"),
				mu.APGroup(mu.BsonOptionalExtension(full, bson.M{
					"_id": mu.APOCond(
						"$cluster_name",
						mu.APOConcat("cluster_§$#$§_", "$cluster_name"),
						mu.APOConcat("hostname_§$#$§_", "$hostname"),
					),
					"license":     mu.APOMaxAggr("$databases.count"),
					"cluster_cpu": mu.APOMaxAggr("$cluster_cpu"),
				}, bson.M{
					"hosts": mu.APOPush(bson.M{
						"hostname":  "$hostname",
						"databases": "$databases.dbs",
					}),
				})),
				mu.APSet(bson.M{
					"license": mu.APOCond(
						"$cluster_cpu",
						mu.APODivide("$cluster_cpu", 2),
						"$license",
					),
				}),
				mu.APGroup(mu.BsonOptionalExtension(full, bson.M{
					"_id":   0,
					"value": mu.APOSum("$license"),
				}, bson.M{
					"hosts": mu.APOPush("$hosts"),
				})),
				mu.APOptionalStage(full, mu.MAPipeline(
					mu.APUnwind("$hosts"),
					mu.APUnwind("$hosts"),
					mu.APGroup(bson.M{
						"_id":   0,
						"value": mu.APOMaxAggr("$value"),
						"hosts": mu.APOPush("$hosts"),
					}),
				)),
			)),
			mu.APSet(bson.M{
				"used": mu.APOArrayElemAt("$used", 0),
			}),
			mu.APOptionalStage(full, mu.APSet(bson.M{
				"hosts": mu.APOIfNull("$used.hosts", bson.A{}),
			})),
			mu.APSet(bson.M{
				"used": mu.APOIfNull(mu.APOCeil("$used.value"), 0),
			}),
			mu.APSet(bson.M{
				"compliance": mu.APOGreaterOrEqual("$count", "$used"),
			}),
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

// SetLicenseCount set the count of a certain license
func (md *MongoDatabase) SetLicenseCount(name string, count int) utils.AdvancedErrorInterface {
	//Find the informations
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").UpdateOne(context.TODO(), bson.M{
		"_id": name,
	}, bson.M{
		"$set": bson.M{
			"count": count,
		},
	})
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Check the existance of the result
	if res.MatchedCount == 0 {
		return utils.AerrLicenseNotFound
	} else {
		return nil
	}
}
