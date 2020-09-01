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
	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertOracleDatabaseAgreement insert a Oracle/Database agreement into the database
func (md *MongoDatabase) InsertOracleDatabaseAgreement(aggreement model.OracleDatabaseAgreement) (*mongo.InsertOneResult, utils.AdvancedErrorInterface) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("agreements_oracle_database").InsertOne(context.TODO(), aggreement)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return res, nil
}

// ListOracleDatabaseAgreements lists the Oracle/Database agreements
func (md *MongoDatabase) ListOracleDatabaseAgreements() ([]apimodel.OracleDatabaseAgreementsFE, utils.AdvancedErrorInterface) {
	var out []apimodel.OracleDatabaseAgreementsFE = make([]apimodel.OracleDatabaseAgreementsFE, 0)

	//Find the matching alerts
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("agreements_oracle_database").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APSet(bson.M{
				"availableCount": -1,
				"licensesCount":  mu.APOCond(mu.APOOr(mu.APOEqual("$metrics", "Processor Perpetual"), mu.APOEqual("$metrics", "Computer Perpetual")), "$count", 0),
				"usersCount":     mu.APOCond(mu.APOEqual("$metrics", "Named User Plus Perpetual"), "$count", 0),
				"hosts": mu.APOMap("$hosts", "hn", bson.M{
					"hostname":                  "$$hn",
					"coveredLicensesCount":      -1,
					"totalCoveredLicensesCount": -1,
				}),
			}),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	if err = cur.All(context.TODO(), &out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
	}
	return out, nil
}
