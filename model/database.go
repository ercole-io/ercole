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

// Database holds information about a database.
type Database struct {
	InstanceNumber  string           `bson:"InstanceNumber"`
	Name            string           `bson:"Name"`
	UniqueName      string           `bson:"UniqueName"`
	Status          string           `bson:"Status"`
	Version         string           `bson:"Version"`
	Platform        string           `bson:"Platform"`
	Archivelog      string           `bson:"Archivelog"`
	Charset         string           `bson:"Charset"`
	NCharset        string           `bson:"NCharset"`
	BlockSize       string           `bson:"BlockSize"`
	CPUCount        string           `bson:"CPUCount"`
	SGATarget       string           `bson:"SGATarget"`
	PGATarget       string           `bson:"PGATarget"`
	MemoryTarget    string           `bson:"MemoryTarget"`
	SGAMaxSize      string           `bson:"SGAMaxSize"`
	SegmentsSize    string           `bson:"SegmentsSize"`
	Used            string           `bson:"Used"`
	Allocated       string           `bson:"Allocated"`
	Elapsed         string           `bson:"Elapsed"`
	DBTime          string           `bson:"DBTime"`
	DailyCPUUsage   string           `bson:"DailyCPUUsage"`
	Work            string           `bson:"Work"`
	ASM             bool             `bson:"ASM"`
	Dataguard       bool             `bson:"Dataguard"`
	Patches         []Patch          `bson:"Patches"`
	Tablespaces     []Tablespace     `bson:"Tablespaces"`
	Schemas         []Schema         `bson:"Schemas"`
	Features        []Feature        `bson:"Features"`
	Licenses        []License        `bson:"Licenses"`
	ADDMs           []Addm           `bson:"ADDMs"`
	SegmentAdvisors []SegmentAdvisor `bson:"SegmentAdvisors"`
	LastPSUs        []PSU            `bson:"LastPSUs"`
	Backups         []Backup         `bson:"Backups"`
}

// DatabaseBsonValidatorRules contains mongodb validation rules for database
var DatabaseBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"InstanceNumber",
		"Name",
		"UniqueName",
		"Status",
		"Version",
		"Platform",
		"Archivelog",
		"Charset",
		"NCharset",
		"BlockSize",
		"CPUCount",
		"SGATarget",
		"PGATarget",
		"MemoryTarget",
		"SGAMaxSize",
		"SegmentsSize",
		"Used",
		"Allocated",
		"Elapsed",
		"DBTime",
		"DailyCPUUsage",
		"Work",
		"ASM",
		"Dataguard",
		"Patches",
		"Tablespaces",
		"Schemas",
		"Features",
		"Licenses",
		"ADDMs",
		"SegmentAdvisors",
		"LastPSUs",
		"Backups",
	}},
	{"properties", bson.D{
		{"InstanceNumber", bson.D{
			{"bsonType", "string"},
		}},
		{"Name", bson.D{
			{"bsonType", "string"},
		}},
		{"UniqueName", bson.D{
			{"bsonType", "string"},
		}},
		{"Status", bson.D{
			{"bsonType", "string"},
		}},
		{"Version", bson.D{
			{"bsonType", "string"},
		}},
		{"Platform", bson.D{
			{"bsonType", "string"},
		}},
		{"Archivelog", bson.D{
			{"bsonType", "string"},
		}},
		{"Charset", bson.D{
			{"bsonType", "string"},
		}},
		{"NCharset", bson.D{
			{"bsonType", "string"},
		}},
		{"BlockSize", bson.D{
			{"bsonType", "string"},
		}},
		{"CPUCount", bson.D{
			{"bsonType", "string"},
		}},
		{"SGATarget", bson.D{
			{"bsonType", "string"},
		}},
		{"PGATarget", bson.D{
			{"bsonType", "string"},
		}},
		{"MemoryTarget", bson.D{
			{"bsonType", "string"},
		}},
		{"SGAMaxSize", bson.D{
			{"bsonType", "string"},
		}},
		{"SegmentsSize", bson.D{
			{"bsonType", "string"},
		}},
		{"Used", bson.D{
			{"bsonType", "string"},
		}},
		{"Allocated", bson.D{
			{"bsonType", "string"},
		}},
		{"Elapsed", bson.D{
			{"bsonType", "string"},
		}},
		{"DBTime", bson.D{
			{"bsonType", "string"},
		}},
		{"DailyCPUUsage", bson.D{
			{"bsonType", "string"},
		}},
		{"Work", bson.D{
			{"bsonType", "string"},
		}},
		{"ASM", bson.D{
			{"bsonType", "bool"},
		}},
		{"Dataguard", bson.D{
			{"bsonType", "bool"},
		}},
		{"Patches", bson.D{
			{"bsonType", "array"},
			{"items", PatchBsonValidatorRules},
		}},
		{"Tablespaces", bson.D{
			{"bsonType", "array"},
			{"items", TablespaceBsonValidatorRules},
		}},
		{"Schemas", bson.D{
			{"bsonType", "array"},
			{"items", SchemaBsonValidatorRules},
		}},
		{"Features", bson.D{
			{"bsonType", "array"},
			{"items", FeatureBsonValidatorRules},
		}},
		{"Licenses", bson.D{
			{"bsonType", "array"},
			{"items", LicenseBsonValidatorRules},
		}},
		{"ADDMs", bson.D{
			{"bsonType", "array"},
			{"items", AddmBsonValidatorRules},
		}},
		{"SegmentAdvisors", bson.D{
			{"bsonType", "array"},
			{"items", SegmentAdvisorBsonValidatorRules},
		}},
		{"LastPSUs", bson.D{
			{"bsonType", "array"},
			{"items", PSUBsonValidatorRules},
		}},
		{"Backups", bson.D{
			{"bsonType", "array"},
			{"items", BackupBsonValidatorRules},
		}},
	}},
}

// DatabasesArrayAsMap return the equivalent map of the database array with Database.Name as Key
func DatabasesArrayAsMap(dbs []Database) map[string]Database {
	out := make(map[string]Database)
	for _, db := range dbs {
		out[db.Name] = db
	}
	return out
}

// HasEnterpriseLicense return true if the database has enterprise license.
func HasEnterpriseLicense(db Database) bool {
	//The database may not support the "license" feature
	if db.Licenses == nil {
		return false
	}

	//Search for a enterprise license
	for _, lic := range db.Licenses {
		if (lic.Name == "Oracle ENT" || lic.Name == "oracle ENT" || lic.Name == "Oracle EXT" || lic.Name == "oracle EXT") && lic.Count > 0 {
			return true
		}
	}

	return false
}
