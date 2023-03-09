// Copyright (c) 2023 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// GetHostsCountUsingTechnologies return a map that contains the number of usages for every features
func (md *MongoDatabase) GetHostsCountUsingTechnologies(location string, environment string, olderThan time.Time) (map[string]float64, error) {
	var out map[string]float64 = make(map[string]float64)

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(location, environment),
			FilterByOldnessSteps(olderThan),
			mu.APGroup(bson.M{
				"_id": 1,
				model.TechnologyOracleDatabase: mu.APOSum(
					mu.APOCond(mu.APOGreater(mu.APOSize(mu.APOIfNull("$features.oracle.database.databases", bson.A{})), 0), 1, 0),
				),
				model.TechnologyOracleExadata: mu.APOSum(
					mu.APOCond(mu.APOGreater(mu.APOSize(mu.APOIfNull("$features.oracle.exadata.components", bson.A{})), 0), 1, 0),
				),
				model.TechnologyOracleMySQL: mu.APOSum(
					mu.APOCond(mu.APOGreater(mu.APOSize(mu.APOIfNull("$features.mysql.instances", bson.A{})), 0), 1, 0),
				),
				model.TechnologyMicrosoftSQLServer: mu.APOSum(
					mu.APOCond(mu.APOGreater(mu.APOSize(mu.APOIfNull("$features.microsoft.sqlServer.instances", bson.A{})), 0), 1, 0),
				),
				model.TechnologyPostgreSQLPostgreSQL: mu.APOSum(
					mu.APOCond(mu.APOGreater(mu.APOSize(mu.APOIfNull("$features.postgresql.instances", bson.A{})), 0), 1, 0),
				),
				model.TechnologyMongoDBMongoDB: mu.APOSum(
					mu.APOCond(mu.APOGreater(mu.APOSize(mu.APOIfNull("$features.mongodb.dbStats", bson.A{})), 0), 1, 0),
				),
			}),
			mu.APUnset("_id"),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return out, nil
	}

	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return out, nil
}
