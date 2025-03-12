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

// Package service is a package that provides methods for querying data
package service

import (
	"errors"
	"sort"
	"strings"

	"github.com/ercole-io/ercole/v2/utils"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) GetDatabaseConnectionStatus() bool {
	err := as.Database.CheckStatusMongodb()
	return err == nil
}

func (as *APIService) SearchDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	type getter func(filter dto.GlobalFilter) ([]dto.Database, error)

	getters := []getter{as.getOracleDatabases, as.getMySQLDatabases, as.getSqlServerDatabases, as.getPostgreSqlDatabases, as.getMongoDBDatabases}

	dbs := make([]dto.Database, 0)

	for _, get := range getters {
		thisDbs, err := get(filter)
		if err != nil {
			return nil, err
		}

		dbs = append(dbs, thisDbs...)
	}

	return dbs, nil
}

func (as *APIService) getOracleDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	sodf := dto.SearchOracleDatabasesFilter{
		GlobalFilter: filter,
		PageNumber:   -1,
		PageSize:     -1,
	}

	oracleDbs, err := as.SearchOracleDatabases(sodf)
	if err != nil {
		return nil, err
	}

	dbs := make([]dto.Database, 0)

	for _, oracleDb := range oracleDbs.Content {
		db := dto.Database{
			Name:             oracleDb.Name,
			Type:             model.TechnologyOracleDatabase,
			Version:          oracleDb.Version,
			Hostname:         oracleDb.Hostname,
			Environment:      oracleDb.Environment,
			Location:         oracleDb.Location,
			Charset:          oracleDb.Charset,
			Memory:           oracleDb.Memory,
			DatafileSize:     oracleDb.DatafileSize,
			SegmentsSize:     oracleDb.SegmentsSize,
			Archivelog:       oracleDb.Archivelog,
			HighAvailability: oracleDb.Ha,
			DisasterRecovery: oracleDb.Dataguard,
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (as *APIService) getMySQLDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	mysqlInstances, err := as.Database.SearchMySQLInstances(filter)
	if err != nil {
		return nil, err
	}

	dbs := make([]dto.Database, 0)

	for _, instance := range mysqlInstances {
		segmentsSize := 0.0
		for _, ts := range instance.TableSchemas {
			segmentsSize += ts.Allocation
		}

		db := dto.Database{
			Name:             instance.Name,
			Type:             model.TechnologyOracleMySQL,
			Version:          instance.Version,
			Hostname:         instance.Hostname,
			Environment:      instance.Environment,
			Location:         instance.Location,
			Charset:          instance.CharsetServer,
			Memory:           instance.BufferPoolSize / 1024,
			DatafileSize:     0,
			SegmentsSize:     segmentsSize / 1024,
			Archivelog:       instance.LogBin,
			HighAvailability: instance.HighAvailability,
			DisasterRecovery: instance.IsMaster || instance.IsSlave,
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (as *APIService) getSqlServerDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	sodf := dto.SearchSqlServerInstancesFilter{
		GlobalFilter: filter,
		PageNumber:   -1,
		PageSize:     -1,
	}

	sqlServerInstances, err := as.SearchSqlServerInstances(sodf)
	if err != nil {
		return nil, err
	}

	dbs := make([]dto.Database, 0)

	for _, instance := range sqlServerInstances.Content {
		db := dto.Database{
			Name:        instance.Name,
			Type:        model.TechnologyMicrosoftSQLServer,
			Version:     instance.Version,
			Hostname:    instance.Hostname,
			Environment: instance.Environment,
			Location:    instance.Location,
			Charset:     instance.CollationName,
		}
		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (as *APIService) getPostgreSqlDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	sodf := dto.SearchPostgreSqlInstancesFilter{
		GlobalFilter: filter,
		PageNumber:   -1,
		PageSize:     -1,
	}

	postgreSqlInstances, err := as.SearchPostgreSqlInstances(sodf)
	if err != nil {
		return nil, err
	}

	dbs := make([]dto.Database, 0)

	for _, instance := range postgreSqlInstances.Content {
		db := dto.Database{
			Name:        instance.Name,
			Type:        model.TechnologyPostgreSQLPostgreSQL,
			Version:     instance.Version,
			Hostname:    instance.Hostname,
			Environment: instance.Environment,
			Location:    instance.Location,
			Charset:     instance.Charset,
		}
		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (as *APIService) getMongoDBDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	sodf := dto.SearchMongoDBInstancesFilter{
		GlobalFilter: filter,
		PageNumber:   -1,
		PageSize:     -1,
	}

	mongoDBInstances, err := as.SearchMongoDBInstances(sodf)
	if err != nil {
		return nil, err
	}

	dbs := make([]dto.Database, 0)
	setUnique := make(map[string]dto.MongoDBInstance)

	for _, instance := range mongoDBInstances.Content {
		if _, ok := setUnique[instance.InstanceName]; !ok {
			db := dto.Database{
				Name:        instance.InstanceName,
				Type:        model.TechnologyMongoDBMongoDB,
				Version:     instance.Version,
				Hostname:    instance.Hostname,
				Environment: instance.Environment,
				Location:    instance.Location,
				Charset:     instance.Charset,
			}
			dbs = append(dbs, db)
			setUnique[instance.InstanceName] = instance
		}
	}

	return dbs, nil
}

func (as *APIService) SearchDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	databases, err := as.SearchDatabases(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Databases"
	headers := []string{
		"Name",
		"Type",
		"Version",
		"Hostname",
		"Environment",
		"Location",
		"Charset",
		"Memory",
		"Datafile Size",
		"Segments Size",
	}

	file, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)
	for _, val := range databases {
		nextAxis := axisHelp.NewRow()

		file.SetCellValue(sheet, nextAxis(), val.Name)
		file.SetCellValue(sheet, nextAxis(), val.Type)
		file.SetCellValue(sheet, nextAxis(), val.Version)
		file.SetCellValue(sheet, nextAxis(), val.Hostname)
		file.SetCellValue(sheet, nextAxis(), val.Environment)
		file.SetCellValue(sheet, nextAxis(), val.Location)
		file.SetCellValue(sheet, nextAxis(), val.Charset)
		file.SetCellValue(sheet, nextAxis(), val.Memory)
		file.SetCellValue(sheet, nextAxis(), val.DatafileSize)
		file.SetCellValue(sheet, nextAxis(), val.SegmentsSize)
	}

	return file, nil
}

func (as *APIService) GetDatabasesStatistics(filter dto.GlobalFilter) (*dto.DatabasesStatistics, error) {
	dbs, err := as.SearchDatabases(filter)
	if err != nil {
		return nil, err
	}

	stats := new(dto.DatabasesStatistics)
	for _, db := range dbs {
		stats.TotalMemorySize += db.Memory * 1024 * 1024 * 1024         // From GBytes to bytes
		stats.TotalSegmentsSize += db.SegmentsSize * 1024 * 1024 * 1024 // From GBytes to bytes
	}

	return stats, nil
}

func (as *APIService) GetUsedLicensesPerDatabases(hostname string, filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	type getter func(hostname string, filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error)

	getters := []getter{as.getOracleDatabasesUsedLicenses, as.getMySQLUsedLicenses, as.getSqlServerDatabasesUsedLicenses}

	usedLicenses := make([]dto.DatabaseUsedLicense, 0)

	for _, get := range getters {
		thisDbs, err := get(hostname, filter)
		if err != nil {
			return nil, err
		}

		usedLicenses = append(usedLicenses, thisDbs...)
	}

	return usedLicenses, nil
}

func (as *APIService) clusterLicenses(license dto.DatabaseUsedLicense, clusters []dto.Cluster) (float64, *dto.Cluster, error) {
	clusterByHostnames := make(map[string]*dto.Cluster)

	for i := range clusters {
		for j := range clusters[i].VMs {
			clusterByHostnames[clusters[i].VMs[j].Hostname] = &clusters[i]
		}
	}

	cluster, found := clusterByHostnames[license.Hostname]
	if !found {
		return 0, nil, utils.ErrHostNotInCluster
	}

	return float64(cluster.CPU) * 0.5, cluster, nil
}

func (as *APIService) veritasClusterLicenses(hostdata *model.HostDataBE, hostdatasPerHostname map[string]*model.HostDataBE) (float64, string, string, error) {
	clusterCores, err := hostdata.GetClusterCores(hostdatasPerHostname)

	if errors.Is(err, utils.ErrHostNotInCluster) {
		return 0, "", "", utils.ErrHostNotInCluster
	} else if err != nil {
		return 0, "", "", err
	}

	hostnames := hostdata.ClusterMembershipStatus.VeritasClusterHostnames
	sort.Slice(hostnames, func(i, j int) bool {
		return hostnames[i] < hostnames[j]
	})

	clusterName := strings.Join(hostnames, ",")

	return float64(clusterCores) * hostdata.CoreFactor(), clusterName, "VeritasCluster", nil
}

func (as *APIService) GetUsedLicensesPerDatabasesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	licenses, err := as.GetUsedLicensesPerDatabases("", filter)
	if err != nil {
		return nil, err
	}

	sheet := "Licenses Used"
	headers := []string{
		"Hostname",
		"DB Name",
		"Part Number",
		"Description",
		"Metric",
		"Used Licenses",
		"Cluster Licenses",
		"Ignored",
		"Ignored Comment",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range licenses {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.DbName)
		sheets.SetCellValue(sheet, nextAxis(), val.LicenseTypeID)
		sheets.SetCellValue(sheet, nextAxis(), val.Description)
		sheets.SetCellValue(sheet, nextAxis(), val.Metric)
		sheets.SetCellValue(sheet, nextAxis(), val.UsedLicenses)
		sheets.SetCellValue(sheet, nextAxis(), val.ClusterLicenses)
		sheets.SetCellValue(sheet, nextAxis(), val.Ignored)
		sheets.SetCellValue(sheet, nextAxis(), val.IgnoredComment)
	}

	return sheets, err
}

func (as *APIService) getSqlServerDatabasesUsedLicenses(hostname string, filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	sqlServerLics, err := as.GetSqlServerUsedLicenses(hostname, filter)
	if err != nil {
		return nil, err
	}

	licenseTypes, err := as.GetSqlServerDatabaseLicenseTypesAsMap()
	if err != nil {
		return nil, err
	}

	genericLics := make([]dto.DatabaseUsedLicense, 0, len(sqlServerLics.Content))

	for _, lic := range sqlServerLics.Content {
		lt := licenseTypes[lic.LicenseTypeID]

		g := dto.DatabaseUsedLicense{
			Hostname:       lic.Hostname,
			DbName:         lic.DbName,
			LicenseTypeID:  lic.LicenseTypeID,
			Description:    lt.ItemDescription,
			Metric:         lic.ContractType,
			UsedLicenses:   lic.UsedLicenses,
			Ignored:        lic.Ignored,
			IgnoredComment: lic.IgnoredComment,
		}

		genericLics = append(genericLics, g)
	}

	return genericLics, nil
}

func (as *APIService) getOracleDatabasesUsedLicenses(hostname string, filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	oracleLics, err := as.Database.SearchOracleDatabaseUsedLicenses(hostname, "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	licenseTypes, err := as.GetOracleDatabaseLicenseTypesAsMap()
	if err != nil {
		return nil, err
	}

	usedLicenses := make([]dto.DatabaseUsedLicense, 0, len(oracleLics.Content))

	for _, o := range oracleLics.Content {
		lt := licenseTypes[o.LicenseTypeID]

		g := dto.DatabaseUsedLicense{
			Hostname:       o.Hostname,
			DbName:         o.DbName,
			LicenseTypeID:  o.LicenseTypeID,
			Description:    lt.ItemDescription,
			Metric:         lt.Metric,
			UsedLicenses:   o.UsedLicenses,
			Ignored:        o.Ignored,
			IgnoredComment: o.IgnoredComment,
		}

		usedLicenses = append(usedLicenses, g)
	}

	hostdatas, err := as.Database.GetHostDatas(dto.GlobalFilter{
		OlderThan: utils.MAX_TIME,
	})
	if err != nil {
		return nil, err
	}

	hostdatasPerHostname := make(map[string]*model.HostDataBE, len(hostdatas))
	hostdatasMap := make(map[string]model.HostDataBE, len(hostdatas))

	for i := range hostdatas {
		hd := &hostdatas[i]
		hostdatasPerHostname[hd.Hostname] = hd
		hostdatasMap[hd.Hostname] = *hd
	}

	clusters, err := as.Database.GetClusters(dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	})
	if err != nil {
		return nil, err
	}

	clustersMap := make(map[string]dto.Cluster, len(clusters))
	for _, cluster := range clusters {
		clustersMap[cluster.Name] = cluster
	}

	hypervisorLicenses := make([]dto.DatabaseUsedLicense, 0)

	for i, l := range usedLicenses {
		if usedLicenses[i].Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
			usedLicenses[i].UsedLicenses *= model.GetFactorByMetric(usedLicenses[i].Metric)
		}

		hostdata, found := hostdatasPerHostname[l.Hostname]
		if !found {
			as.Log.Errorf("%v: %s", utils.ErrHostNotFound, l.Hostname)
			continue
		}

		consumedLicenses, cluster, err := as.clusterLicenses(l, clusters)
		if err != nil && !errors.Is(err, utils.ErrHostNotInCluster) {
			return nil, err
		} else if !errors.Is(err, utils.ErrHostNotInCluster) {
			usedLicenses[i].ClusterLicenses = consumedLicenses * model.GetFactorByMetric(usedLicenses[i].Metric)
			usedLicenses[i].ClusterName = cluster.Name
			usedLicenses[i].ClusterType = cluster.Type

			isCapped, err := as.manageLicenseWithCappedCPU(usedLicenses[i], clustersMap, hostdatasMap)
			if err != nil {
				return nil, err
			}

			usedLicenses[i].OlvmCapped = isCapped

			hypervisorLicenses = append(hypervisorLicenses, usedLicenses[i])

			continue
		}

		consumedLicenses, clusterName, clusterType, err := as.veritasClusterLicenses(hostdata, hostdatasPerHostname)
		if err != nil && !errors.Is(err, utils.ErrHostNotInCluster) {
			return nil, err
		} else if !errors.Is(err, utils.ErrHostNotInCluster) {
			usedLicenses[i].ClusterLicenses = consumedLicenses * model.GetFactorByMetric(usedLicenses[i].Metric)
			usedLicenses[i].ClusterName = clusterName
			usedLicenses[i].ClusterType = clusterType
			continue
		}
	}

	usedLicenses = as.removeLicensesByDependencies(usedLicenses, hostdatasPerHostname, clusters)

	usedLicenses = as.manageStandardDBVersionLicenses(usedLicenses, clusters, hostdatasPerHostname)

	as.CalcVeritasClusterLicenses(usedLicenses)

	errHypervisor := as.checkOlvmCappedHypervisorLicenses(hypervisorLicenses, clustersMap)
	if errHypervisor != nil {
		return nil, errHypervisor
	}

	return usedLicenses, nil
}

var goldenGateIds []string = []string{"L75978", "L75967"}
var activeDataguardIds []string = []string{"L47210", "L47217"}

var racIds []string = []string{"L10005", "A90619"}
var racOneNodeIds []string = []string{"L76084", "L76094"}

func (as *APIService) removeLicensesByDependencies(usedLicenses []dto.DatabaseUsedLicense, hostdatasPerHostname map[string]*model.HostDataBE, clusters []dto.Cluster) []dto.DatabaseUsedLicense {
	dependencies := []struct {
		given  []string // If a "given" licenseTypeID is found
		remove []string // Remove any "remove" licenseTypeID from host and cluster
	}{
		{
			given:  goldenGateIds,
			remove: activeDataguardIds,
		},
		{
			given:  racIds,
			remove: racOneNodeIds,
		},
	}

	for _, d := range dependencies {
		indexHosts := make(map[string]bool)

		for i := range usedLicenses {
			for _, givenId := range d.given {
				if usedLicenses[i].LicenseTypeID == givenId {
					indexHosts[usedLicenses[i].Hostname] = true
				}
			}
		}

		for hostname := range indexHosts {
		clusters:
			for _, cluster := range clusters {
				for _, vm := range cluster.VMs {
					if vm.Hostname == hostname {
						for _, x := range cluster.VMs {
							indexHosts[x.Hostname] = true
						}
						break clusters
					}
				}
			}
		}

		for hostname := range indexHosts {
			hostdata, ok := hostdatasPerHostname[hostname]

			if !ok || hostdata == nil {
				continue
			}

			if hostdata.ClusterMembershipStatus.VeritasClusterServer {
				for _, hostVeritasCluster := range hostdata.ClusterMembershipStatus.VeritasClusterHostnames {
					indexHosts[hostVeritasCluster] = true
				}
			}
		}

	licenses:
		for i := 0; i < len(usedLicenses); {
			l := &usedLicenses[i]

			if _, ok := indexHosts[l.Hostname]; !ok {
				i++
				continue
			}

			for _, r := range d.remove {
				if l.LicenseTypeID == r {
					usedLicenses = append(usedLicenses[:i], usedLicenses[i+1:]...)
					continue licenses
				}
			}

			i++
		}
	}

	return usedLicenses
}

func (as *APIService) manageStandardDBVersionLicenses(usedLicenses []dto.DatabaseUsedLicense, clusters []dto.Cluster, hostdatas map[string]*model.HostDataBE) []dto.DatabaseUsedLicense {
	clustersMap := make(map[string]dto.Cluster, len(clusters))
	for _, cluster := range clusters {
		clustersMap[cluster.Name] = cluster
	}

	for i, usedlicense := range usedLicenses {
		if usedlicense.ClusterName == "" {
			continue
		}

		host, ok := hostdatas[usedlicense.Hostname]
		if !ok {
			as.Log.Warnf("%s : %s", utils.ErrHostNotFound, usedlicense.Hostname)
			continue
		}

		if host != nil &&
			host.Features.Oracle != nil &&
			host.Features.Oracle.Database != nil &&
			host.Features.Oracle.Database.Databases != nil {
			cluster, ok := clustersMap[usedlicense.ClusterName]
			if !ok {
				// as.Log.Warnf("%s : %s", utils.ErrClusterNotFound, usedlicense.ClusterName)
				continue
			}

			databases := host.Features.Oracle.Database.Databases
			for _, database := range databases {
				for _, license := range database.Licenses {
					if license.LicenseTypeID == usedlicense.LicenseTypeID &&
						database.Name == usedlicense.DbName &&
						database.Edition() == model.OracleDatabaseEditionStandard {
						usedLicenses[i].ClusterLicenses = float64(cluster.Sockets) * model.GetFactorByMetric(usedlicense.Metric)
					}
				}
			}
		}
	}

	return usedLicenses
}

func (as *APIService) CalcVeritasClusterLicenses(usedLicenses []dto.DatabaseUsedLicense) {
	for i := 0; i < len(usedLicenses); i++ {
		ul := &usedLicenses[i]

		if ul.LicenseTypeID == "L47837" && ul.ClusterType == "VeritasCluster" {
			used := float64(len(strings.Split(ul.ClusterName, ",")))
			ul.UsedLicenses = 1
			ul.ClusterLicenses = used
		}
	}
}

func (as *APIService) getMySQLUsedLicenses(hostname string, filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	mysqlLics, err := as.GetMySQLUsedLicenses(hostname, filter)
	if err != nil {
		return nil, err
	}

	genericLics := make([]dto.DatabaseUsedLicense, 0, len(mysqlLics))

	for _, lic := range mysqlLics {
		g := dto.DatabaseUsedLicense{
			Hostname:       lic.Hostname,
			DbName:         lic.InstanceName,
			LicenseTypeID:  lic.LicenseTypeID,
			Description:    lic.InstanceEdition,
			Metric:         lic.ContractType,
			UsedLicenses:   lic.UsedLicenses,
			Ignored:        lic.Ignored,
			IgnoredComment: lic.IgnoredComment,
		}

		genericLics = append(genericLics, g)
	}

	return genericLics, nil
}

func (as *APIService) GetDatabaseLicensesCompliance(locations []string) ([]dto.LicenseCompliance, error) {
	licenses := make([]dto.LicenseCompliance, 0)

	oracle, err := as.GetOracleDatabaseLicensesCompliance(locations)
	if err != nil {
		return nil, err
	}

	licenses = append(licenses, oracle...)

	mysql, err := as.GetMySQLDatabaseLicensesCompliance(locations)
	if err != nil {
		return nil, err
	}

	licenses = append(licenses, mysql...)

	sqlServer, err := as.GetSqlServerDatabaseLicensesCompliance(locations)
	if err != nil {
		return nil, err
	}

	licenses = append(licenses, sqlServer...)

	for i := 0; i < len(licenses); {
		l := licenses[i]

		if l.Covered == 0 && l.Consumed == 0 {
			licenses = append(licenses[0:i], licenses[i+1:]...)
			continue
		}

		i++
	}

	return licenses, nil
}

func (as *APIService) GetDatabaseLicensesComplianceAsXLSX(locations []string) (*excelize.File, error) {
	licenses, err := as.GetDatabaseLicensesCompliance(locations)
	if err != nil {
		return nil, err
	}

	sheet := "Licenses Compliance"
	headers := []string{
		"Part Number",
		"Description",
		"Metric",
		"License Available",
		"Purchased",
		"Consumed",
		"Covered",
		"Compliance",
		"ULA",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range licenses {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.LicenseTypeID)
		sheets.SetCellValue(sheet, nextAxis(), val.ItemDescription)
		sheets.SetCellValue(sheet, nextAxis(), val.Metric)
		sheets.SetCellValue(sheet, nextAxis(), val.Available)
		sheets.SetCellValue(sheet, nextAxis(), val.Purchased)
		sheets.SetCellValue(sheet, nextAxis(), val.Consumed)
		sheets.SetCellValue(sheet, nextAxis(), val.Covered)
		sheets.SetCellValue(sheet, nextAxis(), val.Compliance)
		sheets.SetCellValue(sheet, nextAxis(), val.Unlimited)
	}

	return sheets, err
}

func (as *APIService) GetUsedLicensesPerHostAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	usedLicenses, err := as.GetUsedLicensesPerHost(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Licenses Used Per Host"
	headers := []string{
		"Hostname",
		"Databases",
		"Database Names",
		"Part Number",
		"Description",
		"Metric",
		"Used Licenses",
		"Cluster Licenses",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range usedLicenses {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), len(val.DatabaseNames))
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.DatabaseNames, ", "))
		sheets.SetCellValue(sheet, nextAxis(), val.LicenseTypeID)
		sheets.SetCellValue(sheet, nextAxis(), val.Description)
		sheets.SetCellValue(sheet, nextAxis(), val.Metric)
		sheets.SetCellValue(sheet, nextAxis(), val.UsedLicenses)
		sheets.SetCellValue(sheet, nextAxis(), val.ClusterLicenses)
	}

	return sheets, err
}

