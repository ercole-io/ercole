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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alert holds informations about a alert
type Alert struct {
	ID                      primitive.ObjectID     `json:"id" bson:"_id"`
	AlertCategory           string                 `json:"alertCategory" bson:"alertCategory"`
	AlertAffectedTechnology *string                `json:"alertAffectedTechnology" bson:"alertAffectedTechnology"`
	AlertCode               string                 `json:"alertCode" bson:"alertCode"`
	AlertSeverity           string                 `json:"alertSeverity" bson:"alertSeverity"`
	AlertStatus             string                 `json:"alertStatus" bson:"alertStatus"`
	Description             string                 `json:"description" bson:"description"`
	Date                    time.Time              `json:"date" bson:"date"`
	OtherInfo               map[string]interface{} `json:"otherInfo" bson:"otherInfo"`
}

const (
	AlertCategoryEngine  string = "ENGINE"
	AlertCategoryAgent   string = "AGENT"
	AlertCategoryLicense string = "LICENSE"
)

func getAlertCategories() []string {
	return []string{AlertCategoryEngine, AlertCategoryAgent, AlertCategoryLicense}
}

const (
	// ENGINE

	AlertCodeNewServer               string = "NEW_SERVER"
	AlertCodeUnlistedRunningDatabase string = "UNLISTED_RUNNING_DATABASE"
	AlertCodeMissingPrimaryDatabase  string = "MISSING_PRIMARY_DATABASE"

	// AGENT

	AlertCodeNoData string = "NO_DATA"

	// LICENSE

	AlertCodeNewDatabase string = "NEW_DATABASE"
	AlertCodeNewLicense  string = "NEW_LICENSE"
	AlertCodeNewOption   string = "NEW_OPTION"
)

func getAlertCodes() []string {
	return []string{
		AlertCodeNewServer, AlertCodeUnlistedRunningDatabase, AlertCodeMissingPrimaryDatabase,
		AlertCodeNoData,
		AlertCodeNewDatabase, AlertCodeNewLicense, AlertCodeNewOption,
	}
}

const (
	AlertSeverityInfo     string = "INFO"
	AlertSeverityWarning  string = "WARNING"
	AlertSeverityCritical string = "CRITICAL"
)

func getAlertSeverities() []string {
	return []string{AlertSeverityInfo, AlertSeverityWarning, AlertSeverityCritical}
}

// Alert status
const (
	// New contains string "NEW"
	AlertStatusNew string = "NEW"
	// Ack contains string "ACK"
	AlertStatusAck string = "ACK"
)

func getAlertStatuses() []string {
	return []string{AlertStatusNew, AlertStatusAck}
}

func (alert Alert) IsValid() bool {
	fields := make(map[string][]string)
	fields[alert.AlertCategory] = getAlertCategories()
	fields[alert.AlertCode] = getAlertCodes()
	fields[alert.AlertSeverity] = getAlertSeverities()
	fields[alert.AlertStatus] = getAlertStatuses()

fields:
	for thisValue, allValidValues := range fields {
		for _, validValue := range allValidValues {
			if thisValue == validValue {
				continue fields
			}
		}

		return false
	}

	return true
}
