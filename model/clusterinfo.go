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

package model

import (
	"reflect"

	godynstruct "github.com/amreo/go-dyn-struct"
	"go.mongodb.org/mongo-driver/bson"
)

//ClusterInfo hold informations about a cluster
type ClusterInfo struct {
	Name      string                 `bson:"Name"`
	Type      string                 `bson:"Type"`
	CPU       int                    `bson:"CPU"`
	Sockets   int                    `bson:"Sockets"`
	VMs       []VMInfo               `bson:"VMs"`
	OtherInfo map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v ClusterInfo) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *ClusterInfo) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v ClusterInfo) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *ClusterInfo) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// ClusterInfoBsonValidatorRules contains mongodb validation rules for clusterInfo
var ClusterInfoBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Name",
		"Type",
		"CPU",
		"Sockets",
		"VMs",
	},
	"properties": bson.M{
		"Name": bson.M{
			"bsonType": "string",
		},
		"Hour": bson.M{
			"bsonType": "string",
		},
		"Type": bson.M{
			"bsonType": "string",
		},
		"CPU": bson.M{
			"bsonType": "number",
		},
		"Sockets": bson.M{
			"bsonType": "number",
		},
		"VMs": bson.M{
			"bsonType": "array",
			"items":    VMInfoBsonValidatorRules,
		},
	},
}
