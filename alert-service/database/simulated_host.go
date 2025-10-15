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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

const (
	simulatedHostCollection = "simulated_hosts"
	hostCollection          = "hosts"
)

func (md *MongoDatabase) GetSimulatedHosts() ([]model.SimulatedHost, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(simulatedHostCollection).
		Find(context.Background(), bson.D{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	simulatedHosts := make([]model.SimulatedHost, 0)

	err = cur.All(context.Background(), &simulatedHosts)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return simulatedHosts, nil
}

func (md *MongoDatabase) UpdateHostCores(hostname string, cores int) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).
		UpdateOne(
			context.Background(),
			bson.M{"hostname": hostname, "archived": false},
			bson.M{
				"$set": bson.M{
					"info.cpuCores": cores,
				},
			},
		)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDatabase) RemoveSimulatedHost(id primitive.ObjectID) error {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(simulatedHostCollection).
		DeleteOne(context.TODO(), bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if res.DeletedCount != 1 {
		return utils.NewError(utils.ErrGroupNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateLicenseCount(hostname string, licenseCount int) error {
	filter := bson.M{
		"archived": false,
		"hostname": hostname,
	}

	update := bson.M{
		"$set": bson.M{
			"features.oracle.database.databases.$[].licenses.$[].count": licenseCount,
		},
	}

	result, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(hostCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return utils.ErrLicenseNotFound
	}

	return nil
}
