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
	"math"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const AwsRecommendationCollection = "aws_recommendations"

func (md *MongoDatabase) AddAwsObject(m interface{}, collection string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(collection).InsertOne(context.TODO(), m)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) AddAwsObjects(m []interface{}, collection string) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(collection).InsertMany(context.TODO(), m)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetAwsRecommendationsByProfiles(profileIDs []string) ([]model.AwsRecommendation, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"seqValue": -1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsRecommendationCollection).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	awsRecommendations := make([]model.AwsRecommendation, 0)
	if err := cur.All(context.TODO(), &awsRecommendations); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	awsRecommendation1 := make([]model.AwsRecommendation, 0)

	if len(awsRecommendations) > 0 {
		inCondition := bson.M{"$in": profileIDs}

		filter := bson.M{"seqValue": awsRecommendations[0].SeqValue, "profileID": inCondition}

		cur1, err1 := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsRecommendationCollection).Find(ctx, filter)
		if err1 != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}

		if err := cur1.All(context.TODO(), &awsRecommendation1); err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	}

	return awsRecommendation1, nil
}

func (md *MongoDatabase) GetLastAwsSeqValue() (uint64, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"seqValue": -1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsRecommendationCollection).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	awsRecommendations := make([]model.AwsRecommendation, 0)
	if err := cur.All(context.TODO(), &awsRecommendations); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	var retVal uint64
	if len(awsRecommendations) == 0 {
		retVal = 0
	} else {
		retVal = awsRecommendations[0].SeqValue
	}

	return retVal, nil
}

func (md *MongoDatabase) GetAwsRecommendationsBySeqValue(seqValue uint64) ([]model.AwsRecommendation, error) {
	ctx := context.TODO()

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsRecommendationCollection).Find(ctx, bson.M{"seqValue": seqValue})
	if err != nil {
		return nil, err
	}

	awsRecommendations := make([]model.AwsRecommendation, 0)
	if err := cur.All(context.TODO(), &awsRecommendations); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return awsRecommendations, nil
}
