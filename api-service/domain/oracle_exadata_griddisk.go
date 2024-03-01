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

import "github.com/ercole-io/ercole/v2/model"

type OracleExadataGridDisk struct {
	Type          string
	Hostname      string
	GridDisk      string
	CellDisk      string
	Size          *OracleExadataMeasurement
	Status        string
	ErrorCount    int
	CachingPolicy string
	AsmDiskName   string
	AsmDiskGroup  string
	AsmDiskSize   *OracleExadataMeasurement
	AsmDiskStatus string
}

func ToOracleExadataGridDisk(m model.OracleExadataGridDisk) (*OracleExadataGridDisk, error) {
	res := &OracleExadataGridDisk{
		Type:          m.Type,
		Hostname:      m.Hostname,
		GridDisk:      m.GridDisk,
		CellDisk:      m.CellDisk,
		Status:        m.Status,
		ErrorCount:    m.ErrorCount,
		CachingPolicy: m.CachingPolicy,
		AsmDiskName:   m.AsmDiskName,
		AsmDiskGroup:  m.AsmDiskGroup,
		AsmDiskStatus: m.AsmDiskStatus,
	}

	size, err := StringToOracleExadataMeasurement(m.Size)
	if err != nil {
		return nil, err
	}

	res.Size = size

	asmdisksize, err := StringToOracleExadataMeasurement(m.AsmDiskSize)
	if err != nil {
		return nil, err
	}

	res.AsmDiskSize = asmdisksize

	return res, nil
}
