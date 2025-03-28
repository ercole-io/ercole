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

	"github.com/ercole-io/ercole/v2/utils"
)

// SearchOracleDatabaseAddms search addms
func (md *MongoDatabase) SearchOracleDatabaseAddms(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error) {
	var out = make([]map[string]interface{}, 0)
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
				"hostname":       true,
				"location":       true,
				"environment":    true,
				"createdAt":      true,
				"database.name":  true,
				"database.addms": true,
			}),
			mu.APUnwind("$database.addms"),
			mu.APProject(bson.M{
				"hostname":       true,
				"location":       true,
				"environment":    true,
				"createdAt":      true,
				"dbname":         "$database.name",
				"action":         "$database.addms.action",
				"benefit":        "$database.addms.benefit",
				"finding":        "$database.addms.finding",
				"recommendation": "$database.addms.recommendation",
			}),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}

		out = append(out, item)
	}

	return out, nil
}
