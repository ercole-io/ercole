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

// OracleExadataComponent holds informations about a device in a exadata
type OracleExadataComponent struct {
	Hostname             string                  `bson:"Hostname"`
	ServerType           string                  `bson:"ServerType"`
	Model                string                  `bson:"Model"`
	SwVersion            string                  `bson:"SwVersion"`
	SwReleaseDate        string                  `bson:"SwReleaseDate"`
	RunningCPUCount      int                     `bson:"RunningCPUCount"`
	TotalCPUCount        int                     `bson:"TotalCPUCount"`
	Memory               int                     `bson:"Memory"`
	Status               string                  `bson:"Status"`
	RunningPowerSupply   int                     `bson:"RunningPowerSupply"`
	TotalPowerSupply     int                     `bson:"TotalPowerSupply"`
	PowerStatus          string                  `bson:"PowerStatus"`
	RunningFanCount      int                     `bson:"RunningFanCount"`
	TotalFanCount        int                     `bson:"TotalFanCount"`
	FanStatus            string                  `bson:"FanStatus"`
	TempActual           float32                 `bson:"TempActual"`
	TempStatus           string                  `bson:"TempStatus"`
	CellsrvServiceStatus string                  `bson:"CellsrvServiceStatus"`
	MsServiceStatus      string                  `bson:"MsServiceStatus"`
	RsServiceStatus      string                  `bson:"RsServiceStatus"`
	FlashcacheMode       string                  `bson:"FlashcacheMode"`
	CellDisks            []OracleExadataCellDisk `bson:"CellDisks"`
	OtherInfo            map[string]interface{}  `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleExadataComponent) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleExadataComponent) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleExadataComponent) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleExadataComponent) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// OracleExadataComponentBsonValidatorRules contains mongodb validation rules for OracleExadataComponentBsonValidatorRules
var OracleExadataComponentBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Hostname",
		"ServerType",
		"Model",
		"SwVersion",
		"SwReleaseDate",
		"RunningCPUCount",
		"TotalCPUCount",
		"Memory",
		"Status",
		"RunningPowerSupply",
		"TotalPowerSupply",
		"PowerStatus",
		"RunningFanCount",
		"TotalFanCount",
		"FanStatus",
		"TempActual",
		"TempStatus",
		"CellsrvServiceStatus",
		"MsServiceStatus",
		"RsServiceStatus",
		"FlashcacheMode",
		"CellDisks",
	},
	"properties": bson.M{
		"Hostname": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 253,
			"format":    "idn-hostname",
		},
		"ServerType": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"StorageServer",
				"DBServer",
				"IBSwitch",
			},
		},
		"Model": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"SwVersion": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"SwReleaseDate": bson.M{
			"bsonType": "string",
			"format":   "date",
		},
		"RunningCPUCount": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"TotalCPUCount": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"Memory": bson.M{
			"bsonType": "string",
			"minimum":  0,
			"$comment": "Memory in GB",
		},
		"Status": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"online",
				"offline",
			},
		},
		"RunningPowerSupply": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"TotalPowerSupply": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"PowerStatus": bson.M{
			"bsonType": "string",
		},
		"RunningFanCount": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"TotalFanCount": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"FanStatus": bson.M{
			"bsonType": "string",
		},
		"TempActual": bson.M{
			"bsonType": "number",
		},
		"TempStatus": bson.M{
			"bsonType": "string",
		},
		"CellsrvServiceStatus": bson.M{
			"bsonType": "string",
		},
		"MsServiceStatus": bson.M{
			"bsonType": "string",
		},
		"RsServiceStatus": bson.M{
			"bsonType": "string",
		},
		"FlashcacheMode": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"WriteBack",
				"WriteThrough",
			},
		},
		"CellDisks": bson.M{
			"bsonType": bson.A{"array", "null"},
			"items":    OracleExadataCellDiskBsonValidatorRules,
		},
	},
}
