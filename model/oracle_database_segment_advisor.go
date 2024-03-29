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

// OracleDatabaseSegmentAdvisor holds information about a segment advisor
type OracleDatabaseSegmentAdvisor struct {
	SegmentOwner   string  `json:"segmentOwner" bson:"segmentOwner"`
	SegmentName    string  `json:"segmentName" bson:"segmentName"`
	SegmentType    string  `json:"segmentType" bson:"segmentType"`
	PartitionName  string  `json:"partitionName" bson:"partitionName"`
	Reclaimable    float64 `json:"reclaimable" bson:"reclaimable"`
	Recommendation string  `json:"recommendation" bson:"recommendation"`
}
