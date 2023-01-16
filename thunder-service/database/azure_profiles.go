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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const AzureProfile_collection = "azure_profiles"

func (md *MongoDatabase) AddAzureProfile(profile model.AzureProfile) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AzureProfile_collection).
		InsertOne(
			context.TODO(),
			profile,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateAzureProfile(profile model.AzureProfile) error {
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

	if profile.ClientSecret == nil {
		delete(p, "clientsecret")
	}

	cur, err = md.Client.Database(md.Config.Mongodb.DBName).Collection(AzureProfile_collection).UpdateOne(context.TODO(), bson.M{"_id": profile.ID}, bson.M{"$set": p})

	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetAzureProfiles(hidePrivateKey bool) ([]model.AzureProfile, error) {
	ctx := context.TODO()

	opts := options.Find()
	if hidePrivateKey {
		opts.SetProjection(bson.M{"clientsecret": 0})
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AzureProfile_collection).Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	profiles := make([]model.AzureProfile, 0)
	if err := cur.All(context.TODO(), &profiles); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return profiles, nil
}

func (md *MongoDatabase) GetMapAzureProfiles() (map[primitive.ObjectID]model.AzureProfile, error) {
	profiles, err := md.GetAzureProfiles(false)

	var retProfiles = make(map[primitive.ObjectID]model.AzureProfile)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	for _, profile := range profiles {
		retProfiles[profile.ID] = profile
	}

	return retProfiles, nil
}

func (md *MongoDatabase) DeleteAzureProfile(id primitive.ObjectID) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AzureProfile_collection).
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

func (md *MongoDatabase) SelectAzureProfile(profileId string, selected bool) error {
	var id primitive.ObjectID

	var err error

	if id, err = primitive.ObjectIDFromHex(profileId); err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AzureProfile_collection).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"selected": selected}})

	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetSelectedAzureProfiles() ([]primitive.ObjectID, error) {
	ctx := context.TODO()

	opts := options.Find()
	filter := bson.M{"selected": true}

	opts.SetProjection(bson.M{"profile": 1})

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(AzureProfile_collection).Find(ctx, filter, opts)
	if err != nil {
		return nil, utils.NewError(cur.Err(), "DB ERROR")
	}

	var selected []profileId

	var selectedProfiles []primitive.ObjectID

	if err := cur.All(context.TODO(), &selected); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	for _, s := range selected {
		selectedProfiles = append(selectedProfiles, s.ID)
	}

	return selectedProfiles, nil
}
