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

// SearchOracleExadata search exadata
func (md *MongoDatabase) SearchOracleExadata(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{} = make([]interface{}, 0)

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"Features.Oracle.Exadata": mu.QONotEqual(nil),
			}),
			mu.APSearchFilterStage([]interface{}{
				"$Hostname",
				"$Features.Oracle.Exadata.Components.Hostname",
			}, keywords),
			mu.APProject(bson.M{
				"Hostname":    true,
				"Location":    true,
				"Environment": true,
				"CreatedAt":   true,
				"DBServers": mu.APOMap(
					mu.APOFilter("$Features.Oracle.Exadata.Components", "dev", mu.APOEqual("$$dev.ServerType", "DBServer")),
					"dev",
					mu.BsonOptionalExtension(full,
						bson.M{
							"Hostname":     "$$dev.Hostname",
							"Model":        "$$dev.Model",
							"ExaSwVersion": "$$dev.ExaSwVersion",
							"CPUEnabled":   "$$dev.CPUEnabled",
							"Memory":       "$$dev.Memory",
							"PowerCount":   "$$dev.PowerCount",
							"TempActual":   "$$dev.TempActual",
						},
						bson.M{
							"Status":         "$$dev.Status",
							"PowerStatus":    "$$dev.PowerStatus",
							"FanCount":       "$$dev.FanCount",
							"FanStatus":      "$$dev.FanStatus",
							"TempStatus":     "$$dev.TempStatus",
							"CellsrvService": "$$dev.CellsrvService",
							"MsService":      "$$dev.MsService",
							"RsService":      "$$dev.RsService",
						},
					),
				),
				"StorageServers": mu.APOMap(
					mu.APOFilter("$Features.Oracle.Exadata.Components", "dev", mu.APOEqual("$$dev.ServerType", "StorageServer")),
					"dev",
					mu.BsonOptionalExtension(full,
						bson.M{
							"Hostname":     "$$dev.Hostname",
							"Model":        "$$dev.Model",
							"ExaSwVersion": "$$dev.ExaSwVersion",
							"CPUEnabled":   "$$dev.CPUEnabled",
							"Memory":       "$$dev.Memory",
							"PowerCount":   "$$dev.PowerCount",
							"TempActual":   "$$dev.TempActual",
						},
						bson.M{
							"Status":         "$$dev.Status",
							"PowerStatus":    "$$dev.PowerStatus",
							"FanCount":       "$$dev.FanCount",
							"FanStatus":      "$$dev.FanStatus",
							"TempStatus":     "$$dev.TempStatus",
							"CellsrvService": "$$dev.CellsrvService",
							"MsService":      "$$dev.MsService",
							"RsService":      "$$dev.RsService",
							"FlashcacheMode": "$$dev.FlashcacheMode",
							"CellDisks":      "$$dev.CellDisks",
						},
					),
				),
				"IBSwitches": mu.APOMap(
					mu.APOFilter("$Features.Oracle.Exadata.Components", "dev", mu.APOEqual("$$dev.ServerType", "IBSwitch")),
					"dev",
					bson.M{
						"Hostname":     "$$dev.Hostname",
						"Model":        "$$dev.Model",
						"ExaSwVersion": "$$dev.ExaSwVersion",
					},
				),
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
