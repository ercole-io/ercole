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

// Host contains info about the host
type Host struct {
	Hostname       string
	Environment    string
	Location       string
	CPUModel       string `bson:"cpu_model"`
	CPUCores       int    `bson:"cpu_cores"`
	CPUThreads     int    `bson:"cpu_threads"`
	Socket         int
	Type           string
	Virtual        bool
	Kernel         string
	OS             string
	MemoryTotal    int  `bson:"memory_total"`
	SwapTotal      int  `bson:"swap_total"`
	OracleCluster  bool `bson:"oracle_cluster"`
	VeritasCluster bool `bson:"veritas_cluster"`
	SunCluster     bool `bson:"sun_cluster"`
}

// HostBsonValidatorRules contains mongodb validation rules for host
var HostBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"hostname",
		"environment",
		"location",
		"cpu_model",
		"cpu_cores",
		"cpu_threads",
		"socket",
		"type",
		"virtual",
		"kernel",
		"os",
		"memory_total",
		"swap_total",
		"oracle_cluster",
		"veritas_cluster",
		"sun_cluster",
	}},
	{"properties", bson.D{
		{"hostname", bson.D{
			{"bsonType", "string"},
			{"pattern", "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"},
		}},
		{"environment", bson.D{
			{"bsonType", "string"},
		}},
		{"location", bson.D{
			{"bsonType", "string"},
		}},
		{"cpu_model", bson.D{
			{"bsonType", "string"},
		}},
		{"cpu_cores", bson.D{
			{"bsonType", "int"},
		}},
		{"cpu_threads", bson.D{
			{"bsonType", "int"},
		}},
		{"socket", bson.D{
			{"bsonType", "int"},
		}},
		{"type", bson.D{
			{"bsonType", "string"},
		}},
		{"virtual", bson.D{
			{"bsonType", "bool"},
		}},
		{"kernel", bson.D{
			{"bsonType", "string"},
		}},
		{"os", bson.D{
			{"bsonType", "string"},
		}},
		{"memory_total", bson.D{
			{"bsonType", "int"},
		}},
		{"swap_total", bson.D{
			{"bsonType", "int"},
		}},
		{"oracle_cluster", bson.D{
			{"bsonType", "bool"},
		}},
		{"veritas_cluster", bson.D{
			{"bsonType", "bool"},
		}},
		{"sun_cluster", bson.D{
			{"bsonType", "bool"},
		}},
	}},
}
