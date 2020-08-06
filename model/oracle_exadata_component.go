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
	Hostname             string                   `json:"hostname" bson:"hostname"`
	ServerType           string                   `json:"serverType" bson:"serverType"`
	Model                string                   `json:"model" bson:"model"`
	SwVersion            string                   `json:"swVersion" bson:"swVersion"`
	SwReleaseDate        string                   `json:"swReleaseDate" bson:"swReleaseDate"`
	RunningCPUCount      *int                     `json:"runningCPUCount" bson:"runningCPUCount"`
	TotalCPUCount        *int                     `json:"totalCPUCount" bson:"totalCPUCount"`
	Memory               *int                     `json:"memory" bson:"memory"`
	Status               *string                  `json:"status" bson:"status"`
	RunningPowerSupply   *int                     `json:"runningPowerSupply" bson:"runningPowerSupply"`
	TotalPowerSupply     *int                     `json:"totalPowerSupply" bson:"totalPowerSupply"`
	PowerStatus          *string                  `json:"powerStatus" bson:"powerStatus"`
	RunningFanCount      *int                     `json:"runningFanCount" bson:"runningFanCount"`
	TotalFanCount        *int                     `json:"totalFanCount" bson:"totalFanCount"`
	FanStatus            *string                  `json:"fanStatus" bson:"fanStatus"`
	TempActual           *float64                 `json:"tempActual" bson:"tempActual"`
	TempStatus           *string                  `json:"tempStatus" bson:"tempStatus"`
	CellsrvServiceStatus *string                  `json:"cellsrvServiceStatus" bson:"cellsrvServiceStatus"`
	MsServiceStatus      *string                  `json:"msServiceStatus" bson:"msServiceStatus"`
	RsServiceStatus      *string                  `json:"rsServiceStatus" bson:"rsServiceStatus"`
	FlashcacheMode       *string                  `json:"flashcacheMode" bson:"flashcacheMode"`
	CellDisks            *[]OracleExadataCellDisk `json:"cellDisks" bson:"cellDisks"`
	OtherInfo            map[string]interface{}   `json:"-" bson:"-"`
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
		"hostname",
		"serverType",
		"model",
		"swVersion",
		"swReleaseDate",
		"runningCPUCount",
		"totalCPUCount",
		"memory",
		"status",
		"runningPowerSupply",
		"totalPowerSupply",
		"powerStatus",
		"runningFanCount",
		"totalFanCount",
		"fanStatus",
		"tempActual",
		"tempStatus",
		"cellsrvServiceStatus",
		"msServiceStatus",
		"rsServiceStatus",
		"flashcacheMode",
		"cellDisks",
	},
	"properties": bson.M{
		"hostname": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 253,
			"pattern":   `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-_]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-_]*[A-Za-z0-9])$`,
		},
		"serverType": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"StorageServer",
				"DBServer",
				"IBSwitch",
			},
		},
		"model": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"swVersion": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"swReleaseDate": bson.M{
			"bsonType": "string",
			"pattern":  "[0-9]{4}-[0-9]{2}-[0-9]{2}",
		},
		"runningCPUCount": bson.M{
			"bsonType": bson.A{"null", "number"},
			"minimum":  1,
		},
		"totalCPUCount": bson.M{
			"bsonType": bson.A{"null", "number"},
			"minimum":  1,
		},
		"memory": bson.M{
			"bsonType": bson.A{"null", "number"},
			"minimum":  0,
		},
		"status": bson.M{
			"anyOf": bson.A{
				bson.M{"bsonType": "null"},
				bson.M{
					"bsonType": "string",
					"enum": bson.A{
						"online",
						"offline",
					},
				},
			},
		},
		"runningPowerSupply": bson.M{
			"bsonType": bson.A{"null", "number"},
			"minimum":  1,
		},
		"totalPowerSupply": bson.M{
			"bsonType": bson.A{"null", "number"},
			"minimum":  1,
		},
		"powerStatus": bson.M{
			"bsonType": bson.A{"null", "string"},
		},
		"runningFanCount": bson.M{
			"bsonType": bson.A{"null", "number"},
			"minimum":  1,
		},
		"totalFanCount": bson.M{
			"bsonType": bson.A{"null", "number"},
			"minimum":  1,
		},
		"fanStatus": bson.M{
			"bsonType": bson.A{"null", "string"},
		},
		"tempActual": bson.M{
			"bsonType": bson.A{"null", "number"},
		},
		"tempStatus": bson.M{
			"bsonType": bson.A{"null", "string"},
		},
		"cellsrvServiceStatus": bson.M{
			"bsonType": bson.A{"null", "string"},
		},
		"msServiceStatus": bson.M{
			"bsonType": bson.A{"null", "string"},
		},
		"rsServiceStatus": bson.M{
			"bsonType": bson.A{"null", "string"},
		},
		"flashcacheMode": bson.M{
			"anyOf": bson.A{
				bson.M{"bsonType": "null"},
				bson.M{
					"bsonType": "string",
					"enum": bson.A{
						"WriteBack",
						"WriteThrough",
					},
				},
			},
		},
		"cellDisks": bson.M{
			"bsonType": bson.A{"array", "null"},
			"items":    OracleExadataCellDiskBsonValidatorRules,
		},
	},
}
