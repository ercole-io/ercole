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

import "go.mongodb.org/mongo-driver/bson"

func MongoAggregationOptionalStep(optional bool, step bson.M) bson.M {
	if optional {
		return step
	}
	return bson.M{"$skip": 0}
}

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
