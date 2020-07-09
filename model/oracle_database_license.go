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

// OracleDatabaseLicense holds information about a Oracle database license
type OracleDatabaseLicense struct {
	Name      string                 `json:"name"`
	Count     float64                `json:"count"`
	OtherInfo map[string]interface{} `json:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleDatabaseLicense) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleDatabaseLicense) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleDatabaseLicense) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleDatabaseLicense) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// OracleDatabaseLicenseBsonValidatorRules contains mongodb validation rules for OracleDatabaseLicense
var OracleDatabaseLicenseBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"name",
		"count",
	},
	"properties": bson.M{
		"name": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"count": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
	},
}

// DiffFeature status of each feature
const (
	// DiffFeatureInactive is used when the feature changes from (0/-) to 0
	DiffFeatureInactive int = -2
	// DiffFeatureDeactivated is used when the feature changes from 1 to (0/-)
	DiffFeatureDeactivated int = -1
	// DiffFeatureMissing is used when a feature is missing in the diff
	DiffFeatureMissing int = 0
	// DiffFeatureActivated is used when the feature changes from (0/-) to 1
	DiffFeatureActivated int = 1
	// DiffFeatureInactive is used when the feature changes from 1 to 1
	DiffFeatureActive int = 2
)

// DiffLicenses return a map that contains the difference of status between the oldLicenses and newLicenses
func DiffLicenses(oldLicenses []OracleDatabaseLicense, newLicenses []OracleDatabaseLicense) map[string]int {
	result := make(map[string]int)

	//Add the features to the result assuming that the all new features are inactive
	for _, license := range oldLicenses {
		if license.Count > 0 {
			result[license.Name] = DiffFeatureDeactivated
		} else {
			result[license.Name] = DiffFeatureInactive
		}
	}

	//Activate/deactivate missing feature
	for _, license := range newLicenses {
		if (result[license.Name] == DiffFeatureInactive || result[license.Name] == DiffFeatureMissing) && license.Count <= 0 {
			result[license.Name] = DiffFeatureInactive
		} else if (result[license.Name] == DiffFeatureDeactivated) && license.Count <= 0 {
			result[license.Name] = DiffFeatureDeactivated
		} else if (result[license.Name] == DiffFeatureInactive || result[license.Name] == DiffFeatureMissing) && license.Count > 0 {
			result[license.Name] = DiffFeatureActivated
		} else if (result[license.Name] == DiffFeatureDeactivated) && license.Count > 0 {
			result[license.Name] = DiffFeatureActive
		}
	}

	return result
}
