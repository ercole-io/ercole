// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"math"
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) SearchMongoDBInstances(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.MongoDBInstanceResponse, error) {
	var mongoDBInstanceResponse dto.MongoDBInstanceResponse

	var pagePaging, pagePagingSize int

	if pageSize > 0 {
		pagePagingSize = pageSize
	} else {
		pagePagingSize = math.MaxInt64
	}

	if !(page >= 0) {
		pagePaging = 0
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.mongodb.instances"),
			mu.APUnwind("$features.mongodb.instances.dbStats"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"instance":    "$features.mongodb.instances",
			}),
			mu.APSearchFilterStage([]interface{}{"$hostname", "$name"}, keywords),
			mu.APAddFields(bson.M{
				"name":    "$instance.name",
				"dbName":  "$instance.dbStats.dbName",
				"charset": "$instance.dbStats.charset",
				"version": "$instance.version",
			}),
			mu.APReplaceWith(mu.APOMergeObjects("$$ROOT", "$instance")),
			mu.APUnset("instance"),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APLimit(pagePagingSize),
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	err = cur.All(context.TODO(), &mongoDBInstanceResponse.Content)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	if mongoDBInstanceResponse.Content == nil {
		mongoDBInstanceResponse.Content = []dto.MongoDBInstance{}
	}

	md.Client.Database(md.Config.Mongodb.DBName)
	cur1, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),

		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.mongodb.instances"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"instance":    "$features.mongodb.instances",
			}),
			mu.APSearchFilterStage([]interface{}{"$hostname", "$name"}, keywords),
			mu.APFacet(bson.M{
				"metadata": mu.MAPipeline(
					mu.APCount("totalElements"),
				),
			},
			),
			mu.APSet(bson.M{
				"metadata": mu.APOIfNull(mu.APOArrayElemAt("$metadata", 0), bson.M{
					"totalElements": 0,
				}),
			}),

			mu.APSet(bson.M{
				"metadata.totalPages": "$metadata",
			}),
			mu.APAddFields(bson.M{
				"metadata.totalPages": mu.APOFloor(mu.APODivide("$metadata.totalElements", pagePagingSize)),
				"metadata.size":       mu.APOMin(pagePagingSize, mu.APOSubtract("$metadata.totalElements", pagePagingSize*pagePaging)),
				"metadata.number":     pagePaging,
			}),
			mu.APAddFields(bson.M{
				"metadata.empty": mu.APOEqual("$metadata.size", 0),
				"metadata.first": pagePaging == 0,
				"metadata.last":  mu.APOGreaterOrEqual(pagePaging, mu.APOSubtract("$metadata.totalPages", 1)),
			}),
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	cur1.Next(context.TODO())

	if err := cur1.Decode(&mongoDBInstanceResponse); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return &mongoDBInstanceResponse, nil
}
