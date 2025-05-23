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
	"strings"

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

func (as *APIService) GetOracleDatabaseLicensesCompliance(locations []string) ([]dto.LicenseCompliance, error) {
	filter := dto.NewGetOracleDatabaseContractsFilter()
	filter.Locations = locations

	contracts, err := as.Database.ListOracleDatabaseContracts(filter)
	if err != nil {
		return nil, err
	}

	usages, err := as.getLicensesUsage(filter.Locations)
	if err != nil {
		return nil, err
	}

	if err := as.assignOracleDatabaseContractsToHosts(contracts, usages); err != nil {
		return nil, utils.NewError(err, "can't assign contracts to hosts")
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

		if as.Config.APIService.RoundLicenses {
			license.Consumed = math.Round(license.Consumed)
		}
	}

	availableLicenses := make(map[string]float64)
	// get coverage values from contracts
	for _, contract := range contracts {
		license, ok := licenses[contract.LicenseTypeID]
		if !ok {
			license = getLicenseCompliance(contract.LicenseTypeID)
			licenses[license.LicenseTypeID] = license
		}

		if contract.Unlimited {
			license.Unlimited = true
		}

		if contract.Metric == model.LicenseTypeMetricNamedUserPlusPerpetual {
			availableLicenses[contract.LicenseTypeID] += contract.AvailableLicensesPerUser
		} else {
			availableLicenses[contract.LicenseTypeID] += contract.AvailableLicensesPerCore
		}

		if as.Config.APIService.RoundLicenses {
			license.Covered += math.Round(contract.CoveredLicenses)
			license.Purchased += math.Round((contract.LicensesPerCore + contract.LicensesPerUser))
		} else {
			license.Covered += contract.CoveredLicenses
			license.Purchased += (contract.LicensesPerCore + contract.LicensesPerUser)
		}
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

		license.Available = math.Floor(availableLicenses[license.LicenseTypeID])

		result = append(result, *license)
	}

	drLicenses, err := as.getVeritasClusterLicensesDR()
	if err != nil {
		return nil, err
	}

	if len(drLicenses) > 0 {
		result = as.mergeDrLicensesCompliances(result, drLicenses)
	}

	return result, nil
}

func (as *APIService) getVeritasClusterLicensesDR() (map[string]dto.LicenseCompliance, error) {
	licenses, err := as.GetClusterVeritasLicenses(dto.GlobalFilter{OlderThan: utils.MAX_TIME})
	if err != nil {
		return nil, err
	}

	drLicenses := make([]dto.ClusterVeritasLicense, 0)

	for _, l := range licenses {
		if strings.Contains(l.ID, "_DR") {
			var licenseExist bool

			drHostnames := strings.Split(l.ID, "-")

			for _, drHostname := range drHostnames {
				var realHost string

				exists, err := as.Database.ExistHostdata(drHostname)
				if err != nil {
					return nil, err
				}

				realHost = drHostname

				if !exists {
					hostname := strings.Replace(drHostname, "_DR", "", 1)

					existsDR, err := as.Database.ExistHostdata(hostname)
					if err != nil {
						return nil, err
					}

					if existsDR {
						realHost = hostname
					}
				}

				isIgnored, err := as.Database.IsOracleLicenseIgnored(realHost, l.LicenseTypeID)
				if err != nil {
					return nil, err
				}

				licenseExists, err := as.Database.OracleLicenseExists(realHost, l.LicenseTypeID)
				if err != nil {
					return nil, err
				}

				if licenseExist {
					continue
				}

				licenseExist = !isIgnored && licenseExists
			}

			if licenseExist {
				drLicenses = append(drLicenses, l)
			}
		}
	}

	grouped := make(map[string]dto.LicenseCompliance)

	for _, drLicense := range drLicenses {
		grouped[drLicense.LicenseTypeID] = dto.LicenseCompliance{
			LicenseTypeID:   drLicense.LicenseTypeID,
			ItemDescription: drLicense.Description,
			Metric:          drLicense.Metric,
			Consumed:        grouped[drLicense.LicenseTypeID].Consumed + drLicense.Count,
		}
	}

	return grouped, nil
}

