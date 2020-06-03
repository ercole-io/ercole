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

// Feature holds information about Oracle database feature
type Feature struct {
	Name       string                 `bson:"Name"`
	Status     bool                   `bson:"Status"`
	_otherInfo map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Feature) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Feature) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v Feature) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Feature) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// FeatureBsonValidatorRules contains mongodb validation rules for feature
var FeatureBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Name",
		"Status",
	},
	"properties": bson.M{
		"Name": bson.M{
			"bsonType": "string",
		},
		"Status": bson.M{
			"bsonType": "bool",
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

// DiffFeature return a map that contains the difference of status between the oldFeature and newFeature
func DiffFeature(oldFeatures []Feature, newFeatures []Feature) map[string]int {
	result := make(map[string]int)

	//Add the features to the result assuming that the all new features are inactive
	for _, feature := range oldFeatures {
		if feature.Status {
			result[feature.Name] = DiffFeatureDeactivated
		} else {
			result[feature.Name] = DiffFeatureInactive
		}
	}

	//Activate/deactivate missing feature
	for _, feature := range newFeatures {
		if (result[feature.Name] == DiffFeatureInactive || result[feature.Name] == DiffFeatureMissing) && !feature.Status {
			result[feature.Name] = DiffFeatureInactive
		} else if (result[feature.Name] == DiffFeatureDeactivated) && !feature.Status {
			result[feature.Name] = DiffFeatureDeactivated
		} else if (result[feature.Name] == DiffFeatureInactive || result[feature.Name] == DiffFeatureMissing) && feature.Status {
			result[feature.Name] = DiffFeatureActivated
		} else if (result[feature.Name] == DiffFeatureDeactivated) && feature.Status {
			result[feature.Name] = DiffFeatureActive
		}
	}

	return result
}
