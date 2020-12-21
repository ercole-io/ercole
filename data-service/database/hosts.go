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
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO return value, not mongo struct
func (md *MongoDatabase) ArchiveHost(hostname string) (*mongo.UpdateResult, utils.AdvancedErrorInterface) {
	if res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").UpdateOne(context.TODO(), bson.M{
		"hostname": hostname,
		"archived": false,
	}, mu.UOSet(bson.M{
		"archived": true,
	})); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	} else {
		return res, nil
	}
}

// TODO return value, not mongo struct
func (md *MongoDatabase) InsertHostData(hostData model.HostDataBE) (*mongo.InsertOneResult, utils.AdvancedErrorInterface) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").InsertOne(context.TODO(), hostData)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return res, nil
}

// FindOldCurrentHosts return the list of current hosts that haven't sent hostdata after time t
func (md *MongoDatabase) FindOldCurrentHosts(t time.Time) ([]string, utils.AdvancedErrorInterface) {
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"hostname",
		bson.M{
			"archived":  false,
			"createdAt": mu.QOLessThan(t),
		})
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	var hosts []string = make([]string, 0)
	for _, val := range values {
		hosts = append(hosts, val.(string))
	}

	return hosts, nil
}

// FindOldArchivedHosts return the list of archived hosts older than t
func (md *MongoDatabase) FindOldArchivedHosts(t time.Time) ([]primitive.ObjectID, utils.AdvancedErrorInterface) {
	values, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Distinct(
		context.TODO(),
		"_id",
		bson.M{
			"archived":  true,
			"createdAt": mu.QOLessThan(t),
		})
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	var ids []primitive.ObjectID = make([]primitive.ObjectID, 0)
	for _, val := range values {
		ids = append(ids, val.(primitive.ObjectID))
	}

	return ids, nil
}

func (md *MongoDatabase) DeleteHostData(id primitive.ObjectID) utils.AdvancedErrorInterface {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").DeleteOne(
		context.TODO(),
		bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return nil
}
