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

	"github.com/amreo/mu"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const oracleDbAgreementsCollection = "oracle_database_agreements"

// InsertOracleDatabaseAgreement insert an Oracle/Database agreement into the database
func (md *MongoDatabase) InsertOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) (*mongo.InsertOneResult, utils.AdvancedErrorInterface) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		InsertOne(context.TODO(), agreement)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return res, nil
}

// GetOracleDatabaseAgreement return the agreement specified by id
func (md *MongoDatabase) GetOracleDatabaseAgreement(agreementID string) (*model.OracleDatabaseAgreement, utils.AdvancedErrorInterface) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		FindOne(context.TODO(), bson.M{
			"agreementID": agreementID,
		})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.AerrOracleDatabaseAgreementNotFound
	} else if res.Err() != nil {
		return nil, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
	}

	var out model.OracleDatabaseAgreement
	if err := res.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
	}
	return &out, nil
}

// GetOracleDatabaseAgreementByAssociatedLicenseType return the agreement specified by an associated part id
func (md *MongoDatabase) GetOracleDatabaseAgreementByAssociatedLicenseType(associateLicenseTypeID primitive.ObjectID) (*model.OracleDatabaseAgreement, utils.AdvancedErrorInterface) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		FindOne(context.TODO(), bson.M{
			"licenseTypes._id": associateLicenseTypeID,
		})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.AerrOracleDatabaseAgreementNotFound
	} else if res.Err() != nil {
		return nil, utils.NewAdvancedErrorPtr(res.Err(), "DB ERROR")
	}

	var out model.OracleDatabaseAgreement
	if err := res.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
	}
	return &out, nil
}

// UpdateOracleDatabaseAgreement update an Oracle/Database agreement in the database
func (md *MongoDatabase) UpdateOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) utils.AdvancedErrorInterface {
	result, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		ReplaceOne(context.TODO(), bson.M{
			"_id": agreement.ID,
		}, agreement)
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	if result.MatchedCount != 1 {
		return utils.AerrOracleDatabaseAgreementNotFound
	}

	return nil
}

// RemoveOracleDatabaseAgreement remove an Oracle/Database agreement from the database
func (md *MongoDatabase) RemoveOracleDatabaseAgreement(id primitive.ObjectID) utils.AdvancedErrorInterface {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		DeleteOne(context.TODO(), bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	if res.DeletedCount == 0 {
		return utils.AerrOracleDatabaseAgreementNotFound
	}
	return nil
}

// ListOracleDatabaseAgreements lists the Oracle/Database agreements
func (md *MongoDatabase) ListOracleDatabaseAgreements() ([]dto.OracleDatabaseAgreementFE, utils.AdvancedErrorInterface) {
	var out []dto.OracleDatabaseAgreementFE = make([]dto.OracleDatabaseAgreementFE, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		Aggregate(
			context.TODO(),
			mu.MAPipeline(
				mu.APUnwind("$licenseTypes"),
				mu.APLookupSimple("oracle_database_license_types", "licenseTypes.licenseTypeID", "_id", "licenseType"),
				mu.APUnwind("$licenseType"),
				mu.APSet(bson.M{
					"_id": "$licenseTypes._id",

					"licenseTypeID":   "$licenseTypes.licenseTypeID",
					"itemDescription": "$licenseType.itemDescription",
					"metric":          "$licenseType.metric",

					"referenceNumber": "$licenseTypes.referenceNumber",
					"unlimited":       "$licenseTypes.unlimited",
					"count":           "$licenseTypes.count",
					"catchAll":        "$licenseTypes.catchAll",
					"restricted":      "$licenseTypes.restricted",

					"hosts": mu.APOMap("$licenseTypes.hosts", "hn", bson.M{
						"hostname": "$$hn",
					}),

					"availableCount": "$licenseTypes.count",
					//TODO And other licenses types?
					"licensesCount": mu.APOCond(
						mu.APOOr(
							mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricProcessorPerpetual),
							mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricComputerPerpetual)),
						"$licenseTypes.count",
						0),
					"usersCount": mu.APOCond(
						mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricNamedUserPlusPerpetual), "$licenseTypes.count", 0),
				}),
				mu.APUnset("licenseType"),
				mu.APUnset("licenseTypes"),
			),
		)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	if err = cur.All(context.TODO(), &out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
	}

	return out, nil
}

// ListHostUsingOracleDatabaseLicenses lists the hosts/clusters that need to be licensed by Oracle/Database agreements
func (md *MongoDatabase) ListHostUsingOracleDatabaseLicenses() ([]dto.HostUsingOracleDatabaseLicenses, utils.AdvancedErrorInterface) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		Aggregate(
			context.TODO(),
			mu.MAPipeline(
				mu.APMatch(bson.M{
					"archived":                 false,
					"features.oracle.database": mu.QONotEqual(nil),
					"$expr":                    mu.APOGreater(mu.APOSize("$features.oracle.database.databases"), 0),
				}),
				mu.APProject(bson.M{
					"hostname": true,
					"licenses": "$features.oracle.database.databases.licenses",
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
				mu.APSet(bson.M{
					"cluster": mu.APOArrayElemAt("$cluster", 0),
				}),
				mu.APAddFields(bson.M{
					"cluster":    "$cluster.name",
					"clusterCpu": "$cluster.cpu",
				}),
				mu.APUnwind("$licenses"),
				mu.APUnwind("$licenses"),
				mu.APGroup(bson.M{
					"_id": bson.M{
						"hostname":      "$hostname",
						"cluster":       "$cluster",
						"clusterCpu":    "$clusterCpu",
						"licenseTypeID": "$licenses.licenseTypeID",
					},
					"licenseCount": mu.APOMaxAggr("$licenses.count"),
				}),
				mu.APMatch(bson.M{
					"licenseCount": bson.M{
						"$gt": 0,
					},
				}),
				mu.APGroup(bson.M{
					"_id": bson.M{
						"licenseTypeID": "$_id.licenseTypeID",
						"object": mu.APOCond(
							"$_id.cluster",
							bson.M{
								"name": "$_id.cluster",
								"type": "cluster",
							},
							bson.M{
								"name": "$_id.hostname",
								"type": "host",
							},
						),
					},
					"licenseCount": mu.APOMaxAggr(mu.APOCond(
						"$_id.cluster",
						mu.APODivide("$_id.clusterCpu", 2),
						"$licenseCount",
					)),
				}),
				mu.APProject(bson.M{
					"_id":           0,
					"name":          "$_id.object.name",
					"type":          "$_id.object.type",
					"licenseTypeID": "$_id.licenseTypeID",
					"licenseCount":  1,
					"originalCount": "$licenseCount",
				}),
			),
		)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	var out []dto.HostUsingOracleDatabaseLicenses = make([]dto.HostUsingOracleDatabaseLicenses, 0)

	if err := cur.All(context.TODO(), &out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
	}

	return out, nil
}
