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
func (as *APIService) GetOracleDatabaseLicenseType(id string) (*model.OracleDatabaseLicenseType, error) {
	out, err := as.Database.GetOracleDatabaseLicenseType(id)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	if out == nil {
		return nil, utils.ErrOracleDatabaseLicenseTypeIDNotFound
	} else {
		return out, nil
	}
}

func (as *APIService) GetOracleDatabaseLicensesCompliance() ([]dto.LicenseCompliance, error) {
	agreements, err := as.Database.ListOracleDatabaseAgreements()
	if err != nil {
		return nil, err
	}

	usages, err := as.getLicensesUsage()
	if err != nil {
		return nil, err
	}

	if err := as.assignOracleDatabaseAgreementsToHosts(agreements, usages); err != nil {
		return nil, utils.NewError(err, "can't assign agreements to hosts")
	}

	getLicenseCompliance, err := as.getterNewLicenseCompliance()
	if err != nil {
		return nil, err
	}

	licenses := make(map[string]*dto.LicenseCompliance)

	for _, usage := range usages {
		license, ok := licenses[usage.LicenseTypeID]

		if !ok {
			license = getLicenseCompliance(usage.LicenseTypeID)
			licenses[license.LicenseTypeID] = license
		}

		license.Consumed += usage.OriginalCount
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
		license.Consumed = math.Round(license.Consumed)
		license.Covered = math.Round(license.Covered)
		license.Purchased = math.Round(license.Purchased)

		if license.Unlimited || license.Consumed == 0 || license.Cost == 0 {
			license.Compliance = 1
		} else {
			license.Compliance = license.Covered / license.Consumed
		}

		license.Available = availableLicenses[license.LicenseTypeID]

		result = append(result, *license)
	}

	return result, nil
}

func (as *APIService) getLicensesUsage() ([]dto.HostUsingOracleDatabaseLicenses, error) {
	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses, err := as.getOracleDatabasesUsedLicenses("", filter)
	if err != nil {
		return nil, err
	}

	usages := make([]dto.HostUsingOracleDatabaseLicenses, 0, len(usedLicenses))
	hostnamesPerLicense := make(map[string]map[string]bool)

	for _, usedLicense := range usedLicenses {
		if usedLicense.Ignored {
			continue
		}

		var typeClusterHost, name string

		var licensesCount float64

		if usedLicense.ClusterName != "" {
			typeClusterHost = "cluster"
			name = usedLicense.ClusterName
			licensesCount = usedLicense.ClusterLicenses

			if usedLicense.ClusterType != "VeritasCluster" {
				cluster, err := as.GetCluster(usedLicense.ClusterName, utils.MAX_TIME)
				if err != nil {
					return nil, err
				}

				var capped, licenseCapped, notlicenseCapped bool

				vms := make(map[string]bool)

				for _, hostNameVM := range cluster.VMs {
					if !hostNameVM.CappedCPU {
						if _, ok := vms[hostNameVM.Hostname]; !ok {
							vms[hostNameVM.Hostname] = false
						}

						continue
					} else {
						if _, ok := vms[hostNameVM.Hostname]; !ok {
							vms[hostNameVM.Hostname] = true
						}
						capped = true
					}
				}

				if capped {
					for vm, cap := range vms {
						host, err := as.GetHost(vm, utils.MAX_TIME, false)
						if err != nil {
							continue
						}

						if host != nil && host.Features.Oracle != nil && host.Features.Oracle.Database != nil && host.Features.Oracle.Database.Databases != nil {
							databases := host.Features.Oracle.Database.Databases

							for _, database := range databases {
								for _, license := range database.Licenses {
									if license.LicenseTypeID == usedLicense.LicenseTypeID && cap {
										licenseCapped = true
									} else if license.LicenseTypeID == usedLicense.LicenseTypeID && !cap {
										notlicenseCapped = true
									}
								}
							}
						}
					}
				}

				if !notlicenseCapped && licenseCapped {
					typeClusterHost = "host"
					name = usedLicense.Hostname
					licensesCount = usedLicense.UsedLicenses
				}
			}

			_, found := hostnamesPerLicense[name]
			if !found {
				hostnamesPerLicense[name] = make(map[string]bool)
			}

			alreadyUsed := hostnamesPerLicense[name][usedLicense.LicenseTypeID]
			if alreadyUsed {
				continue
			}

			hostnamesPerLicense[name][usedLicense.LicenseTypeID] = true
		} else {
			typeClusterHost = "host"
			name = usedLicense.Hostname
			licensesCount = usedLicense.UsedLicenses

			_, found := hostnamesPerLicense[name]
			if !found {
				hostnamesPerLicense[name] = make(map[string]bool)
			}

			alreadyUsed := hostnamesPerLicense[name][usedLicense.LicenseTypeID]
			if alreadyUsed {
				continue
			}

			hostnamesPerLicense[name][usedLicense.LicenseTypeID] = true
		}

		g := dto.HostUsingOracleDatabaseLicenses{
			LicenseTypeID: usedLicense.LicenseTypeID,
			Name:          name,
			Type:          typeClusterHost,
			LicenseCount:  licensesCount,
			OriginalCount: licensesCount,
		}

		usages = append(usages, g)
	}

	return usages, nil
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
			Cost:            l.Cost,
		}
	}

	return getter, nil
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
