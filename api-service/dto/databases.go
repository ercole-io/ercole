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

package dto

type Database struct {
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Version          string  `json:"version"`
	Hostname         string  `json:"hostname"`
	Environment      string  `json:"environment"`
	Charset          string  `json:"charset"`
	Memory           float64 `json:"memory"`       // in GB
	DatafileSize     float64 `json:"datafileSize"` // in GB
	SegmentsSize     float64 `json:"segmentSize"`  // in GB
	Archivelog       bool    `json:"archivelog"`
	HighAvailability bool    `json:"highAvailability"`
	DisasterRecovery bool    `json:"disasterRecovery"`
}

type DatabasesStatistics struct {
	TotalMemorySize   float64 `json:"total-memory-size"`   // in GB
	TotalSegmentsSize float64 `json:"total-segments-size"` // in GB
}
