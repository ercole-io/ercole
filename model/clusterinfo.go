// Copyright (c) 2019 Sorint.lab S.p.A.
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

package model

import "go.mongodb.org/mongo-driver/bson"

//ClusterInfo hold informations about a cluster
type ClusterInfo struct {
	Name                        string
	Type                        string
	CPU                         int
	Sockets                     int
	VMs                         []VMInfo
	HostnameAgentVirtualization string `bson:"hostname_agent_virtualization"`
}

// ClusterInfoBsonValidatorRules contains mongodb validation rules for clusterInfo
var ClusterInfoBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"name",
		"type",
		"cpu",
		"sockets",
		"vms",
		"hostname_agent_virtualization",
	}},
	{"properties", bson.D{
		{"name", bson.D{
			{"bsonType", "string"},
		}},
		{"type", bson.D{
			{"bsonType", "string"},
		}},
		{"cpu", bson.D{
			{"bsonType", "int"},
		}},
		{"sockets", bson.D{
			{"bsonType", "int"},
		}},
		{"vms", bson.D{
			{"bsonType", "array"},
			{"items", VMInfoBsonValidatorRules},
		}},
		{"hostname_agent_virtualization", bson.D{
			{"bsonType", "string"},
			{"pattern", "^.+$"},
		}},
	}},
}
