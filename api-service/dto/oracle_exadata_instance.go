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
	"errors"
	"time"

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
	TotalMemory int                      `json:"totalMemory"`
	UsedMemory  int                      `json:"usedMemory"`
	FreeMemory  int                      `json:"freeMemory"`
	TotalCPU    int                      `json:"totalCPU"`
	UsedCPU     int                      `json:"usedCPU"`
	FreeCPU     int                      `json:"freeCPU"`
	RDMA        *model.OracleExadataRdma `json:"rdma"`
}

func ToOracleExadataInstance(inst *model.OracleExadataInstance) (*OracleExadataInstance, error) {
	if inst == nil {
		return nil, errors.New("cannot convert nil model to OracleExadataInstance dto")
	}

	res := OracleExadataInstance{
		Hostname:    inst.Hostname,
		RackID:      inst.RackID,
		Location:    inst.Location,
		Environment: inst.Environment,
		CreatedAt:   inst.CreatedAt,
		UpdateAt:    inst.UpdatedAt,
		RDMA:        inst.RDMA,
	}

	for _, cmp := range inst.Components {
		res.TotalMemory += cmp.Memory
		res.TotalCPU += cmp.TotalCPU

		if cmp.HostType == model.BARE_METAL {
			res.UsedCPU += cmp.CPUEnabled
		} else if cmp.HostType == model.DOM0 || cmp.HostType == model.KVM_HOST {
			for _, vm := range cmp.VMs {
				if vm.Type == model.VM_KVM || vm.Type == model.VM_XEN {
					res.UsedMemory += (vm.RamCurrent / 1000) + (vm.RamOnline / 1000)
					res.UsedCPU += vm.CPUCurrent + vm.CPUOnline
				}
			}
		}

		d, err := ToOracleExadataComponent(&cmp)
		if err != nil {
			return nil, err
		}

		res.Components = append(res.Components, *d)
	}

	// Calculate free memory and CPU
	res.FreeMemory = res.TotalMemory - res.UsedMemory
	res.FreeCPU = res.TotalCPU - res.UsedCPU

	return &res, nil
}

func ToOracleExadataInstances(instancesModel []model.OracleExadataInstance) ([]OracleExadataInstance, error) {
	res := make([]OracleExadataInstance, 0, len(instancesModel))

	for _, instance := range instancesModel {
		dto, err := ToOracleExadataInstance(&instance)
		if err != nil {
			return nil, err
		}

		res = append(res, *dto)
	}

	return res, nil
}
