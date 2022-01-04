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

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		if err := create_index_hosts(db); err != nil {
			return err
		}
		return nil

	}, func(db *mongo.Database) error {
		return nil
	})
}

func create_index_hosts(db *mongo.Database) error {
	if _, err := db.Collection("hosts").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"createdAt", -1},
		},
	}); err != nil {
		return err
	}

	return nil
}
