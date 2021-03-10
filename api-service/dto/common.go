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
	"time"

	"github.com/ercole-io/ercole/v2/utils"
)

type GlobalFilter struct {
	Location    string
	Environment string
	OlderThan   time.Time
}

func GetGlobalFilter(r *http.Request) (f *GlobalFilter, err error) {
	f = new(GlobalFilter)

	f.Location = r.URL.Query().Get("location")
	f.Environment = r.URL.Query().Get("environment")

	if f.OlderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		return nil, err
	}

	return
}
