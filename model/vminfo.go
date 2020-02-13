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

// VMInfo holds info about the vm
type VMInfo struct {
	Name         string `bson:"Name"`
	ClusterName  string `bson:"ClusterName"`
	Hostname     string `bson:"Hostname"` //Hostname or IP address
	CappedCPU    bool   `bson:"CappedCPU"`
	PhysicalHost string `bson:"PhysicalHost"`
}

// VMInfoBsonValidatorRules contains mongodb validation rules for VMInfo
var VMInfoBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Name",
		"ClusterName",
		"Hostname",
		"CappedCPU",
		"PhysicalHost",
	}},
	{"properties", bson.D{
		{"Name", bson.D{
			{"bsonType", "string"},
		}},
		{"ClusterName", bson.D{
			{"bsonType", "string"},
		}},
		{"Hostname", bson.D{
			{"bsonType", "string"},
		}},
		{"CappedCPU", bson.D{
			{"bsonType", "bool"},
		}},
		{"PhysicalHost", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
