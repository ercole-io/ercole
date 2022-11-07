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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"context"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (md *MongoDatabase) FindAllOracleDatabasePdbs(filter dto.GlobalFilter) ([]dto.OracleDatabasePluggableDatabase, error) {
	ctx := context.TODO()
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		ctx,
		mu.MAPipeline(
			FilterByOldnessSteps(filter.OlderThan),
			FilterByLocationAndEnvironmentSteps(filter.Location, filter.Environment),
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases"}},
			bson.M{"$unwind": bson.M{"path": "$features.oracle.database.databases.pdbs"}},
			bson.M{"$project": bson.M{
				"hostname": 1,
				"pdb":      "$features.oracle.database.databases.pdbs",
			}},
		),
	)

	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	out := make([]dto.OracleDatabasePluggableDatabase, 0)
	if err = cur.All(ctx, &out); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return out, nil
}
