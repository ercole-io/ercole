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
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (hds *HostDataService) mySqlDatabasesChecks(previousHostdata, hostdata *model.HostDataBE) {
	if hostdata.Features.MySQL.Instances == nil {
		return
	}

	licenseTypes, err := hds.getMySqlDatabaseLicenseTypes()
	if err != nil {
		hds.Log.Error(err)

		licenseTypes = make([]model.MySqlLicenseType, 0)
	}

	hds.setMySqlLicenseTypes(hostdata, licenseTypes)

	hds.ignoreMySqlPreviousLicences(previousHostdata, hostdata)
}

func (hds *HostDataService) getMySqlDatabaseLicenseTypes() ([]model.MySqlLicenseType, error) {
	licenseTypes, err := hds.ApiSvcClient.GetMySqlDatabaseLicenseTypes()
	if err != nil {
		return nil, utils.NewError(err, "Can't retrieve licenseTypes")
	}

	return licenseTypes, nil
}

func (hds *HostDataService) setMySqlLicenseTypes(hostdata *model.HostDataBE, licenseTypes []model.MySqlLicenseType) {
	for i, instance := range hostdata.Features.MySQL.Instances {
		license := &hostdata.Features.MySQL.Instances[i].License

		for _, licenseType := range licenseTypes {
			if instance.Edition == model.MySQLEditionEnterprise {
				license.Count = 1
				license.LicenseTypeID = licenseType.ID
				license.Name = licenseType.ItemDescription
			}
		}

		hostdata.Features.MySQL.Instances[i].License = *license
	}
}

func (hds *HostDataService) ignoreMySqlPreviousLicences(previous, new *model.HostDataBE) {
	if previous == nil || previous.Features.MySQL == nil {
		return
	}

	type ignoredLicense struct {
		licenseTypeID string
		ignored       bool
		comment       string
	}

	ignoredDbLicenses := make(map[string][]ignoredLicense)

	for _, db := range previous.Features.MySQL.Instances {
		licenses := make([]ignoredLicense, 0)

		if db.License.Ignored {
			ignored := ignoredLicense{ignored: true, licenseTypeID: db.License.LicenseTypeID, comment: db.License.IgnoredComment}
			licenses = append(licenses, ignored)
		}

		if len(licenses) > 0 {
			ignoredDbLicenses[db.UUID] = licenses
		}
	}

	for i, db := range new.Features.MySQL.Instances {
		if ignoredDbLicense, ok := ignoredDbLicenses[db.UUID]; ok {
			for _, v := range ignoredDbLicense {
				if db.License.LicenseTypeID == v.licenseTypeID {
					new.Features.MySQL.Instances[i].License.Ignored = v.ignored
					new.Features.MySQL.Instances[i].License.IgnoredComment = v.comment
				}
			}
		}
	}
}
