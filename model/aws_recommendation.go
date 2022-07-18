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
)
