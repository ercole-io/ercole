// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) FindClusterVeritasLicenses(filter dto.GlobalFilter) ([]dto.ClusterVeritasLicense, error) {
	ctx := context.TODO()

	pipeline := mu.MAPipeline(
		FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
		FilterByOldnessSteps(filter.OlderThan),
		bson.A{
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "archived", Value: false},
						{Key: "clusterMembershipStatus.veritasClusterServer", Value: true},
					},
				},
			},
			bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
			bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.licenses"}}}},
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "oracle_database_license_types"},
						{Key: "localField", Value: "features.oracle.database.databases.licenses.licenseTypeID"},
						{Key: "foreignField", Value: "_id"},
						{Key: "as", Value: "lic"},
					},
				},
			},
			bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$lic"}}}},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "clusterHosts", Value: "$clusterMembershipStatus.veritasClusterHostnames"},
						{Key: "hostname", Value: 1},
						{Key: "licenseTypeID", Value: "$features.oracle.database.databases.licenses.licenseTypeID"},
						{Key: "description", Value: "$lic.itemDescription"},
						{Key: "metric", Value: "$lic.metric"},
						{Key: "cpuCores", Value: "$info.cpuCores"},
					},
				},
			},
			bson.D{
				{Key: "$group",
					Value: bson.D{
						{Key: "_id",
							Value: bson.D{
								{Key: "hostname", Value: "$hostname"},
								{Key: "licenseTypeID", Value: "$licenseTypeID"},
							},
						},
						{Key: "clusterHosts", Value: bson.D{{Key: "$first", Value: "$clusterHosts"}}},
						{Key: "description", Value: bson.D{{Key: "$first", Value: "$description"}}},
						{Key: "metric", Value: bson.D{{Key: "$first", Value: "$metric"}}},
						{Key: "cpuCores", Value: bson.D{{Key: "$max", Value: "$cpuCores"}}},
					},
				},
			},
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "clusterHosts", Value: bson.D{{Key: "$ne", Value: nil}}},
					},
				},
			},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "licenseTypeID", Value: "$_id.licenseTypeID"},
						{Key: "hostnames", Value: "$clusterHosts"},
						{Key: "description", Value: 1},
						{Key: "metric", Value: 1},
						{Key: "count",
							Value: bson.D{
								{Key: "$cond",
									Value: bson.D{
										{Key: "if",
											Value: bson.D{
												{Key: "$eq",
													Value: bson.A{
														"$_id.licenseTypeID",
														"L47837",
													},
												},
											},
										},
										{Key: "then", Value: bson.D{{Key: "$size", Value: "$clusterHosts"}}},
										{Key: "else",
											Value: bson.D{
												{Key: "$multiply",
													Value: bson.A{
														bson.D{{Key: "$size", Value: "$clusterHosts"}},
														bson.D{
															{Key: "$divide",
																Value: bson.A{
																	"$cpuCores",
																	2,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{Key: "cpuCores", Value: 1},
						{Key: "id",
							Value: bson.D{
								{Key: "$reduce",
									Value: bson.D{
										{Key: "input", Value: "$clusterHosts"},
										{Key: "initialValue", Value: ""},
										{Key: "in", Value: bson.D{
											{Key: "$cond", Value: bson.A{
												bson.D{{Key: "$eq", Value: bson.A{"$$value", ""}}},
												"$$this",
												bson.D{{Key: "$concat", Value: bson.A{"$$value", "-", "$$this"}}},
											}},
										}},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ClusterVeritasLicense, 0)

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
