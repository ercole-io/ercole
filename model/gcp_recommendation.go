// Copyright (c) 2024 Sorint.lab S.p.A.
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

type GcpRecommendation struct {
	SeqValue        uint64             `json:"seqValue" bson:"seqValue"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	ProfileID       primitive.ObjectID `json:"profileID" bson:"profileID"`
	ResourceID      uint64             `json:"resourceID" bson:"resourceID"`
	ResourceName    string             `json:"resourceName" bson:"resourceName"`
	Category        string             `json:"category" bson:"category"`
	Suggestion      string             `json:"suggestion" bson:"suggestion"`
	ProjectID       string             `json:"projectID" bson:"projectID"`
	ProjectName     string             `json:"projectName" bson:"projectName"`
	ObjectType      string             `json:"objectType" bson:"objectType"`
	Details         map[string]string  `json:"details" bson:"details"`
	ResolutionLevel string             `json:"resolutionLevel" bson:"resolutionLevel"`
}
