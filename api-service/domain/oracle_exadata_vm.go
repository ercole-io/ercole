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

package domain

import (
	"github.com/ercole-io/ercole/v2/model"
)

type OracleExadataVM struct {
	Type         string
	PhysicalHost string
	Status       string
	Name         string
	CPUCurrent   int
	CPURestart   int
	RamCurrent   *OracleExadataMeasurement
	RamRestart   *OracleExadataMeasurement
	CPUOnline    int
	CPUMaxUsable int
	RamOnline    *OracleExadataMeasurement
	RamMaxUsable *OracleExadataMeasurement
	ClusterName  string
}

func ToOracleExadataVm(m model.OracleExadataVM) (*OracleExadataVM, error) {
	res := &OracleExadataVM{
		Type:         m.Type,
		PhysicalHost: m.PhysicalHost,
		Status:       m.Status,
		Name:         m.Name,
		CPUCurrent:   m.CPUCurrent,
		CPURestart:   m.CPURestart,
		CPUOnline:    m.CPUOnline,
		CPUMaxUsable: m.CPUMaxUsable,
		ClusterName:  m.ClusterName,
	}

	ramcurrent, err := IntToOracleExadataMeasurement(m.RamCurrent, "MB")
	if err != nil {
		return nil, err
	}

	res.RamCurrent = ramcurrent

	ramrestart, err := IntToOracleExadataMeasurement(m.RamRestart, "MB")
	if err != nil {
		return nil, err
	}

	res.RamRestart = ramrestart

	ramonline, err := IntToOracleExadataMeasurement(m.RamOnline, "MB")
	if err != nil {
		return nil, err
	}

	res.RamOnline = ramonline

	rammaxusable, err := IntToOracleExadataMeasurement(m.RamMaxUsable, "MB")
	if err != nil {
		return nil, err
	}

	res.RamMaxUsable = rammaxusable

	return res, nil
}
