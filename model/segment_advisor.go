// Copyright (c) 2019 Sorint.lab S.p.A.
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

import "go.mongodb.org/mongo-driver/bson"

// SegmentAdvisor holds information about a segment advisor
type SegmentAdvisor struct {
	SegmentOwner   string `bson:"SegmentOwner"`
	SegmentName    string `bson:"SegmentName"`
	SegmentType    string `bson:"SegmentType"`
	PartitionName  string `bson:"PartitionName"`
	Reclaimable    string `bson:"Reclaimable"`
	Recommendation string `bson:"Recommendation"`
}

// SegmentAdvisorBsonValidatorRules contains mongodb validation rules for segmentAdvisor
var SegmentAdvisorBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"SegmentOwner",
		"SegmentName",
		"SegmentType",
		"PartitionName",
		"Reclaimable",
		"Recommendation",
	}},
	{"properties", bson.D{
		{"SegmentOwner", bson.D{
			{"bsonType", "string"},
		}},
		{"SegmentName", bson.D{
			{"bsonType", "string"},
		}},
		{"SegmentType", bson.D{
			{"bsonType", "string"},
		}},
		{"PartitionName", bson.D{
			{"bsonType", "string"},
		}},
		{"Reclaimable", bson.D{
			{"bsonType", "string"},
		}},
		{"Recommendation", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
