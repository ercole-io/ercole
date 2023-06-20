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

// StorageProvisioning holds various informations about the input/output operations per seconds of the host/db.
type StorageProvisioning struct {
	TimeStart  *time.Time `json:"timeStart" bson:"timeStart"`
	TimeEnd    *time.Time `json:"timeEnd" bson:"timeEnd"`
	CpuDbAvg   float64    `json:"cpuDbAvg" bson:"cpuDbAvg"`
	CpuDbMax   float64    `json:"cpuDbMax" bson:"cpuDbMax"`
	CpuHostAvg float64    `json:"cpuHostAvg" bson:"cpuHostAvg"`
	CpuHostMax float64    `json:"cpuHostMax" bson:"cpuHostMax"`
	IopsAvg    float64    `json:"iopsAvg" bson:"iopsAvg"`
	IopsMax    float64    `json:"iopsMax" bson:"iopsMax"`
	IombAvg    float64    `json:"iombAvg" bson:"iombAvg"`
	IombMax    float64    `json:"iombMax" bson:"iombMax"`
}

func (sp *StorageProvisioning) IsRange() bool {
	if sp.TimeStart != nil && sp.TimeEnd != nil {
		diff := sp.TimeEnd.Sub(*sp.TimeStart)

		return int64(diff.Hours()/24) > 1
	}

	return false
}
