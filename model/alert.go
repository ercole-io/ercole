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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alert holds informations about a alert
type Alert struct {
	ID                      primitive.ObjectID     `json:"id" bson:"_id"`
	AlertCategory           string                 `json:"alertCategory bson:alertCategory"`
	AlertAffectedTechnology *string                `json:"alertAffectedTechnology bson:alertAffectedTechnology"`
	AlertCode               string                 `json:"alertCode bson:alertCode"`
	AlertSeverity           string                 `json:"alertSeverity bson:alertSeverity"`
	AlertStatus             string                 `json:"alertStatus bson:alertStatus"`
	Description             string                 `json:"description"`
	Date                    time.Time              `json:"date"`
	OtherInfo               map[string]interface{} `json:"otherInfo bson:otherInfo"`
}

// Alert codes
const (
	// AlertCategoryEngine contains string "ENGINE"
	AlertCategoryEngine string = "ENGINE"
	// AlertCategoryAgent contains string "AGENT"
	AlertCategoryAgent string = "AGENT"
	// AlertCategoryLicense contains string "LICENSE"
	AlertCategoryLicense string = "LICENSE"
)

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
	// AlertSeverityWarning contains string "WARNING"
	AlertSeverityWarning string = "WARNING"
	// AlertSeverityCritical contains string "CRITICAL"
	AlertSeverityCritical string = "CRITICAL"
	// AlertSeverityInfo contains string "INFO"
	AlertSeverityInfo string = "INFO"
)

// Alert status
const (
	// New contains string "NEW"
	AlertStatusNew string = "NEW"
	// Ack contains string "ACK"
	AlertStatusAck string = "ACK"
)

// AlertBsonValidatorRules contains mongodb validation rules for alert
var AlertBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"_id",
		"alertCategory",
		"alertAffectedTechnology",
		"alertCode",
		"alertSeverity",
		"alertStatus",
		"description",
		"date",
	},
	"properties": bson.M{
		"alertCategory": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				AlertCategoryEngine,
				AlertCategoryAgent,
				AlertCategoryLicense,
			},
		},
		"alertAffectedTechnology": bson.M{
			"bsonType": bson.A{"null", "string"},
			"enum": bson.A{
				nil,
				TechnologyOracleDatabase,
				TechnologyOracleExadata,
			},
		},
		"alertCode": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				AlertCodeNewDatabase,
				AlertCodeNewOption,
				AlertCodeNewLicense,
				AlertCodeNewServer,
				AlertCodeNoData,
			},
		},
		"alertSeverity": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				AlertSeverityWarning,
				AlertSeverityCritical,
				AlertSeverityInfo,
			},
		},
		"alertStatus": bson.M{
			"bsonType": "string",
			"enum": bson.A{
				AlertStatusNew,
				AlertStatusAck,
			},
		},
		"description": bson.M{
			"bsonType": "string",
		},
		"date": bson.M{
			"bsonType": "date",
		},
		"OtherInfo": bson.M{
			"bsonType": "object",
		},
	},
}