func (as *APIService) GetUsedLicensesPerHost(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicensePerHost, error) {
	licenses, err := as.GetUsedLicensesPerDatabases("", filter)
	if err != nil {
		return nil, err
	}

	hostdatas, err := as.Database.GetHostDatas(dto.GlobalFilter{
		OlderThan: utils.MAX_TIME,
	})
	if err != nil {
		return nil, err
	}

	hostdatasPerHostname := make(map[string]*model.HostDataBE, len(hostdatas))
	hostdatasMap := make(map[string]model.HostDataBE, len(hostdatas))

	for i := range hostdatas {
		hd := &hostdatas[i]
		hostdatasPerHostname[hd.Hostname] = hd
		hostdatasMap[hd.Hostname] = *hd
	}

	var licensesPerHost []dto.DatabaseUsedLicensePerHost

licenses:
	for _, v := range licenses {
		if v.Ignored {
			continue
		}

		for i, v2 := range licensesPerHost {
			if v.Hostname == v2.Hostname && v.LicenseTypeID == v2.LicenseTypeID {
				licensesPerHost[i].DatabaseNames = append(licensesPerHost[i].DatabaseNames, v.DbName)
				continue licenses
			}
		}

		var clusterLicenses float64

		clustersMap := make(map[string]dto.Cluster, 0)

		if v.ClusterName != "" && v.ClusterType != "VeritasCluster" {
			cluster, err := as.GetCluster(v.ClusterName, utils.MAX_TIME)
			if err != nil {
				continue licenses
			}

			clustersMap[cluster.Name] = *cluster

			for _, hostVM := range cluster.VMs {
				if hostVM.CappedCPU {
					host, err := as.GetHost(hostVM.Hostname, utils.MAX_TIME, false)
					if err != nil {
						continue
					}
					if host != nil &&
						host.Features.Oracle != nil &&
						host.Features.Oracle.Database != nil &&
						host.Features.Oracle.Database.Databases != nil {

						databases := host.Features.Oracle.Database.Databases
						for _, database := range databases {
							for _, license := range database.Licenses {
								if license.LicenseTypeID == v.LicenseTypeID &&
									database.Name == v.DbName {
									if database.Edition() == model.OracleDatabaseEditionStandard {
										clusterLicenses = float64(cluster.Sockets) * model.GetFactorByMetric(v.Metric)
									} else {
										clusterLicenses = 0
									}

								}
							}
						}
					}

				} else {
					clusterLicenses = v.ClusterLicenses
					break
				}

			}
		}

		isCapped, err := as.manageLicenseWithCappedCPU(v, clustersMap, hostdatasMap)
		if err != nil {
			return nil, err
		}

		licensesPerHost = append(licensesPerHost,
			dto.DatabaseUsedLicensePerHost{
				Hostname:        v.Hostname,
				DatabaseNames:   []string{v.DbName},
				LicenseTypeID:   v.LicenseTypeID,
				Description:     v.Description,
				Metric:          v.Metric,
				UsedLicenses:    v.UsedLicenses,
				ClusterLicenses: clusterLicenses,
				OlvmCapped:      isCapped,
			},
		)
	}

	return licensesPerHost, nil
}

