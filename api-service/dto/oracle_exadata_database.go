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
	"time"

	"github.com/ercole-io/ercole/v2/api-service/domain"
)

type OracleExadataDatabase struct {
	Type            string     `json:"type"`
	DbName          string     `json:"dbName"`
	Cell            string     `json:"cell"`
	DbID            int        `json:"dbID"`
	FlashCacheLimit int        `json:"flashCacheLimit"`
	IormShare       int        `json:"iormShare"`
	LastIOReq       *time.Time `json:"lastIOReq"`
}

func ToOracleExadataDatabase(d domain.OracleExadataDatabase) (*OracleExadataDatabase, error) {
	return &OracleExadataDatabase{
		Type:            d.Type,
		DbName:          d.DbName,
		Cell:            d.Cell,
		DbID:            d.DbID,
		FlashCacheLimit: d.FlashCacheLimit,
		IormShare:       d.IormShare,
		LastIOReq:       d.LastIOReq,
	}, nil
}
