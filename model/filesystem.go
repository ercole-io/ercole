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

// Filesystem holds information about mounted filesystem and used space
type Filesystem struct {
	Filesystem string `bson:"Filesystem"`
	FsType     string `bson:"FsType"`
	Size       string `bson:"Size"`
	Used       string `bson:"Used"`
	Available  string `bson:"Available"`
	UsedPerc   string `bson:"UsedPerc"`
	MountedOn  string `bson:"MountedOn"`
	_otherInfo map[string]interface{}
}

// MarshalJSON return the JSON rappresentation of this
func (v Filesystem) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Filesystem) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v Filesystem) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Filesystem) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// FilesystemBsonValidatorRules contains mongodb validation rules for filesystem
var FilesystemBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Filesystem",
		"FsType",
		"Size",
		"Used",
		"Available",
		"UsedPerc",
		"MountedOn",
	}},
	{"properties", bson.D{
		{"Filesystem", bson.D{
			{"bsonType", "string"},
		}},
		{"FsType", bson.D{
			{"bsonType", "string"},
		}},
		{"Size", bson.D{
			{"bsonType", "string"},
		}},
		{"Used", bson.D{
			{"bsonType", "string"},
		}},
		{"Available", bson.D{
			{"bsonType", "string"},
		}},
		{"UsedPerc", bson.D{
			{"bsonType", "string"},
		}},
		{"MountedOn", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
