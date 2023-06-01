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

import "github.com/ercole-io/ercole/v2/model"

type OracleExadataInstance struct {
	Hostname    string                   `json:"hostname"`
	Environment string                   `json:"environment"`
	Location    string                   `json:"location"`
	RackID      string                   `json:"rackID"`
	Components  []OracleExadataComponent `json:"components"`
	TotalMemory int                      `json:"totalMemory"`
	UsedMemory  int                      `json:"usedMemory"`
	FreeMemory  int                      `json:"freeMemory"`
	TotalCPU    int                      `json:"totalCPU"`
	UsedCPU     int                      `json:"usedCPU"`
	FreeCPU     int                      `json:"freeCPU"`
}

func ToOracleExadataInstance(inst *model.OracleExadataInstance) OracleExadataInstance {
	if inst != nil {
		res := OracleExadataInstance{
			Hostname:    inst.Hostname,
			RackID:      inst.RackID,
			Location:    inst.Location,
			Environment: inst.Environment,
		}

		for _, cmp := range inst.Components {
			if cmp.HostType == model.DOM0 || cmp.HostType == model.KVM_HOST {
				res.TotalMemory += cmp.Memory
				res.TotalCPU += cmp.TotalCPU

				for _, vm := range cmp.VMs {
					if vm.Type == model.VM_KVM || vm.Type == model.VM_XEN {
						res.UsedMemory += vm.RamCurrent + vm.RamOnline
						res.UsedCPU += vm.CPUCurrent + vm.CPUOnline
					}
				}
			}

			res.Components = append(res.Components, ToOracleExadataComponent(&cmp))
		}

		res.FreeMemory = res.TotalMemory - res.UsedMemory
		res.FreeCPU = res.TotalCPU - res.UsedCPU

		return res
	}

	return OracleExadataInstance{}
}

func ToOracleExadataInstances(instancesModel []model.OracleExadataInstance) []OracleExadataInstance {
	res := make([]OracleExadataInstance, 0, len(instancesModel))

	for _, instance := range instancesModel {
		res = append(res, ToOracleExadataInstance(&instance))
	}

	return res
}
