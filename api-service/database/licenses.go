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
	"errors"
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// SearchLicenses search licenses
func (md *MongoDatabase) SearchLicenses(mode string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	var isFull bool
	if mode == "full" {
		isFull = true
	} else if mode == "summary" {
		isFull = false
	} else {
		return nil, utils.NewAdvancedErrorPtr(errors.New("Wrong mode value"), "")
	}

	//Find the informations
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APLookupPipeline("hosts",
				bson.M{
					"ln": "$_id",
				},
				"used",
				mu.MAPipeline(
					FilterByOldnessSteps(olderThan),
					FilterByLocationAndEnvironmentSteps(location, environment),
					mu.APProject(bson.M{
						"hostname": 1,
						"databases": mu.APOReduce(
							mu.APOFilter(
								mu.APOMap("$features.oracle.database.databases",
									"db",
									bson.M{
										"name": "$$db.name",
										"count": mu.APOLet(
											bson.M{
												"val": mu.APOArrayElemAt(mu.APOFilter("$$db.licenses", "lic", mu.APOEqual("$$lic.name", "$$ln")), 0),
											},
											"$$val.count",
										),
									}),
								"db",
								mu.APOGreater("$$db.count", 0),
							),
							bson.M{"count": 0, "dbs": bson.A{}},
							bson.M{
								"count": mu.APOMax("$$value.count", "$$this.count"),
								"dbs": bson.M{
									"$concatArrays": bson.A{
										"$$value.dbs",
										bson.A{"$$this.name"},
									},
								},
							},
						),
					}),
					mu.APMatch(bson.M{
						"databases.count": bson.M{
							"$gt": 0,
						},
					}),
					AddAssociatedClusterNameAndVirtualizationNode(olderThan),
					mu.APGroup(mu.BsonOptionalExtension(isFull, bson.M{
						"_id": mu.APOCond(
							"$clusterName",
							mu.APOConcat("cluster_§$#$§_", "$clusterName"),
							mu.APOConcat("hostname_§$#$§_", "$hostname"),
						),
						"license":    mu.APOMaxAggr("$databases.count"),
						"clusterCpu": mu.APOMaxAggr("$clusterCpu"),
					}, bson.M{
						"hosts": mu.APOPush(bson.M{
							"hostname":  "$hostname",
							"databases": "$databases.dbs",
						}),
					})),
					mu.APSet(bson.M{
						"license": mu.APOCond(
							"$clusterCpu",
							mu.APODivide("$clusterCpu", 2),
							"$license",
						),
					}),
					mu.APGroup(mu.BsonOptionalExtension(isFull, bson.M{
						"_id":   0,
						"value": mu.APOSum("$license"),
					}, bson.M{
						"hosts": mu.APOPush("$hosts"),
					})),
					mu.APOptionalStage(isFull, mu.MAPipeline(
						mu.APUnwind("$hosts"),
						mu.APUnwind("$hosts"),
						mu.APGroup(bson.M{
							"_id":   0,
							"value": mu.APOMaxAggr("$value"),
							"hosts": mu.APOPush("$hosts"),
						}),
					)),
				)),
			mu.APSet(bson.M{
				"used": mu.APOArrayElemAt("$used", 0),
			}),
			mu.APOptionalStage(isFull, mu.APSet(bson.M{
				"hosts": mu.APOIfNull("$used.hosts", bson.A{}),
			})),
			mu.APSet(bson.M{
				"used": mu.APOIfNull(mu.APOCeil("$used.value"), 0),
			}),
			mu.APSet(bson.M{
				"compliance": mu.APOGreaterOrEqual(
					mu.APOCond("$unlimited", "$used", "$count"),
					"$used",
				),
				"totalCost": bson.M{
					"$multiply": bson.A{"$used", "$costPerProcessor"},
				},
				"paidCost": bson.M{
					"$multiply": bson.A{
						mu.APOCond("$unlimited", "$used", "$count"),
						"$costPerProcessor",
					},
				},
			}),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	var out []interface{} = make([]interface{}, 0)

	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// ListLicenses list licenses
func (md *MongoDatabase) ListLicenses(sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	//Find the informations
	cursor, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$features.oracle.database.databases"),
			mu.APUnwind("$features.oracle.database.databases.licenses"),
			mu.APMatch(bson.M{"features.oracle.database.databases.licenses.count": bson.M{"$gt": 0}}),
			mu.APProject(
				bson.M{
					"_id":              0,
					"hostname":         1,
					"dbName":           "$features.oracle.database.databases.name",
					"licenseName":      "$features.oracle.database.databases.licenses.name",
					"purchasedLicense": "$features.oracle.database.databases.licenses.count",
				},
			),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)

	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	var out []interface{} = make([]interface{}, 0)

	for cursor.Next(context.TODO()) {
		var item map[string]interface{}
		if cursor.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, &item)
	}
	return out, nil
}

