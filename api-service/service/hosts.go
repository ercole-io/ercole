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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package service is a package that provides methods for querying data
package service

import (
	"fmt"
	"math"
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
	sheetDatabaseEbsDbTier := "Database_&_EBS_DB_Tier"
	sheetHostAdded := "Hosts_added"
	sheetHostDismissed := "Hosts_dismissed"
	j, z := 4, 4 // offset for headers (HostAdded and HostDismissed)
	headerHostCreated, headerHostDismissed := false, false

	hosts, err := as.Database.SearchHosts("lms", filters.SearchHostsFilters)
	if err != nil {
		return nil, utils.NewError(err, "")
	}

	csiByHostname, err := as.getCSIsByHostname()
	if err != nil {
		return nil, utils.NewError(err, "")
	}

	lms, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_lms.xlsm")
	if err != nil {
		aerr := utils.NewError(err, "READ_TEMPLATE")
		return nil, aerr
	}

	if filters.From != utils.MIN_TIME || filters.To != utils.MAX_TIME {
		//HostAdded management
		createdHosts, err := as.Database.GetListValidHostsByRangeDates(filters.From, filters.To)
		if err != nil {
			return nil, utils.NewError(err, "")
		}

		for _, cHostName := range createdHosts {
			createdDate, err := as.Database.GetHostMinValidCreatedAtDate(cHostName)
			if err != nil {
				return nil, utils.NewError(err, "")
			}

			if createdDate.After(filters.From) || createdDate.Equal(filters.From) {
				cFilters := filters.SearchHostsFilters
				cFilters.OlderThan = filters.To
				cFilters.Hostname = cHostName

				cHosts, err := as.Database.SearchHosts("lms", cFilters)
				if err != nil {
					return nil, utils.NewError(err, "")
				}

				for i := 0; i < len(cHosts); i++ {
					if !headerHostCreated {
						indexsheetHostAdded := lms.NewSheet(sheetHostAdded)
						indexSheetDatabaseEbsDbTier := lms.GetSheetIndex(sheetDatabaseEbsDbTier)

						errs := lms.CopySheet(indexSheetDatabaseEbsDbTier, indexsheetHostAdded)
						if errs != nil {
							return nil, errs
						}

						lms.SetActiveSheet(indexSheetDatabaseEbsDbTier)

						headerHostCreated = true
					}

					if cHosts[i]["usingLicenseCount"] != 0.0 {
						setCellValueLMS(lms, sheetHostAdded, j, csiByHostname, cHosts[i])
						j++
					}
				}
			}
		}

		//HostDismissed management
		dismissedHosts, err := as.Database.GetListDismissedHostsByRangeDates(filters.From, filters.To)
		if err != nil {
			return nil, utils.NewError(err, "")
		}

		for _, dHostName := range dismissedHosts {
			existHost, err := as.Database.ExistHostdata(dHostName)
			if err != nil {
				return nil, utils.NewError(err, "")
			}

			if !existHost {
				dFilters := filters.SearchHostsFilters
				dFilters.OlderThan = filters.To
				dFilters.Hostname = dHostName
				dHosts, err := as.Database.SearchHosts("lms", dFilters)

				if err != nil {
					return nil, utils.NewError(err, "")
				}

				for i := 0; i < len(dHosts); i++ {
					if !headerHostDismissed {
						indexsheetHostDismissed := lms.NewSheet(sheetHostDismissed)
						indexSheetDatabaseEbsDbTier := lms.GetSheetIndex(sheetDatabaseEbsDbTier)

						errs := lms.CopySheet(indexSheetDatabaseEbsDbTier, indexsheetHostDismissed)
						if errs != nil {
							return nil, errs
						}

						lms.SetActiveSheet(indexSheetDatabaseEbsDbTier)

						headerHostDismissed = true
					}

					if dHosts[i]["usingLicenseCount"] != 0.0 {
						setCellValueLMS(lms, sheetHostDismissed, z, csiByHostname, dHosts[i])
						z++
					}
				}
			}
		}
	}

	indexRow := 4 // offset for headers

	for i := 0; i < len(hosts); i++ {
		if hosts[i]["usingLicenseCount"] != 0.0 {
			setCellValueLMS(lms, sheetDatabaseEbsDbTier, indexRow, csiByHostname, hosts[i])
			indexRow++
		}
	}

	return lms, nil
}

