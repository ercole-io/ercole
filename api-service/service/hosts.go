// Copyright (c) 2021 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) SearchHosts(mode string, filters dto.SearchHostsFilters) ([]map[string]interface{}, error) {
	return as.Database.SearchHosts(mode, filters)
}

func (as *APIService) SearchHostsAsLMS(filters dto.SearchHostsAsLMS) (*excelize.File, error) {
	hosts, err := as.Database.SearchHosts("lms", filters.SearchHostsFilters)
	if err != nil {
		return nil, utils.NewError(err, "")
	}

	sheetDatabaseEbsDbTier := "Database_&_EBS_DB_Tier"
	sheetHostAdded := "Hosts_added"
	sheetHostDismissed := "Hosts_dismissed"

	csiByHostname, err := as.getCSIsByHostname()
	if err != nil {
		return nil, utils.NewError(err, "")
	}

	lms, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_lms.xlsm")
	if err != nil {
		aerr := utils.NewError(err, "READ_TEMPLATE")
		return nil, aerr
	}

	j, z := 4, 4 // offset for headers (HostAdded and HostDismissed)
	for i, val := range hosts {
		i += 4 // offset for headers
		setCellValueLMS(lms, sheetDatabaseEbsDbTier, i, csiByHostname, val)
		createdDate := val["createdAt"].(primitive.DateTime).Time().UTC()
		var dismissedAt time.Time
		if val["dismissedAt"] != nil {
			dismissedAt = val["dismissedAt"].(primitive.DateTime).Time().UTC()
		}
		//HostAdded management
		if (filters.From != utils.MIN_TIME ||
			filters.To != utils.MAX_TIME) &&
			createdDate.After(filters.From) &&
			createdDate.Before(filters.To) {
			if j == 4 {
				indexsheetHostAdded := lms.NewSheet(sheetHostAdded)
				indexSheetDatabaseEbsDbTier := lms.GetSheetIndex(sheetDatabaseEbsDbTier)
				errs := lms.CopySheet(indexSheetDatabaseEbsDbTier, indexsheetHostAdded)
				if errs != nil {
					return nil, errs
				}
				lms.SetActiveSheet(indexSheetDatabaseEbsDbTier)
			}
			setCellValueLMS(lms, sheetHostAdded, j, csiByHostname, val)
			j++
		}

		//HostDismissed management
		if (filters.From != utils.MIN_TIME ||
			filters.To != utils.MAX_TIME) &&
			dismissedAt.After(filters.From) &&
			dismissedAt.Before(filters.To) {
			if z == 4 {
				indexsheetHostDismissed := lms.NewSheet(sheetHostDismissed)
				indexSheetDatabaseEbsDbTier := lms.GetSheetIndex(sheetDatabaseEbsDbTier)
				errs := lms.CopySheet(indexSheetDatabaseEbsDbTier, indexsheetHostDismissed)
				if errs != nil {
					return nil, errs
				}
				lms.SetActiveSheet(indexSheetDatabaseEbsDbTier)
			}
			setCellValueLMS(lms, sheetHostDismissed, z, csiByHostname, val)
			z++
		}
	}

	return lms, nil
}

func setCellValueLMS(lms *excelize.File, sheetName string, i int, csiByHostname map[string][]string, val map[string]interface{}) {
	lms.SetCellValue(sheetName, fmt.Sprintf("B%d", i), val["physicalServerName"])
	lms.SetCellValue(sheetName, fmt.Sprintf("C%d", i), val["virtualServerName"])
	lms.SetCellValue(sheetName, fmt.Sprintf("D%d", i), val["virtualizationTechnology"])
	lms.SetCellValue(sheetName, fmt.Sprintf("E%d", i), val["dbInstanceName"])
	lms.SetCellValue(sheetName, fmt.Sprintf("F%d", i), val["pluggableDatabaseName"])
	lms.SetCellValue(sheetName, fmt.Sprintf("G%d", i), val["environment"])
	lms.SetCellValue(sheetName, fmt.Sprintf("H%d", i), val["options"])
	lms.SetCellValue(sheetName, fmt.Sprintf("I%d", i), val["usedManagementPacks"])
	lms.SetCellValue(sheetName, fmt.Sprintf("N%d", i), val["productVersion"])
	lms.SetCellValue(sheetName, fmt.Sprintf("O%d", i), val["productLicenseAllocated"])
	lms.SetCellValue(sheetName, fmt.Sprintf("P%d", i), val["licenseMetricAllocated"])
	lms.SetCellValue(sheetName, fmt.Sprintf("Q%d", i), val["usingLicenseCount"])

	hostname := val["physicalServerName"].(string)
	if len(hostname) == 0 {
		hostname = val["virtualServerName"].(string)
	}
	if csi, ok := csiByHostname[hostname]; ok {
		lms.SetCellValue(sheetName, fmt.Sprintf("R%d", i), strings.Join(csi, ", "))
	}

	lms.SetCellValue(sheetName, fmt.Sprintf("AC%d", i), val["processorModel"])
	lms.SetCellValue(sheetName, fmt.Sprintf("AD%d", i), val["processors"])
	lms.SetCellValue(sheetName, fmt.Sprintf("AE%d", i), val["coresPerProcessor"])
	lms.SetCellValue(sheetName, fmt.Sprintf("AF%d", i), val["physicalCores"])
	lms.SetCellValue(sheetName, fmt.Sprintf("AG%d", i), val["threadsPerCore"])
	lms.SetCellValue(sheetName, fmt.Sprintf("AH%d", i), val["processorSpeed"])
	lms.SetCellValue(sheetName, fmt.Sprintf("AJ%d", i), val["operatingSystem"])
}

