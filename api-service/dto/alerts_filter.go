// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlertsFilter struct {
	IDs                     []primitive.ObjectID   `json:"ids" bson:"_id,omitempty"`
	AlertCategory           *string                `json:"alertCategory" bson:"alertCategory,omitempty"`
	AlertAffectedTechnology *string                `json:"alertAffectedTechnology" bson:"alertAffectedTechnology,omitempty"`
	AlertCode               *string                `json:"alertCode" bson:"alertCode,omitempty"`
	AlertSeverity           *string                `json:"alertSeverity" bson:"alertSeverity,omitempty"`
	AlertStatus             *string                `json:"alertStatus" bson:"alertStatus,omitempty"`
	Description             *string                `json:"description" bson:"description,omitempty"`
	Date                    time.Time              `json:"date" bson:"date,omitempty"`
	OtherInfo               map[string]interface{} `json:"otherInfo" bson:"otherInfo,omitempty"`
}
