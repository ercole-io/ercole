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
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
)

type OracleExadataComponent struct {
	RackID            string                     `json:"rackID"`
	HostType          string                     `json:"hostType"`
	Hostname          string                     `json:"hostname"`
	HostID            string                     `json:"hostID"`
	CPUEnabled        int                        `json:"cpuEnabled"`
	TotalCPU          int                        `json:"totalCPU"`
	Memory            int                        `json:"memory"`
	ImageVersion      string                     `json:"imageVersion"`
	Kernel            string                     `json:"kernel"`
	Model             string                     `json:"model"`
	FanUsed           int                        `json:"fanUsed"`
	FanTotal          int                        `json:"fanTotal"`
	PsuUsed           int                        `json:"psuUsed"`
	PsuTotal          int                        `json:"psuTotal"`
	MsStatus          string                     `json:"msStatus"`
	RsStatus          string                     `json:"rsStatus"`
	CellServiceStatus string                     `json:"cellServiceStatus"`
	SwVersion         string                     `json:"swVersion"`
	VMs               []model.OracleExadataVM    `json:"vms,omitempty"`
	StorageCells      []OracleExadataStorageCell `json:"storageCells,omitempty"`

	UsedRAM int `json:"usedRAM"`
	FreeRAM int `json:"freeRAM"`
	UsedCPU int `json:"usedCPU"`
	FreeCPU int `json:"freeCPU"`

	FreeSizePercentage string `json:"freeSizePercentage"`
}

func ToOracleExadataComponent(componentModel *model.OracleExadataComponent) (*OracleExadataComponent, error) {
	if componentModel != nil {
		storagedtos, err := ToOracleExadataStorageCells(componentModel.StorageCells)
		if err != nil {
			return nil, err
		}

		res := OracleExadataComponent{
			RackID:            componentModel.RackID,
			HostType:          componentModel.HostType,
			Hostname:          componentModel.Hostname,
			HostID:            componentModel.HostID,
			CPUEnabled:        componentModel.CPUEnabled,
			TotalCPU:          componentModel.TotalCPU,
			Memory:            componentModel.Memory,
			ImageVersion:      componentModel.ImageVersion,
			Kernel:            componentModel.Kernel,
			Model:             componentModel.Model,
			FanUsed:           componentModel.FanUsed,
			FanTotal:          componentModel.FanTotal,
			PsuUsed:           componentModel.PsuUsed,
			PsuTotal:          componentModel.PsuTotal,
			MsStatus:          componentModel.MsStatus,
			RsStatus:          componentModel.RsStatus,
			CellServiceStatus: componentModel.CellServiceStatus,
			SwVersion:         componentModel.SwVersion,
			VMs:               componentModel.VMs,
			StorageCells:      storagedtos,
		}

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

		perc, err := res.GetFreeSpacePercentage()
		if err != nil {
			return nil, err
		}

		res.FreeSizePercentage = perc

		return &res, nil
	}

	return nil, errors.New("cannot convert model OracleExadataComponent dto")
}

func (c *OracleExadataComponent) GetFreeSpacePercentage() (string, error) {
	totsize := 0.0
	totFreeSpace := 0.0

	for _, storageCell := range c.StorageCells {
		sizeTb, err := storageCell.Size.ToTb()
		if err != nil {
			return "", err
		}

		totsize += sizeTb.Quantity

		freeSpaceTb, err := storageCell.FreeSpace.ToTb()
		if err != nil {
			return "", err
		}

		totFreeSpace += freeSpaceTb.Quantity
	}

	return fmt.Sprintf("%v%%", (totFreeSpace*100)/totsize), nil
}
