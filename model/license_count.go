package model

import "go.mongodb.org/mongo-driver/bson"

// LicenseCount holds information about Oracle database license
type LicenseCount struct {
	Name  string `bson:"_id"`
	Count uint32
}

// LicenseCountBsonValidatorRules contains mongodb validation rules for licenseCount
var LicenseCountBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"_id",
		"count",
	}},
	{"properties", bson.D{
		{"_id", bson.D{
			{"bsonType", "string"},
		}},
		{"count", bson.D{
			{"bsonType", "int"},
		}},
	}},
}