func (as *APIService) GetUsedLicensesPerCluster(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicensePerCluster, error) {
	licenses, err := as.GetUsedLicensesPerDatabases("", filter)
	if err != nil {
		return nil, err
	}

	clusters, err := as.Database.GetClusters(filter)
	if err != nil {
		return nil, err
	}

	clusterByHostnames := make(map[string]*dto.Cluster)

	for i := range clusters {
		for j := range clusters[i].VMs {
			clusterByHostnames[clusters[i].VMs[j].Hostname] = &clusters[i]
		}
	}

	// By cluster.Hostname and by LicenseTypeID
	m := make(map[string]map[string]*dto.DatabaseUsedLicensePerCluster)

licenses:
	for _, l := range licenses {
		c, ok := clusterByHostnames[l.Hostname]
		if !ok {
			continue licenses
		}

		clusterLicenses, ok := m[c.Name]
		if !ok {
			clusterLicenses = make(map[string]*dto.DatabaseUsedLicensePerCluster)
			m[c.Name] = clusterLicenses
		}

		ll, ok := clusterLicenses[l.LicenseTypeID]
		if !ok {
			ll = &dto.DatabaseUsedLicensePerCluster{
				Cluster:       c.Name,
				Hostnames:     []string{},
				LicenseTypeID: l.LicenseTypeID,
				Description:   l.Description,
				Metric:        l.Metric,
				UsedLicenses:  l.ClusterLicenses,
			}

			clusterLicenses[l.LicenseTypeID] = ll
		}

		for _, h := range ll.Hostnames {
			if l.Hostname == h {
				continue licenses
			}
		}
		ll.Hostnames = append(ll.Hostnames, l.Hostname)
	}

	result := make([]dto.DatabaseUsedLicensePerCluster, 0)

	for i := range m {
		for j := range m[i] {
			result = append(result, *m[i][j])
		}
	}

	return result, nil
}

