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

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
)

const AwsObjectsCollection = "aws_objects"

func (md *MongoDatabase) GetAwsObjectsBySeqValue(seqValue uint64) ([]model.AwsObject, error) {
	ctx := context.TODO()

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsObjectsCollection).
		Aggregate(ctx, bson.A{
			bson.M{"$match": bson.M{
				"seqValue": seqValue,
			}},
		})
	if err != nil {
		return nil, err
	}

	awsObjects := make([]model.AwsObject, 0)
	if err := cur.All(context.TODO(), &awsObjects); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return awsObjects, nil
}
