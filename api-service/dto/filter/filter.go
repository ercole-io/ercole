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

package filter

import (
	"net/http"
	"strconv"
)

const limitDefault = 25

const page = 1

type Filter struct {
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
	Search *string `json:"search"`
}

func (f *Filter) GetSkip() int {
	return (f.Page - 1) * f.Limit
}

func New() Filter {
	return Filter{
		Limit: limitDefault,
		Page:  page,
	}
}

func Get(r *http.Request) (*Filter, error) {
	f := New()

	search := r.URL.Query().Get("search")
	if search != "" {
		f.Search = &search
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}

		f.Page = p
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}

		f.Limit = l
	}

	return &f, nil
}
