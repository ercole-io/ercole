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

// Host contains info about the host
type Host struct {
	Hostname       string  `bson:"Hostname"`
	Environment    string  `bson:"Environment"`
	Location       string  `bson:"Location"`
	CPUModel       string  `bson:"CPUModel"`
	CPUCores       int     `bson:"CPUCores"`
	CPUThreads     int     `bson:"CPUThreads"`
	Socket         int     `bson:"Socket"`
	Type           string  `bson:"Type"`
	Virtual        bool    `bson:"Virtual"`
	Kernel         string  `bson:"Kernel"`
	OS             string  `bson:"OS"`
	MemoryTotal    float32 `bson:"MemoryTotal"`
	SwapTotal      float32 `bson:"SwapTotal"`
	OracleCluster  bool    `bson:"OracleCluster"`
	VeritasCluster bool    `bson:"VeritasCluster"`
	SunCluster     bool    `bson:"SunCluster"`
	AixCluster     bool    `bson:"AixCluster"`
	_otherInfo     map[string]interface{}
}

// MarshalJSON return the JSON rappresentation of this
func (v Host) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Host) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v Host) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Host) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// HostBsonValidatorRules contains mongodb validation rules for host
var HostBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Hostname",
		"Environment",
		"Location",
		"CPUModel",
		"CPUCores",
		"CPUThreads",
		"Socket",
		"Type",
		"Virtual",
		"Kernel",
		"OS",
		"MemoryTotal",
		"SwapTotal",
		"OracleCluster",
		"VeritasCluster",
		"SunCluster",
		"AixCluster",
	}},
	{"properties", bson.D{
		{"Hostname", bson.D{
			{"bsonType", "string"},
			{"pattern", "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"},
		}},
		{"Environment", bson.D{
			{"bsonType", "string"},
		}},
		{"Location", bson.D{
			{"bsonType", "string"},
		}},
		{"CPUModel", bson.D{
			{"bsonType", "string"},
		}},
		{"CPUCores", bson.D{
			{"bsonType", "number"},
		}},
		{"CPUThreads", bson.D{
			{"bsonType", "number"},
		}},
		{"Socket", bson.D{
			{"bsonType", "number"},
		}},
		{"Type", bson.D{
			{"bsonType", "string"},
		}},
		{"Virtual", bson.D{
			{"bsonType", "bool"},
		}},
		{"Kernel", bson.D{
			{"bsonType", "string"},
		}},
		{"OS", bson.D{
			{"bsonType", "string"},
		}},
		{"MemoryTotal", bson.D{
			{"bsonType", "double"},
		}},
		{"SwapTotal", bson.D{
			{"bsonType", "double"},
		}},
		{"OracleCluster", bson.D{
			{"bsonType", "bool"},
		}},
		{"VeritasCluster", bson.D{
			{"bsonType", "bool"},
		}},
		{"SunCluster", bson.D{
			{"bsonType", "bool"},
		}},
		{"AixCluster", bson.D{
			{"bsonType", "bool"},
		}},
	}},
}
