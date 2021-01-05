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

import "time"

type SearchHostsFilters struct {
	Search      []string
	SortBy      string
	SortDesc    bool
	Location    string
	Environment string
	OlderThan   time.Time
	PageNumber  int
	PageSize    int

	Hostname                      string
	Database                      string
	Technology                    string
	HardwareAbstractionTechnology string
	Cluster                       *string
	VirtualizationNode            string
	OperatingSystem               string
	Kernel                        string
	LTEMemoryTotal                float64
	GTEMemoryTotal                float64
	LTESwapTotal                  float64
	GTESwapTotal                  float64
	IsMemberOfCluster             *bool
	CPUModel                      string
	LTECPUCores                   int
	GTECPUCores                   int
	LTECPUThreads                 int
	GTECPUThreads                 int
}
