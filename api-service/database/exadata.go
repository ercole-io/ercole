// Copyright (c) 2023 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const exadataCollection = "exadatas"

func (md *MongoDatabase) ListExadataInstances(f dto.GlobalFilter) ([]model.OracleExadataInstance, error) {
	ctx := context.TODO()

	result := make([]model.OracleExadataInstance, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).Aggregate(ctx,
		bson.A{
			Filter(f),
		})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) GetExadataInstance(rackID string) (*model.OracleExadataInstance, error) {
	ctx := context.TODO()

	filter := bson.M{"rackID": rackID}

	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		FindOne(ctx, filter)
	if res.Err() != nil && res.Err() != mongo.ErrNoDocuments {
		return nil, res.Err()
	}

	instance := &model.OracleExadataInstance{}

	if err := res.Decode(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

func (md *MongoDatabase) UpdateExadataInstance(instance model.OracleExadataInstance) error {
	filter := bson.M{"rackID": instance.RackID}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		ReplaceOne(context.TODO(), filter, instance)
	if err != nil {
		return err
	}

	return md.updateExadataTime(instance.RackID)
}

func (md *MongoDatabase) updateExadataTime(rackID string) error {
	now := md.TimeNow()

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		UpdateOne(context.TODO(), bson.M{"rackID": rackID},
			bson.M{"$set": bson.M{"updateAt": now}})
	if err != nil {
		return err
	}

	return nil
}
