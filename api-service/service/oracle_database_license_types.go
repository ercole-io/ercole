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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// GetOracleDatabaseLicenseTypes return the list of OracleDatabaseLicenseType
func (as *APIService) GetOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error) {
	parts, err := as.Database.GetOracleDatabaseLicenseTypes()
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	return parts, nil
}

// GetOracleDatabaseLicenseTypesAsMap return the list of OracleDatabaseLicenseType as map by ID
func (as *APIService) GetOracleDatabaseLicenseTypesAsMap() (map[string]model.OracleDatabaseLicenseType, error) {
	parts, err := as.GetOracleDatabaseLicenseTypes()
	if err != nil {
		return nil, err
	}

	partsMap := make(map[string]model.OracleDatabaseLicenseType)
	for _, part := range parts {
		partsMap[part.ID] = part
	}

	return partsMap, nil
}

// GetOracleDatabaseLicenseType return a LicenseType by ID
func (as *APIService) GetOracleDatabaseLicenseType(id string) (*model.OracleDatabaseLicenseType, error) {
	parts, err := as.Database.GetOracleDatabaseLicenseTypes()
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	for _, part := range parts {
		if id == part.ID {
			return &part, nil
		}
	}

	return nil, utils.ErrOracleDatabaseLicenseTypeIDNotFound
}

func (as *APIService) GetOracleDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {
	agreements, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	hosts, err := as.Database.ListHostUsingOracleDatabaseLicenses()
	if err != nil {
		return nil, err
	}

	if err := as.assignOracleDatabaseAgreementsToHosts(agreements, hosts); err != nil {
		return nil, utils.NewError(err, "can't assign agreements to hosts")
	}

	getLicensesConsumedByHost := as.getterLicensesConsumedByHost()

	licenses := make(map[string]*dto.LicenseCompliance)
	for _, host := range hosts {
		license, ok := licenses[host.LicenseTypeID]
		if !ok {
			license = &dto.LicenseCompliance{
				LicenseTypeID: host.LicenseTypeID,
			}

			licenses[license.LicenseTypeID] = license
		}

		if consumedLicenses, err := getLicensesConsumedByHost(host); err != nil {
			as.Log.Error(err)
			license.Consumed += host.OriginalCount
		} else {
			license.Consumed += consumedLicenses
		}
	}

	for _, agreement := range agreements {
		license, ok := licenses[agreement.LicenseTypeID]
		if !ok {
			license = &dto.LicenseCompliance{
				LicenseTypeID: agreement.LicenseTypeID,
			}

			licenses[license.LicenseTypeID] = license
		}

		if agreement.Unlimited {
			license.Unlimited = true
		}

		for _, host := range agreement.Hosts {
			license.Covered += host.CoveredLicensesCount
		}
	}

	licenseTypes, err := as.GetOracleDatabaseLicenseTypesAsMap()
	if err != nil {
		return nil, err
	}

	result := make([]dto.LicenseCompliance, 0, len(licenses))
	for _, license := range licenses {
		if license.Unlimited || license.Consumed == 0 {
			license.Compliance = 1
		} else {
			license.Compliance = license.Covered / license.Consumed
		}

		license.ItemDescription = licenseTypes[license.LicenseTypeID].ItemDescription
		license.Metric = licenseTypes[license.LicenseTypeID].Metric

		result = append(result, *license)
	}

	return result, nil
}

func (as *APIService) getterLicensesConsumedByHost() func(host dto.HostUsingOracleDatabaseLicenses) (float64, error) {
	// map to keep history if a certain host per a certain licence as already be counted
	// by another host in its veritas cluster
	hostLicenseAlreadyCounted := make(map[string]map[string]bool)

	return func(host dto.HostUsingOracleDatabaseLicenses) (float64, error) {
		return as.getLicensesConsumedByHost(host, hostLicenseAlreadyCounted)
	}
}

func (as *APIService) getLicensesConsumedByHost(host dto.HostUsingOracleDatabaseLicenses,
	hostnamesPerLicense map[string]map[string]bool) (float64, error) {

	hostdata, err := as.Database.GetHostData(host.Name, utils.MAX_TIME)
	if err != nil {
		return 0, err
	}

	cms := hostdata.ClusterMembershipStatus
	if !cms.VeritasClusterServer ||
		(cms.VeritasClusterServer && len(cms.VeritasClusterHostnames) <= 2) {
		return host.OriginalCount, nil
	}

	_, found := hostnamesPerLicense[host.Name]
	if !found {
		hostnamesPerLicense[host.Name] = make(map[string]bool)
	}

	alreadyUsed := hostnamesPerLicense[host.Name][host.LicenseTypeID]
	if alreadyUsed {
		return 0, nil
	}

	var sumClusterCores int
	for _, h := range cms.VeritasClusterHostnames {
		_, found := hostnamesPerLicense[h]
		if !found {
			hostnamesPerLicense[h] = make(map[string]bool)
		}
		hostnamesPerLicense[h][host.LicenseTypeID] = true

		anotherHostdata, err := as.Database.GetHostData(h, utils.MAX_TIME)
		if err != nil {
			as.Log.Error(err)
			continue
		}
		sumClusterCores += anotherHostdata.Info.CPUCores
	}

	hostnamesPerLicense[host.Name][host.LicenseTypeID] = true
	consumedLicenses := float64(sumClusterCores) * 0.5 // core factor

	return consumedLicenses, nil
}
