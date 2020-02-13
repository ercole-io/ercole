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

// PSU holds information about a PSU
type PSU struct {
	Date        string `bson:"Date"`
	Description string `bson:"Description"`
}

// PSUBsonValidatorRules contains mongodb validation rules for psu
var PSUBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Date",
		"Description",
	}},
	{"properties", bson.D{
		{"Date", bson.D{
			{"bsonType", "string"},
		}},
		{"Description", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
