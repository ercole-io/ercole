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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const exdataCollection = "exadatas"

func (md *MongoDatabase) FindExadataByRackID(rackID string) (*model.OracleExadataInstance, error) {
	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		FindOne(context.TODO(), bson.M{"rackID": rackID})
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

func (md *MongoDatabase) UpdateExadataHostname(rackID, hostname string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		UpdateOne(context.TODO(), bson.M{"rackID": rackID},
			bson.M{"$set": bson.M{"hostname": hostname}})
	if err != nil {
		return err
	}

	return md.updateExadataTime(rackID)
}

func (md *MongoDatabase) PushComponentToExadataInstance(rackID string, component model.OracleExadataComponent) error {
	filter := bson.M{"rackID": rackID}
	update := bson.M{"$push": bson.M{"components": component}}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return md.updateExadataTime(rackID)
}

func (md *MongoDatabase) SetExadataComponent(rackID string, component model.OracleExadataComponent) error {
	filter := bson.M{"rackID": rackID}

	update := bson.M{
		"$set": bson.M{
			"components.$[elem].hostname":          component.Hostname,
			"components.$[elem].hostType":          component.HostType,
			"components.$[elem].cpuEnabled":        component.CPUEnabled,
			"components.$[elem].totalCPU":          component.TotalCPU,
			"components.$[elem].memory":            component.Memory,
			"components.$[elem].imageVersion":      component.ImageVersion,
			"components.$[elem].kernel":            component.Kernel,
			"components.$[elem].model":             component.Model,
			"components.$[elem].fanUsed":           component.FanUsed,
			"components.$[elem].fanTotal":          component.FanTotal,
			"components.$[elem].psuUsed":           component.PsuUsed,
			"components.$[elem].psuTotal":          component.PsuTotal,
			"components.$[elem].msStatus":          component.MsStatus,
			"components.$[elem].rsStatus":          component.RsStatus,
			"components.$[elem].cellServiceStatus": component.CellServiceStatus,
			"components.$[elem].swVersion":         component.SwVersion,
			"components.$[elem].vms":               component.VMs,
			"components.$[elem].storageCells":      component.StorageCells,
		},
	}

	arrayFilters := []interface{}{
		bson.M{"elem.hostname": component.Hostname},
	}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		UpdateOne(
			context.TODO(),
			filter,
			update,
			options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters}),
		)
	if err != nil {
		return err
	}

	return md.updateExadataTime(rackID)
}

func (md *MongoDatabase) UpdateExadataHidden(rackID string, hidden bool) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		UpdateOne(context.TODO(), bson.M{"rackID": rackID},
			bson.M{"$set": bson.M{"hidden": hidden}})
	if err != nil {
		return err
	}

	return md.updateExadataTime(rackID)
}

func (md *MongoDatabase) updateExadataTime(rackID string) error {
	now := md.TimeNow()

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exdataCollection).
		UpdateOne(context.TODO(), bson.M{"rackID": rackID},
			bson.M{"$set": bson.M{"updateAt": now}})
	if err != nil {
		return err
	}

	return nil
}
