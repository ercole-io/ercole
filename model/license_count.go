package model

import "go.mongodb.org/mongo-driver/bson"

type LicenseCount struct {
	Name  string `bson:"_id"`
	Count uint32
}

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
