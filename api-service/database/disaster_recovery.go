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
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (md *MongoDatabase) CreateDR(hostname string, clusterVeritasHostnames []string) (string, error) {
	filter := bson.M{"archived": false, "isDR": false, "hostname": hostname}

	res := md.Client.Database(md.Config.Mongodb.DBName).
		Collection(hostCollection).FindOne(context.Background(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return "", utils.ErrHostNotFound
		}

		return "", res.Err()
	}

	var host model.HostDataBE
	if err := res.Decode(&host); err != nil {
		return "", res.Err()
	}

	host.ID = primitive.NewObjectID()
	host.Hostname = fmt.Sprintf("%s_DR", hostname)
	host.IsDR = true

	if host.ClusterMembershipStatus.VeritasClusterServer {
		host.ClusterMembershipStatus.VeritasClusterHostnames = clusterVeritasHostnames
	}

	_, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection(hostCollection).InsertOne(context.Background(), host)
	if err != nil {
		return "", err
	}

	return host.Hostname, nil
}
