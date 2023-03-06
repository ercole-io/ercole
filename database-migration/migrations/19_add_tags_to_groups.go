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
	tagAdmin := []string{"admin"}
	tagDba := []string{"dba"}
	tagProcurement := []string{"procurement"}

	err := migrate.Register(func(db *mongo.Database) error {
		if err := addTagsToGroup(tagAdmin, "admin", db); err != nil {
			return err
		}

		if err := addTagsToGroup(tagDba, "dba", db); err != nil {
			return err
		}

		if err := addTagsToGroup(tagProcurement, "procurement", db); err != nil {
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

func addTagsToGroup(tags []string, groupName string, client *mongo.Database) error {
	collectionName := "groups"
	ctx := context.TODO()

	if cols, err := client.ListCollectionNames(ctx, bson.D{}); err != nil {
		return err
	} else if !utils.Contains(cols, collectionName) {
		return fmt.Errorf("%s not found", collectionName)
	}

	if _, err := client.Collection(collectionName).UpdateOne(ctx, bson.M{"name": groupName}, bson.M{"$set": bson.M{"tags": tags}}); err != nil {
		return err
	}

	return nil
}
