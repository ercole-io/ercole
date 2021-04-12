// Copyright (c) 2020 Sorint.lab S.p.A.
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
	"net/http"

	"github.com/ercole-io/ercole/v2/utils"
)

type SearchOracleDatabasesFilter struct {
	GlobalFilter

	Full       bool
	Search     string
	SortBy     string
	SortDesc   bool
	PageNumber int
	PageSize   int
}

func GetSearchOracleDatabasesFilter(r *http.Request) (f *SearchOracleDatabasesFilter, err error) {
	f = new(SearchOracleDatabasesFilter)

	gf, err := GetGlobalFilter(r)
	if err != nil {
		return nil, err
	}
	f.GlobalFilter = *gf

	if f.Full, err = utils.Str2bool(r.URL.Query().Get("full"), false); err != nil {
		return nil, err
	}

	f.Search = r.URL.Query().Get("search")
	f.SortBy = r.URL.Query().Get("sort-by")
	if f.SortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		return nil, err
	}

	if f.PageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		return nil, err
	}
	if f.PageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		return nil, err
	}

	return
}

type OracleDatabasesStatistics struct {
	TotalMemorySize   float64 `json:"total-memory-size"`   // in bytes
	TotalSegmentsSize float64 `json:"total-segments-size"` // in bytes
	TotalDatafileSize float64 `json:"total-datafile-size"` // in bytes
	TotalWork         float64 `json:"total-work"`
}
