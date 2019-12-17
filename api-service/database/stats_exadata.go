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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTotalExadataMemorySizeStats return the total size of memory of exadata
func (md *MongoDatabase) GetTotalExadataMemorySizeStats(location string, environment string) (float32, utils.AdvancedErrorInterface) {
	var out map[string]float32

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		utils.MongoAggegationPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			bson.M{"$match": bson.M{
				"archived": false,
			}},
			bson.M{"$group": bson.M{
				"_id": 0,
				"value": bson.M{
					"$sum": bson.M{
						"$reduce": bson.M{
							"input":        "$extra.exadata.devices",
							"initialValue": 0,
							"in": bson.M{
								"$add": bson.A{
									"$$value",
									bson.M{
										"$let": bson.M{
											"vars": bson.M{
												"match": bson.M{
													"$regexFind": bson.M{
														"input": "$$this.memory",
														"regex": primitive.Regex{Pattern: "^(\\d+)GB$", Options: "i"},
													},
												},
											},
											"in": bson.M{
												"$convert": bson.M{
													"input": bson.M{"$arrayElemAt": bson.A{
														"$$match.captures",
														0,
													}},
													"to":     "double",
													"onNull": 0,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}},
		),
	)
	if err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out["value"], nil
}

// GetTotalExadataCPUStats return the total cpu of exadata
func (md *MongoDatabase) GetTotalExadataCPUStats(location string, environment string) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}

	//Calculate the stats
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		utils.MongoAggegationPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			bson.M{"$match": bson.M{
				"archived": false,
			}},
			bson.M{"$project": bson.M{
				"value": bson.M{
					"$reduce": bson.M{
						"input":        "$extra.exadata.devices",
						"initialValue": bson.M{"enabled": 0, "total": 0},
						"in": bson.M{
							"$let": bson.M{
								"vars": bson.M{
									"match": bson.M{
										"$regexFind": bson.M{
											"input": "$$this.cpu_enabled",
											"regex": primitive.Regex{Pattern: "^(\\d+)/(\\d+)$", Options: "i"},
										},
									},
								},
								"in": bson.M{
									"enabled": bson.M{
										"$add": bson.A{
											"$$value.enabled",
											bson.M{
												"$convert": bson.M{
													"input": bson.M{"$arrayElemAt": bson.A{
														"$$match.captures",
														0,
													}},
													"to":     "double",
													"onNull": 0,
												},
											},
										},
									},
									"total": bson.M{
										"$add": bson.A{
											"$$value.total",
											bson.M{
												"$convert": bson.M{
													"input": bson.M{"$arrayElemAt": bson.A{
														"$$match.captures",
														0,
													}},
													"to":     "double",
													"onNull": 0,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}},

			bson.M{"$group": bson.M{
				"_id": 0,
				"enabled": bson.M{
					"$sum": "$value.enabled",
				},
				"total": bson.M{
					"$sum": "$value.total",
				},
			}},
			bson.M{"$unset": bson.A{"_id"}},
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return 0, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}
