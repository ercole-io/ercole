// Copyright (c) 2020 Sorint.lab S.p.A.
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

package database

import (
	"context"

	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (md *MongoDatabase) SearchMySQLInstances(filter dto.GlobalFilter) ([]dto.MySQLInstance, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			mu.APUnwind("$features.mysql.instances"),
			mu.APProject(bson.M{
				"hostname":    1,
				"location":    1,
				"environment": 1,
				"instance":    "$features.mysql.instances",
			}),
			mu.APReplaceWith(mu.APOMergeObjects("$$ROOT", "$instance")),
			mu.APUnset("instance"),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var out []dto.MySQLInstance
	err = cur.All(context.TODO(), &out)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}

func (md *MongoDatabase) GetMySQLUsedLicenses(filter dto.GlobalFilter) ([]dto.MySQLUsedLicense, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			mu.APUnwind("$features.mysql.instances"),
			mu.APMatch(bson.M{
				// Only ENTERPRISE MySQL db are considered as licenses
				"features.mysql.instances.edition": model.MySQLEditionEnterprise,
			}),
			mu.APProject(bson.M{
				"hostname":        1,
				"instanceName":    "$features.mysql.instances.name",
				"instanceEdition": "$features.mysql.instances.edition",
			}),
			mu.APReplaceWith(mu.APOMergeObjects("$$ROOT", "$instance")),
			mu.APUnset("instance"),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var out []dto.MySQLUsedLicense
	err = cur.All(context.TODO(), &out)
	if err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}
