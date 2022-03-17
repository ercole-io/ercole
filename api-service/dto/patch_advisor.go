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
