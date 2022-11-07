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

func (md *MongoDatabase) GetOracleOptionList(filter dto.GlobalFilter) ([]dto.OracleDatabaseFeatureUsageStatDto, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabaseFeatureUsageStatDto, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.featureUsageStats"}},
			bson.M{"$project": bson.M{
				"hostname":                               1,
				"location":                               1,
				"environment":                            1,
				"createdAt":                              1,
				"databasename":                           "$features.oracle.database.databases.name",
				"oracleDatabaseFeatureUsageStat.product": "$features.oracle.database.databases.featureUsageStats.product",
				"oracleDatabaseFeatureUsageStat.feature": "$features.oracle.database.databases.featureUsageStats.feature",
				"oracleDatabaseFeatureUsageStat.detectedUsages":   "$features.oracle.database.databases.featureUsageStats.detectedUsages",
				"oracleDatabaseFeatureUsageStat.currentlyUsed":    "$features.oracle.database.databases.featureUsageStats.currentlyUsed",
				"oracleDatabaseFeatureUsageStat.firstUsageDate":   "$features.oracle.database.databases.featureUsageStats.firstUsageDate",
				"oracleDatabaseFeatureUsageStat.lastUsageDate":    "$features.oracle.database.databases.featureUsageStats.lastUsageDate",
				"oracleDatabaseFeatureUsageStat.extraFeatureInfo": "$features.oracle.database.databases.featureUsageStats.extraFeatureInfo",
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
