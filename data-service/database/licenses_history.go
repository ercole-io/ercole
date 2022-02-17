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

// Package database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

//TODO rename this collection with a more generic name
const licensesHistoryCollection = "oracle_database_licenses_history"

func (md *MongoDatabase) HistoricizeLicensesCompliance(licenses []dto.LicenseCompliance) error {
	now := md.TimeNow()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for _, license := range licenses {
		done, err := md.updateLicenseComplianceHistoric(license, today)
		if err != nil {
			return err
		}

		if !done {
			err := md.insertLicenseComplianceHistoric(license, today)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (md *MongoDatabase) updateLicenseComplianceHistoric(license dto.LicenseCompliance, today time.Time) (done bool, err error) {
	filter := bson.M{
		"licenseTypeID": license.LicenseTypeID,
		"history":       bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "date", Value: today}}}},
	}

	if len(license.LicenseTypeID) == 0 {
		filter["itemDescription"] = license.ItemDescription
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "history.$.consumed", Value: license.Consumed},
				{Key: "history.$.covered", Value: license.Covered},
				{Key: "history.$.purchased", Value: license.Purchased},
			},
		}}

	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(licensesHistoryCollection).
		UpdateOne(context.TODO(),
			filter,
			update)
	if err != nil {
		return false, utils.NewError(err, "DB ERROR")
	}

	if res.MatchedCount < 1 {
		return false, nil
	}

	return true, nil
}

func (md *MongoDatabase) insertLicenseComplianceHistoric(license dto.LicenseCompliance, today time.Time) error {
	filter := bson.M{
		"licenseTypeID": license.LicenseTypeID,
	}

	if len(license.LicenseTypeID) == 0 {
		filter["itemDescription"] = license.ItemDescription
	}

	updateOptions := options.Update()
	updateOptions.SetUpsert(true)

	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{
					Key: "history",
					Value: bson.D{
						{Key: "date", Value: today},
						{Key: "consumed", Value: license.Consumed},
						{Key: "covered", Value: license.Covered},
						{Key: "purchased", Value: license.Purchased},
					},
				}}}}

	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(licensesHistoryCollection).
		UpdateMany(context.TODO(),
			filter,
			update,
			updateOptions)

	if (res.ModifiedCount < 1 && res.UpsertedCount < 1) || err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}
