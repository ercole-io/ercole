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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const nodesCollection = "nodes"

func (md *MongoDatabase) GetNodesByRoles(roles []string) ([]model.Node, error) {
	ctx := context.TODO()

	result := make([]model.Node, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(nodesCollection).
		Aggregate(ctx, bson.A{bson.M{"$match": bson.M{"roles": bson.M{"$in": roles}}}})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) GetNodeByName(name string) (*model.Node, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(nodesCollection).
		FindOne(context.TODO(), bson.M{"name": name})
	if res.Err() != nil {
		return nil, res.Err()
	}

	result := &model.Node{}

	if err := res.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) AddNode(node model.Node) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(nodesCollection).InsertOne(context.TODO(), node)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) UpdateNode(node model.Node) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(nodesCollection).
		UpdateOne(
			context.TODO(),
			bson.M{"name": node.Name},
			bson.D{{Key: "$set", Value: bson.D{
				primitive.E{Key: "roles", Value: node.Roles},
				primitive.E{Key: "parent", Value: node.Parent},
			}}},
		)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) RemoveNode(name string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(nodesCollection).
		DeleteOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		return err
	}

	return nil
}
