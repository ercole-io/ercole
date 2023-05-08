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

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
)

const exdataCollection = "exadatas"

func (md *MongoDatabase) FindExadataByRackID(racID string) (*model.OracleExadataInstance, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		FindOne(context.TODO(), bson.M{"racID": racID})
	if res.Err() != nil {
		return nil, res.Err()
	}

	result := &model.OracleExadataInstance{}

	if err := res.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) AddExadata(exadata model.OracleExadataInstance) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		InsertOne(context.TODO(), exadata)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) UpdateExadata(exadata model.OracleExadataInstance) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		UpdateOne(context.TODO(), bson.M{"rackID": exadata.RackID},
			bson.M{"$set": bson.M{"hostname": exadata.Hostname, "components": exadata.Components}})
	if err != nil {
		return err
	}

	return nil
}
