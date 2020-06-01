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

// Tablespace holds the informations about a tablespace.
type Tablespace struct {
	Database   string `bson:"Database"`
	Name       string `bson:"Name"`
	MaxSize    string `bson:"MaxSize"`
	Total      string `bson:"Total"`
	Used       string `bson:"Used"`
	UsedPerc   string `bson:"UsedPerc"`
	Status     string `bson:"Status"`
	_otherInfo map[string]interface{}
}

// MarshalJSON return the JSON rappresentation of this
func (v Tablespace) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Tablespace) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v Tablespace) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Tablespace) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// TablespaceBsonValidatorRules contains mongodb validation rules for tablespace
var TablespaceBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Database",
		"Name",
		"MaxSize",
		"Total",
		"Used",
		"UsedPerc",
		"Status",
	}},
	{"properties", bson.D{
		{"Database", bson.D{
			{"bsonType", "string"},
		}},
		{"Name", bson.D{
			{"bsonType", "string"},
		}},
		{"MaxSize", bson.D{
			{"bsonType", "string"},
		}},
		{"Total", bson.D{
			{"bsonType", "string"},
		}},
		{"Used", bson.D{
			{"bsonType", "string"},
		}},
		{"UsedPerc", bson.D{
			{"bsonType", "string"},
		}},
		{"Status", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
