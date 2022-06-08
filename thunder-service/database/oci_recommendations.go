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
	"math"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const OciRecommendation_collection = "oci_recommendations"

func (md *MongoDatabase) AddOciRecommendation(ercoleRecommendation model.OciRecommendation) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendation_collection).
		InsertOne(
			context.TODO(),
			ercoleRecommendation,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) AddOciRecommendations(ercoleRecommendations []model.OciRecommendation) error {
	recToDB := make([]interface{}, len(ercoleRecommendations))

	for i := range ercoleRecommendations {
		recToDB[i] = ercoleRecommendations[i]
	}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendation_collection).
		InsertMany(
			context.TODO(),
			recToDB,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetOciRecommendationsByProfiles(profileIDs []string) ([]model.OciRecommendation, error) {
	ctx := context.TODO()

	inCondition := bson.M{"$in": profileIDs}
	filter := bson.M{"profileID": inCondition}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendation_collection).Find(ctx, filter)
	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	ercoleRecommendations := make([]model.OciRecommendation, 0)
	if err := cur.All(context.TODO(), &ercoleRecommendations); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return ercoleRecommendations, nil
}

func (md *MongoDatabase) GetOciRecommendations(profileIDs []string) ([]model.OciRecommendation, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"seqValue": -1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendation_collection).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	ociRecommendations := make([]model.OciRecommendation, 0)
	if err := cur.All(context.TODO(), &ociRecommendations); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	ociRecommendation1 := make([]model.OciRecommendation, 0)
	if len(ociRecommendations) > 0 {
		inCondition := bson.M{"$in": profileIDs}

		filter := bson.M{"seqValue": ociRecommendations[0].SeqValue, "profileID": inCondition}

		cur1, err1 := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendation_collection).Find(ctx, filter)
		if err1 != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}

		if err := cur1.All(context.TODO(), &ociRecommendation1); err != nil {
			return nil, utils.NewError(err, "DB ERROR")
		}
	}

	return ociRecommendation1, nil
}

func (md *MongoDatabase) DeleteOldOciRecommendations(dateFrom time.Time) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendation_collection).
		DeleteMany(
			context.TODO(),
			bson.M{"createdAt": bson.M{"$lt": dateFrom}},
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetLastOciSeqValue() (uint64, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"seqValue": -1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciRecommendation_collection).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	ociRecommendations := make([]model.OciRecommendation, 0)
	if err := cur.All(context.TODO(), &ociRecommendations); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return math.MaxUint64, utils.NewError(err, "DB ERROR")
	}

	var retVal uint64
	if len(ociRecommendations) == 0 {
		retVal = 0
	} else {
		retVal = ociRecommendations[0].SeqValue
	}

	return retVal, nil
}
