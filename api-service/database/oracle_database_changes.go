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
	"context"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) GetOracleChanges(filter dto.GlobalFilter) ([]dto.OracleChangesDto, error) {
	ctx := context.TODO()

	result := make([]dto.OracleChangesDto, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			mu.APMatch(bson.M{
				"features.oracle.database.databases": bson.M{
					"$ne": nil,
				},
			}),
			mu.APLookupPipeline(
				"hosts",
				bson.M{
					"hn": "$hostname",
					"ca": "$createdAt",
				},
				"history",
				mu.MAPipeline(
					mu.APMatch(mu.QOExpr(mu.APOAnd(mu.APOEqual("$hostname", "$$hn"), mu.APOGreaterOrEqual("$$ca", "$createdAt")))),
					mu.APProject(bson.M{
						"createdAt": 1,
						"features.oracle.database.databases.name":          1,
						"features.oracle.database.databases.datafileSize":  1,
						"features.oracle.database.databases.segmentsSize":  1,
						"features.oracle.database.databases.allocable":     1,
						"features.oracle.database.databases.dailyCPUUsage": 1,
						"totalDailyCPUUsage":                               mu.APOSumReducer("$features.oracle.database.databases", mu.APOConvertToDoubleOrZero("$$this.dailyCPUUsage")),
					}),
				),
			),
			mu.APSet(bson.M{
				"features.oracle.database.databases": mu.APOMap(
					"$features.oracle.database.databases",
					"db",
					mu.APOMergeObjects(
						"$$db",
						bson.M{
							"changes": mu.APOFilter(
								mu.APOMap("$history", "hh", mu.APOMergeObjects(
									bson.M{"updated": "$$hh.createdAt"},
									mu.APOArrayElemAt(mu.APOFilter("$$hh.features.oracle.database.databases", "hdb", mu.APOEqual("$$hdb.name", "$$db.name")), 0),
								)),
								"time_frame",
								"$$time_frame.segmentsSize",
							),
						},
					),
				),
			}),
			mu.APUnset(
				"features.oracle.database.databases.changes.name",
				"history.features",
			),
			mu.APSet(bson.M{
				"features.oracle": mu.APOCond(mu.APOEqual("$features.oracle.database.databases", nil), nil, "$features.oracle"),
			}),
			mu.APProject(bson.M{
				"hostname": 1,
				"oracleChangesDBs": mu.APOMap(
					"$features.oracle.database.databases",
					"changesDB",
					bson.M{
						"databasename": "$$changesDB.name",
						"oracleChanges": mu.APOMap(
							"$$changesDB.changes",
							"changes",
							bson.M{
								"dailyCPUUsage": "$$changes.dailyCPUUsage",
								"segmentsSize":  "$$changes.segmentsSize",
								"updated":       "$$changes.updated",
								"datafileSize":  "$$changes.datafileSize",
								"allocable":     "$$changes.allocable",
							},
						),
					},
				),
			}),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}
