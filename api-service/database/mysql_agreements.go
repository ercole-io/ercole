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
)

const mySQLAgreementCollection = "mysql_agreements"

func (md *MongoDatabase) AddMySQLAgreement(agreement model.MySQLAgreement) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLAgreementCollection).
		InsertOne(
			context.TODO(),
			agreement,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateMySQLAgreement(agreement model.MySQLAgreement) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLAgreementCollection).
		ReplaceOne(
			context.TODO(),
			bson.M{"_id": agreement.ID},
			agreement,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetMySQLAgreements() ([]model.MySQLAgreement, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLAgreementCollection).
		Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	agreements := make([]model.MySQLAgreement, 0)
	err = cur.All(context.TODO(), &agreements)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return agreements, nil
}

func (md *MongoDatabase) DeleteMySQLAgreement(id primitive.ObjectID) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLAgreementCollection).
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
