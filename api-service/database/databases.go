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

// SearchDatabases search databases
func (md *MongoDatabase) SearchDatabases(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$extra.databases"),
			mu.APAddFields(bson.M{
				"extra.databases.ha": mu.APOOr("$info.sun_cluster", "$info.veritas_cluster", "$info.oracle_cluster", "$info.aix_cluster"),
			}),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"created_at":  1,
				"database":    "$extra.databases",
			}),
			mu.APSearchFilterStage([]string{"hostname", "database.name"}, keywords),
			mu.APProject(bson.M{
				"hostname":    true,
				"location":    true,
				"environment": true,
				"created_at":  true,
				"database":    true,
			}),
			mu.APAddFields(bson.M{
				"database.memory": mu.APOAdd(
					mu.APOConvertToDoubleOrZero("$database.pga_target"),
					mu.APOConvertToDoubleOrZero("$database.sga_target"),
					mu.APOConvertToDoubleOrZero("$database.memory_target"),
				),
				"datafile_size":      "$database.used",
				"archive_log_status": mu.APOEqual("$database.archive_log", "ARCHIVELOG"),
				"rac": mu.APOAny("$database.features", "fe", mu.APOAnd(
					mu.APOEqual("$$fe.name", "Real Application Clusters"),
					mu.APOEqual("$$fe.status", true),
				)),
			}),
			mu.APReplaceWith(mu.APOMergeObjects("$$ROOT", "$database")),
			mu.APOptionalStage(!full, mu.APProject(bson.M{
				"hostname":           true,
				"location":           true,
				"environment":        true,
				"created_at":         true,
				"name":               true,
				"unique_name":        true,
				"version":            true,
				"status":             true,
				"charset":            true,
				"block_size":         true,
				"cpu_count":          true,
				"work":               true,
				"memory":             true,
				"datafile_size":      true,
				"segments_size":      true,
				"archive_log_status": true,
				"dataguard":          true,
				"rac":                true,
				"ha":                 true,
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
		out = append(out, &item)
	}
	return out, nil
}
