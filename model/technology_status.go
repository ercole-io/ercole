// Copyright (c) 2022 Sorint.lab S.p.A.
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

// TechnologyStatus represent usage status of a technology
type TechnologyStatus struct {
	Product            string  `json:"product"`
	ConsumedByHosts    float64 `json:"-"`
	CoveredByContracts float64 `json:"-"`
	TotalCost          float64 `json:"-"`
	PaidCost           float64 `json:"-"`
	Compliance         float64 `json:"compliance"`
	UnpaidDues         float64 `json:"unpaidDues"`
	HostsCount         int     `json:"hostsCount"`
	NewAvgPercentage   float64 `json:"newAvgPercentage"`
}
