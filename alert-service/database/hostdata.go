// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"context"
	"time"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// FindHostData find a host data
func (md *MongoDatabase) FindHostData(id primitive.ObjectID) (model.HostDataBE, error) {
	//Find the hostdata
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").FindOne(context.TODO(), bson.M{
		"_id": id,
	})
	if res.Err() != nil {
		return model.HostDataBE{}, utils.NewError(res.Err(), "DB ERROR")
	}

	//Decode the data

	var out model.HostDataBE
	if err := res.Decode(&out); err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DB ERROR")
	}

	//Return it!
	return out, nil
}

//TODO RM?
// FindMostRecentHostDataOlderThan return the most recest hostdata that is older than t
func (md *MongoDatabase) FindMostRecentHostDataOlderThan(hostname string, t time.Time) (model.HostDataBE, error) {
	var out model.HostDataBE

	//Find the most recent HostData older than t
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"hostname":  hostname,
				"createdAt": mu.QOLessThan(t),
			}),
			mu.APSort(bson.M{
				"createdAt": -1,
			}),
			mu.APLimit(1),
		),
	)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DB ERROR")
	}
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return model.HostDataBE{}, nil
	}

	if err := cur.Decode(&out); err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DB ERROR")
	}

	return out, nil
}
