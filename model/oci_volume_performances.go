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

type OciVolumePerformance struct {
	Vpu          int              `json:"vpu" bson:"vpu"`
	Performances []OciPerformance `json:"performances" bson:"performances"`
}

type OciPerformance struct {
	Size   int           `json:"size" bson:"size"`
	Values OciPerfValues `json:"values" bson:"values"`
}

type OciPerfValues struct {
	MaxThroughput float64 `json:"maxThroughput" bson:"maxThroughput"`
	MaxIOPS       int     `json:"maxIOPS" bson:"maxIOPS"`
}
