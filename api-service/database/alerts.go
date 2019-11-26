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
	"regexp"
	"strings"

	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchAlerts search alerts
func (md *MongoDatabase) SearchAlerts(keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	var quotedKeywords []string
	for _, k := range keywords {
		quotedKeywords = append(quotedKeywords, regexp.QuoteMeta(k))
	}
	//Find the matching alerts
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("alerts").Aggregate(
		context.TODO(),
		bson.A{
			bson.D{{"$match", bson.D{
				{"alert_status", model.AlertStatusNew},
				{"$or", bson.A{
					bson.D{{"description", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"alert_code", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"alert_severity", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"other_info.hostname", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"other_info.dbname", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"other_info.features", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
				}},
			}}},
			bson.D{{"$unset", bson.A{
				"other_info",
			}}},
			optionalSortingStep(sortBy, sortDesc),
			optionalPagingStep(page, pageSize),
		},
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
		out = append(out, &item)
	}
	return out, nil
}
