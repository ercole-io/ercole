// Copyright (c) 2025 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (hds *HostDataService) createDR(hostdata model.HostDataBE) error {
	drname := fmt.Sprintf("%s_DR", hostdata.Hostname)

	if !hds.Database.ExistsDR(drname) {
		return nil
	}

	if err := hds.Database.DismissHost(drname); err != nil {
		return err
	}

	hostdata.ID = primitive.NewObjectID()
	hostdata.Hostname = drname
	hostdata.IsDR = true

	if hostdata.ClusterMembershipStatus.VeritasClusterServer {
		totalLicenses, err := hds.Database.GetClusterVeritasLicenseByHostnames(hostdata.ClusterMembershipStatus.VeritasClusterHostnames)
		if err != nil {
			return err
		}

		for i := 0; i < len(hostdata.ClusterMembershipStatus.VeritasClusterHostnames); i++ {
			hostdata.ClusterMembershipStatus.VeritasClusterHostnames[i] = fmt.Sprintf("%s_DR", hostdata.ClusterMembershipStatus.VeritasClusterHostnames[i])
		}

		if len(hostdata.Features.Oracle.Database.Databases) == 0 {
			hostdata.Features.Oracle.Database.Databases = append(hostdata.Features.Oracle.Database.Databases, model.OracleDatabase{
				Name:     "ERC999",
				Licenses: totalLicenses,
			})
		} else {
			drLicenses := make([]model.OracleDatabaseLicense, 0)
			for _, db := range hostdata.Features.Oracle.Database.Databases {
				for _, l := range db.Licenses {
					if !hds.containsLicense(drLicenses, l.LicenseTypeID) {
						drLicenses = append(drLicenses, l)
					}
				}
			}

			diffLicenses := hds.diffLicenses(totalLicenses, drLicenses)
			hostdata.Features.Oracle.Database.Databases[0].Licenses = append(hostdata.Features.Oracle.Database.Databases[0].Licenses, diffLicenses...)
		}
	}

	return hds.Database.InsertHostData(hostdata)
}

func (hds *HostDataService) diffLicenses(real, dr []model.OracleDatabaseLicense) []model.OracleDatabaseLicense {
	result := make([]model.OracleDatabaseLicense, 0)

	for _, r := range real {
		if !hds.containsLicense(dr, r.LicenseTypeID) {
			result = append(result, r)
		}
	}

	return result
}

func (hds *HostDataService) containsLicense(lics []model.OracleDatabaseLicense, id string) bool {
	for _, l := range lics {
		if l.LicenseTypeID == id {
			return true
		}
	}

	return false
}
