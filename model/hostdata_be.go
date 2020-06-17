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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServerSchemaVersion contains the version of the schema
const ServerSchemaVersion int = 1

// HostDataBE holds all informations about a host & services
type HostDataBE struct {
	ID        primitive.ObjectID `bson:"_id"`
	Archived  bool               `bson:"Archived"`
	CreatedAt time.Time          `bson:"CreatedAt"`

	Hostname                string                  `bson:"Hostname"`
	Location                string                  `bson:"Location"`
	Environment             string                  `bson:"Environment"`
	Tags                    []string                `bson:"Tags"`
	AgentVersion            string                  `bson:"AgentVersion"`
	SchemaVersion           int                     `bson:"SchemaVersion"`
	Info                    Host                    `bson:"Info"`
	ClusterMembershipStatus ClusterMembershipStatus `bson:"ClusterMembershipStatus"`
	Features                Features                `bson:"Features"`
	Filesystems             []Filesystem            `bson:"Filesystems"`
	Clusters                []ClusterInfo           `bson:"Clusters"`
	OtherInfo               map[string]interface{}  `bson:"OtherInfo"`
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
