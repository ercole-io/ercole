// Copyright (c) 2025 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package database

import (
	"context"

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) ExistsDR(hostname string) bool {
	filter := bson.M{"archived": false, "isDR": true, "hostname": hostname}

	res := md.Client.Database(md.Config.Mongodb.DBName).
		Collection("hosts").FindOne(context.Background(), filter)

	return res.Err() == nil
}

func (md *MongoDatabase) FindDR(hostname string) (*model.HostDataBE, error) {
	filter := bson.D{
		{Key: "archived", Value: false},
		{Key: "isDR", Value: true},
		{Key: "hostname", Value: hostname},
	}

	res := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").FindOne(context.Background(), filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	dr := model.HostDataBE{}

	if err := res.Decode(&dr); err != nil {
		return nil, err
	}

	return &dr, nil
}
