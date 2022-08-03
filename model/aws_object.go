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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AwsObject struct {
	SeqValue     uint64             `json:"seqValue" bson:"seqValue"`
	ProfileID    primitive.ObjectID `json:"-" bson:"profileID"`
	ProfileName  string             `json:"profileName" bson:"profileName"`
	ObjectsCount []AwsObjectCount   `json:"ObjectsCount" bson:"ObjectsCount"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

type AwsObjectCount struct {
	Name  string `json:"name" bson:"name"`
	Count int    `json:"count" bson:"count"`
}
