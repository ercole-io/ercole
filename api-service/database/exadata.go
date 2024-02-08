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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const exadataCollection = "exadatas"

func (md *MongoDatabase) ListExadataInstances(f dto.GlobalFilter) ([]dto.ExadataInstanceResponse, error) {
	ctx := context.TODO()

	result := make([]dto.ExadataInstanceResponse, 0)

	projection := bson.D{{Key: "rackID", Value: 1}, {Key: "hostname", Value: 1}}

	opts := options.Find().SetProjection(projection)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		Find(ctx, Filter(f), opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) FindExadataInstance(rackID string) (*model.OracleExadataInstance, error) {
	ctx := context.TODO()

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "rackID", Value: rackID}}}},
		bson.D{{Key: "$unwind", Value: "$components"}},
		bson.D{{Key: "$unwind", Value: "$components.vms"}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "exadata_vm_clusternames"},
					{Key: "localField", Value: "components.vms.name"},
					{Key: "foreignField", Value: "vmname"},
					{Key: "as", Value: "matchedDocument"},
				},
			},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "components.vms.clusterName",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.D{
									{Key: "if",
										Value: bson.D{
											{Key: "$eq",
												Value: bson.A{
													bson.D{{Key: "$size", Value: "$matchedDocument"}},
													1,
												},
											},
										},
									},
									{Key: "then",
										Value: bson.D{
											{Key: "$arrayElemAt",
												Value: bson.A{
													"$matchedDocument.clustername",
													0,
												},
											},
										},
									},
									{Key: "else", Value: ""},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id",
						Value: bson.D{
							{Key: "_id", Value: "$_id"},
							{Key: "hostname", Value: "$hostname"},
							{Key: "environment", Value: "$environment"},
							{Key: "location", Value: "$location"},
							{Key: "rackID", Value: "$rackID"},
						},
					},
					{Key: "components",
						Value: bson.D{
							{Key: "$push",
								Value: bson.D{
									{Key: "rackID", Value: "$components.rackID"},
									{Key: "hostType", Value: "$components.hostType"},
									{Key: "hostname", Value: "$components.hostname"},
									{Key: "hostID", Value: "$components.hostID"},
									{Key: "cpuEnabled", Value: "$components.cpuEnabled"},
									{Key: "totalCPU", Value: "$components.totalCPU"},
									{Key: "memory", Value: "$components.memory"},
									{Key: "imageVersion", Value: "$components.imageVersion"},
									{Key: "kernel", Value: "$components.kernel"},
									{Key: "model", Value: "$components.model"},
									{Key: "fanUsed", Value: "$components.fanUsed"},
									{Key: "fanTotal", Value: "$components.fanTotal"},
									{Key: "psuUsed", Value: "$components.psuUsed"},
									{Key: "psuTotal", Value: "$components.psuTotal"},
									{Key: "msStatus", Value: "$components.msStatus"},
									{Key: "rsStatus", Value: "$components.rsStatus"},
									{Key: "cellServiceStatus", Value: "$components.cellServiceStatus"},
									{Key: "swVersion", Value: "$components.swVersion"},
									{Key: "vms",
										Value: bson.A{
											bson.D{
												{Key: "type", Value: "$components.vms.type"},
												{Key: "physicalHost", Value: "$components.vms.physicalHost"},
												{Key: "status", Value: "$components.vms.status"},
												{Key: "name", Value: "$components.vms.name"},
												{Key: "cpuCurrent", Value: "$components.vms.cpuCurrent"},
												{Key: "cpuRestart", Value: "$components.vms.cpuRestart"},
												{Key: "ramCurrent", Value: "$components.vms.ramCurrent"},
												{Key: "ramRestart", Value: "$components.vms.ramRestart"},
												{Key: "cpuOnline", Value: "$components.vms.cpuOnline"},
												{Key: "cpuMaxUsable", Value: "$components.vms.cpuMaxUsable"},
												{Key: "ramOnline", Value: "$components.vms.ramOnline"},
												{Key: "ramMaxUsable", Value: "$components.vms.ramMaxUsable"},
												{Key: "clusterName", Value: "$components.vms.clusterName"},
											},
										},
									},
									{Key: "storageCells", Value: "$components.storageCells"},
									{Key: "clusterNames", Value: "$components.clusterNames"},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: "$_id._id"},
					{Key: "hostname", Value: "$_id.hostname"},
					{Key: "environment", Value: "$_id.environment"},
					{Key: "location", Value: "$_id.location"},
					{Key: "rackID", Value: "$_id.rackID"},
					{Key: "components", Value: "$components"},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]model.OracleExadataInstance, 0)

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	if len(res) > 0 {
		return &res[0], nil
	}

	return nil, mongo.ErrNoDocuments
}

func (md *MongoDatabase) UpdateExadataInstance(instance model.OracleExadataInstance) error {
	filter := bson.M{"rackID": instance.RackID}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		ReplaceOne(context.TODO(), filter, instance)
	if err != nil {
		return err
	}

	return md.updateExadataTime(instance.RackID)
}

func (md *MongoDatabase) updateExadataTime(rackID string) error {
	now := md.TimeNow()

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataCollection).
		UpdateOne(context.TODO(), bson.M{"rackID": rackID},
			bson.M{"$set": bson.M{"updateAt": now}})
	if err != nil {
		return err
	}

	return nil
}
