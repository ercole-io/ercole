package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Alert holds informations about a alert
type Alert struct {
	ID            string `bson:"_id"`
	HostName      string
	AlertCode     string `bson:"alert_code"`
	AlertSeverity string `bson:"alert_severity"`
	AlertStatus   string `bson:"alert_status"`
	Description   string
	Date          time.Time
}

// Alert codes
const (
	// NewDatabase contains string "NEW_DATABASE"
	NewDatabase string = "NEW_DATABASE"
	// NewOption contains string "NEW_OPTION"
	NewOption string = "NEW_OPTION"
	// NewLicense contains string "NEW_LICENSE"
	NewLicense string = "NEW_LICENSE"
	// NewServer contains string "NEW_SERVER"
	NewServer string = "NEW_SERVER"
	// NoData contains string "NO_DATA"
	NoData string = "NO_DATA"
)

// Alert severity
const (
	// Minor contains string "MINOR"
	Minor string = "MINOR"
	// Warning contains string "WARNING"
	Warning string = "WARNING"
	// Major contains string "MAJOR"
	Major string = "MAJOR"
	// Critical contains string "CRITICAL"
	Critical string = "CRITICAL"
	// Notice contains string "NOTICE"
	Notice string = "NOTICE"
)

// Alert status
const (
	// New contains string "NEW"
	New string = "NEW"
	// Ack contains string "ACK"
	Ack string = "ACK"
)

// AlertBsonValidatorRules contains mongodb validation rules for alert
var AlertBsonValidatorRules = bson.D{
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
				NewDatabase,
				NewOption,
				NewLicense,
				NewServer,
				NoData,
			}},
		}},
		{"alert_severity", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				Minor,
				Warning,
				Major,
				Critical,
				Notice,
			}},
		}},
		{"alert_status", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				New,
				Ack,
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
