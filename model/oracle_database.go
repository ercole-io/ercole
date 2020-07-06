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
	InstanceName      string                            `bson:"InstanceName"`
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
	SGATarget         float64                           `bson:"SGATarget"`
	PGATarget         float64                           `bson:"PGATarget"`
	MemoryTarget      float64                           `bson:"MemoryTarget"`
	SGAMaxSize        float64                           `bson:"SGAMaxSize"`
	SegmentsSize      float64                           `bson:"SegmentsSize"`
	DatafileSize      float64                           `bson:"DatafileSize"`
	Allocated         float64                           `bson:"Allocated"`
	Elapsed           *float64                          `bson:"Elapsed"`
	DBTime            *float64                          `bson:"DBTime"`
	DailyCPUUsage     *float64                          `bson:"DailyCPUUsage"`
	Work              *float64                          `bson:"Work"`
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

// OracleDatabaseBsonValidatorRules contains mongodb validation rules for OracleDatabase
var OracleDatabaseBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"InstanceNumber",
		"InstanceName",
		"Name",
		"UniqueName",
		"Status",
		"IsCDB",
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
		"DatafileSize",
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
		"PSUs",
		"Backups",
		"FeatureUsageStats",
		"PDBs",
		"Services",
	},
	"properties": bson.M{
		"InstanceNumber": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"InstanceName": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"Name": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"UniqueName": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"Status": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"OPEN",
				"MOUNTED",
			},
		},
		"IsCDB": bson.M{
			"bsonType": "bool",
		},
		"Version": bson.M{
			"bsonType":  "string",
			"minLength": 8,
			"maxLength": 64,
		},
		"Platform": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"Archivelog": bson.M{
			"bsonType": "bool",
		},
		"Charset": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"NCharset": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"BlockSize": bson.M{
			"bsonType": "string",
			"minimum":  1,
		},
		"CPUCount": bson.M{
			"bsonType": "string",
			"minimum":  1,
		},
		"SGATarget": bson.M{
			"bsonType": "number",
		},
		"PGATarget": bson.M{
			"bsonType": "number",
		},
		"MemoryTarget": bson.M{
			"bsonType": "number",
		},
		"SGAMaxSize": bson.M{
			"bsonType": "number",
		},
		"SegmentsSize": bson.M{
			"bsonType": "number",
		},
		"DatafileSize": bson.M{
			"bsonType": "number",
		},
		"Allocated": bson.M{
			"bsonType": "number",
		},
		"Elapsed": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"DBTime": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"DailyCPUUsage": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"Work": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"ASM": bson.M{
			"bsonType": "bool",
		},
		"Dataguard": bson.M{
			"bsonType": "bool",
		},
		"Patches": bson.M{
			"bsonType": "array",
			"items":    OracleDatabasePatchBsonValidatorRules,
		},
		"Tablespaces": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseTablespaceBsonValidatorRules,
		},
		"Schemas": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseSchemaBsonValidatorRules,
		},
		"Licenses": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseLicenseBsonValidatorRules,
		},
		"ADDMs": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseAddmBsonValidatorRules,
		},
		"SegmentAdvisors": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseSegmentAdvisorBsonValidatorRules,
		},
		"PSUs": bson.M{
			"bsonType": "array",
			"items":    OracleDatabasePSUBsonValidatorRules,
		},
		"Backups": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseBackupBsonValidatorRules,
		},
		"FeatureUsageStats": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseFeatureUsageStatBsonValidatorRules,
		},
		"Services": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseServiceBsonValidatorRules,
		},
		"PDBs": bson.M{
			"bsonType": "array",
			"items":    OracleDatabasePluggableDatabaseBsonValidatorRules,
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
