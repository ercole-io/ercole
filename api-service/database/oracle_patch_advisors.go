// Copyright (c) 2020 Sorint.lab S.p.A.
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
	"time"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

// SearchOracleDatabasePatchAdvisors search patch advisors
func (md *MongoDatabase) SearchOracleDatabasePatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) (*dto.PatchAdvisorResponse, error) {
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			ExcludeDR(),
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"createdAt":   1,
				"database":    "$features.oracle.database.databases",
			}),
			mu.APSearchFilterStage([]interface{}{"$hostname", "$database.name"}, keywords),
			mu.APProject(bson.M{
				"hostname":         true,
				"location":         true,
				"environment":      true,
				"createdAt":        true,
				"database.name":    true,
				"database.version": true,
				"database.psus":    true,
			}),
			mu.APSet(bson.M{
				"database.psus": mu.APOReduce(
					mu.APOMap("$database.psus", "psu", mu.APOMergeObjects(
						"$$psu",
						bson.M{
							"date": mu.APODateFromString("$$psu.date", "%Y-%m-%d"),
						},
					)),
					nil,
					mu.APOCond(mu.APOEqual("$$value", nil), "$$this", mu.APOMaxWithCmpExpr("$$value.date", "$$this.date", "$$value", "$$this")),
				),
			}),
			mu.APProject(bson.M{
				"hostname":     true,
				"location":     true,
				"environment":  true,
				"createdAt":    true,
				"dbname":       "$database.name",
				"dbver":        "$database.version",
				"description":  mu.APOCond("$database.psus.description", "$database.psus.description", ""),
				"date":         mu.APOCond("$database.psus.date", "$database.psus.date", time.Unix(0, 0)),
				"status":       mu.APOCond(mu.APOGreater("$database.psus.date", windowTime), "OK", "KO"),
				"fourMonths":   mu.APOGreaterOrEqual("$database.psus.date", time.Now().AddDate(0, -4, 0)),
				"sixMonths":    mu.APOGreaterOrEqual("$database.psus.date", time.Now().AddDate(0, -6, 0)),
				"twelveMonths": mu.APOGreaterOrEqual("$database.psus.date", time.Now().AddDate(0, -12, 0)),
			}),
			mu.APOptionalStage(status != "", mu.APMatch(bson.M{
				"status": status,
			})),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			PagingMetadataStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var patchAdvisorResponse dto.PatchAdvisorResponse

	cur.Next(context.TODO())

	if err := cur.Decode(&patchAdvisorResponse); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return &patchAdvisorResponse, nil
}
