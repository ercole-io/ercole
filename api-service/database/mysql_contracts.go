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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

const mySQLContractCollection = "mysql_contracts"

func (md *MongoDatabase) AddMySQLContract(contract model.MySQLContract) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLContractCollection).
		InsertOne(
			context.TODO(),
			contract,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) UpdateMySQLContract(contract model.MySQLContract) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLContractCollection).
		ReplaceOne(
			context.TODO(),
			bson.M{"_id": contract.ID},
			contract,
		)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if cur.MatchedCount != 1 {
		return utils.NewError(utils.ErrNotFound, "DB ERROR")
	}

	return nil
}

func (md *MongoDatabase) GetMySQLContracts(locations []string) ([]model.MySQLContract, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLContractCollection).
		Aggregate(context.TODO(), filterExistingLocations(locations))
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	contracts := make([]model.MySQLContract, 0)

	err = cur.All(context.TODO(), &contracts)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return contracts, nil
}

func (md *MongoDatabase) DeleteMySQLContract(id primitive.ObjectID) error {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(mySQLContractCollection).
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
