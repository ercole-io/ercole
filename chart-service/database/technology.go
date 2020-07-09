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

// Package database contains methods used to perform CRUD operations to the MongoDB database
package database

import (
	"context"
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTechnologyCount return the number of occurence per technology
func (md *MongoDatabase) GetTechnologyCount(location string, environment string, olderThan time.Time) (map[string]float64, utils.AdvancedErrorInterface) {
	var out map[string]float64
	//Create the operating system technology detector
	var technologyDetector bson.M = bson.M{}
	var technologyCounter bson.M = bson.M{}
	var unknownOSMatcher bson.A = bson.A{}

	// operating system
	for _, v := range md.OperatingSystemAggregationRules {
		technologyDetector[v.Product] = mu.APOCond(bson.M{
			"$regexMatch": bson.M{
				"input": mu.APOConcat("$Info.OS", " ", "$Info.OSVersion"),
				"regex": v.Regex,
			},
		}, 1, 0)

		unknownOSMatcher = append(unknownOSMatcher, bson.M{
			"$not": bson.M{
				"$regexMatch": bson.M{
					"input": mu.APOConcat("$Info.OS", " ", "$Info.OSVersion"),
					"regex": v.Regex,
				},
			},
		})
	}
	technologyDetector[model.TechnologyUnknownOperatingSystem] = mu.APOCond(mu.APOAnd(unknownOSMatcher...), 1, 0)
	// database
	technologyDetector[model.TechnologyOracleDatabase] = mu.APOSize(mu.APOIfNull("$Features.Oracle.Database.Databases", bson.A{}))

	// build the technology counter
	technologyCounter["_id"] = 0
	for k := range technologyDetector {
		technologyCounter[k] = mu.APOSum("$" + k)
	}

	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APProject(technologyDetector),
			mu.APGroup(technologyCounter),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Next the cursor. If there is no document return a empty document
	hasNext := cur.Next(context.TODO())
	if !hasNext {
		return map[string]float64{}, nil
	}

	//Decode the document
	if err := cur.Decode(&out); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return out, nil
}
