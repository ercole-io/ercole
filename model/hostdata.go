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
)

// HostData holds the whole information that will be sent to the server.
type HostData struct {
	Hostname              string
	Environment           string
	Location              string
	HostType              string `bson:"host_type"`
	Version               string
	HostDataSchemaVersion int    `bson:"-"`
	ServerVersion         string `bson:"server_version"`
	SchemaVersion         int    `bson:"schema_version"`
	Databases             string
	Schemas               string
	Info                  Host
	Extra                 ExtraInfo
	Archived              bool
	CreatedAt             time.Time `bson:"created_at"`
}

var HostDataBsonValidatorRules bson.D = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"hostname",
		"environment",
		"location",
		"host_type",
		// "version",
		"server_version",
		"schema_version",
		"info",
		"extra",
		"archived",
		"created_at",
	}},
	{"properties", bson.D{
		{"hostname", bson.D{
			{"bsonType", "string"},
		}},
		{"environment", bson.D{
			{"bsonType", "string"},
		}},
		{"location", bson.D{
			{"bsonType", "string"},
		}},
		{"host_type", bson.D{
			{"enum", bson.A{
				"oracledb",
				"virtualization",
			}},
		}},
		{"version", bson.D{
			{"bsonType", "string"},
		}},
		{"server_version", bson.D{
			{"bsonType", "string"},
		}},
		{"schema_version", bson.D{
			{"bsonType", "int"},
		}},
		{"databases", bson.D{
			{"bsonType", "string"},
		}},
		{"schemas", bson.D{
			{"bsonType", "string"},
		}},
		{"info", HostBsonValidatorRules},
		{"extra", ExtraInfoBsonValidatorRules},
		{"archived", bson.D{
			{"bsonType", "bool"},
		}},
		{"created_at", bson.D{
			{"bsonType", "date"},
		}},
	}},
}

const SchemaVersion int = 1
