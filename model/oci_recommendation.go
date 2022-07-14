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

import "time"

// Recommendation holds informations about a recommendation
type OciNativeRecommendation struct {
	TenancyOCID         string `json:"tenancyOCID"`
	Name                string `json:"name"`
	NumPending          string `json:"numPending"`
	EstimatedCostSaving string `json:"estimatedCostSaving"`
	Status              string `json:"status"`
	Importance          string `json:"importance"`
	RecommendationId    string `json:"recommendationId"`
}

type OciRecommendation struct {
	SeqValue        uint64      `json:"seqValue" bson:"seqValue"`
	ProfileID       string      `json:"profileID" bson:"profileID"`
	Category        string      `json:"category" bson:"category"`
	Suggestion      string      `json:"suggestion" bson:"suggestion"`
	CompartmentID   string      `json:"compartmentID" bson:"compartmentID"`
	CompartmentName string      `json:"compartmentName" bson:"compartmentName"`
	Name            string      `json:"name" bson:"name"`
	ResourceID      string      `json:"resourceID" bson:"resourceID"`
	ObjectType      string      `json:"objectType" bson:"objectType"`
	Details         []RecDetail `json:"details" bson:"details"`
	CreatedAt       time.Time   `json:"createdAt" bson:"createdAt"`
}

type RecDetail struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

const (
	OciBlockStorageRightsizing       = "Block Storage Rightsizing"
	OciComputeInstanceIdle           = "Compute Instance Rightsizing"
	OciInstanceRightsizing           = "Compute Instance Rightsizing"
	OciInstanceWithoutMonitoring     = "Compute Instance Without Monitoring"
	OciOldSnapshot                   = "Old Snapshot"
	OciUnusedResource                = "Unused Resource"
	OciUnusedStorage                 = "Unused Storage"
	OciSISRightsizing                = "Software Infracstructure Service Rightsizing"
	OciSISRightsizing1               = "Software Infracstructure Service Rightsizing 1"
	OciObjectStorageOptimization     = "Object Storage Optimization"
	OciUnusedServiceDecommisioning   = "Unused Service Decommisioning"
	OciComputeInstanceDecommisioning = "Compute Instance Decommisioning"
)

const (
	OciObjectTypeBlockStorage      = "Block Storage"
	OciObjectTypeComputeInstance   = "Compute Instance"
	OciObjectTypeDatabase          = "Database"
	OciObjectTypeLoadBalancer      = "Load Balancer"
	OciObjectTypeSnapshot          = "Snapshot"
	OciObjectTypeClusterKubernetes = "Cluster Kubernetes"
)

const (
	OciResizeOversizedBlockStorage      = "Resize oversized Block Storage"
	OciDeleteComputeInstanceNotActive   = "Delete Compute Instance not active"
	OciDeleteComputeInstanceNotUsed     = "Delete Compute Instance not used"
	OciResizeOversizedComputeInstance   = "Resize oversized compute instance"
	OciEnableBucketAutoTiering          = "Enable bucket auto tiering"
	OciDeleteSnapshotOlder              = "Delete snapshot older than 30 days"
	OciResizeOversizedDatabaseInstance  = "Resize oversized Database instance"
	OciResizeOversizedKubernetesCluster = "Resize oversized Kubernetes Cluster"
	OciDeleteLoadBalancerNotActive      = "Delete Load balancer not active"
	OciDeleteKubernetesNodeNotActive    = "Delete kubernetes node not active"
	OciDeleteKubernetesNodeNotUsed      = "Delete Kubernetes node not used"
	OciDeleteDatabaseInstanceNotActive  = "Delete Database instance not active"
	OciDeleteBlockStorageNotUsed        = "Delete Block Storage not used"
)
