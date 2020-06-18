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

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTechnologiesUsage return a map that contains the number of usages for every features
func (md *MongoDatabase) GetTechnologiesUsage(location string, environment string, olderThan time.Time) (map[string]float64, utils.AdvancedErrorInterface) {
	var out map[string]float64 = make(map[string]float64)

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			FilterByOldnessSteps(olderThan),
			mu.APGroup(bson.M{
				"_id": 1,
				"Oracle/Database_HostsCount": mu.APOSum(
					mu.APOCond(mu.APOGreater(mu.APOSize(mu.APOIfNull("$Features.Oracle.Database.Databases", bson.A{})), 0), 1, 0),
				),
				"Oracle/Exadata": mu.APOSum(
					mu.APOCond(
						mu.APOEqual(bson.M{
							"$type": "$Features.Oracle.Exadata",
						}, "object"),
						1,
						0,
					),
				),
			}),
			mu.APUnset("_id"),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return out, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}
