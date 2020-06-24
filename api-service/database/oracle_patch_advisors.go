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
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// SearchOracleDatabasePatchAdvisors search patch advisors
func (md *MongoDatabase) SearchOracleDatabasePatchAdvisors(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, windowTime time.Time, location string, environment string, olderThan time.Time, status string) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	var out []map[string]interface{} = make([]map[string]interface{}, 0)
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"Hostname":    1,
				"Environment": 1,
				"Location":    1,
				"CreatedAt":   1,
				"Database":    "$Features.Oracle.Database.Databases",
			}),
			mu.APSearchFilterStage([]interface{}{"$Hostname", "$Database.Name"}, keywords),
			mu.APProject(bson.M{
				"Hostname":         true,
				"Location":         true,
				"Environment":      true,
				"CreatedAt":        true,
				"Database.Name":    true,
				"Database.Version": true,
				"Database.PSUs":    true,
			}),
			mu.APSet(bson.M{
				"Database.PSUs": mu.APOReduce(
					mu.APOMap("$Database.PSUs", "psu", mu.APOMergeObjects(
						"$$psu",
						bson.M{
							"Date": mu.APODateFromString("$$psu.Date", "%Y-%m-%d"),
						},
					)),
					nil,
					mu.APOCond(mu.APOEqual("$$value", nil), "$$this", mu.APOMaxWithCmpExpr("$$value.Date", "$$this.Date", "$$value", "$$this")),
				),
			}),
			mu.APProject(bson.M{
				"Hostname":    true,
				"Location":    true,
				"Environment": true,
				"CreatedAt":   true,
				"Dbname":      "$Database.Name",
				"Dbver":       "$Database.Version",
				"Description": mu.APOCond("$Database.PSUs.Description", "$Database.PSUs.Description", ""),
				"Date":        mu.APOCond("$Database.PSUs.Date", "$Database.PSUs.Date", time.Unix(0, 0)),
				"Status":      mu.APOCond(mu.APOGreater("$Database.PSUs.Date", windowTime), "OK", "KO"),
			}),
			mu.APOptionalStage(status != "", mu.APMatch(bson.M{
				"Status": status,
			})),
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
		out = append(out, item)
	}
	return out, nil
}
