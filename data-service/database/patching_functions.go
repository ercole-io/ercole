// Copyright (c) 2020 Sorint.lab S.p.A.
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
)

// FindPatchingFunction find the the patching function associated to the hostname in the database
func (md *MongoDatabase) FindPatchingFunction(hostname string) (model.PatchingFunction, error) {
	var out model.PatchingFunction

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("patching_functions").Find(context.TODO(), bson.M{
		"hostname": hostname,
	})
	if err != nil {
		return model.PatchingFunction{}, utils.NewError(err, "DB ERROR")
	}

	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return model.PatchingFunction{}, nil
	}

	if err := cur.Decode(&out); err != nil {
		return model.PatchingFunction{}, utils.NewError(err, "DB ERROR")
	}

	return out, nil
}
