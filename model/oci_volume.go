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

// OciVolume holds informations about an Oracle Cloud Volume
type OciVolume struct {
	CompartmentID      string `json:"compartmentID"`
	ResourceID         string `json:"resourceID"`
	Name               string `json:"name"`
	Size               string `json:"size"`
	VpusPerGB          string `json:"vpuPerGB"`
	AvailabilityDomain string `json:"availabilityDomain"`
	State              string `json:"availabilityDomain"`
}

type OciResourcePerformance struct {
	ResourceID string  `json:"resourceID"`
	Name       string  `json:"name"`
	Size       int     `json:"size"`
	VpusPerGB  int     `json:"vpuPerGB"`
	Throughput float64 `json:"throughput"`
	Iops       int     `json:"iops"`
}
