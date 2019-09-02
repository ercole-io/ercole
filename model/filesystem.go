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

// Filesystem holds information about mounted filesystems
// and used space
type Filesystem struct {
	Filesystem string
	FsType     string
	Size       string
	Used       string
	Available  string
	UsedPerc   string
	MountedOn  string
}

var FilesystemBsonValidatorRules bson.D = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"filesystem",
		"fs_type",
		"size",
		"used",
		"available",
		"used_perc",
		"mounted_on",
	}},
	{"properties", bson.D{
		{"filesystem", bson.D{
			{"bsonType", "string"},
		}},
		{"fs_type", bson.D{
			{"bsonType", "string"},
		}},
		{"size", bson.D{
			{"bsonType", "string"},
		}},
		{"used", bson.D{
			{"bsonType", "string"},
		}},
		{"available", bson.D{
			{"bsonType", "string"},
		}},
		{"used_perc", bson.D{
			{"bsonType", "string"},
		}},
		{"mounted_on", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