func (as *APIService) SearchHostsAsXLSX(filters dto.SearchHostsFilters) (*excelize.File, error) {
	hosts, err := as.Database.GetHostDataSummaries(filters)
	if err != nil {
		return nil, err
	}

	sheet := "Hosts"
	headers := []string{
		"Hostname",
		"Platform",
		"Cluster",
		"Node",
		"Processor Model",
		"Threads",
		"Cores",
		"Socket",
		"Version",
		"Updated",
		"Environment",
		"Databases",
		"Technology",
		"Operating System",
		"Clust",
		"Kernel",
		"Memory",
		"Swap",
	}

	file, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)
	for _, val := range hosts {
		nextAxis := axisHelp.NewRow()
		file.SetCellValue(sheet, nextAxis(), val.Hostname)

		if val.Info.HardwareAbstractionTechnology == "PH" {
			file.SetCellValue(sheet, nextAxis(), "Bare metal")
		} else {
			file.SetCellValue(sheet, nextAxis(), val.Info.HardwareAbstractionTechnology)
		}

		file.SetCellValue(sheet, nextAxis(), val.Cluster)
		file.SetCellValue(sheet, nextAxis(), val.VirtualizationNode)
		file.SetCellValue(sheet, nextAxis(), val.Info.CPUModel)
		file.SetCellValue(sheet, nextAxis(), val.Info.CPUThreads)
		file.SetCellValue(sheet, nextAxis(), val.Info.CPUCores)
		file.SetCellValue(sheet, nextAxis(), val.Info.CPUSockets)
		file.SetCellValue(sheet, nextAxis(), val.AgentVersion)
		file.SetCellValue(sheet, nextAxis(), val.CreatedAt)
		file.SetCellValue(sheet, nextAxis(), val.Environment)

		databases := strings.Builder{}
		technology := strings.Builder{}
		for k, v := range val.Databases {
			databases.WriteString(strings.Join(v, " "))
			technology.WriteString(k)
		}
		file.SetCellValue(sheet, nextAxis(), databases.String())
		file.SetCellValue(sheet, nextAxis(), technology.String())
		file.SetCellValue(sheet, nextAxis(), val.Info.OS)
		file.SetCellValue(sheet, nextAxis(), val.ClusterMembershipStatus.OracleClusterware)
		file.SetCellValue(sheet, nextAxis(), val.Info.KernelVersion)
		file.SetCellValue(sheet, nextAxis(), val.Info.MemoryTotal)
		file.SetCellValue(sheet, nextAxis(), val.Info.SwapTotal)
	}

	return file, nil
}

func (as *APIService) getCSIsByHostname() (res map[string][]string, err error) {
	agreements, aerr := as.Database.ListOracleDatabaseAgreements()
	if aerr != nil {
		return nil, aerr
	}

	res = make(map[string][]string)

	for i, a := range agreements {
		for _, h := range a.Hosts {
			this := &agreements[i].CSI

			if this != nil && len(*this) > 0 {
				res[h.Hostname] = append(res[h.Hostname], *this)
			}
		}
	}

	return res, nil
}

func (as *APIService) GetHostDataSummaries(filters dto.SearchHostsFilters) ([]dto.HostDataSummary, error) {
	return as.Database.GetHostDataSummaries(filters)
}

// GetHost return the host specified in the hostname param
func (as *APIService) GetHost(hostname string, olderThan time.Time, raw bool) (*dto.HostData, error) {
	host, err := as.Database.GetHost(hostname, olderThan, raw)
	if err != nil {
		return nil, err
	}

	var realApplicationClusters bool
	var indexList []int
	if host.Features.Oracle != nil && host.Features.Oracle.Database != nil && host.Features.Oracle.Database.Databases != nil {
		for i := range host.Features.Oracle.Database.Databases {
			db := &host.Features.Oracle.Database.Databases[i]
			for j := range db.Licenses {
				lic := db.Licenses[j]
				if lic.LicenseTypeID != "" {
					licType, err := as.GetOracleDatabaseLicenseType(lic.LicenseTypeID)
					if err != nil {
						return nil, err
					}
					if licType.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
						count := lic.Count
						count *= float64(model.GetFactorByMetric(licType.Metric))
						host.Features.Oracle.Database.Databases[i].Licenses[j].Count = count
					}

					if lic.Name == "Real Application Clusters" && lic.Count > 0 {
						realApplicationClusters = true
					}

					if lic.Name == "Real Application Clusters One Node" && lic.Count > 0 {
						indexList = append(indexList, j)
					}
				}
			}

			if len(indexList) > 0 && realApplicationClusters {
				for _, k := range indexList {
					host.Features.Oracle.Database.Databases[i].Licenses[k].Count = 0
				}
				realApplicationClusters = false
				indexList = nil
			}
		}
	}

	return host, nil
}

// ListLocations list locations
func (as *APIService) ListLocations(location string, environment string, olderThan time.Time) ([]string, error) {
	return as.Database.ListLocations(location, environment, olderThan)
}

// ListEnvironments list environments
func (as *APIService) ListEnvironments(location string, environment string, olderThan time.Time) ([]string, error) {
	return as.Database.ListEnvironments(location, environment, olderThan)
}

// DismissHost dismiss the specified host
func (as *APIService) DismissHost(hostname string) error {
	filter := dto.AlertsFilter{OtherInfo: map[string]interface{}{"hostname": hostname}}
	if err := as.AckAlerts(filter); err != nil {
		as.Log.Errorf("Can't ack hostname %s alerts by filter", hostname)
	}

	return as.Database.DismissHost(hostname)
}
