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

package service

import (
	"strings"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (hds *HostDataService) sqlServerDatabasesChecks(previousHostdata, hostdata *model.HostDataBE) {
	if hostdata.Features.Microsoft.SQLServer == nil || hostdata.Features.Microsoft.SQLServer.Instances == nil {
		return
	}

	licenseTypes, err := hds.getSqlServerDatabaseLicenseTypes()
	if err != nil {
		hds.Log.Error(err)

		licenseTypes = make([]model.SqlServerDatabaseLicenseType, 0)
	}

	hostdataWithVersion := hds.setSqlServerVersion(hostdata)

	hds.setSqlServerLicenseTypes(hostdataWithVersion, licenseTypes)

	hds.ignoreSqlServerPreviousLicences(previousHostdata, hostdataWithVersion)
}

func (hds *HostDataService) getSqlServerDatabaseLicenseTypes() ([]model.SqlServerDatabaseLicenseType, error) {
	licenseTypes, err := hds.ApiSvcClient.GetSQLServerDatabaseLicenseTypes()
	if err != nil {
		return nil, utils.NewError(err, "Can't retrieve licenseTypes")
	}

	return licenseTypes, nil
}

func (hds *HostDataService) setSqlServerVersion(hostdata *model.HostDataBE) *model.HostDataBE {
	versionsMap := map[string]string{
		"15.0":  "2019",
		"14.0":  "2017",
		"13.0":  "2016",
		"12.0":  "2014",
		"11.0":  "2012",
		"10.50": "2008 R2",
		"10.0":  "2008",
		"9.0":   "2005",
		"8.0":   "2000",
		"7.0":   "7.0",
		"6.50":  "6.5",
		"6.00":  "6.0",
	}

	for i, instance := range hostdata.Features.Microsoft.SQLServer.Instances {
		hostdata.Features.Microsoft.SQLServer.Instances[i].Version = versionsMap[getMajorRel(instance.Version)]
	}

	return hostdata
}

func (hds *HostDataService) setSqlServerLicenseTypes(hostdata *model.HostDataBE, licenseTypes []model.SqlServerDatabaseLicenseType) {
	var licenseCount float64

	var listEnt = make(map[string]string)

	info := hostdata.Info

	if info.HardwareAbstraction == "PH" {
		result := info.CPUCores / info.CPUSockets
		if result < 4 {
			licenseCount = float64(info.CPUSockets * 4)
		} else {
			licenseCount = float64(info.CPUCores)
		}
	} else if info.HardwareAbstraction == "VIRT" {
		if info.CPUThreads < 4 {
			licenseCount = 4
		} else {
			licenseCount = float64(info.CPUThreads)
		}
	} else {
		return
	}

	for i, instance := range hostdata.Features.Microsoft.SQLServer.Instances {
		license := &hostdata.Features.Microsoft.SQLServer.Instances[i].License

		for _, licenseType := range licenseTypes {
			if licenseType.Edition == instance.Edition && licenseType.Version == instance.Version {
				license.LicenseTypeID = licenseType.ID
				license.Name = licenseType.ItemDescription

				if instance.Edition == "ENT" {
					listEnt[licenseType.ID] = licenseType.ItemDescription
				}

				if instance.Edition != "ENT" && instance.Edition != "STD" {
					license.Count = 0
				} else {
					license.Count = licenseCount
				}
			}
		}

		hostdata.Features.Microsoft.SQLServer.Instances[i].License = *license
	}

	if len(listEnt) > 0 {
		licenseTypeId, name := findMaxLicenseId(listEnt)

		for i, instance := range hostdata.Features.Microsoft.SQLServer.Instances {
			if instance.Edition == "STD" {
				instance.Edition = "ENT"
				instance.License.LicenseTypeID = licenseTypeId
				instance.License.Name = name
				hostdata.Features.Microsoft.SQLServer.Instances[i] = instance
			}
		}
	}
}

func (hds *HostDataService) ignoreSqlServerPreviousLicences(previous, new *model.HostDataBE) {
	if previous == nil || previous.Features.Microsoft == nil ||
		previous.Features.Microsoft.SQLServer == nil {
		return
	}

	ignoredDbLicenses := make(map[int][]string)

	for _, db := range previous.Features.Microsoft.SQLServer.Instances {
		licenses := make([]string, 0)

		if db.License.Ignored {
			licenses = append(licenses, db.License.LicenseTypeID)
		}

		if len(licenses) > 0 {
			ignoredDbLicenses[db.DatabaseID] = licenses
		}
	}

	for i, db := range new.Features.Microsoft.SQLServer.Instances {
		if licenseTypeID, ok := ignoredDbLicenses[db.DatabaseID]; ok {
			new.Features.Microsoft.SQLServer.Instances[i].License.Ignored = utils.Contains(licenseTypeID, db.License.LicenseTypeID)
		}
	}
}

func findMaxLicenseId(lics map[string]string) (string, string) {
	var maxLicId, maxName string

	for licId, name := range lics {
		if licId > maxLicId {
			maxLicId = licId
			maxName = name
		}
	}

	return maxLicId, maxName
}

func getMajorRel(version string) string {
	var val1, val2 = "", ""

	res := strings.Split(version, ".")

	for i, v := range res {
		if i == 0 {
			val1 = v
		} else if i == 1 {
			val2 = v
		} else {
			break
		}
	}

	return val1 + "." + val2
}
