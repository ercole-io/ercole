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

	"github.com/ercole-io/ercole/v2/model"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	err := migrate.Register(func(db *mongo.Database) error {
		if err := insertNodes(db); err != nil {
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

func insertNodes(client *mongo.Database) error {
	collectionName := "nodes"

	if collectionNames, err := client.ListCollectionNames(context.TODO(), bson.D{{Key: "name", Value: collectionName}}); len(collectionNames) > 0 {
		if err != nil {
			return err
		}

		return nil
	}

	nodes := getNodes()

	if err := client.RunCommand(context.TODO(), bson.D{
		{Key: "create", Value: collectionName},
	}).Err(); err != nil {
		return err
	}

	if _, err := client.Collection(collectionName).InsertMany(context.TODO(), nodes); err != nil {
		return err
	}

	return nil
}

func getNodes() []interface{} {
	return []interface{}{
		model.Node{
			Name:   "Dashboard",
			Roles:  []string{"admin"},
			Parent: "",
		},

		model.Node{
			Name:   "Hosts",
			Roles:  []string{"admin"},
			Parent: "",
		},
		model.Node{
			Name:   "Databases",
			Roles:  []string{"admin"},
			Parent: "",
		},
		model.Node{
			Name:   "All Technologies",
			Roles:  []string{"admin"},
			Parent: "Databases",
		},
		model.Node{
			Name:   "Oracle",
			Roles:  []string{"admin"},
			Parent: "Databases",
		},
		model.Node{
			Name:   "DB List",
			Roles:  []string{"admin"},
			Parent: "Databases",
		},
		model.Node{
			Name:   "ADDM",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Segment Advisor",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Patch Advisor",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "DBA Role",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Patch",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Options",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Tablespaces",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Backups",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Services",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "DB Growth",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Schemas",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Pluggable DBs",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "CPU Time",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "MySQL",
			Roles:  []string{"admin"},
			Parent: "Databases",
		},
		model.Node{
			Name:   "DB List",
			Roles:  []string{"admin"},
			Parent: "MySQL",
		},
		model.Node{
			Name:   "DB List",
			Roles:  []string{"admin"},
			Parent: "Microsoft",
		},
		model.Node{
			Name:   "DB List",
			Roles:  []string{"admin"},
			Parent: "PostgreSQL",
		},
		model.Node{
			Name:   "Microsoft",
			Roles:  []string{"admin"},
			Parent: "Databases",
		},
		model.Node{
			Name:   "PostgreSQL",
			Roles:  []string{"admin"},
			Parent: "Databases",
		},
		model.Node{
			Name:   "Hypervisors",
			Roles:  []string{"admin"},
			Parent: "",
		},
		model.Node{
			Name:   "Engineered Systems",
			Roles:  []string{"admin"},
			Parent: "",
		},
		model.Node{
			Name:   "Alerts",
			Roles:  []string{"admin"},
			Parent: "",
		},
		model.Node{
			Name:   "Licenses",
			Roles:  []string{"admin"},
			Parent: "",
		},
		model.Node{
			Name:   "Licenses Contract",
			Roles:  []string{"admin"},
			Parent: "Licenses",
		},
		model.Node{
			Name:   "Licenses Compliance",
			Roles:  []string{"admin"},
			Parent: "Licenses",
		},
		model.Node{
			Name:   "Licenses Used",
			Roles:  []string{"admin"},
			Parent: "Licenses",
		},
		model.Node{
			Name:   "Cloud Advisors",
			Roles:  []string{"admin"},
			Parent: "",
		},
		model.Node{
			Name:   "Oracle",
			Roles:  []string{"admin"},
			Parent: "Cloud Advisors",
		},
		model.Node{
			Name:   "AWS",
			Roles:  []string{"admin"},
			Parent: "Cloud Advisors",
		},
		model.Node{
			Name:   "Profile Configuration",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Recommendations",
			Roles:  []string{"admin"},
			Parent: "Oracle",
		},
		model.Node{
			Name:   "Profile Configuration",
			Roles:  []string{"admin"},
			Parent: "AWS",
		},
		model.Node{
			Name:   "Recommendations",
			Roles:  []string{"admin"},
			Parent: "AWS",
		},
		model.Node{
			Name:   "Repository",
			Roles:  []string{"admin"},
			Parent: "Licenses",
		}}
}
