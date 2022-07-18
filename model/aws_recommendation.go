// Copyright (c) 2022 Sorint.lab S.p.A.
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
)

type AwsRecommendation struct {
	SeqValue   uint64                   `json:"seqValue" bson:"seqValue"`
	ProfileID  string                   `json:"profileID" bson:"profileID"`
	Category   string                   `json:"category" bson:"category"`
	Suggestion string                   `json:"suggestion" bson:"suggestion"`
	Name       string                   `json:"name" bson:"name"`
	ResourceID string                   `json:"resourceID" bson:"resourceID"`
	ObjectType string                   `json:"objectType" bson:"objectType"`
	Details    []map[string]interface{} `json:"details" bson:"details"`
	Errors     []map[string]string      `json:"errors" bson:"errors"`
	CreatedAt  time.Time                `json:"createdAt" bson:"createdAt"`
}

const (
	AwsUnusedResource                      = "Unused Resource"
	AwsObjectTypeLoadBalancer              = "Load Balancer"
	AwsPublicID                            = "Public IP"
	AwsDeleteLoadBalancerNotActive         = "Delete Load balancer not active"
	AwsDeletePublicIPAddressNotAssociated  = "Delete public IP address not associated"
	AwsObjectStorageOptimization           = "Object Storage Optimization"
	AwsObjectStorageOptimizationSuggestion = "Enable bucket auto tiering"
	AwsObjectStorageOptimizationType       = "Object Storage"
	AwsDeleteBlockStorageNotUsed           = "Delete Block Storage not used"
	AwsObjectVolume                        = "Object volume"
)
