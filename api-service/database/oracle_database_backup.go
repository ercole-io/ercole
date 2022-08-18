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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) GetOracleBackupList() ([]dto.OracleDatabaseBackupDto, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabaseBackupDto, 0)

	query := bson.A{
		bson.M{"$match": bson.M{
			"dismissedAt": nil,
			"archived":    false,
		}},
		bson.M{
			"$unwind": bson.M{
				"path": "$features.oracle.database.databases",
			}},
		bson.M{"$set": bson.M{
			"database": "$features.oracle.database.databases",
		}},
		bson.M{
			"$unwind": bson.M{
				"path": "$database.backups",
			}},
		bson.M{"$project": bson.M{
			"hostname": "$hostname",
			"dbname":   "$database.name",
			"backups":  "$database.backups",
		}},
		bson.M{"$project": bson.M{
			"oracleDatabaseBackup.backupType": "$backups.backupType",
			"oracleDatabaseBackup.hour":       "$backups.hour",
			"oracleDatabaseBackup.weekDays": bson.M{"$map": bson.M{
				"input": "$backups.weekDays",
				"as":    "days",
				"in":    "$$days",
			}},
			"oracleDatabaseBackup.avgBckSize": "$backups.avgBckSize",
			"oracleDatabaseBackup.retention":  "$backups.retention",
			"hostname":                        "$hostname",
			"databasename":                    "$dbname",
		}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(ctx, query)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}
