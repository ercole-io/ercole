package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PatchAdvisor struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
	Date         primitive.DateTime `json:"date" bson:"date"`
	DbName       string             `json:"dbname" bson:"dbname"`
	Dbver        string             `json:"dbver" bson:"dbver"`
	Description  string             `json:"description" bson:"description"`
	Environment  string             `json:"environment" bson:"environment"`
	Hostname     string             `json:"hostname" bson:"hostname"`
	Location     string             `json:"location" bson:"location"`
	Status       string             `json:"status" bson:"status"`
	FourMonths   bool               `json:"fourMonths" bson:"fourMonths"`
	SixMonths    bool               `json:"sixMonths" bson:"sixMonths"`
	TwelveMonths bool               `json:"twelveMonths" bson:"twelveMonths"`
}

type PatchAdvisors []PatchAdvisor

type PatchAdvisorResponse struct {
	Content  PatchAdvisors  `json:"content" bson:"content"`
	Metadata PagingMetadata `json:"metadata" bson:"metadata"`
}
