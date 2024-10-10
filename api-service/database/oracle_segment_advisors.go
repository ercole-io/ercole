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
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

// SearchOracleDatabaseSegmentAdvisors search segment advisors
func (md *MongoDatabase) SearchOracleDatabaseSegmentAdvisors(keywords []string, sortBy string, sortDesc bool,
	location string, environment string, olderThan time.Time,
) ([]dto.OracleDatabaseSegmentAdvisor, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APProject(bson.M{
				"hostname":    1,
				"environment": 1,
				"location":    1,
				"createdAt":   1,
				"database":    "$features.oracle.database.databases",
			}),
			mu.APSearchFilterStage([]interface{}{"$hostname", "$database.name"}, keywords),
			mu.APProject(bson.M{
				"hostname":                 true,
				"location":                 true,
				"environment":              true,
				"createdAt":                true,
				"database.name":            true,
				"database.segmentAdvisors": true,
				"database.segmentsSize":    true,
			}),
			mu.APUnwind("$database.segmentAdvisors"),
			mu.APProject(bson.M{
				"hostname":       true,
				"location":       true,
				"environment":    true,
				"createdAt":      true,
				"dbname":         "$database.name",
				"reclaimable":    "$database.segmentAdvisors.reclaimable",
				"retrieve":       mu.APODivide("$database.segmentAdvisors.reclaimable", "$database.segmentsSize"),
				"segmentOwner":   "$database.segmentAdvisors.segmentOwner",
				"segmentName":    "$database.segmentAdvisors.segmentName",
				"segmentType":    "$database.segmentAdvisors.segmentType",
				"segmentsSize":   "$database.segmentsSize",
				"partitionName":  "$database.segmentAdvisors.partitionName",
				"recommendation": "$database.segmentAdvisors.recommendation",
			}),
			mu.APOptionalSortingStage(sortBy, sortDesc),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	segmentAdvisors := make([]dto.OracleDatabaseSegmentAdvisor, 0)
	if err := cur.All(context.TODO(), &segmentAdvisors); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return segmentAdvisors, nil
}

func (md *MongoDatabase) SearchOraclePdbSegmentAdvisors(sortBy string, sortDesc bool,
	location string, environment string, olderThan time.Time) ([]dto.OracleDatabaseSegmentAdvisor, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.pdbs"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.pdbs.segmentAdvisors"}},
			bson.M{
				"$project": bson.M{
					"hostname":       1,
					"location":       1,
					"environment":    1,
					"createdAt":      1,
					"dbname":         "$features.oracle.database.databases.name",
					"reclaimable":    "$features.oracle.database.databases.pdbs.segmentAdvisors.reclaimable",
					"segmentOwner":   "$features.oracle.database.databases.pdbs.segmentAdvisors.segmentOwner",
					"segmentName":    "$features.oracle.database.databases.pdbs.segmentAdvisors.segmentName",
					"segmentType":    "$features.oracle.database.databases.pdbs.segmentAdvisors.segmentType",
					"segmentsSize":   "$features.oracle.database.databases.pdbs.segmentsSize",
					"partitionName":  "$features.oracle.database.databases.pdbs.segmentAdvisors.partitionName",
					"recommendation": "$features.oracle.database.databases.pdbs.segmentAdvisors.recommendation",
					"pdbName":        "$features.oracle.database.databases.pdbs.name",
					"retrieve": bson.M{"$cond": bson.A{
						bson.M{"$eq": bson.A{"$features.oracle.database.databases.pdbs.segmentsSize", 0}},
						0.0,
						bson.M{"$divide": bson.A{"$features.oracle.database.databases.pdbs.segmentAdvisors.reclaimable",
							"$features.oracle.database.databases.pdbs.segmentsSize"}}}},
				},
			},
			bson.M{
				"$project": bson.M{
					"hostname":       1,
					"location":       1,
					"environment":    1,
					"createdAt":      1,
					"dbname":         1,
					"reclaimable":    1,
					"segmentOwner":   1,
					"segmentName":    1,
					"segmentType":    1,
					"segmentsSize":   1,
					"partitionName":  1,
					"recommendation": 1,
					"pdbName":        1,
					"retrieve":       bson.M{"$round": bson.A{"$retrieve", 2}}},
			},
			mu.APOptionalSortingStage(sortBy, sortDesc),
		))
	if err != nil {
		return nil, err
	}

	segmentAdvisors := make([]dto.OracleDatabaseSegmentAdvisor, 0)
	if err := cur.All(context.TODO(), &segmentAdvisors); err != nil {
		return nil, err
	}

	return segmentAdvisors, nil
}
