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
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// SearchOracleExadata search exadata
func (md *MongoDatabase) SearchOracleExadata(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]dto.OracleExadata, error) {

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APMatch(bson.M{
				"features.oracle.exadata": mu.QONotEqual(nil),
			}),
			mu.APSearchFilterStage([]interface{}{
				"$hostname",
				"$features.oracle.exadata.components.hostname",
			}, keywords),
			mu.APProject(bson.M{
				"hostname":    true,
				"location":    true,
				"environment": true,
				"createdAt":   true,
				"dbServers": mu.APOMap(
					mu.APOFilter("$features.oracle.exadata.components", "dev", mu.APOEqual("$$dev.serverType", "DBServer")),
					"dev",
					mu.BsonOptionalExtension(full,
						bson.M{
							"hostname":           "$$dev.hostname",
							"model":              "$$dev.model",
							"swVersion":          "$$dev.swVersion",
							"runningCPUCount":    "$$dev.runningCPUCount",
							"totalCPUCount":      "$$dev.totalCPUCount",
							"memory":             "$$dev.memory",
							"runningPowerSupply": "$$dev.runningPowerSupply",
							"totalPowerSupply":   "$$dev.totalPowerSupply",
							"tempActual":         "$$dev.tempActual",
						},
						bson.M{
							"status":               "$$dev.status",
							"powerStatus":          "$$dev.powerStatus",
							"runningFanCount":      "$$dev.runningFanCount",
							"totalFanCount":        "$$dev.totalFanCount",
							"fanStatus":            "$$dev.fanStatus",
							"tempStatus":           "$$dev.tempStatus",
							"cellsrvServiceStatus": "$$dev.cellsrvServiceStatus",
							"msServiceStatus":      "$$dev.msServiceStatus",
							"rsServiceStatus":      "$$dev.rsServiceStatus",
						},
					),
				),
				"storageServers": mu.APOMap(
					mu.APOFilter("$features.oracle.exadata.components", "dev", mu.APOEqual("$$dev.serverType", "StorageServer")),
					"dev",
					mu.BsonOptionalExtension(full,
						bson.M{
							"hostname":           "$$dev.hostname",
							"model":              "$$dev.model",
							"swVersion":          "$$dev.swVersion",
							"runningCPUCount":    "$$dev.runningCPUCount",
							"totalCPUCount":      "$$dev.totalCPUCount",
							"memory":             "$$dev.memory",
							"runningPowerSupply": "$$dev.runningPowerSupply",
							"totalPowerSupply":   "$$dev.totalPowerSupply",
							"tempActual":         "$$dev.tempActual",
						},
						bson.M{
							"status":               "$$dev.status",
							"powerStatus":          "$$dev.powerStatus",
							"runningFanCount":      "$$dev.runningFanCount",
							"totalFanCount":        "$$dev.totalFanCount",
							"fanStatus":            "$$dev.fanStatus",
							"tempStatus":           "$$dev.tempStatus",
							"cellsrvServiceStatus": "$$dev.cellsrvServiceStatus",
							"msServiceStatus":      "$$dev.msServiceStatus",
							"rsServiceStatus":      "$$dev.rsServiceStatus",
							"flashcacheMode":       "$$dev.flashcacheMode",
							"cellDisks":            "$$dev.cellDisks",
						},
					),
				),
				"ibSwitches": mu.APOMap(
					mu.APOFilter("$features.oracle.exadata.components", "dev", mu.APOEqual("$$dev.serverType", "IBSwitch")),
					"dev",
					bson.M{
						"hostname":  "$$dev.hostname",
						"model":     "$$dev.model",
						"swVersion": "$$dev.swVersion",
					},
				),
			}),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	out := make([]dto.OracleExadata, 0)

	if err := cur.All(context.TODO(), &out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}
