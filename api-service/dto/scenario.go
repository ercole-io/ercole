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
	ID        string                        `json:"id"`
	Name      string                        `json:"name"`
	CreatedAt time.Time                     `json:"createdAt"`
	Licenses  []LicenseComplianceSimulation `json:"licenses"`
}

type LicenseComplianceSimulation struct {
	Actual *model.LicenseCompliance `json:"actual"`
	Got    *model.LicenseCompliance `json:"got"`
}

func ToScenarioLicenseComplianceResponse(m model.Scenario) ScenarioLicenseComplianceResponse {
	licenses := make([]LicenseComplianceSimulation, 0)

	for _, g := range m.LicenseCompliance.Got {
		got := g

		for _, a := range m.LicenseCompliance.Actual {
			actual := a

			if a.LicenseTypeID == g.LicenseTypeID {
				licenses = append(licenses, LicenseComplianceSimulation{
					Actual: &actual,
					Got:    &got,
				})
			}
		}
	}

	return ScenarioLicenseComplianceResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Licenses:  licenses,
	}
}

type ScenarioLicenseUsedDatabaseResponse struct {
	ID        string                          `json:"id"`
	Name      string                          `json:"name"`
	CreatedAt time.Time                       `json:"createdAt"`
	Licenses  []LicenseUsedDatabaseSimulation `json:"licenses"`
}

type LicenseUsedDatabaseSimulation struct {
	Actual *model.LicenseUsedDatabase `json:"actual"`
	Got    *model.LicenseUsedDatabase `json:"got"`
}

func ToScenarioLicenseUsedDatabaseResponse(m model.Scenario) ScenarioLicenseUsedDatabaseResponse {
	licenses := make([]LicenseUsedDatabaseSimulation, 0)

	for _, g := range m.LicenseUsed.LicenseDatabase.Got {
		got := g

		for _, a := range m.LicenseUsed.LicenseDatabase.Actual {
			actual := a

			if a.Hostname == g.Hostname && a.DbName == g.DbName && a.LicenseTypeID == g.LicenseTypeID {
				licenses = append(licenses, LicenseUsedDatabaseSimulation{
					Actual: &actual,
					Got:    &got,
				})
			}
		}
	}

	return ScenarioLicenseUsedDatabaseResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Licenses:  licenses,
	}
}

type ScenarioLicenseUsedHostResponse struct {
	ID        string                      `json:"id"`
	Name      string                      `json:"name"`
	CreatedAt time.Time                   `json:"createdAt"`
	Licenses  []LicenseUsedHostSimulation `json:"licenses"`
}

type LicenseUsedHostSimulation struct {
	Actual *model.LicenseUsedHost `json:"actual"`
	Got    *model.LicenseUsedHost `json:"got"`
}

func ToScenarioLicenseUsedHostResponse(m model.Scenario) ScenarioLicenseUsedHostResponse {
	licenses := make([]LicenseUsedHostSimulation, 0)

	for _, g := range m.LicenseUsed.LicenseHost.Got {
		got := g

		for _, a := range m.LicenseUsed.LicenseHost.Got {
			actual := a

			if a.Hostname == g.Hostname && a.LicenseTypeID == g.LicenseTypeID {
				licenses = append(licenses, LicenseUsedHostSimulation{
					Actual: &actual,
					Got:    &got,
				})
			}
		}
	}

	return ScenarioLicenseUsedHostResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Licenses:  licenses,
	}
}

type ScenarioLicenseUsedClusterResponse struct {
	ID        string                         `json:"id"`
	Name      string                         `json:"name"`
	CreatedAt time.Time                      `json:"createdAt"`
	Licenses  []LicenseUsedClusterSimulation `json:"licenses"`
}

type LicenseUsedClusterSimulation struct {
	Actual *model.LicenseUsedCluster `json:"actual"`
	Got    *model.LicenseUsedCluster `json:"got"`
}

func ToScenarioLicenseUsedClusterResponse(m model.Scenario) ScenarioLicenseUsedClusterResponse {
	licenses := make([]LicenseUsedClusterSimulation, 0)

	for _, g := range m.LicenseUsed.LicenseHypervisorCluster.Got {
		got := g

		for _, a := range m.LicenseUsed.LicenseHypervisorCluster.Actual {
			actual := a

			if a.Cluster == g.Cluster && a.LicenseTypeID == g.LicenseTypeID {
				licenses = append(licenses, LicenseUsedClusterSimulation{
					Actual: &actual,
					Got:    &got,
				})
			}
		}
	}

	return ScenarioLicenseUsedClusterResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Licenses:  licenses,
	}
}

type ScenarioLicenseUsedClusterVeritasResponse struct {
	ID        string                                `json:"id"`
	Name      string                                `json:"name"`
	CreatedAt time.Time                             `json:"createdAt"`
	Licenses  []LicenseUsedClsuterVeritasSimulation `json:"licenses"`
}

type LicenseUsedClsuterVeritasSimulation struct {
	Actual *model.LicenseUsedClusterVeritas `json:"actual"`
	Got    *model.LicenseUsedClusterVeritas `json:"got"`
}

func ToScenarioLicenseUsedClusterVeritasResponse(m model.Scenario) ScenarioLicenseUsedClusterVeritasResponse {
	licenses := make([]LicenseUsedClsuterVeritasSimulation, 0)

	for _, g := range m.LicenseUsed.LicenseClusterVeritas.Got {
		got := g

		for _, a := range m.LicenseUsed.LicenseClusterVeritas.Actual {
			actual := a

			if a.ID == g.ID && a.LicenseTypeID == g.LicenseTypeID {
				licenses = append(licenses, LicenseUsedClsuterVeritasSimulation{
					Actual: &actual,
					Got:    &got,
				})
			}
		}
	}

	return ScenarioLicenseUsedClusterVeritasResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Licenses:  licenses,
	}
}
