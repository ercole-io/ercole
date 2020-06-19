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

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
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

// SearchOracleDatabaseLicenseModifiers search license modifiers
func (md *MongoDatabase) SearchOracleDatabaseLicenseModifiers(keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	var out []map[string]interface{} = make([]map[string]interface{}, 0)

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("patching_functions").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APProject(bson.M{
				"Hostname": 1,
				"LicenseModifiers": bson.M{
					"$objectToArray": "$Vars.LicenseModifiers",
				},
			}),
			mu.APUnwind("$LicenseModifiers"),
			mu.APProject(bson.M{
				"Hostname":     1,
				"DatabaseName": "$LicenseModifiers.k",
				"License": bson.M{
					"$objectToArray": "$LicenseModifiers.v",
				},
			}),
			mu.APUnwind("$License"),
			mu.APProject(bson.M{
				"Hostname":     1,
				"DatabaseName": 1,
				"LicenseName":  "$License.k",
				"NewValue":     "$License.v",
			}),
			mu.APSearchFilterStage([]interface{}{
				"$Hostname",
				"$DatabaseName",
				"$LicenseName",
			}, keywords),
			mu.APOptionalSortingStage(sortBy, sortDesc),
			mu.APOptionalPagingStage(page, pageSize),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		if cur.Decode(&item) != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, item)
	}
	return out, nil
}

// DeletePatchingFunction delete the patching function
func (md *MongoDatabase) DeletePatchingFunction(hostname string) utils.AdvancedErrorInterface {
	//Find the informations
	_, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("patching_functions").DeleteOne(context.TODO(), bson.M{
		"Hostname": hostname,
	}, nil)
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}
	return nil
}
