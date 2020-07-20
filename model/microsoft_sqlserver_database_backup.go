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

type MicrosoftSQLServerDatabaseBackup struct {
	BackupType string                 `json:"backupType" bson:"backupType"`
	Hour       string                 `json:"hour" bson:"hour"`
	WeekDays   []string               `json:"weekDays" bson:"weekDays"`
	AvgBckSize float64                `json:"avgBckSize" bson:"avgBckSize"`
	OtherInfo  map[string]interface{} `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v MicrosoftSQLServerDatabaseBackup) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerDatabaseBackup) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v MicrosoftSQLServerDatabaseBackup) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerDatabaseBackup) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MicrosoftSQLServerDatabaseBackupBsonValidatorRules contains mongodb validation rules for MicrosoftSQLServerDatabaseBackup
var MicrosoftSQLServerDatabaseBackupBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"backupType",
		"hour",
		"weekDays",
		"avgBckSize",
	},
	"properties": bson.M{
		"backupType": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"Database",
				"Log",
				"File or filegroup",
				"Differential database",
				"Differential file",
				"Differential partial",
				"Partial",
			},
		},
		"hour": bson.M{
			"bsonType":  "string",
			"minLength": 5,
			"maxLength": 5,
			"pattern":   "^[0-9]{2}:[0-9]{2}$",
		},
		"weekDays": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"Monday",
				"Tuesday",
				"Wednesday",
				"Thursday",
				"Friday",
				"Saturday",
				"Sunday",
			},
			"uniqueItems": true,
		},
		"avgBckSize": bson.M{
			"bsonType": "number",
			"minimum":  0,
		},
	},
}
