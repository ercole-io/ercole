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
	SegmentOwner   string `bson:"segment_owner"`
	SegmentName    string `bson:"segment_name"`
	SegmentType    string `bson:"segment_type"`
	PartitionName  string `bson:"partition_name"`
	Reclaimable    string
	Recommendation string
}

// SegmentAdvisorBsonValidatorRules contains mongodb validation rules for segmentAdvisor
var SegmentAdvisorBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"segment_owner",
		"segment_name",
		"segment_type",
		"partition_name",
		"reclaimable",
		"recommendation",
	}},
	{"properties", bson.D{
		{"segment_owner", bson.D{
			{"bsonType", "string"},
		}},
		{"segment_name", bson.D{
			{"bsonType", "string"},
		}},
		{"segment_type", bson.D{
			{"bsonType", "string"},
		}},
		{"partition_name", bson.D{
			{"bsonType", "string"},
		}},
		{"reclaimable", bson.D{
			{"bsonType", "string"},
		}},
		{"recommendation", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
