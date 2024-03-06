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

type OracleExadataStorageCell struct {
	Type               string
	Hostname           string
	CellDisk           string
	Cell               string
	Size               *OracleExadataMeasurement
	FreeSpace          *OracleExadataMeasurement
	UsedSize           *OracleExadataMeasurement
	Status             string
	ErrorCount         int
	GridDisks          []OracleExadataGridDisk
	Databases          []OracleExadataDatabase
	UsedSizePercentage string
}

func ToOracleExadataStorageCell(m model.OracleExadataStorageCell) (*OracleExadataStorageCell, error) {
	res := &OracleExadataStorageCell{
		Type:       m.Type,
		Hostname:   m.Hostname,
		CellDisk:   m.CellDisk,
		Cell:       m.Cell,
		Status:     m.Status,
		ErrorCount: m.ErrorCount,
	}

	size, err := StringToOracleExadataMeasurement(m.Size)
	if err != nil {
		return nil, err
	}

	res.Size = size

	freespace, err := StringToOracleExadataMeasurement(m.FreeSpace)
	if err != nil {
		return nil, err
	}

	res.FreeSpace = freespace

	usedsize := OracleExadataMeasurement{
		Symbol:   res.Size.Symbol,
		Quantity: res.Size.Quantity,
	}

	usedsize.Sub(*res.FreeSpace)
	res.UsedSize = &usedsize

	griddisks, err := ToUpperLevelLayers[model.OracleExadataGridDisk, OracleExadataGridDisk](m.GridDisks, ToOracleExadataGridDisk)
	if err != nil {
		return nil, err
	}

	res.GridDisks = griddisks

	databases, err := ToUpperLevelLayers[model.OracleExadataDatabase, OracleExadataDatabase](m.Databases, ToOracleExadataDatabase)
	if err != nil {
		return nil, err
	}

	res.Databases = databases

	res.UsedSizePercentage = GetPercentage(*res.UsedSize, *res.Size)

	return res, nil
}
