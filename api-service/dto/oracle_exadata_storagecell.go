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

type OracleExadataStorageCell struct {
	Type               string                  `json:"type"`
	Hostname           string                  `json:"hostname"`
	CellDisk           string                  `json:"cellDisk"`
	Cell               string                  `json:"cell"`
	Size               string                  `json:"size"`
	FreeSpace          string                  `json:"freeSpace"`
	Status             string                  `json:"status"`
	ErrorCount         int                     `json:"errorCount"`
	GridDisks          []OracleExadataGridDisk `json:"gridDisks,omitempty"`
	Databases          []OracleExadataDatabase `json:"databases"`
	FreeSizePercentage string                  `json:"freeSizePercentage"`
}

func ToOracleExadataStorageCell(d domain.OracleExadataStorageCell) (*OracleExadataStorageCell, error) {
	res := &OracleExadataStorageCell{
		Type:               d.Type,
		Hostname:           d.Hostname,
		CellDisk:           d.CellDisk,
		Cell:               d.Cell,
		Status:             d.Status,
		ErrorCount:         d.ErrorCount,
		FreeSizePercentage: d.FreeSizePercentage,
	}

	if d.Size != nil {
		hsize, err := d.Size.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.Size = hsize
	}

	if d.FreeSpace != nil {
		hfreespace, err := d.FreeSpace.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.FreeSpace = hfreespace
	}

	griddisks, err := domain.ToUpperLevelLayers[domain.OracleExadataGridDisk, OracleExadataGridDisk](d.GridDisks, ToOracleExadataGridDisk)
	if err != nil {
		return nil, err
	}

	res.GridDisks = griddisks

	databases, err := domain.ToUpperLevelLayers[domain.OracleExadataDatabase, OracleExadataDatabase](d.Databases, ToOracleExadataDatabase)
	if err != nil {
		return nil, err
	}

	res.Databases = databases

	return res, nil
}
