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

package utils

import (
	"reflect"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoAggegationPipeline return a aggregation pipeline joining the steps that could be a single step or a slice of steps
func MongoAggegationPipeline(steps ...interface{}) bson.A {
	out := bson.A{}
	for _, step := range steps {
		if reflect.TypeOf(step).Kind() == reflect.Slice {
			for _, item := range SliceToSliceOfInterface(step) {
				out = append(out, item)
			}
		} else {
			out = append(out, step)
		}
	}

	return out
}

// MongoAggregationOptionalStep return the step if optional is true, otherwise return a null step
func MongoAggregationOptionalStep(optional bool, step interface{}) interface{} {
	if optional {
		return step
	}
	return bson.M{"$skip": 0}
}

// MongoAggregationOptionalSortingStep return a step that sort documents by the criteria in the params
func MongoAggregationOptionalSortingStep(sortBy string, sortDesc bool) bson.M {
	if sortBy == "" {
		return bson.M{"$skip": 0}
	}

	sortOrder := 0
	if sortDesc {
		sortOrder = -1
	} else {
		sortOrder = 1
	}

	return bson.M{"$sort": bson.M{
		sortBy: sortOrder,
	}}
}

// MongoAggregationOptionalPagingStep return a step that turn a stream of documents into a page that contains the documents plus some metadata
func MongoAggregationOptionalPagingStep(page int, size int) bson.M {
	if page == -1 || size == -1 {
		return bson.M{"$skip": 0}
	}

	return bson.M{"$facet": bson.M{
		"content": bson.A{
			bson.M{"$skip": page * size},
			bson.M{"$limit": size},
		},
		"metadata": bson.A{
			bson.M{"$count": "total_elements"},
			bson.M{"$addFields": bson.M{
				"total_pages": bson.M{
					"$floor": bson.M{
						"$divide": bson.A{
							"$total_elements",
							size,
						},
					},
				},
				"size": bson.M{
					"$min": bson.A{
						size,
						bson.M{"$subtract": bson.A{
							"$total_elements",
							size * page,
						}},
					},
				},
				"number": page,
			}},
			bson.M{"$addFields": bson.M{
				"empty": bson.M{
					"$eq": bson.A{
						"$size",
						0,
					},
				},
				"first": page == 0,
				"last": bson.M{
					"$eq": bson.A{
						page,
						bson.M{"$subtract": bson.A{
							"$total_pages",
							1,
						}},
					},
				},
			}},
		},
	}}
}

// MongoAggregationSearchFilterStep return a aggregation step that filter the documents when any field match any keyword
func MongoAggregationSearchFilterStep(fields []string, keywords []string) interface{} {
	//Build the search pattern
	quotedKeywords := []string{}
	for _, k := range keywords {
		quotedKeywords = append(quotedKeywords, regexp.QuoteMeta(k))
	}
	pattern := strings.Join(quotedKeywords, "|")

	//Build the $or conditions
	conditions := []interface{}{}
	for _, f := range fields {
		conditions = append(conditions, bson.M{f: bson.M{
			"$regex": primitive.Regex{Pattern: pattern, Options: "i"},
		}})
	}

	//Return the matching step
	return bson.M{"$match": bson.M{
		"$or": conditions,
	}}
}

// OptionalBsonMExtension return a bson with the same key-values pairs as the orig and extension, if extend is true, otherwise return the orig bson
func OptionalBsonMExtension(extend bool, orig bson.M, extension bson.M) bson.M {
	if extend {
		for k, v := range extension {
			orig[k] = v
		}
	}

	return orig
}
