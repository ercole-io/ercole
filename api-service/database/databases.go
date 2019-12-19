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

// SearchCurrentDatabases search current databases
func (md *MongoDatabase) SearchCurrentDatabases(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		utils.MongoAggegationPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			utils.MongoAggregationSearchFilterStep([]string{"hostname", "database.name"}, keywords),
			bson.M{"$project": bson.M{
				"hostname":    true,
				"location":    true,
				"environment": true,
				"created_at":  true,
				"database":    true,
			}},
			bson.M{"$addFields": bson.M{
				"database.memory": utils.MongoAggregationAdd(
					bson.M{
						"$convert": bson.M{
							"input":   "$database.pga_target",
							"to":      "double",
							"onError": 0,
						},
					},
					bson.M{
						"$convert": bson.M{
							"input":   "$database.sga_target",
							"to":      "double",
							"onError": 0,
						},
					},
					bson.M{
						"$convert": bson.M{
							"input":   "$database.memory_target",
							"to":      "double",
							"onError": 0,
						},
					},
				),
				"datafile_size": "$database.used",
				"archive_log_status": bson.M{
					"$eq": bson.A{
						"$database.archive_log",
						"ARCHIVELOG",
					},
				},
				"rac": bson.M{
					"$gt": bson.A{
						bson.M{
							"$size": bson.M{
								"$filter": bson.M{
									"input": "$database.features",
									"as":    "fe",
									"cond": bson.M{
										"$and": bson.A{
											bson.M{
												"$eq": bson.A{
													"$$fe.name",
													"Real Application Clusters",
												},
											},
											bson.M{
												"$eq": bson.A{
													"$$fe.status",
													true,
												},
											},
										},
									},
								},
							},
						},
						0,
					},
				},
			}},
			bson.M{"$replaceWith": bson.M{
				"$mergeObjects": bson.A{
					"$$ROOT",
					"$database",
				},
			}},
			utils.MongoAggregationOptionalStep(!full, bson.M{"$project": bson.M{
				"hostname":           true,
				"location":           true,
				"environment":        true,
				"created_at":         true,
				"name":               true,
				"unique_name":        true,
				"version":            true,
				"status":             true,
				"charset":            true,
				"block_size":         true,
				"cpu_count":          true,
				"work":               true,
				"memory":             true,
				"datafile_size":      true,
				"segments_size":      true,
				"archive_log_status": true,
				"dataguard":          true,
				"rac":                true,
				"ha":                 true,
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
