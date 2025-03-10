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

	"github.com/ercole-io/ercole/v2/model"
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

func NewSearchHostsFilters() SearchHostsFilters {
	return SearchHostsFilters{
		Search:                        []string{},
		SortBy:                        "",
		SortDesc:                      false,
		Location:                      "",
		Environment:                   "",
		OlderThan:                     utils.MAX_TIME,
		PageNumber:                    -1,
		PageSize:                      -1,
		Hostname:                      "",
		Database:                      "",
		Technology:                    "",
		HardwareAbstractionTechnology: "",
		Cluster:                       new(string),
		VirtualizationNode:            "",
		OperatingSystem:               "",
		Kernel:                        "",

		LTEMemoryTotal:    -1,
		GTEMemoryTotal:    -1,
		LTESwapTotal:      -1,
		GTESwapTotal:      -1,
		IsMemberOfCluster: nil,
		CPUModel:          "",
		LTECPUCores:       -1,
		GTECPUCores:       -1,
		LTECPUThreads:     -1,
		GTECPUThreads:     -1,
	}
}

func GetSearchHostFilters(r *http.Request) (*SearchHostsFilters, error) {
	f := SearchHostsFilters{}

	var err error

	f.Search = strings.Split(r.URL.Query().Get("search"), " ")

	f.SortBy = r.URL.Query().Get("sort-by")

	if f.SortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		return nil, err
	}

	f.Location = r.URL.Query().Get("location")
	f.Environment = r.URL.Query().Get("environment")

	if f.OlderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		return nil, err
	}

	if f.PageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		return nil, err
	}

	if f.PageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		return nil, err
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

	if f.LTEMemoryTotal, err = utils.Str2float64(r.URL.Query().Get("memory-total-lte"), -1); err != nil {
		return nil, err
	}

	if f.GTEMemoryTotal, err = utils.Str2float64(r.URL.Query().Get("memory-total-gte"), -1); err != nil {
		return nil, err
	}

	if f.LTESwapTotal, err = utils.Str2float64(r.URL.Query().Get("swap-total-lte"), -1); err != nil {
		return nil, err
	}

	if f.GTESwapTotal, err = utils.Str2float64(r.URL.Query().Get("swap-total-gte"), -1); err != nil {
		return nil, err
	}

	if r.URL.Query().Get("is-member-of-cluster") == "" {
		f.IsMemberOfCluster = nil
	} else {
		f.IsMemberOfCluster = new(bool)
		if *f.IsMemberOfCluster, err = utils.Str2bool(r.URL.Query().Get("is-member-of-cluster"), false); err != nil {
			return nil, err
		}
	}

	f.CPUModel = r.URL.Query().Get("cpu-model")

	if f.LTECPUCores, err = utils.Str2int(r.URL.Query().Get("cpu-cores-lte"), -1); err != nil {
		return nil, err
	}

	if f.GTECPUCores, err = utils.Str2int(r.URL.Query().Get("cpu-cores-gte"), -1); err != nil {
		return nil, err
	}

	if f.LTECPUThreads, err = utils.Str2int(r.URL.Query().Get("cpu-threads-lte"), -1); err != nil {
		return nil, err
	}

	if f.GTECPUThreads, err = utils.Str2int(r.URL.Query().Get("cpu-threads-gte"), -1); err != nil {
		return nil, err
	}

	return &f, nil
}

type HostDataSummary struct {
	ID                      string                        `json:"id" bson:"_id"`
	CreatedAt               time.Time                     `json:"createdAt" bson:"createdAt"`
	Hostname                string                        `json:"hostname" bson:"hostname"`
	Location                string                        `json:"location" bson:"location"`
	Environment             string                        `json:"environment" bson:"environment"`
	AgentVersion            string                        `json:"agentVersion" bson:"agentVersion"`
	Info                    model.Host                    `json:"info" bson:"info"`
	ClusterMembershipStatus model.ClusterMembershipStatus `json:"clusterMembershipStatus" bson:"clusterMembershipStatus"`
	VirtualizationNode      string                        `json:"virtualizationNode" bson:"virtualizationNode"`
	Cluster                 string                        `json:"cluster" bson:"cluster"`
	Databases               map[string][]string           `json:"databases" bson:"databases"` // map[Technology] []database names
	MissingDatabases        []model.MissingDatabase       `json:"missingDatabases,omitempty"`
	Technology              string                        `json:"technology" bson:"technology"`
}

type VirtualHostWithoutCluster struct {
	Hostname                      string `json:"hostname" bson:"hostname"`
	HardwareAbstractionTechnology string `json:"hardwareAbstractionTechnology" bson:"hardwareAbstractionTechnology"`
}
