// Copyright (c) 2022 Sorint.lab S.p.A.
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

type PostgreSqlInstanceResponse struct {
	Content  []PostgreSqlInstance `json:"content" bson:"content"`
	Metadata PagingMetadata       `json:"metadata" bson:"metadata"`
}

type PostgreSqlInstance struct {
	Hostname    string `json:"hostname" bson:"hostname"`
	Environment string `json:"environment" bson:"environment"`
	Name        string `json:"name" bson:"name"`
	Charset     string `json:"charset" bson:"charset"`
	Version     string `json:"version" bson:"version"`
}

type SearchPostgreSqlInstancesFilter struct {
	GlobalFilter

	Search     string
	SortBy     string
	SortDesc   bool
	PageNumber int
	PageSize   int
}

func GetPostgreSqlServerInstancesFilter(r *http.Request) (f *SearchPostgreSqlInstancesFilter, err error) {
	f = new(SearchPostgreSqlInstancesFilter)

	gf, err := GetGlobalFilter(r)
	if err != nil {
		return nil, err
	}

	f.GlobalFilter = *gf

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
