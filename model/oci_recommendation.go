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

// Recommendation holds informations about a recommendation
type OciRecommendation struct {
	TenancyOCID         string `json:"tenancyOCID"`
	Name                string `json:"name"`
	NumPending          string `json:"numPending"`
	EstimatedCostSaving string `json:"estimatedCostSaving"`
	Status              string `json:"status"`
	Importance          string `json:"importance"`
	RecommendationId    string `json:"recommendationId"`
}

type OciErcoleRecommendation struct {
	Category        string      `json:"category"`
	Suggestion      string      `json:"suggestion"`
	CompartmentID   string      `json:"compartmentID"`
	CompartmentName string      `json:"compartmentName"`
	Name            string      `json:"name"`
	ResourceID      string      `json:"resourceID"`
	ObjectType      string      `json:"objectType"`
	Details         []RecDetail `json:"details"`
}

type RecDetail struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const (
	RecommendationTypeBlockStorageRightsizing       = "Block Storage Rightsizing"
	RecommendationTypeComputeInstanceIdle           = "Compute Instance Rightsizing"
	RecommendationTypeInstanceRightsizing           = "Compute Instance Rightsizing"
	RecommendationTypeInstanceWithoutMonitoring     = "Compute Instance Without Monitoring"
	RecommendationTypeOldSnapshot                   = "Old Snapshot"
	RecommendationTypeUnusedResource                = "Unused Resource"
	RecommendationTypeUnusedStorage                 = "Unused Storage"
	RecommendationTypeSISRightsizing                = "Software Infracstructure Service Rightsizing"
	RecommendationTypeSISRightsizing1               = "Software Infracstructure Service Rightsizing 1"
	RecommendationTypeObjectStorageOptimization     = "Object Storage Optimization"
	RecommendationTypeUnusedServiceDecommisioning   = "Unused Service Decommisioning"
	RecommendationTypeComputeInstanceDecommisioning = "Compute Instance Decommisioning"
)

const (
	ObjectTypeBlockStorage      = "Block Storage"
	ObjectTypeComputeInstance   = "Compute Instance"
	ObjectTypeDatabase          = "Database"
	ObjectTypeLoadBalancer      = "Load Balancer"
	ObjectTypeSnapshot          = "Snapshot"
	ObjectTypeClusterKubernetes = "Cluster Kubernetes"
)

const (
	ResizeOversizedBlockStorage = "Resize oversized Block Storage"
	DeleteComputeInstanceNotActive = "Delete Compute Instance not active"
	DeleteComputeInstanceNotUsed = "Delete Compute Instance not used"
	ResizeOversizedComputeInstance = "Resize oversized compute instance"
	EnableBucketAutoTiering = "Enable bucket auto tiering"
	DeleteSnapshotOlder = "Delete snapshot older than 30 days"
	ResizeOversizedDatabaseInstance = "Resize oversized Database instance"
	ResizeOversizedKubernetesCluster = "Resize oversized Kubernetes Cluster"
	DeleteLoadBalancerNotActive = "Delete Load balancer not active"
	DeleteKubernetesNodeNotActive = "Delete kubernetes node not active"
	DeleteKubernetesNodeNotUsed = "Delete Kubernetes node not used"
	DeleteDatabaseInstanceNotActive = "Delete Database instance not active"
	DeleteBlockStorageNotUsed = "Delete Block Storage not used"
)
