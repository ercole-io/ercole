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
	"errors"
	"strings"

	"github.com/ercole-io/ercole/v2/utils"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) SearchDatabases(filter dto.GlobalFilter) ([]dto.Database, error) {
	type getter func(filter dto.GlobalFilter) ([]dto.Database, error)
	getters := []getter{as.getOracleDatabases, as.getMySQLDatabases}

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

func (as *APIService) GetDatabasesUsedLicenses(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	type getter func(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error)
	getters := []getter{as.getOracleDatabasesUsedLicenses, as.getMySQLUsedLicenses}

	usedLicenses := make([]dto.DatabaseUsedLicense, 0)
	for _, get := range getters {
		thisDbs, err := get(filter)
		if err != nil {
			return nil, err
		}
		usedLicenses = append(usedLicenses, thisDbs...)
	}

	hostdatas, err := as.Database.GetHostDatas(utils.MAX_TIME)
	if err != nil {
		return nil, err
	}

	hostdatasPerHostname := make(map[string]*model.HostDataBE, len(hostdatas))
	for i := range hostdatas {
		hd := &hostdatas[i]
		hostdatasPerHostname[hd.Hostname] = hd
	}

	for i, l := range usedLicenses {

		hostdata, found := hostdatasPerHostname[l.Hostname]
		if usedLicenses[i].Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
			usedLicenses[i].UsedLicenses *= model.GetFactorByMetric(usedLicenses[i].Metric)
		}
		if !found {
			as.Log.Errorf("%w: %s", utils.ErrHostNotFound, l.Hostname)
			continue
		}

		if hostdata.Features.Oracle != nil && hostdata.Features.Oracle.Database != nil && hostdata.Features.Oracle.Database.Databases != nil {
			for x := range hostdata.Features.Oracle.Database.Databases {
				if hostdata.Features.Oracle.Database.Databases[x].Name == usedLicenses[i].DbName {
					for j := range hostdata.Features.Oracle.Database.Databases[x].Licenses {
						if hostdata.Features.Oracle.Database.Databases[x].Licenses[j].LicenseTypeID == usedLicenses[i].LicenseTypeID {
							usedLicenses[i].Ignored = hostdata.Features.Oracle.Database.Databases[x].Licenses[j].Ignored
							break
						}
					}
				}
			}
		}
		clusterCores, err := hostdata.GetClusterCores(hostdatasPerHostname)
		if errors.Is(err, utils.ErrHostNotInCluster) {
			continue
		} else if err != nil {
			return nil, err
		}
		consumedLicenses := float64(clusterCores) * hostdata.CoreFactor()

		usedLicenses[i].ClusterLicenses = consumedLicenses

	}

	return usedLicenses, nil
}

func (as *APIService) GetDatabasesUsedLicensesAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	licenses, err := as.GetDatabasesUsedLicenses(filter)
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
	}
	return sheets, err
}

func (as *APIService) getOracleDatabasesUsedLicenses(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	oracleLics, err := as.Database.SearchOracleDatabaseUsedLicenses("", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan)
	if err != nil {
		return nil, err
	}

	licenseTypes, err := as.GetOracleDatabaseLicenseTypesAsMap()
	if err != nil {
		return nil, err
	}

	genericLics := make([]dto.DatabaseUsedLicense, 0, len(oracleLics.Content))
	for _, o := range oracleLics.Content {
		lt := licenseTypes[o.LicenseTypeID]

		g := dto.DatabaseUsedLicense{
			Hostname:      o.Hostname,
			DbName:        o.DbName,
			LicenseTypeID: o.LicenseTypeID,
			Description:   lt.ItemDescription,
			Metric:        lt.Metric,
			UsedLicenses:  o.UsedLicenses,
		}

		genericLics = append(genericLics, g)
	}

	return genericLics, nil
}

func (as *APIService) getMySQLUsedLicenses(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicense, error) {
	mysqlLics, err := as.GetMySQLUsedLicenses(filter)
	if err != nil {
		return nil, err
	}

	genericLics := make([]dto.DatabaseUsedLicense, 0, len(mysqlLics))

	for _, lic := range mysqlLics {
		g := dto.DatabaseUsedLicense{
			Hostname:      lic.Hostname,
			DbName:        lic.InstanceName,
			LicenseTypeID: "",
			Description:   "MySQL " + lic.InstanceEdition,
			Metric:        lic.AgreementType,
			UsedLicenses:  1,
		}

		genericLics = append(genericLics, g)
	}

	return genericLics, nil
}

func (as *APIService) GetDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {

	licenses := make([]dto.LicenseCompliance, 0)

	oracle, err := as.GetOracleDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}
	licenses = append(licenses, oracle...)

	mysql, err := as.GetMySQLDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}
	licenses = append(licenses, mysql...)

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

func (as *APIService) GetDatabaseLicensesComplianceAsXLSX() (*excelize.File, error) {
	licenses, err := as.GetDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}

	sheet := "Licenses Compliance"
	headers := []string{
		"Part Number",
		"Description",
		"Metric",
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
		sheets.SetCellValue(sheet, nextAxis(), val.Consumed)
		sheets.SetCellValue(sheet, nextAxis(), val.Covered)
		sheets.SetCellValue(sheet, nextAxis(), val.Compliance)
		sheets.SetCellValue(sheet, nextAxis(), val.Unlimited)
	}
	return sheets, err
}

func (as *APIService) GetDatabasesUsedLicensesPerHostAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	usedLicenses, err := as.GetDatabasesUsedLicensesPerHost(filter)
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

func (as *APIService) GetDatabasesUsedLicensesPerHost(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicensePerHost, error) {
	licenses, err := as.GetDatabasesUsedLicenses(filter)
	if err != nil {
		return nil, err
	}

	var licensesPerHost []dto.DatabaseUsedLicensePerHost

licenses:
	for _, v := range licenses {
		for i, v2 := range licensesPerHost {
			if v.Hostname == v2.Hostname && v.LicenseTypeID == v2.LicenseTypeID {
				licensesPerHost[i].DatabaseNames = append(licensesPerHost[i].DatabaseNames, v.DbName)
				continue licenses
			}
		}

		licensesPerHost = append(licensesPerHost, dto.DatabaseUsedLicensePerHost{
			Hostname:        v.Hostname,
			DatabaseNames:   []string{v.DbName},
			LicenseTypeID:   v.LicenseTypeID,
			Description:     v.Description,
			Metric:          v.Metric,
			UsedLicenses:    v.UsedLicenses,
			ClusterLicenses: v.ClusterLicenses,
		})
	}

	return licensesPerHost, nil
}

func (as *APIService) GetDatabasesUsedLicensesPerCluster(filter dto.GlobalFilter) ([]dto.DatabaseUsedLicensePerCluster, error) {
	licenses, err := as.GetDatabasesUsedLicenses(filter)
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
				UsedLicenses:  float64(c.CPU) * 0.5,
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

func (as *APIService) GetDatabasesUsedLicensesPerClusterAsXLSX(filter dto.GlobalFilter) (*excelize.File, error) {
	usedLicenses, err := as.GetDatabasesUsedLicensesPerCluster(filter)
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
