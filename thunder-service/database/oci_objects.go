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

const OciObject_collection = "oci_objects"

func (md *MongoDatabase) AddOciObjects(objects model.OciObjects) error {
	filter := bson.M{"profileID": objects.ProfileID}
	update := bson.M{"$set": bson.M{"archived": true}}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciObject_collection).UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	_, err1 := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciObject_collection).
		InsertOne(
			context.TODO(),
			objects,
		)
	if err1 != nil {
		return utils.NewError(err1, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetOciObjects() ([]model.OciObjects, error) {
	ctx := context.TODO()

	ABSort := bson.M{"$sort": bson.M{"profileID": 1, "createdAt": -1}}
	ABGroup := bson.M{"$group": bson.M{
		"_id":         bson.M{"profileID": "$profileID"},
		"lastChanged": bson.M{"$first": "$$ROOT"}},
	}

	APprojection := bson.M{"$project": bson.M{"_id": "$lastChanged._id", "profileID": "$lastChanged.profileID", "objects": "$lastChanged.objects", "createdAt": "$lastChanged.createdAt", "archived": "$lastChanged.archived", "error": "$lastChanged.error"}}

	pipeline := []bson.M{ABSort, ABGroup, APprojection}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciObject_collection).Aggregate(ctx, pipeline)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	profiles := make([]model.OciObjects, 0)
	if err := cur.All(context.TODO(), &profiles); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err := cur.Err(); err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return profiles, nil
}

func (md *MongoDatabase) DeleteOldOciObjects(dateFrom time.Time) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(OciObject_collection).
		DeleteMany(
			context.TODO(),
			bson.M{"createdAt": bson.M{"$lt": dateFrom}},
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}
