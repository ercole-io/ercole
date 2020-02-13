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

// Addm contains info about addm
type Addm struct {
	Finding        string `bson:"Finding"`
	Recommendation string `bson:"Recommendation"`
	Action         string `bson:"Action"`
	Benefit        string `bson:"Benefit"`
}

// AddmBsonValidatorRules contains mongodb validation rules for addm
var AddmBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Finding",
		"Recommendation",
		"Action",
		"Benefit",
	}},
	{"properties", bson.D{
		{"Finding", bson.D{
			{"bsonType", "string"},
		}},
		{"Recommendation", bson.D{
			{"bsonType", "string"},
		}},
		{"Action", bson.D{
			{"bsonType", "string"},
		}},
		{"Benefit", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
