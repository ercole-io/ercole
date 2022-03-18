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
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const hosts_collection = "hosts"

func (md *MongoDatabase) GetErcoleDatabases() ([]model.ErcoleDatabase, error) {
	ctx := context.TODO()

	opts := options.Find()
	opts.SetProjection(bson.M{"hostname": 1, "info.cpuThreads": 1, "archived": 1, "features.oracle.database.databases.name": 1, "features.oracle.database.databases.uniqueName": 1, "features.oracle.database.databases.work": 1})

	now := md.TimeNow().UTC()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -31)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hosts_collection).Find(ctx, bson.M{"createdAt": bson.M{"$gte": start, "$lt": end}}, opts)

	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	databases := make([]model.ErcoleDatabase, 0)
	err = cur.All(context.TODO(), &databases)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return databases, nil
}

func (md *MongoDatabase) GetErcoleActiveDatabases() ([]model.ErcoleDatabase, error) {
	ctx := context.TODO()

	opts := options.Find()
	opts.SetProjection(bson.M{"hostname": 1, "info.cpuThreads": 1, "features.oracle.database.databases.name": 1, "features.oracle.database.databases.uniqueName": 1, "features.oracle.database.databases.work": 1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hosts_collection).Find(ctx, bson.M{"archived": false, "dismissedAt": nil}, opts)

	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	databases := make([]model.ErcoleDatabase, 0)
	err = cur.All(context.TODO(), &databases)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return databases, nil
}
