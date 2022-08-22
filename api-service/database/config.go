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

	"github.com/ercole-io/ercole/v2/config"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) FindConfig() (*config.Configuration, error) {
	ctx := context.TODO()

	res := config.Configuration{}
	if err := md.Client.Database(md.Config.Mongodb.DBName).Collection("config").FindOne(ctx, bson.D{}).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (md *MongoDatabase) ChangeConfig(config config.Configuration) error {
	ctx := context.TODO()

	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("config").DeleteMany(ctx, bson.D{}); err != nil {
		return err
	}

	if _, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("config").InsertOne(ctx, config); err != nil {
		return err
	}

	return nil
}
