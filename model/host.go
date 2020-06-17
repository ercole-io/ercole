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

// Host contains info about the host
type Host struct {
	Hostname                      string                 `bson:"Hostname"`
	CPUModel                      string                 `bson:"CPUModel"`
	CPUFrequency                  string                 `bson:"CPUFrequency"`
	CPUSockets                    int                    `bson:"CPUSockets"`
	CPUCores                      int                    `bson:"CPUCores"`
	CPUThreads                    int                    `bson:"CPUThreads"`
	ThreadsPerCore                int                    `bson:"ThreadsPerCore"`
	CoresPerSocket                int                    `bson:"CoresPerSocket"`
	HardwareAbstraction           string                 `bson:"HardwareAbstraction"`
	HardwareAbstractionTechnology string                 `bson:"HardwareAbstractionTechnology"`
	Kernel                        string                 `bson:"Kernel"`
	KernelVersion                 string                 `bson:"KernelVersion"`
	OS                            string                 `bson:"OS"`
	OSVersion                     string                 `bson:"OSVersion"`
	MemoryTotal                   float32                `bson:"MemoryTotal"`
	SwapTotal                     float32                `bson:"SwapTotal"`
	OtherInfo                     map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Host) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Host) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v Host) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Host) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// HostBsonValidatorRules contains mongodb validation rules for host
var HostBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
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
	},
	"properties": bson.M{
		"Hostname": bson.M{
			"bsonType": "string",
			"pattern":  "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$",
		},
		"Environment": bson.M{
			"bsonType": "string",
		},
		"Location": bson.M{
			"bsonType": "string",
		},
		"CPUModel": bson.M{
			"bsonType": "string",
		},
		"CPUCores": bson.M{
			"bsonType": "number",
		},
		"CPUThreads": bson.M{
			"bsonType": "number",
		},
		"Socket": bson.M{
			"bsonType": "number",
		},
		"Type": bson.M{
			"bsonType": "string",
		},
		"Virtual": bson.M{
			"bsonType": "bool",
		},
		"Kernel": bson.M{
			"bsonType": "string",
		},
		"OS": bson.M{
			"bsonType": "string",
		},
		"MemoryTotal": bson.M{
			"bsonType": "double",
		},
		"SwapTotal": bson.M{
			"bsonType": "double",
		},
		"OracleCluster": bson.M{
			"bsonType": "bool",
		},
		"VeritasCluster": bson.M{
			"bsonType": "bool",
		},
		"SunCluster": bson.M{
			"bsonType": "bool",
		},
		"AixCluster": bson.M{
			"bsonType": "bool",
		},
	},
}
