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
	"github.com/ercole-io/ercole/chart-service/chartmodel"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetOracleDatabaseChartByVersion return the chart data about oracle database version
func (md *MongoDatabase) GetOracleDatabaseChartByVersion(location string, environment string, olderThan time.Time) ([]chartmodel.ChartBubble, utils.AdvancedErrorInterface) {
	var out []chartmodel.ChartBubble = make([]chartmodel.ChartBubble, 0)
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APProject(bson.M{
				"_id":     0,
				"version": "$Features.Oracle.Database.Databases.Version",
			}),
			mu.APGroupAndCountStages("name", "size", "$version"),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item chartmodel.ChartBubble
		if err := cur.Decode(&item); err != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, item)
	}
	return out, nil
}

// GetOracleDatabaseChartByWork return the chart data about the work of all database
func (md *MongoDatabase) GetOracleDatabaseChartByWork(location string, environment string, olderThan time.Time) ([]chartmodel.ChartBubble, utils.AdvancedErrorInterface) {
	var out []chartmodel.ChartBubble = make([]chartmodel.ChartBubble, 0)
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APUnwind("$Features.Oracle.Database.Databases"),
			mu.APMatch(bson.M{
				"Features.Oracle.Database.Databases.Work": mu.QONotEqual(nil),
			}),
			mu.APProject(bson.M{
				"_id":  0,
				"name": mu.APOConcat("$Hostname", "/", "$Features.Oracle.Database.Databases.Name"),
				"size": "$Features.Oracle.Database.Databases.Work",
			}),
		),
	)
	if err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	//Decode the documents
	for cur.Next(context.TODO()) {
		var item chartmodel.ChartBubble
		if err := cur.Decode(&item); err != nil {
			return nil, utils.NewAdvancedErrorPtr(err, "Decode ERROR")
		}
		out = append(out, item)
	}
	return out, nil
}
