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

package model

import (
	"reflect"
	"time"

	godynstruct "github.com/amreo/go-dyn-struct"
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
	ServerVersion       string             `json:"serverVersion" bson:"serverVersion"`
	ServerSchemaVersion int                `json:"serverSchemaVersion" bson:"serverSchemaVersion"`

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
	OtherInfo               map[string]interface{}  `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v HostDataBE) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *HostDataBE) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v HostDataBE) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *HostDataBE) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
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

//TODO Deduplicate this method with HostData method
func (v *HostDataBE) CoreFactor() float64 {
	if v.Cloud.Membership == CloudMembershipAws {
		return 1
	}

	return 0.5
}
