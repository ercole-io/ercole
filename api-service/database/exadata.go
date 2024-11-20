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
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const exadataCollection = "exadatas"

var getExadataInstancePipeline = bson.A{
	bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "exadata_vm_clusternames"},
				{Key: "localField", Value: "components.vms.name"},
				{Key: "foreignField", Value: "vmname"},
				{Key: "as", Value: "matchedDocument"},
			},
		},
	},
	bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "components",
					Value: bson.D{
						{Key: "$map",
							Value: bson.D{
								{Key: "input", Value: "$components"},
								{Key: "as", Value: "component"},
								{Key: "in",
									Value: bson.D{
										{Key: "$mergeObjects",
											Value: bson.A{
												"$$component",
												bson.D{
													{Key: "vms",
														Value: bson.D{
															{Key: "$map",
																Value: bson.D{
																	{Key: "input", Value: "$$component.vms"},
																	{Key: "as", Value: "vm"},
																	{Key: "in",
																		Value: bson.D{
																			{Key: "$mergeObjects",
																				Value: bson.A{
																					"$$vm",
																					bson.D{
																						{Key: "clusterName",
																							Value: bson.D{
																								{Key: "$let",
																									Value: bson.D{
																										{Key: "vars",
																											Value: bson.D{
																												{Key: "matchedDocument",
																													Value: bson.D{
																														{Key: "$filter",
																															Value: bson.D{
																																{Key: "input", Value: "$matchedDocument"},
																																{Key: "as", Value: "match"},
																																{Key: "cond",
																																	Value: bson.D{
																																		{Key: "$and",
																																			Value: bson.A{
																																				bson.D{
																																					{Key: "$eq",
																																						Value: bson.A{
																																							"$$match.instancerackid",
																																							"$rackID",
																																						},
																																					},
																																				},
																																				bson.D{
																																					{Key: "$eq",
																																						Value: bson.A{
																																							"$$match.vmname",
																																							"$$vm.name",
																																						},
																																					},
																																				},
																																				bson.D{
																																					{Key: "$eq",
																																						Value: bson.A{
																																							"$$match.hostid",
																																							"$$component.hostID",
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
																										{Key: "in",
																											Value: bson.D{
																												{Key: "$cond",
																													Value: bson.D{
																														{Key: "if",
																															Value: bson.D{
																																{Key: "$gt",
																																	Value: bson.A{
																																		bson.D{{Key: "$size", Value: "$$matchedDocument"}},
																																		0,
																																	},
																																},
																															},
																														},
																														{Key: "then",
																															Value: bson.D{
																																{Key: "$arrayElemAt",
																																	Value: bson.A{
																																		"$$matchedDocument.clustername",
																																		0,
																																	},
																																},
																															},
																														},
																														{Key: "else", Value: ""},
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
	bson.D{{Key: "$unset", Value: "matchedDocument"}},
}

func (md *MongoDatabase) ListExadataInstances(f dto.GlobalFilter, hidden bool) ([]dto.ExadataInstanceResponse, error) {
	ctx := context.TODO()

	result := make([]dto.ExadataInstanceResponse, 0)

	projection := bson.D{{Key: "rackID", Value: 1}, {Key: "hostname", Value: 1}}

	opts := options.Find().SetProjection(projection)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		Find(ctx, FilterExadata(f, hidden), opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) FindExadataInstance(rackID string, hidden bool) (*model.OracleExadataInstance, error) {
	ctx := context.TODO()

	pipeline := append(getExadataInstancePipeline, bson.D{{
		Key: "$match",
		Value: bson.D{
			{Key: "rackID", Value: rackID},
			{Key: "$or",
				Value: bson.A{
					bson.D{{Key: "hidden", Value: bson.D{{Key: "$exists", Value: hidden}}}},
					bson.D{{Key: "hidden", Value: hidden}},
				},
			},
		}}})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]model.OracleExadataInstance, 0)

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	if len(res) > 0 {
		return &res[0], nil
	}

	return nil, mongo.ErrNoDocuments
}

func (md *MongoDatabase) UpdateExadataInstance(instance model.OracleExadataInstance) error {
	filter := bson.M{"rackID": instance.RackID}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		ReplaceOne(context.TODO(), filter, instance)
	if err != nil {
		return err
	}

	return md.updateExadataTime(instance.RackID)
}

func (md *MongoDatabase) FindAllExadataInstances(hidden bool) ([]model.OracleExadataInstance, error) {
	ctx := context.TODO()

	result := make([]model.OracleExadataInstance, 0)

	condition := bson.D{
		{Key: "$or",
			Value: bson.A{
				bson.D{{Key: "hidden", Value: bson.D{{Key: "$exists", Value: hidden}}}},
				bson.D{{Key: "hidden", Value: hidden}},
			},
		},
	}

	if hidden {
		condition = bson.D{{Key: "hidden", Value: hidden}}
	}

	pipeline := append(getExadataInstancePipeline,
		bson.D{{Key: "$match", Value: condition}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "hostname", Value: 1}}}})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) FindExadataClusterViews() ([]dto.OracleExadataClusterView, error) {
	ctx := context.TODO()

	cvPipeline := bson.A{
		bson.D{{Key: "$unwind", Value: "$components"}},
		bson.D{{Key: "$unwind", Value: "$components.vms"}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "components.vms.clusterName", Value: bson.D{{Key: "$ne", Value: ""}}}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$components.vms.clusterName"},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
					{Key: "virtualNode",
						Value: bson.D{
							{Key: "$addToSet",
								Value: bson.D{
									{Key: "hostname", Value: "$hostname"},
									{Key: "rackID", Value: "$rackID"},
									{Key: "hostType", Value: "$components.hostType"},
									{Key: "clustername", Value: "$components.vms.clusterName"},
									{Key: "vmname", Value: "$components.vms.name"},
									{Key: "totalRAM",
										Value: bson.D{
											{Key: "$sum",
												Value: bson.D{
													{Key: "$cond",
														Value: bson.A{
															bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		"$components.hostType",
																		model.DOM0,
																	},
																},
															},
															"$components.vms.ramOnline",
															"$components.vms.ramCurrent",
														},
													},
												},
											},
										},
									},
									{Key: "totalCPU",
										Value: bson.D{
											{Key: "$sum",
												Value: bson.D{
													{Key: "$cond",
														Value: bson.A{
															bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		"$components.hostType",
																		model.DOM0,
																	},
																},
															},
															"$components.vms.cpuOnline",
															"$components.vms.cpuCurrent",
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
		bson.D{{Key: "$match", Value: bson.D{{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}}}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$$ROOT"}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$virtualNode.clustername"},
					{Key: "virtualNode",
						Value: bson.D{
							{Key: "$addToSet",
								Value: bson.D{
									{Key: "Hostname", Value: bson.D{{Key: "$first", Value: "$virtualNode.hostname"}}},
									{Key: "RackID", Value: bson.D{{Key: "$first", Value: "$virtualNode.rackID"}}},
									{Key: "HostType", Value: bson.D{{Key: "$first", Value: "$virtualNode.hostType"}}},
									{Key: "Clustername", Value: bson.D{{Key: "$first", Value: "$virtualNode.clustername"}}},
									{Key: "VmNames", Value: "$virtualNode.vmname"},
									{Key: "HostnameVm2", Value: bson.D{{Key: "$last", Value: "$virtualNode.vmname"}}},
									{Key: "TotalRAM", Value: bson.D{{Key: "$sum", Value: "$virtualNode.totalRAM"}}},
									{Key: "TotalCPU", Value: bson.D{{Key: "$sum", Value: "$virtualNode.totalCPU"}}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: "$virtualNode"}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "virtualNode", Value: 1},
				},
			},
		},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$virtualNode"}}}},
	}

	pipeline := append(getExadataInstancePipeline, cvPipeline...)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]dto.OracleExadataClusterView, 0)

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (md *MongoDatabase) updateExadataTime(rackID string) error {
	now := md.TimeNow()

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		UpdateOne(context.TODO(), bson.M{"rackID": rackID},
			bson.M{"$set": bson.M{"updateAt": now}})
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) FindExadataPatchAdvisorsByRackID(rackID string) ([]dto.OracleExadataPatchAdvisor, error) {
	ctx := context.TODO()

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "hidden", Value: false},
					{Key: "rackID", Value: rackID},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$components"}}}},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "matchesPattern",
						Value: bson.D{
							{Key: "$regexFind",
								Value: bson.D{
									{Key: "input", Value: "$components.imageVersion"},
									{Key: "regex", Value: primitive.Regex{Pattern: `\b\d{6}\b`, Options: ""}},
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
					{Key: "extractedDate", Value: "$matchesPattern.match"},
				},
			},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "releaseDate",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.A{
									"$matchesPattern",
									bson.D{
										{Key: "$dateFromString",
											Value: bson.D{
												{Key: "dateString",
													Value: bson.D{
														{Key: "$concat",
															Value: bson.A{
																"20",
																"$extractedDate",
															},
														},
													},
												},
												{Key: "format", Value: "%Y%m%d"},
											},
										},
									},
									primitive.Null{},
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
					{Key: "rackID", Value: "$components.rackID"},
					{Key: "hostname", Value: "$components.hostname"},
					{Key: "imageVersion", Value: "$components.imageVersion"},
					{Key: "releaseDate", Value: "$releaseDate"},
					{Key: "fourMonths", Value: bson.D{
						{Key: "$gte",
							Value: bson.A{
								"$releaseDate",
								time.Now().AddDate(0, -4, 0),
							},
						},
					}},
					{Key: "sixMonths", Value: bson.D{
						{Key: "$gte",
							Value: bson.A{
								"$releaseDate",
								time.Now().AddDate(0, -6, 0),
							},
						},
					}},
					{Key: "twelveMonths", Value: bson.D{
						{Key: "$gte",
							Value: bson.A{
								"$releaseDate",
								time.Now().AddDate(0, -12, 0),
							},
						},
					}},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]dto.OracleExadataPatchAdvisor, 0)

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (md *MongoDatabase) FindAllExadataPatchAdvisors() ([]dto.OracleExadataPatchAdvisor, error) {
	ctx := context.TODO()

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "hidden", Value: false},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$components"}}}},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "matchesPattern",
						Value: bson.D{
							{Key: "$regexFind",
								Value: bson.D{
									{Key: "input", Value: "$components.imageVersion"},
									{Key: "regex", Value: primitive.Regex{Pattern: `\b\d{6}\b`, Options: ""}},
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
					{Key: "extractedDate", Value: "$matchesPattern.match"},
				},
			},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "releaseDate",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.A{
									"$matchesPattern",
									bson.D{
										{Key: "$dateFromString",
											Value: bson.D{
												{Key: "dateString",
													Value: bson.D{
														{Key: "$concat",
															Value: bson.A{
																"20",
																"$extractedDate",
															},
														},
													},
												},
												{Key: "format", Value: "%Y%m%d"},
											},
										},
									},
									primitive.Null{},
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
					{Key: "rackID", Value: "$components.rackID"},
					{Key: "hostname", Value: "$components.hostname"},
					{Key: "imageVersion", Value: "$components.imageVersion"},
					{Key: "releaseDate", Value: "$releaseDate"},
					{Key: "fourMonths", Value: bson.D{
						{Key: "$gte",
							Value: bson.A{
								"$releaseDate",
								time.Now().AddDate(0, -4, 0),
							},
						},
					}},
					{Key: "sixMonths", Value: bson.D{
						{Key: "$gte",
							Value: bson.A{
								"$releaseDate",
								time.Now().AddDate(0, -6, 0),
							},
						},
					}},
					{Key: "twelveMonths", Value: bson.D{
						{Key: "$gte",
							Value: bson.A{
								"$releaseDate",
								time.Now().AddDate(0, -12, 0),
							},
						},
					}},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]dto.OracleExadataPatchAdvisor, 0)

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}
