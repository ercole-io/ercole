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

// OracleDatabase holds information about a oracle database.
type OracleDatabase struct {
	InstanceNumber    int                               `bson:"InstanceNumber"`
	Name              string                            `bson:"Name"`
	UniqueName        string                            `bson:"UniqueName"`
	Status            string                            `bson:"Status"`
	IsCDB             bool                              `bson:"IsCDB"`
	Version           string                            `bson:"Version"`
	Platform          string                            `bson:"Platform"`
	Archivelog        bool                              `bson:"Archivelog"`
	Charset           string                            `bson:"Charset"`
	NCharset          string                            `bson:"NCharset"`
	BlockSize         int                               `bson:"BlockSize"`
	CPUCount          int                               `bson:"CPUCount"`
	SGATarget         float32                           `bson:"SGATarget"`
	PGATarget         float32                           `bson:"PGATarget"`
	MemoryTarget      float32                           `bson:"MemoryTarget"`
	SGAMaxSize        float32                           `bson:"SGAMaxSize"`
	SegmentsSize      float32                           `bson:"SegmentsSize"`
	DatafileSize      float32                           `bson:"DatafileSize"`
	Allocated         float32                           `bson:"Allocated"`
	Elapsed           *float32                          `bson:"Elapsed"`
	DBTime            *float32                          `bson:"DBTime"`
	DailyCPUUsage     *float32                          `bson:"DailyCPUUsage"`
	Work              *float32                          `bson:"Work"`
	ASM               bool                              `bson:"ASM"`
	Dataguard         bool                              `bson:"Dataguard"`
	Patches           []OracleDatabasePatch             `bson:"Patches"`
	Tablespaces       []OracleDatabaseTablespace        `bson:"Tablespaces"`
	Schemas           []OracleDatabaseSchema            `bson:"Schemas"`
	Licenses          []OracleDatabaseLicense           `bson:"Licenses"`
	ADDMs             []OracleDatabaseAddm              `bson:"ADDMs"`
	SegmentAdvisors   []OracleDatabaseSegmentAdvisor    `bson:"SegmentAdvisors"`
	PSUs              []OracleDatabasePSU               `bson:"PSUs"`
	Backups           []OracleDatabaseBackup            `bson:"Backups"`
	FeatureUsageStats []OracleDatabaseFeatureUsageStat  `bson:"FeatureUsageStats"`
	PDBs              []OracleDatabasePluggableDatabase `bson:"PDBs"`
	Services          []OracleDatabaseService           `bson:"Services"`
	OtherInfo         map[string]interface{}            `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleDatabase) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleDatabase) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleDatabase) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleDatabase) UnmarshalBSON(data []byte) error {
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
func DatabasesArrayAsMap(dbs []OracleDatabase) map[string]OracleDatabase {
	out := make(map[string]OracleDatabase)
	for _, db := range dbs {
		out[db.Name] = db
	}
	return out
}

// HasEnterpriseLicense return true if the database has enterprise license.
func HasEnterpriseLicense(db OracleDatabase) bool {
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
