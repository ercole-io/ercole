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

type OracleExadataComponent struct {
	model.OracleExadataComponent
	UsedRAM int `json:"usedRAM"`
	FreeRAM int `json:"freeRAM"`
	UsedCPU int `json:"usedCPU"`
	FreeCPU int `json:"freeCPU"`
}

func ToOracleExadataComponent(componentModel *model.OracleExadataComponent) OracleExadataComponent {
	if componentModel != nil {
		res := OracleExadataComponent{OracleExadataComponent: *componentModel}

		for _, vm := range componentModel.VMs {
			if vm.Type == model.VM_KVM || vm.Type == model.VM_XEN {
				res.UsedRAM += (vm.RamCurrent / 1000) + (vm.RamOnline / 1000)
				res.UsedCPU += vm.CPUCurrent + vm.CPUOnline
			}
		}

		if res.HostType == model.DOM0 || res.HostType == model.KVM_HOST {
			res.FreeRAM = res.Memory - res.UsedRAM
			res.FreeCPU = res.TotalCPU - res.UsedCPU
		}

		return res
	}

	return OracleExadataComponent{}
}
