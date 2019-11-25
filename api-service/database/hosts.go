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

	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchCurrentHosts search current hosts
func (md *MongoDatabase) SearchCurrentHosts(full bool, keywords []string, sortBy string, sortDesc bool, page int, pageSize int) ([]interface{}, utils.AdvancedErrorInterface) {
	var out []interface{}
	var quotedKeywords []string
	for _, k := range keywords {
		quotedKeywords = append(quotedKeywords, regexp.QuoteMeta(k))
	}
	//Find the matching hostdata
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		bson.A{
			bson.D{{"$match", bson.D{
				{"archived", false},
				{"$or", bson.A{
					bson.D{{"hostname", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"extra.databases.name", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"extra.databases.unique_name", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
					bson.D{{"extra.clusters.name", bson.D{
						{"$regex", primitive.Regex{Pattern: strings.Join(quotedKeywords, "|"), Options: "i"}},
					}}},
				}},
			}}},
			optionalStep(!full, bson.D{{"$project", bson.D{
				{"hostname", true},
				{"environment", true},
				{"host_type", true},
				{"cluster", ""},
				{"physical_host", ""},
				{"created_at", true},
				{"databases", true},
				{"os", "$info.os"},
				{"kernel", "$info.kernel"},
				{"oracle_cluster", "$info.oracle_cluster"},
				{"sun_cluster", "$info.sun_cluster"},
				{"veritas_cluster", "$info.veritas_cluster"},
				{"virtual", "$info.virtual"},
				{"type", "$info.type"},
				{"cpu_threads", "$info.cpu_threads"},
				{"cpu_cores", "$info.cpu_cores"},
				{"socket", "$info.socket"},
				{"mem_total", "$info.memory_total"},
				{"swap_total", "$info.swap_total"},
				{"cpu_model", "$info.cpu_model"},
			}}}),
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
