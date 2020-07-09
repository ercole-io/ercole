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
	"go.mongodb.org/mongo-driver/bson"
)

// OracleDatabaseFeatureUsageStat holds information about a oracle database feature usage stat.
type OracleDatabaseFeatureUsageStat struct {
	Product          string                 `json:"product"`
	Feature          string                 `json:"feature"`
	DetectedUsages   int64                  `json:"detectedUsages bson:detectedUsages"`
	CurrentlyUsed    bool                   `json:"currentlyUsed bson:currentlyUsed"`
	FirstUsageDate   time.Time              `json:"firstUsageDate bson:firstUsageDate"`
	LastUsageDate    time.Time              `json:"lastUsageDate bson:lastUsageDate"`
	ExtraFeatureInfo string                 `json:"extraFeatureInfo bson:extraFeatureInfo"`
	OtherInfo        map[string]interface{} `json:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleDatabaseFeatureUsageStat) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleDatabaseFeatureUsageStat) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleDatabaseFeatureUsageStat) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleDatabaseFeatureUsageStat) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// OracleDatabaseFeatureUsageStatBsonValidatorRules contains mongodb validation rules for OracleDatabaseFeatureUsageStat
var OracleDatabaseFeatureUsageStatBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"product",
		"feature",
		"detectedUsages",
		"currentlyUsed",
		"firstUsageDate",
		"lastUsageDate",
		"extraFeatureInfo",
	},
	"properties": bson.M{
		"product": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"feature": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"detectedUsages": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
		"currentlyUsed": bson.M{
			"bsonType": "bool",
		},
		"firstUsageDate": bson.M{
			"bsonType": "date",
		},
		"lastUsageDate": bson.M{
			"bsonType": "date",
		},
		"extraFeatureInfo": bson.M{
			"bsonType":  "string",
			"maxLength": 64,
		},
	},
}
