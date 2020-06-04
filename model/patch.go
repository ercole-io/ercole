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

// Patch holds information about a Oracle patch
type Patch struct {
	Database    string                 `bson:"Database"`
	Version     string                 `bson:"Version"`
	PatchID     string                 `bson:"PatchID"`
	Action      string                 `bson:"Action"`
	Description string                 `bson:"Description"`
	Date        string                 `bson:"Date"`
	OtherInfo   map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Patch) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Patch) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v Patch) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Patch) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// PatchBsonValidatorRules contains mongodb validation rules for patch
var PatchBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Database",
		"Version",
		"PatchID",
		"Action",
		"Description",
		"Date",
	},
	"properties": bson.M{
		"Database": bson.M{
			"bsonType": "string",
		},
		"Version": bson.M{
			"bsonType": "string",
		},
		"PatchID": bson.M{
			"bsonType": "string",
		},
		"Action": bson.M{
			"bsonType": "string",
		},
		"Description": bson.M{
			"bsonType": "string",
		},
		"Date": bson.M{
			"bsonType": "string",
		},
	},
}
