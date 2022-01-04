// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"errors"
	"reflect"

	godynstruct "github.com/amreo/go-dyn-struct"
	"github.com/hashicorp/go-multierror"
)

// SchemaVersion contains the version of the schema
const SchemaVersion int = 1

// HostData holds all informations about a host & services
type HostData struct {
	Hostname                string                  `json:"hostname"`
	Location                string                  `json:"location"`
	Environment             string                  `json:"environment"`
	Tags                    []string                `json:"tags"`
	AgentVersion            string                  `json:"agentVersion"`
	SchemaVersion           int                     `json:"schemaVersion"`
	Info                    Host                    `json:"info"`
	ClusterMembershipStatus ClusterMembershipStatus `json:"clusterMembershipStatus"`
	Features                Features                `json:"features"`
	Filesystems             []Filesystem            `json:"filesystems"`
	Clusters                []ClusterInfo           `json:"clusters"`
	Cloud                   Cloud                   `json:"cloud"`
	Errors                  []AgentError            `json:"errors"`
	OtherInfo               map[string]interface{}  `json:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v HostData) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *HostData) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v HostData) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *HostData) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

func (v *HostData) AddErrors(errs ...error) {
	for _, e := range errs {
		if e == nil {
			continue
		}

		var merr *multierror.Error
		if errors.As(e, &merr) {
			for _, merre := range merr.Errors {
				v.Errors = append(v.Errors, NewAgentError(merre))
			}

			continue
		}

		v.Errors = append(v.Errors, NewAgentError(e))
	}
}
