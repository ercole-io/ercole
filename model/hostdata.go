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
	"time"

	godynstruct "github.com/amreo/go-dyn-struct"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SchemaVersion contains the version of the schema
const SchemaVersion int = 1

// HostData holds all informations about a host & services
type HostData struct {
	ID            primitive.ObjectID     `bson:"_id"`
	Hostname      string                 `bson:"Hostname"`
	Environment   string                 `bson:"Environment"`
	Location      string                 `bson:"Location"`
	Version       string                 `bson:"Version"`
	ServerVersion string                 `bson:"ServerVersion"`
	SchemaVersion int                    `bson:"SchemaVersion"`
	Info          Host                   `bson:"Info"`
	Extra         ExtraInfo              `bson:"Extra"`
	Archived      bool                   `bson:"Archived"`
	CreatedAt     time.Time              `bson:"CreatedAt"`
	_otherInfo    map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v HostData) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *HostData) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v HostData) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *HostData) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// HostDataBsonValidatorRules contains mongodb validation rules for hostData
var HostDataBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"Hostname",
		"Environment",
		"Location",
		"Version",
		"ServerVersion",
		"SchemaVersion",
		"Info",
		"Extra",
		"Archived",
		"CreatedAt",
	},
	"properties": bson.M{
		"Hostname": bson.M{
			"bsonType": "string",
		},
		"Environment": bson.M{
			"bsonType": "string",
		},
		"Location": bson.M{
			"bsonType": "string",
		},
		"Version": bson.M{
			"bsonType": "string",
		},
		"ServerVersion": bson.M{
			"bsonType": "string",
		},
		"SchemaVersion": bson.M{
			"bsonType": "number",
		},
		"Info":  HostBsonValidatorRules,
		"Extra": ExtraInfoBsonValidatorRules,
		"Archived": bson.M{
			"bsonType": "bool",
		},
		"CreatedAt": bson.M{
			"bsonType": "date",
		},
	},
}
