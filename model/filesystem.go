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

// Filesystem holds information about mounted filesystem and used space
type Filesystem struct {
	Filesystem     string `json:"filesystem" bson:"filesystem"`
	Type           string `json:"type" bson:"type"`
	Size           int64  `json:"size" bson:"size"`                     // in kilobytes
	UsedSpace      int64  `json:"usedSpace" bson:"usedSpace"`           // in kilobytes
	AvailableSpace int64  `json:"availableSpace" bson:"availableSpace"` // in kilobytes
	MountedOn      string `json:"mountedOn" bson:"mountedOn"`
}
