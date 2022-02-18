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
	"errors"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ercole-io/ercole/v2/utils"
)

func init() {
	err := migrate.Register(unwind_oracle_agreements, group_oracle_agreementy_by_license_type_id)

	if err != nil {
		panic(err)
	}
}

func unwind_oracle_agreements(db *mongo.Database) error {
	collection := "oracle_database_agreements"
	ctx := context.TODO()

	for _, index := range []string{"agreementID_1", "licenseTypes._id_1"} {
		_, err := db.Collection(collection).
			Indexes().DropOne(ctx, index)
		if err != nil {
			return utils.NewError(err, "Can't drop index:", index)
		}
	}

	cursor, err := db.Collection(collection).
		Aggregate(ctx,
			bson.A{
				bson.M{
					"$unwind": bson.M{"path": "$licenseTypes"},
				},
				bson.M{
					"$replaceRoot": bson.M{
						"newRoot": bson.M{
							"$mergeObjects": bson.A{
								bson.M{
									"agreementID": "$agreementID",
									"csi":         "$csi",
								},
								"$licenseTypes",
							},
						},
					},
				},
			})
	if err != nil {
		return utils.NewError(err, "Can't aggregate", collection)
	}

	if err := db.Collection(collection).Drop(ctx); err != nil {
		return utils.NewError(err, "Can't drop", collection)
	}

	if cursor.RemainingBatchLength() > 0 {
		var agrs []interface{}
		if err := cursor.All(ctx, &agrs); err != nil {
			return utils.NewError(err, "Can't decode cursor")
		}

		if _, err := db.Collection(collection).InsertMany(ctx, agrs); err != nil {
			return utils.NewError(err, "Can't insert all agreements")
		}
	}

	return nil
}

func group_oracle_agreementy_by_license_type_id(db *mongo.Database) error {
	return utils.NewError(errors.New("Not yet implemented"))
}
