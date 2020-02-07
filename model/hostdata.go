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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SchemaVersion contains the version of the schema
const SchemaVersion int = 1

// HostData holds all informations about a host & services
type HostData struct {
	ID            primitive.ObjectID `bson:"_id"`
	Hostname      string             `bson:"Hostname"`
	Environment   string             `bson:"Environment"`
	Location      string             `bson:"Location"`
	Version       string             `bson:"Version"`
	ServerVersion string             `bson:"ServerVersion"`
	SchemaVersion int                `bson:"SchemaVersion"`
	Info          Host               `bson:"Info"`
	Extra         ExtraInfo          `bson:"Extra"`
	Archived      bool               `bson:"Archived"`
	CreatedAt     time.Time          `bson:"CreatedAt"`
}

// HostDataBsonValidatorRules contains mongodb validation rules for hostData
var HostDataBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
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
	}},
	{"properties", bson.D{
		{"Hostname", bson.D{
			{"bsonType", "string"},
		}},
		{"Environment", bson.D{
			{"bsonType", "string"},
		}},
		{"Location", bson.D{
			{"bsonType", "string"},
		}},
		{"Version", bson.D{
			{"bsonType", "string"},
		}},
		{"ServerVersion", bson.D{
			{"bsonType", "string"},
		}},
		{"SchemaVersion", bson.D{
			{"bsonType", "number"},
		}},
		{"Info", HostBsonValidatorRules},
		{"Extra", ExtraInfoBsonValidatorRules},
		{"Archived", bson.D{
			{"bsonType", "bool"},
		}},
		{"CreatedAt", bson.D{
			{"bsonType", "date"},
		}},
	}},
}
