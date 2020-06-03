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

// VMInfo holds info about the vm
type VMInfo struct {
	Name         string                 `bson:"Name"`
	ClusterName  string                 `bson:"ClusterName"`
	Hostname     string                 `bson:"Hostname"` //Hostname or IP address
	CappedCPU    bool                   `bson:"CappedCPU"`
	PhysicalHost string                 `bson:"PhysicalHost"`
	_otherInfo   map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v VMInfo) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *VMInfo) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v VMInfo) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *VMInfo) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// VMInfoBsonValidatorRules contains mongodb validation rules for VMInfo
var VMInfoBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Name",
		"ClusterName",
		"Hostname",
		"CappedCPU",
		"PhysicalHost",
	},
	"properties": bson.M{
		"Name": bson.M{
			"bsonType": "string",
		},
		"ClusterName": bson.M{
			"bsonType": "string",
		},
		"Hostname": bson.M{
			"bsonType": "string",
		},
		"CappedCPU": bson.M{
			"bsonType": "bool",
		},
		"PhysicalHost": bson.M{
			"bsonType": "string",
		},
	},
}
