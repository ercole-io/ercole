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
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const sqlServerDbContractsCollection = "sql_server_database_contracts"

// InsertSqlServerDatabaseContract insert an SqlServer/Database contract into the database
func (md *MongoDatabase) InsertSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) error {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(sqlServerDbContractsCollection).
		InsertOne(context.TODO(), contract)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	return nil
}

// // GetOracleDatabaseContract return the contract specified by id
// func (md *MongoDatabase) GetOracleDatabaseContract(id primitive.ObjectID) (*model.OracleDatabaseContract, error) {
// 	res := md.Client.Database(md.Config.Mongodb.DBName).Collection(oracleDbContractsCollection).
// 		FindOne(context.TODO(), bson.M{
// 			"_id": id,
// 		})
// 	if res.Err() == mongo.ErrNoDocuments {
// 		return nil, utils.ErrOracleDatabaseContractNotFound
// 	} else if res.Err() != nil {
// 		return nil, utils.NewError(res.Err(), "DB ERROR")
// 	}

// 	var out model.OracleDatabaseContract

// 	if err := res.Decode(&out); err != nil {
// 		return nil, utils.NewError(err, "Decode ERROR")
// 	}

// 	return &out, nil
// }

func (md *MongoDatabase) UpdateSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) error {
	result, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(sqlServerDbContractsCollection).
		ReplaceOne(context.TODO(), bson.M{
			"_id": contract.ID,
		}, contract)
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if result.MatchedCount != 1 {
		return utils.ErrContractNotFound
	}

	return nil
}

func (md *MongoDatabase) RemoveSqlServerDatabaseContract(id primitive.ObjectID) error {
	res, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(sqlServerDbContractsCollection).
		DeleteOne(context.TODO(), bson.M{
			"_id": id,
		})
	if err != nil {
		return utils.NewError(err, "DB ERROR")
	}

	if res.DeletedCount == 0 {
		return utils.ErrContractNotFound
	}

	return nil
}

func (md *MongoDatabase) ListSqlServerDatabaseContracts() ([]model.SqlServerDatabaseContract, error) {
	ctx := context.TODO()
	out := make([]model.SqlServerDatabaseContract, 0)

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection(sqlServerDbContractsCollection).
		Find(ctx, bson.M{})
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if err = cur.All(ctx, &out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}
