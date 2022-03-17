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

package dto

import "time"

type OracleExadataResponse struct {
	Content  []OracleExadata `json:"content" bson:"content"`
	Metadata PagingMetadata  `json:"metadata" bson:"metadata"`
}

type OracleExadata struct {
	Id             string           `json:"_id" bson:"_id"`
	CreatedAt      time.Time        `json:"createdAt" bson:"createdAt"`
	DbServers      []DbServers      `json:"dbServers" bson:"dbServers"`
	Environment    string           `json:"environment" bson:"environment"`
	Hostname       string           `json:"hostname" bson:"hostname"`
	IbSwitches     []IbSwitches     `json:"ibSwitches" bson:"ibSwitches"`
	Location       string           `json:"location" bson:"location"`
	StorageServers []StorageServers `json:"storageServers" bson:"storageServers"`
}

type DbServers struct {
	Hostname           string `json:"hostname" bson:"hostname"`
	Memory             int    `json:"memory" bson:"memory"`
	Model              string `json:"model" bson:"model"`
	RunningCPUCount    int    `json:"runningCPUCount" bson:"runningCPUCount"`
	RunningPowerSupply int    `json:"runningPowerSupply" bson:"runningPowerSupply"`
	SwVersion          string `json:"swVersion" bson:"swVersion"`
	TempActual         int    `json:"tempActual" bson:"tempActual"`
	TotalCPUCount      int    `json:"totalCPUCount" bson:"totalCPUCount"`
	TotalPowerSupply   int    `json:"totalPowerSupply" bson:"totalPowerSupply"`
}

type IbSwitches struct {
	Hostname  string `json:"hostname" bson:"hostname"`
	Model     string `json:"model" bson:"model"`
	SwVersion string `json:"swVersion" bson:"swVersion"`
}

type StorageServers struct {
	Hostname           string `json:"hostname" bson:"hostname"`
	Memory             int    `json:"memory" bson:"memory"`
	Model              string `json:"model" bson:"model"`
	RunningCPUCount    int    `json:"runningCPUCount" bson:"runningCPUCount"`
	RunningPowerSupply int    `json:"runningPowerSupply" bson:"runningPowerSupply"`
	SwVersion          string `json:"swVersion" bson:"swVersion"`
	TempActual         int    `json:"tempActual" bson:"tempActual"`
	TotalCPUCount      int    `json:"totalCPUCount" bson:"totalCPUCount"`
	TotalPowerSupply   int    `json:"totalPowerSupply" bson:"totalPowerSupply"`
}
