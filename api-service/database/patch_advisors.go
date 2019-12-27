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
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
)

// SearchCurrentPatchAdvisors search current patch advisors
func (md *MongoDatabase) SearchCurrentPatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("currentDatabases").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APSearchFilterStage([]string{"hostname", "database.name"}, keywords),
			mu.APProject(bson.M{
				"hostname":           true,
				"location":           true,
				"environment":        true,
				"created_at":         true,
				"database.name":      true,
				"database.version":   true,
				"database.last_psus": true,
			}),
			mu.APSet(bson.M{
				"database.last_psus": mu.APOReduce(
					mu.APOMap("$database.last_psus", "psu", mu.APOMergeObjects(
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
				"hostname":    true,
				"location":    true,
				"environment": true,
				"created_at":  true,
				"dbname":      "$database.name",
				"dbver":       "$database.version",
				"description": mu.APOCond("$database.last_psus.description", "$database.last_psus.description", ""),
				"date":        mu.APOCond("$database.last_psus.date", "$database.last_psus.date", time.Unix(0, 0)),
				"status":      mu.APOCond(mu.APOGreater("$database.last_psus.date", windowTime), "OK", "KO"),
			}),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
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
