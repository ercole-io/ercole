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

// Backup holds informations about a backup
type Backup struct {
	BackupType string                 `bson:"BackupType"`
	Hour       string                 `bson:"Hour"`
	WeekDays   string                 `bson:"WeekDays"`
	AvgBckSize string                 `bson:"AvgBckSize"`
	Retention  string                 `bson:"Retention"`
	_otherInfo map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Backup) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Backup) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v Backup) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Backup) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// BackupBsonValidatorRules contains mongodb validation rules for backup
var BackupBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"BackupType",
		"Hour",
		"WeekDays",
		"AvgBckSize",
		"Retention",
	},
	"properties": bson.M{
		"BackupType": bson.M{
			"bsonType": "string",
		},
		"Hour": bson.M{
			"bsonType": "string",
		},
		"WeekDays": bson.M{
			"bsonType": "string",
		},
		"AvgBckSize": bson.M{
			"bsonType": "string",
		},
		"Retention": bson.M{
			"bsonType": "string",
		},
	},
}
