// Copyright (c) 2021 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/utils"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	collection := "oracle_database_agreements"

	migrate.Register(func(db *mongo.Database) error {
		ctx := context.TODO()

		_, err := db.Collection(collection).
			Indexes().DropOne(ctx, "agreementID_1")
		if err != nil {
			return utils.NewError(err, "Can't drop agreementID_1 index")
		}

		if _, err := db.Collection(collection).
			Indexes().
			CreateMany(context.TODO(),
				[]mongo.IndexModel{
					{

						Keys: bson.D{
							{Key: "agreementID", Value: 1},
							{Key: "csi", Value: 1},
						},
						Options: options.Index().SetUnique(true),
					},
					{
						Keys: bson.D{
							{Key: "licenseTypes._id", Value: 1},
						},
						Options: options.Index().SetUnique(true),
					}},
			); err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		ctx := context.TODO()

		_, err := db.Collection(collection).
			Indexes().DropOne(ctx, "agreementID_1_csi_1")
		if err != nil {
			return utils.NewError(err, "Can't drop agreementID_1_csi_1 index")
		}

		if _, err := db.Collection(collection).
			Indexes().
			CreateMany(context.TODO(),
				[]mongo.IndexModel{
					{

						Keys: bson.D{
							{Key: "agreementID", Value: 1},
						},
						Options: options.Index().SetUnique(true),
					},
					{
						Keys: bson.D{
							{Key: "licenseTypes._id", Value: 1},
						},
						Options: options.Index().SetUnique(true),
					}},
			); err != nil {
			return err
		}

		return nil
	})
}
