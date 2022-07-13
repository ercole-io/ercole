package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AwsRecommendation struct {
	ID         primitive.ObjectID  `json:"id" bson:"_id"`
	SeqValue   uint64              `json:"seqValue" bson:"seqValue"`
	ProfileID  string              `json:"profileID" bson:"profileID"`
	Category   string              `json:"category" bson:"category"`
	Suggestion string              `json:"suggestion" bson:"suggestion"`
	Name       string              `json:"name" bson:"name"`
	ResourceID string              `json:"resourceID" bson:"resourceID"`
	ObjectType string              `json:"objectType" bson:"objectType"`
	Details    []map[string]string `json:"details" bson:"details"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
}
