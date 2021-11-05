// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"time"

	"github.com/ercole-io/ercole/v2/utils"
)

type SearchHostsAsLMS struct {
	SearchHostsFilters
	From time.Time
	To   time.Time
}

func GetSearchHostsAsLMSFilters(r *http.Request) (*SearchHostsAsLMS, error) {
	flms := SearchHostsAsLMS{}
	var err error

	hf, errhf := GetSearchHostFilters(r)
	if errhf != nil {
		return nil, errhf
	}

	flms.SearchHostsFilters = *hf

	if flms.From, err = utils.Str2time(r.URL.Query().Get("from"), utils.MIN_TIME); err != nil {
		return nil, err
	}

	if flms.To, err = utils.Str2time(r.URL.Query().Get("to"), utils.MAX_TIME); err != nil {
		return nil, err
	}

	return &flms, nil
}
