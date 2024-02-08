// Copyright (c) 2024 Sorint.lab S.p.A.
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
	err := migrate.Register(create_index_exadata_vm_clustername, nil)

	if err != nil {
		panic(err)
	}
}

func create_index_exadata_vm_clustername(db *mongo.Database) error {
	if _, err := db.Collection("exadata_vm_clusternames").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "instancerackid", Value: -1},
			{Key: "hostid", Value: -1},
			{Key: "vmname", Value: -1},
		},
	}); err != nil {
		return err
	}

	return nil
}
