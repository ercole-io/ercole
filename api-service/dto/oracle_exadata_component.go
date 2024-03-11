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

import "github.com/ercole-io/ercole/v2/api-service/domain"

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
	VMs               []OracleExadataVM          `json:"vms,omitempty"`
	StorageCells      []OracleExadataStorageCell `json:"storageCells,omitempty"`
	ClusterNames      []string                   `json:"clusterNames"`

	UsedRAM           int    `json:"usedRAM"`
	FreeRAM           int    `json:"freeRAM"`
	UsedRAMPercentage string `json:"usedRAMPercentage"`

	UsedCPU           int    `json:"usedCPU"`
	FreeCPU           int    `json:"freeCPU"`
	UsedCPUPercentage string `json:"usedCPUPercentage"`

	TotalSize          int    `json:"totalSize"`
	TotalFreeSpace     int    `json:"totalFreeSpace"`
	UsedSizePercentage string `json:"usedSizePercentage"`
}

func ToOracleExadataComponent(d domain.OracleExadataComponent) (*OracleExadataComponent, error) {
	res := &OracleExadataComponent{
		RackID:             d.RackID,
		HostType:           d.HostType,
		Hostname:           d.Hostname,
		HostID:             d.HostID,
		CPUEnabled:         d.CPUEnabled,
		TotalCPU:           d.TotalCPU,
		ImageVersion:       d.ImageVersion,
		Kernel:             d.Kernel,
		Model:              d.Model,
		FanUsed:            d.FanUsed,
		FanTotal:           d.FanTotal,
		PsuUsed:            d.PsuUsed,
		PsuTotal:           d.PsuTotal,
		MsStatus:           d.MsStatus,
		RsStatus:           d.RsStatus,
		CellServiceStatus:  d.CellServiceStatus,
		SwVersion:          d.SwVersion,
		ClusterNames:       d.ClusterNames,
		UsedCPU:            d.UsedCPU,
		FreeCPU:            d.FreeCPU,
		UsedSizePercentage: d.UsedSizePercentage,
		UsedRAMPercentage:  d.UsedRAMPercentage,
		UsedCPUPercentage:  d.UsedCPUPercentage,
	}

	if d.Memory != nil {
		rmem, err := d.Memory.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.Memory = rmem
	}

	if d.UsedRAM != nil {
		rusedram, err := d.UsedRAM.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.UsedRAM = rusedram
	}

	if d.FreeRAM != nil {
		rfreeram, err := d.FreeRAM.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.FreeRAM = rfreeram
	}

	if d.TotalSize != nil {
		rtotalsize, err := d.TotalSize.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.TotalSize = rtotalsize
	}

	if d.TotalFreeSpace != nil {
		rtotalfreespace, err := d.TotalFreeSpace.RoundedGiB()
		if err != nil {
			return nil, err
		}

		res.TotalFreeSpace = rtotalfreespace
	}

	vms, err := domain.ToUpperLevelLayers[domain.OracleExadataVM, OracleExadataVM](d.VMs, ToOracleExadataVM)
	if err != nil {
		return nil, err
	}

	res.VMs = vms

	storagecells, err := domain.ToUpperLevelLayers[domain.OracleExadataStorageCell, OracleExadataStorageCell](d.StorageCells, ToOracleExadataStorageCell)
	if err != nil {
		return nil, err
	}

	res.StorageCells = storagecells

	return res, nil
}
