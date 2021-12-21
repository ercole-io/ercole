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
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

const OciVolumePerformance_collection = "oci_volume_performance"

func (md *MongoDatabase) GetOciVolumePerformances() ([]model.OciVolumePerformance, error) {
	ctx := context.TODO()

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciVolumePerformance_collection).Find(ctx, bson.D{})

	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	volumePerformances := make([]model.OciVolumePerformance, 0)
	err = cur.All(context.TODO(), &volumePerformances)

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}
	return volumePerformances, nil
}

func (md *MongoDatabase) AddOciVolumePerformance(volumePerformance model.OciVolumePerformance) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciVolumePerformance_collection).
		InsertOne(
			context.TODO(),
			volumePerformance,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateOciVolumePerformance(volumePerformance model.OciVolumePerformance) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciVolumePerformance_collection).
		ReplaceOne(
			context.TODO(),
			bson.M{"_id": volumePerformance.ID},
			volumePerformance,
		)

	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) DeleteOciVolumePerformance(id primitive.ObjectID) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciVolumePerformance_collection).
		DeleteOne(
			context.TODO(),
			bson.M{"_id": id},
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.DeletedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetOciVolumePerformance(vpu int, size int) (*model.OciVolumePerformance, error) {
	var valRet model.OciVolumePerformance

	ctx := context.TODO()

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciVolumePerformance_collection).Find(ctx, bson.D{})

	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	volumePerformances := make([]model.OciVolumePerformance, 0)
	err = cur.All(context.TODO(), &volumePerformances)

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var perfTmp model.OciPerformance
	var valTmp model.OciPerfValues
out:
	for _, volPerf := range volumePerformances {
		if volPerf.Vpu == vpu {
			for _, perf := range volPerf.Performances {
				if size == perf.Size {
					valRet.ID = volPerf.ID
					perfTmp.Size = perf.Size
					perfTmp.Values = perf.Values
					valRet.Performances = append(valRet.Performances, perfTmp)
					valRet.Vpu = vpu
					break out
				} else if size > perf.Size {
					valTmp.MaxThroughput = float64(perf.Values.MaxThroughput) * float64(size) / float64(perf.Size)
					valTmp.MaxIOPS = perf.Values.MaxIOPS * size / perf.Size
					perfTmp.Size = size
					perfTmp.Values = valTmp
					valRet.ID = volPerf.ID
					valRet.Performances = append(valRet.Performances, perfTmp)
					valRet.Vpu = vpu
					break out
				}
			}

		}

	}
	return &valRet, nil
}
