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
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

const hostCollection = "hosts"

func (md *MongoDatabase) SearchHosts(mode string, filters dto.SearchHostsFilters) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, 0)
	if err := md.getHosts(mode, filters, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (md *MongoDatabase) GetHostDataSummaries(filters dto.SearchHostsFilters) ([]dto.HostDataSummary, error) {
	filters.PageNumber, filters.PageSize = -1, -1
	out := make([]dto.HostDataSummary, 0)

	if err := md.getHosts("summary", filters, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// out must be a pointer to a slice
func (md *MongoDatabase) getHosts(mode string, filters dto.SearchHostsFilters, out interface{}) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByLocationAndEnvironmentSteps(filters.Location, filters.Environment),
			FilterByOldnessSteps(filters.OlderThan),
			mu.MAPipeline(
				mu.APOptionalStage(filters.Hostname != "", mu.APMatch(bson.M{
					"hostname": primitive.Regex{Pattern: regexp.QuoteMeta(filters.Hostname), Options: "i"},
				})),
				mu.APOptionalStage(filters.Database != "", mu.APMatch(bson.M{
					"features.oracle.database.databases.Name": primitive.Regex{Pattern: regexp.QuoteMeta(filters.Database), Options: "i"},
				})),
				mu.APOptionalStage(filters.HardwareAbstractionTechnology != "", mu.APMatch(bson.M{
					"info.hardwareAbstractionTechnology": primitive.Regex{Pattern: regexp.QuoteMeta(filters.HardwareAbstractionTechnology), Options: "i"},
				})),
				mu.APOptionalStage(filters.OperatingSystem != "", mu.APMatch(bson.M{
					"info.os": primitive.Regex{Pattern: regexp.QuoteMeta(filters.OperatingSystem), Options: "i"},
				})),
				mu.APOptionalStage(filters.Kernel != "", mu.APMatch(bson.M{
					"info.kernel": primitive.Regex{Pattern: regexp.QuoteMeta(filters.Kernel), Options: "i"},
				})),
				mu.APOptionalStage(filters.LTEMemoryTotal != -1, mu.APMatch(bson.M{
					"info.memoryTotal": mu.QOLessThanOrEqual(filters.LTEMemoryTotal),
				})),
				mu.APOptionalStage(filters.GTEMemoryTotal != -1, mu.APMatch(bson.M{
					"info.memoryTotal": bson.M{
						"$gte": filters.GTEMemoryTotal,
					},
				})),
				mu.APOptionalStage(filters.LTESwapTotal != -1, mu.APMatch(bson.M{
					"info.swapTotal": mu.QOLessThanOrEqual(filters.LTESwapTotal),
				})),
				mu.APOptionalStage(filters.GTESwapTotal != -1, mu.APMatch(bson.M{
					"info.swapTotal": bson.M{
						"$gte": filters.GTESwapTotal,
					},
				})),
				getIsMemberOfClusterFilterStep(filters.IsMemberOfCluster),
				mu.APOptionalStage(filters.CPUModel != "", mu.APMatch(bson.M{
					"info.cpuModel": primitive.Regex{Pattern: regexp.QuoteMeta(filters.CPUModel), Options: "i"},
				})),
				mu.APOptionalStage(filters.LTECPUCores != -1, mu.APMatch(bson.M{
					"info.cpuCores": mu.QOLessThanOrEqual(filters.LTECPUCores),
				})),
				mu.APOptionalStage(filters.GTECPUCores != -1, mu.APMatch(bson.M{
					"info.cpuCores": bson.M{
						"$gte": filters.GTECPUCores,
					},
				})),
				mu.APOptionalStage(filters.LTECPUThreads != -1, mu.APMatch(bson.M{
					"info.cpuThreads": mu.QOLessThanOrEqual(filters.LTECPUThreads),
				})),
				mu.APOptionalStage(filters.GTECPUThreads != -1, mu.APMatch(bson.M{
					"info.cpuThreads": bson.M{
						"$gte": filters.GTECPUThreads,
					},
				})),
			),
			mu.APSearchFilterStage([]interface{}{
				"$hostname",
				"$features.oracle.database.databases.name",
				"$features.oracle.database.databases.uniqueName",
				"$clusters.name",
			}, filters.Search),
			AddAssociatedClusterNameAndVirtualizationNode(filters.OlderThan),
			getClusterFilterStep(filters.Cluster),
			mu.APOptionalStage(filters.VirtualizationNode != "", mu.APMatch(bson.M{
				"virtualizationNode": primitive.Regex{Pattern: regexp.QuoteMeta(filters.VirtualizationNode), Options: "i"},
			})),
			mu.APOptionalStage(mode == "mongo" || mode == "hostnames", mu.APUnset("cluster", "virtualizationNode")),
			mu.APOptionalStage(mode == "hostnames", mu.MAPipeline(
				mu.APProject(bson.M{
					"_id":      0,
					"hostname": 1,
				}),
				mu.APOptionalSortingStage(filters.SortBy, filters.SortDesc),
			)),
			mu.APOptionalStage(mode != "mongo" && mode != "hostnames", mu.MAPipeline(
				mu.APOptionalStage(mode == "lms", mu.APMatch(
					mu.QOExpr(mu.APOGreater(mu.APOSize(mu.APOIfNull("$features.oracle.database.databases", bson.A{})), 0))),
				),
				mu.APOptionalStage(mode == "summary", mu.APProject(bson.M{
					"_id":                     true,
					"createdAt":               true,
					"hostname":                true,
					"location":                true,
					"environment":             true,
					"agentVersion":            true,
					"info":                    true,
					"consumptions":            true,
					"clusterMembershipStatus": true,
					"virtualizationNode":      true,
					"cluster":                 true,
					"databases": bson.M{
						model.TechnologyOracleDatabase:       "$features.oracle.database.databases.name",
						model.TechnologyMicrosoftSQLServer:   "$features.microsoft.sqlServer.instances.name",
						model.TechnologyOracleMySQL:          "$features.mysql.instances.name",
						model.TechnologyPostgreSQLPostgreSQL: "$features.postgresql.instances.name",
						model.TechnologyMongoDBMongoDB:       "$features.mongodb.instances.name",
					},
					"missingDatabases": "$features.oracle.database.missingDatabases",
					"technology": bson.D{
						{Key: "$switch",
							Value: bson.D{
								{Key: "branches",
									Value: bson.A{
										bson.D{
											{Key: "case",
												Value: bson.D{
													{Key: "$or",
														Value: bson.A{
															bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		bson.D{{Key: "$type", Value: "$features.oracle"}},
																		"object",
																	},
																},
															},
															bson.D{
																{Key: "$and",
																	Value: bson.A{
																		bson.D{{Key: "$isArray", Value: "$features.oracle.database.databases"}},
																		bson.D{
																			{Key: "$gt",
																				Value: bson.A{
																					bson.D{{Key: "$size", Value: "$features.oracle.database.databases"}},
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
											{Key: "then", Value: model.TechnologyOracleDatabase},
										},
										bson.D{
											{Key: "case",
												Value: bson.D{
													{Key: "$or",
														Value: bson.A{
															bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		bson.D{{Key: "$type", Value: "$features.microsoft"}},
																		"object",
																	},
																},
															},
															bson.D{
																{Key: "$and",
																	Value: bson.A{
																		bson.D{{Key: "$isArray", Value: "$features.microsoft.sqlServer.instances"}},
																		bson.D{
																			{Key: "$gt",
																				Value: bson.A{
																					bson.D{{Key: "$size", Value: "$features.microsoft.sqlServer.instances"}},
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
											{Key: "then", Value: model.TechnologyMicrosoftSQLServer},
										},
										bson.D{
											{Key: "case",
												Value: bson.D{
													{Key: "$or",
														Value: bson.A{
															bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		bson.D{{Key: "$type", Value: "$features.mysql"}},
																		"object",
																	},
																},
															},
															bson.D{
																{Key: "$and",
																	Value: bson.A{
																		bson.D{{Key: "$isArray", Value: "$features.mysql.instances"}},
																		bson.D{
																			{Key: "$gt",
																				Value: bson.A{
																					bson.D{{Key: "$size", Value: "$features.mysql.instances"}},
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
											{Key: "then", Value: model.TechnologyOracleMySQL},
										},
										bson.D{
											{Key: "case",
												Value: bson.D{
													{Key: "$or",
														Value: bson.A{
															bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		bson.D{{Key: "$type", Value: "$features.postgresql"}},
																		"object",
																	},
																},
															},
															bson.D{
																{Key: "$and",
																	Value: bson.A{
																		bson.D{{Key: "$isArray", Value: "$features.postgresql.instances"}},
																		bson.D{
																			{Key: "$gt",
																				Value: bson.A{
																					bson.D{{Key: "$size", Value: "$features.postgresql.instances"}},
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
											{Key: "then", Value: model.TechnologyPostgreSQLPostgreSQL},
										},
										bson.D{
											{Key: "case",
												Value: bson.D{
													{Key: "$or",
														Value: bson.A{
															bson.D{
																{Key: "$eq",
																	Value: bson.A{
																		bson.D{{Key: "$type", Value: "$features.mongodb"}},
																		"object",
																	},
																},
															},
															bson.D{
																{Key: "$and",
																	Value: bson.A{
																		bson.D{{Key: "$isArray", Value: "$features.mongodb.instances"}},
																		bson.D{
																			{Key: "$gt",
																				Value: bson.A{
																					bson.D{{Key: "$size", Value: "$features.mongodb.instances"}},
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
											{Key: "then", Value: model.TechnologyMongoDBMongoDB},
										},
									},
								},
								{Key: "default", Value: primitive.Null{}},
							},
						},
					},
				})),
				mu.APOptionalStage(mode == "lms", mu.MAPipeline(
					mu.APMatch(mu.QOExpr(mu.APOGreater(mu.APOSize("$features.oracle.database.databases"), 0))),
					mu.APUnwind("$features.oracle.database.databases"),
					mu.APSet(bson.M{
						"database": "$features.oracle.database.databases",
					}),
					mu.APUnset("features"),
					mu.APSet(bson.M{
						"isVirtualServer": mu.APOEqual("$info.hardwareAbstraction", model.HardwareAbstractionVirtual),
						"database.pdbs": mu.APOCond("$database.isCDB", bson.M{
							"$concatArrays": bson.A{
								bson.A{""},
								mu.APOMap("$database.pdbs", "pdb", "$$pdb.name"),
							},
						}, bson.A{""}),
					}),
					mu.APUnwind("$database.pdbs"),
					mu.APProject(bson.M{
						// "Database":           1,
						"createdAt":          "$createdAt",
						"dismissedAt":        "$dismissedAt",
						"physicalServerName": mu.APOCond("$isVirtualServer", mu.APOIfNull("$cluster", ""), "$hostname"),
						"virtualServerName":  mu.APOCond("$isVirtualServer", "$hostname", mu.APOIfNull("$cluster", "")),
						"virtualizationTechnology": bson.M{
							"$switch": bson.M{
								"branches": bson.A{
									bson.M{"case": mu.APOEqual("$info.hardwareAbstractionTechnology", model.HardwareAbstractionTechnologyPhysical), "then": ""},
									bson.M{
										"case": mu.APOEqual("$info.hardwareAbstractionTechnology", model.HardwareAbstractionTechnologyOvm),
										"then": mu.APOCond(bson.M{
											"$regexMatch": bson.M{
												"input": "$info.cpuModel",
												"regex": primitive.Regex{
													Options: "i",
													Pattern: "sparc",
												},
											},
										}, "OVM Server for SPARC", "OVM Server for x86"),
									},
									bson.M{"case": mu.APOEqual("$info.hardwareAbstractionTechnology", model.HardwareAbstractionTechnologyVmware), "then": "VMware"},
									bson.M{"case": mu.APOEqual("$info.hardwareAbstractionTechnology", model.HardwareAbstractionTechnologyHyperv), "then": "Hyper-V"},
									bson.M{"case": mu.APOEqual("$info.hardwareAbstractionTechnology", model.HardwareAbstractionTechnologyXen), "then": "Xen"},
									bson.M{"case": mu.APOEqual("$info.hardwareAbstractionTechnology", model.HardwareAbstractionTechnologyHpvirt), "then": "HP Integrity Virtual Machine"},
								},
								"default": mu.APOConcat("$info.hardwareAbstractionTechnology"),
							},
						},
						"dbInstanceName":        "$database.name",
						"pluggableDatabaseName": "$database.pdbs",
						"environment":           "$environment",
						"options": mu.APOJoin(mu.APOMap(
							mu.APOFilter("$database.licenses", "lic",
								mu.APOAnd(
									mu.APOGreater("$$lic.count", 0),
									mu.APONotEqual("$$lic.name", "Oracle STD"),
									mu.APONotEqual("$$lic.name", "Oracle EXE"),
									mu.APONotEqual("$$lic.name", "Oracle ENT"),
									mu.APOEqual("$$lic.ignored", false),
								),
							),
							"lic",
							"$$lic.name",
						), ", "),
						"usedManagementPacks": mu.APOJoin(mu.APOMap(
							mu.APOFilter("$database.licenses", "lic",
								mu.APOAnd(
									mu.APOGreater("$$lic.count", 0),
									mu.APOOr(
										mu.APOEqual("$$lic.name", "Diagnostics Pack"),
										mu.APOEqual("$$lic.name", "Tuning Pack"),
									),
									mu.APOEqual("$$lic.ignored", false),
								),
							),
							"lic",
							"$$lic.name",
						), ", "),
						"productVersion": mu.APOArrayElemAt(mu.APOSplit("$database.version", " "), 0),
						"productLicenseAllocated": mu.APOLet(
							bson.M{
								"edition": mu.APOArrayElemAt(mu.APOSplit("$database.version", " "), 1),
							},
							bson.M{
								"$switch": bson.M{
									"branches": bson.A{
										bson.M{"case": mu.APOEqual("$$edition", "Enterprise"), "then": "EE"},
										bson.M{"case": mu.APOEqual("$$edition", "Standard"), "then": "SE"},
									},
									"default": mu.APOConcat("$$edition"),
								},
							},
						),
						"licenseMetricAllocated": "processor",
						"usingLicenseCount": mu.APOIfNull(mu.APOArrayElemAt(
							mu.APOMap(
								mu.APOFilter("$database.licenses", "lic",
									mu.APOAnd(
										mu.APOGreater("$$lic.count", 0),
										mu.APOOr(
											mu.APOEqual("$$lic.name", "Oracle STD"),
											mu.APOEqual("$$lic.name", "Oracle EXE"),
											mu.APOEqual("$$lic.name", "Oracle ENT"),
										),
										mu.APOEqual("$$lic.ignored", false),
									),
								),
								"lic",
								"$$lic.count",
							),
							0,
						), 0.0),
						"processorModel":    "$info.cpuModel",
						"processors":        "$info.cpuSockets",
						"coresPerProcessor": "$info.coresPerSocket",
						"threadsPerCore": mu.APOCond(
							mu.APOGreaterOrEqual(mu.APOIndexOfCp("$info.cpuModel", "SPARC"), 0),
							8,
							2,
						),
						"processorSpeed":  "$info.cpuFrequency",
						"operatingSystem": "$info.os",
					}),
					mu.APSet(bson.M{
						"physicalCores": mu.APOCond(mu.APOEqual("$info.cpuSockets", 0), "$coresPerProcessor", bson.M{
							"$multiply": bson.A{"$coresPerProcessor", "$processors"},
						}),
					}),
				)),
				mu.APOptionalSortingStage(filters.SortBy, filters.SortDesc),
				mu.APOptionalPagingStage(filters.PageNumber, filters.PageSize),
			)),
		),
	)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(context.TODO(), out); err != nil {
		return utils.NewError(err, "Decode ERROR")
	}

	return nil
}

func getClusterFilterStep(cl *string) interface{} {
	if cl == nil {
		return mu.APMatch(bson.M{
			"cluster": nil,
		})
	} else if *cl != "" {
		return mu.APMatch(bson.M{
			"cluster": primitive.Regex{Pattern: regexp.QuoteMeta(*cl), Options: "i"},
		})
	} else {
		return bson.A{}
	}
}

func getIsMemberOfClusterFilterStep(member *bool) interface{} {
	if member != nil {
		return mu.APMatch(mu.QOExpr(
			mu.APOEqual(*member, mu.APOOr("$clusterMembershipStatus.oracleClusterware", "$clusterMembershipStatus.veritasClusterServer", "$clusterMembershipStatus.sunCluster", "$clusterMembershipStatus.hacmp")),
		))
	}

	return bson.A{}
}

// GetHost fetch all informations about a host in the database
func (md *MongoDatabase) GetHost(hostname string, olderThan time.Time, raw bool) (*dto.HostData, error) {
	var cur *mongo.Cursor

	var err error

	//Get host technology
	technology, errTech := md.getHostTechnology(hostname, olderThan)
	if errTech != nil {
		return nil, utils.NewError(errTech, "DB ERROR")
	}

	//Find the matching hostdata
	switch technology {
	case model.TechnologyOracleDatabase, model.TechnologyOracleExadata:
		cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
			context.TODO(),
			mu.MAPipeline(
				getBaseHost(hostname, technology, olderThan, raw),
				mu.APOptionalStage(!raw, mu.MAPipeline(
					mu.APLookupPipeline(
						"hosts",
						bson.M{
							"hn": "$hostname",
							"ca": "$createdAt",
						},
						"history",
						mu.MAPipeline(
							mu.APMatch(mu.QOExpr(mu.APOAnd(mu.APOEqual("$hostname", "$$hn"), mu.APOGreaterOrEqual("$$ca", "$createdAt")))),
							mu.APProject(bson.M{
								"createdAt": 1,
								"features.oracle.database.databases.name":          1,
								"features.oracle.database.databases.datafileSize":  1,
								"features.oracle.database.databases.segmentsSize":  1,
								"features.oracle.database.databases.allocable":     1,
								"features.oracle.database.databases.dailyCPUUsage": 1,
								"totalDailyCPUUsage":                               mu.APOSumReducer("$features.oracle.database.databases", mu.APOConvertToDoubleOrZero("$$this.dailyCPUUsage")),
							}),
						),
					),
					mu.APSet(bson.M{
						"features.oracle.database.databases": mu.APOMap(
							"$features.oracle.database.databases",
							"db",
							mu.APOMergeObjects(
								"$$db",
								bson.M{
									"changes": mu.APOFilter(
										mu.APOMap("$history", "hh", mu.APOMergeObjects(
											bson.M{"updated": "$$hh.createdAt"},
											mu.APOArrayElemAt(mu.APOFilter("$$hh.features.oracle.database.databases", "hdb", mu.APOEqual("$$hdb.name", "$$db.name")), 0),
										)),
										"time_frame",
										"$$time_frame.segmentsSize",
									),
								},
							),
						),
					}),
					mu.APUnset(
						"features.oracle.database.databases.changes.name",
						"history.features",
					),
					mu.APSet(bson.M{
						"features.oracle": mu.APOCond(mu.APOEqual("$features.oracle.database.databases", nil), nil, "$features.oracle"),
					}),
				)),
			),
		)
		if err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	case model.TechnologyMicrosoftSQLServer:
		cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
			context.TODO(),
			mu.MAPipeline(
				getBaseHost(hostname, technology, olderThan, raw),
			),
		)
		if err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	case model.TechnologyOracleMySQL:
		cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
			context.TODO(),
			mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				getBaseHost(hostname, technology, olderThan, raw),
				mu.APOptionalStage(!raw, mu.MAPipeline(
					mu.APLookupPipeline(
						"hosts",
						bson.M{
							"hn": "$hostname",
							"ca": "$createdAt",
						},
						"history",
						mu.MAPipeline(
							mu.APMatch(mu.QOExpr(mu.APOAnd(mu.APOEqual("$hostname", "$$hn"), mu.APOGreaterOrEqual("$$ca", "$createdAt")))),
							mu.APUnwind("$features.mysql.instances"),
							mu.APProject(bson.M{
								"createdAt":                1,
								"features.mysql.instances": bson.A{"$features.mysql.instances"},
								"totalAllocation":          mu.APOSum("$features.mysql.instances.tableSchemas.allocation"),
							}),
							mu.APProject(bson.M{
								"createdAt":                     1,
								"features.mysql.instances.name": 1,
								"totalAllocation":               1,
							}),
						),
					),
					mu.APSet(bson.M{
						"features.mysql.instances": mu.APOMap(
							"$features.mysql.instances",
							"db",
							mu.APOMergeObjects(
								"$$db",
								bson.M{
									"changes": mu.APOFilter(
										mu.APOMap("$history", "hh", mu.APOMergeObjects(
											bson.M{"updated": "$$hh.createdAt"},
											bson.M{"allocation": "$$hh.totalAllocation"},
											mu.APOArrayElemAt(mu.APOFilter("$$hh.features.mysql.instances", "hdb", mu.APOEqual("$$hdb.name", "$$db.name")), 0),
										)),
										"time_frame",
										"$$time_frame.name",
									),
								},
							),
						),
					}),
					mu.APUnset(
						"features.mysql.instances.changes.name",
						"history.features",
					),
					mu.APSet(bson.M{
						"features.mysql": mu.APOCond(mu.APOEqual("$features.mysql.instances", nil), nil, "$features.mysql"),
					}),
				)),
			),
		)
		if err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	case model.TechnologyPostgreSQLPostgreSQL:
		cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
			context.TODO(),
			mu.MAPipeline(
				getBaseHost(hostname, technology, olderThan, raw),
			),
		)
		if err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	case model.TechnologyMongoDBMongoDB:
		cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
			context.TODO(),
			mu.MAPipeline(
				getBaseHost(hostname, technology, olderThan, raw),
			),
		)
		if err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	default:
		cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
			context.TODO(),
			mu.MAPipeline(
				getBaseHost(hostname, technology, olderThan, raw),
			),
		)
		if err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return nil, utils.ErrHostNotFound
	}

	//Decode the document
	var host dto.HostData
	if err := cur.Decode(&host); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return &host, nil
}

func getBaseHost(hostname string, technology string, olderThan time.Time, raw bool) bson.A {
	return mu.MAPipeline(
		FilterByOldnessSteps(olderThan),
		mu.APMatch(bson.M{
			"hostname": hostname,
		}),
		mu.APAddFields(
			bson.M{
				"technology": technology,
			}),
		mu.APOptionalStage(!raw, mu.MAPipeline(
			mu.APLookupPipeline("alerts", bson.M{"hn": "$hostname"}, "alerts", mu.MAPipeline(
				mu.APMatch(mu.QOExpr(mu.APOEqual("$otherInfo.hostname", "$$hn"))),
			)),
			AddAssociatedClusterNameAndVirtualizationNode(olderThan),
		)),
	)
}

func (md *MongoDatabase) GetHostData(hostname string, olderThan time.Time) (*model.HostDataBE, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		Aggregate(
			context.TODO(),
			mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APMatch(bson.M{
					"hostname": hostname,
				}),
			),
		)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return nil, fmt.Errorf("%w: %s", utils.ErrHostNotFound, hostname)
	}

	var hostdata model.HostDataBE
	if err := cur.Decode(&hostdata); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return &hostdata, nil
}

func (md *MongoDatabase) GetHostDatas(filter dto.GlobalFilter) ([]model.HostDataBE, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		Aggregate(
			context.TODO(),
			mu.MAPipeline(
				FilterByOldnessSteps(filter.OlderThan),
				FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			),
		)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	hostdatas := make([]model.HostDataBE, 0)
	if err := cur.All(context.TODO(), &hostdatas); err != nil {
		return nil, utils.NewError(err, "DECODE ERROR")
	}

	return hostdatas, nil
}

// ListAllLocations list all available locations
func (md *MongoDatabase) ListAllLocations(location string, environment string, olderThan time.Time) ([]string, error) {
	var out = make([]string, 0)

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id": "$location",
			}),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]string
		if cur.Decode(&item) != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}

		out = append(out, item["_id"])
	}

	return out, nil
}

// ListEnvironments list environments
func (md *MongoDatabase) ListEnvironments(location string, environment string, olderThan time.Time) ([]string, error) {
	var out = make([]string, 0)

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(bson.M{
				"_id": "$environment",
			}),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]string
		if cur.Decode(&item) != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}

		out = append(out, item["_id"])
	}

	return out, nil
}

// FindHostData find the current hostdata with a certain hostname
func (md *MongoDatabase) FindHostData(hostname string) (model.HostDataBE, error) {
	//Find the hostdata
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").FindOne(context.TODO(), bson.M{
		"dismissedAt": nil,
		"hostname":    hostname,
		"archived":    false,
	})
	if res.Err() == mongo.ErrNoDocuments {
		return model.HostDataBE{}, utils.ErrHostNotFound
	} else if res.Err() != nil {
		return model.HostDataBE{}, utils.NewError(res.Err(), "DB ERROR")
	}

	//Decode the data
	var out model.HostDataBE
	if err := res.Decode(&out); err != nil {
		return model.HostDataBE{}, utils.NewError(res.Err(), "DB ERROR")
	}

	//Return it!
	return out, nil
}

// ReplaceHostData adds a new hostdata to the database
func (md *MongoDatabase) ReplaceHostData(hostData model.HostDataBE) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").ReplaceOne(context.TODO(),
		bson.M{
			"_id": hostData.ID,
		},
		hostData,
	)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

// ExistHostdata return true if exist a non-dismissed hostdata with the hostname equal hostname
func (md *MongoDatabase) ExistHostdata(hostname string) (bool, error) {
	//Count the number of new NO_DATA alerts associated to the host
	val, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").CountDocuments(context.TODO(), bson.M{
		"dismissedAt": nil,
		"archived":    false,
		"hostname":    hostname,
	}, &options.CountOptions{
		Limit: utils.Intptr(1),
	})
	if err != nil {
		return false, utils.NewError(err, "DB ERROR")
	}

	//Return true if the count > 0
	return val > 0, nil
}

// DismissHost dismiss the specified host
func (md *MongoDatabase) DismissHost(hostname string) error {
	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").UpdateMany(context.TODO(), bson.M{
		"hostname":    hostname,
		"dismissedAt": nil,
	}, mu.UOSet(bson.M{
		"dismissedAt": time.Now(),
		"archived":    true,
	})); err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

// GetHostMinValidCreatedAtDate get the host's minimun valid CreatedAt date
func (md *MongoDatabase) GetHostMinValidCreatedAtDate(hostname string) (time.Time, error) {
	var createdAt map[string]interface{}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: 1}})
	findOptions.SetLimit(1)
	findOptions.SetProjection(bson.M{"_id": 0, "createdAt": 1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Find(context.TODO(), bson.M{
		"hostname":    hostname,
		"dismissedAt": nil,
	},
		findOptions,
	)
	if err != nil {
		return time.Time{}, utils.NewError(err, "DB ERROR")
	}

	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return time.Time{}, utils.NewError(errors.New("Invalid result"), "DB ERROR")
	}

	if err := cur.Decode(&createdAt); err != nil {
		return time.Time{}, utils.NewError(err, "DB ERROR")
	}

	return createdAt["createdAt"].(primitive.DateTime).Time().UTC(), nil
}

// GetListValidHostsByRangeDates get list of valid hosts by range dates
func (md *MongoDatabase) GetListValidHostsByRangeDates(from time.Time, to time.Time) ([]string, error) {
	var hosts = make([]string, 0)

	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"hostname",
		bson.M{
			"dismissedAt": nil,
			"createdAt":   bson.M{"$gte": from, "$lte": to},
		},
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	for _, val := range values {
		hosts = append(hosts, val.(string))
	}

	return hosts, nil
}

// GetListDismissedHostsByRangeDates get list of dismissed hosts by range dates
func (md *MongoDatabase) GetListDismissedHostsByRangeDates(from time.Time, to time.Time) ([]string, error) {
	var hosts = make([]string, 0)

	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"hostname",
		bson.M{
			"dismissedAt": bson.M{"$gte": from, "$lte": to},
		},
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	for _, val := range values {
		hosts = append(hosts, val.(string))
	}

	return hosts, nil
}

// ExistNotInClusterHost return true if the host specified by hostname exist and it is not in cluster, otherwise false
func (md *MongoDatabase) ExistNotInClusterHost(hostname string) (bool, error) {
	//check that the host exist
	var out = make([]struct{}, 0)

	//Find the matching alerts
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"dismissedAt": nil,
				"archived":    false,
				"hostname":    hostname,
			}),
			mu.APProject(bson.M{
				"hostname": true,
			}),
			mu.APLookupPipeline("hosts", bson.M{"hn": "$hostname"}, "cluster", mu.MAPipeline(
				mu.APMatch(bson.M{
					"dismissedAt": nil,
					"archived":    false,
				}),
				mu.APUnwind("$clusters"),
				mu.APReplaceWith("$clusters"),
				mu.APUnwind("$vms"),
				mu.APSet(bson.M{
					"vms.clusterName": "$name",
				}),
				mu.APMatch(mu.QOExpr(mu.APOEqual("$vms.hostname", "$$hn"))),
				mu.APLimit(1),
			)),
			mu.APMatch(mu.QOExpr(mu.APOEqual(mu.APOSize("$cluster"), 0))),
		),
	)
	if err != nil {
		return false, utils.NewError(err, "DB ERROR")
	}

	//Decode the documents
	if err = cur.All(context.TODO(), &out); err != nil {
		return false, utils.NewError(err, "Decode ERROR")
	}

	return len(out) > 0, nil
}

func (md *MongoDatabase) getHostTechnology(hostname string, olderThan time.Time) (string, error) {
	var result = make(map[string]bool)

	var out string

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "hostname", Value: hostname},
					{Key: "archived", Value: false},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "Oracle/Database",
						Value: bson.D{
							{Key: "$eq",
								Value: bson.A{
									bson.D{{Key: "$type", Value: "$features.oracle"}},
									"object",
								},
							},
						},
					},
					{Key: "Microsoft/SQLServer",
						Value: bson.D{
							{Key: "$eq",
								Value: bson.A{
									bson.D{{Key: "$type", Value: "$features.microsoft"}},
									"object",
								},
							},
						},
					},
					{Key: "Oracle/MySQL",
						Value: bson.D{
							{Key: "$eq",
								Value: bson.A{
									bson.D{{Key: "$type", Value: "$features.mysql"}},
									"object",
								},
							},
						},
					},
					{Key: "PostgreSQL/PostgreSQL",
						Value: bson.D{
							{Key: "$eq",
								Value: bson.A{
									bson.D{{Key: "$type", Value: "$features.postgresql"}},
									"object",
								},
							},
						},
					},
					{Key: "MongoDB/MongoDB",
						Value: bson.D{
							{Key: "$eq",
								Value: bson.A{
									bson.D{{Key: "$type", Value: "$features.mongodb"}},
									"object",
								},
							},
						},
					},
				},
			},
		},
		bson.D{{Key: "$unset", Value: "_id"}},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		pipeline)
	if err != nil {
		return "", utils.NewError(err, "DB ERROR")
	}

	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return out, nil
	}

	if err := cur.Decode(&result); err != nil {
		return "", utils.NewError(err, "DB ERROR")
	}

	for technology, exist := range result {
		if exist {
			out = technology
			break
		}
	}

	return out, nil
}

func (md *MongoDatabase) SearchHostMysqlLMS(filter dto.SearchHostsAsLMS) ([]dto.MySqlHostLMS, error) {
	ctx := context.TODO()

	result := make([]dto.MySqlHostLMS, 0)

	pipeline := bson.A{}

	if !filter.From.IsZero() {
		pipeline = append(pipeline, bson.D{{Key: "$match",
			Value: bson.D{
				{Key: "createdAt",
					Value: bson.D{
						{Key: "$gte", Value: filter.From},
					},
				},
			},
		}})
	}

	if !filter.To.IsZero() {
		pipeline = append(pipeline, bson.D{{Key: "$match",
			Value: bson.D{
				{Key: "createdAt",
					Value: bson.D{
						{Key: "$lt", Value: filter.To},
					},
				},
			},
		}})
	}

	if filter.Location != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "location",
						Value: bson.D{
							{Key: "$in",
								Value: strings.Split(filter.Location, ","),
							},
						},
					},
				},
			},
		})
	}

	if filter.Environment != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "environment", Value: filter.Environment}}}})
	}

	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "archived", Value: false}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.mysql.instances"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "physicalServerName",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.D{
									{Key: "if",
										Value: bson.D{
											{Key: "$in",
												Value: bson.A{
													"$info.hardwareAbstractionTechnology",
													bson.A{
														"KVM",
														"OVM",
														"VMWARE",
													},
												},
											},
										},
									},
									{Key: "then", Value: "$clusters.name"},
									{Key: "else", Value: ""},
								},
							},
						},
					},
					{Key: "virtualServerName", Value: "$hostname"},
					{Key: "virtualization", Value: "$info.hardwareAbstractionTechnology"},
					{Key: "dbInstanceName", Value: "$features.mysql.instances.name"},
					{Key: "environmentUsage", Value: "$environment"},
					{Key: "productVersion", Value: "$features.mysql.instances.version"},
					{Key: "productLicenseAllocated", Value: "$features.mysql.instances.edition"},
					{Key: "licenseMetricAllocated", Value: "HOST"},
					{Key: "numberOfLicenseUsed",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.D{
									{Key: "if",
										Value: bson.D{
											{Key: "$eq",
												Value: bson.A{
													"$features.mysql.instances.license.ignored",
													false,
												},
											},
										},
									},
									{Key: "then", Value: "$features.mysql.instances.license.count"},
									{Key: "else", Value: 0},
								},
							},
						},
					},
					{Key: "processorModel", Value: "$info.cpuModel"},
					{Key: "sockets", Value: "$info.cpuSockets"},
					{Key: "physicalCores", Value: "$info.cpuCores"},
					{Key: "threadsPerCore", Value: "$info.threadsPerCore"},
					{Key: "os", Value: "$info.os"},
				},
			},
		})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) FindVirtualHostWithoutCluster() ([]dto.VirtualHostWithoutCluster, error) {
	res := make([]dto.VirtualHostWithoutCluster, 0)

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "archived", Value: false},
					{Key: "info.hardwareAbstraction", Value: "VIRT"},
					{Key: "info.hardwareAbstractionTechnology", Value: bson.D{{Key: "$ne", Value: "HPVRT"}}},
				},
			},
		},
	}

	pipeline = append(pipeline, AddAssociatedClusterNameAndVirtualizationNode(utils.MAX_TIME)...)
	pipeline = append(pipeline, bson.D{
		{Key: "$match",
			Value: bson.D{
				{Key: "cluster", Value: primitive.Null{}},
			},
		},
	})
	pipeline = append(pipeline, bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "hostname", Value: 1},
				{Key: "hardwareAbstractionTechnology", Value: "$info.hardwareAbstractionTechnology"},
			},
		},
	})

	ctx := context.TODO()

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}
