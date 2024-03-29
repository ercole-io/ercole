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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const groupCollection = "groups"

// InsertGroup insert a group into the database
func (md *MongoDatabase) InsertGroup(group model.Group) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(groupCollection).
		InsertOne(
			context.TODO(),
			group,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

// GetGroup return the group specified by group name
func (md *MongoDatabase) GetGroup(name string) (*model.Group, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(groupCollection).
		FindOne(context.TODO(), bson.M{
			"name": name,
		})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.NewError(utils.ErrGroupNotFound, "DB ERROR")
	} else if res.Err() != nil {
		return nil, utils.NewError(res.Err(), "DB ERROR")
	}

	var out model.Group

	if err := res.Decode(&out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return &out, nil
}

func (md *MongoDatabase) GetGroupByTag(tag string) (*model.Group, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(groupCollection).
		FindOne(context.TODO(), bson.M{"tags": tag})
	if res.Err() != nil {
		return nil, res.Err()
	}

	result := &model.Group{}

	if err := res.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateGroup update a group in the database
func (md *MongoDatabase) UpdateGroup(group model.Group) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(groupCollection).
		UpdateOne(
			context.TODO(),
			bson.M{"name": group.Name},
			bson.M{"$set": group},
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrGroupNotFound, "DB ERROR")
	}

	return nil
}

// DeleteGroup delete a group from the database
func (md *MongoDatabase) DeleteGroup(name string) error {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(groupCollection).
		DeleteOne(context.TODO(), bson.M{
			"name": name,
		})
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if res.DeletedCount != 1 {
		return utils.NewError(utils.ErrGroupNotFound, "DB ERROR")
	}

	return nil
}

// GetGroups lists groups
func (md *MongoDatabase) GetGroups() ([]model.Group, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(groupCollection).
		Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	groups := make([]model.Group, 0)

	err = cur.All(context.TODO(), &groups)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return groups, nil
}
