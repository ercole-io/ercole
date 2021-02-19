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

package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cluster struct {
	ID                          primitive.ObjectID        `json:"id" bson:"_id"`
	CPU                         int                       `json:"cpu" bson:"cpu"`
	CreatedAt                   time.Time                 `json:"createdAt" bson:"createdAt"`
	Environment                 string                    `json:"environment" bson:"environment"`
	FetchEndpoint               string                    `json:"fetchEndpoint" bson:"fetchEndpoint"`
	Hostname                    string                    `json:"hostname" bson:"hostname"`
	HostnameAgentVirtualization string                    `json:"hostnameAgentVirtualization" bson:"hostnameAgentVirtualization"`
	Location                    string                    `json:"location" bson:"location"`
	Name                        string                    `json:"name" bson:"name"`
	Sockets                     int                       `json:"sockets" bson:"sockets"`
	Type                        string                    `json:"type" bson:"type"`
	VirtualizationNodes         []string                  `json:"virtualizationNodes" bson:"virtualizationNodes"`
	VirtualizationNodesCount    int                       `json:"virtualizationNodesCount" bson:"virtualizationNodesCount"`
	VirtualizationNodesStats    []VirtualizationNodesStat `json:"virtualizationNodesStats" bson:"virtualizationNodesStats"`
	VMs                         []VM                      `json:"vms" bson:"vms"`
	VMsCount                    int                       `json:"vmsCount" bson:"vmsCount"`
	VMsErcoleAgentCount         int                       `json:"vmsErcoleAgentCount" bson:"vmsErcoleAgentCount"`
}

type VirtualizationNodesStat struct {
	TotalVMsCount                   int    `json:"totalVMsCount" bson:"totalVMsCount"`
	TotalVMsWithErcoleAgentCount    int    `json:"totalVMsWithErcoleAgentCount" bson:"totalVMsWithErcoleAgentCount"`
	TotalVMsWithoutErcoleAgentCount int    `json:"totalVMsWithoutErcoleAgentCount" bson:"totalVMsWithoutErcoleAgentCount"`
	VirtualizationNode              string `json:"virtualizationNode" bson:"virtualizationNode"`
}

type VM struct {
	CappedCPU          bool   `json:"cappedCPU" bson:"cappedCPU"`
	Hostname           string `json:"hostname" bson:"hostname"`
	Name               string `json:"name" bson:"name"`
	VirtualizationNode string `json:"virtualizationNode" bson:"virtualizationNode"`
}
