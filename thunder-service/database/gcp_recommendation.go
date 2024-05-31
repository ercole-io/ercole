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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	GcpRecommendationCollection = "gcp_recommendations"
	GcpErrorCollection          = "gcp_errors"
)

func (md *MongoDatabase) GetLastGcpSeqValue() (uint64, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"seqValue": -1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpRecommendationCollection).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	gcprecommendation := make([]model.GcpRecommendation, 0)
	if err := cur.All(ctx, &gcprecommendation); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	var res uint64

	if len(gcprecommendation) == 0 {
		res = 0
	} else {
		res = gcprecommendation[0].SeqValue
	}

	return res, nil
}

func (md *MongoDatabase) ListGcpRecommendationsByProfiles(profileIDs []primitive.ObjectID) ([]model.GcpRecommendation, error) {
	ctx := context.Background()

	seqValue, err := md.GetLastGcpSeqValue()
	if err != nil {
		return nil, err
	}

	result := make([]model.GcpRecommendation, 0)

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "seqValue", Value: seqValue},
					{Key: "profileID",
						Value: bson.D{
							{Key: "$in", Value: profileIDs},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "gcp_profiles"},
					{Key: "localField", Value: "profileID"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "profile"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "seqValue", Value: 1},
					{Key: "createdAt", Value: 1},
					{Key: "profileID", Value: 1},
					{Key: "instanceID", Value: 1},
					{Key: "category", Value: 1},
					{Key: "suggestion", Value: 1},
					{Key: "projectID", Value: 1},
					{Key: "projectName", Value: 1},
					{Key: "objectType", Value: 1},
					{Key: "details", Value: 1},
					{Key: "profileName", Value: "$profile.name"},
				},
			},
		},
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpRecommendationCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return result, nil
}

func (md *MongoDatabase) AddGcpRecommendation(gcprecommendation interface{}) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpRecommendationCollection).
		InsertOne(context.TODO(), gcprecommendation)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) AddGcpError(gcperror interface{}) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(GcpErrorCollection).
		InsertOne(context.Background(), gcperror)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}
