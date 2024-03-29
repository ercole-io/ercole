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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := manageIndexesAddUsers(db); err != nil {
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

func manageIndexesAddUsers(client *mongo.Database) error {
	collectionName := "users"
	indexName := "username_1_provider_1"
	ctx := context.TODO()

	_, err := client.Collection(collectionName).
		Indexes().DropOne(ctx, indexName)
	if err != nil {
		return utils.NewError(err, "Can't drop index:", indexName)
	}

	_, errIndex := client.Collection(collectionName).Indexes().
		CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{
				{Key: "username", Value: 1},
				{Key: "provider", Value: 1},
			},
			Options: (&options.IndexOptions{
				Collation: &options.Collation{Locale: "en", Strength: 2},
			}).SetUnique(true),
		})
	if errIndex != nil {
		return errIndex
	}

	return nil
}
