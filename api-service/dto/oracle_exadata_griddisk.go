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

import (
	"strconv"

	"github.com/ercole-io/ercole/v2/api-service/domain"
)

type OracleExadataGridDisk struct {
	Type          string `json:"type"`
	Hostname      string `json:"hostname"`
	GridDisk      string `json:"gridDisk"`
	CellDisk      string `json:"cellDisk"`
	Size          int    `json:"size"`
	Status        string `json:"status"`
	ErrorCount    int    `json:"errorCount"`
	CachingPolicy string `json:"cachingPolicy"`
	AsmDiskName   string `json:"asmDiskName"`
	AsmDiskGroup  string `json:"asmDiskGroup"`
	AsmDiskSize   string `json:"asmDiskSize"`
	AsmDiskStatus string `json:"asmDiskStatus"`
}

func ToOracleExadataGridDisk(d domain.OracleExadataGridDisk) (*OracleExadataGridDisk, error) {
	res := &OracleExadataGridDisk{
		Type:          d.Type,
		Hostname:      d.Hostname,
		GridDisk:      d.GridDisk,
		CellDisk:      d.CellDisk,
		Status:        d.Status,
		ErrorCount:    d.ErrorCount,
		CachingPolicy: d.CachingPolicy,
		AsmDiskName:   d.AsmDiskName,
		AsmDiskGroup:  d.AsmDiskGroup,
		AsmDiskStatus: d.AsmDiskStatus,
	}

	if d.Size != nil {
		rsize, err := d.Size.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.Size = rsize
	}

	if d.AsmDiskSize != nil {
		rasmdisksize, err := d.AsmDiskSize.RoundedGiB()
		if err != nil {
			if err.Error() == domain.UNKNOWN_VALUE {
				res.AsmDiskSize = domain.UNKNOWN_VALUE
			}
		} else {
			res.AsmDiskSize = strconv.Itoa(rasmdisksize)
		}
	}

	return res, nil
}
