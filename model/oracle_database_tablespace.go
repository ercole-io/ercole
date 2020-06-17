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

// OracleDatabaseTablespace holds the informations about a tablespace.
type OracleDatabaseTablespace struct {
	Name      string                 `bson:"Name"`
	MaxSize   float32                `bson:"MaxSize"`
	Total     float32                `bson:"Total"`
	Used      float32                `bson:"Used"`
	UsedPerc  float32                `bson:"UsedPerc"`
	Status    string                 `bson:"Status"`
	OtherInfo map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleDatabaseTablespace) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleDatabaseTablespace) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleDatabaseTablespace) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleDatabaseTablespace) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// OracleDatabaseTablespaceBsonValidatorRules contains mongodb validation rules for OracleDatabaseTablespace
var OracleDatabaseTablespaceBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Name",
		"MaxSize",
		"Total",
		"Used",
		"UsedPerc",
		"Status",
	},
	"properties": bson.M{
		"Name": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"MaxSize": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"Total": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"Used": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"UsedPerc": bson.M{
			"bsonType": "number",
			"minimum":  0,
			"maximum":  100,
		},
		"Status": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"ONLINE",
				"OFFLINE",
			},
		},
	},
}