func (as *APIService) GetHostsMysqlAsLMS(filters dto.SearchHostsAsLMS) (*excelize.File, error) {
	sheetDatabaseEbsDbTier := "Database_&_EBS_DB_Tier"

	hosts, err := as.Database.SearchHostMysqlLMS(filters)
	if err != nil {
		return nil, err
	}

	lms, err := excelize.OpenFile(as.Config.ResourceFilePath + "/templates/template_lms.xlsm")
	if err != nil {
		aerr := utils.NewError(err, "READ_TEMPLATE")
		return nil, aerr
	}

	for i, host := range hosts {
		i += 4
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("B%d", i), host.PhysicalServerName)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("C%d", i), host.VirtualServerName)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("D%d", i), host.Virtualization)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("E%d", i), host.DbInstanceName)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("G%d", i), host.EnvironmentUsage)

		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("N%d", i), host.ProductVersion)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("O%d", i), host.ProductLicenseAllocated)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("P%d", i), host.LicenseMetricAllocated)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("Q%d", i), host.NumberOfLicenseUsed)

		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("AC%d", i), host.ProcessorModel)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("AD%d", i), host.Sockets)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("AF%d", i), host.PhysiclaCores)
		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("AG%d", i), host.ThreadsPerCore)

		lms.SetCellValue(sheetDatabaseEbsDbTier, fmt.Sprintf("AJ%d", i), host.Os)
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
		"Location",
		"Databases",
		"Missing Databases",
		"Technology",
		"Operating System",
		"Clust",
		"Kernel",
		"Memory",
		"Swap",
		"Cluster Type",
		"Cluster Veritas Nodes Count",
		"Cluster Veritas Nodes",
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
		file.SetCellValue(sheet, nextAxis(), val.Location)

		databases := strings.Builder{}
		technology := strings.Builder{}
		missingdbs := strings.Builder{}

		for k, v := range val.Databases {
			databases.WriteString(strings.Join(v, " "))
			technology.WriteString(k)
		}

		for _, v := range val.MissingDatabases {
			missingdbs.WriteString(fmt.Sprintf("%s ", v.Name))
		}

		file.SetCellValue(sheet, nextAxis(), databases.String())
		file.SetCellValue(sheet, nextAxis(), missingdbs.String())
		file.SetCellValue(sheet, nextAxis(), technology.String())

		var os string

		if strings.Contains(val.Info.OS, "Red Hat") {
			os = val.Info.OS + " - " + val.Info.OSVersion
		} else {
			os = val.Info.OS
		}

		file.SetCellValue(sheet, nextAxis(), os)

		isClustered := false
		if val.ClusterMembershipStatus.HACMP ||
			val.ClusterMembershipStatus.OracleClusterware ||
			val.ClusterMembershipStatus.SunCluster ||
			val.ClusterMembershipStatus.VeritasClusterServer {
			isClustered = true
		}

		file.SetCellValue(sheet, nextAxis(), isClustered)

		file.SetCellValue(sheet, nextAxis(), val.Info.KernelVersion)
		file.SetCellValue(sheet, nextAxis(), val.Info.MemoryTotal)
		file.SetCellValue(sheet, nextAxis(), val.Info.SwapTotal)

		clusterType := ""
		clusterVeritasNodesCount := 0
		clusterVeritasNodes := ""

		if val.ClusterMembershipStatus.OracleClusterware {
			clusterType = "ClusterWare"
		}

		if val.ClusterMembershipStatus.VeritasClusterServer {
			clusterType = "VeritasCluster"
			clusterVeritasNodesCount = len(val.ClusterMembershipStatus.VeritasClusterHostnames)
			clusterVeritasNodes = strings.Join(val.ClusterMembershipStatus.VeritasClusterHostnames, ",")
		}

		file.SetCellValue(sheet, nextAxis(), clusterType)
		file.SetCellValue(sheet, nextAxis(), clusterVeritasNodesCount)
		file.SetCellValue(sheet, nextAxis(), clusterVeritasNodes)
	}

	return file, nil
}

func (as *APIService) getCSIsByHostname() (res map[string][]string, err error) {
	contracts, aerr := as.Database.ListOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	if aerr != nil {
		return nil, aerr
	}

	res = make(map[string][]string)

	for i, a := range contracts {
		for _, h := range a.Hosts {
			this := &contracts[i].CSI

			if this != nil && len(*this) > 0 {
				res[h.Hostname] = append(res[h.Hostname], *this)
			}
		}
	}

	return res, nil
}

