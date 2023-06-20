// Copyright (c) 2023 Sorint.lab S.p.A.
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

import "time"

// Consumption holds all informations about cpu consumption of a specific host
type Consumption struct {
	TimeStart *time.Time `json:"timeStart" bson:"timeStart"`
	TimeEnd   *time.Time `json:"timeEnd" bson:"timeEnd"`
	CpuAvg    float64    `json:"cpuAvg" bson:"cpuAvg"`
}
