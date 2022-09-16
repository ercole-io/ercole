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
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	alert_filter "github.com/ercole-io/ercole/v2/api-service/dto/filter"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

const alertsCollection = "alerts"

func (md *MongoDatabase) SearchAlerts(alertFilter alert_filter.Alert) (*dto.Pagination, error) {
	offset := int64(alertFilter.Filter.Limit * (alertFilter.Filter.Page - 1))
	limit := int64(alertFilter.Filter.Limit)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(alertsCollection).Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APOptionalStage(alertFilter.Status != "", mu.APMatch(bson.M{
				"alertStatus": alertFilter.Status,
			})),
			mu.APOptionalStage(alertFilter.Severity != "", mu.APMatch(bson.M{
				"alertSeverity": alertFilter.Severity,
			})),
			mu.APMatch(bson.M{
				"date": bson.M{
					"$gte": alertFilter.From,
					"$lt":  alertFilter.To,
				},
			}),
			mu.APSearchFilterStage([]interface{}{
				"$description",
				"$alertCode",
				"$alertSeverity",
				"$alertCategory",
				"$otherInfo.Hostname",
				"$otherInfo.Dbname",
				"$otherInfo.Features",
			}, alertFilter.Keywords),
			mu.APSet(bson.M{
				"hostname": "$otherInfo.hostname",
			}),

			mu.APOptionalStage(len(alertFilter.Location) > 0 || len(alertFilter.Environment) > 0,
				mu.APLookupPipeline(
					"hosts",
					bson.M{"hn": "$otherInfo.hostname"},
					"host",
					mu.MAPipeline(
						mu.APMatch(bson.M{
							"$expr":       bson.M{"$eq": bson.A{"$hostname", "$$hn"}},
							"dismissedAt": nil,
							"archived":    false,
						}),
						mu.APProject(bson.M{
							"_id":         0,
							"location":    1,
							"environment": 1,
						}),
					),
				),
			),
			mu.APOptionalStage(len(alertFilter.Location) > 0 || len(alertFilter.Environment) > 0,
				bson.M{
					"$unwind": bson.M{"path": "$host", "preserveNullAndEmptyArrays": true},
				},
			),
			mu.APOptionalStage(len(alertFilter.Location) > 0,
				mu.APMatch(bson.M{
					"$or": bson.A{
						bson.M{"host.location": alertFilter.Location},
						bson.M{"host": bson.M{"$exists": false}},
					},
				}),
			),
			mu.APOptionalStage(len(alertFilter.Environment) > 0,
				mu.APMatch(bson.M{
					"$or": bson.A{
						bson.M{"host.environment": alertFilter.Environment},
						bson.M{"host": bson.M{"$exists": false}},
					},
				})),
			mu.APUnset("host"),

			mu.APOptionalStage(alertFilter.Mode == "aggregated-code-severity", mu.MAPipeline(
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

			mu.APOptionalStage(alertFilter.Mode == "aggregated-category-severity", mu.MAPipeline(
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

			mu.APOptionalSortingStage(alertFilter.SortBy, alertFilter.SortDesc),

			bson.M{
				"$facet": bson.M{
					"items":      bson.A{bson.M{"$skip": offset}, bson.M{"$limit": limit}},
					"totalCount": bson.A{bson.M{"$count": "totalCount"}},
				},
			},

			bson.M{
				"$unwind": bson.M{
					"path": "$totalCount",
				},
			},
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

	var count int32

	lenOut := len(out)
	if lenOut > 0 {
		for _, value := range out[lenOut-1]["totalCount"].(map[string]interface{}) {
			count = value.(int32)
		}

		return dto.ToPagination((out[lenOut-1]["items"]), int(count), alertFilter.Filter.Limit, alertFilter.Filter.Page), nil
	}

	return dto.ToPagination(nil, int(count), alertFilter.Filter.Limit, alertFilter.Filter.Page), nil
}

func (md *MongoDatabase) CountAlertsNODATA(alertsFilter dto.AlertsFilter) (int64, error) {
	data, err := bson.Marshal(alertsFilter)
	if err != nil {
		return 0, err
	}

	var filter map[string]interface{}

	err = bson.Unmarshal(data, &filter)
	if err != nil {
		return 0, err
	}

	if len(filter) < 1 {
		return 0, nil //Do not acknowledge anything
	}

	if v, ok := filter["otherInfo"]; ok {
		if v, ok := v.(map[string]interface{}); ok {
			for k, vv := range v {
				filter["otherInfo."+k] = vv
			}
		}

		delete(filter, "otherInfo")
	}

	ids := alertsFilter.IDs
	if len(ids) >= 1 {
		filter["_id"] = bson.M{"$in": ids}
	}

	filter["alertCode"] = bson.M{"$eq": "NO_DATA"}

	count, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(alertsCollection).
		CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, utils.NewError(err, "DB ERROR")
	}

	return count, nil
}

func (md *MongoDatabase) UpdateAlertsStatus(alertsFilter dto.AlertsFilter, newStatus string) error {
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

	ids := alertsFilter.IDs
	if len(ids) >= 1 {
		filter["_id"] = bson.M{"$in": ids}
	}

	if alertsFilter.AlertCode != nil {
		alertCode := *alertsFilter.AlertCode
		if alertCode == model.AlertStatusDismissed {
			return utils.NewError(errors.New("Invalid status"), "Invalid status")
		}
	} else {
		filter["alertStatus"] = bson.M{"$ne": model.AlertStatusDismissed}
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

func (md *MongoDatabase) RemoveAlertsNODATA(alertsFilter dto.AlertsFilter) error {
	filter := bson.D{{Key: "$and", Value: []interface{}{bson.D{{Key: "otherInfo.hostname", Value: alertsFilter.OtherInfo["hostname"]}}, bson.D{{Key: "alertCode", Value: model.AlertCodeNoData}}}}}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection(alertsCollection).
		DeleteMany(
			context.TODO(),
			filter,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}