func (as *APIService) mergeDrLicensesCompliances(compliances []dto.LicenseCompliance, drlicenses map[string]dto.LicenseCompliance) []dto.LicenseCompliance {
	res := make([]dto.LicenseCompliance, 0, len(drlicenses))
	unmatched := make([]dto.LicenseCompliance, 0)

	for _, c := range compliances {
		if v, ok := drlicenses[c.LicenseTypeID]; ok {
			c.Consumed = math.Round(v.Consumed + c.Consumed)

			if c.Unlimited || c.Consumed == 0 || c.Cost == 0 {
				c.Compliance = 1
			} else {
				c.Compliance = c.Covered / c.Consumed
			}
		} else {
			unmatched = append(unmatched, v)
		}

		res = append(res, c)
	}

	res = append(res, unmatched...)

	return res
}

func (as *APIService) getLicensesUsage(locations []string) ([]dto.HostUsingOracleDatabaseLicenses, error) {
	filter := dto.GlobalFilter{
		Location:    strings.Join(locations, ","),
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	usedLicenses, err := as.getOracleDatabasesUsedLicenses("", filter)
	if err != nil {
		return nil, err
	}

	usages := make([]dto.HostUsingOracleDatabaseLicenses, 0, len(usedLicenses))
	hostnamesPerLicense := make(map[string]map[string]bool)

	hostdatas, err := as.Database.GetHostDatas(dto.GlobalFilter{
		OlderThan: utils.MAX_TIME,
		Location:  strings.Join(locations, ","),
	})
	if err != nil {
		return nil, err
	}

	clusters, err := as.Database.GetClusters(dto.GlobalFilter{
		Location:    strings.Join(locations, ","),
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	})
	if err != nil {
		return nil, err
	}

	hostdatasMap := make(map[string]model.HostDataBE, len(hostdatas))
	for _, hostdata := range hostdatas {
		hostdatasMap[hostdata.Hostname] = hostdata
	}

	clustersMap := make(map[string]dto.Cluster, len(clusters))
	for _, cluster := range clusters {
		clustersMap[cluster.Name] = cluster
	}

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

			isCapped, err := as.manageLicenseWithCappedCPU(usedLicense, clustersMap, hostdatasMap)
			if err != nil {
				return nil, err
			}

			host, ok := hostdatasMap[usedLicense.Hostname]
			if !ok {
				continue
			}

			if isCapped &&
				host.Features.Oracle != nil &&
				host.Features.Oracle.Database != nil &&
				host.Features.Oracle.Database.Databases != nil {
				databases := host.Features.Oracle.Database.Databases

				for _, database := range databases {
					for _, license := range database.Licenses {
						if license.LicenseTypeID == usedLicense.LicenseTypeID &&
							database.Name == usedLicense.DbName {
							if database.Edition() == model.OracleDatabaseEditionStandard {
								licensesCount = usedLicense.ClusterLicenses
							} else {
								licensesCount = usedLicense.UsedLicenses
							}
						}
					}
				}

				typeClusterHost = "host"
				name = usedLicense.Hostname
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

func (as *APIService) manageLicenseWithCappedCPU(usedLicense dto.DatabaseUsedLicense, clustersMap map[string]dto.Cluster, hostadatasMap map[string]model.HostDataBE) (bool, error) {
	if usedLicense.ClusterType != "VeritasCluster" && usedLicense.ClusterName != "" {
		cluster, ok := clustersMap[usedLicense.ClusterName]
		if !ok {
			return false, utils.ErrClusterNotFound
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
				host, ok := hostadatasMap[vm]
				if !ok {
					continue
				}

				if host.Features.Oracle != nil && host.Features.Oracle.Database != nil && host.Features.Oracle.Database.Databases != nil {
					databases := host.Features.Oracle.Database.Databases

					for _, database := range databases {
						for _, license := range database.Licenses {
							if !license.Ignored {
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
		}

		if !notlicenseCapped && licenseCapped {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, nil
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
				LicenseTypeID: l.ID,
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
