// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"math"
	"strings"
	"time"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// FilterByLocationAndEnvironmentSteps return the steps required to filter the data by the location and environment field.
func FilterByLocationAndEnvironmentSteps(location string, environment string) interface{} {
	return bson.A{
		mu.APOptionalStage(location != "" && location != model.AllLocation, mu.APMatch(bson.M{
			"location": bson.M{"$in": strings.Split(location, ",")},
		})),
		mu.APOptionalStage(environment != "", mu.APMatch(bson.M{
			"environment": environment,
		})),
	}
}

func FilterByOldnessSteps(olderThan time.Time) bson.A {
	return mu.MAPipeline(
		mu.APOptionalStage(olderThan == utils.MAX_TIME, mu.APMatch(bson.M{
			"dismissedAt": nil,
			"archived":    false,
		})),
		mu.APOptionalStage(olderThan != utils.MAX_TIME, bson.A{
			mu.APMatch(bson.M{
				"createdAt": mu.QOLessThanOrEqual(olderThan),
			}),

			bson.M{
				"$sort": bson.M{"createdAt": 1},
			},
			bson.M{
				"$group": bson.M{
					"_id": "$hostname",
					"hostdata": bson.M{
						"$max": bson.M{
							"$mergeObjects": bson.A{
								bson.M{"createdAt": "$createdAt"},
								"$$ROOT",
							},
						},
					},
				},
			},

			bson.M{
				"$replaceRoot": bson.M{

					"newRoot": "$hostdata",
				},
			},
		}),
	)
}

func ExcludeDR() bson.A {
	return bson.A{
		bson.M{
			"$match": bson.M{"isDR": false},
		},
	}
}

func AddHardwareAbstraction(field string) bson.A {
	return mu.MAPipeline(mu.APAddFields(bson.M{
		field: mu.APOOr("$clusterMembershipStatus.oracleClusterware", "$clusterMembershipStatus.veritasClusterServer", "$clusterMembershipStatus.sunCluster", "$clusterMembershipStatus.hacmp"),
	}))
}

func AddAssociatedClusterNameAndVirtualizationNode(olderThan time.Time) bson.A {
	return mu.MAPipeline(
		mu.APLookupPipeline("hosts", bson.M{"hn": "$hostname"}, "vm", mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			mu.APUnwind("$clusters"),
			mu.APReplaceWith("$clusters"),
			mu.APUnwind("$vms"),
			mu.APSet(bson.M{
				"vms.clusterName": "$name",
			}),
			mu.APReplaceWith("$vms"),
			mu.APMatch(mu.QOExpr(mu.APOEqual("$hostname", "$$hn"))),
			mu.APLimit(1),
		)),
		mu.APSet(bson.M{
			"vm": mu.APOArrayElemAt("$vm", 0),
		}),
		mu.APAddFields(bson.M{
			"cluster":            mu.APOIfNull("$vm.clusterName", nil),
			"virtualizationNode": mu.APOIfNull("$vm.virtualizationNode", nil),
		}),
		mu.APUnset("vm"),
	)
}

// PagingMetadataStage insert PagingStage
func PagingMetadataStage(page int, size int) interface{} {
	if page >= 0 && size > 0 {
		return mu.APOptionalPagingStage(page, size)
	}

	return mu.APOptionalPagingStage(0, math.MaxInt64)
}

func FindByHostname(hostname string) bson.A {
	return mu.MAPipeline(mu.APOptionalStage(hostname != "", mu.APMatch(bson.M{"hostname": hostname})))
}

func FilterExadata(filter dto.GlobalFilter, hidden bool) bson.D {
	res := bson.D{
		{Key: "$or",
			Value: bson.A{
				bson.D{{Key: "hidden", Value: bson.D{{Key: "$exists", Value: hidden}}}},
				bson.D{{Key: "hidden", Value: hidden}},
			},
		},
	}

	if hidden {
		res = bson.D{
			{Key: "hidden", Value: hidden},
		}
	}

	if filter.Location != "" {
		res = append(res, bson.E{Key: "location", Value: filter.Location})
	}

	if filter.Environment != "" {
		res = append(res, bson.E{Key: "environment", Value: filter.Environment})
	}

	if filter.OlderThan != utils.MAX_TIME {
		res = append(res, bson.E{Key: "createdAt", Value: bson.M{"$lte": filter.OlderThan}})
	}

	return res
}

func filterExistingLocations(locations []string) bson.A {
	if len(locations) > 0 && !utils.Contains(locations, "") {
		return bson.A{
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "location",
							Value: bson.D{
								{Key: "$in",
									Value: locations,
								},
							},
						},
					},
				},
			},
		}
	}

	return bson.A{}
}
