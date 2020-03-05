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
func (md *MongoDatabase) SearchDatabases(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	var out []map[string]interface{} = make([]map[string]interface{}, 0)
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Extra.Databases"),
			mu.APAddFields(bson.M{
				"Extra.Databases.HA": mu.APOOr("$Info.SunCluster", "$Info.VeritasCluster", "$Info.OracleCluster", "$Info.AixCluster"),
			}),
			mu.APProject(bson.M{
				"Hostname":    1,
				"Environment": 1,
				"Location":    1,
				"CreatedAt":   1,
				"Database":    "$Extra.Databases",
			}),
			mu.APSearchFilterStage([]string{"Hostname", "Database.Name"}, keywords),
			mu.APAddFields(bson.M{
				"Database.Memory": mu.APOAdd(
					mu.APOConvertToDoubleOrZero("$Database.PGATarget"),
					mu.APOConvertToDoubleOrZero("$Database.SGATarget"),
					mu.APOConvertToDoubleOrZero("$Database.MemoryTarget"),
				),
				"Database.DatafileSize":     "$Database.Used",
				"Database.ArchiveLogStatus": mu.APOEqual("$Database.ArchiveLog", "ARCHIVELOG"),
				"Database.RAC": mu.APOAny("$Database.Features", "fe", mu.APOAnd(
					mu.APOEqual("$$fe.Name", "Real Application Clusters"),
					mu.APOEqual("$$fe.Status", true),
				)),
			}),
			mu.APReplaceWith(mu.APOMergeObjects("$$ROOT", "$Database")),
			mu.APUnset("Database"),
			mu.APOptionalStage(!full, mu.APProject(bson.M{
				"Hostname":         true,
				"Location":         true,
				"Environment":      true,
				"CreatedAt":        true,
				"Name":             true,
				"UniqueName":       true,
				"Version":          true,
				"Status":           true,
				"Charset":          true,
				"BlockSize":        true,
				"CPUCount":         true,
				"Work":             true,
				"Memory":           true,
				"DatafileSize":     true,
				"SegmentsSize":     true,
				"ArchiveLogStatus": true,
				"Dataguard":        true,
				"RAC":              true,
				"HA":               true,
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
