package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Alert struct {
	ID            string `bson:"_id"`
	AlertCode     string
	AlertSeverity string
	AlertStatus   string
	Description   string
	Date          time.Time
}

const NEW_DATABASE string = "NEW_DATABASE"
const NEW_OPTION string = "NEW_OPTION"
const NEW_LICENSE string = "NEW_LICENSE"
const NEW_SERVER string = "NEW_SERVER"
const NO_DATA string = "NO_DATA"

const MINOR string = "MINOR"
const WARNING string = "WARNING"
const MAJOR string = "MAJOR"
const CRITICAL string = "CRITICAL"
const NOTICE string = "NOTICE"

var AlertBsonValidatorRules bson.D = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"_id",
		"alert_code",
		"alert_severity",
		"alert_status",
		"description",
		"date",
	}},
	{"properties", bson.D{
		{"_id", bson.D{
			{"bsonType", "string"},
		}},
		{"alert_code", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				"NEW_DATABASE",
				"NEW_OPTION",
				"NEW_LICENSE",
				"NEW_SERVER",
				"NO_DATA",
			}},
		}},
		{"alert_severity", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				"MINOR",
				"WARNING",
				"MAJOR",
				"CRITICAL",
				"NOTICE",
			}},
		}},
		{"alert_status", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				"NEW",
				"ACK",
			}},
		}},
		{"description", bson.D{
			{"bsonType", "string"},
		}},
		{"date", bson.D{
			{"bsonType", "date"},
		}},
	}},
}
