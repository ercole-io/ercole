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
package dto

import (
	"time"

	"github.com/ercole-io/ercole/v2/model"
)

type CreateScenarioRequest struct {
	Name  string                      `json:"name"`
	Hosts []CreateHostScenarioRequest `json:"hosts"`
}

type CreateHostScenarioRequest struct {
	Hostname string `json:"hostname"`
	Core     int    `json:"core"`
}

type ScenarioResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"createdAt"`
	Hosts     []SimulatedHost `json:"hosts"`
}

type SimulatedHost struct {
	ID            string `json:"id"`
	Hostname      string `json:"hostname"`
	OriginalCore  int    `json:"originalCore"`
	SimulatedCore int    `json:"simulatedCore"`
}

func ToSimulatedHost(m model.SimulatedHost) SimulatedHost {
	return SimulatedHost{
		ID:            m.ID.Hex(),
		Hostname:      m.Host.Hostname,
		OriginalCore:  m.Host.Info.CPUCores,
		SimulatedCore: m.Core,
	}
}

func ToScenarioResponse(m model.Scenario) ScenarioResponse {
	simulatedHosts := make([]SimulatedHost, 0)
	for _, host := range m.Hosts {
		simulatedHosts = append(simulatedHosts, ToSimulatedHost(host))
	}

	return ScenarioResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Hosts:     simulatedHosts,
	}
}

func ToScenariosResponse(scenarios []model.Scenario) []ScenarioResponse {
	res := make([]ScenarioResponse, 0)
	for _, s := range scenarios {
		res = append(res, ToScenarioResponse(s))
	}

	return res
}

type ScenarioLicenseComplianceResponse struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	CreatedAt time.Time                 `json:"createdAt"`
	Actual    []model.LicenseCompliance `json:"actual"`
	Got       []model.LicenseCompliance `json:"got"`
}

func ToScenarioLicenseComplianceResponse(m model.Scenario) ScenarioLicenseComplianceResponse {
	return ScenarioLicenseComplianceResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Actual:    m.LicenseCompliance.Actual,
		Got:       m.LicenseCompliance.Got,
	}
}

type ScenarioLicenseUsedDatabaseResponse struct {
	ID        string                      `json:"id"`
	Name      string                      `json:"name"`
	CreatedAt time.Time                   `json:"createdAt"`
	Actual    []model.LicenseUsedDatabase `json:"actual"`
	Got       []model.LicenseUsedDatabase `json:"got"`
}

func ToScenarioLicenseUsedDatabaseResponse(m model.Scenario) ScenarioLicenseUsedDatabaseResponse {
	return ScenarioLicenseUsedDatabaseResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Actual:    m.LicenseUsed.LicenseDatabase.Actual,
		Got:       m.LicenseUsed.LicenseDatabase.Got,
	}
}

type ScenarioLicenseUsedHostResponse struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	CreatedAt time.Time               `json:"createdAt"`
	Actual    []model.LicenseUsedHost `json:"actual"`
	Got       []model.LicenseUsedHost `json:"got"`
}

func ToScenarioLicenseUsedHostResponse(m model.Scenario) ScenarioLicenseUsedHostResponse {
	return ScenarioLicenseUsedHostResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Actual:    m.LicenseUsed.LicenseHost.Actual,
		Got:       m.LicenseUsed.LicenseHost.Got,
	}
}

type ScenarioLicenseUsedClusterResponse struct {
	ID        string                     `json:"id"`
	Name      string                     `json:"name"`
	CreatedAt time.Time                  `json:"createdAt"`
	Actual    []model.LicenseUsedCluster `json:"actual"`
	Got       []model.LicenseUsedCluster `json:"got"`
}

func ToScenarioLicenseUsedClusterResponse(m model.Scenario) ScenarioLicenseUsedClusterResponse {
	return ScenarioLicenseUsedClusterResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Actual:    m.LicenseUsed.LicenseHypervisorCluster.Actual,
		Got:       m.LicenseUsed.LicenseHypervisorCluster.Got,
	}
}

type ScenarioLicenseUsedClusterVeritasResponse struct {
	ID        string                            `json:"id"`
	Name      string                            `json:"name"`
	CreatedAt time.Time                         `json:"createdAt"`
	Actual    []model.LicenseUsedClusterVeritas `json:"actual"`
	Got       []model.LicenseUsedClusterVeritas `json:"got"`
}

func ToScenarioLicenseUsedClusterVeritasResponse(m model.Scenario) ScenarioLicenseUsedClusterVeritasResponse {
	return ScenarioLicenseUsedClusterVeritasResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Actual:    m.LicenseUsed.LicenseClusterVeritas.Actual,
		Got:       m.LicenseUsed.LicenseClusterVeritas.Got,
	}
}