func (as *APIService) GetUsedLicensesPerClusterAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	usedLicenses, err := as.GetUsedLicensesPerCluster(filter)
	if err != nil {
		return nil, err
	}

	sheet := "Licenses Used Per Cluster"
	headers := []string{
		"Cluster",
		"Part Number",
		"Description",
		"Metric",
		"Hostnames",
		"Used Licenses",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range usedLicenses {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Cluster)
		sheets.SetCellValue(sheet, nextAxis(), val.LicenseTypeID)
		sheets.SetCellValue(sheet, nextAxis(), val.Description)
		sheets.SetCellValue(sheet, nextAxis(), val.Metric)
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.Hostnames, ", "))
		sheets.SetCellValue(sheet, nextAxis(), val.UsedLicenses)
	}

	return sheets, err
}

func (as *APIService) checkOlvmCappedHypervisorLicenses(licenses []dto.DatabaseUsedLicense, clustersmap map[string]dto.Cluster) error {
	olvmCapped := true

	for _, license := range licenses {
		if cluster, ok := clustersmap[license.ClusterName]; ok {
			for _, vm := range cluster.VMs {
				vmExist, err := as.Database.ExistHostdata(vm.Hostname)
				if err != nil {
					return err
				}

				if !vm.CappedCPU && vmExist {
					if !license.Ignored {
						olvmCapped = false
					}
				}
			}
		}
	}

	for i := 0; i < len(licenses); i++ {
		licenses[i].OlvmCapped = olvmCapped
	}

	return nil
}
