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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	scenarioCollection      = "scenarios"
	simulatedHostCollection = "simulated_hosts"
)

func (md *MongoDatabase) CreateScenario(scenario *model.Scenario) (*model.Scenario, error) {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(scenarioCollection).InsertOne(context.Background(), scenario)
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		scenario.ID = oid
	}

	return scenario, nil
}

func (md *MongoDatabase) CreateSimulatedHosts(hosts ...model.SimulatedHost) error {
	var docs []any
	for _, h := range hosts {
		docs = append(docs, h)
	}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(simulatedHostCollection).InsertMany(context.Background(), docs)
	if err != nil {
		return err
	}

	return nil
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
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(simulatedHostCollection).
		DeleteOne(context.Background(), bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetScenarios() ([]model.Scenario, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(scenarioCollection).
		Find(context.Background(), bson.D{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	scenarios := make([]model.Scenario, 0)

	err = cur.All(context.Background(), &scenarios)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return scenarios, nil
}

func (md *MongoDatabase) GetScenario(id primitive.ObjectID) (*model.Scenario, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(scenarioCollection).
		FindOne(context.Background(), bson.M{
			"_id": id,
		})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, utils.NewError(utils.ErrScenarioNotFound, "DB ERROR")
	} else if res.Err() != nil {
		return nil, res.Err()
	}

	var out model.Scenario

	if err := res.Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (md *MongoDatabase) RemoveScenario(id primitive.ObjectID) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(scenarioCollection).
		DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
