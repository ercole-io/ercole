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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package service

import (
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as *APIService) CreateScenario(req dto.CreateScenarioRequest, locations []string, filter dto.GlobalFilter) (*model.Scenario, error) {
	scenario := &model.Scenario{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	for _, h := range req.Hosts {
		host, err := as.Database.FindHostData(h.Hostname)
		if err != nil {
			return nil, err
		}

		scenario.Hosts = append(scenario.Hosts, model.SimulatedHost{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			Host:      host,
			Core:      h.Core,
		})
	}

	actualCompliance, err := as.getLicenseCompliance(locations)
	if err != nil {
		return nil, err
	}

	scenario.LicenseCompliance.Actual = actualCompliance

	actualUsedPerDatabase, err := as.getLicenseUsedPerDatabase("", filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseDatabase.Actual = actualUsedPerDatabase

	actualUsedPerHost, err := as.getLicenseUsedPerHost(filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseHost.Actual = actualUsedPerHost

	actualUsedCluster, err := as.getLicenseUsedPerCluster(filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseHypervisorCluster.Actual = actualUsedCluster

	actualUsedClusterVeritas, err := as.getLicenseUsedPerClusterVeritas(filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseClusterVeritas.Actual = actualUsedClusterVeritas

	err = as.Database.CreateSimulatedHosts(scenario.Hosts...)
	if err != nil {
		return nil, err
	}

	for _, simulatedHost := range scenario.Hosts {
		if err := as.UpdateHostLicenseCount(simulatedHost.Host.Hostname, simulatedHost.Core); err != nil {
			return nil, err
		}
	}

	gotCompliance, err := as.getLicenseCompliance(locations)
	if err != nil {
		return nil, err
	}

	scenario.LicenseCompliance.Got = gotCompliance

	gotUsedPerDatabase, err := as.getLicenseUsedPerDatabase("", filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseDatabase.Got = gotUsedPerDatabase

	gotUsedPerHost, err := as.getLicenseUsedPerHost(filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseHost.Got = gotUsedPerHost

	gotUsedCluster, err := as.getLicenseUsedPerCluster(filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseHypervisorCluster.Got = gotUsedCluster

	gotUsedClusterVeritas, err := as.getLicenseUsedPerClusterVeritas(filter)
	if err != nil {
		return nil, err
	}

	scenario.LicenseUsed.LicenseClusterVeritas.Got = gotUsedClusterVeritas

	res, err := as.Database.CreateScenario(scenario)
	if err != nil {
		return nil, err
	}

	for _, simulatedHost := range scenario.Hosts {
		if err := as.UpdateHostLicenseCount(simulatedHost.Host.Hostname, simulatedHost.Host.Info.CPUCores); err != nil {
			return nil, err
		}

		err = as.Database.RemoveSimulatedHost(simulatedHost.ID)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (as *APIService) getLicenseCompliance(locations []string) ([]model.LicenseCompliance, error) {
	licensesCompliance, err := as.GetDatabaseLicensesCompliance(locations)
	if err != nil {
		return nil, err
	}

	res := make([]model.LicenseCompliance, 0, len(licensesCompliance))
	for _, lc := range licensesCompliance {
		res = append(res, lc.ToModel())
	}

	return res, nil
}

func (as *APIService) getLicenseUsedPerDatabase(hostname string, f dto.GlobalFilter) ([]model.LicenseUsedDatabase, error) {
	licenseUsedDatabase, err := as.GetUsedLicensesPerDatabases(hostname, f)
	if err != nil {
		return nil, err
	}

	res := make([]model.LicenseUsedDatabase, 0, len(licenseUsedDatabase))
	for _, lud := range licenseUsedDatabase {
		res = append(res, lud.ToModel())
	}

	return res, nil
}

func (as *APIService) getLicenseUsedPerHost(f dto.GlobalFilter) ([]model.LicenseUsedHost, error) {
	licenseUsedHost, err := as.GetUsedLicensesPerHost(f)
	if err != nil {
		return nil, err
	}

	res := make([]model.LicenseUsedHost, 0, len(licenseUsedHost))
	for _, l := range licenseUsedHost {
		res = append(res, l.ToModel())
	}

	return res, nil
}

func (as *APIService) getLicenseUsedPerCluster(f dto.GlobalFilter) ([]model.LicenseUsedCluster, error) {
	licenseUsedCluster, err := as.GetUsedLicensesPerCluster(f)
	if err != nil {
		return nil, err
	}

	res := make([]model.LicenseUsedCluster, 0, len(licenseUsedCluster))
	for _, l := range licenseUsedCluster {
		res = append(res, l.ToModel())
	}

	return res, nil
}

func (as *APIService) getLicenseUsedPerClusterVeritas(f dto.GlobalFilter) ([]model.LicenseUsedClusterVeritas, error) {
	licenseUsedClusterVeritas, err := as.GetClusterVeritasLicenses(f)
	if err != nil {
		return nil, err
	}

	res := make([]model.LicenseUsedClusterVeritas, 0, len(licenseUsedClusterVeritas))
	for _, l := range licenseUsedClusterVeritas {
		res = append(res, l.ToModel())
	}

	return res, nil
}

func (as *APIService) GetScenarios() ([]model.Scenario, error) {
	return as.Database.GetScenarios()
}

func (as *APIService) GetScenario(id primitive.ObjectID) (*model.Scenario, error) {
	return as.Database.GetScenario(id)
}

func (as *APIService) RemoveScenario(id primitive.ObjectID) error {
	return as.Database.RemoveScenario(id)
}

func (as *APIService) UpdateHostLicenseCount(hostname string, core int) error {
	licenseCount := core / 2

	err := as.Database.UpdateHostCores(hostname, core)
	if err != nil {
		return err
	}

	err = as.Database.UpdateLicenseCount(hostname, licenseCount)
	if err != nil {
		return err
	}

	return nil
}
