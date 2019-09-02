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

// Database holds information about the database.
type Database struct {
	InstanceNumber  string
	Name            string
	UniqueName      string
	Status          string
	Version         string
	Platform        string
	Archivelog      string
	Charset         string
	NCharset        string
	BlockSize       string
	CPUCount        string
	SGATarget       string
	PGATarget       string
	MemoryTarget    string
	SGAMaxSize      string
	SegmentsSize    string
	Used            string
	Allocated       string
	Elapsed         string
	DBTime          string
	Work            string
	ASM             bool
	Dataguard       bool
	Patches         []Patch
	Tablespaces     []Tablespace
	Schemas         []Schema
	Features        []Feature
	Licenses        []License
	ADDMs           []Addm
	SegmentAdvisors []SegmentAdvisor
	LastPSUs        []PSU
	Backups         []Backup
}

var DatabaseBsonValidatorRules bson.D = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"instance_number",
		"name",
		"unique_name",
		"status",
		"version",
		"platform",
		"archive_log",
		"charset",
		"ncharset",
		"block_size",
		"cpu_count",
		"sga_target",
		"pga_target",
		"memory_target",
		"sga_max_size",
		"segments_size",
		"used",
		"allocated",
		"elapsed",
		"db_time",
		"work",
		"asm",
		"dataguard",
		"patches",
		"tablespaces",
		"schemas",
		"features",
		"licenses",
		"addms",
		"segment_advisor",
		"last_psus",
		"backups",
	}},
	{"properties", bson.D{
		{"instance_number", bson.D{
			{"bsonType", "string"},
		}},
		{"name", bson.D{
			{"bsonType", "string"},
		}},
		{"unique_name", bson.D{
			{"bsonType", "string"},
		}},
		{"status", bson.D{
			{"bsonType", "string"},
		}},
		{"version", bson.D{
			{"bsonType", "string"},
		}},
		{"platform", bson.D{
			{"bsonType", "string"},
		}},
		{"archive_log", bson.D{
			{"bsonType", "string"},
		}},
		{"charset", bson.D{
			{"bsonType", "string"},
		}},
		{"ncharset", bson.D{
			{"bsonType", "string"},
		}},
		{"block_size", bson.D{
			{"bsonType", "string"},
		}},
		{"cpu_count", bson.D{
			{"bsonType", "string"},
		}},
		{"sga_target", bson.D{
			{"bsonType", "string"},
		}},
		{"pga_target", bson.D{
			{"bsonType", "string"},
		}},
		{"memory_target", bson.D{
			{"bsonType", "string"},
		}},
		{"sga_max_size", bson.D{
			{"bsonType", "string"},
		}},
		{"segments_size", bson.D{
			{"bsonType", "string"},
		}},
		{"used", bson.D{
			{"bsonType", "string"},
		}},
		{"allocated", bson.D{
			{"bsonType", "string"},
		}},
		{"elapsed", bson.D{
			{"bsonType", "string"},
		}},
		{"dbtime", bson.D{
			{"bsonType", "string"},
		}},
		{"work", bson.D{
			{"bsonType", "string"},
		}},
		{"asm", bson.D{
			{"bsonType", "bool"},
		}},
		{"dataguard", bson.D{
			{"bsonType", "bool"},
		}},
		{"patches", bson.D{
			{"bsonType", "array"},
			{"items", PatchBsonValidatorRules},
		}},
		{"tablespaces", bson.D{
			{"bsonType", "array"},
			{"items", TablespaceBsonValidatorRules},
		}},
		{"schemas", bson.D{
			{"bsonType", "array"},
			{"items", SchemaBsonValidatorRules},
		}},
		{"features", bson.D{
			{"bsonType", "array"},
			{"items", FeatureBsonValidatorRules},
		}},
		{"licenses", bson.D{
			{"bsonType", "array"},
			{"items", LicenseBsonValidatorRules},
		}},
		{"addms", bson.D{
			{"bsonType", "array"},
			{"items", AddmBsonValidatorRules},
		}},
		{"segment_advisor", bson.D{
			{"bsonType", "array"},
			{"items", SegmentAdvisorBsonValidatorRules},
		}},
		{"last_psus", bson.D{
			{"bsonType", "array"},
			{"items", PSUBsonValidatorRules},
		}},
		{"backups", bson.D{
			{"bsonType", "array"},
			{"items", BackupBsonValidatorRules},
		}},
	}},
}
