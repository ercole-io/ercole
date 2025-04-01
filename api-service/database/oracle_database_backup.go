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

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) GetOracleBackupList(filter dto.GlobalFilter) ([]dto.OracleDatabaseBackupDto, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabaseBackupDto, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.backups"}},
			bson.M{"$project": bson.M{
				"hostname":                        1,
				"location":                        1,
				"environment":                     1,
				"createdAt":                       1,
				"databasename":                    "$features.oracle.database.databases.name",
				"oracleDatabaseBackup.backupType": "$features.oracle.database.databases.backups.backupType",
				"oracleDatabaseBackup.hour":       "$features.oracle.database.databases.backups.hour",
				"oracleDatabaseBackup.weekDays": bson.M{"$map": bson.M{
					"input": "$features.oracle.database.databases.backups.weekDays",
					"as":    "days",
					"in":    "$$days",
				}},
				"oracleDatabaseBackup.avgBckSize": "$features.oracle.database.databases.backups.avgBckSize",
				"oracleDatabaseBackup.retention":  "$features.oracle.database.databases.backups.retention",
			}},
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}
