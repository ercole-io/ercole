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
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
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
