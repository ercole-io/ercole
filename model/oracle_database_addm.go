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

// OracleDatabaseAddm contains info about a addm
type OracleDatabaseAddm struct {
	Finding        string                 `json:"finding"`
	Recommendation string                 `json:"recommendation"`
	Action         string                 `json:"action"`
	Benefit        float64                `json:"benefit"`
	OtherInfo      map[string]interface{} `json:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleDatabaseAddm) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleDatabaseAddm) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleDatabaseAddm) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleDatabaseAddm) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// OracleDatabaseAddmBsonValidatorRules contains mongodb validation rules for OracleDatabaseAddm
var OracleDatabaseAddmBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"finding",
		"recommendation",
		"action",
		"benefit",
	},
	"properties": bson.M{
		"finding": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 256,
		},
		"recommendation": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"action": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 256,
		},
		"benefit": bson.M{
			"bsonType": "string",
			"minimum":  0,
			"maximum":  100,
		},
	},
}
