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
	"time"

	"github.com/ercole-io/ercole/v2/model"
)

type OracleExadataDatabase struct {
	Type            string
	DbName          string
	Cell            string
	DbID            int
	FlashCacheLimit int
	IormShare       int
	LastIOReq       *time.Time
}

func ToOracleExadataDatabase(m model.OracleExadataDatabase) (*OracleExadataDatabase, error) {
	return &OracleExadataDatabase{
		Type:            m.Type,
		DbName:          m.DbName,
		Cell:            m.Cell,
		DbID:            m.DbID,
		FlashCacheLimit: m.FlashCacheLimit,
		IormShare:       m.IormShare,
		LastIOReq:       m.LastIOReq,
	}, nil
}
