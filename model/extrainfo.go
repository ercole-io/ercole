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

// ExtraInfo holds various informations.
type ExtraInfo struct {
	Databases   []Database    `bson:"Databases"`
	Filesystems []Filesystem  `bson:"Filesystems"`
	Clusters    []ClusterInfo `bson:"Clusters"`
	Exadata     *Exadata      `bson:"Exadata"`
	_otherInfo  map[string]interface{}
}

// MarshalJSON return the JSON rappresentation of this
func (v ExtraInfo) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *ExtraInfo) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v ExtraInfo) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *ExtraInfo) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// ExtraInfoBsonValidatorRules contains mongodb validation rules for extraInfo
var ExtraInfoBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Filesystems",
	}},
	{"properties", bson.D{
		{"Databases", bson.D{
			{"anyOf", bson.A{
				bson.D{
					{"bsonType", "array"},
					{"items", DatabaseBsonValidatorRules},
				},
				bson.D{{"type", "null"}},
			}},
		}},
		{"Filesystems", bson.D{
			{"bsonType", "array"},
			{"items", FilesystemBsonValidatorRules},
		}},
		{"Clusters", bson.D{
			{"anyOf", bson.A{
				bson.D{
					{"bsonType", "array"},
					{"items", ClusterInfoBsonValidatorRules},
				},
				bson.D{{"type", "null"}},
			}},
		}},
		{"Exadata", bson.D{
			{"anyOf", bson.A{
				ExadataBsonValidatorRules,
				bson.D{{"type", "null"}},
			}},
		}},
	}},
}

// ExadataBsonValidatorRules
