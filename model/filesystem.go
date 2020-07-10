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

	godynstruct "github.com/amreo/go-dyn-struct"
	"go.mongodb.org/mongo-driver/bson"
)

// Filesystem holds information about mounted filesystem and used space
type Filesystem struct {
	Filesystem     string                 `json:"filesystem" bson:"filesystem"`
	Type           string                 `json:"type" bson:"type"`
	Size           int64                  `json:"size" bson:"size"`
	UsedSpace      int64                  `json:"usedSpace" bson:"usedSpace"`
	AvailableSpace int64                  `json:"availableSpace" bson:"availableSpace"`
	MountedOn      string                 `json:"mountedOn" bson:"mountedOn"`
	OtherInfo      map[string]interface{} `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Filesystem) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Filesystem) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v Filesystem) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Filesystem) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// FilesystemBsonValidatorRules contains mongodb validation rules for Filesystem
var FilesystemBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"filesystem",
		"type",
		"size",
		"usedSpace",
		"availableSpace",
		"mountedOn",
	},
	"properties": bson.M{
		"filesystem": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"type": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"size": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"usedSpace": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"availableSpace": bson.M{
			"bsonType": "number",
		},
		"mountedOn": bson.M{
			"bsonType": "string",
		},
	},
}
