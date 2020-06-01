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

// ExadataDevice holds informations about a device in a exadata
type ExadataDevice struct {
	Hostname       string                 `bson:"Hostname"`
	ServerType     string                 `bson:"ServerType"`
	Model          string                 `bson:"Model"`
	ExaSwVersion   string                 `bson:"ExaSwVersion"`
	CPUEnabled     string                 `bson:"CPUEnabled"`
	Memory         string                 `bson:"Memory"`
	Status         string                 `bson:"Status"`
	PowerCount     string                 `bson:"PowerCount"`
	PowerStatus    string                 `bson:"PowerStatus"`
	FanCount       string                 `bson:"FanCount"`
	FanStatus      string                 `bson:"FanStatus"`
	TempActual     string                 `bson:"TempActual"`
	TempStatus     string                 `bson:"TempStatus"`
	CellsrvService string                 `bson:"CellsrvService"`
	MsService      string                 `bson:"MsService"`
	RsService      string                 `bson:"RsService"`
	FlashcacheMode string                 `bson:"FlashcacheMode"`
	CellDisks      []ExadataCellDisk      `bson:"CellDisks"`
	_otherInfo     map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v ExadataDevice) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *ExadataDevice) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v ExadataDevice) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *ExadataDevice) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// ExadataDeviceBsonValidatorRules contains mongodb validation rules for exadata device
var ExadataDeviceBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Hostname",
		"ServerType",
		"Model",
		"ExaSwVersion",
		"CPUEnabled",
		"Memory",
		"Status",
		"PowerCount",
		"PowerStatus",
		"FanCount",
		"FanStatus",
		"TempActual",
		"TempStatus",
		"CellsrvService",
		"MsService",
		"RsService",
		"FlashcacheMode",
	}},
	{"properties", bson.D{
		{"Hostname", bson.D{
			{"bsonType", "string"},
		}},
		{"ServerType", bson.D{
			{"enum", bson.A{
				"StorageServer",
				"DBServer",
				"IBSwitch",
			}},
		}},
		{"Model", bson.D{
			{"bsonType", "string"},
		}},
		{"ExaSwVersion", bson.D{
			{"bsonType", "string"},
		}},
		{"CPUEnabled", bson.D{
			{"bsonType", "string"},
		}},
		{"Memory", bson.D{
			{"bsonType", "string"},
		}},
		{"Status", bson.D{
			{"bsonType", "string"},
		}},
		{"PowerCount", bson.D{
			{"bsonType", "string"},
		}},
		{"PowerStatus", bson.D{
			{"bsonType", "string"},
		}},
		{"FanCount", bson.D{
			{"bsonType", "string"},
		}},
		{"FanStatus", bson.D{
			{"bsonType", "string"},
		}},
		{"TempActual", bson.D{
			{"bsonType", "string"},
		}},
		{"TempStatus", bson.D{
			{"bsonType", "string"},
		}},
		{"CellsrvService", bson.D{
			{"bsonType", "string"},
		}},
		{"MsService", bson.D{
			{"bsonType", "string"},
		}},
		{"RsService", bson.D{
			{"bsonType", "string"},
		}},
		{"FlashcacheMode", bson.D{
			{"bsonType", "string"},
		}},
		{"CellDisks", bson.D{
			{"bsonType", bson.A{"array", "null"}},
			{"items", ExadataCellDiskBsonValidatorRules},
		}},
	}},
}
