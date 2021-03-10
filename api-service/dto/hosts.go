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
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/utils"
)

type SearchHostsFilters struct {
	Search      []string
	SortBy      string
	SortDesc    bool
	Location    string
	Environment string
	OlderThan   time.Time
	PageNumber  int
	PageSize    int

	Hostname                      string
	Database                      string
	Technology                    string
	HardwareAbstractionTechnology string
	Cluster                       *string
	VirtualizationNode            string
	OperatingSystem               string
	Kernel                        string
	LTEMemoryTotal                float64
	GTEMemoryTotal                float64
	LTESwapTotal                  float64
	GTESwapTotal                  float64
	IsMemberOfCluster             *bool
	CPUModel                      string
	LTECPUCores                   int
	GTECPUCores                   int
	LTECPUThreads                 int
	GTECPUThreads                 int
}

func GetSearchHostFilters(r *http.Request) (*SearchHostsFilters, utils.AdvancedErrorInterface) {
	f := SearchHostsFilters{}
	var aerr utils.AdvancedErrorInterface

	f.Search = strings.Split(r.URL.Query().Get("search"), " ")

	f.SortBy = r.URL.Query().Get("sort-by")

	if f.SortDesc, aerr = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); aerr != nil {
		return nil, aerr
	}

	f.Location = r.URL.Query().Get("location")
	f.Environment = r.URL.Query().Get("environment")

	if f.OlderThan, aerr = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); aerr != nil {
		return nil, aerr
	}

	if f.PageNumber, aerr = utils.Str2int(r.URL.Query().Get("page"), -1); aerr != nil {
		return nil, aerr
	}

	if f.PageSize, aerr = utils.Str2int(r.URL.Query().Get("size"), -1); aerr != nil {
		return nil, aerr
	}

	f.Hostname = r.URL.Query().Get("hostname")
	f.Database = r.URL.Query().Get("database")
	f.Technology = r.URL.Query().Get("technology")
	f.HardwareAbstractionTechnology = r.URL.Query().Get("hardware-abstraction-technology")
	if r.URL.Query().Get("cluster") == "NULL" {
		f.Cluster = nil
	} else {
		f.Cluster = new(string)
		*f.Cluster = r.URL.Query().Get("cluster")
	}
	f.VirtualizationNode = r.URL.Query().Get("virtualization-node")
	f.OperatingSystem = r.URL.Query().Get("operating-system")
	f.Kernel = r.URL.Query().Get("kernel")
	if f.LTEMemoryTotal, aerr = utils.Str2float64(r.URL.Query().Get("memory-total-lte"), -1); aerr != nil {
		return nil, aerr
	}
	if f.GTEMemoryTotal, aerr = utils.Str2float64(r.URL.Query().Get("memory-total-gte"), -1); aerr != nil {
		return nil, aerr
	}
	if f.LTESwapTotal, aerr = utils.Str2float64(r.URL.Query().Get("swap-total-lte"), -1); aerr != nil {
		return nil, aerr
	}
	if f.GTESwapTotal, aerr = utils.Str2float64(r.URL.Query().Get("swap-total-gte"), -1); aerr != nil {
		return nil, aerr
	}
	if r.URL.Query().Get("is-member-of-cluster") == "" {
		f.IsMemberOfCluster = nil
	} else {
		f.IsMemberOfCluster = new(bool)
		if *f.IsMemberOfCluster, aerr = utils.Str2bool(r.URL.Query().Get("is-member-of-cluster"), false); aerr != nil {
			return nil, aerr
		}
	}
	f.CPUModel = r.URL.Query().Get("cpu-model")
	if f.LTECPUCores, aerr = utils.Str2int(r.URL.Query().Get("cpu-cores-lte"), -1); aerr != nil {
		return nil, aerr
	}
	if f.GTECPUCores, aerr = utils.Str2int(r.URL.Query().Get("cpu-cores-gte"), -1); aerr != nil {
		return nil, aerr
	}
	if f.LTECPUThreads, aerr = utils.Str2int(r.URL.Query().Get("cpu-threads-lte"), -1); aerr != nil {
		return nil, aerr
	}
	if f.GTECPUThreads, aerr = utils.Str2int(r.URL.Query().Get("cpu-threads-gte"), -1); aerr != nil {
		return nil, aerr
	}

	return &f, nil
}
