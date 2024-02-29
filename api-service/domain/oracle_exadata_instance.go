// Copyright (c) 2024 Sorint.lab S.p.A.
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

package domain

import (
	"time"

	"github.com/ercole-io/ercole/v2/model"
)

type OracleExadataInstance struct {
	Hostname    string
	Environment string
	Location    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RackID      string
	Components  []OracleExadataComponent
	RDMA        *model.OracleExadataRdma

	TotalMemory          *OracleExadataMeasurement
	UsedMemory           *OracleExadataMeasurement
	FreeMemory           *OracleExadataMeasurement
	UsedMemoryPercentage string

	TotalCPU          int
	UsedCPU           int
	FreeCPU           int
	UsedCPUPercentage string

	TotalSize          *OracleExadataMeasurement
	UsedSize           *OracleExadataMeasurement
	FreeSpace          *OracleExadataMeasurement
	UsedSizePercentage string
}

func ToOracleExadataInstance(m model.OracleExadataInstance) (*OracleExadataInstance, error) {
	res := &OracleExadataInstance{
		Hostname:    m.Hostname,
		Environment: m.Environment,
		Location:    m.Location,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		RackID:      m.RackID,
		RDMA:        m.RDMA,
	}

	res.TotalMemory = NewOracleExadataMeasurement()
	res.UsedMemory = NewOracleExadataMeasurement()
	res.FreeMemory = NewOracleExadataMeasurement()
	res.TotalSize = NewOracleExadataMeasurement()
	res.UsedSize = NewOracleExadataMeasurement()
	res.FreeSpace = NewOracleExadataMeasurement()

	components, err := ToUpperLevelLayers[model.OracleExadataComponent, OracleExadataComponent](m.Components, ToOracleExadataComponent)
	if err != nil {
		return nil, err
	}

	res.Components = components

	for _, c := range res.Components {
		if c.HostType != model.STORAGE_CELL {
			res.TotalMemory.Add(c.Memory.Quantity, c.Memory.Symbol)

			res.TotalCPU += c.TotalCPU
		}

		if c.HostType == model.BARE_METAL {
			res.UsedCPU += c.CPUEnabled
		}

		if c.HostType == model.DOM0 || c.HostType == model.KVM_HOST {
			for _, vm := range c.VMs {
				if vm.Type == model.VM_KVM || vm.Type == model.VM_XEN {
					res.UsedMemory.Add(vm.RamCurrent.Quantity, vm.RamCurrent.Symbol)
					res.UsedMemory.Add(vm.RamOnline.Quantity, vm.RamOnline.Symbol)

					res.UsedCPU += vm.CPUCurrent + vm.CPUOnline
				}
			}
		}

		res.TotalSize.Add(c.TotalSize.Quantity, c.TotalSize.Symbol)
		res.FreeSpace.Add(c.TotalFreeSpace.Quantity, c.TotalFreeSpace.Symbol)
	}

	res.FreeMemory.Quantity = res.TotalMemory.Quantity
	res.FreeMemory.Symbol = res.TotalMemory.Symbol

	if res.UsedMemory != nil {
		res.FreeMemory.Sub(*res.UsedMemory)
	}

	// ram usage percentage
	res.UsedMemoryPercentage = GetPercentage(*res.UsedMemory, *res.TotalMemory)

	// vcpu usage percentage
	res.FreeCPU = res.TotalCPU - res.UsedCPU
	res.UsedCPUPercentage = getUsedCpuPercentage(res.UsedCPU, res.TotalCPU)

	// storage usage percentage
	usedsize := OracleExadataMeasurement{
		Symbol:   res.TotalSize.Symbol,
		Quantity: res.TotalSize.Quantity,
	}

	if res.FreeSpace != nil {
		usedsize.Sub(*res.FreeSpace)
	}

	res.UsedSize = &usedsize

	res.UsedSizePercentage = GetPercentage(*res.UsedSize, *res.TotalSize)

	return res, nil
}
