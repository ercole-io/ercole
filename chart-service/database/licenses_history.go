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

// Package database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func (md *MongoDatabase) GetLicenseComplianceHistory(start, end time.Time) ([]dto.LicenseComplianceHistory, error) {
	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "history.date", Value: bson.D{{Key: "$gt", Value: start}}},
					{Key: "history.date", Value: bson.D{{Key: "$lt", Value: end}}},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "licenseTypeID", Value: 1},
					{Key: "itemDescription", Value: 1},
					{Key: "metric", Value: 1},
					{Key: "history",
						Value: bson.D{
							{Key: "$filter",
								Value: bson.D{
									{Key: "input", Value: "$history"},
									{Key: "as", Value: "item"},
									{Key: "cond",
										Value: bson.D{
											{Key: "$and",
												Value: bson.A{
													bson.D{
														{Key: "$gt",
															Value: bson.A{
																"$$item.date",
																start,
															},
														},
													},
													bson.D{
														{Key: "$lt",
															Value: bson.A{
																"$$item.date",
																end,
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
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection("database_licenses_history").
		Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var items []dto.LicenseComplianceHistory

	err = cur.All(context.TODO(), &items)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return items, nil
}
