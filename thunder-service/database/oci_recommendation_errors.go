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
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const OciRecommendationError_collection = "oci_recommendation_errors"

func (md *MongoDatabase) AddOciRecommendationError(ociRecommendationError model.OciRecommendationError) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendationError_collection).
		InsertOne(
			context.TODO(),
			ociRecommendationError,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) AddOciRecommendationErrors(ociRecommendationErrors []model.OciRecommendationError) error {
	recToDB := make([]interface{}, len(ociRecommendationErrors))

	for i := range ociRecommendationErrors {
		recToDB[i] = ociRecommendationErrors[i]
	}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendationError_collection).
		InsertMany(
			context.TODO(),
			recToDB,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetOciRecommendationErrors(seqNum uint64) ([]model.OciRecommendationError, error) {
	ociRecommendationErrors := make([]model.OciRecommendationError, 0)
	ctx := context.TODO()

	filter := bson.M{"seqValue": seqNum}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendationError_collection).Find(ctx, filter)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.All(context.TODO(), &ociRecommendationErrors); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return ociRecommendationErrors, nil
}

func (md *MongoDatabase) GetOciRecommendationErrorsByProfiles(profileIDs []string) ([]model.OciRecommendationError, error) {
	ctx := context.TODO()

	inCondition := bson.M{"$in": profileIDs}
	filter := bson.M{"profileID": inCondition}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendationError_collection).Find(ctx, filter)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	ociRecommendationErrors := make([]model.OciRecommendationError, 0)
	if err := cur.All(context.TODO(), &ociRecommendationErrors); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return ociRecommendationErrors, nil
}

func (md *MongoDatabase) DeleteOldOciRecommendationErrors(dateFrom time.Time) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendationError_collection).
		DeleteMany(
			context.TODO(),
			bson.M{"createdAt": bson.M{"$lt": dateFrom}},
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}
