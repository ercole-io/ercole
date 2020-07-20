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

type MicrosoftSQLServerDatabaseTablespace struct {
	Filename   string                 `json:"filename" bson:"filename"`
	Alloc      int                    `json:"alloc" bson:"alloc"`
	Used       int                    `json:"used" bson:"used"`
	Growth     float64                `json:"growth" bson:"growth"`
	GrowthUnit string                 `json:"growthUnit" bson:"growthUnit"`
	FileType   string                 `json:"fileType" bson:"fileType"`
	Status     string                 `json:"status" bson:"status"`
	OtherInfo  map[string]interface{} `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v MicrosoftSQLServerDatabaseTablespace) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerDatabaseTablespace) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v MicrosoftSQLServerDatabaseTablespace) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerDatabaseTablespace) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MicrosoftSQLServerDatabaseTablespaceBsonValidatorRules contains mongodb validation rules for MicrosoftSQLServerDatabaseTablespace
var MicrosoftSQLServerDatabaseTablespaceBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"filename",
		"alloc",
		"used",
		"growth",
		"growthUnit",
		"fileType",
		"status",
	},
	"properties": bson.M{
		"filename": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"alloc": bson.M{
			"bsonType": "integer",
			"minimum":  0,
		},
		"used": bson.M{
			"bsonType": "integer",
			"minimum":  0,
		},
		"growth": bson.M{
			"bsonType": "number",
		},
		"growthUnit": bson.M{
			"bsonType": "number",
			"enum": bson.A{
				"%",
				"MB",
			},
		},
		"fileType": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"ROWS",
				"LOG",
				"FILESTREAM",
				"FULLTEXT",
			},
		},
		"status": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"ONLINE",
				"RESTORING",
				"RECOVERING",
				"RECOVERY_PENDING",
				"SUSPECT",
				"OFFLINE",
				"DEFUNCT",
			},
		},
	},
}
