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
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// FilterByLocationAndEnvironmentSteps return the steps required to filter the data by the location and environment field.
func FilterByLocationAndEnvironmentSteps(location string, environment string) interface{} {
	return bson.A{
		utils.MongoAggregationOptionalStep(location != "", bson.M{"$match": bson.M{
			"location": location,
		}}),
		utils.MongoAggregationOptionalStep(environment != "", bson.M{"$match": bson.M{
			"environment": environment,
		}}),
	}
}
