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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

func (md *MongoDatabase) ListUsers() ([]model.User, error) {
	ctx := context.TODO()

	result := make([]model.User, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).Aggregate(ctx, bson.A{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) AddUser(user model.User) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) GetUser(username string) (*model.User, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		FindOne(context.TODO(), bson.M{"username": username})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.ErrUserNotFound
	} else if res.Err() != nil {
		return nil, res.Err()
	}

	result := &model.User{}

	if err := res.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) UpdateUserGroups(user model.User) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		UpdateOne(
			context.TODO(),
			bson.M{"username": user.Username},
			bson.D{{Key: "$set", Value: bson.D{
				primitive.E{Key: "groups", Value: user.Groups},
			}}},
		)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) RemoveUser(username string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		DeleteOne(context.TODO(), bson.M{"username": username})
	if err != nil {
		return err
	}

	return nil
}
