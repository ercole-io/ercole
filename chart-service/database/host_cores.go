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
	"time"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func (md *MongoDatabase) GetHostCores(location, environment string, olderThan, newerThan time.Time) ([]dto.HostCores, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"createdAt": bson.M{
					"$gte": newerThan,
					"$lte": olderThan,
				},
			}),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(
				bson.M{
					"_id":   bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}},
					"cores": bson.M{"$sum": "$info.cpuCores"},
				},
			),
			mu.APSort(bson.M{
				"_id": 1,
			}),
			mu.APProject(bson.M{"date": bson.M{"$dateFromString": bson.M{"dateString": "$_id"}}, "cores": 1, "_id": 0}),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var items = make([]dto.HostCores, 0)
	if err := cur.All(context.TODO(), &items); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return items, nil
}