// GetLicense get a certain license
func (md *MongoDatabase) GetLicense(name string, olderThan time.Time) (interface{}, utils.AdvancedErrorInterface) {
	var out map[string]interface{}
	//Find the informations

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"_id": name,
			}),
			mu.APLookupPipeline("hosts", bson.M{
				"ln": "$_id",
			}, "used", mu.MAPipeline(
				FilterByOldnessSteps(olderThan),
				mu.APProject(bson.M{
					"hostname": 1,
					"databases": mu.APOReduce(
						mu.APOFilter(
							mu.APOMap("$features.oracle.database.databases", "db", bson.M{
								"Name": "$$db.name",
								"count": mu.APOLet(
									bson.M{
										"val": mu.APOArrayElemAt(mu.APOFilter("$$db.licenses", "lic", mu.APOEqual("$$lic.name", "$$ln")), 0),
									},
									"$$val.count",
								),
							}),
							"db",
							mu.APOGreater("$$db.count", 0),
						),
						bson.M{"count": 0, "dbs": bson.A{}},
						bson.M{
							"count": mu.APOMax("$$value.count", "$$this.count"),
							"dbs": bson.M{
								"$concatArrays": bson.A{
									"$$value.dbs",
									bson.A{"$$this.name"},
								},
							},
						},
					),
				}),
				mu.APMatch(bson.M{
					"databases.count": bson.M{
						"$gt": 0,
					},
				}),
				AddAssociatedClusterNameAndVirtualizationNode(olderThan),
				mu.APGroup(bson.M{
					"_id": mu.APOCond(
						"$clusterName",
						mu.APOConcat("cluster_§$#$§_", "$clusterName"),
						mu.APOConcat("hostname_§$#$§_", "$hostname"),
					),
					"license":    mu.APOMaxAggr("$databases.count"),
					"clusterCpu": mu.APOMaxAggr("$clusterCpu"),
					"hosts": mu.APOPush(bson.M{
						"hostname":  "$hostname",
						"databases": "$databases.dbs",
					}),
				}),
				mu.APSet(bson.M{
					"license": mu.APOCond(
						"$clusterCpu",
						mu.APODivide("$clusterCpu", 2),
						"$license",
					),
				}),
				mu.APGroup(bson.M{
					"_id":   0,
					"value": mu.APOSum("$license"),
					"hosts": mu.APOPush("$hosts"),
				}),
				mu.APUnwind("$hosts"),
				mu.APUnwind("$hosts"),
				mu.APGroup(bson.M{
					"_id":   0,
					"value": mu.APOMaxAggr("$value"),
					"hosts": mu.APOPush("$hosts"),
				}),
			)),
			mu.APSet(bson.M{
				"used": mu.APOArrayElemAt("$used", 0),
			}),
			mu.APSet(bson.M{
				"hosts": mu.APOIfNull("$used.hosts", bson.A{}),
			}),
			mu.APSet(bson.M{
				"used": mu.APOIfNull(mu.APOCeil("$used.value"), 0),
			}),
			mu.APSet(bson.M{
				"compliance": mu.APOGreaterOrEqual(
					mu.APOCond("$unlimited", "$used", "$count"),
					"$used",
				),
				"totalCost": bson.M{
					"$multiply": bson.A{"$used", "$costPerProcessor"},
				},
				"paidCost": bson.M{
					"$multiply": bson.A{
						mu.APOCond("$unlimited", "$used", "$count"),
						"$costPerProcessor",
					},
				},
			}),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return nil, utils.AerrLicenseNotFound
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}

// SetLicenseCount set the count of a certain license
func (md *MongoDatabase) SetLicenseCount(name string, count int) utils.AdvancedErrorInterface {
	//Find the informations
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").UpdateOne(context.TODO(), bson.M{
		"_id": name,
	}, mu.UOSet(bson.M{
		"count": count,
	}))
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Check the existance of the result
	if res.MatchedCount == 0 {
		return utils.AerrLicenseNotFound
	} else {
		return nil
	}
}

// SetLicenseCostPerProcessor set the cost per processor of a certain license
func (md *MongoDatabase) SetLicenseCostPerProcessor(name string, count float64) utils.AdvancedErrorInterface {
	//Find the informations
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").UpdateOne(context.TODO(), bson.M{
		"_id": name,
	}, mu.UOSet(bson.M{
		"costPerProcessor": count,
	}))
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Check the existance of the result
	if res.MatchedCount == 0 {
		return utils.AerrLicenseNotFound
	} else {
		return nil
	}
}

// SetLicenseUnlimitedStatus set the unlimited status of a certain license
func (md *MongoDatabase) SetLicenseUnlimitedStatus(name string, unlimitedStatus bool) utils.AdvancedErrorInterface {
	//Find the informations
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("licenses").UpdateOne(context.TODO(), bson.M{
		"_id": name,
	}, mu.UOSet(bson.M{
		"unlimited": unlimitedStatus,
	}))
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Check the existance of the result
	if res.MatchedCount == 0 {
		return utils.AerrLicenseNotFound
	} else {
		return nil
	}
}
