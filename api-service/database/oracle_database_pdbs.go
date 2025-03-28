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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (md *MongoDatabase) FindAllOracleDatabasePdbs(filter dto.GlobalFilter) ([]dto.OracleDatabasePluggableDatabase, error) {
	ctx := context.TODO()

	pipeline := mu.MAPipeline(
		FilterByOldnessSteps(filter.OlderThan),
		FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
		bson.A{
			bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}, {Key: "isDR", Value: false}}}},
			bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
			bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs"}}}},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "hostname", Value: 1},
						{Key: "dbname", Value: "$features.oracle.database.databases.name"},
						{Key: "pdb", Value: "$features.oracle.database.databases.pdbs"},
						{Key: "filteredMigrability",
							Value: bson.D{
								{Key: "$cond",
									Value: bson.D{
										{Key: "if",
											Value: bson.D{
												{Key: "$eq",
													Value: bson.A{
														"$features.oracle.database.databases.pdbs.pgsqlMigrability",
														primitive.Null{},
													},
												},
											},
										},
										{Key: "then", Value: primitive.Null{}},
										{Key: "else",
											Value: bson.D{
												{Key: "$filter",
													Value: bson.D{
														{Key: "input", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability"},
														{Key: "as", Value: "migrability"},
														{Key: "cond",
															Value: bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		"$$migrability.metric",
																		"PLSQL LINES",
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
					},
				},
			},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "hostname", Value: 1},
						{Key: "dbname", Value: "$dbname"},
						{Key: "pdb", Value: "$pdb"},
						{Key: "color",
							Value: bson.D{
								{Key: "$switch",
									Value: bson.D{
										{Key: "branches",
											Value: bson.A{
												bson.D{
													{Key: "case",
														Value: bson.D{
															{Key: "$ne",
																Value: bson.A{
																	"$filteredMigrability",
																	primitive.Null{},
																},
															},
														},
													},
													{Key: "then",
														Value: bson.D{
															{Key: "$switch",
																Value: bson.D{
																	{Key: "branches",
																		Value: bson.A{
																			bson.D{
																				{Key: "case",
																					Value: bson.D{
																						{Key: "$lt",
																							Value: bson.A{
																								bson.D{
																									{Key: "$arrayElemAt",
																										Value: bson.A{
																											"$filteredMigrability.count",
																											0,
																										},
																									},
																								},
																								1000,
																							},
																						},
																					},
																				},
																				{Key: "then", Value: "green"},
																			},
																			bson.D{
																				{Key: "case",
																					Value: bson.D{
																						{Key: "$and",
																							Value: bson.A{
																								bson.D{
																									{Key: "$lte",
																										Value: bson.A{
																											bson.D{
																												{Key: "$arrayElemAt",
																													Value: bson.A{
																														"$filteredMigrability.count",
																														0,
																													},
																												},
																											},
																											10000,
																										},
																									},
																								},
																								bson.D{
																									{Key: "$gte",
																										Value: bson.A{
																											bson.D{
																												{Key: "$arrayElemAt",
																													Value: bson.A{
																														"$filteredMigrability.count",
																														0,
																													},
																												},
																											},
																											1000,
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																				{Key: "then", Value: "yellow"},
																			},
																			bson.D{
																				{Key: "case",
																					Value: bson.D{
																						{Key: "$gt",
																							Value: bson.A{
																								bson.D{
																									{Key: "$arrayElemAt",
																										Value: bson.A{
																											"$filteredMigrability.count",
																											0,
																										},
																									},
																								},
																								10000,
																							},
																						},
																					},
																				},
																				{Key: "then", Value: "red"},
																			},
																		},
																	},
																	{Key: "default", Value: ""},
																},
															},
														},
													},
												},
											},
										},
										{Key: "default", Value: ""},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(ctx, pipeline)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	out := make([]dto.OracleDatabasePluggableDatabase, 0)
	if err = cur.All(ctx, &out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}

func (md *MongoDatabase) PdbExist(hostname, dbname, pdbname string) (bool, error) {
	filter := bson.D{
		{Key: "archived", Value: false},
		{Key: "isDR", Value: false},
		{Key: "hostname", Value: hostname},
		{Key: "features.oracle.database.databases.name", Value: dbname},
		{Key: "features.oracle.database.databases.pdbs.name", Value: pdbname},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return cur > 0, nil
}
