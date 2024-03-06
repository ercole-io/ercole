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
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
)

type OracleExadataComponent struct {
	RackID            string
	HostType          string
	Hostname          string
	HostID            string
	CPUEnabled        int
	TotalCPU          int
	Memory            *OracleExadataMeasurement
	ImageVersion      string
	Kernel            string
	Model             string
	FanUsed           int
	FanTotal          int
	PsuUsed           int
	PsuTotal          int
	MsStatus          string
	RsStatus          string
	CellServiceStatus string
	SwVersion         string
	VMs               []OracleExadataVM
	StorageCells      []OracleExadataStorageCell
	ClusterNames      []string

	UsedRAM           *OracleExadataMeasurement
	FreeRAM           *OracleExadataMeasurement
	UsedRAMPercentage string

	UsedCPU           int
	FreeCPU           int
	UsedCPUPercentage string

	TotalSize          *OracleExadataMeasurement
	TotalFreeSpace     *OracleExadataMeasurement
	UsedSizePercentage string
}

func ToOracleExadataComponent(m model.OracleExadataComponent) (*OracleExadataComponent, error) {
	res := &OracleExadataComponent{
		RackID:            m.RackID,
		HostType:          m.HostType,
		Hostname:          m.Hostname,
		HostID:            m.HostID,
		CPUEnabled:        m.CPUEnabled,
		TotalCPU:          m.TotalCPU,
		ImageVersion:      m.ImageVersion,
		Kernel:            m.Kernel,
		Model:             m.Model,
		FanUsed:           m.FanUsed,
		FanTotal:          m.FanTotal,
		PsuUsed:           m.PsuUsed,
		PsuTotal:          m.PsuTotal,
		MsStatus:          m.MsStatus,
		RsStatus:          m.RsStatus,
		CellServiceStatus: m.CellServiceStatus,
		SwVersion:         m.SwVersion,
		ClusterNames:      m.ClusterNames,
	}

	res.UsedRAM = NewOracleExadataMeasurement()
	res.FreeRAM = NewOracleExadataMeasurement()
	res.TotalSize = NewOracleExadataMeasurement()
	res.TotalFreeSpace = NewOracleExadataMeasurement()

	memory, err := IntToOracleExadataMeasurement(m.Memory, "GB")
	if err != nil {
		return nil, err
	}

	res.Memory = memory

	vms, err := ToUpperLevelLayers[model.OracleExadataVM, OracleExadataVM](m.VMs, ToOracleExadataVm)
	if err != nil {
		return nil, err
	}

	res.VMs = vms

	storagecells, err := ToUpperLevelLayers[model.OracleExadataStorageCell, OracleExadataStorageCell](m.StorageCells, ToOracleExadataStorageCell)
	if err != nil {
		return nil, err
	}

	res.StorageCells = storagecells

	for _, vm := range res.VMs {
		if vm.Type == model.VM_KVM || vm.Type == model.VM_XEN {
			res.UsedRAM.Add(vm.RamCurrent.Quantity, vm.RamCurrent.Symbol)
			res.UsedRAM.Add(vm.RamOnline.Quantity, vm.RamOnline.Symbol)

			res.UsedCPU += vm.CPUCurrent + vm.CPUOnline
		}
	}

	if res.HostType == model.DOM0 || res.HostType == model.KVM_HOST {
		freeram := OracleExadataMeasurement{
			Symbol:   res.Memory.Symbol,
			Quantity: res.Memory.Quantity,
		}

		if res.UsedRAM != nil {
			freeram.Sub(*res.UsedRAM)
		}

		res.FreeRAM = &freeram

		res.FreeCPU = res.TotalCPU - res.UsedCPU
	}

	if res.HostType == model.BARE_METAL {
		res.UsedCPU += m.CPUEnabled
	}

	for _, sc := range res.StorageCells {
		res.TotalSize.Add(sc.Size.Quantity, sc.Size.Symbol)
		res.TotalFreeSpace.Add(sc.FreeSpace.Quantity, sc.FreeSpace.Symbol)
	}

	// storage usage percentage
	usedsize := OracleExadataMeasurement{
		Symbol:   res.TotalSize.Symbol,
		Quantity: res.TotalSize.Quantity,
	}
	usedsize.Sub(*res.TotalFreeSpace)
	res.UsedSizePercentage = GetPercentage(usedsize, *res.TotalSize)

	// vcpu usage percentage
	res.UsedCPUPercentage = getUsedCpuPercentage(res.UsedCPU, res.TotalCPU)

	// ram usage percentage
	res.UsedRAMPercentage = GetPercentage(*res.UsedRAM, *res.Memory)

	return res, nil
}

func getUsedCpuPercentage(used, total int) string {
	if total != 0 {
		perc := (used * 100) / total

		return fmt.Sprintf("%d%%", perc)
	}

	return fmt.Sprintf("0%%")
}
