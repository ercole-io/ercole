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

import "go.mongodb.org/mongo-driver/bson"

// ExadataDevice holds informations about a device in a exadata
type ExadataDevice struct {
	Hostname       string
	ServerType     string `bson:"server_type"`
	Model          string
	ExaSwVersion   string `bson:"exa_sw_version"`
	CPUEnabled     string `bson:"cpu_enabled"`
	Memory         string
	Status         string
	PowerCount     string            `bson:"power_count"`
	PowerStatus    string            `bson:"power_status"`
	FanCount       string            `bson:"fan_count"`
	FanStatus      string            `bson:"fan_status"`
	TempActual     string            `bson:"temp_actual"`
	TempStatus     string            `bson:"temp_status"`
	CellsrvService string            `bson:"cellsrv_service"`
	MsService      string            `bson:"ms_service"`
	RsService      string            `bson:"rs_service"`
	FlashcacheMode string            `bson:"flashcache_mode"`
	CellDisks      []ExadataCellDisk `bson:"cell_disks"`
}

// ExadataDeviceBsonValidatorRules contains mongodb validation rules for exadata device
var ExadataDeviceBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"hostname",
		"server_type",
		"model",
		"exa_sw_version",
		"cpu_enabled",
		"memory",
		"status",
		"power_count",
		"power_status",
		"fan_count",
		"fan_status",
		"temp_actual",
		"temp_status",
		"cellsrv_service",
		"ms_service",
		"rs_service",
		"flashcache_mode",
	}},
	{"properties", bson.D{
		{"hostname", bson.D{
			{"bsonType", "string"},
		}},
		{"server_type", bson.D{
			{"enum", bson.A{
				"StorageServer",
				"DBServer",
				"IBSwitch",
			}},
		}},
		{"model", bson.D{
			{"bsonType", "string"},
		}},
		{"exa_sw_version", bson.D{
			{"bsonType", "string"},
		}},
		{"cpu_enabled", bson.D{
			{"bsonType", "string"},
		}},
		{"memory", bson.D{
			{"bsonType", "string"},
		}},
		{"status", bson.D{
			{"bsonType", "string"},
		}},
		{"power_count", bson.D{
			{"bsonType", "string"},
		}},
		{"power_status", bson.D{
			{"bsonType", "string"},
		}},
		{"fan_count", bson.D{
			{"bsonType", "string"},
		}},
		{"fan_status", bson.D{
			{"bsonType", "string"},
		}},
		{"temp_actual", bson.D{
			{"bsonType", "string"},
		}},
		{"temp_status", bson.D{
			{"bsonType", "string"},
		}},
		{"cellsrv_service", bson.D{
			{"bsonType", "string"},
		}},
		{"ms_service", bson.D{
			{"bsonType", "string"},
		}},
		{"rs_service", bson.D{
			{"bsonType", "string"},
		}},
		{"flashcache_mode", bson.D{
			{"bsonType", "string"},
		}},
		{"cell_disks", bson.D{
			{"bsonType", bson.A{"array", "null"}},
			{"items", ExadataCellDiskBsonValidatorRules},
		}},
	}},
}
