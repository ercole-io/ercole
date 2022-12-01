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
		if err := rename_oracle_database_licenses_history_collection(db); err != nil {
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

func rename_oracle_database_licenses_history_collection(db *mongo.Database) error {
	source := "oracle_database_licenses_history"
	dest := "database_licenses_history"
	ctx := context.TODO()

	if cols, err := db.ListCollectionNames(ctx, bson.D{}); err != nil {
		return err
	} else if utils.Contains(cols, source) {
		if err := db.CreateCollection(ctx, dest); err != nil {
			return err
		}

		documents, err := db.Collection(source).Find(ctx, bson.D{})
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}

		args := make([]interface{}, 0)
		if err := documents.All(ctx, &args); err != nil {
			return err
		}

		if _, err := db.Collection(dest).InsertMany(ctx, args); err != nil {
			return err
		}

		if err := db.Collection(source).Drop(ctx); err != nil {
			return err
		}
	}

	return nil
}
