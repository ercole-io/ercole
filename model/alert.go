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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alert holds informations about a alert
type Alert struct {
	ID            primitive.ObjectID `bson:"_id"`
	AlertCode     string             `bson:"alert_code"`
	AlertSeverity string             `bson:"alert_severity"`
	AlertStatus   string             `bson:"alert_status"`
	Description   string
	Date          time.Time
	OtherInfo     map[string]interface{} `bson:"other_info"`
}

// Alert codes
const (
	// NewDatabase contains string "NEW_DATABASE"
	AlertCodeNewDatabase string = "NEW_DATABASE"
	// NewOption contains string "NEW_OPTION"
	AlertCodeNewOption string = "NEW_OPTION"
	// NewLicense contains string "NEW_LICENSE"
	AlertCodeNewLicense string = "NEW_LICENSE"
	// NewServer contains string "NEW_SERVER"
	AlertCodeNewServer string = "NEW_SERVER"
	// NoData contains string "NO_DATA"
	AlertCodeNoData string = "NO_DATA"
)

// Alert severity
const (
	// Minor contains string "MINOR"
	AlertSeverityMinor string = "MINOR"
	// Warning contains string "WARNING"
	AlertSeverityWarning string = "WARNING"
	// Major contains string "MAJOR"
	AlertSeverityMajor string = "MAJOR"
	// Critical contains string "CRITICAL"
	AlertSeverityCritical string = "CRITICAL"
	// Notice contains string "NOTICE"
	AlertSeverityNotice string = "NOTICE"
)

// Alert status
const (
	// New contains string "NEW"
	AlertStatusNew string = "NEW"
	// Ack contains string "ACK"
	AlertStatusAck string = "ACK"
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
		{"alert_code", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				AlertCodeNewDatabase,
				AlertCodeNewOption,
				AlertCodeNewLicense,
				AlertCodeNewServer,
				AlertCodeNoData,
			}},
		}},
		{"alert_severity", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				AlertSeverityMinor,
				AlertSeverityWarning,
				AlertSeverityMajor,
				AlertSeverityCritical,
				AlertSeverityNotice,
			}},
		}},
		{"alert_status", bson.D{
			{"bsonType", "string"},
			{"enum", bson.A{
				AlertStatusNew,
				AlertStatusAck,
			}},
		}},
		{"description", bson.D{
			{"bsonType", "string"},
		}},
		{"date", bson.D{
			{"bsonType", "date"},
		}},
		{"other_info", bson.D{
			{"bsonType", "object"},
		}},
	}},
}
