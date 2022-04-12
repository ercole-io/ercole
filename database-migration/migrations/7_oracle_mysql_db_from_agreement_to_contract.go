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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package migrations

import (
	"context"
	"fmt"

	"github.com/ercole-io/ercole/v2/utils"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := rename_oracle_agreements_collection_field(db); err != nil {
			return err
		}
		if err := rename_mysql_agreements_collection_field(db); err != nil {
			return err
		}

		return nil

	}, func(db *mongo.Database) error {
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func rename_oracle_agreements_collection_field(db *mongo.Database) error {
	collectionFrom := "oracle_database_agreements"
	collectionTo := "oracle_database_contracts"
	ctx := context.TODO()
	filter := bson.M{}
	update := bson.M{"$rename": bson.M{"agreementID": "contractID"}}

	names, errColl := db.ListCollectionNames(ctx, filter)
	if errColl != nil {
		return errColl
	}

	for _, name := range names {
		if name == collectionFrom {
			coll := db.Collection(collectionFrom)
			if coll == nil {
				return nil
			}

			docs, err := coll.Find(ctx, bson.D{})
			if err != nil {
				return err
			}

			errDest := db.CreateCollection(ctx, collectionTo)
			if errDest != nil {
				return errDest
			}

			if docs.RemainingBatchLength() > 0 {
				var agrs []interface{}
				if err := docs.All(ctx, &agrs); err != nil {
					return utils.NewError(err, "Can't decode cursor")
				}

				if _, err := db.Collection(collectionTo).InsertMany(ctx, agrs); err != nil {
					return utils.NewError(err, "Can't insert all contracts")
				}
			}

			_, errField := db.Collection(collectionTo).UpdateMany(ctx, filter, update)
			if errField != nil {
				return utils.NewError(err, "Can't rename from agreementID to contractID", collectionTo)
			}

			errDrop := db.Collection(collectionFrom).Drop(ctx)
			if errDrop != nil {
				return utils.NewError(err, "Can't drop collection", collectionFrom)
			}
		}
	}

	return nil
}

func rename_mysql_agreements_collection_field(db *mongo.Database) error {
	collectionFrom := "mysql_agreements"
	collectionTo := "mysql_contracts"
	ctx := context.TODO()
	filter := bson.M{}
	update := bson.M{"$rename": bson.M{"agreementID": "contractID"}}

	names, errColl := db.ListCollectionNames(ctx, filter)
	if errColl != nil {
		return errColl
	}

	for _, name := range names {
		if name == collectionFrom {
			coll := db.Collection(collectionFrom)
			if coll == nil {
				return nil
			}

			docs, err := coll.Find(ctx, bson.D{})
			if err != nil {
				return err
			}

			errDest := db.CreateCollection(ctx, collectionTo)
			if errDest != nil {
				return errDest
			}

			if docs.RemainingBatchLength() > 0 {
				var agrs []interface{}
				if err := docs.All(ctx, &agrs); err != nil {
					return utils.NewError(err, "Can't decode cursor")
				}

				if _, err := db.Collection(collectionTo).InsertMany(ctx, agrs); err != nil {
					return utils.NewError(err, "Can't insert all contracts")
				}
			}

			_, errField := db.Collection(collectionTo).UpdateMany(ctx, filter, update)
			if errField != nil {
				return utils.NewError(err, "Can't rename from agreementID to contractID", collectionTo)
			}

			errDrop := db.Collection(collectionFrom).Drop(ctx)
			if errDrop != nil {
				return utils.NewError(err, "Can't drop collection", collectionFrom)
			}
		}
	}

	return nil
}
