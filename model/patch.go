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

// Patch holds information about a Oracle patch
type Patch struct {
	Database    string `bson:"Database"`
	Version     string `bson:"Version"`
	PatchID     string `bson:"PatchID"`
	Action      string `bson:"Action"`
	Description string `bson:"Description"`
	Date        string `bson:"Date"`
}

// PatchBsonValidatorRules contains mongodb validation rules for patch
var PatchBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Database",
		"Version",
		"PatchID",
		"Action",
		"Description",
		"Date",
	}},
	{"properties", bson.D{
		{"Database", bson.D{
			{"bsonType", "string"},
		}},
		{"Version", bson.D{
			{"bsonType", "string"},
		}},
		{"PatchID", bson.D{
			{"bsonType", "string"},
		}},
		{"Action", bson.D{
			{"bsonType", "string"},
		}},
		{"Description", bson.D{
			{"bsonType", "string"},
		}},
		{"Date", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
