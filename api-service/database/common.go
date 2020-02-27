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
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/amreo/mu"
	"go.mongodb.org/mongo-driver/bson"
)

// FilterByLocationAndEnvironmentSteps return the steps required to filter the data by the location and environment field.
func FilterByLocationAndEnvironmentSteps(location string, environment string) interface{} {
	return bson.A{
		mu.APOptionalStage(location != "", mu.APMatch(bson.M{
			"Location": location,
		})),
		mu.APOptionalStage(environment != "", mu.APMatch(bson.M{
			"Environment": environment,
		})),
	}
}

func FilterByOldnessSteps(olderThan time.Time) bson.A {
	return mu.MAPipeline(
		mu.APOptionalStage(olderThan == utils.MAX_TIME, mu.APMatch(bson.M{
			"Archived": false,
		})),
		mu.APOptionalStage(olderThan != utils.MAX_TIME, bson.A{
			mu.APMatch(bson.M{
				"Created_at": mu.QOLessThanOrEqual(olderThan),
			}),
			mu.APLookupPipeline("hosts", bson.M{"hn": "$Hostname", "ca": "$CreatedAt"}, "Check", mu.MAPipeline(
				mu.APProject(bson.M{
					"Hostname":  1,
					"CreatedAt": 1,
				}),
				mu.APMatch(bson.M{
					"$expr": mu.APOAnd(mu.APOEqual("$Hostname", "$$hn"), mu.APOGreater("$CreatedAt", "$$ca"), mu.APOGreaterOrEqual(olderThan, "$CreatedAt")),
				}),
				mu.APLimit(1),
			)),
			mu.APMatch(bson.M{
				"Check": bson.M{
					"$size": 0,
				},
			}),
			mu.APUnset("Check"),
		}),
	)
}
