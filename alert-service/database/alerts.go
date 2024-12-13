// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// InsertAlert insert the alert in the database
func (md *MongoDatabase) InsertAlert(alert model.Alert) (*mongo.InsertOneResult, error) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").InsertOne(context.TODO(), alert)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return res, nil
}

// ExistNoDataAlertByHost return true if the host has associated a new NO_DATA alert
func (md *MongoDatabase) ExistNoDataAlertByHost(hostname string) (bool, error) {
	//Count the number of new NO_DATA alerts associated to the host
	val, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").CountDocuments(context.TODO(), bson.M{
		"alertCode":          model.AlertCodeNoData,
		"alertStatus":        model.AlertStatusNew,
		"otherInfo.hostname": hostname,
	}, &options.CountOptions{
		Limit: utils.Intptr(1),
	})
	if err != nil {
		return false, utils.NewError(err, "DB ERROR")
	}

	//Return true if the count > 0
	return val > 0, nil
}

func (md *MongoDatabase) AckOldAlerts(dueDays int) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	expiredDate := md.TimeNow().AddDate(0, 0, -dueDays)
	filter := bson.M{"date": bson.M{"$lt": expiredDate}}
	update := bson.M{"$set": bson.M{"alertStatus": "ACK"}}

	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").
		UpdateMany(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (md *MongoDatabase) RemoveOldAlerts(dueDays int) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	expiredDate := md.TimeNow().AddDate(0, 0, -dueDays)
	filter := bson.M{
		"date":        bson.M{"$lt": expiredDate},
		"alertStatus": "ACK",
	}

	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").
		DeleteMany(ctx, filter)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (md *MongoDatabase) FindAlertsByDate(startDate, endDate time.Time) ([]model.Alert, error) {
	ctx := context.TODO()

	start := time.Date(
		startDate.Year(),
		startDate.Month(),
		startDate.Day(),
		0, 0, 0, 0, startDate.Location())

	end := time.Date(
		endDate.Year(),
		endDate.Month(),
		endDate.Day(),
		23, 59, 59, 9999, endDate.Location())

	filter := bson.D{
		{Key: "date",
			Value: bson.D{
				{Key: "$gte", Value: start},
				{Key: "$lte", Value: end},
			},
		},
		{Key: "alertStatus", Value: "NEW"},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	alerts := make([]model.Alert, 0)

	if err := cur.All(ctx, &alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}
