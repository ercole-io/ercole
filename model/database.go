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

// Database holds information about a database.
type Database struct {
	InstanceNumber  string                 `bson:"InstanceNumber"`
	Name            string                 `bson:"Name"`
	UniqueName      string                 `bson:"UniqueName"`
	Status          string                 `bson:"Status"`
	Version         string                 `bson:"Version"`
	Platform        string                 `bson:"Platform"`
	Archivelog      string                 `bson:"Archivelog"`
	Charset         string                 `bson:"Charset"`
	NCharset        string                 `bson:"NCharset"`
	BlockSize       string                 `bson:"BlockSize"`
	CPUCount        string                 `bson:"CPUCount"`
	SGATarget       string                 `bson:"SGATarget"`
	PGATarget       string                 `bson:"PGATarget"`
	MemoryTarget    string                 `bson:"MemoryTarget"`
	SGAMaxSize      string                 `bson:"SGAMaxSize"`
	SegmentsSize    string                 `bson:"SegmentsSize"`
	Used            string                 `bson:"Used"`
	Allocated       string                 `bson:"Allocated"`
	Elapsed         string                 `bson:"Elapsed"`
	DBTime          string                 `bson:"DBTime"`
	DailyCPUUsage   string                 `bson:"DailyCPUUsage"`
	Work            string                 `bson:"Work"`
	ASM             bool                   `bson:"ASM"`
	Dataguard       bool                   `bson:"Dataguard"`
	Patches         []Patch                `bson:"Patches"`
	Tablespaces     []Tablespace           `bson:"Tablespaces"`
	Schemas         []Schema               `bson:"Schemas"`
	Licenses        []License              `bson:"Licenses"`
	ADDMs           []Addm                 `bson:"ADDMs"`
	SegmentAdvisors []SegmentAdvisor       `bson:"SegmentAdvisors"`
	LastPSUs        []PSU                  `bson:"LastPSUs"`
	Backups         []Backup               `bson:"Backups"`
	OtherInfo       map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v Database) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *Database) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v Database) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *Database) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// DatabaseBsonValidatorRules contains mongodb validation rules for database
var DatabaseBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
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
		"Licenses",
		"ADDMs",
		"SegmentAdvisors",
		"LastPSUs",
		"Backups",
	},
	"properties": bson.M{
		"InstanceNumber": bson.M{
			"bsonType": "string",
		},
		"Name": bson.M{
			"bsonType": "string",
		},
		"UniqueName": bson.M{
			"bsonType": "string",
		},
		"Status": bson.M{
			"bsonType": "string",
		},
		"Version": bson.M{
			"bsonType": "string",
		},
		"Platform": bson.M{
			"bsonType": "string",
		},
		"Archivelog": bson.M{
			"bsonType": "string",
		},
		"Charset": bson.M{
			"bsonType": "string",
		},
		"NCharset": bson.M{
			"bsonType": "string",
		},
		"BlockSize": bson.M{
			"bsonType": "string",
		},
		"CPUCount": bson.M{
			"bsonType": "string",
		},
		"SGATarget": bson.M{
			"bsonType": "string",
		},
		"PGATarget": bson.M{
			"bsonType": "string",
		},
		"MemoryTarget": bson.M{
			"bsonType": "string",
		},
		"SGAMaxSize": bson.M{
			"bsonType": "string",
		},
		"SegmentsSize": bson.M{
			"bsonType": "string",
		},
		"Used": bson.M{
			"bsonType": "string",
		},
		"Allocated": bson.M{
			"bsonType": "string",
		},
		"Elapsed": bson.M{
			"bsonType": "string",
		},
		"DBTime": bson.M{
			"bsonType": "string",
		},
		"DailyCPUUsage": bson.M{
			"bsonType": "string",
		},
		"Work": bson.M{
			"bsonType": "string",
		},
		"ASM": bson.M{
			"bsonType": "bool",
		},
		"Dataguard": bson.M{
			"bsonType": "bool",
		},
		"Patches": bson.M{
			"bsonType": "array",
			"items":    PatchBsonValidatorRules,
		},
		"Tablespaces": bson.M{
			"bsonType": "array",
			"items":    TablespaceBsonValidatorRules,
		},
		"Schemas": bson.M{
			"bsonType": "array",
			"items":    SchemaBsonValidatorRules,
		},
		"Licenses": bson.M{
			"bsonType": "array",
			"items":    LicenseBsonValidatorRules,
		},
		"ADDMs": bson.M{
			"bsonType": "array",
			"items":    AddmBsonValidatorRules,
		},
		"SegmentAdvisors": bson.M{
			"bsonType": "array",
			"items":    SegmentAdvisorBsonValidatorRules,
		},
		"LastPSUs": bson.M{
			"bsonType": "array",
			"items":    PSUBsonValidatorRules,
		},
		"Backups": bson.M{
			"bsonType": "array",
			"items":    BackupBsonValidatorRules,
		},
	},
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
