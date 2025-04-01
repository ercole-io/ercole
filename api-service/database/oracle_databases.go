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
	"math"
	"time"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

// SearchOracleDatabases search databases
func (md *MongoDatabase) SearchOracleDatabases(keywords []string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) (*dto.OracleDatabaseResponse, error) {
	//Find the matching hostdata
	var oracleDatabaseResponse dto.OracleDatabaseResponse

	var pagePaging, pagePagingSize int

	if pageSize > 0 {
		pagePagingSize = pageSize
	} else {
		pagePagingSize = math.MaxInt64
	}

	if !(page >= 0) {
		pagePaging = 0
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			AddHardwareAbstraction("features.oracle.database.databases.ha"),
			mu.APProject(bson.M{
				"hostname":           1,
				"environment":        1,
				"location":           1,
				"name":               1,
				"uniqueName":         1,
				"status":             1,
				"version":            1,
				"archivelog":         1,
				"charset":            1,
				"blockSize":          1,
				"cpuCount":           1,
				"memoryTarget":       1,
				"segmentsSize":       1,
				"datafileSize":       1,
				"work":               1,
				"dataguard":          1,
				"dbID":               1,
				"role":               1,
				"database":           "$features.oracle.database.databases",
				"clusterwareVersion": "$info.clusterwareVersion",
			}),
			mu.APSearchFilterStage([]interface{}{"$hostname", "$database.name"}, keywords),
			mu.APAddFields(bson.M{
				"database.memory": mu.APOAdd(
					"$database.pgaTarget",
					"$database.sgaTarget",
					"$database.memoryTarget",
				),
				"database.rac": mu.APOAny("$database.licenses", "lic", mu.APOAnd(
					mu.APOEqual("$$lic.name", "Real Application Clusters"),
					mu.APOGreater("$$lic.count", 0),
				)),
				"database.isCDB":    "$database.isCDB",
				"database.services": "$database.services",
				"database.licenses": "$database.licenses",
			}),
			mu.APSet(bson.M{
				"database.pdbs": mu.APOCond("$database.isCDB", bson.M{
					"$concatArrays": bson.A{
						mu.APOMap("$database.pdbs", "pdb", "$$pdb"),
					},
				}, []string{}),
			}),
			mu.APReplaceWith(mu.APOMergeObjects("$$ROOT", "$database")),
			mu.APUnset("database"),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APLimit(pagePagingSize),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	err = cur.All(context.TODO(), &oracleDatabaseResponse.Content)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	if oracleDatabaseResponse.Content == nil {
		oracleDatabaseResponse.Content = []dto.OracleDatabase{}
	}

	md.Client.Database(md.Config.Mongodb.DBName)
	cur1, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),

		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"database":    "$features.oracle.database.databases",
			}),
			mu.APSearchFilterStage([]interface{}{"$hostname", "$database.name"}, keywords),
			mu.APFacet(bson.M{
				"metadata": mu.MAPipeline(
					mu.APCount("totalElements"),
				),
			},
			),
			mu.APSet(bson.M{
				"metadata": mu.APOIfNull(mu.APOArrayElemAt("$metadata", 0), bson.M{
					"totalElements": 0,
				}),
			}),

			mu.APSet(bson.M{
				"metadata.totalPages": "$metadata",
			}),
			mu.APAddFields(bson.M{
				"metadata.totalPages": mu.APOFloor(mu.APODivide("$metadata.totalElements", pagePagingSize)),
				"metadata.size":       mu.APOMin(pagePagingSize, mu.APOSubtract("$metadata.totalElements", pagePagingSize*pagePaging)),
				"metadata.number":     pagePaging,
			}),
			mu.APAddFields(bson.M{
				"metadata.empty": mu.APOEqual("$metadata.size", 0),
				"metadata.first": pagePaging == 0,
				"metadata.last":  mu.APOGreaterOrEqual(pagePaging, mu.APOSubtract("$metadata.totalPages", 1)),
			}),
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	cur1.Next(context.TODO())

	if err := cur1.Decode(&oracleDatabaseResponse); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return &oracleDatabaseResponse, nil
}

func (md *MongoDatabase) DbExist(hostname, dbname string) (bool, error) {
	filter := bson.D{
		{Key: "archived", Value: false},
		{Key: "hostname", Value: hostname},
		{Key: "features.oracle.database.databases.name", Value: dbname},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return cur > 0, nil
}
