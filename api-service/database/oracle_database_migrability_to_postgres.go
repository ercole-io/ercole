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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (md *MongoDatabase) FindPsqlMigrabilities(hostname, dbname string) ([]model.PgsqlMigrability, error) {
	ctx := context.TODO()

	result := make([]model.PgsqlMigrability, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx,
			bson.A{
				bson.D{
					{Key: "$match",
						Value: bson.D{
							{Key: "archived", Value: false},
							{Key: "hostname", Value: hostname},
							{Key: "isDR", Value: false},
						},
					},
				},
				bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
				bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.name", Value: dbname}}}},
				bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pgsqlMigrability"}}}},
				bson.D{
					{Key: "$project",
						Value: bson.D{
							{Key: "metric", Value: "$features.oracle.database.databases.pgsqlMigrability.metric"},
							{Key: "count", Value: "$features.oracle.database.databases.pgsqlMigrability.count"},
							{Key: "schema", Value: "$features.oracle.database.databases.pgsqlMigrability.schema"},
							{Key: "objectType", Value: "$features.oracle.database.databases.pgsqlMigrability.objectType"},
						},
					},
				},
			})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) FindPdbPsqlMigrabilities(hostname, dbname, pdbname string) ([]model.PgsqlMigrability, error) {
	ctx := context.TODO()

	result := make([]model.PgsqlMigrability, 0)

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "archived", Value: false},
					{Key: "hostname", Value: hostname},
					{Key: "isDR", Value: false},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.name", Value: dbname}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.pdbs.name", Value: pdbname}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "metric", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.metric"},
					{Key: "count", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.count"},
					{Key: "schema", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.schema"},
					{Key: "objectType", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability.objectType"},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) ListOracleDatabasePsqlMigrabilities() ([]dto.OracleDatabasePgsqlMigrability, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabasePgsqlMigrability, 0)

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}, {Key: "isDR", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.pgsqlMigrability", Value: bson.D{{Key: "$ne", Value: primitive.Null{}}}}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "hostname", Value: 1},
					{Key: "dbname", Value: "$features.oracle.database.databases.name"},
					{Key: "metrics", Value: "$features.oracle.database.databases.pgsqlMigrability"},
				},
			},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "flag",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.A{
									bson.D{
										{Key: "$gt",
											Value: bson.A{
												bson.D{
													{Key: "$size",
														Value: bson.D{
															{Key: "$filter",
																Value: bson.D{
																	{Key: "input", Value: "$metrics"},
																	{Key: "as", Value: "item"},
																	{Key: "cond",
																		Value: bson.D{
																			{Key: "$and",
																				Value: bson.A{
																					bson.D{
																						{Key: "$eq",
																							Value: bson.A{
																								"$$item.metric",
																								"PLSQL LINES",
																							},
																						},
																					},
																					bson.D{
																						{Key: "$lt",
																							Value: bson.A{
																								"$$item.count",
																								1000,
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
												0,
											},
										},
									},
									"green",
									bson.D{
										{Key: "$cond",
											Value: bson.A{
												bson.D{
													{Key: "$gt",
														Value: bson.A{
															bson.D{
																{Key: "$size",
																	Value: bson.D{
																		{Key: "$filter",
																			Value: bson.D{
																				{Key: "input", Value: "$metrics"},
																				{Key: "as", Value: "item"},
																				{Key: "cond",
																					Value: bson.D{
																						{Key: "$and",
																							Value: bson.A{
																								bson.D{
																									{Key: "$eq",
																										Value: bson.A{
																											"$$item.metric",
																											"PLSQL LINES",
																										},
																									},
																								},
																								bson.D{
																									{Key: "$gte",
																										Value: bson.A{
																											"$$item.count",
																											1000,
																										},
																									},
																								},
																								bson.D{
																									{Key: "$lte",
																										Value: bson.A{
																											"$$item.count",
																											10000,
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
															0,
														},
													},
												},
												"yellow",
												"red",
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
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection(hostCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) ListOracleDatabasePdbPsqlMigrabilities() ([]dto.OracleDatabasePdbPgsqlMigrability, error) {
	ctx := context.TODO()

	result := make([]dto.OracleDatabasePdbPgsqlMigrability, 0)

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}, {Key: "isDR", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.pdbs"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.pdbs.pgsqlMigrability", Value: bson.D{{Key: "$ne", Value: primitive.Null{}}}}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "hostname", Value: 1},
					{Key: "dbname", Value: "$features.oracle.database.databases.name"},
					{Key: "pdbname", Value: "$features.oracle.database.databases.pdbs.name"},
					{Key: "metrics", Value: "$features.oracle.database.databases.pdbs.pgsqlMigrability"},
				},
			},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "flag",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.A{
									bson.D{
										{Key: "$gt",
											Value: bson.A{
												bson.D{
													{Key: "$size",
														Value: bson.D{
															{Key: "$filter",
																Value: bson.D{
																	{Key: "input", Value: "$metrics"},
																	{Key: "as", Value: "item"},
																	{Key: "cond",
																		Value: bson.D{
																			{Key: "$and",
																				Value: bson.A{
																					bson.D{
																						{Key: "$eq",
																							Value: bson.A{
																								"$$item.metric",
																								"PLSQL LINES",
																							},
																						},
																					},
																					bson.D{
																						{Key: "$lt",
																							Value: bson.A{
																								"$$item.count",
																								1000,
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
												0,
											},
										},
									},
									"green",
									bson.D{
										{Key: "$cond",
											Value: bson.A{
												bson.D{
													{Key: "$gt",
														Value: bson.A{
															bson.D{
																{Key: "$size",
																	Value: bson.D{
																		{Key: "$filter",
																			Value: bson.D{
																				{Key: "input", Value: "$metrics"},
																				{Key: "as", Value: "item"},
																				{Key: "cond",
																					Value: bson.D{
																						{Key: "$and",
																							Value: bson.A{
																								bson.D{
																									{Key: "$eq",
																										Value: bson.A{
																											"$$item.metric",
																											"PLSQL LINES",
																										},
																									},
																								},
																								bson.D{
																									{Key: "$gte",
																										Value: bson.A{
																											"$$item.count",
																											1000,
																										},
																									},
																								},
																								bson.D{
																									{Key: "$lte",
																										Value: bson.A{
																											"$$item.count",
																											10000,
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
															0,
														},
													},
												},
												"yellow",
												"red",
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
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection(hostCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
