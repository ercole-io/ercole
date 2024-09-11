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

package service

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/domain"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as APIService) GetAllExadataInstanceAsXlsx() (*excelize.File, error) {
	exadataInstances, err := as.getAllExadataInstance()
	if err != nil {
		return nil, err
	}

	memorysheet := "Memory-CPU"
	storagesheet := "Storage"
	clustersheet := "Cluster View"

	file, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_exadatas.xlsx")
	if err != nil {
		return nil, err
	}

	axisHelperMem := exutils.NewAxisHelper(1)
	axisHelperStorage := exutils.NewAxisHelper(1)
	axisHelperUnk := exutils.NewAxisHelper(1)

	for _, instance := range exadataInstances {
		for _, component := range instance.Components {
			switch {
			case component.HostType != model.STORAGE_CELL:
				nextAxis := axisHelperMem.NewRow()
				file.SetCellValue(memorysheet, nextAxis(), instance.Hostname)
				file.SetCellValue(memorysheet, nextAxis(), instance.RackID)
				file.SetCellValue(memorysheet, nextAxis(), component.Hostname)
				file.SetCellValue(memorysheet, nextAxis(), component.HostType)
				file.SetCellValue(memorysheet, nextAxis(), component.HostID)
				file.SetCellValue(memorysheet, nextAxis(), component.Memory)
				file.SetCellValue(memorysheet, nextAxis(), component.UsedRAM)
				file.SetCellValue(memorysheet, nextAxis(), component.ReservedMemory)
				file.SetCellValue(memorysheet, nextAxis(), component.UsedRAMPercentage)
				file.SetCellValue(memorysheet, nextAxis(), component.TotalCPU)
				file.SetCellValue(memorysheet, nextAxis(), component.UsedCPU)
				file.SetCellValue(memorysheet, nextAxis(), component.ReservedCPU)
				file.SetCellValue(memorysheet, nextAxis(), component.UsedCPUPercentage)

			case component.HostType == model.STORAGE_CELL:
				nextAxis := axisHelperStorage.NewRow()
				file.SetCellValue(storagesheet, nextAxis(), instance.Hostname)
				file.SetCellValue(storagesheet, nextAxis(), instance.RackID)
				file.SetCellValue(storagesheet, nextAxis(), component.Hostname)
				file.SetCellValue(storagesheet, nextAxis(), component.HostType)
				file.SetCellValue(storagesheet, nextAxis(), component.HostID)
				file.SetCellValue(storagesheet, nextAxis(), component.TotalSize)
				file.SetCellValue(storagesheet, nextAxis(), component.UsedSizePercentage)
				file.SetCellValue(storagesheet, nextAxis(), component.TotalFreeSpace)
			}
		}
	}

	clusterViews, err := as.Database.FindExadataClusterViews()
	if err != nil {
		return nil, err
	}

	initcolumn := 7
	newcolumn := 0

	for _, cv := range clusterViews {
		nextAxisCls := axisHelperUnk.NewRow()
		file.SetCellValue(clustersheet, nextAxisCls(), cv.Hostname)
		file.SetCellValue(clustersheet, nextAxisCls(), cv.RackID)
		file.SetCellValue(clustersheet, nextAxisCls(), cv.HostType)
		file.SetCellValue(clustersheet, nextAxisCls(), cv.TotalRAM)
		file.SetCellValue(clustersheet, nextAxisCls(), cv.TotalCPU)
		file.SetCellValue(clustersheet, nextAxisCls(), cv.Clustername)

		for i, vmname := range cv.VmNames {
			if i > 1 {
				newcolumn = i
			}

			file.SetCellValue(clustersheet, nextAxisCls(), vmname.(string))
		}
	}

	for i := 1; i < newcolumn; i++ {
		col := initcolumn + i
		file.SetCellValue(clustersheet, fmt.Sprintf("%s1", string(rune('A'+col))), fmt.Sprintf("Phys Nodes %d", i+2))
	}

	return file, err
}

func (as APIService) getAllExadataInstance() ([]dto.OracleExadataInstance, error) {
	instances, err := as.Database.FindAllExadataInstances(false)
	if err != nil {
		return nil, err
	}

	doms, err := domain.ToUpperLevelLayers[model.OracleExadataInstance, domain.OracleExadataInstance](instances, domain.ToOracleExadataInstance)
	if err != nil {
		return nil, err
	}

	res, err := domain.ToUpperLevelLayers[domain.OracleExadataInstance, dto.OracleExadataInstance](doms, dto.ToOracleExadataInstance)
	if err != nil {
		return nil, err
	}

	return res, err
}
