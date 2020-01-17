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

// SearchExadata search exadata
func (md *MongoDatabase) SearchExadata(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"archived": false,
				"extra.exadata": bson.M{
					"$ne": nil,
				},
			}),
			mu.APSearchFilterStage([]string{
				"hostname",
				"extra_info.exadata.devices.hostname",
			}, keywords),
			mu.APProject(bson.M{
				"hostname":    true,
				"location":    true,
				"environment": true,
				"created_at":  true,
				"db_servers": mu.APOMap(
					mu.APOFilter("$extra.exadata.devices", "dev", mu.APOEqual("$$dev.server_type", "DBServer")),
					"dev",
					mu.BsonOptionalExtension(full,
						bson.M{
							"hostname":       "$$dev.hostname",
							"model":          "$$dev.model",
							"exa_sw_version": "$$dev.exa_sw_version",
							"cpu_enabled":    "$$dev.cpu_enabled",
							"memory":         "$$dev.memory",
							"power_count":    "$$dev.power_count",
							"temp_actual":    "$$dev.temp_actual",
						},
						bson.M{
							"status":          "$$dev.status",
							"power_status":    "$$dev.power_status",
							"fan_count":       "$$dev.fan_count",
							"fan_status":      "$$dev.fan_status",
							"temp_status":     "$$dev.temp_status",
							"cellsrv_service": "$$dev.cellserv_service",
							"ms_service":      "$$dev.ms_service",
							"rs_service":      "$$dev.rs_service",
						},
					),
				),
				"storage_servers": mu.APOMap(
					mu.APOFilter("$extra.exadata.devices", "dev", mu.APOEqual("$$dev.server_type", "StorageServer")),
					"dev",
					mu.BsonOptionalExtension(full,
						bson.M{
							"hostname":       "$$dev.hostname",
							"model":          "$$dev.model",
							"exa_sw_version": "$$dev.exa_sw_version",
							"cpu_enabled":    "$$dev.cpu_enabled",
							"memory":         "$$dev.memory",
							"power_count":    "$$dev.power_count",
							"temp_actual":    "$$dev.temp_actual",
						},
						bson.M{
							"status":          "$$dev.status",
							"power_status":    "$$dev.power_status",
							"fan_count":       "$$dev.fan_count",
							"fan_status":      "$$dev.fan_status",
							"temp_status":     "$$dev.temp_status",
							"cellsrv_service": "$$dev.cellserv_service",
							"ms_service":      "$$dev.ms_service",
							"rs_service":      "$$dev.rs_service",
							"flashcache_mode": "$$dev.flashcache_mode",
							"cell_disks":      "$$dev.cell_disks",
						},
					),
				),
				"ib_switchs": mu.APOMap(
					mu.APOFilter("$extra.exadata.devices", "dev", mu.APOEqual("$$dev.server_type", "IBSwitch")),
					"dev",
					bson.M{
						"hostname":       "$$dev.hostname",
						"model":          "$$dev.model",
						"exa_sw_version": "$$dev.exa_sw_version",
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
