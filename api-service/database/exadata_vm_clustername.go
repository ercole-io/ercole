// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const exadataVmClusterCollection = "exadata_vm_clusternames"

func (md *MongoDatabase) InsertExadataVmClustername(rackID, hostID, vmname, clustername string) error {
	document := model.OracleExadataVmClustername{
		InstanceRackID: rackID,
		HostID:         hostID,
		VmName:         vmname,
		Clustername:    clustername,
	}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataVmClusterCollection).
		InsertOne(context.Background(), document)
	if err != nil {
		return err
	}

	return md.updateExadataTime(rackID)
}

func (md *MongoDatabase) FindExadataVmClustername(rackID, hostID, vmname string) (*model.OracleExadataVmClustername, error) {
	filter := bson.D{
		{Key: "instancerackid", Value: rackID},
		{Key: "hostid", Value: hostID},
		{Key: "vmname", Value: vmname},
	}

	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataVmClusterCollection).
		FindOne(context.Background(), filter, nil)
	if res.Err() != nil {
		return nil, res.Err()
	}

	result := &model.OracleExadataVmClustername{}

	if err := res.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (md *MongoDatabase) UpdateExadataVmClustername(rackID, hostID, vmname, clustername string) error {
	filter := bson.D{
		{Key: "instancerackid", Value: rackID},
		{Key: "hostid", Value: hostID},
		{Key: "vmname", Value: vmname},
	}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(exadataVmClusterCollection).
		UpdateOne(
			context.Background(),
			filter,
			bson.D{{Key: "$set", Value: bson.D{
				primitive.E{Key: "clustername", Value: clustername},
			}}},
		)
	if err != nil {
		return err
	}

	return nil
}
