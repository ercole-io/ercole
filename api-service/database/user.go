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

func (md *MongoDatabase) GetUser(username string, provider string) (*model.User, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		FindOne(context.TODO(), bson.M{"username": username, "provider": provider})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.ErrInvalidUser
	} else if res.Err() != nil {
		return nil, res.Err()
	}

	result := &model.User{}

	if err := res.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) UpdateUserGroups(username string, provider string, groups []string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		UpdateOne(
			context.TODO(),
			bson.M{"username": username, "provider": provider},
			bson.D{{Key: "$set", Value: bson.D{
				primitive.E{Key: "groups", Value: groups},
			}}},
		)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) UpdateUserLastLogin(user model.User) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		UpdateOne(
			context.TODO(),
			bson.M{
				"username": user.Username,
				"provider": user.Provider,
			},
			bson.M{"$set": bson.M{
				"lastLogin": user.LastLogin,
			}},
		)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) RemoveUser(username string, provider string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		DeleteOne(context.TODO(), bson.M{"username": username, "provider": provider})
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) UpdatePassword(username string, password string, salt string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(userCollection).
		UpdateOne(
			context.TODO(),
			bson.M{"username": username},
			bson.D{{Key: "$set", Value: bson.D{
				primitive.E{Key: "password", Value: password},
				primitive.E{Key: "salt", Value: salt},
			}}},
		)
	if err != nil {
		return err
	}

	return nil
}
