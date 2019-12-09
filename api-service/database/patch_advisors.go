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
	"time"

	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchCurrentPatchAdvisors search current patch advisors
func (md *MongoDatabase) SearchCurrentPatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	var quotedKeywords []string
	for _, k := range keywords {
		quotedKeywords = append(quotedKeywords, regexp.QuoteMeta(k))
	}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		bson.A{
			bson.M{"$match": bson.M{
				"$or": bson.A{
					bson.M{"hostname": bson.M{
						"$regex": primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"},
					}},
					bson.M{"database.name": bson.M{
						"$regex": primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"},
					}},
				},
			}},
			bson.M{"$project": bson.M{
				"hostname":           true,
				"location":           true,
				"environment":        true,
				"created_at":         true,
				"database.name":      true,
				"database.version":   true,
				"database.last_psus": true,
			}},
			bson.M{"$set": bson.M{
				"database.last_psus": bson.M{
					"$reduce": bson.M{
						"input": bson.M{
							"$map": bson.M{
								"input": "$database.last_psus",
								"as":    "psu",
								"in": bson.M{
									"$mergeObjects": bson.A{
										"$$psu",
										bson.M{
											"date": bson.M{
												"$dateFromString": bson.M{
													"dateString": "$$psu.date",
													"format":     "%Y-%m-%d",
												},
											},
										},
									},
								},
							},
						},
						"initialValue": nil,
						"in": bson.M{
							"$cond": bson.M{
								"if": bson.M{
									"$eq": bson.A{
										"$$value",
										nil,
									},
								},
								"then": "$$this",
								"else": bson.M{
									"$cond": bson.M{
										"if": bson.M{
											"$gt": bson.A{
												"$$value.date",
												"$$this.date",
											},
										},
										"then": "$$value",
										"else": "$$this",
									},
								},
							},
						},
					},
				},
			}},
			bson.M{"$project": bson.M{
				"hostname":    true,
				"location":    true,
				"environment": true,
				"created_at":  true,
				"dbname":      "$database.name",
				"dbver":       "$database.version",
				"description": bson.M{
					"$cond": bson.M{
						"if":   "$database.last_psus.description",
						"then": "$database.last_psus.description",
						"else": "",
					},
				},
				"date": bson.M{
					"$cond": bson.M{
						"if":   "$database.last_psus.date",
						"then": "$database.last_psus.date",
						"else": time.Unix(0, 0),
					},
				},
				"status": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$gt": bson.A{
								"$database.last_psus.date",
								windowTime,
							},
						},
						"then": "OK",
						"else": "KO",
					},
				},
			}},
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
