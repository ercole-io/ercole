// Copyright (c) 2024 Sorint.lab S.p.A.
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

import "github.com/ercole-io/ercole/v2/api-service/domain"

type OracleExadataVM struct {
	Type         string `json:"type"`
	PhysicalHost string `json:"physicalHost"`
	Status       string `json:"status"`
	Name         string `json:"name"`
	CPUCurrent   int    `json:"cpuCurrent"`
	CPURestart   int    `json:"cpuRestart"`
	RamCurrent   int    `json:"ramCurrent"`
	RamRestart   int    `json:"ramRestart"`
	CPUOnline    int    `json:"cpuOnline"`
	CPUMaxUsable int    `json:"cpuMaxUsable"`
	RamOnline    int    `json:"ramOnline"`
	RamMaxUsable int    `json:"ramMaxUsable"`
	ClusterName  string `json:"clusterName"`
}

func ToOracleExadataVM(d domain.OracleExadataVM) (*OracleExadataVM, error) {
	res := &OracleExadataVM{
		Type:         d.Type,
		PhysicalHost: d.PhysicalHost,
		Status:       d.Status,
		Name:         d.Name,
		CPUCurrent:   d.CPUCurrent,
		CPURestart:   d.CPURestart,
		CPUOnline:    d.CPUOnline,
		CPUMaxUsable: d.CPUMaxUsable,
		ClusterName:  d.ClusterName,
	}

	if d.RamCurrent != nil {
		rramcurrent, err := d.RamCurrent.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.RamCurrent = rramcurrent
	}

	if d.RamRestart != nil {
		rramrestart, err := d.RamRestart.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.RamRestart = rramrestart
	}

	if d.RamOnline != nil {
		rramonline, err := d.RamOnline.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.RamOnline = rramonline
	}

	if d.RamMaxUsable != nil {
		rrammaxusable, err := d.RamMaxUsable.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.RamMaxUsable = rrammaxusable
	}

	return res, nil
}
