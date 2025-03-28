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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (md *MongoDatabase) FindOracleDatabasePoliciesAudit(hostname, dbname string) (*dto.OraclePoliciesAudit, error) {
	ctx := context.TODO()

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "archived", Value: false},
					{Key: "isDR", Value: false},
					{Key: "hostname", Value: hostname},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "features.oracle.database.databases.name", Value: dbname}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.policiesAudit"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "policiesAudit", Value: "$features.oracle.database.databases.policiesAudit"},
				},
			},
		},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$policiesAudit"}}}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)

	for cur.Next(ctx) {
		var item map[string]string
		if cur.Decode(&item) != nil {
			return nil, err
		}

		list = append(list, item["_id"])
	}

	return &dto.OraclePoliciesAudit{List: list}, nil
}

func (md *MongoDatabase) ListOracleDatabasePoliciesAudit() ([]dto.OraclePoliciesAuditListResponse, error) {
	ctx := context.TODO()

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}, {Key: "isDR", Value: false}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "config"},
					{Key: "pipeline",
						Value: bson.A{
							bson.D{
								{Key: "$project",
									Value: bson.D{
										{Key: "_id", Value: 0},
										{Key: "policiesAuditConfig", Value: "$apiservice.oracledatabasepoliciesaudit"},
									},
								},
							},
						},
					},
					{Key: "as", Value: "policiesAuditConfig"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: "$policiesAuditConfig"}},
		bson.D{{Key: "$unwind", Value: "$policiesAuditConfig.policiesAuditConfig"}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$_id"},
					{Key: "document", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
					{Key: "policiesAuditConfig", Value: bson.D{{Key: "$push", Value: "$policiesAuditConfig.policiesAuditConfig"}}},
				},
			},
		},
		bson.D{
			{Key: "$replaceRoot",
				Value: bson.D{
					{Key: "newRoot",
						Value: bson.D{
							{Key: "$mergeObjects",
								Value: bson.A{
									"$document",
									bson.D{{Key: "policiesAuditConfig", Value: "$policiesAuditConfig"}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "policiesaudit",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.A{
									bson.D{
										{Key: "$and",
											Value: bson.A{
												bson.D{
													{Key: "$ne",
														Value: bson.A{
															bson.D{{Key: "$type", Value: "$features.oracle.database.databases.policiesAudit"}},
															"missing",
														},
													},
												},
												bson.D{
													{Key: "$ne",
														Value: bson.A{
															bson.D{{Key: "$type", Value: "$features.oracle.database.databases.policiesAudit"}},
															"null",
														},
													},
												},
												bson.D{
													{Key: "$gt",
														Value: bson.A{
															bson.D{{Key: "$size", Value: "$features.oracle.database.databases.policiesAudit"}},
															0,
														},
													},
												},
											},
										},
									},
									"$features.oracle.database.databases.policiesAudit",
									bson.A{},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "policiesAuditConfigExistsAndNotEmpty",
						Value: bson.D{
							{Key: "$and",
								Value: bson.A{
									bson.D{
										{Key: "$ne",
											Value: bson.A{
												"$policiesAuditConfig",
												nil,
											},
										},
									},
									bson.D{
										{Key: "$gt",
											Value: bson.A{
												bson.D{{Key: "$size", Value: "$policiesAuditConfig"}},
												0,
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
			{Key: "$addFields",
				Value: bson.D{
					{Key: "matched",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.D{
									{Key: "if", Value: "$policiesAuditConfigExistsAndNotEmpty"},
									{Key: "then",
										Value: bson.D{
											{Key: "$filter",
												Value: bson.D{
													{Key: "input", Value: "$policiesAuditConfig"},
													{Key: "as", Value: "item"},
													{Key: "cond",
														Value: bson.D{
															{Key: "$in",
																Value: bson.A{
																	"$$item",
																	"$policiesaudit",
																},
															},
														},
													},
												},
											},
										},
									},
									{Key: "else", Value: bson.A{}},
								},
							},
						},
					},
					{Key: "notmatched",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.D{
									{Key: "if", Value: "$policiesAuditConfigExistsAndNotEmpty"},
									{Key: "then",
										Value: bson.D{
											{Key: "$filter",
												Value: bson.D{
													{Key: "input", Value: "$policiesAuditConfig"},
													{Key: "as", Value: "item"},
													{Key: "cond",
														Value: bson.D{
															{Key: "$not",
																Value: bson.D{
																	{Key: "$in",
																		Value: bson.A{
																			"$$item",
																			"$policiesaudit",
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
									{Key: "else", Value: bson.A{}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "flag",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.D{
									{Key: "if",
										Value: bson.D{
											{Key: "$or",
												Value: bson.A{
													bson.D{
														{Key: "$eq",
															Value: bson.A{
																"$policiesaudit",
																primitive.Null{},
															},
														},
													},
													bson.D{
														{Key: "$eq",
															Value: bson.A{
																"$policiesaudit",
																bson.A{},
															},
														},
													},
													bson.D{{Key: "$not", Value: "$policiesAuditConfigExistsAndNotEmpty"}},
												},
											},
										},
									},
									{Key: "then", Value: "N/A"},
									{Key: "else",
										Value: bson.D{
											{Key: "$cond",
												Value: bson.D{
													{Key: "if",
														Value: bson.D{
															{Key: "$setIsSubset",
																Value: bson.A{
																	"$policiesAuditConfig",
																	"$policiesaudit",
																},
															},
														},
													},
													{Key: "then", Value: "green"},
													{Key: "else", Value: "red"},
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
					{Key: "dbName", Value: "$features.oracle.database.databases.name"},
					{Key: "policiesAuditConfigured", Value: "$policiesAuditConfig"},
					{Key: "policiesAudit", Value: "$policiesaudit"},
					{Key: "matched", Value: "$matched"},
					{Key: "notmatched", Value: "$notmatched"},
					{Key: "flag", Value: 1},
				},
			},
		},
	}

	result := make([]dto.OraclePoliciesAuditListResponse, 0)

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
