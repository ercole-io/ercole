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
	"math"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const AwsRDSCollection = "aws_rds"

func (md *MongoDatabase) GetLastAwsRDSSeqValue() (uint64, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"seqValue": -1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsRDSCollection).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	awsRds := make([]model.AwsRDS, 0)
	if err := cur.All(context.TODO(), &awsRds); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	var retVal uint64
	if len(awsRds) == 0 {
		retVal = 0
	} else {
		retVal = awsRds[0].SeqValue
	}

	return retVal, nil
}

func (md *MongoDatabase) GetAwsRDS() ([]model.AwsRDS, error) {
	seqValue, err := md.GetLastAwsRDSSeqValue()
	if err != nil {
		return nil, err
	}

	result := make([]model.AwsRDS, 0)

	filter := bson.M{"seqValue": seqValue}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsRDSCollection).Find(context.Background(), filter)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(context.Background(), &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}