func (as *APIService) GetHostDataSummaries(filters dto.SearchHostsFilters) ([]dto.HostDataSummary, error) {
	hosts, err := as.Database.GetHostDataSummaries(filters)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

// GetHost return the host specified in the hostname param
func (as *APIService) GetHost(hostname string, olderThan time.Time, raw bool) (*dto.HostData, error) {
	host, err := as.Database.GetHost(hostname, olderThan, raw)
	if err != nil {
		return nil, err
	}

	if host.Features.Oracle != nil &&
		host.Features.Oracle.Database != nil &&
		host.Features.Oracle.Database.Databases != nil {
		for i, n := range host.Features.Oracle.Database.Databases {
			if n.Partitionings == nil {
				host.Features.Oracle.Database.Databases[i].Partitionings = []model.OracleDatabasePartitioning{}
			}

			if n.PDBs != nil {
				for x, y := range n.PDBs {
					if y.SegmentAdvisors == nil {
						host.Features.Oracle.Database.Databases[i].PDBs[x].SegmentAdvisors = []model.OracleDatabaseSegmentAdvisor{}
					}

					if y.Partitionings == nil {
						host.Features.Oracle.Database.Databases[i].PDBs[x].Partitionings = []model.OracleDatabasePartitioning{}
					}
				}
			}

			db := &host.Features.Oracle.Database.Databases[i]
			db.PgaSum = db.GetPgaSum()
			db.SgaSum = db.GetSgaSum()
		}

		as.setMemoryTarget(host)
	}

	return host, nil
}

// ListAllLocations list all available locations
func (as *APIService) ListAllLocations(location string, environment string, olderThan time.Time) ([]string, error) {
	return as.Database.ListAllLocations(location, environment, olderThan)
}

// ListLocations list locations
func (as *APIService) ListLocations(user interface{}) ([]string, error) {
	u := user.(model.User)

	if u.IsAdmin() {
		return as.Database.ListAllLocations("", "", utils.MAX_TIME)
	}

	return as.Database.GetUserLocations(u.Username)
}

// ListEnvironments list environments
func (as *APIService) ListEnvironments(location string, environment string, olderThan time.Time) ([]string, error) {
	return as.Database.ListEnvironments(location, environment, olderThan)
}

// DismissHost dismiss the specified host
func (as *APIService) DismissHost(hostname string) error {
	filter := dto.AlertsFilter{OtherInfo: map[string]interface{}{"hostname": hostname}}
	if err := as.RemoveAlertsNODATA(filter); err != nil {
		as.Log.Errorf("Can't delete alerts by %s", hostname)
	}

	if err := as.AckAlerts(filter); err != nil {
		as.Log.Errorf("Can't ack hostname %s alerts by filter", hostname)
	}

	if err := as.UpdateAlertsStatus(filter, model.AlertStatusDismissed); err != nil {
		as.Log.Errorf("Can't dismiss hostname %s alerts by filter", hostname)
	}

	if err := as.DeleteHostFromOracleDatabaseContracts(hostname); err != nil {
		as.Log.Errorf("Can't remove hostname %s contracts", hostname)
	}

	if err := as.Database.DismissHost(hostname); err != nil {
		return err
	}

	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(as.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryEngine,
		AlertCode:               model.AlertCodeDismissHost,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusNew,
		Date:                    as.TimeNow(),
		Description:             fmt.Sprintf("Host %s was dismissed", hostname),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
		},
	}

	if err := as.AlertSvcClient.ThrowNewAlert(alr); err != nil {
		as.Log.Errorf("Dismiss alert was not added: %s", err)
	}

	return nil
}

func (as *APIService) GetVirtualHostWithoutCluster() ([]dto.VirtualHostWithoutCluster, error) {
	return as.Database.FindVirtualHostWithoutCluster()
}

func (as *APIService) setMemoryTarget(host *dto.HostData) {
	sum := host.GetPGASGASum() + host.GetMemoryTargetSum()
	host.MemorySum = sum
	percentage := (sum * 100) / host.Info.MemoryTotal
	host.MemorySumPercentage = math.Floor(percentage*100) / 100
	flag := "red"

	if as.Config.APIService.PGASGASumTargetPercentage > int(percentage) {
		flag = "green"
	}

	if val, ok := as.Config.APIService.PGASGASumTargetPercentagePerHosts[host.Hostname]; ok {
		flag = "red"

		if val > int(percentage) {
			flag = "green"
		}
	}

	host.MemorySumFlag = flag
}

func checkHosts(as *APIService, hosts []string) error {
	commonFilters := dto.NewSearchHostsFilters()
	notInClusterHosts, err := as.SearchHosts("hostnames",
		commonFilters)

	if err != nil {
		return utils.NewError(err, "")
	}

	notInClusterHostnames := make([]string, len(notInClusterHosts))
	for i, h := range notInClusterHosts {
		notInClusterHostnames[i] = h["hostname"].(string)
	}

hosts_loop:
	for _, host := range hosts {
		for _, notInClusterHostname := range notInClusterHostnames {
			if host == notInClusterHostname {
				continue hosts_loop
			}
		}

		return utils.ErrHostNotFound
	}

	return nil
}
