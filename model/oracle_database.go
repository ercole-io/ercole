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
	InstanceNumber    int                               `json:"instanceNumber" bson:"instanceNumber"`
	InstanceName      string                            `json:"instanceName" bson:"instanceName"`
	Name              string                            `json:"name"`
	UniqueName        string                            `json:"uniqueName" bson:"uniqueName"`
	Status            string                            `json:"status"`
	IsCDB             bool                              `json:"isCDB" bson:"isCDB"`
	Version           string                            `json:"version"`
	Platform          string                            `json:"platform"`
	Archivelog        bool                              `json:"archivelog"`
	Charset           string                            `json:"charset"`
	NCharset          string                            `json:"nCharset" bson:"nCharset"`
	BlockSize         int                               `json:"blockSize" bson:"blockSize"`
	CPUCount          int                               `json:"cpuCount" bson:"cpuCount"`
	SGATarget         float64                           `json:"sgaTarget" bson:"sgaTarget"`
	PGATarget         float64                           `json:"pgaTarget" bson:"pgaTarget"`
	MemoryTarget      float64                           `json:"memoryTarget" bson:"memoryTarget"`
	SGAMaxSize        float64                           `json:"sgaMaxSize" bson:"sgaMaxSize"`
	SegmentsSize      float64                           `json:"segmentsSize" bson:"segmentsSize"`
	DatafileSize      float64                           `json:"datafileSize" bson:"datafileSize"`
	Allocated         float64                           `json:"allocated"`
	Elapsed           *float64                          `json:"elapsed"`
	DBTime            *float64                          `json:"dbTime" bson:"dbTime"`
	DailyCPUUsage     *float64                          `json:"dailyCPUUsage" bson:"dailyCPUUsage"`
	Work              *float64                          `json:"work"`
	ASM               bool                              `json:"asm"`
	Dataguard         bool                              `json:"dataguard"`
	Patches           []OracleDatabasePatch             `json:"patches"`
	Tablespaces       []OracleDatabaseTablespace        `json:"tablespaces"`
	Schemas           []OracleDatabaseSchema            `json:"schemas"`
	Licenses          []OracleDatabaseLicense           `json:"licenses"`
	ADDMs             []OracleDatabaseAddm              `json:"addms"`
	SegmentAdvisors   []OracleDatabaseSegmentAdvisor    `json:"segmentAdvisors" bson:"segmentAdvisors"`
	PSUs              []OracleDatabasePSU               `json:"psus"`
	Backups           []OracleDatabaseBackup            `json:"backups"`
	FeatureUsageStats []OracleDatabaseFeatureUsageStat  `json:"featureUsageStats" bson:"featureUsageStats"`
	PDBs              []OracleDatabasePluggableDatabase `json:"pdbs"`
	Services          []OracleDatabaseService           `json:"services"`
	OtherInfo         map[string]interface{}            `json:"-"`
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
		"instanceNumber",
		"instanceName",
		"name",
		"uniqueName",
		"status",
		"isCDB",
		"version",
		"platform",
		"archivelog",
		"charset",
		"nCharset",
		"blockSize",
		"cpuCount",
		"sgaTarget",
		"pgatarget",
		"memoryTarget",
		"sgamaxSize",
		"segmentsSize",
		"datafileSize",
		"allocated",
		"elapsed",
		"dbtime",
		"dailyCPUUsage",
		"work",
		"asm",
		"dataguard",
		"patches",
		"tablespaces",
		"schemas",
		"licenses",
		"addms",
		"segmentAdvisors",
		"psus",
		"backups",
		"featureUsageStats",
		"pdbs",
		"services",
	},
	"properties": bson.M{
		"instanceNumber": bson.M{
			"bsonType": "number",
			"minimum":  1,
		},
		"instanceName": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"name": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"uniqueName": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"status": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				"OPEN",
				"MOUNTED",
			},
		},
		"isCDB": bson.M{
			"bsonType": "bool",
		},
		"version": bson.M{
			"bsonType":  "string",
			"minLength": 8,
			"maxLength": 64,
		},
		"platform": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"archivelog": bson.M{
			"bsonType": "bool",
		},
		"charset": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"nCharset": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"blockSize": bson.M{
			"bsonType": "string",
			"minimum":  1,
		},
		"cpuCount": bson.M{
			"bsonType": "string",
			"minimum":  1,
		},
		"sgaTarget": bson.M{
			"bsonType": "number",
		},
		"pgatarget": bson.M{
			"bsonType": "number",
		},
		"memoryTarget": bson.M{
			"bsonType": "number",
		},
		"sgamaxSize": bson.M{
			"bsonType": "number",
		},
		"segmentsSize": bson.M{
			"bsonType": "number",
		},
		"datafileSize": bson.M{
			"bsonType": "number",
		},
		"allocated": bson.M{
			"bsonType": "number",
		},
		"elapsed": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"dbtime": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"dailyCPUUsage": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"work": bson.M{
			"anyOf": bson.A{
				bson.M{"type": "null"},
				bson.M{"type": "number"},
			},
		},
		"ASM": bson.M{
			"bsonType": "bool",
		},
		"dataguard": bson.M{
			"bsonType": "bool",
		},
		"patches": bson.M{
			"bsonType": "array",
			"items":    OracleDatabasePatchBsonValidatorRules,
		},
		"tablespaces": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseTablespaceBsonValidatorRules,
		},
		"schemas": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseSchemaBsonValidatorRules,
		},
		"licenses": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseLicenseBsonValidatorRules,
		},
		"addms": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseAddmBsonValidatorRules,
		},
		"segmentAdvisors": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseSegmentAdvisorBsonValidatorRules,
		},
		"psus": bson.M{
			"bsonType": "array",
			"items":    OracleDatabasePSUBsonValidatorRules,
		},
		"backups": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseBackupBsonValidatorRules,
		},
		"featureUsageStats": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseFeatureUsageStatBsonValidatorRules,
		},
		"services": bson.M{
			"bsonType": "array",
			"items":    OracleDatabaseServiceBsonValidatorRules,
		},
		"pdbs": bson.M{
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
