// Copyright (c) 2025 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package database

import (
	"context"

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) ExistsDR(hostname string) bool {
	filter := bson.M{"archived": false, "isDR": true, "hostname": hostname}

	res := md.Client.Database(md.Config.Mongodb.DBName).
		Collection("hosts").FindOne(context.Background(), filter)

	return res.Err() == nil
}

func (md *MongoDatabase) GetClusterVeritasLicenseByHostnames(hostnames []string) ([]model.OracleDatabaseLicense, error) {
	pipeline :=
		bson.A{
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "archived", Value: false},
						{Key: "clusterMembershipStatus.veritasClusterHostnames",
							Value: bson.D{
								{Key: "$in",
									Value: hostnames,
								},
							},
						},
					},
				},
			},
			bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases"}}}},
			bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$features.oracle.database.databases.licenses"}}}},
			bson.D{
				{Key: "$group",
					Value: bson.D{
						{Key: "_id",
							Value: bson.D{
								{Key: "licenseTypeID", Value: "$features.oracle.database.databases.licenses.licenseTypeID"},
								{Key: "name", Value: "$features.oracle.database.databases.licenses.name"},
								{Key: "count", Value: "$features.oracle.database.databases.licenses.count"},
								{Key: "ignored", Value: "$features.oracle.database.databases.licenses.ignored"},
								{Key: "ignoredComment", Value: "$features.oracle.database.databases.licenses.ignoredComment"},
							},
						},
					},
				},
			},
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "_id.licenseTypeID",
							Value: bson.D{
								{Key: "$exists", Value: true},
								{Key: "$ne", Value: ""},
							},
						},
					},
				},
			},
			bson.D{
				{Key: "$replaceWith",
					Value: bson.D{
						{Key: "licenseTypeID", Value: "$_id.licenseTypeID"},
						{Key: "name", Value: "$_id.name"},
						{Key: "count", Value: "$_id.count"},
						{Key: "ignored", Value: "$_id.ignored"},
						{Key: "ignoredComment", Value: "$_id.ignoredComment"},
					},
				},
			},
		}

	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	licenses := make([]model.OracleDatabaseLicense, 0)

	for cur.Next(context.Background()) {
		license := model.OracleDatabaseLicense{}
		if err := cur.Decode(&license); err != nil {
			return nil, err
		}

		licenses = append(licenses, license)
	}

	return licenses, nil
}
