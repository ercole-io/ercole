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

type MicrosoftSQLServerInstance struct {
	Status          string                       `json:"status" bson:"status"`
	Name            string                       `json:"name" bson:"name"`
	DisplayName     string                       `json:"displayName" bson:"displayName"`
	ServerName      string                       `json:"serverName" bson:"serverName"`
	DatabaseID      int                          `json:"databaseID" bson:"databaseID"`
	StateDesc       string                       `json:"stateDesc" bson:"stateDesc"`
	Version         string                       `json:"version" bson:"version"`
	Platform        string                       `json:"platform" bson:"platform"`
	RecoveryModel   string                       `json:"recoveryModel" bson:"recoveryModel"`
	CollationModel  string                       `json:"collationModel" bson:"collationModel"`
	BlockSize       int                          `json:"blockSize" bson:"blockSize"`
	SchedulersCount int                          `json:"schedulersCount" bson:"schedulersCount"`
	AffinityMask    int                          `json:"affinityMask" bson:"affinityMask"`
	MinServerMemory int                          `json:"minServerMemory" bson:"minServerMemory"`
	MaxServerMemory int                          `json:"maxServerMemory" bson:"maxServerMemory"`
	CTP             int                          `json:"ctp" bson:"ctp"`
	MaxDop          int                          `json:"maxDop" bson:"maxDop"`
	Alloc           float64                      `json:"alloc" bson:"alloc"`
	Edition         string                       `json:"edition" bson:"edition"`
	ProductCode     string                       `json:"productCode" bson:"productCode"`
	LicensingInfo   string                       `json:"licensingInfo" bson:"licensingInfo"`
	Databases       []MicrosoftSQLServerDatabase `json:"databases" bson:"databases"`
	Patches         []MicrosoftSQLServerPatch    `json:"patches" bson:"patches"`
	OtherInfo       map[string]interface{}       `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v MicrosoftSQLServerInstance) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerInstance) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v MicrosoftSQLServerInstance) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerInstance) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MicrosoftSQLServerInstanceBsonValidatorRules contains mongodb validation rules for MicrosoftSQLServerInstance
var MicrosoftSQLServerInstanceBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"status",
		"name",
		"displayName",
		"serverName",
		"databaseID",
		"databaseName",
		"stateDesc",
		"version",
		"platform",
		"recovertModel",
		"collationName",
		"blockSize",
		"schedulersCount",
		"affinityMask",
		"minServerMemory",
		"maxServerMemory",
		"ctp",
		"maxDop",
		"alloc",
		"edition",
		"productVersion",
		"editionType",
		"productCode",
		"licensingInfo",
		"databases",
		"patches",
	},
	"properties": bson.M{
		"status": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"Running",
				"Stopped",
				"ContinuePending",
				"Paused",
				"PausePending",
				"StartPending",
				"StopPending",
			},
		},
		"name": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"displayName": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"serverName": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"databaseID": bson.M{
			"bsonType": "integer",
			"minimum":  1,
		},
		"stateDesc": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"ONLINE",
				"RESTORING",
				"RECOVERING",
				"RECOVERY_PENDING",
				"SUSPECT",
				"EMERGENCY",
				"OFFLINE",
				"COPYING",
				"OFFLINE_SECONDARY",
			},
		},
		"version": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"platform": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"recoveryModel": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"FULL",
				"BULK_LOGGED",
				"SIMPLE",
			},
		},
		"collationName": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"blockSize": bson.M{
			"bsonType": "integer",
			"minimum":  1,
		},
		"schedulersCount": bson.M{
			"bsonType": "integer",
			"minimum":  1,
		},
		"affinityMask": bson.M{
			"bsonType": "integer",
			"minimum":  0,
		},
		"minServerMemory": bson.M{
			"bsonType": "integer",
			"minimum":  1,
		},
		"maxServerMemory": bson.M{
			"bsonType": "integer",
			"minimum":  1,
		},
		"ctp": bson.M{
			"bsonType": "integer",
			"minimum":  1,
		},
		"maxDop": bson.M{
			"bsonType": "integer",
			"minimum":  0,
		},
		"alloc": bson.M{
			"bsonType": "integer",
			"minimum":  0,
		},
		"edition": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"ENT",
				"EXP",
				"STD",
				"BI",
				"DEV",
				"WEB",
				"AZU",
			},
		},
		"productCode": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
			"pattern":   "^\\{[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}\\}$",
		},
		"licensingInfo": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 512,
		},
		"databases": bson.M{
			"bsonType": "array",
			"items":    MicrosoftSQLServerDatabaseBsonValidatorRules,
		},
		"patches": bson.M{
			"bsonType": "array",
			"items":    MicrosoftSQLServerPatchBsonValidatorRules,
		},
	},
}
