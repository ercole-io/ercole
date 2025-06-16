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
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as *APIService) CreateDR(hostname string) (string, error) {
	host, err := as.Database.FindHostData(hostname)
	if err != nil {
		return "", err
	}

	dr, err := utils.DeepCopy(host)
	if err != nil {
		return "", err
	}

	dr.ID = primitive.NewObjectID()
	dr.Hostname = fmt.Sprintf("%s_DR", hostname)
	dr.IsDR = true

	if host.ClusterMembershipStatus.VeritasClusterServer {
		for i := 0; i < len(host.ClusterMembershipStatus.VeritasClusterHostnames); i++ {
			dr.ClusterMembershipStatus.VeritasClusterHostnames[i] = fmt.Sprintf("%s_DR", host.ClusterMembershipStatus.VeritasClusterHostnames[i])
		}

		totalLicenses, err := as.Database.GetClusterVeritasLicenseByHostnames(host.ClusterMembershipStatus.VeritasClusterHostnames)
		if err != nil {
			return "", err
		}

		if len(dr.Features.Oracle.Database.Databases) == 0 {
			dr.Features.Oracle.Database.Databases = append(dr.Features.Oracle.Database.Databases, model.OracleDatabase{
				Name:     "ERC999",
				Licenses: totalLicenses,
			})
		} else {
			drLicenses := make([]model.OracleDatabaseLicense, 0)
			for _, db := range dr.Features.Oracle.Database.Databases {
				for _, l := range db.Licenses {
					if !as.containsLicense(drLicenses, l.LicenseTypeID) {
						drLicenses = append(drLicenses, l)
					}
				}
			}

			diffLicenses := as.diffLicenses(totalLicenses, drLicenses)
			dr.Features.Oracle.Database.Databases[0].Licenses = append(dr.Features.Oracle.Database.Databases[0].Licenses, diffLicenses...)
		}
	}

	err = as.Database.InsertHostdata(dr)
	if err != nil {
		return "", err
	}

	return dr.Hostname, nil
}

func (as *APIService) diffLicenses(real, dr []model.OracleDatabaseLicense) []model.OracleDatabaseLicense {
	result := make([]model.OracleDatabaseLicense, 0)

	for _, r := range real {
		if !as.containsLicense(dr, r.LicenseTypeID) {
			result = append(result, r)
		}
	}

	return result
}

func (as *APIService) containsLicense(lics []model.OracleDatabaseLicense, id string) bool {
	for _, l := range lics {
		if l.LicenseTypeID == id {
			return true
		}
	}

	return false
}
