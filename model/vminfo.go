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
	Name        string
	ClusterName string
	Hostname    string //Hostname or IP address
	CappedCPU   int
}

var VMInfoBsonValidatorRules bson.D = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"name",
		"cluster_name",
		"hostname",
		"capped_cpu",
	}},
	{"properties", bson.D{
		{"name", bson.D{
			{"bsonType", "string"},
		}},
		{"cluster_name", bson.D{
			{"bsonType", "string"},
		}},
		{"hostname", bson.D{
			{"bsonType", "string"},
		}},
		{"capped_cpu", bson.D{
			{"bsonType", "int"},
		}},
	}},
}
