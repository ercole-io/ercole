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
	CPU                         int                       `json,bson:"cpu"`
	CreatedAt                   time.Time                 `json,bson:"createdAt"`
	Environment                 string                    `json,bson:"environment"`
	FetchEndpoint               string                    `json,bson:"fetchEndPoint"`
	Hostname                    string                    `json,bson:"hostname"`
	HostnameAgentVirtualization string                    `json,bson:"hostnameAgentVirtualization"`
	Location                    string                    `json,bson:"location"`
	Name                        string                    `json,bson:"name"`
	Sockets                     int                       `json,bson:"sockets"`
	Type                        string                    `json,bson:"type"`
	VirtualizationNodes         []string                  `json,bson:"virtualizationNodes"`
	VirtualizationNodesCount    int                       `json,bson:"virtualizationNodesCount"`
	VirtualizationNodesStats    []VirtualizationNodesStat `json,bson:"virtualizationNodesStats"`
	VMs                         []VM                      `json,bson:"vms"`
	VMsCount                    int                       `json,bson:"vmsCount"`
	VMsErcoleAgentCount         int                       `json,bson:"vmsErcoleAgentCount"`
}

type VirtualizationNodesStat struct {
	TotalVMsCount                   int    `json,bson:"totalVMsCount "`
	TotalVMsWithErcoleAgentCount    int    `json,bson:"totalVMsWithErcoleAgentCount"`
	TotalVMsWithoutErcoleAgentCount int    `json,bson:"totalVMsWithoutErcoleAgentCount"`
	VirtualizationNode              string `json,bson:"virtualizationNode"`
}

type VM struct {
	CappedCPU          bool   `json,bson:"cappedCPU"`
	Hostname           string `json,bson:"hostname"`
	Name               string `json,bson:"name"`
	VirtualizationNode string `json,bson:"virtualizationNode"`
}
