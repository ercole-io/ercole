// Copyright (c) 2020 Sorint.lab S.p.A.
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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

// SearchClusters search clusters
func (as *APIService) SearchClusters(full bool, search string, sortBy string, sortDesc bool, page int, pageSize int, location string, environment string, olderThan time.Time) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchClusters(full, strings.Split(search, " "), sortBy, sortDesc, page, pageSize, location, environment, olderThan)
}

// GetCluster return the cluster specified in the clusterName param
func (as *APIService) GetCluster(clusterName string, olderThan time.Time) (*dto.Cluster, utils.AdvancedErrorInterface) {
	return as.Database.GetCluster(clusterName, olderThan)
}

// GetClusterXLSX return  cluster vms as xlxs file
func (as *APIService) GetClusterXLSX(clusterName string, olderThan time.Time) (*excelize.File, error) {
	cluster, aerr := as.Database.GetCluster(clusterName, olderThan)
	if aerr != nil {
		return nil, aerr
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
		xlsx.SetCellValue("VMs", fmt.Sprintf("D%d", i), strconv.FormatBool(val.CappedCPU))
	}

	return xlsx, nil
}
