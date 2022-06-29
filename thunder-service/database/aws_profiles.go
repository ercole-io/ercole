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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const AwsProfile_collection = "aws_profiles"

func (md *MongoDatabase) AddAwsProfile(profile model.AwsProfile) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsProfile_collection).
		InsertOne(
			context.TODO(),
			profile,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateAwsProfile(profile model.AwsProfile) error {
	var err error

	var cur *mongo.UpdateResult

	var p bson.M

	data, err := bson.Marshal(profile)
	if err != nil {
		return utils.NewError(err, "Unable to mashal profile")
	}

	err = bson.Unmarshal(data, &p)
	if err != nil {
		return utils.NewError(err, "Unable to unmarshal profile")
	}

	if profile.SecretAccessKey == nil {
		delete(p, "secretaccesskey")
	}

	cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsProfile_collection).UpdateOne(context.TODO(), bson.M{"_id": profile.ID}, bson.M{"$set": p})

	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetAwsProfiles(hidePrivateKey bool) ([]model.AwsProfile, error) {
	ctx := context.TODO()

	opts := options.Find()
	if hidePrivateKey {
		opts.SetProjection(bson.M{"secretaccesskey": 0})
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsProfile_collection).Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	profiles := make([]model.AwsProfile, 0)
	if err := cur.All(context.TODO(), &profiles); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return profiles, nil
}

func (md *MongoDatabase) GetMapAwsProfiles() (map[primitive.ObjectID]model.AwsProfile, error) {
	profiles, err := md.GetAwsProfiles(false)

	var retProfiles = make(map[primitive.ObjectID]model.AwsProfile)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	for _, profile := range profiles {
		retProfiles[profile.ID] = profile
	}

	return retProfiles, nil
}

func (md *MongoDatabase) DeleteAwsProfile(id primitive.ObjectID) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsProfile_collection).
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

func (md *MongoDatabase) SelectAwsProfile(profileId string, selected bool) error {
	var id primitive.ObjectID

	var err error

	if id, err = primitive.ObjectIDFromHex(profileId); err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsProfile_collection).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"selected": selected}})

	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetSelectedAwsProfiles() ([]string, error) {
	ctx := context.TODO()

	opts := options.Find()
	filter := bson.M{"selected": true}

	opts.SetProjection(bson.M{"profile": 1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AwsProfile_collection).Find(ctx, filter, opts)
	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	var selected []profileId

	var selectedProfiles []string

	if err := cur.All(context.TODO(), &selected); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	for _, s := range selected {
		selectedProfiles = append(selectedProfiles, s.ID.Hex())
	}

	return selectedProfiles, nil
}
