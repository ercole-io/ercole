// Copyright (c) 2024 Sorint.lab S.p.A.
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
)

const GcpProfileCollection = "gcp_profiles"

func (md *MongoDatabase) ListGcpProfiles() ([]model.GcpProfile, error) {
	ctx := context.TODO()

	result := make([]model.GcpProfile, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpProfileCollection).Aggregate(ctx, bson.A{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) AddGcpProfile(profile model.GcpProfile) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpProfileCollection).
		InsertOne(context.Background(), profile)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetActiveGcpProfiles() ([]model.GcpProfile, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpProfileCollection).
		Find(context.TODO(), bson.M{"selected": true})
	if err != nil {
		return nil, err
	}

	result := make([]model.GcpProfile, 0)

	if err := cur.All(context.Background(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) SelectGcpProfile(id primitive.ObjectID, selected bool) error {
	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpProfileCollection).
		UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"selected": selected}}); err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateGcpProfile(id primitive.ObjectID, profile model.GcpProfile) error {
	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpProfileCollection).
		UpdateOne(
			context.TODO(),
			bson.M{"_id": id},
			bson.D{{Key: "$set", Value: bson.D{
				primitive.E{Key: "name", Value: profile.Name},
				primitive.E{Key: "privatekey", Value: profile.PrivateKey},
				primitive.E{Key: "clientemail", Value: profile.ClientEmail},
			}}},
		); err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) RemoveGcpProfile(id primitive.ObjectID) error {
	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpProfileCollection).
		DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}
