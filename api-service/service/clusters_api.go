// Copyright (c) 2022 Sorint.lab S.p.A.
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
// along with this program.  If not, see <http://www.gn+èu.org/licenses/>.

// Package service is a package that provides methods for querying data
package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/utils/exutils"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ercole-io/ercole/v2/api-service/dto"
)

// SearchClusters search clusters
func (as *APIService) SearchClusters(mode string, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, error) {
	return as.Database.SearchClusters(mode, strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// GetCluster return the cluster specified in the clusterName param
func (as *APIService) GetCluster(clusterName string, olderThan time.Time) (*dto.Cluster, error) {
	cluster, err := as.Database.GetCluster(clusterName, olderThan)
	if err != nil {
		return nil, err
	}

	var errEH error

	for i, vm := range cluster.VMs {
		exist, err := as.Database.ExistHostdata(vm.Hostname)
		if err != nil {
			as.Log.Error(err)
			return nil, err
		}

		if !exist {
			res := strings.Split(vm.Hostname, ".")
			if len(res) > 0 {
				exist, errEH = as.Database.ExistHostdata(res[0])
				if errEH != nil {
					as.Log.Error(errEH)
					return nil, errEH
				}

				if !exist {
					exist, errEH = as.Database.ExistHostdata(strings.ToUpper(res[0]))
					if errEH != nil {
						as.Log.Error(errEH)
						return nil, errEH
					}

					if exist {
						cluster.VMs[i].Hostname = strings.ToUpper(res[0])
					}
				} else {
					cluster.VMs[i].Hostname = res[0]
				}
			} else {
				return nil, errors.New("Hostname not correct")
			}
		}

		cluster.VMs[i].IsErcoleInstalled = exist
	}

	return cluster, nil
}

// GetClusterXLSX return  cluster vms as xlxs file
func (as *APIService) GetClusterXLSX(clusterName string, olderThan time.Time) (*excelize.File, error) {
	cluster, err := as.Database.GetCluster(clusterName, olderThan)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_cluster.xlsx")
	if err != nil {
		return nil, err
	}

	for i, val := range cluster.VMs {
		i += 2
		xlsx.SetCellValue("VMs", fmt.Sprintf("A%d", i), val.Name)
		xlsx.SetCellValue("VMs", fmt.Sprintf("B%d", i), val.Hostname)
		xlsx.SetCellValue("VMs", fmt.Sprintf("C%d", i), val.VirtualizationNode)
		xlsx.SetCellValue("VMs", fmt.Sprintf("D%d", i), val.PhysicalServerModelName)
		xlsx.SetCellValue("VMs", fmt.Sprintf("E%d", i), strconv.FormatBool(val.CappedCPU))
	}

	return xlsx, nil
}

// SearchClustersAsXLSX return  clusters vms as xlxs file
func (as *APIService) SearchClustersAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	clusters, err := as.Database.SearchClusters("full", []string{}, "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	sheet := "Hypervisor"
	headers := []string{
		"Name",
		"Type",
		"Cores",
		"Socket",
		"Physical Host",
		"Total VM",
		"Total VM Ercole",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range clusters {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val["name"])
		sheets.SetCellValue(sheet, nextAxis(), val["type"])
		sheets.SetCellValue(sheet, nextAxis(), val["cpu"])
		sheets.SetCellValue(sheet, nextAxis(), val["sockets"])
		sheets.SetCellValue(sheet, nextAxis(), val["virtualizationNodes"])
		sheets.SetCellValue(sheet, nextAxis(), val["vmsCount"])
		sheets.SetCellValue(sheet, nextAxis(), val["vmsErcoleAgentCount"])
	}

	return sheets, nil
}
