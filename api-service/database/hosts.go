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
	"go.mongodb.org/mongo-driver/bson"
)

// SearchCurrentHosts search current hosts
func (md *MongoDatabase) SearchCurrentHosts(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		utils.MongoAggegationPipeline(
			utils.MongoAggregationOptionalStep(location != "", bson.M{"$match": bson.M{
				"location": location,
			}}),
			utils.MongoAggregationOptionalStep(environment != "", bson.M{"$match": bson.M{
				"environment": environment,
			}}),
			bson.M{"$match": bson.M{
				"archived": false,
			}},
			utils.MongoAggregationSearchFilterStep([]string{
				"hostname",
				"extra.databases.name",
				"extra.databases.unique_name",
				"extra.clusters.name",
			}, keywords),
			bson.M{"$lookup": bson.M{
				"from":         "currentClusters",
				"localField":   "hostname",
				"foreignField": "cluster.vms.hostname",
				"as":           "cluster",
			}},
			bson.M{"$set": bson.M{
				"cluster": bson.M{
					"$arrayElemAt": bson.A{
						"$cluster",
						0,
					},
				},
			}},
			bson.M{"$set": bson.M{
				"cluster": bson.M{
					"$arrayElemAt": bson.A{
						bson.M{
							"$filter": bson.M{
								"input": "$cluster.cluster.vms",
								"as":    "vm",
								"cond": bson.M{
									"$eq": bson.A{
										"$$vm.hostname",
										"$hostname",
									},
								},
							},
						},
						0,
					},
				},
			}},
			bson.M{"$addFields": bson.M{
				"cluster": bson.M{
					"$ifNull": bson.A{
						"$cluster.cluster_name",
						nil,
					},
				},
				"physical_host": bson.M{
					"$ifNull": bson.A{
						"$cluster.physical_host",
						nil,
					},
				},
			}},
			utils.MongoAggregationOptionalStep(!full, bson.M{"$project": bson.M{
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
			}}),
			utils.MongoAggregationOptionalSortingStep(sortBy, sortDesc),
			utils.MongoAggregationOptionalPagingStep(page, pageSize),
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
		utils.MongoAggegationPipeline(
			bson.M{"$match": bson.M{
				"archived": false,
				"hostname": hostname,
			}},
			bson.M{"$lookup": bson.M{
				"from":         "alerts",
				"localField":   "hostname",
				"foreignField": "other_info.hostname",
				"as":           "alerts",
			}},
			bson.M{"$lookup": bson.M{
				"from":         "currentClusters",
				"localField":   "hostname",
				"foreignField": "cluster.vms.hostname",
				"as":           "cluster",
			}},
			bson.M{"$set": bson.M{
				"cluster": bson.M{
					"$arrayElemAt": bson.A{
						"$cluster",
						0,
					},
				},
			}},
			bson.M{"$set": bson.M{
				"cluster": bson.M{
					"$arrayElemAt": bson.A{
						bson.M{
							"$filter": bson.M{
								"input": "$cluster.cluster.vms",
								"as":    "vm",
								"cond": bson.M{
									"$eq": bson.A{
										"$$vm.hostname",
										"$hostname",
									},
								},
							},
						},
						0,
					},
				},
			}},
			bson.M{"$addFields": bson.M{
				"cluster": bson.M{
					"$ifNull": bson.A{
						"$cluster.cluster_name",
						nil,
					},
				},
				"physical_host": bson.M{
					"$ifNull": bson.A{
						"$cluster.physical_host",
						nil,
					},
				},
			}},
			bson.M{"$lookup": bson.M{
				"from": "hosts",
				"let": bson.M{
					"hn": "$hostname",
				},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{
						"$expr": bson.M{
							"$eq": bson.A{"$hostname", "$$hn"},
						},
					}},
					bson.M{"$project": bson.M{
						"created_at":                    1,
						"extra.databases.name":          1,
						"extra.databases.used":          1,
						"extra.databases.segments_size": 1,
					}},
				},
				"as": "history",
			}},
			bson.M{"$set": bson.M{
				"extra.databases": bson.M{
					"$map": bson.M{
						"input": "$extra.databases",
						"as":    "db",
						"in": bson.M{
							"$mergeObjects": bson.A{
								"$$db",
								bson.M{
									"changes": bson.M{
										"$filter": bson.M{
											"input": bson.M{"$map": bson.M{
												"input": "$history",
												"as":    "hh",
												"in": bson.M{
													"$mergeObjects": bson.A{
														bson.M{"updated": "$$hh.created_at"},
														bson.M{"$arrayElemAt": bson.A{
															bson.M{
																"$filter": bson.M{
																	"input": "$$hh.extra.databases",
																	"as":    "hdb",
																	"cond":  bson.M{"$eq": bson.A{"$$hdb.name", "$$db.name"}},
																},
															},
															0,
														}},
													},
												},
											}},
											"as":   "time_frame",
											"cond": "$$time_frame.segments_size",
										},
									},
								},
							},
						},
					},
				},
			}},
			bson.M{"$unset": bson.A{
				"extra.databases.changes.name",
				"history.extra",
			}},
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
