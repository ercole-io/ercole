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

type MicrosoftSQLServerProductFeature struct {
	Product    string                 `json:"product" bson:"product"`
	Instance   string                 `json:"instance" bson:"instance"`
	InstanceID string                 `json:"instanceID" bson:"instanceID"`
	Feature    string                 `json:"feature" bson:"feature"`
	Language   string                 `json:"language" bson:"language"`
	Edition    string                 `json:"edition" bson:"edition"`
	Version    string                 `json:"version" bson:"version"`
	Clustered  bool                   `json:"clustered" bson:"clustered"`
	Configured bool                   `json:"configured" bson:"configured"`
	OtherInfo  map[string]interface{} `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v MicrosoftSQLServerProductFeature) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerProductFeature) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v MicrosoftSQLServerProductFeature) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerProductFeature) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MicrosoftSQLServerProductFeatureBsonValidatorRules contains mongodb validation rules for MicrosoftSQLServerProductFeature
var MicrosoftSQLServerProductFeatureBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"product",
		"instance",
		"instanceID",
		"feature",
		"language",
		"edition",
		"version",
		"clustered",
		"configured",
	},
	"properties": bson.M{
		"product": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"instance": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"instanceID": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"feature": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"language": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"edition": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"version": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"clustered": bson.M{
			"bsonType": "boolean",
		},
		"configured": bson.M{
			"bsonType": "boolean",
		},
	},
}
