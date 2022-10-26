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

const roleCollection = "roles"

// GetRole return the role specified by role name
func (md *MongoDatabase) GetRole(name string) (*model.Role, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(roleCollection).
		FindOne(context.TODO(), bson.M{
			"name": name,
		})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.NewError(utils.ErrRoleNotFound, "DB ERROR")
	} else if res.Err() != nil {
		return nil, utils.NewError(res.Err(), "DB ERROR")
	}

	var out model.Role

	if err := res.Decode(&out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return &out, nil
}

// GetRoles lists roles
func (md *MongoDatabase) GetRoles() ([]model.Role, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(roleCollection).
		Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	roles := make([]model.Role, 0)

	err = cur.All(context.TODO(), &roles)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return roles, nil
}

func (md *MongoDatabase) AddRole(role model.Role) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(roleCollection).InsertOne(context.TODO(), role)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) UpdateRole(name string, documents bson.D) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(roleCollection).
		UpdateOne(
			context.TODO(),
			bson.M{"name": name},
			bson.D{{Key: "$set", Value: documents}},
		)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) RemoveRole(roleName string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(roleCollection).
		DeleteOne(context.TODO(), bson.M{"name": roleName})
	if err != nil {
		return err
	}

	return nil
}
