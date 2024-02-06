// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"errors"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

type OracleExadataStorageCell struct {
	Type       string                        `json:"type"`
	Hostname   string                        `json:"hostname"`
	CellDisk   string                        `json:"cellDisk"`
	Cell       string                        `json:"cell"`
	Size       *OracleExadataMeasurement     `json:"size"`
	FreeSpace  *OracleExadataMeasurement     `json:"freeSpace"`
	Status     string                        `json:"status"`
	ErrorCount int                           `json:"errorCount"`
	GridDisks  []model.OracleExadataGridDisk `json:"gridDisks,omitempty"`
	Databases  []model.OracleExadataDatabase `json:"databases"`
	FreeSizePercentage float64 `json:"freeSizePercentage"`
}

func ToOracleExadataStorageCell(m *model.OracleExadataStorageCell) (res *OracleExadataStorageCell, err error) {
	if m != nil {
		res = &OracleExadataStorageCell{
			Type:       m.Type,
			Hostname:   m.Hostname,
			CellDisk:   m.CellDisk,
			Cell:       m.Cell,
			Status:     m.Status,
			ErrorCount: m.ErrorCount,
			GridDisks:  m.GridDisks,
			Databases:  m.Databases,
		}

		res.Size, err = ToOracleExadataMeasurement(m.Size)
		if err != nil {
			return nil, err
		}

		res.FreeSpace, err = ToOracleExadataMeasurement(m.FreeSpace)
		if err != nil {
			return nil, err
		}

		res.FreeSizePercentage, err = res.GetFreeSpacePercentage()
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, errors.New("cannot create OracleExadataStorageCell dto")
}

func ToOracleExadataStorageCells(m []model.OracleExadataStorageCell) ([]OracleExadataStorageCell, error) {
	res := make([]OracleExadataStorageCell, 0, len(m))

	for _, v := range m {
		dtovalue, err := ToOracleExadataStorageCell(&v)
		if err != nil {
			return nil, err
		}

		res = append(res, *dtovalue)
	}

	return res, nil
}

func (s *OracleExadataStorageCell) GetFreeSpacePercentage() (float64, error) {
	realSize, err := s.Size.ToTb()
	if err != nil {
		return 0, err
	}

	realFreeSpace, err := s.FreeSpace.ToTb()
	if err != nil {
		return 0, err
	}

	res := utils.TruncateFloat64((realFreeSpace.Quantity * 100) / realSize.Quantity)

	return res, nil
}
