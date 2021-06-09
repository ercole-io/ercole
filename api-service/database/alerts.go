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
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const alertsCollection = "alerts"

func (md *MongoDatabase) SearchAlerts(mode string, keywords []string, sortBy string, sortDesc bool,
	page, pageSize int, location, environment, severity, status string, from, to time.Time,
) ([]map[string]interface{}, error) {

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(alertsCollection).Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APOptionalStage(status != "", mu.APMatch(bson.M{
				"alertStatus": status,
			})),
			mu.APOptionalStage(severity != "", mu.APMatch(bson.M{
				"alertSeverity": severity,
			})),
			mu.APMatch(bson.M{
				"date": bson.M{
					"$gte": from,
					"$lt":  to,
				},
			}),
			mu.APSearchFilterStage([]interface{}{
				"$description",
				"$alertCode",
				"$alertSeverity",
				"$otherInfo.Hostname",
				"$otherInfo.Dbname",
				"$otherInfo.Features",
			}, keywords),
			mu.APSet(bson.M{
				"hostname": "$otherInfo.hostname",
			}),

			mu.APOptionalStage(len(location) > 0 || len(environment) > 0,
				mu.APLookupPipeline(
					"hosts",
					bson.M{"hn": "$otherInfo.hostname"},
					"host",
					mu.MAPipeline(
						mu.APMatch(bson.M{
							"$expr":    bson.M{"$eq": bson.A{"$hostname", "$$hn"}},
							"archived": false,
						}),
						mu.APProject(bson.M{
							"_id":         0,
							"location":    1,
							"environment": 1,
						}),
					),
				),
			),
			mu.APOptionalStage(len(location) > 0 || len(environment) > 0,
				bson.M{
					"$unwind": bson.M{"path": "$host", "preserveNullAndEmptyArrays": true},
				},
			),
			mu.APOptionalStage(len(location) > 0,
				mu.APMatch(bson.M{
					"$or": bson.A{
						bson.M{"host.location": location},
						bson.M{"host": bson.M{"$exists": false}},
					},
				}),
			),
			mu.APOptionalStage(len(environment) > 0,
				mu.APMatch(bson.M{
					"$or": bson.A{
						bson.M{"host.environment": environment},
						bson.M{"host": bson.M{"$exists": false}},
					},
				})),
			mu.APUnset("host"),

			mu.APOptionalStage(mode == "aggregated-code-severity", mu.MAPipeline(
				mu.APGroup(bson.M{
					"_id": bson.M{
						"code":     "$alertCode",
						"severity": "$alertSeverity",
						"category": "$alertCategory",
					},
					"count": mu.APOSum(1),
					"oldestAlert": bson.M{
						"$min": "$date",
					},
					"affectedHosts": bson.M{
						"$addToSet": "$hostname",
					},
				}),
				mu.APProject(bson.M{
					"_id":           false,
					"category":      "$_id.category",
					"code":          "$_id.code",
					"severity":      "$_id.severity",
					"count":         true,
					"affectedHosts": mu.APOSize("$affectedHosts"),
					"oldestAlert":   true,
				}),
			)),

			mu.APOptionalStage(mode == "aggregated-category-severity", mu.MAPipeline(
				mu.APGroup(bson.M{
					"_id": bson.M{
						"severity": "$alertSeverity",
						"category": "$alertCategory",
					},
					"count": mu.APOSum(1),
					"oldestAlert": bson.M{
						"$min": "$date",
					},
					"affectedHosts": bson.M{
						"$addToSet": "$hostname",
					},
				}),
				mu.APProject(bson.M{
					"_id":           false,
					"category":      "$_id.category",
					"severity":      "$_id.severity",
					"count":         true,
					"affectedHosts": mu.APOSize("$affectedHosts"),
					"oldestAlert":   true,
				}),
			)),

			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var out []map[string]interface{} = make([]map[string]interface{}, 0)

	for cur.Next(context.TODO()) {
		var item map[string]interface{}

		if cur.Decode(&item) != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}

		out = append(out, item)
	}

	return out, nil
}

func (md *MongoDatabase) UpdateAlertsStatus(ids []primitive.ObjectID, newStatus string) error {
	bsonIds := bson.A{}
	for _, id := range ids {
		bsonIds = append(bsonIds, bson.M{"_id": id})
	}
	filter := bson.M{"$or": bsonIds}

	count, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(alertsCollection).
		CountDocuments(context.TODO(), filter)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}
	if count != int64(len(ids)) {
		return utils.ErrAlertNotFound
	}

	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(alertsCollection).
		UpdateMany(context.TODO(),
			filter,
			mu.UOSet(bson.M{
				"alertStatus": newStatus,
			}))
	if err != nil || res.MatchedCount != int64(len(ids)) {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateAlertsStatusByFilter(alertsFilter dto.AlertsFilter, newStatus string) error {
	data, err := bson.Marshal(alertsFilter)
	if err != nil {
		return err
	}

	var filter map[string]interface{}
	err = bson.Unmarshal(data, &filter)
	if err != nil {
		return err
	}
	if len(filter) < 1 {
		return nil //Do not acknowledge anything
	}

	if v, ok := filter["otherInfo"]; ok {
		if v, ok := v.(map[string]interface{}); ok {
			for k, vv := range v {
				filter["otherInfo."+k] = vv
			}
		}

		delete(filter, "otherInfo")
	}

	_, err = md.Client.Database(md.Config.Mongodb.DBName).
		Collection(alertsCollection).
		UpdateMany(
			context.TODO(),
			filter,
			bson.M{"$set": bson.M{
				"alertStatus": newStatus,
			}},
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}
