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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

const oracleDbAgreementsCollection = "oracle_database_agreements"

// InsertOracleDatabaseAgreement insert an Oracle/Database agreement into the database
func (md *MongoDatabase) InsertOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		InsertOne(context.TODO(), agreement)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

// GetOracleDatabaseAgreement return the agreement specified by id
func (md *MongoDatabase) GetOracleDatabaseAgreement(id primitive.ObjectID) (*model.OracleDatabaseAgreement, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		FindOne(context.TODO(), bson.M{
			"_id": id,
		})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.ErrOracleDatabaseAgreementNotFound
	} else if res.Err() != nil {
		return nil, utils.NewError(res.Err(), "DB ERROR")
	}

	var out model.OracleDatabaseAgreement

	if err := res.Decode(&out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return &out, nil
}

// UpdateOracleDatabaseAgreement update an Oracle/Database agreement in the database
func (md *MongoDatabase) UpdateOracleDatabaseAgreement(agreement model.OracleDatabaseAgreement) error {
	result, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		ReplaceOne(context.TODO(), bson.M{
			"_id": agreement.ID,
		}, agreement)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if result.MatchedCount != 1 {
		return utils.ErrOracleDatabaseAgreementNotFound
	}

	return nil
}

// RemoveOracleDatabaseAgreement remove an Oracle/Database agreement from the database
func (md *MongoDatabase) RemoveOracleDatabaseAgreement(id primitive.ObjectID) error {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		DeleteOne(context.TODO(), bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if res.DeletedCount == 0 {
		return utils.ErrOracleDatabaseAgreementNotFound
	}

	return nil
}

// ListOracleDatabaseAgreements lists the Oracle/Database agreements
func (md *MongoDatabase) ListOracleDatabaseAgreements() ([]dto.OracleDatabaseAgreementFE, error) {
	var out []dto.OracleDatabaseAgreementFE = make([]dto.OracleDatabaseAgreementFE, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbAgreementsCollection).
		Aggregate(
			context.TODO(),
			mu.MAPipeline(
				mu.APLookupSimple("oracle_database_license_types", "licenseTypeID", "_id", "licenseType"),
				mu.APUnwind("$licenseType"),
				mu.APSet(bson.M{
					"itemDescription": "$licenseType.itemDescription",
					"metric":          "$licenseType.metric",
					"hosts": mu.APOMap("$hosts", "hn", bson.M{
						"hostname": "$$hn",
					}),

					//TODO And other licenses types?
					"licensesPerCore": mu.APOCond(
						mu.APOOr(
							mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricProcessorPerpetual),
							mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricComputerPerpetual)),
						"$count",
						0),
					"licensesPerUser": mu.APOCond(
						mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricNamedUserPlusPerpetual), "$count", 0),
					"availableLicensesPerCore": mu.APOCond(
						mu.APOOr(
							mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricProcessorPerpetual),
							mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricComputerPerpetual)),
						"$count",
						0),
					"availableLicensesPerUser": mu.APOCond(
						mu.APOEqual("$licenseType.metric", model.LicenseTypeMetricNamedUserPlusPerpetual), "$count", 0),
				}),
				mu.APUnset("licenseType"),
			),
		)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err = cur.All(context.TODO(), &out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}
