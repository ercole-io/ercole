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

package service

import (
	"errors"
	"fmt"
	"math"

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
//TODO Create db method ad hoc
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

	getLicenseCompliance, err := as.getterNewLicenseCompliance()
	if err != nil {
		return nil, err
	}

	licenses := make(map[string]*dto.LicenseCompliance)

	// get consumptions value by hosts
	getLicensesConsumedByHost, err := as.getterLicensesConsumedByHost()
	if err != nil {
		as.Log.Error(err)
		return nil, err
	}

	for _, host := range hosts {
		license, ok := licenses[host.LicenseTypeID]
		if !ok {
			license = getLicenseCompliance(host.LicenseTypeID)
			licenses[license.LicenseTypeID] = license
		}

		consumedLicenses, err := getLicensesConsumedByHost(host)
		if err != nil {
			if errors.Is(err, utils.ErrHostNotFound) {
				as.Log.Warn(err)
			} else {
				as.Log.Error(err)
			}

			consumedLicenses += host.OriginalCount
		}

		license.Consumed += consumedLicenses
	}

	availableLicenses := make(map[string]float64)
	// get coverage values from agreements
	for _, agreement := range agreements {
		license, ok := licenses[agreement.LicenseTypeID]
		if !ok {
			license = getLicenseCompliance(agreement.LicenseTypeID)
			licenses[license.LicenseTypeID] = license
		}

		if agreement.Unlimited {
			license.Unlimited = true
		}

		if agreement.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
			availableLicenses[agreement.LicenseTypeID] += agreement.AvailableLicensesPerUser
		} else {
			availableLicenses[agreement.LicenseTypeID] += agreement.AvailableLicensesPerCore
		}

		license.Covered += agreement.CoveredLicenses
		license.Purchased += (agreement.LicensesPerCore + agreement.LicensesPerUser)
	}

	result := make([]dto.LicenseCompliance, 0, len(licenses))
	for _, license := range licenses {

		if license.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
			license.Consumed *= 25
		}

		license.Consumed = math.Round(license.Consumed)
		license.Covered = math.Round(license.Covered)
		license.Purchased = math.Round(license.Purchased)

		if license.Unlimited || license.Consumed == 0 {
			license.Compliance = 1
		} else {
			license.Compliance = license.Covered / license.Consumed
		}

		license.Available = availableLicenses[license.LicenseTypeID]

		result = append(result, *license)
	}

	return result, nil
}

func (as *APIService) getterNewLicenseCompliance() (func(licenseTypeID string) *dto.LicenseCompliance, error) {
	licenseTypes, err := as.GetOracleDatabaseLicenseTypesAsMap()
	if err != nil {
		return nil, err
	}

	getter := func(licenseTypeID string) *dto.LicenseCompliance {
		l, ok := licenseTypes[licenseTypeID]
		if !ok {
			return &dto.LicenseCompliance{
				LicenseTypeID: l.ID, //TODO
			}
		}

		return &dto.LicenseCompliance{
			LicenseTypeID:   l.ID,
			ItemDescription: l.ItemDescription,
			Metric:          l.Metric,
		}
	}

	return getter, nil
}

func (as *APIService) getterLicensesConsumedByHost() (func(host dto.HostUsingOracleDatabaseLicenses) (float64, error), error) {
	// map to keep history if a certain host per a certain licence as already be counted
	// by another host in its veritas cluster
	hostLicenseAlreadyCounted := make(map[string]map[string]bool)

	hostdatas, err := as.Database.GetHostDatas(utils.MAX_TIME)
	if err != nil {
		return nil, err
	}

	hostdatasPerHostname := make(map[string]*model.HostDataBE, len(hostdatas))
	for i := range hostdatas {
		hd := &hostdatas[i]
		hostdatasPerHostname[hd.Hostname] = hd
	}

	return func(host dto.HostUsingOracleDatabaseLicenses) (float64, error) {
		return as.getLicensesConsumedByHost(host, hostLicenseAlreadyCounted, hostdatasPerHostname)
	}, nil
}

func (as *APIService) getLicensesConsumedByHost(host dto.HostUsingOracleDatabaseLicenses,
	hostnamesPerLicense map[string]map[string]bool,
	hostdatasPerHostname map[string]*model.HostDataBE,
) (float64, error) {

	hostdata, found := hostdatasPerHostname[host.Name]
	if !found {
		return 0, fmt.Errorf("%w: %s", utils.ErrHostNotFound, host.Name)
	}

	_, found = hostnamesPerLicense[host.Name]
	if !found {
		hostnamesPerLicense[host.Name] = make(map[string]bool)
	}

	alreadyUsed := hostnamesPerLicense[host.Name][host.LicenseTypeID]
	if alreadyUsed {
		return 0, nil
	}

	consumedLicenses, err := hostdata.GetClusterCores(hostdatasPerHostname)
	if errors.Is(err, utils.ErrHostNotInCluster) {
		return host.OriginalCount, nil
	} else if err != nil {
		return 0, err
	}

	for _, h := range hostdata.ClusterMembershipStatus.VeritasClusterHostnames {
		_, found := hostnamesPerLicense[h]
		if !found {
			hostnamesPerLicense[h] = make(map[string]bool)
		}
		hostnamesPerLicense[h][host.LicenseTypeID] = true
	}

	hostnamesPerLicense[host.Name][host.LicenseTypeID] = true

	return consumedLicenses, nil
}

func (as *APIService) DeleteOracleDatabaseLicenseType(id string) error {
	return as.Database.RemoveOracleDatabaseLicenseType(id)
}

func (as *APIService) AddOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) (*model.OracleDatabaseLicenseType, error) {
	err := as.Database.InsertOracleDatabaseLicenseType(licenseType)
	if err != nil {
		return nil, err
	}

	return &licenseType, nil
}

func (as *APIService) UpdateOracleDatabaseLicenseType(licenseType model.OracleDatabaseLicenseType) (*model.OracleDatabaseLicenseType, error) {
	if err := as.Database.UpdateOracleDatabaseLicenseType(licenseType); err != nil {
		return nil, err
	}

	return &licenseType, nil
}
