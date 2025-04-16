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
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "hosts"},
						{Key: "let",
							Value: bson.D{
								{Key: "hosts",
									Value: bson.D{
										{Key: "$ifNull",
											Value: bson.A{
												"$clusterMembershipStatus.veritasClusterHostnames",
												bson.A{},
											},
										},
									},
								},
							},
						},
						{Key: "pipeline",
							Value: bson.A{
								bson.D{
									{Key: "$match",
										Value: bson.D{
											{Key: "$expr",
												Value: bson.D{
													{Key: "$in",
														Value: bson.A{
															"$hostname",
															"$$hosts",
														},
													},
												},
											},
										},
									},
								},
								bson.D{
									{Key: "$group",
										Value: bson.D{
											{Key: "_id", Value: "$hostname"},
											{Key: "hostname", Value: bson.D{{Key: "$first", Value: "$hostname"}}},
											{Key: "isDR", Value: bson.D{{Key: "$first", Value: "$isDR"}}},
											{Key: "cpuCores", Value: bson.D{{Key: "$first", Value: "$info.cpuCores"}}},
										},
									},
								},
								bson.D{
									{Key: "$project",
										Value: bson.D{
											{Key: "_id", Value: 0},
											{Key: "hostname", Value: 1},
											{Key: "isDR", Value: 1},
											{Key: "cpuCores", Value: 1},
										},
									},
								},
							},
						},
						{Key: "as", Value: "existingHosts"},
					},
				},
			},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "id",
							Value: bson.D{
								{Key: "$reduce",
									Value: bson.D{
										{Key: "input",
											Value: bson.D{
												{Key: "$map",
													Value: bson.D{
														{Key: "input", Value: "$clusterMembershipStatus.veritasClusterHostnames"},
														{Key: "as", Value: "host"},
														{Key: "in", Value: "$$host"},
													},
												},
											},
										},
										{Key: "initialValue", Value: ""},
										{Key: "in",
											Value: bson.D{
												{Key: "$cond",
													Value: bson.D{
														{Key: "if",
															Value: bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		"$$value",
																		"",
																	},
																},
															},
														},
														{Key: "then", Value: "$$this"},
														{Key: "else",
															Value: bson.D{
																{Key: "$concat",
																	Value: bson.A{
																		"$$value",
																		"-",
																		"$$this",
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
						},
						{Key: "hostname", Value: 1},
						{Key: "licenseTypeID",
							Value: bson.D{
								{Key: "$arrayElemAt",
									Value: bson.A{
										"$lic._id",
										0,
									},
								},
							},
						},
						{Key: "description",
							Value: bson.D{
								{Key: "$arrayElemAt",
									Value: bson.A{
										"$lic.itemDescription",
										0,
									},
								},
							},
						},
						{Key: "metric",
							Value: bson.D{
								{Key: "$arrayElemAt",
									Value: bson.A{
										"$lic.metric",
										0,
									},
								},
							},
						},
						{Key: "cpuCores", Value: "$info.cpuCores"},
						{Key: "isDR", Value: 1},
						{Key: "clusterHosts", Value: "$existingHosts"},
					},
				},
			},
			bson.D{{Key: "$addFields", Value: bson.D{{Key: "sumCpuCores", Value: bson.D{{Key: "$sum", Value: "$clusterHosts.cpuCores"}}}}}},
			bson.D{
				{Key: "$group",
					Value: bson.D{
						{Key: "id", Value: bson.D{{Key: "$first", Value: "$id"}}},
						{Key: "_id",
							Value: bson.D{
								{Key: "hostname", Value: "$hostname"},
								{Key: "licenseTypeID", Value: "$licenseTypeID"},
							},
						},
						{Key: "clusterHosts", Value: bson.D{{Key: "$first", Value: "$clusterHosts"}}},
						{Key: "description", Value: bson.D{{Key: "$first", Value: "$description"}}},
						{Key: "metric", Value: bson.D{{Key: "$first", Value: "$metric"}}},
						{Key: "isDR", Value: bson.D{{Key: "$max", Value: "$isDR"}}},
						{Key: "cpuCores", Value: bson.D{{Key: "$first", Value: "$sumCpuCores"}}},
					},
				},
			},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "id", Value: 1},
						{Key: "licenseTypeID", Value: "$_id.licenseTypeID"},
						{Key: "hostnames", Value: "$clusterHosts.hostname"},
						{Key: "description", Value: 1},
						{Key: "metric", Value: 1},
						{Key: "idDR", Value: 1},
						{Key: "cpuCores", Value: 1},
						{Key: "count",
							Value: bson.D{
								{Key: "$switch",
									Value: bson.D{
										{Key: "branches",
											Value: bson.A{
												bson.D{
													{Key: "case",
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
												},
												bson.D{
													{Key: "case",
														Value: bson.D{
															{Key: "$eq",
																Value: bson.A{
																	"$metric",
																	"Named User Plus Perpetual",
																},
															},
														},
													},
													{Key: "then",
														Value: bson.D{
															{Key: "$multiply",
																Value: bson.A{
																	bson.D{
																		{Key: "$divide",
																			Value: bson.A{
																				"$cpuCores",
																				2,
																			},
																		},
																	},
																	25,
																},
															},
														},
													},
												},
											},
										},
										{Key: "default",
											Value: bson.D{
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
