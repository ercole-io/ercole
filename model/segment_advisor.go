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

import (
	"reflect"

	godynstruct "github.com/amreo/go-dyn-struct"
	"go.mongodb.org/mongo-driver/bson"
)

// SegmentAdvisor holds information about a segment advisor
type SegmentAdvisor struct {
	SegmentOwner   string                 `bson:"SegmentOwner"`
	SegmentName    string                 `bson:"SegmentName"`
	SegmentType    string                 `bson:"SegmentType"`
	PartitionName  string                 `bson:"PartitionName"`
	Reclaimable    string                 `bson:"Reclaimable"`
	Recommendation string                 `bson:"Recommendation"`
	_otherInfo     map[string]interface{} `bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v SegmentAdvisor) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *SegmentAdvisor) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// MarshalBSON return the BSON rappresentation of this
func (v SegmentAdvisor) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v._otherInfo, "_otherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *SegmentAdvisor) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v._otherInfo)
}

// SegmentAdvisorBsonValidatorRules contains mongodb validation rules for segmentAdvisor
var SegmentAdvisorBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"SegmentOwner",
		"SegmentName",
		"SegmentType",
		"PartitionName",
		"Reclaimable",
		"Recommendation",
	},
	"properties": bson.M{
		"SegmentOwner": bson.M{
			"bsonType": "string",
		},
		"SegmentName": bson.M{
			"bsonType": "string",
		},
		"SegmentType": bson.M{
			"bsonType": "string",
		},
		"PartitionName": bson.M{
			"bsonType": "string",
		},
		"Reclaimable": bson.M{
			"bsonType": "string",
		},
		"Recommendation": bson.M{
			"bsonType": "string",
		},
	},
}
