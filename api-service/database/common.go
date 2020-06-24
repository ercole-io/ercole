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
	"time"

	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/utils"
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
				"CreatedAt": mu.QOLessThanOrEqual(olderThan),
			}),
			mu.APLookupPipeline("hosts", bson.M{"hn": "$Hostname", "ca": "$CreatedAt"}, "Check", mu.MAPipeline(
				mu.APProject(bson.M{
					"Hostname":  1,
					"CreatedAt": 1,
				}),
				mu.APMatch(mu.QOExpr(mu.APOAnd(mu.APOEqual("$Hostname", "$$hn"), mu.APOGreater("$CreatedAt", "$$ca"), mu.APOGreaterOrEqual(olderThan, "$CreatedAt")))),
				mu.APLimit(1),
			)),
			mu.APMatch(bson.M{
				"Check": mu.QOSize(0),
			}),
			mu.APUnset("Check"),
		}),
	)
}

func AddHardwareAbstraction(field string) bson.A {
	return mu.MAPipeline(mu.APAddFields(bson.M{
		field: mu.APOOr("$ClusterMembershipStatus.OracleClusterware", "$ClusterMembershipStatus.VeritasClusterServer", "$ClusterMembershipStatus.SunCluster", "$ClusterMembershipStatus.HACMP"),
	}))
}

func AddAssociatedClusterNameAndVirtualizationNode(olderThan time.Time) bson.A {
	return mu.MAPipeline(
		mu.APLookupPipeline("hosts", bson.M{"hn": "$Hostname"}, "VM", mu.MAPipeline(
			FilterByOldnessSteps(olderThan),
			mu.APUnwind("$Clusters"),
			mu.APReplaceWith("$Clusters"),
			mu.APUnwind("$VMs"),
			mu.APSet(bson.M{
				"VMs.ClusterName": "$Name",
			}),
			mu.APReplaceWith("$VMs"),
			mu.APMatch(mu.QOExpr(mu.APOEqual("$Hostname", "$$hn"))),
			mu.APLimit(1),
		)),
		mu.APSet(bson.M{
			"VM": mu.APOArrayElemAt("$VM", 0),
		}),
		mu.APAddFields(bson.M{
			"Cluster":            mu.APOIfNull("$VM.ClusterName", nil),
			"VirtualizationNode": mu.APOIfNull("$VM.VirtualizationNode", nil),
		}),
		mu.APUnset("VM"),
	)
}
