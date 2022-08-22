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

package model

import (
	"time"

	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServerSchemaVersion contains the version of the schema
const ServerSchemaVersion int = 1

// HostDataBE holds all informations about a host & services
type HostDataBE struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	Archived            bool               `json:"archived" bson:"archived"`
	CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	DismissedAt         time.Time          `json:"dismissedAt" bson:"dismissedAt,omitempty"`
	ServerVersion       string             `json:"serverVersion" bson:"serverVersion"`
	ServerSchemaVersion int                `json:"serverSchemaVersion" bson:"serverSchemaVersion"`
	Period              uint               `json:"period" bson:"period"`

	Hostname                string                  `json:"hostname" bson:"hostname"`
	Location                string                  `json:"location" bson:"location"`
	Environment             string                  `json:"environment" bson:"environment"`
	AgentVersion            string                  `json:"agentVersion" bson:"agentVersion"`
	Tags                    []string                `json:"tags" bson:"tags"`
	Info                    Host                    `json:"info" bson:"info"`
	ClusterMembershipStatus ClusterMembershipStatus `json:"clusterMembershipStatus" bson:"clusterMembershipStatus"`
	Features                Features                `json:"features" bson:"features"`
	Filesystems             []Filesystem            `json:"filesystems" bson:"filesystems"`
	Clusters                []ClusterInfo           `json:"clusters" bson:"clusters"`
	Cloud                   Cloud                   `json:"cloud" bson:"cloud"`
	Errors                  []AgentError            `json:"errors" bson:"errors"`
}

func (v *HostDataBE) GetClusterCores(hostdatasPerHostname map[string]*HostDataBE) (int, error) {
	cms := v.ClusterMembershipStatus
	if !cms.VeritasClusterServer ||
		(cms.VeritasClusterServer && len(cms.VeritasClusterHostnames) <= 2) {
		return 0, utils.ErrHostNotInCluster
	}

	var sumClusterCores int

	for _, h := range cms.VeritasClusterHostnames {
		anotherHostdata, found := hostdatasPerHostname[h]
		if !found {
			sumClusterCores += v.Info.CPUCores // Use current hostdata as fallback
			continue
		}

		sumClusterCores += anotherHostdata.Info.CPUCores
	}

	return sumClusterCores, nil
}

func (v *HostDataBE) CoreFactor() float64 {
	if v.Cloud.Membership == CloudMembershipAws {
		return 1
	}

	return 0.5
}
