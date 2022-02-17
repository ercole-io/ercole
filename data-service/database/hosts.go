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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"context"
	"time"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (md *MongoDatabase) DismissHost(hostname string) error {
	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").UpdateOne(context.TODO(), bson.M{
		"hostname":    hostname,
		"dismissedAt": nil,
	}, mu.UOSet(bson.M{
		"dismissedAt": time.Now(),
		"archived":    true,
	})); err != nil {
		return utils.NewError(err, "DB ERROR")
	} else {
		return nil
	}
}

func (md *MongoDatabase) InsertHostData(hostData model.HostDataBE) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").InsertOne(context.TODO(), hostData)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetCurrentHostnames() ([]string, error) {
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"hostname",
		bson.M{
			"dismissedAt": nil,
			"archived":    false,
		})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var hosts []string = make([]string, 0)
	for _, val := range values {
		hosts = append(hosts, val.(string))
	}

	return hosts, nil
}

// FindOldCurrentHostnames return the list of current hosts that haven't sent hostdata after time t
func (md *MongoDatabase) FindOldCurrentHostnames(t time.Time) ([]string, error) {
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"hostname",
		bson.M{
			"dismissedAt": nil,
			"archived":    false,
			"createdAt":   mu.QOLessThan(t),
		})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var hosts []string = make([]string, 0)
	for _, val := range values {
		hosts = append(hosts, val.(string))
	}

	return hosts, nil
}

// FindOldCurrentHosts return the list of current hosts that haven't sent hostdata after time t
func (md *MongoDatabase) FindOldCurrentHostdata(t time.Time) ([]model.HostDataBE, error) {
	filter := bson.M{
		"dismissedAt": nil,
		"archived":    false,
		"createdAt":   mu.QOLessThan(t),
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		Find(context.TODO(), filter)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	//Decode the documents
	hosts := make([]model.HostDataBE, 0)

	for cur.Next(context.TODO()) {
		var host model.HostDataBE

		if cur.Decode(&host) != nil {
			return nil, utils.NewError(err, "Decode ERROR")
		}

		hosts = append(hosts, host)
	}

	//Return it
	return hosts, nil
}

// FindOldArchivedHosts return the list of archived hosts older than t
func (md *MongoDatabase) FindOldArchivedHosts(t time.Time) ([]primitive.ObjectID, error) {
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"_id",
		bson.M{
			"archived":  true,
			"createdAt": mu.QOLessThan(t),
		})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var ids []primitive.ObjectID = make([]primitive.ObjectID, 0)
	for _, val := range values {
		ids = append(ids, val.(primitive.ObjectID))
	}

	return ids, nil
}

func (md *MongoDatabase) DeleteHostData(id primitive.ObjectID) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").DeleteOne(
		context.TODO(),
		bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

// FindMostRecentHostDataOlderThan return the most recest hostdata that is older than t
func (md *MongoDatabase) FindMostRecentHostDataOlderThan(hostname string, t time.Time) (*model.HostDataBE, error) {
	var out model.HostDataBE

	//Find the most recent HostData older than t
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"hostname":  hostname,
				"createdAt": mu.QOLessThan(t),
			}),
			mu.APSort(bson.M{
				"createdAt": -1,
			}),
			mu.APLimit(1),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return nil, nil
	}

	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return &out, nil
}

func (md *MongoDatabase) GetHostnames() ([]string, error) {
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").
		Distinct(
			context.TODO(),
			"hostname",
			bson.M{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var hostnames []string = make([]string, 0, len(values))
	for i := range values {
		hostnames = append(hostnames, values[i].(string))
	}

	return hostnames, nil
}
