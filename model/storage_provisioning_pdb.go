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

// StorageProvisioningPdb holds various informations about the input/output operations per seconds of the pluggable database.
type StorageProvisioningPdb struct {
	TimeStart *time.Time `json:"timeStart" bson:"timeStart"`
	TimeEnd   *time.Time `json:"timeEnd" bson:"timeEnd"`
	CpuDbAvg  float64    `json:"cpuDbAvg" bson:"cpuDbAvg"`
	CpuDbMax  float64    `json:"cpuDbMax" bson:"cpuDbMax"`
	IopsAvg   float64    `json:"iopsAvg" bson:"iopsAvg"`
	IombAvg   float64    `json:"iombAvg" bson:"iombAvg"`
	IombMax   float64    `json:"iombMax" bson:"iombMax"`
}
