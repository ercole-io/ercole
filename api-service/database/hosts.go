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
	"fmt"
	"regexp"
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
					"createdAt":               true,
					"hostname":                true,
					"location":                true,
					"environment":             true,
					"agentVersion":            true,
					"info":                    true,
					"clusterMembershipStatus": true,

					"virtualizationNode": true,
					"cluster":            true,
					"databases": bson.M{
						model.TechnologyOracleDatabase:     "$features.oracle.database.databases.name",
						model.TechnologyMicrosoftSQLServer: "$features.microsoft.sqlServer.instances.name",
						model.TechnologyOracleMySQL:        "$features.mysql.instances.name",
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
						"vmwareOrOVM": mu.APOOr(mu.APOEqual("$info.hardwareAbstractionTechnology", model.HardwareAbstractionTechnologyVmware), mu.APOEqual("$info.hardwareAbstractionPlatform", model.HardwareAbstractionTechnologyOvm)),
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
						"physicalServerName": mu.APOCond("$vmwareOrOVM", mu.APOIfNull("$cluster", ""), "$hostname"),
						"virtualServerName":  mu.APOCond("$vmwareOrOVM", "$hostname", mu.APOIfNull("$cluster", "")),
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
							//TODO Replace with lookup with oracle_database_license_types
							mu.APOFilter("$database.licenses", "lic", mu.APOAnd(mu.APOGreater("$$lic.count", 0), mu.APONotEqual("$$lic.name", "Oracle STD"), mu.APONotEqual("$$lic.name", "Oracle EXE"), mu.APONotEqual("$$lic.name", "Oracle ENT"))),
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
											//TODO Replace with lookup with oracle_database_license_types
											mu.APOEqual("$$lic.name", "Oracle STD"),
											mu.APOEqual("$$lic.name", "Oracle EXE"),
											mu.APOEqual("$$lic.name", "Oracle ENT"),
										),
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
func (md *MongoDatabase) GetHost(hostname string, olderThan time.Time, raw bool) (interface{}, error) {
	var out map[string]interface{}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			mu.APMatch(bson.M{
				"hostname": hostname,
			}),
			mu.APOptionalStage(!raw, mu.MAPipeline(
				mu.APLookupPipeline("alerts", bson.M{"hn": "$hostname"}, "alerts", mu.MAPipeline(
					mu.APMatch(mu.QOExpr(mu.APOEqual("$otherInfo.hostname", "$$hn"))),
				)),
				AddAssociatedClusterNameAndVirtualizationNode(olderThan),
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
			)),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return nil, utils.ErrHostNotFound
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return out, nil
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

// ListLocations list locations
func (md *MongoDatabase) ListLocations(location string, environment string, olderThan time.Time) ([]string, error) {
	var out []string = make([]string, 0)

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
	var out []string = make([]string, 0)

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
		"hostname": hostname,
		"archived": false,
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

	var out2 map[string]interface{}
	if err := res.Decode(&out2); err != nil {
		// return model.HostDataBE{}, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
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

// ExistHostdata return true if exist a non-archived hostdata with the hostname equal hostname
func (md *MongoDatabase) ExistHostdata(hostname string) (bool, error) {
	//Count the number of new NO_DATA alerts associated to the host
	val, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").CountDocuments(context.TODO(), bson.M{
		"archived": false,
		"hostname": hostname,
	}, &options.CountOptions{
		Limit: utils.Intptr(1),
	})
	if err != nil {
		return false, utils.NewError(err, "DB ERROR")
	}

	//Return true if the count > 0
	return val > 0, nil
}

// ArchiveHost archive the specified host
func (md *MongoDatabase) ArchiveHost(hostname string) error {
	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").UpdateOne(context.TODO(), bson.M{
		"hostname": hostname,
		"archived": false,
	}, mu.UOSet(bson.M{
		"archived": true,
	})); err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

// ExistNotInClusterHost return true if the host specified by hostname exist and it is not in cluster, otherwise false
func (md *MongoDatabase) ExistNotInClusterHost(hostname string) (bool, error) {
	//check that the host exist
	var out []struct{} = make([]struct{}, 0)

	//Find the matching alerts
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"archived": false,
				"hostname": hostname,
			}),
			mu.APProject(bson.M{
				"hostname": true,
			}),
			mu.APLookupPipeline("hosts", bson.M{"hn": "$hostname"}, "cluster", mu.MAPipeline(
				mu.APMatch(bson.M{
					"archived": false,
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
