// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindPatchingFunction find the the patching function associated to the hostname in the database
func (md *MongoDatabase) FindPatchingFunction(hostname string) (model.PatchingFunction, utils.AdvancedErrorInterface) {
	var out model.PatchingFunction

	//Find the hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("patching_functions").Find(context.TODO(), bson.M{
		"Hostname": hostname,
	})
	if err != nil {
		return model.PatchingFunction{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return model.PatchingFunction{}, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return model.PatchingFunction{}, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}

// SavePatchingFunction saves the patching function
func (md *MongoDatabase) SavePatchingFunction(pf model.PatchingFunction) utils.AdvancedErrorInterface {
	//Find the informations
	true := true
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("patching_functions").ReplaceOne(context.TODO(), bson.M{
		"_id": pf.ID,
	}, pf, &options.ReplaceOptions{
		Upsert: &true,
	})
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return nil
}
