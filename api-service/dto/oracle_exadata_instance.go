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
	"time"

	"github.com/ercole-io/ercole/v2/api-service/domain"
	"github.com/ercole-io/ercole/v2/model"
)

type OracleExadataInstance struct {
	Hostname    string                   `json:"hostname"`
	Environment string                   `json:"environment"`
	Location    string                   `json:"location"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdateAt    time.Time                `json:"updateAt"`
	RackID      string                   `json:"rackID"`
	Components  []OracleExadataComponent `json:"components"`
	RDMA        *model.OracleExadataRdma `json:"rdma"`

	TotalMemory           string `json:"totalMemory"`
	UsedMemory            string `json:"usedMemory"`
	FreeMemory            string `json:"freeMemory"`
	UsedMemoryMPercentage string `json:"usedMemoryPercentage"`

	TotalCPU          int    `json:"totalCPU"`
	UsedCPU           int    `json:"usedCPU"`
	FreeCPU           int    `json:"freeCPU"`
	UsedCPUPercentage string `json:"usedCPUPercentage"`

	TotalSize          string `json:"totalSize"`
	UsedSize           string `json:"usedSize"`
	FreeSpace          string `json:"freeSpace"`
	UsedSizePercentage string `json:"usedSizePercentage"`
}

func ToOracleExadataInstance(d domain.OracleExadataInstance) (*OracleExadataInstance, error) {
	res := &OracleExadataInstance{
		Hostname:              d.Hostname,
		Environment:           d.Environment,
		Location:              d.Location,
		CreatedAt:             d.CreatedAt,
		UpdateAt:              d.UpdatedAt,
		RackID:                d.RackID,
		RDMA:                  d.RDMA,
		TotalCPU:              d.TotalCPU,
		UsedCPU:               d.UsedCPU,
		FreeCPU:               d.FreeCPU,
		UsedMemoryMPercentage: d.UsedMemoryPercentage,
		UsedCPUPercentage:     d.UsedCPUPercentage,
		UsedSizePercentage:    d.UsedSizePercentage,
	}

	if d.TotalMemory != nil {
		htotalmemory, err := d.TotalMemory.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.TotalMemory = htotalmemory
	}

	if d.UsedMemory != nil {
		husedmemory, err := d.UsedMemory.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.UsedMemory = husedmemory
	}

	if d.FreeMemory != nil {
		hfreememory, err := d.FreeMemory.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.FreeMemory = hfreememory
	}

	if d.TotalSize != nil {
		htotalsize, err := d.TotalSize.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.TotalSize = htotalsize
	}

	if d.UsedSize != nil {
		husedsize, err := d.UsedSize.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.UsedSize = husedsize
	}

	if d.FreeSpace != nil {
		hfreespace, err := d.FreeSpace.Human("GIB")
		if err != nil {
			return nil, err
		}

		res.FreeSpace = hfreespace
	}

	componens, err := domain.ToUpperLevelLayers[domain.OracleExadataComponent, OracleExadataComponent](d.Components, ToOracleExadataComponent)
	if err != nil {
		return nil, err
	}

	res.Components = componens

	return res, nil
}
