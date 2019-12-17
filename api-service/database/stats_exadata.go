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
func (md *MongoDatabase) GetTotalExadataMemorySizeStats(location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

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
	// //Next the cursor. If there is no document return a empty document
	// hasNext := cur.Next(context.TODO())
	// if !hasNext {
	// 	return 0, nil
	// }

	// //Decode the document
	// if err := cur.Decode(&out); err != nil {
	// 	return 0, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	// }

	// return out["value"], nil
}
