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

// OracleDatabasePatch holds information about a Oracle database patch
type OracleDatabasePatch struct {
	Version     string                 `json:"version"`
	PatchID     int                    `json:"patchID bson:patchID"`
	Action      string                 `json:"action"`
	Description string                 `json:"description"`
	Date        string                 `json:"date"`
	OtherInfo   map[string]interface{} `json:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleDatabasePatch) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleDatabasePatch) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleDatabasePatch) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleDatabasePatch) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// OracleDatabasePatchBsonValidatorRules contains mongodb validation rules for OracleDatabasePatch
var OracleDatabasePatchBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"version",
		"patchID",
		"action",
		"description",
		"date",
	},
	"properties": bson.M{
		"version": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"patchID": bson.M{
			"bsonType": "number",
		},
		"action": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 128,
		},
		"description": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 256,
		},
		"date": bson.M{
			"bsonType": "string",
			"pattern":  "[0-9]{4}-[0-9]{2}-[0-9]{2}",
		},
	},
}
